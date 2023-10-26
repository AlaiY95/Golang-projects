package main

import (
	"fmt"
	"strings"
)

// originalLetter is a constant string representing the alphabet in uppercase.
const originalLetter = "ABCDEFGHIJKLMNOPQRSTUVWXYZ"

// hashLetterFn is a function that takes a key and a letter as input and returns a string.
// It reorders the characters in the letter based on the key.
func hashLetterFn(key int, letter string) (result string) {
	// Convert the letter to a slice of runes for manipulation.
	runes := []rune(letter)

	// Calculate the lastLetterKey by shifting the last 'key' characters to the front.
	lastLetterKey := string(runes[len(letter)-key : len(letter)])

	// Calculate the leftOversLetter by taking the remaining characters.
	leftOversLetter := string(runes[0 : len(letter)-key])

	// Concatenate lastLetterKey and leftOversLetter to form the result.
	return fmt.Sprintf(`%s%s`, lastLetterKey, leftOversLetter)
}

// encrypt is a function that takes a key and a plainText string as input and returns an encrypted string.
func encrypt(key int, plainText string) (result string) {

	// Call hashLetterFn to generate the hashed letter based on the key.
	hashLetter := hashLetterFn(key, originalLetter)

	// Initialize an empty string to store the encrypted result.
	var hashedString = ""

	// Define the findOne function that maps each character in plainText to its encrypted counterpart.
	findOne := func(r rune) rune {
		pos := strings.Index(originalLetter, string([]rune{r}))
		if pos != -1 {
			letterPosition := (pos + len(originalLetter)) % len(originalLetter)
			hashedString = hashedString + string(hashLetter[letterPosition])
			return r
		}
		return r
	}
	strings.Map(findOne, plainText)
	return hashedString
}
func decrypt(key int, encrypttedText string) (result string) {
	hashLetter := hashLetterFn(key, originalLetter)
	var hashedString = ""
	findOne := func(r rune) rune {
		pos := strings.Index(hashLetter, string([]rune{r}))
		if pos != -1 {
			letterPosition := (pos + len(originalLetter)) % len(originalLetter)
			hashedString = hashedString + string(originalLetter[letterPosition])
			return r
		}
		return r
	}
	strings.Map(findOne, encrypttedText)
	return hashedString
}

func main() {
	plainText := "HELLOWORLD"
	fmt.Println("Plain Text", plainText)
	encrypted := encrypt(5, plainText)
	fmt.Println("Encrypted", encrypted)
	decrypted := decrypt(5, encrypted)
	fmt.Println("Decrypted", decrypted)
}
