package domain

import (
	"math/rand"
	"regexp"
	"strings"
)

type Int64Slice []int64

func (slice Int64Slice) Search(value int64) bool {
	for i := range slice {
		if slice[i] == value {
			return true
		}
	}
	return false
}

const letterBytes = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789"

func RandStringBytes(n int) string {
	b := make([]byte, n)
	for i := range b {
		b[i] = letterBytes[rand.Intn(len(letterBytes))]
	}
	return string(b)
}

var spaceSymbols = regexp.MustCompile("[:—\\-–]+|\\s+")

func FilterAnswer(answer string) string {
	result := spaceSymbols.ReplaceAllString(answer, " ")
	result = strings.TrimSpace(result)
	result = strings.ToLower(result)
	return result
}
