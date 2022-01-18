package aes

import (
	"crypto/aes"
	"crypto/cipher"
	"crypto/rand"
	"fmt"
	"io"
)

// EncryptCBC AES encrypt func, CBC mode.
func EncryptCBC(key string, plainText []byte) ([]byte, error) {
	if len(plainText) == 0 {
		return nil, fmt.Errorf("plainText is empty")
	}

	if (len(key) != 16) && (len(key) != 24) && (len(key) != 32) {
		return nil, fmt.Errorf("aes key len = %d error, must equal 16/24/32", len(key))
	}

	block, err := aes.NewCipher([]byte(key))
	if err != nil {
		return nil, err
	}

	blockSize := block.BlockSize()
	plainText = PKCS7Padding(plainText, blockSize)
	cipherText := make([]byte, blockSize+len(plainText))

	iv := cipherText[:blockSize]
	if _, err := io.ReadFull(rand.Reader, iv); err != nil {
		return nil, err
	}

	mode := cipher.NewCBCEncrypter(block, iv)
	mode.CryptBlocks(cipherText[blockSize:], plainText)

	return cipherText, nil
}

// DecryptCBC AES decrypt func, CBC mode.
func DecryptCBC(key string, cipherText []byte) ([]byte, error) {
	if len(cipherText) == 0 {
		return nil, fmt.Errorf("plainText is empty")
	}

	if (len(key) != 16) && (len(key) != 24) && (len(key) != 32) {
		return nil, fmt.Errorf("aes key len = %d error, must equal 16/24/32", len(key))
	}

	block, err := aes.NewCipher([]byte(key))
	if err != nil {
		return nil, err
	}

	blockSize := block.BlockSize()
	if len(cipherText) < blockSize {
		return nil, fmt.Errorf("cipherText len is less than blockSize = %d", blockSize)
	}

	iv := cipherText[:blockSize]
	cipherText = cipherText[blockSize:]

	// CBC mode always works in whole blocks.
	if len(cipherText)%blockSize != 0 {
		return nil, fmt.Errorf("cipherText is not a multiple of the block size")
	}

	mode := cipher.NewCBCDecrypter(block, iv)

	// CryptBlocks can work in-place if the two arguments are the same.
	mode.CryptBlocks(cipherText, cipherText)
	cipherText, err = PKCS7UnPadding(cipherText)
	if err != nil {
		return nil, err
	}

	return cipherText, nil
}
