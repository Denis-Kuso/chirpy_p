package handlers

import (
    "strings"
    "fmt"
    "slices"
)
// need to check if there is any flagged keywords
// return "clean version" of input chirp
// example
//  "body": "This is a kerfuffle opinion I need to share with the world"
//  "This is a **** opinion I need to share with the world"
// List of bad words
    // 
    // kerfuffle
    // sharbert
    // fornax

//pseudo:
// split input into words
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
    fmt.Printf("Will return: %s\n",strings.Join(words, " "))
    return strings.Join(words, " ")
}

// separate inputText into words ignoring punctuation
func parseText(inputText string) []string {
    sep := " "
    listOfWords := strings.Split(inputText, sep)
    return listOfWords
}

