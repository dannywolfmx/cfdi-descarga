package verify

import (
	"bytes"
	"fmt"
	"io"

	"github.com/dannywolfmx/cfdi-descarga/certificate"
	"github.com/dannywolfmx/cfdi-descarga/client"
	"github.com/dannywolfmx/cfdi-descarga/signature"
	"github.com/dannywolfmx/cfdi-descarga/xmltemplate"
)

const (
	//downloadURIL is the url to connect with the server
	// idsolicitud and rfcsolicitante are the values to be replaced in the xml template
	digestFormat = `<des:verificasolicituddescarga xmlns:des="http://DescargaMasivaTerceros.sat.gob.mx"><des:solicitud idsolicitud="%s" rfcsolicitante="%s"></des:solicitud></des:verificasolicituddescarga>`
	downloadURL  = `https://cfdidescargamasivasolicitud.clouda.sat.gob.mx/VerificaSolicitudDescargaService.svc`

	//soapAction is the default url required by the server to know what action will be perforer
	soapAction = `http://DescargaMasivaTerceros.sat.gob.mx/IVerificaSolicitudDescargaService/VerificaSolicitudDescarga`

	//signatureFormat represent the XML format valid to SprintF() to embed the digest value in the xml
	// DigestValue is the value to be replaced in the xml template
	signatureFormat = `<SignedInfo xmlns="http://www.w3.org/2000/09/xmldsig#"><CanonicalizationMethod Algorithm="http://www.w3.org/2001/10/xml-exc-c14n#"></CanonicalizationMethod><SignatureMethod Algorithm="http://www.w3.org/2000/09/xmldsig#rsa-sha1"></SignatureMethod><Reference URI=""><Transforms><Transform Algorithm="http://www.w3.org/2001/10/xml-exc-c14n#"></Transform></Transforms><DigestMethod Algorithm="http://www.w3.org/2000/09/xmldsig#sha1"></DigestMethod><DigestValue>%s</DigestValue></Reference></SignedInfo>`

	urlDigestValue = "http://DescargaMasivaTerceros.sat.gob.mx"
)

type verify struct {
	IDSolicitud    string
	RFCSolicitante string
	Certificate    string
	IssuerName     string
	SerialNumber   string
	DigestValue    string
	Signature      string
}

// RequestData implement the RequestData interface to generate a request body
// This structs implement the client.RequestData
type RequestData struct {
	RFCSolicitante, IDSolicitudDeDescarga string
	Cer                                   *certificate.Certificate
	Key                                   *signature.Key
	Token                                 string
}

// SendRequest performe a "verify" request to the SAT server
func (r *RequestData) SendRequest() (*ResponseData, error) {
	req, err := r.getRequest()
	if err != nil {
		return nil, err
	}

	var buff bytes.Buffer
	if err := xmltemplate.TemplateVerify(&buff, req); err != nil {
		return nil, err
	}

	rawResponse, err := send(&buff, r.Token)

	if err != nil {
		return nil, err
	}

	return UnmarshalResponse(rawResponse)
}

// GetRequest return a new struct to send to the verify server
func (r *RequestData) getRequest() (*verify, error) {
	digestXML := fmt.Sprintf(digestFormat, r.IDSolicitudDeDescarga, r.RFCSolicitante)

	//Hash and encode the digest value
	digestValue, err := r.Key.Hash([]byte(digestXML))

	if err != nil {
		return nil, fmt.Errorf("error hashing the digest value: %s", err)
	}

	//Sign the digest value and encode the result to base64
	signedSignature, err := r.Key.SignMessage(fmt.Sprintf(signatureFormat, digestValue))
	if err != nil {
		return nil, fmt.Errorf("error signing the digest value: %s", err)
	}

	return &verify{
		IDSolicitud:    r.IDSolicitudDeDescarga,
		RFCSolicitante: r.RFCSolicitante,
		Certificate:    r.Cer.GetEncodeCertificate(),
		IssuerName:     r.Cer.ExtractIssuer(),
		SerialNumber:   r.Cer.SerialNumber.String(),
		DigestValue:    digestValue,
		Signature:      signedSignature,
	}, nil
}

func send(body io.Reader, token string) ([]byte, error) {
	headers := map[string]string{
		"SOAPAction":    soapAction,
		"Authorization": fmt.Sprintf(`WRAP access_token="%s"`, token),
	}
	request := &client.Request{
		Body:    body,
		URL:     downloadURL,
		Method:  "POST",
		Headers: headers,
	}

	resp, err := request.Send()
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	return io.ReadAll(resp.Body)
}
