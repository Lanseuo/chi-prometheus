package main

import (
	"fmt"
	"math/rand"
	"net/http"
	"time"

	chiprometheus "github.com/Lanseuo/chi-prometheus"
	"github.com/go-chi/chi"
	"github.com/prometheus/client_golang/prometheus/promhttp"
)

func main() {
	n := chi.NewRouter()
	m := chiprometheus.NewMiddleware("serviceName")
	// if you want to use other buckets than the default (300, 1200, 5000) you can run:
	// m := negroniprometheus.NewMiddleware("serviceName", 400, 1600, 700)

	n.Use(m)

	n.Handle("/metrics", promhttp.Handler())
	n.Get("/ok", func(w http.ResponseWriter, r *http.Request) {
		sleep := rand.Intn(4999) + 1
		time.Sleep(time.Duration(sleep) * time.Millisecond)
		w.WriteHeader(http.StatusOK)
		fmt.Fprintf(w, "Slept %d milliseconds\n", sleep)
	})
	n.Get("/notfound", func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusNotFound)
		fmt.Fprintln(w, "Not found")
	})
	n.Get(`/users/{name}`, func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
		fmt.Fprintln(w, "Ok")
	})

	http.ListenAndServe(":3000", n)
}
