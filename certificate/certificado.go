package certificate

import (
	"crypto/x509"
	"fmt"
	"os"
	"strings"

	"github.com/dannywolfmx/cfdi-descarga/util"
)

type Certificate struct {
	*x509.Certificate
}

func (c *Certificate) GetEncodeCertificate() string {
	return util.EncodeBase64(c.Raw)
}

//ExtractIssuer
//The core librery pkix doent return the expectet issuerto string like the SAT expected
//this function try to solve this

var attributeTypeNames = map[string]string{
	"2.5.4.6":              "C",
	"2.5.4.10":             "O",
	"2.5.4.11":             "OU",
	"2.5.4.3":              "CN",
	"2.5.4.5":              "SERIALNUMBER",
	"2.5.4.7":              "L",
	"2.5.4.8":              "ST",
	"2.5.4.9":              "street",
	"2.5.4.17":             "postalCode",
	"1.2.840.113549.1.9.2": "unstructuredName",
	"1.2.840.113549.1.9.1": "emailAddress",
	"2.5.4.45":             "x500UniqueIdentifier",
}

func (c *Certificate) ExtractIssuer() string {
	rdns := c.Issuer.Names

	s := []string{}

	for _, v := range rdns {
		oid, ok := attributeTypeNames[v.Type.String()]
		if ok {
			s = append(s, fmt.Sprintf("%s=%s", oid, v.Value))
		} else {
			s = append(s, fmt.Sprintf("%s=%s", v.Type.String(), v.Value))
		}
	}

	return strings.Join(s, ",")
}

func GetCertificate(path string) (*Certificate, error) {

	contentFile, err := os.ReadFile(path)
	if err != nil {
		return nil, err
	}

	//Get the certificate
	cer, err := x509.ParseCertificate(contentFile)

	if err != nil {
		return nil, err
	}

	return &Certificate{cer}, nil
}
