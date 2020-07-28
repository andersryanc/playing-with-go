// Package gen provides utility functions to generate random strings.
package gen

import (
	"math/rand"
	"time"
)

// DefaultChars is the list of characters that will be used while generating strings when a chars of "" is provided.
const DefaultChars string = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789"

// DefaultLength is the number of characters that will be returned when a legnth of -1 is provided.
const DefaultLength int = 8

// RandomString will generate a random string of the provided length.
//   To use **gen.DefaultChars** provide chars as: ""
//   To use **gen.DefaultLength** provide length as: -1
func RandomString(chars string, length int) string {
	if length == -1 {
		length = DefaultLength
	}
	if chars == "" {
		chars = DefaultChars
	}

	r := rand.New(rand.NewSource(time.Now().UnixNano()))

	str := ""
	for i := 0; i < length; i++ {
		num := r.Intn(len(chars))
		str = str + string(chars[num])
	}
	return str
}

// RandomHexString will generate a random string of the provided length using only the hexidecimal character set: "ABCDEF1234567890"
//   To use the **gen.DefaultLength** provide length as: -1
func RandomHexString(length int) string {
	return RandomString("ABCDEF1234567890", length)
}
