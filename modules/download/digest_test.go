package download

import (
	"testing"
)

//TestGetDigestValueWithOptionals test the returned value
//with optional "RFCEmisor" and "RFCReceptor"
func TestGetDigestValueWithOptionals(t *testing.T) {
	expectedResult := `<des:PeticionDescargaMasivaTercerosEntrada xmlns:des="http://DescargaMasivaTerceros.sat.gob.mx"><des:peticionDescarga IdPaquete="de7f98b6-d20f-44de-ab19-b312b489eec2_01" RfcSolicitante="AUAC4601138F9"></des:peticionDescarga></des:PeticionDescargaMasivaTercerosEntrada>`
	url := "http://DescargaMasivaTerceros.sat.gob.mx"
	digestValueTestData := &digestValue{
		RFCSolicitante: "AUAC4601138F9",
		IDPackage:      `de7f98b6-d20f-44de-ab19-b312b489eec2_01`,
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
