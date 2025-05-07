package main

import (
	"fmt"
	"io/ioutil"
	"log"
	"time"

	"github.com/dannywolfmx/cfdi-descarga/certificate"
	"github.com/dannywolfmx/cfdi-descarga/digest"
	"github.com/dannywolfmx/cfdi-descarga/modules/download"
	"github.com/dannywolfmx/cfdi-descarga/modules/request"
	"github.com/dannywolfmx/cfdi-descarga/modules/verify"
	"github.com/dannywolfmx/cfdi-descarga/signature"
	"github.com/dannywolfmx/cfdi-descarga/util"
	"github.com/spf13/viper"
	// "github.com/dannywolfmx/cfdi-descarga/util" // Assuming util might be needed later
)

const (
	cerPath     = "./files/cer.cer"     // Asegúrate de que este archivo exista en la raíz del proyecto
	keyPath     = "./files/key.key"     // Asegúrate de que este archivo exista en la raíz del proyecto
	initialDate = "2025-01-01T00:00:00" // Fecha de inicio para la solicitud
	endDate     = "2025-01-10T23:59:59" // Fecha de fin para la solicitud
	requestType = "Metadata"            // Puede ser "CFDI" o "Metadata"
)

func main() {
	viper.SetConfigFile(".env")

	if err := viper.ReadInConfig(); err != nil {
		log.Fatalf("Error al leer el archivo de configuración: %s", err)
	}

	//Read the password from the enviroment variable
	cerPassword := viper.GetString("PASSWORD")

	// Check if the password is set
	if cerPassword == "" {
		log.Fatal("Error: La variable de entorno PASSWORD no está configurada.")
	}

	//RFC
	rfc := viper.GetString("RFC")
	if rfc == "" {
		log.Fatal("Error: La variable de entorno RFC no está configurada.")
	}

	log.Println("Iniciando proceso de descarga de CFDI...")

	// --- 1. Cargar Certificado y Llave ---
	log.Println("Cargando certificado y llave...")
	cer, err := certificate.GetCertificate(cerPath)
	if err != nil {
		log.Fatalf("Error al cargar certificado desde %s: %v", cerPath, err)
	}

	key, err := signature.ExtractKey(keyPath, cerPassword)
	if err != nil {
		log.Fatalf("Error al cargar llave desde %s: %v", keyPath, err)
	}
	log.Println("Certificado y llave cargados exitosamente.")

	// --- 2. Autenticación ---
	log.Println("Iniciando autenticación...")
	auth := digest.NewRequestAuth(cer.GetEncodeCertificate(), key, util.GenerateUUIDV4)

	token, err := auth.RequestToken()
	if err != nil {
		log.Fatalf("Error durante la autenticación: %v", err)
	}

	log.Println("Autenticación exitosa. Token obtenido.")

	// --- 3. Solicitud de Descarga ---
	log.Println("Iniciando solicitud de descarga...")
	// Puedes ajustar RFCReceptor y RFCEmisor según necesites
	// var rfcReceptor = rfc
	// var rfcEmisor = rfc
	requestData := request.RequestData{
		InitialDate:    initialDate,
		EndDate:        endDate,
		RFCSolicitante: rfc,
		// RFCReceptor:    &rfcReceptor, // Descomenta si necesitas filtrar por receptor
		RFCEmisor:   &rfc,
		RequestType: requestType,
		Cer:         cer,
		Key:         key,
		Token:       token,
	}

	idSolicitud, err := requestData.RequestIDSolicitud()

	if err != nil {
		log.Fatalf("Error al realizar la solicitud de descarga: %v", err)
	}

	log.Printf("Solicitud de descarga enviada. ID de Solicitud: %s", idSolicitud)

	// --- 4. Verificación de la Solicitud ---
	log.Println("Iniciando verificación periódica de la solicitud...")
	reqDataVerify := verify.RequestData{
		RFCSolicitante:        rfc,
		IDSolicitudDeDescarga: idSolicitud,
		Cer:                   cer,
		Key:                   key,
		// El token se actualizará en cada iteración del bucle
	}
	var packagesIDS []string
	const maxRetries = 20              // Número máximo de intentos de verificación
	const retryDelay = 1 * time.Minute // Tiempo de espera entre intentos de verificación

	for i := 0; i < maxRetries; i++ {
		log.Printf("Intento de verificación %d/%d...", i+1, maxRetries)

		// Re-autenticar para obtener un token fresco si es necesario (depende de la duración del token)
		// En este ejemplo, re-autenticamos en cada intento para asegurar un token válido
		log.Println("Re-autenticando para verificación...")

		token, err = auth.RequestToken()

		if err != nil {
			log.Printf("Error al re-autenticar para verificación: %v. Reintentando en %v...", err, retryDelay)
			time.Sleep(retryDelay)
			continue
		}

		reqDataVerify.Token = token
		log.Println("Re-autenticación exitosa.")

		respuestaVerify, err := reqDataVerify.SendRequest()
		if err != nil {
			log.Printf("Error durante la verificación: %v. Reintentando en %v...", err, retryDelay)
			time.Sleep(retryDelay)
			continue
		}

		if respuestaVerify.RequestCode == nil {
			log.Printf("Error: No se recibió código de estado en la verificación. Reintentando en %v...", retryDelay)
			time.Sleep(retryDelay)
			continue
		}

		log.Printf("Estado de la solicitud: %d (%s)", *respuestaVerify.RequestCode, *respuestaVerify.Message)
		log.Printf("Código de estatus: %d", *respuestaVerify.StatusCode)
		if respuestaVerify.NumCFDI != nil {
			log.Printf("Número de CFDIs: %d", *respuestaVerify.NumCFDI)
		}

		// Estados posibles (según documentación no oficial):
		// 0: Token inválido
		// 1: Aceptada
		// 2: En proceso
		// 3: Terminada (Lista para descargar)
		// 4: Error
		// 5: Rechazada
		// 6: Vencida

		switch *respuestaVerify.RequestCode {
		case 0: // Token inválido
			log.Fatalf("Token inválido. Error: %s", *respuestaVerify.Message)
		case 3: // Terminada
			log.Println("Solicitud terminada. Lista para descargar.")
			packagesIDS = respuestaVerify.PackagesIDS
			goto download // Salir del bucle de verificación
		case 1, 2: // Aceptada o En proceso
			log.Printf("La solicitud aún está en proceso. Esperando %v para la próxima verificación...", retryDelay)
			time.Sleep(retryDelay)
		case 4, 5, 6: // Error, Rechazada o Vencida
			log.Fatalf("La solicitud falló o expiró. Estado: %d, Mensaje: %s", *respuestaVerify.RequestCode, *respuestaVerify.Message)
		default:
			log.Printf("Estado de solicitud desconocido: %d. Reintentando en %v...", *respuestaVerify.RequestCode, retryDelay)
			time.Sleep(retryDelay)
		}
	}

	log.Fatal("Se superó el número máximo de intentos de verificación.")

download:
	if len(packagesIDS) == 0 {
		log.Fatal("No se encontraron paquetes para descargar.")
	}
	log.Printf("Se encontraron %d paquetes para descargar: %v", len(packagesIDS), packagesIDS)

	// --- 5. Descarga de CFDI ---
	log.Println("Iniciando descarga de paquetes...")
	reqDownloadData := download.RequestData{
		RFCSolicitante: rfc,
		Cer:            cer,
		Key:            key,
		// IDPackage y Token se asignarán dentro del bucle
	}

	for _, idPaquete := range packagesIDS {
		log.Printf("Descargando paquete: %s", idPaquete)

		// Re-autenticar para obtener un token fresco
		log.Println("Re-autenticando para descarga...")
		token, err = auth.RequestToken()
		if err != nil {
			log.Printf("Error al re-autenticar para descarga del paquete %s: %v. Omitiendo paquete.", idPaquete, err)
			continue // Saltar al siguiente paquete
		}
		reqDownloadData.Token = token
		reqDownloadData.IDPackage = idPaquete
		log.Println("Re-autenticación exitosa.")

		datosPaquete, err := reqDownloadData.SendRequest()
		if err != nil {
			log.Printf("Error al descargar el paquete %s: %v. Omitiendo paquete.", idPaquete, err)
			continue // Saltar al siguiente paquete
		}

		if len(datosPaquete) == 0 {
			log.Printf("Advertencia: El paquete %s descargado está vacío.", idPaquete)
			continue
		}

		// Guardar el paquete descargado en un archivo .zip
		fileName := fmt.Sprintf("paquete_%s.zip", idPaquete)
		err = ioutil.WriteFile(fileName, datosPaquete, 0644)
		if err != nil {
			log.Printf("Error al guardar el paquete %s en el archivo %s: %v", idPaquete, fileName, err)
		} else {
			log.Printf("Paquete %s guardado exitosamente como %s", idPaquete, fileName)
		}
	}

	log.Println("Proceso de descarga de CFDI completado.")
}
