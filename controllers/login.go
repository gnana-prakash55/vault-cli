package controllers

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"os"
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

	res, err := http.Post("http://localhost:3000/login", "application/json", bytes.NewBuffer(jsonRes))

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

	err = os.MkdirAll("credentails", os.ModePerm)

	if err != nil {
		log.Fatalln(err)
	}

	err = ioutil.WriteFile("credentails/secret.json", jsonToken, 0644)

	if err != nil {
		log.Fatalln(err)
	}

	fmt.Println("User loged in")

}
