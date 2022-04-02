package numberutil

import (
	"testing"
)

func TestFactorial(t *testing.T) {
	type expectedResult struct {
		n        uint32
		expected uint32
	}

	expectedValues := []expectedResult{
		{
			5, 120,
		},
		{
			10, 3628800,
		},
		{
			0, 1,
		},
		{
			1, 1,
		},
	}

	for _, expectedResult := range expectedValues {
		result := Factorial(expectedResult.n)
		if expectedResult.expected != result {
			t.Errorf("Arg = %d,  Exptected %d but %d", expectedResult.n, expectedResult.expected, result)
		}
	}
}

func TestPermutation(t *testing.T) {
	type expectedResult struct {
		n        uint32
		r        uint32
		expected uint32
	}
	expectedValues := []expectedResult{
		{
			4, 2, 12,
		},
		{
			5, 3, 60,
		},
		{
			1, 1, 1,
		},
		{
			0, 0, 1,
		},
		{
			1, 0, 1,
		},
	}
	for _, expectedResult := range expectedValues {
		result, _ := PermutationCount(expectedResult.n, expectedResult.r)
		if expectedResult.expected != result {
			t.Errorf("Arg = %d - %d,  Exptected %d but %d", expectedResult.n, expectedResult.r, expectedResult.expected, result)
		}
	}

	_, err := PermutationCount(1, 2)
	if err == nil {
		t.Errorf("Must throw error on n < r")
	}
}
func TestCombination(t *testing.T) {
	type expectedResult struct {
		n        uint32
		r        uint32
		expected uint32
	}
	expectedValues := []expectedResult{
		{
			5, 2, 10,
		},
		{
			1, 1, 1,
		},
		{
			0, 0, 1,
		},
		{
			1, 0, 1,
		},
	}
	for _, expectedResult := range expectedValues {
		result, _ := CombinationCount(expectedResult.n, expectedResult.r)
		if expectedResult.expected != result {
			t.Errorf("Arg = %d - %d,  Exptected %d but %d", expectedResult.n, expectedResult.r, expectedResult.expected, result)
		}
	}

	_, err := CombinationCount(1, 2)
	if err == nil {
		t.Errorf("Must throw error on n < r")
	}
}
