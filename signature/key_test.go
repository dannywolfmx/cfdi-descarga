package signature

import (
	"bytes"
	"crypto"
	"crypto/rand"
	"crypto/rsa"
	"crypto/sha1"
	"testing"
)

func TestGetKey(t *testing.T) {
	key, err := ExtractKey("../ejemploKey.key", "12345678a")
	if err != nil {
		t.Fatalf("Should get a key: %s", err)
	}
	testMessage := "Hola mundo"
	encryptedBytes, err := rsa.EncryptOAEP(
		sha1.New(),
		rand.Reader,
		&key.PrivateKey.PublicKey,
		[]byte(testMessage),
		nil,
	)
	if err != nil {
		t.Fatalf("Should generate an encrypted message: %s", err)
	}

	_, err = key.PrivateKey.Decrypt(nil, encryptedBytes, &rsa.OAEPOptions{Hash: crypto.SHA1})
	if err != nil {
		t.Fatalf("Should generate an decrypted message: %s", err)
	}
}

// TestSign will check if exist an error
func TestSign(t *testing.T) {
	privateKey, err := ExtractKey("../ejemploKey.key", "12345678a")
	testMessage := []byte("Hola mundo")
	if err != nil {
		t.Fatalf("Should get a certificate: %s", err)
	}
	var buff bytes.Buffer
	_, err = privateKey.Sign(&buff, testMessage)

	if err != nil {
		t.Fatalf("Should get a signature %s", err)
	}
}

// TestSignXML check if the message is signed and converted to base64
func TestSignXML(t *testing.T) {
	message := `<SignedInfo xmlns="http://www.w3.org/2000/09/xmldsig#"><CanonicalizationMethod Algorithm="http://www.w3.org/2001/10/xml-exc-c14n#"></CanonicalizationMethod><SignatureMethod Algorithm="http://www.w3.org/2000/09/xmldsig#rsa-sha1"></SignatureMethod><Reference URI=""><Transforms><Transform Algorithm="http://www.w3.org/2001/10/xml-exc-c14n#"></Transform></Transforms><DigestMethod Algorithm="http://www.w3.org/2000/09/xmldsig#sha1"></DigestMethod><DigestValue>2zavoO/4lAZSUSFVRBpLbiEcfKg=</DigestValue></Reference></SignedInfo>`
	//I cant test the raw result directly, thats why i need to convert the result to base64
	encodeResultTest := `Ir4nciZnmdEiThsnqMWlVDl5ZHrDAao0pvvxoleZDgwhgoTifQ602JtKCiAqsUAPuSfYlnXgPH5mOU2cgY46f+jfspdjivgwZCeJA/lCWt83/v8wg78Xaz1vNRkmDOX0Iu0HsGUISFlVM8ycf5BdkOPQHqvHIQdkHjOnP/j9w48Kn+xXTbODk8F2syM3W87y1K8FatRIltCUZPS/AciTMnv/FjAeBe1w7HJy5iCJEnax1gGdkJo9VpBrU9w5JBzTbHritcHZm97nQ6exZ92sdQtWtvWRrcvIaEUM0iofU186w/I/PYhZbL91xEXon3OGh5+Ra7rkbgCJIbZ8Q1clXA==`

	keyPath := "../test2.key"
	key, err := ExtractKey(keyPath, "12345678a")

	if err != nil {
		t.Fatalf("Should get a key: %s", err)
	}
	signedMessage, err := key.SignMessage(message)

	if err != nil {
		t.Fatalf("Should get a signature err: %s", err)
	}

	if signedMessage != encodeResultTest {
		t.Fatalf("Should the encodeTestResult and the encode result be the same\n encodeTest: \n%s\n result: \n%s\n", encodeResultTest, signedMessage)
	}
}
