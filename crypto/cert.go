package crypto

import (
	"crypto/rand"
	"crypto/rsa"
	"crypto/x509"
	"crypto/x509/pkix"
	"encoding/pem"
	"fmt"
	"io/ioutil"
	"math/big"
	mrand "math/rand"
	"net"
	"os"
	"time"
)

// LoadRootCa load ca.
func LoadRootCa(caFile string, keyFile string) (*x509.Certificate, *rsa.PrivateKey, error) {
	caBytes, err := ioutil.ReadFile(caFile)
	if err != nil {
		return nil, nil, err
	}

	caBlock, _ := pem.Decode(caBytes)
	if caBlock == nil {
		return nil, nil, fmt.Errorf("certificate info block not found")
	}

	cert, err := x509.ParseCertificate(caBlock.Bytes)
	if err != nil {
		return nil, nil, err
	}

	keybytes, err := ioutil.ReadFile(keyFile)
	if err != nil {
		return nil, nil, err
	}

	keyBlock, _ := pem.Decode(keybytes)
	if keyBlock == nil {
		return nil, nil, fmt.Errorf("private key info block not found ")
	}

	praKey, err := x509.ParsePKCS1PrivateKey(keyBlock.Bytes)
	if err != nil {
		return nil, nil, err
	}

	return cert, praKey, nil
}

// KeyUsage: Valid operations supported by keys
// KeyUsageDigitalSignature    数据签名
// KeyUsageContentCommitment   数据验证
// KeyUsageKeyEncipherment     秘钥加密
// KeyUsageDataEncipherment    数据加密
// KeyUsageKeyAgreement        秘钥协商
// KeyUsageCertSign            证书签名
// KeyUsageCRLSign             证书吊销签名
// KeyUsageEncipherOnly        只用于加密
// KeyUsageDecipherOnly        只用于解密

// ExtKeyUsage: Valid extended operations supported by keys
// ExtKeyUsageAny                                任何扩展操作
// ExtKeyUsageServerAuth                         服务端身份验证
// ExtKeyUsageClientAuth                         客户端身份验证
// ExtKeyUsageCodeSigning                        代码签名
// ExtKeyUsageEmailProtection                    安全电子邮件
// ExtKeyUsageIPSECEndSystem                     IPSEC终端系统
// ExtKeyUsageIPSECTunnel                        IPSEC隧道
// ExtKeyUsageIPSECUser                          IPSEC用户身份
// ExtKeyUsageTimeStamping                       时间戳
// ExtKeyUsageOCSPSigning                        OCSP签名
// ExtKeyUsageMicrosoftServerGatedCrypto         微软服务网关加密
// ExtKeyUsageNetscapeServerGatedCrypto          Netscape服务网管加密
// ExtKeyUsageMicrosoftCommercialCodeSigning     微软商业代码签名
// ExtKeyUsageMicrosoftKernelCodeSigning         微软内核代码签名

// IssueCertImpl Issue sub-certificates based on the root certificate
func IssueCertImpl(rootCa *x509.Certificate, rootKey *rsa.PrivateKey) (*pem.Block, *pem.Block, error) {
	template := &x509.Certificate{
		SerialNumber: big.NewInt(mrand.Int63()),
		Subject: pkix.Name{
			Country:            []string{"CN"},
			Organization:       []string{"Easy"},
			OrganizationalUnit: []string{"Easy"},
			Locality:           []string{"city"},
			Province:           []string{"prov"},
			StreetAddress:      []string{""},
			PostalCode:         []string{""},
			SerialNumber:       "",
			CommonName:         "example.com",
		},
		NotBefore:             time.Now(),                  // Certificate validity start time
		NotAfter:              time.Now().AddDate(1, 0, 0), // Certificate validity ends after (year, month, day)
		BasicConstraintsValid: true,                        // basic validity constraints
		IsCA:                  false,                       // Can it be used as a root certificate
		ExtKeyUsage:           []x509.ExtKeyUsage{x509.ExtKeyUsageAny},
		KeyUsage:              x509.KeyUsageDigitalSignature | x509.KeyUsageDataEncipherment,
		DNSNames:              []string{"www.example.com", "example.com"},
		EmailAddresses:        []string{"example@test.com"},
		IPAddresses:           []net.IP{net.ParseIP("192.168.1.1")},
	}

	praKey, pubKey, err := GenerateKeyRSA(2048)
	if err != nil {
		return nil, nil, err
	}

	derCert, err := x509.CreateCertificate(rand.Reader, template, rootCa, pubKey, rootKey)
	if err != nil {
		return nil, nil, err
	}
	pemCert := &pem.Block{Type: "CERTIFICATE", Bytes: derCert}

	buf := x509.MarshalPKCS1PrivateKey(praKey)
	pemKey := &pem.Block{Type: "RSA PRIVATE KEY", Bytes: buf}

	return pemCert, pemKey, err
}

// SaveCert2File Save the certificate as file.
func SaveCert2File(cert *pem.Block, key *pem.Block, certFile string, keyFile string) error {
	certOut, err := os.Create(certFile)
	if err != nil {
		return err
	}

	err = pem.Encode(certOut, cert)
	if err != nil {
		return err
	}

	keyOut, err := os.Create(keyFile)
	if err != nil {
		return err
	}

	err = pem.Encode(keyOut, key)
	if err != nil {
		return err
	}

	return nil
}

// GetCertBytes Save the certificate as bytes.
func GetCertBytes(cert *pem.Block, key *pem.Block) ([]byte, []byte) {
	certBytes := pem.EncodeToMemory(cert)
	keyBytes := pem.EncodeToMemory(key)

	return certBytes, keyBytes
}
