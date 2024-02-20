package main

import (
    "net/http"
    "net/http/httptest"
    "testing"
)

func TestRateLimiting(t *testing.T) {
    // Create a request with a mock handler
    req, err := http.NewRequest("GET", "/test", nil)
    if err != nil {
        t.Fatal(err)
    }

    // Create a ResponseRecorder to record the response
    rr := httptest.NewRecorder()

    // Simulate multiple requests to test rate limiting
    for i := 0; i < 5; i++ {
        handleWithRateLimit(testHandler).ServeHTTP(rr, req)
        if rr.Code != http.StatusOK && i < 3 {
            t.Errorf("Expected status OK, got %v", rr.Code)
        }
        if rr.Code != http.StatusTooManyRequests && i >= 3 {
            t.Errorf("Expected status Too Many Requests, got %v", rr.Code)
        }
        rr = httptest.NewRecorder() // Reset the recorder for the next iteration
    }
}

// Mock handler function for testing
func testHandler(w http.ResponseWriter, r *http.Request) {
    w.WriteHeader(http.StatusOK)
}
