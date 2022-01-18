package crypto

import (
	"encoding/base64"
)

// Base64Encode base64 encode.
func Base64Encode(data []byte) string {
	return base64.StdEncoding.EncodeToString(data)
}

// Base64Decode base64 decode.
func Base64Decode(data string) ([]byte, error) {
	return base64.StdEncoding.DecodeString(data)
}

// Base64EncodeURL base64 url encode.
func Base64EncodeURL(data []byte) string {
	return base64.URLEncoding.EncodeToString(data)
}

// Base64DecodeURL base64 url decode.
func Base64DecodeURL(data string) ([]byte, error) {
	return base64.URLEncoding.DecodeString(data)
}
