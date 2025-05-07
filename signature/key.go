package signature

import (
	"bytes"
	"crypto"
	"crypto/rand"
	"crypto/rsa"
	"crypto/sha1"
	"crypto/x509"
	"encoding/pem"
	"io"
	"os"

	"github.com/dannywolfmx/cfdi-descarga/util"
	"github.com/youmark/pkcs8"
)

type Key struct {
	PrivateKey *rsa.PrivateKey
}

// Sign a message and verify the result with a private key, return raw signed message if not errors
// This method will  verify the signed result message, you dont need to do it again
func (c *Key) Sign(buff io.Writer, message []byte) (int, error) {
	msgHash := sha1.New()
	_, err := msgHash.Write(message)
	if err != nil {
		return 0, err
	}

	hashed := msgHash.Sum(nil)

	signature, err := rsa.SignPKCS1v15(rand.Reader, c.PrivateKey, crypto.SHA1, hashed)
	if err != nil {
		return 0, err
	}

	err = rsa.VerifyPKCS1v15(&c.PrivateKey.PublicKey, crypto.SHA1, hashed, signature)
	if err != nil {
		return 0, err
	}

	return buff.Write(signature)
}

// SignMessage will return a base64 signed xml if not errors
func (key *Key) SignMessage(message string) (string, error) {
	var buff bytes.Buffer
	//Sign content
	_, err := key.Sign(&buff, []byte(message))

	//Verify errors
	if err != nil {
		return "", err
	}
	//Encode the raw result to base64
	return util.EncodeBase64(buff.Bytes()), nil
}

func (c *Key) Hash(message []byte) (string, error) {

	hash := sha1.Sum(message)
	return util.EncodeBase64(hash[:]), nil
}

// ExtractKey will extract the key from an encripted PKCS8 format file
func ExtractKey(path, password string) (*Key, error) {
	contentFile, err := os.ReadFile(path)
	if err != nil {
		return nil, err
	}
	key, err := pkcs8.ParsePKCS8PrivateKeyRSA(contentFile, []byte(password))

	if err != nil {
		return nil, err
	}

	return &Key{
		PrivateKey: key,
	}, nil
}

// ExtractKeyFromPem will extract the key from an unencripted PEM file
func ExtractKeyFromPem(path string) (*Key, error) {
	contentFile, err := os.ReadFile(path)
	if err != nil {
		return nil, err
	}

	block, _ := pem.Decode(contentFile)

	privateKey, err := x509.ParsePKCS8PrivateKey(block.Bytes)

	if err != nil {
		return nil, err
	}

	return &Key{
		PrivateKey: privateKey.(*rsa.PrivateKey),
	}, nil
}
