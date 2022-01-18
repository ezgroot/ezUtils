package crypto

import (
	"crypto/hmac"
	"crypto/sha256"
	"encoding/hex"
)

func GetHmacSha256Hex(key string, data []byte) string {
	h := hmac.New(sha256.New, []byte(key))
	h.Write(data)
	out := h.Sum(nil)
	return hex.EncodeToString(out)
}
