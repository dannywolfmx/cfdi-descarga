package download

import (
	"bytes"
	"fmt"
	"io"

	"github.com/dannywolfmx/cfdi-descarga/certificate"
	"github.com/dannywolfmx/cfdi-descarga/client"
	"github.com/dannywolfmx/cfdi-descarga/signature"
	"github.com/dannywolfmx/cfdi-descarga/util"
	"github.com/dannywolfmx/cfdi-descarga/xmltemplate"
)

const (
	digestFormat    = `<des:peticiondescargamasivatercerosentrada xmlns:des="http://DescargaMasivaTerceros.sat.gob.mx"><des:peticiondescarga idpaquete="%s" rfcsolicitante="%s"></des:peticiondescarga></des:peticiondescargamasivatercerosentrada>`
	url             = `https://cfdidescargamasiva.clouda.sat.gob.mx/DescargaMasivaService.svc`
	signatureFormat = `<SignedInfo xmlns="http://www.w3.org/2000/09/xmldsig#"><CanonicalizationMethod Algorithm="http://www.w3.org/2001/10/xml-exc-c14n#"></CanonicalizationMethod><SignatureMethod Algorithm="http://www.w3.org/2000/09/xmldsig#rsa-sha1"></SignatureMethod><Reference URI=""><Transforms><Transform Algorithm="http://www.w3.org/2001/10/xml-exc-c14n#"></Transform></Transforms><DigestMethod Algorithm="http://www.w3.org/2000/09/xmldsig#sha1"></DigestMethod><DigestValue>%s</DigestValue></Reference></SignedInfo>`
	soapAction      = `http://DescargaMasivaTerceros.sat.gob.mx/IDescargaMasivaTercerosService/Descargar`
	urlDigestValue  = "http://DescargaMasivaTerceros.sat.gob.mx"
)

type download struct {
	IDPackage      string
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
	Cer                              *certificate.Certificate
	Key                              *signature.Key
	RFCSolicitante, IDPackage, Token string
}

// MasiveDownload will performe a download from all the ids
func MasiveDownload(ids []string, r *RequestData) [][]byte {
	files := [][]byte{}
	for _, id := range ids {
		r.IDPackage = id
		file, err := r.SendRequest()
		if err != nil {
			//skip
			continue
		}
		files = append(files, file)
	}

	return files
}

// SendRequest performe a "Download" request to the SAT server
func (r *RequestData) SendRequest() ([]byte, error) {
	rawReq, err := r.getRequest()
	if err != nil {
		return nil, err
	}

	var buff bytes.Buffer

	if err := xmltemplate.TemplateDownload(&buff, rawReq); err != nil {
		return nil, err
	}

	rawResponse, err := send(&buff, r.Token)

	if err != nil {
		return nil, err
	}

	encodeFile, err := ExtractResponse(rawResponse)

	if err != nil {
		return nil, err
	}

	return util.DecodeBase64(encodeFile)
}

// GetRequest return a new struct to send to the download server
func (r *RequestData) getRequest() (interface{}, error) {
	digestXML := fmt.Sprintf(digestFormat, r.IDPackage, r.RFCSolicitante)
	digestValue, err := r.Key.Hash([]byte(digestXML))

	if err != nil {
		return nil, fmt.Errorf("error hashing the digest value: %s", err)
	}

	signedSignature, err := r.Key.SignMessage(fmt.Sprintf(signatureFormat, digestValue))
	if err != nil {
		return nil, fmt.Errorf("error signing the digest value: %s", err)
	}

	return &download{
		IDPackage:      r.IDPackage,
		RFCSolicitante: r.RFCSolicitante,
		Certificate:    r.Cer.GetEncodeCertificate(),
		IssuerName:     r.Cer.ExtractIssuer(),
		SerialNumber:   r.Cer.SerialNumber.String(),
		DigestValue:    digestValue,
		Signature:      signedSignature,
	}, nil
}

// send will perform a request to the server
func send(body io.Reader, token string) ([]byte, error) {
	headers := map[string]string{
		"SOAPAction":    soapAction,
		"Authorization": fmt.Sprintf(`WRAP access_token="%s"`, token),
	}

	request := &client.Request{
		Body:    body,
		URL:     url,
		Method:  "POST",
		Headers: headers,
	}

	response, err := request.Send()
	if err != nil {
		return nil, err
	}

	defer response.Body.Close()
	return io.ReadAll(response.Body)
}
