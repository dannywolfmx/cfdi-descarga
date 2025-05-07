package certificate

import "testing"

func TestGenerateSATIssuer(t *testing.T) {
	expectetResult := `CN=A.C. 2 de pruebas(4096),O=Servicio de Administración Tributaria,OU=Administración de Seguridad de la Información,emailAddress=asisnet@pruebas.sat.gob.mx,street=Av. Hidalgo 77, Col. Guerrero,postalCode=06300,C=MX,ST=Distrito Federal,L=Coyoacán,x500UniqueIdentifier=SAT970701NN3,unstructuredName=Responsable: ACDMA`
	cer, err := GetCertificate("../../test.cer")
	if err != nil {
		t.Fatalf("Should get a certificate: %s", err)
	}
	res := cer.ExtractIssuer()
	if expectetResult != res {
		t.Fatalf("Result diferent to expected, result: \n%s\n expected: \n%s\n", res, expectetResult)
	}
}
