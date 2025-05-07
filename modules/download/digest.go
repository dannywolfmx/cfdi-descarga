package download

import (
	"encoding/xml"
)

//containerDigestValue struct to contain a digest value struct,
//this will generate a parent xml node.
//
//Digest cant be empty because the pointer type
type containerDigestValue struct {
	XMLName xml.Name `xml:"des:PeticionDescargaMasivaTercerosEntrada"`
	URL     string   `xml:"xmlns:des,attr"`
	Digest  *digestValue
}

//digestValue contain the values of the inner xml digest value
//if RFCEmisor and RFCReceptor are empty this attributes will not be generated
type digestValue struct {
	XMLName        xml.Name `xml:"des:peticionDescarga"`
	IDPackage      string   `xml:"IdPaquete,attr"`
	RFCSolicitante string   `xml:"RfcSolicitante,attr"`
}

//getDigestValueXML generate a well formed xml with the struct given
func getDigestValueXML(digestValueData *digestValue, url string) ([]byte, error) {
	content := &containerDigestValue{
		URL:    url,
		Digest: digestValueData,
	}
	return xml.Marshal(content)
}
