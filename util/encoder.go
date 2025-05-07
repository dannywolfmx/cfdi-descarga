package util

import "encoding/base64"

//EncodeBase64 get a message and return a encode slice of bytes
func EncodeBase64(message []byte) string {
	return base64.StdEncoding.EncodeToString(message)
}

//DecodeBase64 get a encode message and return a slice of bytes
func DecodeBase64(message string) ([]byte, error) {
	return base64.StdEncoding.DecodeString(message)
}
