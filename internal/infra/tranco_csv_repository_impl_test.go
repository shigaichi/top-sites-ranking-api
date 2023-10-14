package infra

import (
	"net/http"
	"net/http/httptest"
	"net/url"
	"testing"

	"github.com/google/go-cmp/cmp"
	"github.com/shigaichi/top-sites-ranking-api/internal/domain/model"
)

func TestTrancoCsvImpl_Get(t *testing.T) {
	type arg struct {
		body string
		code int
	}

	tests := []struct {
		name     string
		arg      arg
		expected []model.SiteRanking
		wantErr  bool
	}{
		{
			name:     "3 rows CSV",
			arg:      arg{body: "1,google.com\n2,amazonaws.com\n3,facebook.com\n", code: http.StatusOK},
			expected: []model.SiteRanking{{Rank: 1, Domain: "google.com"}, {Rank: 2, Domain: "amazonaws.com"}, {Rank: 3, Domain: "facebook.com"}},
			wantErr:  false,
		},
		{
			name:     "1 row CSV",
			arg:      arg{body: "1,google.com\n", code: http.StatusOK},
			expected: []model.SiteRanking{{Rank: 1, Domain: "google.com"}},
			wantErr:  false,
		},
		{
			name:     "HTTP error",
			arg:      arg{body: "1,google.com\n", code: http.StatusInternalServerError},
			expected: nil,
			wantErr:  true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			u, teardown := setup(tt.arg.body, tt.arg.code)
			defer teardown()

			p, err := url.Parse(u)
			if err != nil {
				t.Fatalf("mock url parse error: %v", err)
			}

			cli := TrancoCsvImpl{}
			result, err := cli.Get(*p)

			if err != nil {
				if !tt.wantErr {
					t.Fatalf("expected no error, got %v", err)
				}
				return
			}

			if diff := cmp.Diff(result, tt.expected); diff != "" {
				t.Errorf("response is mimatch:\n%s", diff)
			}
		})
	}
}

func setup(body string, code int) (string, func()) {
	mockStatusOK := httptest.NewServer(http.HandlerFunc(
		func(w http.ResponseWriter, r *http.Request) {
			w.WriteHeader(code)
			w.Write([]byte(body))
		},
	))

	teardown := func() {
		mockStatusOK.Close()
	}
	return mockStatusOK.URL, teardown
}
