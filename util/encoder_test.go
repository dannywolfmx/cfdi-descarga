package util

import (
	"bytes"
	"encoding/base64"
	"testing"
)

func TestEncodeBase64(t *testing.T) {
	testMessage := []byte("Hello")

	messageBase64 := EncodeBase64(testMessage)

	decodeMessage, err := base64.StdEncoding.DecodeString(messageBase64)
	if err != nil {
		t.Fatalf("Should get a decoded message: %s", err)
	}

	if bytes.Compare(decodeMessage, testMessage) != 0 {
		t.Fatal("The decoded message dont math to the original message")
	}
}
