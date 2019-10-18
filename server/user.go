package server

import (
	"chat/io"
	"chat/models"
	"encoding/json"
	"errors"
	"io/ioutil"
	"log"
	"net/http"
)

func postHandleUser(response http.ResponseWriter, request *http.Request) {
	var body []byte
	var body_json struct {
		Username string `json:"username"`
		Password string `json:"password"`
	}

	var response_map map[string]interface{}
	var err error

	defer request.Body.Close()
	body, err = ioutil.ReadAll(request.Body)

	if err != nil {
		HandleHTTPErr(response, err, 500)
		return
	}

	err = json.Unmarshal(body, &body_json)

	log.Printf("body_json: %s", body_json)

	log.Printf("passwd: %s", body_json.Password)
	if body_json.Password == "" {
		HandleHTTPErr(response, errors.New("password_missing"), 400)
		return
	}

	log.Printf("uname: %s", body_json.Username)
	if body_json.Username == "" {
		body_json.Username = "Anonymous"
	}

	var user models.User = io.NewUser(body_json.Username, body_json.Password)
	var key string
	key, err = io.NewKey(user.ID, body_json.Password)

	if err != nil {
		HandleHTTPErr(response, err, 500)
	}

	response_map = map[string]interface{}{
		"user_id": user.ID,
		"key":     key,
	}

	SendHTTPJsonResponse(response, response_map)
}

func HandleUser(response http.ResponseWriter, request *http.Request) {
	switch request.Method {
	case "POST":
		postHandleUser(response, request)
	}
}
