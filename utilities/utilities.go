// Copyright Peter Lenson.
// All Rights Reserved

package utilities

import (
	"github.com/satori/go.uuid"
	"log"
	"math/rand"
)

// Generates a unique id
func GetUniqueID() string {
	// Creating UUID Version 4
	// panic on error
	id, err := uuid.NewV4()

	if err != nil {
		log.Printf("Problem generating id %v", err)
	}

	return id.String()
}

// Generates a random string of specified size
func RandomString(n int) string {
	var letter = []rune("abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789")

	b := make([]rune, n)
	for i := range b {
		b[i] = letter[rand.Intn(len(letter))]
	}
	return string(b)
}
