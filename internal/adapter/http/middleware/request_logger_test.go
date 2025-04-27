package middleware

import (
	"bytes"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
	"time"

	log "github.com/sirupsen/logrus"
)

func TestRequestLoggerMiddleware(t *testing.T) {
	var buf bytes.Buffer
	log.SetOutput(&buf)
	defer log.SetOutput(nil)

	log.SetLevel(log.InfoLevel)

	log.SetFormatter(&log.JSONFormatter{})

	mw := RequestLoggerMiddleware([]string{"/status"})

	testHandler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
		_, err := w.Write([]byte("ok"))
		if err != nil {
			t.Fatalf("failed to write response: %v", err)
		}
	})

	handler := mw(testHandler)

	req := httptest.NewRequest(http.MethodGet, "/api/v1/rankings", nil)
	rec := httptest.NewRecorder()

	handler.ServeHTTP(rec, req)

	logOutput := buf.String()

	if !strings.Contains(logOutput, `"method":"GET"`) {
		t.Errorf("expected log to contain method, but got: %s", logOutput)
	}
	if !strings.Contains(logOutput, `"/api/v1/rankings"`) {
		t.Errorf("expected log to contain url, but got: %s", logOutput)
	}
	if !strings.Contains(logOutput, `"status":200`) {
		t.Errorf("expected log to contain status 200, but got: %s", logOutput)
	}
	if !strings.Contains(logOutput, `"handled request"`) {
		t.Errorf("expected log message 'handled request', but got: %s", logOutput)
	}
}

func TestRequestLoggerMiddleware_SkipPath(t *testing.T) {
	var buf bytes.Buffer
	log.SetOutput(&buf)
	defer log.SetOutput(nil)

	log.SetLevel(log.InfoLevel)

	log.SetFormatter(&log.JSONFormatter{})

	mw := RequestLoggerMiddleware([]string{"/status"})

	testHandler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
		_, err := w.Write([]byte("ok"))
		if err != nil {
			t.Fatalf("failed to write response: %v", err)
		}
	})

	handler := mw(testHandler)

	req := httptest.NewRequest(http.MethodGet, "/status", nil)
	rec := httptest.NewRecorder()

	handler.ServeHTTP(rec, req)

	time.Sleep(10 * time.Millisecond)

	logOutput := buf.String()

	if logOutput != "" {
		t.Errorf("expected no log output for skipped path, but got: %s", logOutput)
	}
}
