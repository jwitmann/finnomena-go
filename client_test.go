package finnomena

import (
	"encoding/json"
	"errors"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
	"time"

	"github.com/jwitmann/finnomena-models"
)

func TestNewClient(t *testing.T) {
	client := NewClient()
	if client == nil {
		t.Fatal("NewClient() returned nil")
	}
	if client.baseURL != BaseURL {
		t.Errorf("Expected baseURL %s, got %s", BaseURL, client.baseURL)
	}
	if client.maxRetries != DefaultMaxRetries {
		t.Errorf("Expected maxRetries %d, got %d", DefaultMaxRetries, client.maxRetries)
	}
	if client.retryDelay != DefaultRetryDelay {
		t.Errorf("Expected retryDelay %v, got %v", DefaultRetryDelay, client.retryDelay)
	}
}

func TestNewClientWithTimeout(t *testing.T) {
	timeout := 10 * time.Second
	client := NewClientWithTimeout(timeout)
	if client == nil {
		t.Fatal("NewClientWithTimeout() returned nil")
	}
	if client.httpClient.Timeout != timeout {
		t.Errorf("Expected timeout %v, got %v", timeout, client.httpClient.Timeout)
	}
}

func TestNewClientWithRetry(t *testing.T) {
	maxRetries := 5
	retryDelay := 2 * time.Second
	client := NewClientWithRetry(maxRetries, retryDelay)
	if client == nil {
		t.Fatal("NewClientWithRetry() returned nil")
	}
	if client.maxRetries != maxRetries {
		t.Errorf("Expected maxRetries %d, got %d", maxRetries, client.maxRetries)
	}
	if client.retryDelay != retryDelay {
		t.Errorf("Expected retryDelay %v, got %v", retryDelay, client.retryDelay)
	}
}

func TestSetRetryConfig(t *testing.T) {
	client := NewClient()
	maxRetries := 5
	retryDelay := 2 * time.Second

	client.SetRetryConfig(maxRetries, retryDelay)

	if client.maxRetries != maxRetries {
		t.Errorf("Expected maxRetries %d, got %d", maxRetries, client.maxRetries)
	}
	if client.retryDelay != retryDelay {
		t.Errorf("Expected retryDelay %v, got %v", retryDelay, client.retryDelay)
	}
}

func TestIsRetryableError(t *testing.T) {
	tests := []struct {
		name       string
		err        error
		statusCode int
		expected   bool
	}{
		{
			name:       "network error",
			err:        errors.New("network error"),
			statusCode: 0,
			expected:   true,
		},
		{
			name:       "400 bad request",
			err:        nil,
			statusCode: 400,
			expected:   false,
		},
		{
			name:       "404 not found",
			err:        nil,
			statusCode: 404,
			expected:   false,
		},
		{
			name:       "429 too many requests",
			err:        nil,
			statusCode: 429,
			expected:   false,
		},
		{
			name:       "500 internal server error",
			err:        nil,
			statusCode: 500,
			expected:   true,
		},
		{
			name:       "503 service unavailable",
			err:        nil,
			statusCode: 503,
			expected:   true,
		},
		{
			name:       "200 ok",
			err:        nil,
			statusCode: 200,
			expected:   false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := isRetryableError(tt.err, tt.statusCode)
			if result != tt.expected {
				t.Errorf("isRetryableError() = %v, want %v", result, tt.expected)
			}
		})
	}
}

func TestDoRequestWithRetry_Success(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
		json.NewEncoder(w).Encode(map[string]string{"status": "ok"})
	}))
	defer server.Close()

	client := NewClient()
	client.baseURL = server.URL

	body, err := client.doRequest("/test")
	if err != nil {
		t.Fatalf("doRequest() error = %v", err)
	}
	if len(body) == 0 {
		t.Error("Expected non-empty body")
	}
}

func TestDoRequestWithRetry_Failure(t *testing.T) {
	requestCount := 0
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		requestCount++
		w.WriteHeader(http.StatusInternalServerError)
	}))
	defer server.Close()

	client := NewClientWithRetry(2, 1*time.Millisecond)
	client.baseURL = server.URL

	_, err := client.doRequest("/test")
	if err == nil {
		t.Fatal("Expected error, got nil")
	}

	if requestCount != 2 {
		t.Errorf("Expected 2 requests (maxRetries), got %d", requestCount)
	}
}

func TestDoRequestWithRetry_NoRetryOn4xx(t *testing.T) {
	requestCount := 0
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		requestCount++
		w.WriteHeader(http.StatusBadRequest)
	}))
	defer server.Close()

	client := NewClientWithRetry(3, 1*time.Millisecond)
	client.baseURL = server.URL

	_, err := client.doRequest("/test")
	if err == nil {
		t.Fatal("Expected error, got nil")
	}

	if requestCount != 1 {
		t.Errorf("Expected 1 request (no retry on 4xx), got %d", requestCount)
	}
}

func TestGetFundsList(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path != "/funds" {
			t.Errorf("Expected path /funds, got %s", r.URL.Path)
		}

		response := models.FundsResponse{
			Status: true,
			Data: []models.Fund{
				{
					FundID:    "F000001",
					ShortCode: "TEST-A",
					NameTH:    "กองทุนทดสอบ A",
				},
				{
					FundID:    "F000002",
					ShortCode: "TEST-B",
					NameTH:    "กองทุนทดสอบ B",
				},
			},
		}
		w.WriteHeader(http.StatusOK)
		json.NewEncoder(w).Encode(response)
	}))
	defer server.Close()

	client := NewClient()
	client.baseURL = server.URL

	funds, err := client.GetFundsList()
	if err != nil {
		t.Fatalf("GetFundsList() error = %v", err)
	}

	if len(funds) != 2 {
		t.Errorf("Expected 2 funds, got %d", len(funds))
	}

	if funds[0].ShortCode != "TEST-A" {
		t.Errorf("Expected first fund TEST-A, got %s", funds[0].ShortCode)
	}
}

func TestGetFundsList_Error(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusInternalServerError)
	}))
	defer server.Close()

	client := NewClientWithRetry(1, 1*time.Millisecond)
	client.baseURL = server.URL

	_, err := client.GetFundsList()
	if err == nil {
		t.Fatal("Expected error, got nil")
	}
}

func TestGetHistoricalPrices(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if !strings.HasPrefix(r.URL.Path, "/tv/history") {
			t.Errorf("Expected path /tv/history, got %s", r.URL.Path)
		}

		response := models.BarsResponse{
			Status: "ok",
			Time:   []int64{1609459200, 1609545600},
			Open:   []float64{100.0, 101.0},
			High:   []float64{102.0, 103.0},
			Low:    []float64{99.0, 100.0},
			Close:  []float64{101.0, 102.0},
			Volume: []float64{1000.0, 2000.0},
		}
		w.WriteHeader(http.StatusOK)
		json.NewEncoder(w).Encode(response)
	}))
	defer server.Close()

	client := NewClient()
	client.baseURL = server.URL

	from := time.Date(2021, 1, 1, 0, 0, 0, 0, time.UTC)
	to := time.Date(2021, 1, 2, 0, 0, 0, 0, time.UTC)

	bars, err := client.GetHistoricalPrices("TEST-A", "D", from, to)
	if err != nil {
		t.Fatalf("GetHistoricalPrices() error = %v", err)
	}

	if len(bars.Time) != 2 {
		t.Errorf("Expected 2 bars, got %d", len(bars.Time))
	}

	if bars.Close[0] != 101.0 {
		t.Errorf("Expected first close 101.0, got %f", bars.Close[0])
	}
}

func TestGetFundLatest(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		expectedPath := "/funds/F000001/latest"
		if r.URL.Path != expectedPath {
			t.Errorf("Expected path %s, got %s", expectedPath, r.URL.Path)
		}

		response := models.FundLatestResponse{
			Status: true,
			Data: models.FundLatest{
				FundID:    "F000001",
				ShortCode: "TEST-A",
				Value:     10.5,
			},
		}
		w.WriteHeader(http.StatusOK)
		json.NewEncoder(w).Encode(response)
	}))
	defer server.Close()

	client := NewClient()
	client.baseURL = server.URL

	latest, err := client.GetFundLatest("F000001")
	if err != nil {
		t.Fatalf("GetFundLatest() error = %v", err)
	}

	if latest.ShortCode != "TEST-A" {
		t.Errorf("Expected ShortCode TEST-A, got %s", latest.ShortCode)
	}

	if latest.Value != 10.5 {
		t.Errorf("Expected Value 10.5, got %f", latest.Value)
	}
}
