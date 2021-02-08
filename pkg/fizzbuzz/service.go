package fizzbuzz

import "strconv"

// ComputeSequence returns a list of strings with numbers from 1 to limit,
// where: all multiples of fromFizz are replaced by fizzString,
// all multiples of buzzFizz are replaced by buzzString,
// all multiples of fizzNumber and buzzNumber are replaced by str1str2.
func ComputeSequence(limit int, fizzNumber int, fizzString string, buzzNumber int, buzzString string) []string {

	result := []string{}

	for i := 1; i <= limit; i++ {
		current := computeNumber(i, fizzNumber, fizzString, buzzNumber, buzzString)
		result = append(result, current)
	}

	return result
}

func computeNumber(number int, fizzNumber int, fizzString string, buzzNumber int, buzzString string) string {
	result := ""

	if number%fizzNumber == 0 {
		result += fizzString
	}
	if number%buzzNumber == 0 {
		result += buzzString
	}

	if result == "" {
		result = strconv.Itoa(number)
	}

	return result
}
