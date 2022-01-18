package aes

import (
	"bytes"
	"fmt"
)

// 0 padding mode，Use 0 padding when the data length is not aligned, otherwise no padding,
// When padding with ZeroPadding, there is no way to distinguish between real data and
// padding data, so it is only suitable for encryption and decryption of strings ending with \0.

func ZeroPadding(plainText []byte, blockSize int) []byte {
	padding := blockSize - len(plainText)%blockSize
	paddingText := bytes.Repeat([]byte{0}, padding)

	return append(plainText, paddingText...)
}

func ZeroUnPadding(plainText []byte) []byte {
	return bytes.TrimFunc(plainText,
		func(r rune) bool {
			return r == rune(0)
		})
}

// PKCS5 way to add padding, PKCS5 is only padding for 8 bytes (BlockSize=8), and the padding content is 0x01-0x08.

func PKCS5Padding(plainText []byte, blockSize int) ([]byte, error) {
	if blockSize != 8 {
		return nil, fmt.Errorf("blockSize = %d, must equal 8", blockSize)
	}

	padding := blockSize - len(plainText)%blockSize
	paddingText := bytes.Repeat([]byte{byte(padding)}, padding)

	return append(plainText, paddingText...), nil
}

func PKCS5UnPadding(plainText []byte) ([]byte, error) {
	length := len(plainText)
	if length <= 0 {
		return nil, fmt.Errorf("plaintext len <= 0")
	}

	// get padding length
	paddingSize := int(plainText[length-1])

	if paddingSize >= length {
		return nil, fmt.Errorf("padding len is big than plaintext")
	}

	// remove padding text
	return plainText[:(length - paddingSize)], nil
}

// PKCS7 way to add padding, PKCS7 is compatible with PKCS5.

func PKCS7Padding(plainText []byte, blockSize int) []byte {
	padding := blockSize - len(plainText)%blockSize
	paddingText := bytes.Repeat([]byte{byte(padding)}, padding)

	return append(plainText, paddingText...)
}

func PKCS7UnPadding(plainText []byte) ([]byte, error) {
	length := len(plainText)
	if length <= 0 {
		return nil, fmt.Errorf("plaintext len <= 0")
	}

	paddingSize := int(plainText[length-1])

	if paddingSize >= length {
		return nil, fmt.Errorf("padding len is big than plaintext")
	}

	return plainText[:(length - paddingSize)], nil
}
