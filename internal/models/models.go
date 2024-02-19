package models

import (
	"crypto/sha256"
	"encoding/base64"
)

func NewShortID(value string) string {
	hash := sha256.New()
	hash.Write([]byte(value))

	encodedHash := base64.StdEncoding.EncodeToString(hash.Sum(nil))

	return encodedHash[:10]
}
