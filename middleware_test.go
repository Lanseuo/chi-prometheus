package chiprometheus

import (
	"fmt"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/go-chi/chi"
	"github.com/prometheus/client_golang/prometheus/promhttp"
	"github.com/stretchr/testify/assert"
)

func Test(t *testing.T) {
	recorder := httptest.NewRecorder()

	router := chi.NewRouter()
	middleware := NewMiddleware("test")
	router.Use(middleware)

	router.Handle("/metrics", promhttp.Handler())
	router.Get("/ok", func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
		fmt.Fprintln(w, "Ok")
	})
	router.Get("/users/{name}", func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
		fmt.Fprintln(w, "Hello "+chi.URLParam(r, "name"))
	})

	req, err := http.NewRequest("GET", "/ok", nil)
	if err != nil {
		t.Error(err)
	}
	router.ServeHTTP(recorder, req)

	req, err = http.NewRequest("GET", "/notfound", nil)
	if err != nil {
		t.Error(err)
	}
	router.ServeHTTP(recorder, req)

	req, err = http.NewRequest("GET", "/users/user1", nil)
	if err != nil {
		t.Error(err)
	}
	router.ServeHTTP(recorder, req)

	req, err = http.NewRequest("GET", "/users/user2", nil)
	if err != nil {
		t.Error(err)
	}
	router.ServeHTTP(recorder, req)

	req, err = http.NewRequest("GET", "/metrics", nil)
	if err != nil {
		t.Error(err)
	}
	router.ServeHTTP(recorder, req)

	metricsBody := recorder.Body.String()

	assert.Contains(t, metricsBody, `chi_request_duration_milliseconds_sum{code="200",method="GET",path="/ok",service="test"} `)
	assert.Contains(t, metricsBody, `chi_request_duration_milliseconds_count{code="200",method="GET",path="/ok",service="test"} 1`)

	assert.Contains(t, metricsBody, `chi_request_duration_milliseconds_sum{code="404",method="GET",path="",service="test"} `)
	assert.Contains(t, metricsBody, `chi_request_duration_milliseconds_count{code="404",method="GET",path="",service="test"} 1`)

	assert.Contains(t, metricsBody, `chi_request_duration_milliseconds_sum{code="200",method="GET",path="/users/{name}",service="test"} `)
	assert.Contains(t, metricsBody, `chi_request_duration_milliseconds_count{code="200",method="GET",path="/users/{name}",service="test"} 2`)
}
