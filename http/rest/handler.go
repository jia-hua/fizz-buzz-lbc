package rest

import (
	"net/http"

	"github.com/gin-gonic/gin"

	"github.com/jia-hua/fizz-buzz-lbc/internal/metrics"
	fizzbuzzUseCase "github.com/jia-hua/fizz-buzz-lbc/useCase"
)

// InitHandler associates use cases to routes of the service.
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
			c.JSON(http.StatusBadRequest, gin.H{"invalid format error": err})
			c.Abort()
			return
		}

		metricsService.IncrementFizzBuzz(request.Limit, request.FizzNumber, request.FizzString, request.BuzzNumber, request.BuzzString)

		result, err := fizzbuzzUseCase.ComputeFizzBuzzHandler(request)
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"invalid value error": err})
			c.Abort()
			return
		}

		c.String(http.StatusOK, "%v", result)
	}
}

func metricsHandler(metricsService metrics.Service) func(c *gin.Context) {
	return func(c *gin.Context) {
		metricsService.GetMetricHandler().ServeHTTP(c.Writer, c.Request)
	}
}
