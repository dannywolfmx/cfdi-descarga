package digest

import (
	"bytes"
	"encoding/xml"
	"fmt"
	"io"

	"github.com/dannywolfmx/cfdi-descarga/client"
	"github.com/dannywolfmx/cfdi-descarga/signature"
	"github.com/dannywolfmx/cfdi-descarga/util"
	"github.com/dannywolfmx/cfdi-descarga/xmltemplate"
)

const (
	authDigestFormat    = `<u:Timestamp xmlns:u="http://docs.oasis-open.org/wss/2004/01/oasis-200401-wss-wssecurity-utility-1.0.xsd" u:Id="_0"><u:Created>%s</u:Created><u:Expires>%s</u:Expires></u:Timestamp>`
	authSignatureFormat = `<SignedInfo xmlns="http://www.w3.org/2000/09/xmldsig#"><CanonicalizationMethod Algorithm="http://www.w3.org/2001/10/xml-exc-c14n#"></CanonicalizationMethod><SignatureMethod Algorithm="http://www.w3.org/2000/09/xmldsig#rsa-sha1"></SignatureMethod><Reference URI="#_0"><Transforms><Transform Algorithm="http://www.w3.org/2001/10/xml-exc-c14n#"></Transform></Transforms><DigestMethod Algorithm="http://www.w3.org/2000/09/xmldsig#sha1"></DigestMethod><DigestValue>%s</DigestValue></Reference></SignedInfo>`
	soapAction          = `http://DescargaMasivaTerceros.gob.mx/IAutenticacion/Autentica`
	url                 = `https://cfdidescargamasivasolicitud.clouda.sat.gob.mx/Autenticacion/Autenticacion.svc`
)

const minutes = 60 * 12

var clientRequest = &client.Request{
	Method: "POST",
	Headers: map[string]string{
		"SOAPAction": soapAction,
	},
	URL: url,
}

// Authentication represent the needed data to fill the XML to share with the server
type Authentication struct {
	Created     string
	Expires     string
	Certificate string
	UUID        string
	Digest      *Digest
}

type requestAuth struct {
	EncodedCertificate string
	Key                *signature.Key
	GenerateUUID       UUIDV4Generator
	request            *client.Request
}

func NewRequestAuth(encodedCertificate string, key *signature.Key, uuidGenerator UUIDV4Generator) *requestAuth {
	return &requestAuth{
		EncodedCertificate: encodedCertificate,
		Key:                key,
		GenerateUUID:       uuidGenerator,
		request:            clientRequest,
	}
}

func (d *requestAuth) GetRequest() (*Authentication, error) {
	uuid, err := util.GenerateUUIDV4()
	if err != nil {
		return nil, err
	}

	createdAt, expiresAt := generateDate(minutes)
	digest := DigestGenerator{
		Key:             d.Key,
		DigestFormat:    authDigestFormat,
		SignatureFormat: authSignatureFormat,
		CreatedAt:       createdAt,
		ExpiresAt:       expiresAt,
	}

	digestData, err := digest.New()
	if err != nil {
		return nil, err
	}

	return &Authentication{
		Created:     createdAt,
		Expires:     expiresAt,
		UUID:        fmt.Sprintf("uuid-%s-4", uuid),
		Certificate: d.EncodedCertificate,
		Digest:      digestData,
	}, nil
}

type token struct {
	//The tag represent the node of the token
	Value string `xml:"Body>AutenticaResponse>AutenticaResult"`
}

// SendRequest return a token
func (r *requestAuth) RequestToken() (string, error) {
	req, err := r.GetRequest()
	if err != nil {
		return "", fmt.Errorf("error getting request: %w", err)
	}

	var buff bytes.Buffer
	if err = xmltemplate.TemplateAuth(&buff, req); err != nil {
		return "", fmt.Errorf("error creating XML template: %w", err)
	}

	r.request.Body = &buff

	response, err := r.request.Send()
	if err != nil {
		return "", fmt.Errorf("error sending request: %w", err)
	}

	data, err := io.ReadAll(response.Body)
	if err != nil {
		return "", fmt.Errorf("error reading response body: %w", err)
	}

	var t token
	if err := xml.Unmarshal(data, &t); err != nil {
		return "", fmt.Errorf("error unmarshalling XML: %w", err)
	}

	if t.Value == "" {
		return "", fmt.Errorf("empty token value in response")
	}

	return t.Value, nil
}
