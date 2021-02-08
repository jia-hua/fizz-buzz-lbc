package main

import (
	"github.com/jia-hua/fizz-buzz-lbc/http/rest"
	"github.com/jia-hua/fizz-buzz-lbc/pkg/metrics"

	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promauto"
	"github.com/prometheus/client_golang/prometheus/promhttp"
)

func main() {

	fizzBuzzCounter := promauto.NewCounterVec(prometheus.CounterOpts{
		Namespace: "fizzbuzz",
		Name:      "fizzbuzz_metrics",
		Help:      "fizzbuzz metrics shows statistics of the GET /fizzbuzz endpoint",
	}, []string{"limit", "fizz", "fizzString", "buzz", "buzzString"})

	metricsHandler := promhttp.Handler()

	metricService := metrics.MetricContext{
		FizzBuzzCounter: fizzBuzzCounter,
		MetricsHandler:  metricsHandler,
	}

	router := rest.InitHandler(metricService)

	router.Run() // listen and serve on 0.0.0.0:8080
}
