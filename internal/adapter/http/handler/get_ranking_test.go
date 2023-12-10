package handler

import (
	"context"
	"encoding/json"
	"errors"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/google/go-cmp/cmp"

	"github.com/shigaichi/top-sites-ranking-api/internal/domain/model"
)

type UsecaseMock struct {
	Domain string
	Start  time.Time
	End    time.Time
	Result []model.DailyRank
	Err    error
}

func (m UsecaseMock) GetDailyRanking(ctx context.Context, domain string, start time.Time, end time.Time) ([]model.DailyRank, error) {
	if m.Domain == domain && m.Start == start && m.End == end {
		return m.Result, m.Err
	}
	return nil, errors.New("unexpected parameters")
}

func (m UsecaseMock) GetMonthlyRanking(ctx context.Context, domain string, start time.Time, end time.Time) ([]model.DailyRank, error) {
	if m.Domain == domain && m.Start == start && m.End == end {
		return m.Result, m.Err
	}
	return nil, errors.New("unexpected parameters")
}

func TestGetRankingImpl_GetDailyRanking(t *testing.T) {
	tests := []struct {
		name           string
		mockUsecase    UsecaseMock
		requestURL     string
		expectedStatus int
		expectedDomain string
		expectedRanks  int
		hasCacheHeader bool
	}{
		{
			name: "valid request",
			mockUsecase: UsecaseMock{
				Domain: "example.com",
				Start:  time.Date(2023, 1, 1, 0, 0, 0, 0, time.UTC),
				End:    time.Date(2023, 1, 3, 0, 0, 0, 0, time.UTC),
				Result: []model.DailyRank{
					{Rank: 1, Date: time.Date(2023, 1, 1, 0, 0, 0, 0, time.UTC)},
					{Rank: 10, Date: time.Date(2023, 1, 2, 0, 0, 0, 0, time.UTC)},
					{Rank: 20, Date: time.Date(2023, 1, 3, 0, 0, 0, 0, time.UTC)},
				},
				Err: nil,
			},
			requestURL:     "/api/v1/rankings/daily?domain=example.com&start_date=2023-01-01&end_date=2023-01-03",
			expectedStatus: http.StatusOK,
			expectedDomain: "example.com",
			expectedRanks:  3,
			hasCacheHeader: true,
		},
		{
			name: "valid request. But does not have all ranks",
			mockUsecase: UsecaseMock{
				Domain: "example.com",
				Start:  time.Date(2023, 1, 1, 0, 0, 0, 0, time.UTC),
				End:    time.Date(2023, 1, 3, 0, 0, 0, 0, time.UTC),
				Result: []model.DailyRank{
					{Rank: 1, Date: time.Date(2023, 1, 1, 0, 0, 0, 0, time.UTC)},
					{Rank: 10, Date: time.Date(2023, 1, 2, 0, 0, 0, 0, time.UTC)},
				},
				Err: nil,
			},
			requestURL:     "/api/v1/rankings/daily?domain=example.com&start_date=2023-01-01&end_date=2023-01-03",
			expectedStatus: http.StatusOK,
			expectedDomain: "example.com",
			expectedRanks:  2,
			hasCacheHeader: false,
		},
		{
			name:           "empty start date request",
			mockUsecase:    UsecaseMock{},
			requestURL:     "/api/v1/rankings/daily?domain=example.com&end_date=2023-01-31",
			expectedStatus: http.StatusBadRequest,
			expectedDomain: "",
			expectedRanks:  0,
		},
		{
			name:           "empty end date request",
			mockUsecase:    UsecaseMock{},
			requestURL:     "/api/v1/rankings/daily?domain=example.com&start_date=2023-12-01&end_date=",
			expectedStatus: http.StatusBadRequest,
			expectedDomain: "",
			expectedRanks:  0,
		},
		{
			name:           "invalid start date request",
			mockUsecase:    UsecaseMock{},
			requestURL:     "/api/v1/rankings/daily?domain=example.com&start_date=2023-13-01&end_date=2023-01-31",
			expectedStatus: http.StatusBadRequest,
			expectedDomain: "",
			expectedRanks:  0,
		},
		{
			name:           "invalid end date request",
			mockUsecase:    UsecaseMock{},
			requestURL:     "/api/v1/rankings/daily?domain=example.com&start_date=2023-12-01&end_date=XXX",
			expectedStatus: http.StatusBadRequest,
			expectedDomain: "",
			expectedRanks:  0,
		},
		{
			name:           "start date after end date request",
			mockUsecase:    UsecaseMock{},
			requestURL:     "/api/v1/rankings/daily?domain=example.com&start_date=2023-12-02&end_date=2023-12-01",
			expectedStatus: http.StatusBadRequest,
			expectedDomain: "",
			expectedRanks:  0,
		},
		{
			name: "get error while fetching ranking data",
			mockUsecase: UsecaseMock{
				Domain: "example.com",
				Start:  time.Date(2023, 1, 1, 0, 0, 0, 0, time.UTC),
				End:    time.Date(2023, 1, 31, 0, 0, 0, 0, time.UTC),
				Result: []model.DailyRank{
					{Rank: 1, Date: time.Date(2023, 1, 15, 0, 0, 0, 0, time.UTC)},
				},
				Err: errors.New("test"),
			},
			requestURL:     "/api/v1/rankings/daily?domain=example.com&start_date=2023-01-01&end_date=2023-01-31",
			expectedStatus: http.StatusInternalServerError,
			expectedDomain: "",
			expectedRanks:  0,
		},
		{
			name: "get no data about requested domain",
			mockUsecase: UsecaseMock{
				Domain: "example.com",
				Start:  time.Date(2023, 1, 1, 0, 0, 0, 0, time.UTC),
				End:    time.Date(2023, 1, 31, 0, 0, 0, 0, time.UTC),
				Result: []model.DailyRank{},
				Err:    nil,
			},
			requestURL:     "/api/v1/rankings/daily?domain=example.com&start_date=2023-01-01&end_date=2023-01-31",
			expectedStatus: http.StatusNotFound,
			expectedDomain: "",
			expectedRanks:  0,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			handler := NewGetRankingImpl(tt.mockUsecase)

			req, _ := http.NewRequest("GET", tt.requestURL, nil)
			rr := httptest.NewRecorder()
			handlerFunc := http.HandlerFunc(handler.GetDailyRanking)
			handlerFunc.ServeHTTP(rr, req)

			if status := rr.Code; status != tt.expectedStatus {
				t.Errorf("handler returned wrong status code: got %v want %v", status, tt.expectedStatus)
			}

			if tt.expectedStatus == http.StatusOK {
				var response map[string]interface{}
				err := json.Unmarshal(rr.Body.Bytes(), &response)
				if err != nil {
					t.Fatalf("json.Unmarshal failed: %v", err)
				}

				if domain := response["domain"].(string); domain != tt.expectedDomain {
					t.Errorf("Expected domain to be '%s', got %s", tt.expectedDomain, domain)
				}

				ranks, ok := response["ranks"].([]interface{})
				if !ok || len(ranks) != tt.expectedRanks {
					t.Errorf("Expected %d ranks, got %v", tt.expectedRanks, ranks)
				}
			}

			c := rr.Header().Get("Cache-Control")
			if tt.hasCacheHeader && (len(c) == 0) {
				t.Errorf("header has been expected.")
			} else if !tt.hasCacheHeader && (len(c) > 0) {
				t.Errorf("header has not been expected. but got %s", c)
			}
		})
	}
}

func TestGetRankingImpl_GetMonthlyRanking(t *testing.T) {
	tests := []struct {
		name           string
		mockUsecase    UsecaseMock
		requestURL     string
		expectedStatus int
		expectedBody   string
		hasCacheHeader bool
	}{
		{
			name: "valid request",
			mockUsecase: UsecaseMock{
				Domain: "example.com",
				Start:  getLastDayOfMonth(time.Date(2023, 1, 1, 0, 0, 0, 0, time.UTC)),
				End:    getLastDayOfMonth(time.Date(2023, 3, 1, 0, 0, 0, 0, time.UTC)),
				Result: []model.DailyRank{
					{Rank: 1, Date: getLastDayOfMonth(time.Date(2023, 1, 1, 0, 0, 0, 0, time.UTC))},
					{Rank: 2, Date: getLastDayOfMonth(time.Date(2023, 2, 1, 0, 0, 0, 0, time.UTC))},
					{Rank: 3, Date: getLastDayOfMonth(time.Date(2023, 3, 1, 0, 0, 0, 0, time.UTC))},
				},
			},
			requestURL:     "/api/v1/rankings/monthly?domain=example.com&start_month=2023-01&end_month=2023-03",
			expectedStatus: http.StatusOK,
			expectedBody:   `{"ranks":[{"rank":1,"date":"2023-01-31"},{"rank":2,"date":"2023-02-28"},{"rank":3,"date":"2023-03-31"}],"domain":"example.com"}`,
			hasCacheHeader: true,
		},
		{
			name: "valid request. But does not have all ranks",
			mockUsecase: UsecaseMock{
				Domain: "example.com",
				Start:  getLastDayOfMonth(time.Date(2023, 1, 1, 0, 0, 0, 0, time.UTC)),
				End:    getLastDayOfMonth(time.Date(2023, 3, 1, 0, 0, 0, 0, time.UTC)),
				Result: []model.DailyRank{
					{Rank: 1, Date: getLastDayOfMonth(time.Date(2023, 1, 1, 0, 0, 0, 0, time.UTC))},
					{Rank: 2, Date: getLastDayOfMonth(time.Date(2023, 2, 1, 0, 0, 0, 0, time.UTC))},
				},
			},
			requestURL:     "/api/v1/rankings/monthly?domain=example.com&start_month=2023-01&end_month=2023-03",
			expectedStatus: http.StatusOK,
			expectedBody:   `{"ranks":[{"rank":1,"date":"2023-01-31"},{"rank":2,"date":"2023-02-28"}],"domain":"example.com"}`,
			hasCacheHeader: false,
		},
		{
			name:           "empty start month request",
			mockUsecase:    UsecaseMock{},
			requestURL:     "/api/v1/rankings/monthly?domain=example.com&end_month=2023-12",
			expectedStatus: http.StatusBadRequest,
			expectedBody:   "Bad Request: Missing or invalid query parameters",
		},
		{
			name:           "empty end month request",
			mockUsecase:    UsecaseMock{},
			requestURL:     "/api/v1/rankings/monthly?domain=example.com&start_month=2023-01&end_month=",
			expectedStatus: http.StatusBadRequest,
			expectedBody:   "Bad Request: Missing or invalid query parameters",
		},
		{
			name:           "invalid start month request",
			mockUsecase:    UsecaseMock{},
			requestURL:     "/api/v1/rankings/monthly?domain=example.com&start_month=2023-01-01&end_month=2023-12",
			expectedStatus: http.StatusBadRequest,
			expectedBody:   "Invalid start_month format",
		},
		{
			name:           "invalid end month request",
			mockUsecase:    UsecaseMock{},
			requestURL:     "/api/v1/rankings/monthly?domain=example.com&start_month=2023-01&end_month=XXX",
			expectedStatus: http.StatusBadRequest,
			expectedBody:   "Invalid end_month format",
		},
		{
			name:           "start date after end date request",
			mockUsecase:    UsecaseMock{},
			requestURL:     "/api/v1/rankings/monthly?domain=example.com&start_month=2023-01&end_month=2022-12",
			expectedStatus: http.StatusBadRequest,
			expectedBody:   "start_month should be before or equal to end_month",
		},
		{
			name: "get error while fetching ranking data",
			mockUsecase: UsecaseMock{
				Domain: "example.com",
				Start:  getLastDayOfMonth(time.Date(2023, 1, 1, 0, 0, 0, 0, time.UTC)),
				End:    getLastDayOfMonth(time.Date(2023, 12, 1, 0, 0, 0, 0, time.UTC)),
				Result: []model.DailyRank{},
				Err:    errors.New("test"),
			},
			requestURL:     "/api/v1/rankings/monthly?domain=example.com&start_month=2023-01&end_month=2023-12",
			expectedStatus: http.StatusInternalServerError,
			expectedBody:   http.StatusText(http.StatusInternalServerError),
		},
		{
			name: "get no data about requested domain",
			mockUsecase: UsecaseMock{
				Domain: "example.com",
				Start:  getLastDayOfMonth(time.Date(2023, 1, 1, 0, 0, 0, 0, time.UTC)),
				End:    getLastDayOfMonth(time.Date(2023, 12, 1, 0, 0, 0, 0, time.UTC)),
				Result: []model.DailyRank{},
			},
			requestURL:     "/api/v1/rankings/monthly?domain=example.com&start_month=2023-01&end_month=2023-12",
			expectedStatus: http.StatusNotFound,
			expectedBody:   "No ranks found for the given domain and date range",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			req, _ := http.NewRequest(http.MethodGet, tt.requestURL, nil)

			rec := httptest.NewRecorder()
			handler := NewGetRankingImpl(tt.mockUsecase)
			handlerFunc := http.HandlerFunc(handler.GetMonthlyRanking)
			handlerFunc.ServeHTTP(rec, req)

			if rec.Code != tt.expectedStatus {
				t.Errorf("expected status %d, got %d", tt.expectedStatus, rec.Code)
			}

			if tt.expectedStatus == http.StatusOK {
				var body map[string]interface{}
				if err := json.Unmarshal(rec.Body.Bytes(), &body); err != nil {
					t.Fatal(err)
				}

				expectedBodyMap := map[string]interface{}{}
				if err := json.Unmarshal([]byte(tt.expectedBody), &expectedBodyMap); err != nil {
					t.Fatal(err)
				}

				if diff := cmp.Diff(expectedBodyMap, body); diff != "" {
					t.Errorf("unexpected response (-want +got):\n%s", diff)
				}
			} else {
				if body := rec.Body.String(); body != tt.expectedBody+"\n" {
					t.Errorf("expected body %q, got %q", tt.expectedBody, body)
				}
			}

			c := rec.Header().Get("Cache-Control")
			if tt.hasCacheHeader && (len(c) == 0) {
				t.Errorf("header has been expected.")
			} else if !tt.hasCacheHeader && (len(c) > 0) {
				t.Errorf("header has not been expected. but got %s", c)
			}
		})
	}
}
