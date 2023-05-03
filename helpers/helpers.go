package helpers

import (
	"bufio"
	"bytes"
	"crypto/tls"
	"crypto/x509"
	"encoding/json"
	"log"
	"loginpage/certificate"
	"loginpage/globals"
	"net/http"
	"strings"
)

// Create a custom HTTP client with a custom Transport object
func customHTTPClient() *http.Client {
	// Create a certificate pool
	certPool := x509.NewCertPool()

	// Add the client certificate to the certificate pool
	if ok := certPool.AppendCertsFromPEM([]byte(certificate.ClientCertificate)); !ok {
		log.Fatalf("Failed to add client certificate to certificate pool")
	}

	// Configure the custom HTTP client with the certificate pool
	tlsConfig := &tls.Config{
		RootCAs:            certPool,
		InsecureSkipVerify: true,
	}
	transport := &http.Transport{TLSClientConfig: tlsConfig}
	return &http.Client{Transport: transport}
}

// Manda o Usuario para validar na API de LOGON
func CheckUserPass(username, password string) bool {
	requestData := map[string]string{
		"nome":  username,
		"senha": password,
	}

	// Marshal the JSON data to bytes
	requestDataBytes, err := json.Marshal(requestData)
	if err != nil {
		log.Println(err)
		return false
	}

	// Create a new HTTP request to the API
	requestBody, err := http.NewRequest("POST", certificate.URLAPILogon, bytes.NewBuffer(requestDataBytes))
	if err != nil {
		log.Println(err)
		return false
	}

	// Set the Content-Type header to application/json
	requestBody.Header.Set("Content-Type", "application/json")

	// Send the request using the custom HTTP client
	customClient := customHTTPClient()
	resp, err := customClient.Do(requestBody)
	if err != nil {
		log.Println(err)
		return false
	}
	defer resp.Body.Close()

	var response string
	err = json.NewDecoder(resp.Body).Decode(&response)
	if err != nil {
		log.Println(err)
		return false
	}
	if response == "Authorized" {
		return true
	}
	log.Println(response)
	return false
}

func EmptyUserPass(username, password string) bool {
	return strings.Trim(username, " ") == "" || strings.Trim(password, " ") == ""
}

func Cadastro(nome, cpf, datanascimento, nomecompleto, password string) {
	colaborador := map[string]interface{}{
		"nome":            nome,
		"cpf":             cpf,
		"data-nascimento": datanascimento,
		"nome-completo":   nomecompleto,
		"senha":           password,
	}
	colaboradorJ, err := json.Marshal(colaborador)
	if err != nil {
		log.Println("Error marshaling colaborador:", err)
		return
	}

	// create a POST request with the JSON payload
	req, err := http.NewRequest("POST", certificate.URLAPISignin, bytes.NewBuffer(colaboradorJ))
	if err != nil {
		log.Println("Error creating request:", err)
		return
	}
	req.Header.Set("Content-Type", "application/json")

	// send the request and print the response
	client := customHTTPClient()
	resp, err := client.Do(req)
	if err != nil {
		log.Println("Error sending request:", err)
		return
	}
	defer resp.Body.Close()
}

func BatePonto(nome string) {
	body := []byte(``)
	url := certificate.URLAPIPonto + nome
	req, err := http.NewRequest("POST", url, bytes.NewBuffer(body))
	if err != nil {
		log.Println("Error creating request:", err)
		return
	}
	req.Header.Set("Content-Type", "application/json")

	// send the request
	client := customHTTPClient()
	resp, err := client.Do(req)
	if err != nil {
		log.Println("Error sending request:", err)
		return
	}
	defer resp.Body.Close()
}

// Pegar os ultimos pontos
func UltimosPontos() []globals.Ponto {
	req, err := http.NewRequest("GET", certificate.URLAPIPontos, nil)
	if err != nil {
		log.Println("Error creating request:", err)
		return []globals.Ponto{}
	}
	req.Header.Set("Content-Type", "application/json")

	// send the request
	client := customHTTPClient()
	resp, err := client.Do(req)
	if err != nil {
		log.Println("Error sending request:", err)
		return []globals.Ponto{}
	}
	defer resp.Body.Close()

	// read the response body and close it as well
	body := bufio.NewScanner(resp.Body)
	var buffer []byte
	for body.Scan() {
		buffer = append(buffer, body.Bytes()...)
	}
	pontos := make([]globals.Ponto, 3)
	err = json.Unmarshal(buffer, &pontos)

	if err != nil {
		log.Println("1025-Error unmarhal request:", err)
		return []globals.Ponto{}
	}
	return pontos
}
