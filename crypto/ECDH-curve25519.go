package crypto

import (
	"crypto"
	"crypto/rand"
	"encoding/hex"
	"fmt"
	"io"

	"golang.org/x/crypto/curve25519"
)

type curve25519Impl struct{}

func newCurve25519Impl() *curve25519Impl {
	return &curve25519Impl{}
}

func (e *curve25519Impl) generateKeyPair(rand io.Reader) (crypto.PrivateKey, crypto.PublicKey, error) {
	var pubKey, praKey [32]byte
	var err error
	_, err = io.ReadFull(rand, praKey[:])
	if err != nil {
		return nil, nil, err
	}

	praKey[0] &= 248
	praKey[31] &= 127
	praKey[31] |= 64

	curve25519.ScalarBaseMult(&pubKey, &praKey)

	return &praKey, &pubKey, nil
}

func (e *curve25519Impl) generateSharedSecret(praKey crypto.PrivateKey, pubKey crypto.PublicKey) []byte {
	var pra, pub, secret *[32]byte
	pra = praKey.(*[32]byte)
	pub = pubKey.(*[32]byte)
	secret = new([32]byte)

	curve25519.ScalarMult(secret, pra, pub)

	return secret[:]
}

func (e *curve25519Impl) marshal(p crypto.PublicKey) []byte {
	pub := p.(*[32]byte)
	return pub[:]
}

func (e *curve25519Impl) unmarshal(data []byte) (crypto.PublicKey, error) {
	var pubKey [32]byte
	if len(data) != 32 {
		return nil, fmt.Errorf("Unmarshal pubKeyBuf error len =  %d", len(data))
	}

	copy(pubKey[:], data)

	return &pubKey, nil
}

// Curve22519Ecdh curve22519 ECDH.
type Curve22519Ecdh struct {
	impl       *curve25519Impl
	PublicKey  string
	PrivateKey string
}

// GenerateKeyPair Generate a public and private key pair.
func (c *Curve22519Ecdh) GenerateKeyPair() error {
	pri, pub, err := c.impl.generateKeyPair(rand.Reader)
	if err != nil {
		return err
	}

	c.PrivateKey = hex.EncodeToString(c.impl.marshal(pri))
	c.PublicKey = hex.EncodeToString(c.impl.marshal(pub))

	return err
}

// GetPublicKey get public key.
func (c *Curve22519Ecdh) GetPublicKey() string {
	return c.PublicKey
}

// GetSharedKey get shared key.
func (c *Curve22519Ecdh) GetSharedKey(peerPubKey string) (string, error) {
	hPubKey, err := hex.DecodeString(peerPubKey)
	if err != nil {
		return "", err
	}

	hPriKey, err := hex.DecodeString(c.PrivateKey)
	if err != nil {
		return "", err
	}

	cPeerPubKey, err := c.impl.unmarshal([]byte(hPubKey))
	if err != nil {
		return "", err
	}

	cLocalPriKey, err := c.impl.unmarshal([]byte(hPriKey))
	if err != nil {
		return "", err
	}

	bShareKey := c.impl.generateSharedSecret(cLocalPriKey, cPeerPubKey)
	shareKey := hex.EncodeToString(bShareKey)

	return shareKey, nil
}

// NewCurve22519Ecdh create curve22519 ecdh.
func NewCurve22519Ecdh() *Curve22519Ecdh {
	return &Curve22519Ecdh{impl: newCurve25519Impl()}
}
