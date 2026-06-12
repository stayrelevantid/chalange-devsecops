package crypto

import (
	"crypto/md5"
	"encoding/hex"
)

// HashPassword uses MD5 — this is intentionally insecure for SAST demo (Day 08)
// DO NOT use in production — MD5 is cryptographically broken
func HashPassword(password string) string {
	h := md5.Sum([]byte(password))
	return hex.EncodeToString(h[:])
}
