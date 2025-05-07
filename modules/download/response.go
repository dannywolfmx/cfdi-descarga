package download

import "encoding/xml"

//Struct to get the base64package into the xml
type response struct {
	EncodeFile string `xml:"Body>RespuestaDescargaMasivaTercerosSalida>Paquete"`
}

//ExtractResponse get the token from the xml
func ExtractResponse(data []byte) (string, error) {
	t := response{}
	err := xml.Unmarshal(data, &t)
	if err != nil {
		return "", err
	}

	return t.EncodeFile, nil
}
