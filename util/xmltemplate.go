package util

import (
	"io"
	"strings"
	"text/template"
)

// NewXML read a xml template file and set the data struct into this to generate a xml
// You can use any io.Writer implementation as receiver type. Ejem.: bytes.Buffer
func NewXML(buff io.Writer, data interface{}, xmlTemplate string) error {
	xmlTemplate = strings.TrimSpace(xmlTemplate)

	tmpl, err := template.New("").Parse(xmlTemplate)

	if err != nil {
		return err
	}

	return tmpl.Execute(buff, data)
}
