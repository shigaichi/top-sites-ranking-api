package usecase

import (
	"context"
	"errors"
	"fmt"
	"net/url"
	"testing"
	"time"

	"github.com/google/go-cmp/cmp"

	"github.com/shigaichi/top-sites-ranking-api/internal/domain/model"
	"github.com/shigaichi/top-sites-ranking-api/internal/domain/repository"
	"github.com/shigaichi/tranco"
)

type MockTrancoAPIRepository struct {
	Metadata tranco.ListMetadata
	Err      error
}

func (m *MockTrancoAPIRepository) GetIDByDate(date time.Time) (tranco.ListMetadata, error) {
	return m.Metadata, m.Err
}

type MockTrancoListsRepository struct {
	IsExist     bool
	ExistsIDErr error
	SaveErr     error
}

func (m *MockTrancoListsRepository) ExistsID(ctx context.Context, id string) (bool, error) {
	if id != "X5Y7N" {
		return false, errors.New("unexpected parameters in existsID")
	}
	return m.IsExist, m.ExistsIDErr
}

func (m *MockTrancoListsRepository) Save(ctx context.Context, list model.TrancoList) error {
	if list.ID != "X5Y7N" {
		return errors.New("unexpected parameters in list save")
	}
	return m.SaveErr
}

func (m *MockTrancoListsRepository) FindByCreatedOnLessThan(ctx context.Context, date time.Time) ([]model.TrancoList, error) {
	return nil, errors.New("not implemented")
}

func (m *MockTrancoListsRepository) DeleteByID(ctx context.Context, id string) error {
	return errors.New("not implemented")
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
	ID               int
	GetIDByDomainErr error
	SaveErr          error
}

func (m *MockTrancoDomainsRepository) GetIDByDomain(ctx context.Context, domain string) (int, error) {
	return m.ID, m.GetIDByDomainErr
}

func (m *MockTrancoDomainsRepository) Save(ctx context.Context, domain string) (int, error) {
	return 10, m.SaveErr
}

type MockTrancoRankingsRepository struct {
	Err              error
	ExpectedRankings []model.TrancoRanking
}

func (m *MockTrancoRankingsRepository) Save(ctx context.Context, ranking model.TrancoRanking) error {
	return errors.New("unexpected invoke")
}

func (m *MockTrancoRankingsRepository) BulkSave(ctx context.Context, rankings []model.TrancoRanking) error {
	if m.Err != nil {
		return m.Err
	}

	if diff := cmp.Diff(rankings, m.ExpectedRankings); diff != "" {
		return fmt.Errorf("response is mismatch:\n%s", diff)
	}

	return nil
}

func (m *MockTrancoRankingsRepository) DeleteByListID(ctx context.Context, listID string) error {
	return errors.New("not implemented")
}

func TestStandardWriteInteractor_Write(t *testing.T) {
	tests := []struct {
		name          string
		inputDate     time.Time
		api           repository.TrancoAPIRepository
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
			api:           &MockTrancoAPIRepository{Metadata: tranco.ListMetadata{ListID: "X5Y7N", Download: "https://tranco-list.eu/download/X5Y7N/1000000", CreatedOn: time.Date(2023, 10, 17, 0, 0, 0, 0, time.UTC)}, Err: nil},
			list:          &MockTrancoListsRepository{IsExist: false, ExistsIDErr: nil},
			csv:           &MockTrancoCsvRepository{SiteRankings: []model.SiteRanking{{Domain: "example.com", Rank: 1}}, Err: nil},
			transaction:   &MockTransaction{},
			domain:        &MockTrancoDomainsRepository{ID: 0, GetIDByDomainErr: nil},
			ranking:       &MockTrancoRankingsRepository{ExpectedRankings: []model.TrancoRanking{{DomainID: 10, ListID: "X5Y7N", Ranking: 1}}},
			expectedError: nil,
		},
		{
			name:          "api error",
			inputDate:     time.Now(),
			api:           &MockTrancoAPIRepository{Metadata: tranco.ListMetadata{ListID: "X5Y7N", Download: "https://tranco-list.eu/download/X5Y7N/1000000", CreatedOn: time.Date(2023, 10, 17, 0, 0, 0, 0, time.UTC)}, Err: errors.New("test")},
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
			api:           &MockTrancoAPIRepository{Metadata: tranco.ListMetadata{ListID: "X5Y7N", Download: "https://tranco-list.eu/download/X5Y7N/1000000", CreatedOn: time.Date(2023, 10, 17, 0, 0, 0, 0, time.UTC)}, Err: nil},
			list:          &MockTrancoListsRepository{IsExist: true, ExistsIDErr: nil, SaveErr: nil},
			csv:           nil,
			transaction:   nil,
			domain:        nil,
			ranking:       nil,
			expectedError: nil,
		},
		{
			name:          "csv download error",
			inputDate:     time.Now(),
			api:           &MockTrancoAPIRepository{Metadata: tranco.ListMetadata{ListID: "X5Y7N", Download: "https://tranco-list.eu/download/X5Y7N/1000000", CreatedOn: time.Date(2023, 10, 17, 0, 0, 0, 0, time.UTC)}, Err: nil},
			list:          &MockTrancoListsRepository{IsExist: false, ExistsIDErr: nil, SaveErr: nil},
			csv:           &MockTrancoCsvRepository{SiteRankings: nil, Err: errors.New("test")},
			transaction:   nil,
			domain:        nil,
			ranking:       nil,
			expectedError: errors.New("failed to get csv in writing standard tranco list error: test"),
		},
		{
			name:          "transaction error",
			inputDate:     time.Now(),
			api:           &MockTrancoAPIRepository{Metadata: tranco.ListMetadata{ListID: "X5Y7N", Download: "https://tranco-list.eu/download/X5Y7N/1000000", CreatedOn: time.Date(2023, 10, 17, 0, 0, 0, 0, time.UTC)}, Err: nil},
			list:          &MockTrancoListsRepository{IsExist: false, ExistsIDErr: nil, SaveErr: nil},
			csv:           &MockTrancoCsvRepository{SiteRankings: []model.SiteRanking{{Domain: "example.com", Rank: 1}}, Err: nil},
			transaction:   &MockTransaction{Err: errors.New("test")},
			domain:        nil,
			ranking:       nil,
			expectedError: errors.New("failed to save ranking data in writing standard tranco list and saving operation was rollbacked error: test"),
		},
		{
			name:          "list save error",
			inputDate:     time.Now(),
			api:           &MockTrancoAPIRepository{Metadata: tranco.ListMetadata{ListID: "X5Y7N", Download: "https://tranco-list.eu/download/X5Y7N/1000000", CreatedOn: time.Date(2023, 10, 17, 0, 0, 0, 0, time.UTC)}, Err: nil},
			list:          &MockTrancoListsRepository{IsExist: false, ExistsIDErr: nil, SaveErr: errors.New("test")},
			csv:           &MockTrancoCsvRepository{SiteRankings: []model.SiteRanking{{Domain: "example.com", Rank: 1}}, Err: nil},
			transaction:   &MockTransaction{},
			domain:        nil,
			ranking:       nil,
			expectedError: errors.New("failed to save ranking data in writing standard tranco list and saving operation was rollbacked error: failed to save tranco list with id X5Y7N error: test"),
		},
		{
			name:          "get domain id error",
			inputDate:     time.Now(),
			api:           &MockTrancoAPIRepository{Metadata: tranco.ListMetadata{ListID: "X5Y7N", Download: "https://tranco-list.eu/download/X5Y7N/1000000", CreatedOn: time.Date(2023, 10, 17, 0, 0, 0, 0, time.UTC)}, Err: nil},
			list:          &MockTrancoListsRepository{IsExist: false, ExistsIDErr: nil, SaveErr: nil},
			csv:           &MockTrancoCsvRepository{SiteRankings: []model.SiteRanking{{Domain: "example.com", Rank: 1}}, Err: nil},
			transaction:   &MockTransaction{},
			domain:        &MockTrancoDomainsRepository{ID: 0, GetIDByDomainErr: errors.New("test")},
			ranking:       nil,
			expectedError: errors.New("failed to save ranking data in writing standard tranco list and saving operation was rollbacked error: failed to save list id in writing standard tranco list error: test"),
		},
		{
			name:          "domain exists",
			inputDate:     time.Now(),
			api:           &MockTrancoAPIRepository{Metadata: tranco.ListMetadata{ListID: "X5Y7N", Download: "https://tranco-list.eu/download/X5Y7N/1000000", CreatedOn: time.Date(2023, 10, 17, 0, 0, 0, 0, time.UTC)}, Err: nil},
			list:          &MockTrancoListsRepository{IsExist: false, ExistsIDErr: nil, SaveErr: nil},
			csv:           &MockTrancoCsvRepository{SiteRankings: []model.SiteRanking{{Domain: "example.com", Rank: 1}}, Err: nil},
			transaction:   &MockTransaction{},
			domain:        &MockTrancoDomainsRepository{ID: 10, GetIDByDomainErr: nil},
			ranking:       &MockTrancoRankingsRepository{ExpectedRankings: []model.TrancoRanking{{DomainID: 10, ListID: "X5Y7N", Ranking: 1}}},
			expectedError: nil,
		},
		{
			name:          "domain save error",
			inputDate:     time.Now(),
			api:           &MockTrancoAPIRepository{Metadata: tranco.ListMetadata{ListID: "X5Y7N", Download: "https://tranco-list.eu/download/X5Y7N/1000000", CreatedOn: time.Date(2023, 10, 17, 0, 0, 0, 0, time.UTC)}, Err: nil},
			list:          &MockTrancoListsRepository{IsExist: false, ExistsIDErr: nil, SaveErr: nil},
			csv:           &MockTrancoCsvRepository{SiteRankings: []model.SiteRanking{{Domain: "example.com", Rank: 1}}, Err: nil},
			transaction:   &MockTransaction{},
			domain:        &MockTrancoDomainsRepository{ID: 0, GetIDByDomainErr: nil, SaveErr: errors.New("test")},
			ranking:       nil,
			expectedError: errors.New("failed to save ranking data in writing standard tranco list and saving operation was rollbacked error: failed to save domain in writing standard tranco list error: test"),
		},
		{
			name:          "ranking save error",
			inputDate:     time.Now(),
			api:           &MockTrancoAPIRepository{Metadata: tranco.ListMetadata{ListID: "X5Y7N", Download: "https://tranco-list.eu/download/X5Y7N/1000000", CreatedOn: time.Date(2023, 10, 17, 0, 0, 0, 0, time.UTC)}, Err: nil},
			list:          &MockTrancoListsRepository{IsExist: false, ExistsIDErr: nil, SaveErr: nil},
			csv:           &MockTrancoCsvRepository{SiteRankings: []model.SiteRanking{{Domain: "example.com", Rank: 1}}, Err: nil},
			transaction:   &MockTransaction{},
			domain:        &MockTrancoDomainsRepository{ID: 0, GetIDByDomainErr: nil, SaveErr: nil},
			ranking:       &MockTrancoRankingsRepository{Err: errors.New("test")},
			expectedError: errors.New("failed to save ranking data in writing standard tranco list and saving operation was rollbacked error: failed to bulk save 1 rankings in writing standard tranco list error: test"),
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
