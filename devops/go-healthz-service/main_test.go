// test files are not included in the actual program

package main

import (
	"net/http"
	"net/http/httptest"
	"testing"
)

// this test function is found by go automatically and run
func TestRoutes(t *testing.T) {
	//this is a table-driven test
	tests := []struct {
		name       string
		method     string
		path       string
		wantStatus int
		wantBody   string
	}{
		{"root", http.MethodGet, "/", http.StatusOK, "hello from go-healthz-service\n"},
		{"healthz", http.MethodGet, "/healthz", http.StatusOK, "ok\n"},
		{"unknown path", http.MethodGet, "/nope", http.StatusNotFound, ""},
		{"wrong method", http.MethodPost, "/healthz", http.StatusMethodNotAllowed, ""},
	}

	// this part gets the main.go file, runs it through the tests above each line, each as its own sub-test
	// each tests results are given separately, so we get to understand exactly what scenario failed
	mux := newMux()
	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			// this fakes a visitor in memory for test purposes
			req := httptest.NewRequest(tc.method, tc.path, nil) // creates a pretend request
			rec := httptest.NewRecorder()                       // creates a pretend notebook that captures the answer

			mux.ServeHTTP(rec, req) // with no network involvement, it hands the request to the router as if it came from 8080

			// this part checks the answer and if its wring it gives what it found versus what was expected
			// it checks the answer text but only for the lines that have expected text
			if rec.Code != tc.wantStatus {
				t.Fatalf("status = %d, want %d", rec.Code, tc.wantStatus)
			}
			if tc.wantBody != "" && rec.Body.String() != tc.wantBody {
				t.Fatalf("body = %q, want %q", rec.Body.String(), tc.wantBody)
			}
		})
	}
}
