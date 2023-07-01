package util

import (
	"hash/fnv"
	"time"

	"github.com/google/uuid"
)

func GenerateUUID() string {
	id := uuid.New()
	return id.String()
}

func GenerateUUIDWithDate() string {
	id := uuid.New()
	return id.String() + time.Now().Format("01022006")
}

func HashFnv(s string) uint32 {
	h := fnv.New32a()
	h.Write([]byte(s))
	return h.Sum32()
}
