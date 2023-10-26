package handlers

import (
    "testing"
)

// input for Filter
// Empty string
// no keyword
// multiple keywords
// Keywords with exclamation
// uppercase

type Test struct {
    Name string
    Input string
    Expected string
}

func TestFilterText(t *testing.T) {
    tests := []Test{
        { "empty string",
        "' '",
        "' '"},
        { "No keyword present",
        "I had something interesting for breakfast",
        "I had something interesting for breakfast"},
        { "Keywords with puncuation",
        "Sharbert!",
        "Sharbert!"},
        { "Multiple keywords and case correction",
        "I really need a kerfuffle to go to bed sooner, Fornax !",
        "I really need a **** to go to bed sooner, **** !"},
    }

    for _, datum := range tests {
		result := FilterText(datum.Input) 

		if result != datum.Expected {
            t.Errorf("FAILED on test case: %s.\nExpected: %s , got: %s\n", datum.Name,
				datum.Expected, result)
		}
	}
}


