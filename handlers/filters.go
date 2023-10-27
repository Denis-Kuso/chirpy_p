package handlers

import (
    "strings"
    "fmt"
    "slices"
)

// any words match prohibited words, replace with replacement pattern
var profaneWords = []string{"kerfuffle", "sharbert", "fornax"}
const replacementPattern = "****"

func FilterText(text string) string {
    // TODO
    words := parseText(text)
    var indices []int
    for index, word := range words {
        if slices.Contains(profaneWords,strings.ToLower(word)) {
            indices = append(indices, index)
        }
    }
    for _, indx := range indices {
        words[indx] = replacementPattern
    }
    return strings.Join(words, " ")
}

// separate inputText into words ignoring punctuation
func parseText(inputText string) []string {
    sep := " "
    listOfWords := strings.Split(inputText, sep)
    return listOfWords
}

