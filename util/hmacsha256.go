package util

import (
	"crypto/hmac"
	"crypto/sha256"
	"encoding/hex"
)

func HashSHA256(inputStr string, key string) (string, error) {
	h := hmac.New(sha256.New, []byte(key))

	h.Write([]byte(inputStr))

	sha := hex.EncodeToString(h.Sum(nil))

	return sha, nil
}
