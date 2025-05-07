package signature

import (
	"bytes"
	"testing"
	"text/template"
)

type testFielGenerator struct {
	UUID        string
	Created     string
	Expired     string
	DigestValue string
}

func (f *testFielGenerator) GenerateDate(diference int) (string, string) {
	return f.Created, f.Expired
}

func (f *testFielGenerator) GenerateUUIDV4() (string, error) {
	return f.UUID, nil
}

func (f *testFielGenerator) GenerateDigestValue(created, expired string) (string, error) {
	return f.DigestValue, nil
}

func (f *testFielGenerator) GenerateSignatureValue(digestValue string, key *Key) ([]byte, error) {
	xmlTemplate := `<SignedInfo xmlns="http://www.w3.org/2000/09/xmldsig#"><CanonicalizationMethod Algorithm="http://www.w3.org/2001/10/xml-exc-c14n#"></CanonicalizationMethod><SignatureMethod Algorithm="http://www.w3.org/2000/09/xmldsig#rsa-sha1"></SignatureMethod><Reference URI="#_0"><Transforms><Transform Algorithm="http://www.w3.org/2001/10/xml-exc-c14n#"></Transform></Transforms><DigestMethod Algorithm="http://www.w3.org/2000/09/xmldsig#sha1"></DigestMethod><DigestValue>{{.DigestValue}}</DigestValue></Reference></SignedInfo>`
	//TODO mover esta generacion de template a un singleton para no crearlo cada vez que se ejecutra la funcion
	tmpl, err := template.New("SignatureValue").Parse(xmlTemplate)
	if err != nil {
		return nil, err
	}

	var value bytes.Buffer

	//Well formed data struct with the create and expried dates
	data := struct {
		DigestValue string
	}{
		digestValue,
	}

	if err = tmpl.Execute(&value, data); err != nil {
		return nil, err
	}

	var buff bytes.Buffer
	_, err = key.Sign(&buff, value.Bytes())

	return buff.Bytes(), err
}

func TestGenerateSignatureValue(t *testing.T) {
	key, err := ExtractKey("../ejemplo.key", "12345678a")
	if err != nil {
		t.Fatalf("Should get a key: %s", err)
	}
	generator := testFielGenerator{}
	testDigestValue := "5iiWNNO7aJndihjEU6ROuDK1gzE="
	_, err = generator.GenerateSignatureValue(testDigestValue, key)
	if err != nil {
		t.Fatalf("Should get a signature %s", err)
	}

}
