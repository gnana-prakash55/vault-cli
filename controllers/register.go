package controllers

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
)

// struct for users
type User struct {
	Email    string `json:"email"`
	Username string `json:"username"`
	Password string `json:"password"`
}

// create new user
func CreateNewUser(email string, username string, password string) {

	fmt.Println("Creating new user...")

	// new user object
	user := User{email, username, password}

	// json response
	jsonRes, err := json.Marshal(user)

	// err handler
	if err != nil {
		log.Fatalln(err)
	}

	// make the post request to the server
	res, err := http.Post("http://localhost:3000/user/create", "application/json", bytes.NewBuffer(jsonRes))

	// err handler
	if err != nil {
		log.Fatalln(err)
	}

	defer res.Body.Close()

	bodyBytes, _ := ioutil.ReadAll(res.Body)

	bodyString := string(bodyBytes)

	fmt.Println(bodyString)

}
