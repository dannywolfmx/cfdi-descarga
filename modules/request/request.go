package request

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
	signatureFormat = `<SignedInfo xmlns="http://www.w3.org/2000/09/xmldsig#"><CanonicalizationMethod Algorithm="http://www.w3.org/TR/2001/REC-xml-c14n-20010315"></CanonicalizationMethod><SignatureMethod Algorithm="http://www.w3.org/2000/09/xmldsig#rsa-sha1"></SignatureMethod><Reference URI=""><Transforms><Transform Algorithm="http://www.w3.org/2000/09/xmldsig#enveloped-signature"></Transform></Transforms><DigestMethod Algorithm="http://www.w3.org/2000/09/xmldsig#sha1"></DigestMethod><DigestValue>%s</DigestValue></Reference></SignedInfo>`
	soapAction      = `http://DescargaMasivaTerceros.sat.gob.mx/ISolicitaDescargaService/SolicitaDescarga`
	url             = `https://cfdidescargamasivasolicitud.clouda.sat.gob.mx/SolicitaDescargaService.svc`
	//Url to insert the xml digest value
	urlDigestValue = "http://DescargaMasivaTerceros.sat.gob.mx"
)

// request
// DigestValue is a encode base64
type request struct {
	*digestValue
	Base64DigestValueXML string
	Base64Signature      string
	TipoSolicitud        string
	X509Certificate      string
	X509IssuerName       string
	X509SerialNumber     string
}

// RequestData struct to fill with the request values
// RFCEmisor and RFCReceptor can be nil, without panic
type RequestData struct {
	//DigestValue data
	EndDate        string
	InitialDate    string
	RFCSolicitante string
	RFCEmisor      *string
	RFCReceptor    *string
	RequestType    string

	//Aditional data
	Token string
	Cer   *certificate.Certificate
	Key   *signature.Key
}

// SendRequest performe a "Download" request to the SAT server
func (r *RequestData) RequestIDSolicitud() (string, error) {
	req, err := r.getRequest()
	if err != nil {
		return "", fmt.Errorf("error generating request id solicitud: %w", err)
	}

	var buff bytes.Buffer
	if err := xmltemplate.TemplateRequest(&buff, req); err != nil {
		return "", fmt.Errorf("error generating request xml: %w", err)
	}

	rawResponse, err := send(&buff, r.Token)

	if err != nil {
		return "", fmt.Errorf("error sending request: %w", err)
	}

	idSolicitud, err := ExtractResponse(rawResponse)
	if err != nil {
		return "", fmt.Errorf("error extracting response: %w", err)
	}

	if idSolicitud == "" {
		return "", fmt.Errorf("empty idSolicitud value in response")
	}

	return idSolicitud, nil
}

// GetRequest generate a well formed request
func (r *RequestData) getRequest() (interface{}, error) {

	//Generate a well formed digest xml
	digestData := &digestValue{
		FechaFinal:     r.EndDate,
		FechaInicial:   r.InitialDate,
		RFCSolicitante: r.RFCSolicitante,
		RFCEmisor:      r.RFCEmisor,
		RFCReceptor:    r.RFCReceptor,
		TipoSolicitud:  r.RequestType,
	}

	digestXML, err := getDigestValueXML(digestData, urlDigestValue)

	if err != nil {
		return nil, fmt.Errorf("error generating digest value xml: %s", err)
	}

	//Get hashed base64 value
	digestValue, err := r.Key.Hash(digestXML)

	if err != nil {
		return nil, fmt.Errorf("error hashing the digest value: %s", err)
	}

	//Sign the digest value and encode the result to base64
	signedSignature, err := r.Key.SignMessage(fmt.Sprintf(signatureFormat, digestValue))
	if err != nil {
		return nil, fmt.Errorf("error signing the digest value: %s", err)
	}

	return &request{
		digestValue:          digestData,
		Base64DigestValueXML: digestValue,
		Base64Signature:      signedSignature,
		TipoSolicitud:        r.RequestType,
		X509Certificate:      r.Cer.GetEncodeCertificate(),
		X509IssuerName:       r.Cer.ExtractIssuer(),
		X509SerialNumber:     r.Cer.SerialNumber.String(),
	}, nil
}

// SendRequest performe a "Request" request to the SAT server
func send(body io.Reader, token string) ([]byte, error) {
	headers := map[string]string{
		"SOAPAction":    soapAction,
		"Authorization": fmt.Sprintf(`WRAP access_token="%s"`, token),
	}

	request := &client.Request{
		Body:    body,
		Method:  "POST",
		Headers: headers,
		URL:     url,
	}

	resp, err := request.Send()
	if err != nil {
		return nil, fmt.Errorf("error sending request: %w", err)
	}

	defer resp.Body.Close()
	return io.ReadAll(resp.Body)
}
