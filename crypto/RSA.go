package crypto

import (
	"crypto"
	"crypto/rand"
	"crypto/rsa"
	"crypto/sha256"
	"crypto/x509"
	"encoding/pem"
	"fmt"
	"os"
)

const (
	defaultRsaSize = 2048
)

// GenerateKeyRSA generate RSA public and private key pair.
func GenerateKeyRSA(bits int) (*rsa.PrivateKey, *rsa.PublicKey, error) {
	if bits < defaultRsaSize {
		bits = defaultRsaSize
	}

	privateKey, err := rsa.GenerateKey(rand.Reader, bits)
	if err != nil {
		return nil, nil, err
	}

	return privateKey, &privateKey.PublicKey, nil
}

// EncodePrivateKeyBufferRSA format RSA private key as bytes.
func EncodePrivateKeyBufferRSA(priKey *rsa.PrivateKey) []byte {
	key := x509.MarshalPKCS1PrivateKey(priKey)
	block := &pem.Block{
		Type:  "RSA PRIVATE KEY",
		Bytes: key,
	}

	return pem.EncodeToMemory(block)
}

// EncodePublicKeyBufferRSA format RSA public key as bytes.
func EncodePublicKeyBufferRSA(pubKey *rsa.PublicKey) ([]byte, error) {
	key, err := x509.MarshalPKIXPublicKey(pubKey)
	if err != nil {
		return nil, err
	}

	block := &pem.Block{
		Type:  "PUBLIC KEY",
		Bytes: key,
	}

	return pem.EncodeToMemory(block), nil
}

// EncodePrivateKeyFileRSA format RSA private key as file.
func EncodePrivateKeyFileRSA(priKey *rsa.PrivateKey, fileName string) error {
	keyBytes := x509.MarshalPKCS1PrivateKey(priKey)
	block := &pem.Block{
		Type:  "RSA PRIVATE KEY",
		Bytes: keyBytes,
	}

	file, err := os.Create(fileName)
	if err != nil {
		return err
	}

	err = pem.Encode(file, block)
	if err != nil {
		return err
	}

	return nil
}

// EncodePublicKeyFileRSA format RSA public key as file.
func EncodePublicKeyFileRSA(pubKey *rsa.PublicKey, fileName string) error {
	keyBytes, err := x509.MarshalPKIXPublicKey(pubKey)
	if err != nil {
		return err
	}

	block := &pem.Block{
		Type:  "PUBLIC KEY",
		Bytes: keyBytes,
	}

	file, err := os.Create(fileName)
	if err != nil {
		return err
	}

	err = pem.Encode(file, block)
	if err != nil {
		return err
	}

	return nil
}

// LoadRSAPrivateKeyPKCS1 Parse the private key from the key byte stream (the key must be formatted), using PKCS1.
func LoadRSAPrivateKeyPKCS1(priKey []byte) (*rsa.PrivateKey, error) {
	block, _ := pem.Decode(priKey)
	if block == nil {
		return nil, fmt.Errorf("block is nil")
	}

	privateKey, err := x509.ParsePKCS1PrivateKey([]byte(block.Bytes))
	if err != nil {
		return nil, err
	}

	return privateKey, nil
}

// LoadRSAPrivateKeyPKCS8 Parse the private key from the key byte stream (the key must be formatted), using PKCS8.
func LoadRSAPrivateKeyPKCS8(priKey []byte) (*rsa.PrivateKey, error) {
	block, _ := pem.Decode(priKey)
	if block == nil {
		return nil, fmt.Errorf("block is nil")
	}

	privateKey, err := x509.ParsePKCS8PrivateKey([]byte(block.Bytes))
	if err != nil {
		return nil, err
	}

	return privateKey.(*rsa.PrivateKey), nil
}

// LoadPublicKeyRSA Parse the public key from the public key byte stream (the key must be formatted).
func LoadPublicKeyRSA(pubKey []byte) (*rsa.PublicKey, error) {
	block, _ := pem.Decode(pubKey)
	if block == nil {
		return nil, fmt.Errorf("block is nil")
	}

	publicKey, err := x509.ParsePKIXPublicKey(block.Bytes)
	if err != nil {
		return nil, err
	}

	return publicKey.(*rsa.PublicKey), nil
}

// RsaEncryptRSA RSA encryption, using the peer's public key encryption.
func RsaEncryptRSA(plainText []byte, peerPubKey *rsa.PublicKey) ([]byte, error) {
	label := []byte("")
	sha256hash := sha256.New()
	cipherText, err := rsa.EncryptOAEP(sha256hash, rand.Reader, peerPubKey, plainText, label)

	return cipherText, err
}

// RsaDecryptRSA RSA decryption, decrypt with your own private key
func RsaDecryptRSA(cipherText []byte, myPriKey *rsa.PrivateKey) ([]byte, error) {
	label := []byte("")
	sha256hash := sha256.New()
	plainText, err := rsa.DecryptOAEP(sha256hash, rand.Reader, myPriKey, cipherText, label)
	if err != nil {
		return nil, fmt.Errorf("RSA decrypt failed, error=%s", err)
	}

	return plainText, nil
}

// SignWithSha256RSA Digitally sign the hash value of the data with your own private key.
func SignWithSha256RSA(data []byte, myPriKey []byte) ([]byte, error) {
	privateKey, err := LoadRSAPrivateKeyPKCS1(myPriKey)
	if err != nil {
		return nil, err
	}

	hash, err := Sha256(data)
	if err != nil {
		return nil, err
	}

	return rsa.SignPKCS1v15(rand.Reader, privateKey, crypto.SHA256, hash)
}

// VerySignWithSha256RSA Verifies the hash of the signature using the counterparty's public key.
func VerySignWithSha256RSA(originalDataHash []byte, signedDataHash []byte, peerPubKey []byte) error {
	publicKey, err := LoadPublicKeyRSA(peerPubKey)
	if err != nil {
		return err
	}

	return rsa.VerifyPKCS1v15(publicKey, crypto.SHA256, originalDataHash, signedDataHash)
}
