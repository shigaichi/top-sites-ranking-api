package usecase

import (
	"context"
	"errors"
	"net/url"
	"reflect"
	"testing"
	"time"

	"github.com/shigaichi/top-sites-ranking-api/internal/domain/model"
	"github.com/shigaichi/top-sites-ranking-api/internal/domain/repository"
	"github.com/shigaichi/tranco"
)

type MockTrancoApiRepository struct {
	Metadata tranco.ListMetadata
	Err      error
}

func (m *MockTrancoApiRepository) GetIdByDate(date time.Time) (tranco.ListMetadata, error) {
	return m.Metadata, m.Err
}

type MockTrancoListsRepository struct {
	IsExist     bool
	ExistsIdErr error
	SaveErr     error
}

func (m *MockTrancoListsRepository) ExistsId(ctx context.Context, id string) (bool, error) {
	if id != "X5Y7N" {
		return false, errors.New("unexpected parameters in existsId")
	}
	return m.IsExist, m.ExistsIdErr
}

func (m *MockTrancoListsRepository) Save(ctx context.Context, list model.TrancoList) error {
	if list.Id != "X5Y7N" {
		return errors.New("unexpected parameters in list save")
	}
	return m.SaveErr
}

type MockTrancoCsvRepository struct {
	SiteRankings []model.SiteRanking
	Err          error
}

func (m *MockTrancoCsvRepository) Get(url url.URL) ([]model.SiteRanking, error) {
	expected, _ := url.Parse("https://tranco-list.eu/download/X5Y7N/1000000")
	if *expected != url {
		return nil, errors.New("unexpected parameters in Get")
	}
	return m.SiteRankings, m.Err
}

type MockTransaction struct {
	Err error
}

func (m *MockTransaction) DoInTx(ctx context.Context, txFunc func(context.Context) (interface{}, error)) (interface{}, error) {
	if m.Err != nil {
		return nil, m.Err
	}
	return txFunc(ctx)
}

type MockTrancoDomainsRepository struct {
	Id               int
	GetIdByDomainErr error
	SaveErr          error
}

func (m *MockTrancoDomainsRepository) GetIdByDomain(ctx context.Context, domain string) (int, error) {
	return m.Id, m.GetIdByDomainErr
}

func (m *MockTrancoDomainsRepository) Save(ctx context.Context, domain string) (int, error) {
	return 10, m.SaveErr
}

type MockTrancoRankingsRepository struct {
	Err error
}

func (m *MockTrancoRankingsRepository) Save(ctx context.Context, ranking model.TrancoRanking) error {
	expected := model.TrancoRanking{DomainId: 10, ListId: "X5Y7N", Ranking: 1}
	if !reflect.DeepEqual(expected, ranking) {
		return errors.New("unexpected parameters in ranking save")
	}

	return m.Err
}

func TestStandardWriteInteractor_Write(t *testing.T) {
	tests := []struct {
		name          string
		inputDate     time.Time
		api           repository.TrancoApiRepository
		list          repository.TrancoListsRepository
		csv           repository.TrancoCsvRepository
		transaction   repository.Transaction
		domain        repository.TrancoDomainsRepository
		ranking       repository.TrancoRankingsRepository
		expectedError error
	}{
		{
			name:          "successful write",
			inputDate:     time.Now(),
			api:           &MockTrancoApiRepository{Metadata: tranco.ListMetadata{ListId: "X5Y7N", Download: "https://tranco-list.eu/download/X5Y7N/1000000", CreatedOn: time.Date(2023, 10, 17, 0, 0, 0, 0, time.UTC)}, Err: nil},
			list:          &MockTrancoListsRepository{IsExist: false, ExistsIdErr: nil},
			csv:           &MockTrancoCsvRepository{SiteRankings: []model.SiteRanking{{Domain: "example.com", Rank: 1}}, Err: nil},
			transaction:   &MockTransaction{},
			domain:        &MockTrancoDomainsRepository{Id: 0, GetIdByDomainErr: nil},
			ranking:       &MockTrancoRankingsRepository{},
			expectedError: nil,
		},
		{
			name:          "api error",
			inputDate:     time.Now(),
			api:           &MockTrancoApiRepository{Metadata: tranco.ListMetadata{ListId: "X5Y7N", Download: "https://tranco-list.eu/download/X5Y7N/1000000", CreatedOn: time.Date(2023, 10, 17, 0, 0, 0, 0, time.UTC)}, Err: errors.New("test")},
			list:          nil,
			csv:           nil,
			transaction:   nil,
			domain:        nil,
			ranking:       nil,
			expectedError: errors.New("failed to get tranco list id by date in writing standard tranco list error: test"),
		},
		{
			name:          "list was already saved",
			inputDate:     time.Now(),
			api:           &MockTrancoApiRepository{Metadata: tranco.ListMetadata{ListId: "X5Y7N", Download: "https://tranco-list.eu/download/X5Y7N/1000000", CreatedOn: time.Date(2023, 10, 17, 0, 0, 0, 0, time.UTC)}, Err: nil},
			list:          &MockTrancoListsRepository{IsExist: true, ExistsIdErr: nil, SaveErr: nil},
			csv:           nil,
			transaction:   nil,
			domain:        nil,
			ranking:       nil,
			expectedError: nil,
		},
		{
			name:          "csv download error",
			inputDate:     time.Now(),
			api:           &MockTrancoApiRepository{Metadata: tranco.ListMetadata{ListId: "X5Y7N", Download: "https://tranco-list.eu/download/X5Y7N/1000000", CreatedOn: time.Date(2023, 10, 17, 0, 0, 0, 0, time.UTC)}, Err: nil},
			list:          &MockTrancoListsRepository{IsExist: false, ExistsIdErr: nil, SaveErr: nil},
			csv:           &MockTrancoCsvRepository{SiteRankings: nil, Err: errors.New("test")},
			transaction:   nil,
			domain:        nil,
			ranking:       nil,
			expectedError: errors.New("failed to get csv in writing standard tranco list error: test"),
		},
		{
			name:          "transaction error",
			inputDate:     time.Now(),
			api:           &MockTrancoApiRepository{Metadata: tranco.ListMetadata{ListId: "X5Y7N", Download: "https://tranco-list.eu/download/X5Y7N/1000000", CreatedOn: time.Date(2023, 10, 17, 0, 0, 0, 0, time.UTC)}, Err: nil},
			list:          &MockTrancoListsRepository{IsExist: false, ExistsIdErr: nil, SaveErr: nil},
			csv:           &MockTrancoCsvRepository{SiteRankings: []model.SiteRanking{{Domain: "example.com", Rank: 1}}, Err: nil},
			transaction:   &MockTransaction{Err: errors.New("test")},
			domain:        nil,
			ranking:       nil,
			expectedError: errors.New("failed to save ranking data in writing standard tranco list and saving operation was rollbacked error: test"),
		},
		{
			name:          "list save error",
			inputDate:     time.Now(),
			api:           &MockTrancoApiRepository{Metadata: tranco.ListMetadata{ListId: "X5Y7N", Download: "https://tranco-list.eu/download/X5Y7N/1000000", CreatedOn: time.Date(2023, 10, 17, 0, 0, 0, 0, time.UTC)}, Err: nil},
			list:          &MockTrancoListsRepository{IsExist: false, ExistsIdErr: nil, SaveErr: errors.New("test")},
			csv:           &MockTrancoCsvRepository{SiteRankings: []model.SiteRanking{{Domain: "example.com", Rank: 1}}, Err: nil},
			transaction:   &MockTransaction{},
			domain:        nil,
			ranking:       nil,
			expectedError: errors.New("failed to save ranking data in writing standard tranco list and saving operation was rollbacked error: failed to save list id because it is already exist orr error in writing standard tranco list. error: test"),
		},
		{
			name:          "get domain id error",
			inputDate:     time.Now(),
			api:           &MockTrancoApiRepository{Metadata: tranco.ListMetadata{ListId: "X5Y7N", Download: "https://tranco-list.eu/download/X5Y7N/1000000", CreatedOn: time.Date(2023, 10, 17, 0, 0, 0, 0, time.UTC)}, Err: nil},
			list:          &MockTrancoListsRepository{IsExist: false, ExistsIdErr: nil, SaveErr: nil},
			csv:           &MockTrancoCsvRepository{SiteRankings: []model.SiteRanking{{Domain: "example.com", Rank: 1}}, Err: nil},
			transaction:   &MockTransaction{},
			domain:        &MockTrancoDomainsRepository{Id: 0, GetIdByDomainErr: errors.New("test")},
			ranking:       nil,
			expectedError: errors.New("failed to save ranking data in writing standard tranco list and saving operation was rollbacked error: failed to save list id in writing standard tranco list error: test"),
		},
		{
			name:          "domain exists",
			inputDate:     time.Now(),
			api:           &MockTrancoApiRepository{Metadata: tranco.ListMetadata{ListId: "X5Y7N", Download: "https://tranco-list.eu/download/X5Y7N/1000000", CreatedOn: time.Date(2023, 10, 17, 0, 0, 0, 0, time.UTC)}, Err: nil},
			list:          &MockTrancoListsRepository{IsExist: false, ExistsIdErr: nil, SaveErr: nil},
			csv:           &MockTrancoCsvRepository{SiteRankings: []model.SiteRanking{{Domain: "example.com", Rank: 1}}, Err: nil},
			transaction:   &MockTransaction{},
			domain:        &MockTrancoDomainsRepository{Id: 10, GetIdByDomainErr: nil, SaveErr: errors.New("unexpected invoke")},
			ranking:       &MockTrancoRankingsRepository{},
			expectedError: nil,
		},
		{
			name:          "domain save error",
			inputDate:     time.Now(),
			api:           &MockTrancoApiRepository{Metadata: tranco.ListMetadata{ListId: "X5Y7N", Download: "https://tranco-list.eu/download/X5Y7N/1000000", CreatedOn: time.Date(2023, 10, 17, 0, 0, 0, 0, time.UTC)}, Err: nil},
			list:          &MockTrancoListsRepository{IsExist: false, ExistsIdErr: nil, SaveErr: nil},
			csv:           &MockTrancoCsvRepository{SiteRankings: []model.SiteRanking{{Domain: "example.com", Rank: 1}}, Err: nil},
			transaction:   &MockTransaction{},
			domain:        &MockTrancoDomainsRepository{Id: 0, GetIdByDomainErr: nil, SaveErr: errors.New("test")},
			ranking:       nil,
			expectedError: errors.New("failed to save ranking data in writing standard tranco list and saving operation was rollbacked error: failed to save domain in writing standard tranco list error: test"),
		},
		{
			name:          "ranking save error",
			inputDate:     time.Now(),
			api:           &MockTrancoApiRepository{Metadata: tranco.ListMetadata{ListId: "X5Y7N", Download: "https://tranco-list.eu/download/X5Y7N/1000000", CreatedOn: time.Date(2023, 10, 17, 0, 0, 0, 0, time.UTC)}, Err: nil},
			list:          &MockTrancoListsRepository{IsExist: false, ExistsIdErr: nil, SaveErr: nil},
			csv:           &MockTrancoCsvRepository{SiteRankings: []model.SiteRanking{{Domain: "example.com", Rank: 1}}, Err: nil},
			transaction:   &MockTransaction{},
			domain:        &MockTrancoDomainsRepository{Id: 0, GetIdByDomainErr: nil, SaveErr: nil},
			ranking:       &MockTrancoRankingsRepository{Err: errors.New("test")},
			expectedError: errors.New("failed to save ranking data in writing standard tranco list and saving operation was rollbacked error: failed to save ranking in writing standard tranco list error: test"),
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			interactor := NewStandardWriteInteractor(tt.api, tt.list, tt.csv, tt.transaction, tt.domain, tt.ranking)

			err := interactor.Write(context.Background(), tt.inputDate)

			if tt.expectedError != nil {
				if err == nil {
					t.Errorf("expected error, got nil")
				} else if err.Error() != tt.expectedError.Error() {
					t.Errorf("expected error: %v, got: %v", tt.expectedError, err)
				}
			} else if err != nil {
				t.Errorf("unexpected error: %v", err)
			}
		})
	}
}
