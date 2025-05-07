package request

import "encoding/xml"

//Struct to get the token into the xml
type response struct {
	Response idSolicitud `xml:"Body>SolicitaDescargaResponse>SolicitaDescargaResult"`
}

type idSolicitud struct {
	ID string `xml:"IdSolicitud,attr"`
}

//ExtractResponse get the response from the xml
func ExtractResponse(data []byte) (string, error) {
	t := response{}
	err := xml.Unmarshal(data, &t)
	if err != nil {
		return "", err
	}

	return t.Response.ID, nil
}
