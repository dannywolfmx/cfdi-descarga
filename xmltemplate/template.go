package xmltemplate

import (
	"io"

	_ "embed"

	"github.com/dannywolfmx/cfdi-descarga/util"
)

//go:embed auth.xml
var templateAuthPath string

//go:embed download.xml
var templateDownloadPath string

//go:embed verify.xml
var templateVerifyPath string

//go:embed request.xml
var templateRequestPath string

// TemplateAuth return a filled xml with the params data
func TemplateAuth(buff io.Writer, data interface{}) error {
	return util.NewXML(buff, data, templateAuthPath)
}

// TemplateDownload return a filled xml with the params data
func TemplateDownload(buff io.Writer, data interface{}) error {
	return util.NewXML(buff, data, templateDownloadPath)
}

// TemplateRequest return a filled xml with the params data
func TemplateRequest(buff io.Writer, data interface{}) error {
	return util.NewXML(buff, data, templateRequestPath)
}

// TemplateVerify return a filled xml with the params data
func TemplateVerify(buff io.Writer, data interface{}) error {
	return util.NewXML(buff, data, templateVerifyPath)
}
