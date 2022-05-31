package controllers

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"path/filepath"

	"github.com/gnana-prakash55/vault-cli/utils"
)

// struct for new repo
type NewRepo struct {
	Repo string `json:"repoName"`
}

type ConfigFile struct {
	RepoName string `json:"reponame"`
}

// create new repo for user
func CreateRepo(repoName string) {

	fmt.Println("Creating new repo...")

	repo := NewRepo{repoName}

	jsonRes, err := json.Marshal(repo)

	if err != nil {
		log.Fatalln(err)
	}

	jsonFile, err := os.Open(filepath.Join(".vault", "credentials", "secret.json"))

	if err != nil {
		log.Fatalln(err)
	}

	defer jsonFile.Close()

	value, err := ioutil.ReadAll(jsonFile)

	if err != nil {
		log.Fatalln(err)
	}

	var token Token

	json.Unmarshal(value, &token)

	req, err := http.NewRequest("POST", utils.GoDotEnvVariable("URL")+"/repo/create", bytes.NewBuffer(jsonRes))

	if err != nil {
		log.Fatalln(err)
	}

	req.Header.Add("Content-Type", "application/json")
	req.Header.Add("Authorization", "Bearer "+token.Token)

	client := &http.Client{}

	res, err := client.Do(req)

	if err != nil {
		log.Fatalln(err)
	}

	defer res.Body.Close()

	con, err := ioutil.ReadAll(res.Body)

	if err != nil {
		log.Fatalln(err)
	}

	bodyString := string(con)

	if bodyString == "repo created" {

		config := ConfigFile{repoName}

		jsonConfig, err := json.Marshal(config)

		if err != nil {
			log.Fatalln(err)
		}

		err = ioutil.WriteFile(".vault/config.json", jsonConfig, 0644)

		if err != nil {
			log.Fatalln(err)
		}

	}

	fmt.Println(bodyString)

}
