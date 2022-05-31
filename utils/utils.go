package utils

import (
	"archive/zip"
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"io/ioutil"
	"log"
	"mime/multipart"
	"net/http"
	"os"
	"path/filepath"
	"strings"

	"github.com/joho/godotenv"
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

	if err != nil {
		panic(err)
	}

	return config, nil

}

func UploadFiles(path, token string) (string, error) {

	err := os.MkdirAll("../compressed", os.ModePerm)

	if err != nil {
		return "", err
	}

	config, err := ReadConfig()

	if err != nil {
		return "", err
	}

	err = RecursiveZip(path, filepath.Join("../compressed", config.RepoName+".zip"))

	if err != nil {
		return "", err
	}

	// return "success", nil

	body := &bytes.Buffer{}

	writer := multipart.NewWriter(body)

	writer.WriteField("repo", config.RepoName)

	part, err := writer.CreateFormFile("file", config.RepoName+".zip")

	if err != nil {
		return "", err
	}

	file, err := os.Open(filepath.Join("../compressed", config.RepoName+".zip"))

	io.Copy(part, file)

	if err != nil {
		panic(err)
	}

	writer.Close()

	request, err := http.NewRequest("POST", GoDotEnvVariable("URL")+"/repo/put", body)

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

	err = os.RemoveAll("../compressed")
	if err != nil {
		return "", err
	}

	return string(content), nil

}

func GetFiles(token string) {

	config, err := ReadConfig()

	if err != nil {
		panic(err)
	}

	body, err := json.Marshal(config)

	if err != nil {
		panic(err)
	}

	if err != nil {
		panic(err)
	}

	request, err := http.NewRequest("POST", GoDotEnvVariable("URL")+"/repo/get", bytes.NewBuffer(body))

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

	log.Println("Unzipping files...")

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
			os.MkdirAll(path, os.ModePerm)
		} else {
			os.MkdirAll(filepath.Dir(path), os.ModePerm)
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

func RecursiveZip(pathToZip, destinationPath string) error {
	destinationFile, err := os.Create(destinationPath)
	if err != nil {
		return err
	}
	myZip := zip.NewWriter(destinationFile)
	fmt.Println("Creating Writer!!!")
	err = filepath.Walk(pathToZip, func(filePath string, info os.FileInfo, err error) error {
		if info.IsDir() {
			return nil
		}
		if err != nil {
			return err
		}
		relPath := strings.TrimPrefix(filePath, filepath.Dir(pathToZip))
		zipFile, err := myZip.Create(relPath)
		if err != nil {
			return err
		}
		fsFile, err := os.Open(filePath)
		if err != nil {
			return err
		}

		fmt.Println(filePath)

		_, err = io.Copy(zipFile, fsFile)

		if err != nil {
			return err
		}
		return nil
	})
	if err != nil {
		return err
	}
	err = myZip.Close()
	if err != nil {
		return err
	}

	return nil
}

func GoDotEnvVariable(key string) string {

	// load .env file
	err := godotenv.Load(".env")

	if err != nil {
		log.Fatalf("Error loading .env file")
	}

	return os.Getenv(key)
}
