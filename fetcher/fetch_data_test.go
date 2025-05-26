package fetcher

import (
	"io"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"
)

func TestFetchData_Success(t *testing.T) {
	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
		io.WriteString(w, `[{"id": 1, "title": "test"}]`)
	}))
	defer ts.Close()

	// Override the URL inside the function
	oldURL := apiURL
	apiURL = ts.URL
	defer func() { apiURL = oldURL }()

	data, err := FetchData()
	if err != nil {
		t.Fatalf("Expected no error, got: %v", err)
	}

	if len(data) == 0 {
		t.Fatal("Expected data, got empty response")
	}
}

func TestFetchData_Failure(t *testing.T) {
	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusInternalServerError)
	}))
	defer ts.Close()

	oldURL := apiURL
	apiURL = ts.URL
	defer func() { apiURL = oldURL }()

	_, err := FetchData()
	if err == nil {
		t.Fatal("Expected error, got nil")
	}
}

func TestFetchData_Timeout(t *testing.T) {
	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		time.Sleep(11 * time.Second)
		w.WriteHeader(http.StatusOK)
		io.WriteString(w, `[{"id": 1, "title": "delayed"}]`)
	}))
	defer ts.Close()

	oldURL := apiURL
	apiURL = ts.URL
	defer func() { apiURL = oldURL }()

	_, err := FetchData()
	if err == nil {
		t.Fatal("Expected timeout error,but got nil")
	}
}
