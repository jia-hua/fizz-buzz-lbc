package main

import (
	"github.com/jia-hua/fizz-buzz-lbc/http/rest"
	"github.com/jia-hua/fizz-buzz-lbc/internal/metrics"
)

func main() {

	metricService := metrics.New()

	router := rest.InitHandler(metricService)

	router.Run(":8080") // listen and serve on 0.0.0.0:8080
}
