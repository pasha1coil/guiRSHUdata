package utils

import (
	"crypto/md5"
	"encoding/hex"
)

func GenerateHash(name string) string {
	hasher := md5.New()
	hasher.Write([]byte(name))
	return hex.EncodeToString(hasher.Sum(nil))
}
