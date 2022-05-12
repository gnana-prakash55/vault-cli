package controllers

import (
	"bytes"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
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

	// tokenFile, err := ioutil.ReadFile("/credentails/secret.json")

	req, err := http.NewRequest("POST", "http://localhost:3000/repo/create", bytes.NewBuffer(jsonRes))

	req.Header.Add("Content-Type", "application/json")

	// res, err := http.Post("http://localhost:3000/repo/create", "application/json", bytes.NewBuffer(jsonRes))

	// if err != nil {
	// 	log.Fatalln(err)
	// }

	// defer res.Body.Close()

	// bodyBytes, _ := ioutil.ReadAll(res.Body)

	// bodyString := string(bodyBytes)

}
