package download

import "testing"

//TestExtractResponse will test if the program can extract the "IdSolicitado" attribute
func TestExtractResponse(t *testing.T) {
	xmlTest := `<s:Envelope xmlns:s="http://schemas.xmlsoap.org/soap/envelope/">
	<s:Header>
		<h:respuesta CodEstatus="5000" Mensaje="Solicitud Aceptada" xmlns="http://DescargaMasivaTerceros.sat.gob.mx" xmlns:h="http://DescargaMasivaTerceros.sat.gob.mx" xmlns:xsd="http://www.w3.org/2001/XMLSchema" xmlns:xsi="http://www.w3.org/2001/XMLSchema-instance"/>
	</s:Header>
	<s:Body xmlns:xsd="http://www.w3.org/2001/XMLSchema" xmlns:xsi="http://www.w3.org/2001/XMLSchema-instance">
		<RespuestaDescargaMasivaTercerosSalida xmlns="http://DescargaMasivaTerceros.sat.gob.mx">
			<Paquete>1234</Paquete>
		</RespuestaDescargaMasivaTercerosSalida>
	</s:Body>
</s:Envelope>`

	value, err := ExtractResponse([]byte(xmlTest))

	if err != nil {
		t.Fatalf("Should get extract response %s", err)
	}

	expectedResult := "1234"
	if value != expectedResult {
		t.Fatalf("Unexpected result result: %s, value: %s ", expectedResult, value)
	}
}
