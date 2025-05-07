package verify

import "encoding/xml"

// Struct to get the token into the xml
type response struct {
	Data ResponseData `xml:"Body>VerificaSolicitudDescargaResponse>VerificaSolicitudDescargaResult"`
}

// ResponseData will be filled with the server response message
// To prevent the Go Default values, we add pointes to get a nil, instant default value
type ResponseData struct {
	StatusCode  *int     `xml:"CodEstatus,attr"`
	RequestCode *int     `xml:"EstadoSolicitud,attr"`
	NumCFDI     *int     `xml:"NumeroCFDIs,attr"`
	Message     *string  `xml:"Mensaje,attr"`
	PackagesIDS []string `xml:"IdsPaquetes"`
}

// UnmarshalResponse get the response from the xml
func UnmarshalResponse(data []byte) (*ResponseData, error) {
	var t response
	if err := xml.Unmarshal(data, &t); err != nil {
		return nil, err
	}

	return &t.Data, nil
}
