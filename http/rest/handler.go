package rest

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/go-playground/validator/v10"

	"github.com/jia-hua/fizz-buzz-lbc/pkg/metrics"
	fizzbuzzUseCase "github.com/jia-hua/fizz-buzz-lbc/useCase"
)

// InitHandler associates use cases to routes of the service.
// func InitHandler(metricsService http.Handler, counterMetric *prometheus.CounterVec) *gin.Engine {
func InitHandler(metricsService metrics.Service) *gin.Engine {
	router := gin.New()
	router.Use(gin.Recovery())

	router.GET("/fizzbuzz", computeFizzBuzzHandler(metricsService))

	router.GET("/metrics", metricsHandler(metricsService))

	return router
}

func computeFizzBuzzHandler(metricsService metrics.Service) func(c *gin.Context) {
	return func(c *gin.Context) {
		var request fizzbuzzUseCase.ComputeFizzBuzzRequest

		if err := c.ShouldBindQuery(&request); err != err {
			c.JSON(http.StatusBadRequest, gin.H{"error": err})
			return
		}
		validate := validator.New()
		if err := validate.Struct(&request); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": err})
			c.Abort()
			return
		}

		metricsService.IncrementFizzBuzz(request.Limit, request.FizzNumber, request.FizzString, request.BuzzNumber, request.BuzzString)

		result := fizzbuzzUseCase.ComputeFizzBuzzHandler(request)

		c.String(200, "%v", result)
	}
}

func metricsHandler(metricsService metrics.Service) func(c *gin.Context) {
	return func(c *gin.Context) {
		metricsService.GetMetricHandler().ServeHTTP(c.Writer, c.Request)
	}
}
