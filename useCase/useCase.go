package usecase

import (
	"strings"

	"github.com/go-playground/validator/v10"
	"github.com/jia-hua/fizz-buzz-lbc/pkg/fizzbuzz"
)

// ComputeFizzBuzzRequest represents the computeFizzBuzz use case parameters
type ComputeFizzBuzzRequest struct {
	Limit      int    `form:"limit" validate:"required"`
	FizzNumber int    `form:"fizzNumber" validate:"required"`
	FizzString string `form:"fizzString" validate:"required"`
	BuzzNumber int    `form:"buzzNumber" validate:"required"`
	BuzzString string `form:"buzzString" validate:"required"`
}

// ComputeFizzBuzzResponse represents the computeFizzBuzz use case result
type ComputeFizzBuzzResponse string

// ComputeFizzBuzzHandler is the entrypoint of the use case for computing a fizz buzz sequence
func ComputeFizzBuzzHandler(request ComputeFizzBuzzRequest) (ComputeFizzBuzzResponse, error) {
	validate := validator.New()
	if err := validate.Struct(&request); err != nil {
		return "", err
	}

	sequence := fizzbuzz.ComputeSequence(request.Limit, request.FizzNumber, request.FizzString, request.BuzzNumber, request.BuzzString)

	result := ComputeFizzBuzzResponse(strings.Join(sequence, ","))

	return result, nil
}
