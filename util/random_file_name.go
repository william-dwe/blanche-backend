package util

import (
	"crypto/rand"
	"fmt"
)

func RandomFileName(fileNameLength int) string {
	b := make([]byte, fileNameLength)
	rand.Read(b)
	return fmt.Sprintf("%x", b)
}
