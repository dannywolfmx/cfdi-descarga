package verify

import "testing"

// TestExtractResponse will test if the program can extract the "IdSolicitado" attribute
func TestExtractResponse(t *testing.T) {
	xmlTest := `<s:Envelope xmlns:s="http://schemas.xmlsoap.org/soap/envelope/">
    <s:Body xmlns:xsi="http://www.w3.org/2001/XMLSchema-instance" xmlns:xsd="http://www.w3.org/2001/XMLSchema">
        <VerificaSolicitudDescargaResponse xmlns="http://DescargaMasivaTerceros.sat.gob.mx">
            <VerificaSolicitudDescargaResult CodEstatus="5000" EstadoSolicitud="3" CodigoEstadoSolicitud="5000" NumeroCFDIs="121" Mensaje="Solicitud Aceptada">
                <IdsPaquetes>de7f98b6-d20f-44de-ab19-b312b489eec2_01</IdsPaquetes>
            </VerificaSolicitudDescargaResult>
        </VerificaSolicitudDescargaResponse>
    </s:Body>
</s:Envelope>`

	value, err := UnmarshalResponse([]byte(xmlTest))

	if err != nil {
		t.Fatalf("Should get extract response %s", err)
	}

	statusCode := 5000
	if statusCode != *value.StatusCode {
		t.Fatalf("Unexpected status code, expected %d, result %d", statusCode, *value.StatusCode)
	}

	reqCode := 3
	if reqCode != *value.RequestCode {
		t.Fatalf("Unexpected status code, expected %d, result %d", statusCode, *value.StatusCode)
	}

	numCFDI := 121
	if numCFDI != *value.NumCFDI {
		t.Fatalf("Unexpected NumCFDI, expected %d, result %d", numCFDI, *value.NumCFDI)
	}

	message := "Solicitud Aceptada"
	if message != *value.Message {
		t.Fatalf("Unexpected message, expected %s, result %s", message, *value.Message)
	}
}
