package main

import (
	"bytes"
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestHandler(t *testing.T) {
	payload := bytes.NewBuffer(testdataJSON)
	req, err := http.NewRequest("POST", "/", payload)
	if err != nil {
		t.Fatal(err)
	}

	rr := httptest.NewRecorder()
	handler := http.HandlerFunc(simulationHandler)

	handler.ServeHTTP(rr, req)

	if status := rr.Code; status != http.StatusOK {
		t.Errorf("handler returned wrong status code: got %v want %v",
			status, http.StatusOK)
	}

	// Check the response body size (it should be reasonably big)
	// We can't compare output as it's non-deterministic (each simulation
	// produce different results).
	// TODO(divan): decode response and compare data sizes fields, at least.
	if len(rr.Body.String()) < 100 {
		t.Errorf("handler returned unexpected body: got %v",
			rr.Body.String())
	}
}
