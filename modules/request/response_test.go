package request

import "testing"

//TestExtractResponse will test if the program can extract the "IdSolicitado" attribute
func TestExtractResponse(t *testing.T) {
	xmlTest := `<s:Envelope xmlns:s="http://schemas.xmlsoap.org/soap/envelope/">
    	<s:Body xmlns:xsi="http://www.w3.org/2001/XMLSchema-instance" xmlns:xsd="http://www.w3.org/2001/XMLSchema">
    	    <SolicitaDescargaResponse xmlns="http://DescargaMasivaTerceros.sat.gob.mx">
    	        <SolicitaDescargaResult IdSolicitud="e6e54ed5-22fb-4d17-b03b-f6d80079b68d" CodEstatus="5000" Mensaje="Solicitud Aceptada"/>
    	    </SolicitaDescargaResponse>
    	</s:Body>
	</s:Envelope>`

	idTest := "e6e54ed5-22fb-4d17-b03b-f6d80079b68d"

	value, err := ExtractResponse([]byte(xmlTest))

	if err != nil {
		t.Fatalf("Should get extract response %s", err)
	}

	if value != idTest {
		t.Fatalf("idTest and value should be the same, idTest: %s, value: %s", idTest, value)
	}

}
