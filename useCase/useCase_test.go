package usecase

import (
	"testing"

	"github.com/stretchr/testify/require"
)

func TestComputeFizzBuzzShouldHaveGoodBehavior(t *testing.T) {
	// GIVEN good parameters
	goodRequest := ComputeFizzBuzzRequest{
		Limit:      10,
		FizzNumber: 2,
		FizzString: "fizz",
		BuzzNumber: 5,
		BuzzString: "buzz",
	}

	// WHEN I call the function ComputeFizzBuzzHandler
	result, err := ComputeFizzBuzzHandler(goodRequest)

	// THEN it should behave correctly
	expectedResult := ComputeFizzBuzzResponse("1,fizz,3,fizz,buzz,fizz,7,fizz,9,fizzbuzz")
	require.NoError(t, err)
	require.Equal(t, expectedResult, result)

}

func TestComputeFizzBuzzShouldThrowError(t *testing.T) {
	// GIVEN bad parameters (empty value here)
	badRequest := ComputeFizzBuzzRequest{
		Limit:      10,
		FizzNumber: 2,
		FizzString: "fizz",
		BuzzNumber: 5,
		BuzzString: "",
	}

	// WHEN I call the function ComputeFizzBuzzHandler
	_, err := ComputeFizzBuzzHandler(badRequest)

	// THEN it should return an error
	require.Error(t, err)
}
