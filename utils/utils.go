package utils

import (
	"archive/zip"
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
	"strings"
)

const URL = "http://localhost:3000/repo/put"
const GET_URL = "http://localhost:3000/repo/get"

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

func GetFiles(token string) {

	config, err := ReadConfig()

	body, err := json.Marshal(config)

	if err != nil {
		panic(err)
	}

	if err != nil {
		panic(err)
	}

	request, err := http.NewRequest("POST", GET_URL, bytes.NewBuffer(body))

	if err != nil {
		panic(err)
	}

	request.Header.Add("Content-Type", "application/json")
	request.Header.Add("Authorization", "Bearer "+token)
	client := &http.Client{}

	response, err := client.Do(request)

	if err != nil {
		panic(err)
	}

	log.Println("Logging in...")

	defer response.Body.Close()

	content, err := ioutil.ReadAll(response.Body)

	if err != nil {
		panic(err)
	}

	err = Unzip(content, "content")

	if err != nil {
		panic(err)
	}
}

func Unzip(src []byte, dest string) error {
	r, err := zip.NewReader(bytes.NewReader(src), int64(len(src)))
	if err != nil {
		return err
	}

	os.MkdirAll(dest, 0755)

	// Closure to address file descriptors issue with all the deferred .Close() methods
	extractAndWriteFile := func(f *zip.File) error {
		rc, err := f.Open()
		if err != nil {
			return err
		}
		defer func() {
			if err := rc.Close(); err != nil {
				panic(err)
			}
		}()

		path := filepath.Join(dest, f.Name)

		// Check for ZipSlip (Directory traversal)
		if !strings.HasPrefix(path, filepath.Clean(dest)+string(os.PathSeparator)) {
			return fmt.Errorf("illegal file path: %s", path)
		}

		if f.FileInfo().IsDir() {
			os.MkdirAll(path, f.Mode())
		} else {
			os.MkdirAll(filepath.Dir(path), f.Mode())
			f, err := os.OpenFile(path, os.O_WRONLY|os.O_CREATE|os.O_TRUNC, f.Mode())
			if err != nil {
				return err
			}
			defer func() {
				if err := f.Close(); err != nil {
					panic(err)
				}
			}()

			_, err = io.Copy(f, rc)
			if err != nil {
				return err
			}
		}
		return nil
	}

	for _, f := range r.File {
		err := extractAndWriteFile(f)
		if err != nil {
			return err
		}
	}

	return nil
}
