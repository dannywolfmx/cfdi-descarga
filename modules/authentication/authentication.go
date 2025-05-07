package authentication

import (
	"github.com/dannywolfmx/cfdi-descarga/digest"
)

// RequestData implement the RequestData interface to generate a request body
// This structs implement the client.RequestData
type RequestData struct {
	digest.RequestAuth
}
