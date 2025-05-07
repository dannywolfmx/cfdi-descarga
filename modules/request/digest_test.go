package request

import (
	"testing"
)

//TestGetDigestValueWithOptionals test the returned value
//with optional "RFCEmisor" and "RFCReceptor"
func TestGetDigestValueWithOptionals(t *testing.T) {
	expectedResult := `<des:SolicitaDescarga xmlns:des="http://DescargaMasivaTerceros.sat.gob.mx"><des:solicitud RfcEmisor="ABE" RfcReceptor="ABD" RfcSolicitante="ABC" FechaInicial="2018-10-10T00:00:00" FechaFinal="2018-10-20T00:00:00" TipoSolicitud="CFDI"></des:solicitud></des:SolicitaDescarga>`
	url := "http://DescargaMasivaTerceros.sat.gob.mx"
	rfcReceptor := "ABD"
	rfcEmisor := "ABE"
	digestValueTestData := &digestValue{
		RFCSolicitante: "ABC",
		RFCReceptor:    &rfcReceptor,
		RFCEmisor:      &rfcEmisor,
		FechaInicial:   "2018-10-10T00:00:00",
		FechaFinal:     "2018-10-20T00:00:00",
		TipoSolicitud:  "CFDI",
	}
	value, err := getDigestValueXML(digestValueTestData, url)

	if err != nil {
		t.Fatalf("Should get a digest value: %s", err)
	}
	result := string(value)
	//Check if the result is correct
	if expectedResult != result {
		t.Fatalf("expected and the result does not match, expected: %s, result:%s", expectedResult, result)
	}
}

//TestGetDigestValueEmptyData test the returned value with empty data
func TestGetDigestValueEmptyData(t *testing.T) {
	expectedResult := `<des:SolicitaDescarga xmlns:des="http://DescargaMasivaTerceros.sat.gob.mx"><des:solicitud RfcSolicitante="" FechaInicial="" FechaFinal="" TipoSolicitud=""></des:solicitud></des:SolicitaDescarga>`
	digestValueTestData := &digestValue{}
	url := "http://DescargaMasivaTerceros.sat.gob.mx"
	value, err := getDigestValueXML(digestValueTestData, url)

	if err != nil {
		t.Fatalf("Should get a digest value: %s", err)
	}
	result := string(value)
	//Check if the result is correct
	if expectedResult != result {
		t.Fatalf("expected and the result does not match, expected: %s, result:%s", expectedResult, result)
	}
}

//TestGetDigestValueWithoutOptionals test the returned
//value without optional "RFCEmisor" and "RFCReceptor"
func TestGetDigestValueWithoutOptionals(t *testing.T) {
	expectedResult := `<des:SolicitaDescarga xmlns:des="http://DescargaMasivaTerceros.sat.gob.mx"><des:solicitud RfcEmisor="AUAC4601138F9" RfcReceptor="AUAC4601138F9" RfcSolicitante="AUAC4601138F9" FechaInicial="2018-10-10T00:00:00" FechaFinal="2018-10-20T00:00:00" TipoSolicitud="CFDI"></des:solicitud></des:SolicitaDescarga>`
	url := "http://DescargaMasivaTerceros.sat.gob.mx"
	rfc := "AUAC4601138F9"
	digestValueTestData := &digestValue{
		RFCEmisor:      &rfc,
		RFCReceptor:    &rfc,
		RFCSolicitante: rfc,
		FechaInicial:   "2018-10-10T00:00:00",
		FechaFinal:     "2018-10-20T00:00:00",
		TipoSolicitud:  "CFDI",
	}

	value, err := getDigestValueXML(digestValueTestData, url)

	if err != nil {
		t.Fatalf("Should get a digest value: %s", err)
	}
	result := string(value)
	//Check if the result is correct
	if expectedResult != result {
		t.Fatalf("expected and the result does not match, expected: %s, result:%s", expectedResult, result)
	}
}
