package helpers

import (
	"bufio"
	"bytes"
	"encoding/json"
	"log"
	"loginpage/globals"
	"net/http"
	"strings"
)

// Manda o Usuario para validar na API de LOGON
func CheckUserPass(username, password string) bool {
	requestData := map[string]string{
		"nome":  username,
		"senha": password,
	}

	//Marshal the Json data to bytes
	requestDataBytes, err := json.Marshal(requestData)
	if err != nil {
		log.Println(err)
		return false
	}

	// create a new http request la para a api de logon
	requestBody, err := http.NewRequest("POST", "http://srcdymw896.execute-api.us-east-1.amazonaws.com/api-login/logon", bytes.NewBuffer(requestDataBytes))
	if err != nil {
		log.Println(err)
		return false
	}

	//set the Content-Type header to application/json
	requestBody.Header.Set("Content-Type", "application/json")

	//send the request using the default Http Client
	resp, err := http.DefaultClient.Do(requestBody)
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
	url := "http://srcdymw896.execute-api.us-east-1.amazonaws.com/api-login/signin"
	req, err := http.NewRequest("POST", url, bytes.NewBuffer(colaboradorJ))
	if err != nil {
		log.Println("Error creating request:", err)
		return
	}
	req.Header.Set("Content-Type", "application/json")

	// send the request and print the response
	client := http.DefaultClient
	resp, err := client.Do(req)
	if err != nil {
		log.Println("Error sending request:", err)
		return
	}
	defer resp.Body.Close()
}

func BatePonto(nome string) {
	body := []byte(``)
	url := "http://vqief2ixwg.execute-api.us-east-1.amazonaws.com/api-ponto/ponto/" + nome
	req, err := http.NewRequest("POST", url, bytes.NewBuffer(body))
	if err != nil {
		log.Println("Error creating request:", err)
		return
	}
	req.Header.Set("Content-Type", "application/json")

	// send the request
	client := http.DefaultClient
	resp, err := client.Do(req)
	if err != nil {
		log.Println("Error sending request:", err)
		return
	}
	defer resp.Body.Close()
}

// Pegar os ultimos pontos
func UltimosPontos() []globals.Ponto {
	url := "http://vqief2ixwg.execute-api.us-east-1.amazonaws.com/api-ponto/pontos"
	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		log.Println("Error creating request:", err)
		return []globals.Ponto{}
	}
	req.Header.Set("Content-Type", "application/json")
	// send the request
	client := http.DefaultClient
	resp, err := client.Do(req)
	if err != nil {
		log.Println("Error sending request:", err)
		return []globals.Ponto{}
	}
	defer resp.Body.Close()
	body := bufio.NewScanner(resp.Body)
	var buffer []byte
	for body.Scan() {
		buffer = append(buffer, body.Bytes()...)
	}
	data := make([]globals.Ponto, 3)
	err = json.Unmarshal(buffer, &data)
	if err != nil {
		log.Println("1025-Error unmarhal request:", err)
		return []globals.Ponto{}
	}
	return data
}
