package learning

import (
	"crypto/rand"
	"crypto/rsa"
	"crypto/x509"
	"encoding/pem"
	"log"
	"strings"
	"testing"
)

func TestPrivateKeyGenerationInPEMFormat(t *testing.T) {
	key, err := rsa.GenerateKey(rand.Reader, 1024)
	if err != nil {
		t.Fatal(err)
	}
	block := &pem.Block{
		Type:  "RSA PRIVATE KEY",
		Bytes: x509.MarshalPKCS1PrivateKey(key),
	}
	pemdata := string(pem.EncodeToMemory(block))
	if !strings.HasPrefix(pemdata, "-----BEGIN RSA PRIVATE KEY-----") {
		t.Fatal("Expecting PEM data to start with BEGIN RSA PRIVATE KEY")
	}
	if !strings.HasSuffix(pemdata, "-----END RSA PRIVATE KEY-----\n") {
		t.Fatal("Expecting PEM data to end with END RSA PRIVATE KEY")
	}
}

func TestEncryptedPrivateKeyGenerationInPEMFormat(t *testing.T) {
	key, err := rsa.GenerateKey(rand.Reader, 1024)
	if err != nil {
		t.Fatal(err)
	}
	block, err := x509.EncryptPEMBlock(rand.Reader, "RSA PRIVATE KEY", x509.MarshalPKCS1PrivateKey(key),
		[]byte("passphrase"), x509.PEMCipherAES256)
	if err != nil {
		log.Fatal(err)
	}
	pemdata := string(pem.EncodeToMemory(block))
	if !strings.HasPrefix(pemdata, "-----BEGIN RSA PRIVATE KEY-----\nProc-Type: 4,ENCRYPTED") {
		t.Fatal("Expecting PEM data to start with BEGIN RSA PRIVATE KEY")
	}
	if !strings.HasSuffix(pemdata, "-----END RSA PRIVATE KEY-----\n") {
		t.Fatal("Expecting PEM data to end with END RSA PRIVATE KEY")
	}
}
