package crypto

import (
	"crypto/md5"
	"crypto/sha1"
	"crypto/sha256"
	"crypto/sha512"
	"hash/crc32"
	"hash/crc64"
)

// CRC32Cp hash.
func CRC32Cp(data []byte) uint32 {
	var tab = crc32.MakeTable(crc32.Castagnoli)
	return crc32.Checksum(data, tab)
}

// CRC32IEEE hash.
func CRC32IEEE(data []byte) uint32 {
	var tab = crc32.MakeTable(crc32.IEEE)
	return crc32.Checksum(data, tab)
}

// CRC32Kp hash.
func CRC32Kp(data []byte) uint32 {
	var tab = crc32.MakeTable(crc32.Koopman)
	return crc32.Checksum(data, tab)
}

// CRC64ISO hash.
func CRC64ISO(data []byte) uint64 {
	var tab = crc64.MakeTable(crc64.ISO)
	return crc64.Checksum(data, tab)
}

// CRC64ECMA hash.
func CRC64ECMA(data []byte) uint64 {
	var tab = crc64.MakeTable(crc64.ECMA)
	return crc64.Checksum(data, tab)
}

// Sha1 hash.
func Sha1(data []byte) ([]byte, error) {
	h := sha1.New()
	_, err := h.Write(data)
	if err != nil {
		return nil, err
	}
	hash := h.Sum(nil)

	return hash, err
}

// Sha256 hash.
func Sha256(data []byte) ([]byte, error) {
	h := sha256.New()
	_, err := h.Write(data)
	if err != nil {
		return nil, err
	}
	hash := h.Sum(nil)

	return hash, err
}

// Sha512 hash.
func Sha512(data []byte) ([]byte, error) {
	h := sha512.New()
	_, err := h.Write(data)
	if err != nil {
		return nil, err
	}
	hash := h.Sum(nil)

	return hash, err
}

// MD5 hash.
func MD5(data []byte) ([]byte, error) {
	m := md5.New()
	_, err := m.Write(data)
	if err != nil {
		return nil, err
	}

	hash := m.Sum(nil)

	return hash, nil
}
