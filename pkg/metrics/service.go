package metrics

import (
	"net/http"
	"strconv"

	"github.com/prometheus/client_golang/prometheus"
)

// MetricContext contains dependencies
type MetricContext struct {
	FizzBuzzCounter *prometheus.CounterVec
	MetricsHandler  http.Handler
}

// Service to collect metrics
type Service interface {
	IncrementFizzBuzz(limit int, fizzNumber int, fizzString string, buzzNumber int, buzzString string)
	GetMetricHandler() http.Handler
}

// IncrementFizzBuzz increments the count for the fizzbuzz endpoint
func (c MetricContext) IncrementFizzBuzz(limit int, fizzNumber int, fizzString string, buzzNumber int, buzzString string) {
	c.FizzBuzzCounter.WithLabelValues(
		strconv.Itoa(limit),
		strconv.Itoa(fizzNumber),
		fizzString,
		strconv.Itoa(buzzNumber),
		buzzString).Inc()
}

// GetMetricHandler returns the http handler
func (c MetricContext) GetMetricHandler() http.Handler {
	return c.MetricsHandler
}
