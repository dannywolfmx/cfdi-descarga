package authentication

import "encoding/xml"

//Struct to get the token into the xml
type token struct {
	//The tag represent the node of the token
	Value string `xml:"Body>AutenticaResponse>AutenticaResult"`
}

//ExtractToken get the token from the xml
func ExtractResponse(data []byte) (string, error) {
	t := token{}
	err := xml.Unmarshal(data, &t)
	if err != nil {
		return "", err
	}

	return t.Value, nil
}
