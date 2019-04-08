package main

import (
	"bytes"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestGenerateFactsHandle(t *testing.T) {
	router := setupRouter()
	tt := []struct {
		name       string
		q          string
		httpStatus int
	}{
		{"t1", "/api/generate", http.StatusBadRequest},
		{"t2", "/api/generate?n=10&m=12", http.StatusOK},
		{"t2", "/api/generate?n=10&m=wrong", http.StatusBadRequest},
	}

	for _, tc := range tt {
		t.Run(tc.name, func(t *testing.T) {
			req, _ := http.NewRequest("GET", tc.q, nil)
			w := httptest.NewRecorder()
			router.ServeHTTP(w, req)
			assert.Equal(t, tc.httpStatus, w.Code)
			var f Facts
			if w.Code == http.StatusOK {
				if err := f.UnmarshalBinary([]byte(w.Body.String())); err != nil {
					t.Fatal(err)
				}
			}

		})
	}
}
func TestGenerateFactsHandle2(t *testing.T) {
	router := setupRouter()
	req, _ := http.NewRequest("GET", "/api/generate?n=10&m=d12", nil)
	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)
	assert.Equal(t, http.StatusBadRequest, w.Code)
}

func TestFindHandle(t *testing.T) {
	router := setupRouter()
	facts := Facts{Costs: [][]float64{
		{8, 3, 1},
		{4, 7, 4},
		{5, 2, 6}},

		Demands:  []float64{70, 60, 30},
		Supplies: []float64{30, 90, 50},
	}
	b, _ := facts.MarshalBinary()
	req, _ := http.NewRequest("POST", "/api/find", bytes.NewReader(b))
	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)
	assert.Equal(t, http.StatusOK, w.Code)
	var d Decision
	if err := d.UnmarshalBinary([]byte(w.Body.String())); err != nil {
		t.Fatal(err)
	}
}

func TestFindHandle1(t *testing.T) {
	router := setupRouter()
	req, _ := http.NewRequest("POST", "/api/find", bytes.NewReader([]byte("fake body")))
	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)
	assert.Equal(t, http.StatusBadRequest, w.Code)
}
