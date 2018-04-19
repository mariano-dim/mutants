package main

import (
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
)

func TestMutantCheckMutant(t *testing.T) {

	dnaJSON := []string{
		`{"dna":["TTAAAT","TACTCC","ATACAC","AAGACT","CCCCTT","AAAAAT"]}`, // dos hallazgos horizontales
		`{"dna":["TTAAAA","TACTCC","ATACAC","AAGACT","CCCCTT","AAAAAT"]}`, // tres hallazgos horizontales
		`{"dna":["TTAAAA","TACTCC","ATACAC","AAAACT","CCCCTT","AAAAAT"]}`, // cuatro hallazgos horizontales
		`{"dna":["TTAAAA","TACTCC","ATCCCC","AAAACT","CCCCTT","AAAAAT"]}`, // cinco hallazgos horizontales
		`{"dna":["TTAAAA","TACCCC","ATCCCC","AAAACT","CCCCTT","AAAAAT"]}`, // seis hallazgos horizontales
		`{"dna":["TTAAAT","TACTCC","ATACAC","AAGACT","CCCCTT","AAAAAT"]}`, // dos hallazgos horizontales m trspuesta
		`{"dna":["TTAAAA","TACTCC","ATACAC","AAGACT","CCCCTT","AAAAAT"]}`, // tres hallazgos horizontales m trspuesta
		`{"dna":["TTAAAA","TACTCC","ATACAC","AAAACT","CCCCTT","AAAAAT"]}`, // cuatro hallazgos horizontales m trspuesta
		`{"dna":["TTAAAA","TACTCC","ATCCCC","AAAACT","CCCCTT","AAAAAT"]}`, // cinco hallazgos horizontales m trspuesta
		`{"dna":["TTAAAA","TACCCC","ATCCCC","AAAACT","CCCCTT","AAAAAT"]}`, // seis hallazgos horizontales m trspuesta
		`{"dna":["AAAAAT","AGACCG","ATCGCT","AAATCC","ATAGAT","AGAACG"]}`,
	}

	for _, dna := range dnaJSON {
		reader := strings.NewReader(dna)

		req, err := http.NewRequest("POST", "/mutants", reader)
		if err != nil {
			t.Fatal(err)
		}
		rr := httptest.NewRecorder()

		handler := http.HandlerFunc(MutantCheck)

		handler.ServeHTTP(rr, req)

		if status := rr.Code; status != http.StatusOK {
			t.Errorf("handler returned wrong status code: got %v want %v",
				status, http.StatusOK)
		}
	}

}

func TestPing(t *testing.T) {
	req, err := http.NewRequest("GET", "/ping", nil)
	if err != nil {
		t.Fatal(err)
	}
	rr := httptest.NewRecorder()

	handler := http.HandlerFunc(Ping)

	handler.ServeHTTP(rr, req)

	// Check the status code is what we expect.
	if status := rr.Code; status != http.StatusOK {
		t.Errorf("handler returned wrong status code: got %v want %v",
			status, http.StatusOK)
	}
}

func TestMutantCheckInvalid(t *testing.T) {

	dnaJSON := []string{
		`{"dna" :["TTAAAT", "TACTCC", "ATACAC", "AAGACT", "CCCXTT", "AAAAAT"]}`, // Con un elemento extraño
		`{"dna" :["TTAAAT", "TACTCC", "ATaCAC", "AAGACT", "CCACTT", "ATGAAT"]}`,
		`{"dna" :["asAAAT", "TACTCC", "ATTCAC", "AAGACT", "CCACTT", "ATGAAT"]}`,
		`{"dna" :["TTAAAT", "TACTCC", "123456", "AAGACT", "CCACTT", "ATGAAT"]}`,
		`{"dna" :["TTAAAT", "123TCC", "ATaCAC", "AAGACT", "CCACTT", "ATGAATASDA"]}`,
	}

	for _, dna := range dnaJSON {
		reader := strings.NewReader(dna)

		req, err := http.NewRequest("POST", "/mutants", reader)
		if err != nil {
			t.Fatal(err)
		}
		rr := httptest.NewRecorder()

		handler := http.HandlerFunc(MutantCheck)

		handler.ServeHTTP(rr, req)

		if status := rr.Code; status != http.StatusForbidden {
			t.Errorf("handler returned wrong status code: got %v want %v",
				status, http.StatusOK)
		}
	}

}

func TestMutantCheckNOTMutant(t *testing.T) {

	dnaJSON := []string{
		`{"dna":["ATCGAT","CGATCG","ATCGAT","CGATCG","ATCGAT","CGATCG"]}`, // Todos elem distintos
		`{"dna":["CGATCG","ATCGAT","ATTGAT","CGATCG","ATCGAT","CGATCG"]}`,
	}

	for _, dna := range dnaJSON {
		reader := strings.NewReader(dna)

		req, err := http.NewRequest("POST", "/mutants", reader)
		if err != nil {
			t.Fatal(err)
		}
		rr := httptest.NewRecorder()

		handler := http.HandlerFunc(MutantCheck)

		handler.ServeHTTP(rr, req)

		if status := rr.Code; status != http.StatusForbidden {
			t.Errorf("handler returned wrong status code: got %v want %v",
				status, http.StatusOK)
		}

	}

}

func TestStatsMutants(t *testing.T) {

	req, err := http.NewRequest("GET", "/stats", nil)
	if err != nil {
		t.Fatal(err)
	}
	rr := httptest.NewRecorder()

	handler := http.HandlerFunc(StatsMutants)

	handler.ServeHTTP(rr, req)

	// Check the status code is what we expect.
	if status := rr.Code; status != http.StatusOK {
		t.Errorf("handler returned wrong status code: got %v want %v",
			status, http.StatusOK)
	}
}
