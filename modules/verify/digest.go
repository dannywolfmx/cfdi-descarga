package verify

import (
	"encoding/xml"
)

//containerDigestValue struct to contain a digest value struct,
//this will generate a parent xml node.
//
//Digest cant be empty because the pointer type
type containerDigestValue struct {
	XMLName xml.Name `xml:"des:VerificaSolicitudDescarga"`
	URL     string   `xml:"xmlns:des,attr"`
	Digest  *digestValue
}

//digestValue contain the values of the inner xml digest value
//if RFCEmisor and RFCReceptor are empty this attributes will not be generated
type digestValue struct {
	XMLName        xml.Name `xml:"des:solicitud"`
	IDSolicitud    string   `xml:"IdSolicitud,attr"`
	RFCSolicitante string   `xml:"RfcSolicitante,attr"`
}

//getDigestValueXML generate a well formed xml with the struct given
func getDigestValueXML(rfcSolicitante, idRequest, url string) ([]byte, error) {
	content := &containerDigestValue{
		URL: url,
		Digest: &digestValue{
			IDSolicitud:    idRequest,
			RFCSolicitante: rfcSolicitante,
		},
	}
	return xml.Marshal(content)
}
