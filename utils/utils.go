package utils

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"io/fs"
	"io/ioutil"
	"log"
	"mime/multipart"
	"net/http"
	"os"
	"path/filepath"
)

const URL = "http://localhost:3000/repo/upload"

type Token struct {
	Token string `json:"token"`
}

type Config struct {
	RepoName string `json:"repoName"`
}

func ReadToken() (string, error) {

	jsonFile, err := os.Open(filepath.Join(".vault", "credentials", "secret.json"))
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

func ReadConfig() (Config, error) {
	jsonFile, err := os.Open(filepath.Join(".vault", "config.json"))
	if err != nil {
		return Config{}, err
	}

	configValue, _ := ioutil.ReadAll(jsonFile)
	var config Config

	err = json.Unmarshal(configValue, &config)

	return config, nil

}

func UploadFiles(path, token string) (string, error) {

	body := &bytes.Buffer{}
	writer := multipart.NewWriter(body)

	config, err := ReadConfig()

	writer.WriteField("repo", config.RepoName)

	err = filepath.Walk(path, func(path string, info fs.FileInfo, err error) error {
		if err != nil {
			fmt.Println(err)
			return err
		}

		// fmt.Printf("dir: %v: name: %s\n", info.IsDir(), path)

		if !info.IsDir() {

			file, err := os.Open(path)
			if err != nil {
				return err
			}

			defer file.Close()

			fmt.Println(file.Name())
			part, err := writer.CreateFormFile("file", filepath.Base(file.Name()))

			if err != nil {
				return err
			}

			io.Copy(part, file)

		}

		return nil
	})

	writer.Close()

	if err != nil {
		panic(err)
	}

	request, err := http.NewRequest("POST", URL, body)

	if err != nil {
		return "", err
	}

	request.Header.Add("Content-Type", writer.FormDataContentType())
	request.Header.Add("Authorization", "Bearer "+token)
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
