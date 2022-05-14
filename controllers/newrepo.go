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
)

// struct for new repo
type NewRepo struct {
	Repo string `json:"repoName"`
}

// create new repo for user
func CreateRepo(repoName string) {

	fmt.Println("Creating new repo...")

	repo := NewRepo{repoName}

	jsonRes, err := json.Marshal(repo)

	if err != nil {
		log.Fatalln(err)
	}

	jsonFile, err := os.Open(filepath.Join(".credentials", "secret.json"))

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

	req, err := http.NewRequest("POST", "http://13.232.12.225:3000/repo/create", bytes.NewBuffer(jsonRes))

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

	fmt.Println(bodyString)

}
