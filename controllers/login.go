package controllers

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"os"

	"github.com/gnana-prakash55/vault-cli/utils"
)

// struct for user login
type LoginUser struct {
	Email    string `json:"email"`
	Password string `json:"password"`
}

type Token struct {
	Token string `json:"token"`
}

// login user
func Login(email string, password string) {

	fmt.Println("Login user...")

	login := LoginUser{email, password}

	jsonRes, err := json.Marshal(login)

	if err != nil {
		log.Fatalln(err)
	}

	res, err := http.Post(utils.GoDotEnvVariable("URL")+"/login", "application/json", bytes.NewBuffer(jsonRes))

	if err != nil {
		log.Fatalln(err)
	}

	defer res.Body.Close()

	bodyBytes, _ := ioutil.ReadAll(res.Body)

	bodyString := string(bodyBytes)

	if bodyString == "wrong username or password" {
		log.Fatalln(bodyString)
	}

	token := Token{bodyString}

	jsonToken, err := json.Marshal(token)

	if err != nil {
		log.Fatalln(err)
	}

	err = os.MkdirAll(".vault/credentials", os.ModePerm)

	if err != nil {
		log.Fatalln(err)
	}

	err = ioutil.WriteFile(".vault/credentials/secret.json", jsonToken, 0644)

	if err != nil {
		log.Fatalln(err)
	}

	fmt.Println("User loged in")

}
