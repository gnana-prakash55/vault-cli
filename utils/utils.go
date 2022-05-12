package utils

import (
	"bytes"
	"encoding/json"
	"io"
	"io/ioutil"
	"log"
	"mime/multipart"
	"net/http"
	"os"
	"path/filepath"
)

const URL = "http://localhost:8888/submit"

type Token struct {
	Token string `json:"token"`
}

func ReadToken() (string, error) {

	jsonFile, err := os.Open(filepath.Join(".credentials", "secret.json"))
	if err != nil {
		return "", err
	}

	log.Println("Opening credentials....")

	defer jsonFile.Close()

	value, _ := ioutil.ReadAll(jsonFile)

	var token Token

	json.Unmarshal(value, &token)

	return token.Token, nil

}

func UploadFiles(filename, token string) (string, error) {
	file, err := os.Open(filename)

	if err != nil {
		return "", err
	}

	defer file.Close()

	body := &bytes.Buffer{}
	writer := multipart.NewWriter(body)
	part, err := writer.CreateFormFile("fileUploadName", filepath.Base(file.Name()))

	if err != nil {
		return "", err
	}

	io.Copy(part, file)
	writer.Close()

	request, err := http.NewRequest("POST", URL, body)

	if err != nil {
		return "", err
	}

	request.Header.Add("Content-Type", writer.FormDataContentType())
	request.Header.Add("Authorizaton", "Bearer "+token)
	client := &http.Client{}

	response, err := client.Do(request)

	if err != nil {
		return "", err
	}

	log.Println("Logging in...")

	defer response.Body.Close()

	content, err := ioutil.ReadAll(response.Body)

	if err != nil {
		return "", err
	}

	return string(content), nil

}
