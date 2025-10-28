package utils

import (
	"crypto/sha256"
	"encoding/base64"
)

func HashToken256(token string) string {
	hash := sha256.Sum256([]byte(token))
	return base64.RawURLEncoding.EncodeToString(hash[:])
}
