package verify

import (
	"testing"
)

//TestGetDigestValueWithOptionals test the returned value
//with optional "RFCEmisor" and "RFCReceptor"
func TestGetDigestValueWithOptionals(t *testing.T) {
	expectedResult := `<des:VerificaSolicitudDescarga xmlns:des="http://DescargaMasivaTerceros.sat.gob.mx"><des:solicitud IdSolicitud="de7f98b6-d20f-44de-ab19-b312b489eec2" RfcSolicitante="AUAC4601138F9"></des:solicitud></des:VerificaSolicitudDescarga>`
	url := "http://DescargaMasivaTerceros.sat.gob.mx"
	rfcSolicitante := "AUAC4601138F9"
	idSolicitud := `de7f98b6-d20f-44de-ab19-b312b489eec2`
	value, err := getDigestValueXML(rfcSolicitante, idSolicitud, url)

	if err != nil {
		t.Fatalf("Should get a digest value: %s", err)
	}
	result := string(value)
	//Check if the result is correct
	if expectedResult != result {
		t.Fatalf("expected and the result does not match, expected: %s, result:%s", expectedResult, result)
	}
}
