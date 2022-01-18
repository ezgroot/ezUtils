package aes

import (
	"crypto/aes"
	"crypto/cipher"
	"fmt"
)

type ecb struct {
	b         cipher.Block
	blockSize int
}

func newECB(b cipher.Block) *ecb {
	return &ecb{b: b, blockSize: b.BlockSize()}
}

// *********************************** AES ECB Encrypt impl ********************************************

type ecbEncrypt ecb

// newECBEncrypt returns a BlockMode which encrypts in electronic code book mode, using the given Block.
func newECBEncrypt(b cipher.Block) cipher.BlockMode {
	return (*ecbEncrypt)(newECB(b))
}

func (x *ecbEncrypt) BlockSize() int {
	return x.blockSize
}

func (x *ecbEncrypt) CryptBlocks(dst []byte, src []byte) {
	for len(src) > 0 {
		x.b.Encrypt(dst, src[:x.blockSize])
		src = src[x.blockSize:]
		dst = dst[x.blockSize:]
	}
}

// *********************************** AES ECB Decrypt impl ********************************************

type ecbDecrypt ecb

// newECBDecrypt returns a BlockMode which decrypts in electronic code book mode, using the given Block.
func newECBDecrypt(b cipher.Block) cipher.BlockMode {
	return (*ecbDecrypt)(newECB(b))
}

func (x *ecbDecrypt) BlockSize() int {
	return x.blockSize
}

func (x *ecbDecrypt) CryptBlocks(dst []byte, src []byte) {
	for len(src) > 0 {
		x.b.Decrypt(dst, src[:x.blockSize])
		src = src[x.blockSize:]
		dst = dst[x.blockSize:]
	}
}

// *********************************** AES ECB API ********************************************

// EncryptECB AES encrypt, ECB mode.
func EncryptECB(key string, plainText []byte) ([]byte, error) {
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

	ecb := newECBEncrypt(block)
	content := plainText
	content = PKCS7Padding(content, block.BlockSize())
	if len(content)%block.BlockSize() != 0 {
		return nil, fmt.Errorf("date len=%d invalid, must to be divisible by BlockSize=%d", len(key), block.BlockSize())
	}

	cipherText := make([]byte, len(content))
	if len(cipherText) < len(content) {
		return nil, fmt.Errorf("output smaller than input")
	}

	ecb.CryptBlocks(cipherText, content)

	return cipherText, nil
}

// DecryptECB AES decrypt, ECB mode.
func DecryptECB(key string, cipherText []byte) ([]byte, error) {
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

	blockMode := newECBDecrypt(block)

	if len(cipherText)%blockMode.BlockSize() != 0 {
		return nil, fmt.Errorf("date len=%d invalid, must to be divisible by BlockSize=%d", len(key), blockMode.BlockSize())
	}

	plainText := make([]byte, len(cipherText))
	if len(plainText) < len(plainText) {
		return nil, fmt.Errorf("output smaller than input")
	}

	blockMode.CryptBlocks(plainText, cipherText)
	plainText, err = PKCS7UnPadding(plainText)
	if err != nil {
		return nil, err
	}

	return plainText, nil
}
