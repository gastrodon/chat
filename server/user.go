package server

import (
	"chat/io"
	"chat/models"
	"encoding/json"
	"errors"
	"io/ioutil"
	"net/http"
)

func putHandleUser(response http.ResponseWriter, request *http.Request) {
	var body []byte
	var json_body struct{
		Username string `json:"username"`
	}

	var response_map map[string]interface{}
	var user models.User
	var exists bool
	var err error

	defer request.Body.Close()

	if len(request.URL.Query()["key"]) == 0 {
		HandleHTTPErr(response, errors.New("no_key"), 401)
	}

	var key string = request.URL.Query()["key"][0]

	user, exists, err = io.UserFromKey(key)

	if err != nil {
		HandleHTTPErr(response, err, 500)
		return
	}

	if !exists {
		HandleHTTPErr(response, errors.New("bad_key"), 401)
		return
	}


	body, err = ioutil.ReadAll(request.Body)

	if err != nil {
		HandleHTTPErr(response, err, 500)
		return
	}

	err = json.Unmarshal(body, &json_body)

	if json_body.Username == "" {
		json_body.Username = "Anonymous"
	}

	user, err = io.UpdateUname(user.ID, json_body.Username)

	if err != nil {
		HandleHTTPErr(response, err, 500)
		return
	}

	response_map = map[string]interface{} {
		"username": user.Name,
	}

	SendHTTPJsonResponse(response, response_map)
	return
}

func postHandleUser(response http.ResponseWriter, request *http.Request) {
	var body []byte
	var json_body struct {
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

	err = json.Unmarshal(body, &json_body)

	if json_body.Password == "" {
		HandleHTTPErr(response, errors.New("password_missing"), 400)
		return
	}

	if json_body.Username == "" {
		json_body.Username = "Anonymous"
	}

	var user models.User = io.NewUser(json_body.Username, json_body.Password)
	var key string
	key, err = io.NewKey(user.ID, json_body.Password)

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
		return
	case "PUT":
		putHandleUser(response, request)
		return
	}
}
