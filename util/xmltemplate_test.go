package util

import (
	"bytes"
	"testing"
)

func TestNewXML(t *testing.T) {
	filePath := "../xmltemplate/test.xml"

	testXML := "<test>test</test>"
	testStruct := struct {
		Data string
	}{
		"test",
	}

	var buff bytes.Buffer

	err := NewXML(&buff, testStruct, filePath)

	if err != nil {
		t.Fatalf("Should get a xml. Err: %s", err)
	}

	//Remove EOF in unix like and windows
	value := buff.String()

	if value != testXML {
		t.Fatalf("Should be	the same xml and testXML. xml: %s, testxml: %s", value, testXML)
	}
}
