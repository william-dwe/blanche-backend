package util

import (
	"math/rand"
	"time"
)

func GenerateRandomNumber(max int, min int) int {
	rand.Seed(time.Now().UnixNano())
	return rand.Intn(max-min) + min
}
