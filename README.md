# CFDI-Descarga

Biblioteca escrita en Go para descargar CFDIs de forma masiva desde el portal del SAT (Sistema de Administración Tributaria de México).

![Go Version](https://img.shields.io/badge/Go-1.20+-blue)


# Nota: 

El servidor del SAT va muy pero muy regular, por lo que aunque el proyecto funciona casi en todas las ocasiones que he hecho pruebas el servidor no funciona,
una alternativa más viable es no utilizar este webservice del SAT y mejor hacer scraping a la página que esa por algun motivo admite la descarga en segundos 
(Por el momento estoy descepcionado por el SAT asi que ocupo motivación para poder hacer este proyecto secundario de scraping)

Libero este proyecto con licencia MIT, por lo que sientanse libres de usarlo, modificarlo y compartirlo.

## Descripción

CFDI-Descarga permite automatizar la obtención de CFDIs (Comprobantes Fiscales Digitales por Internet) directamente desde los servidores del SAT, utilizando los servicios web oficiales. Esta herramienta soporta el flujo completo del proceso de descarga masiva:

1. **Autenticación** con certificados (.cer) y llaves (.key)
2. **Solicitud** de paquetes de CFDIs
3. **Verificación** del estado de la solicitud
4. **Descarga** de los archivos resultantes

## Requisitos previos

Para utilizar esta biblioteca necesitas:

- Go versión 1.20 o superior
- Certificado Digital (.cer)
- Llave privada (.key)
- Contraseña de la llave privada
- RFC del solicitante

## Instalación

```bash
go get github.com/dannywolfmx/cfdi-descarga
```

## Uso básico

En el Main hay un proyecto de ejemplo funcional, que les puede ser de utilidad para comprender como funciona la librería.

El proceso para descargar CFDIs consta de cuatro pasos secuenciales:

### 1. Autenticación

```go
// Cargar certificado y llave
cer, err := certificate.GetCertificate("ruta/al/certificado.cer")
if err != nil {
    log.Fatalf("Error al cargar certificado: %v", err)
}

key, err := signature.ExtractKey("ruta/a/llave.key", "contraseña")
if err != nil {
    log.Fatalf("Error al cargar llave: %v", err)
}

// Crear request de autenticación
reqAuth := authentication.RequestData{
    EncodedCertificate: cer.GetEncodeCertificate(),
    Key: key,
}

// Obtener token de autenticación
token, err := reqAuth.SendRequest()
if err != nil {
    log.Fatalf("Error en autenticación: %v", err)
}
```

### 2. Solicitud de CFDIs

```go
// Preparar solicitud
requestData := request.RequestData{
    InitialDate:    "2023-01-01T00:00:00",
    EndDate:        "2023-01-31T23:59:59",
    RFCSolicitante: "XXXX010101XXX",
    RequestType:    "CFDI",  // Puede ser "CFDI" o "Metadata"
    Cer:            cer,
    Key:            key,
    Token:          token,
}

// Enviar solicitud
idSolicitud, err := requestData.SendRequest()
if err != nil {
    log.Fatalf("Error al solicitar CFDIs: %v", err)
}
```

### 3. Verificación del estado

```go
// Preparar verificación
reqDataVerify := verify.RequestData{
    RFCSolicitante: "XXXX010101XXX",
    IDRequest:      idSolicitud,
    Cer:            cer,
    Key:            key,
    Token:          token,
}

// Verificar estado
respuesta, err := reqDataVerify.SendRequest()
if err != nil {
    log.Fatalf("Error al verificar solicitud: %v", err)
}

// Revisar el código de estado
// *respuesta.RequestCode == 3 indica solicitud lista para descargar
```

### 4. Descarga de paquetes

```go
// Preparar descarga
reqDownloadData := download.RequestData{
    RFCSolicitante: "XXXX010101XXX",
    IDPackage:      idPaquete, // Obtenido de respuesta.PackagesIDS
    Cer:            cer,
    Key:            key,
    Token:          token,
}

// Descargar paquete
datos, err := reqDownloadData.SendRequest()
if err != nil {
    log.Fatalf("Error al descargar paquete: %v", err)
}

// Guardar archivo
err = ioutil.WriteFile("cfdi_paquete.zip", datos, 0644)
if err != nil {
    log.Fatalf("Error al guardar archivo: %v", err)
}
```

## Documentación de referencia

- [Documentación oficial del SAT sobre Descarga Masiva](https://www.sat.gob.mx/cs/Satellite?blobcol=urldata&blobkey=id&blobtable=MungoBlobs&blobwhere=1579314915253&ssbinary=true)
- [Guía de consumo del Web Service (SW)](https://developers.sw.com.mx/knowledge-base/consumo-webservice-descarga-masiva-sat/)

## Recursos para pruebas

- [CSD de prueba para personas morales y físicas](http://omawww.sat.gob.mx/tramitesyservicios/Paginas/documentos/RFC-PAC-SC.zip)
- [Certificados de prueba vigentes](https://developers.sw.com.mx/knowledge-base/donde-encuentro-csd-de-pruebas-vigentes/)

- [Guia documentada Solicitud de descarga masiva](https://www.sat.gob.mx/cs/Satellite?blobcol=urldata&blobkey=id&blobtable=MungoBlobs&blobwhere=1461175195160&ssbinary=true)

