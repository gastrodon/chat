package server

import (
	"chat/io"
	"chat/models"
	"chat/util"
	"encoding/json"
	"io/ioutil"
	"net/http"
	"strings"
)

func putHandleUser(response http.ResponseWriter, request *http.Request) {
	if len(request.URL.Query()["key"]) == 0 {
		HandleHTTPErr(response, "no_key", 401)
		return
	}

	var user models.User
	var exists bool
	var err error
	user, exists, err = io.UserFromKey(request.URL.Query()["key"][0])

	if err != nil {
		HandleHTTPErr(response, "internal_err", 500)
		util.LogHandlerError(request, err)
		return
	}

	if !exists {
		HandleHTTPErr(response, "bad_key", 401)
		return
	}

	var body []byte
	body, err = ioutil.ReadAll(request.Body)

	if err != nil {
		HandleHTTPErr(response, "internal_err", 500)
		util.LogHandlerError(request, err)
		return
	}

	var json_body struct {
		Username string `json:"username"`
	}
	err = json.Unmarshal(body, &json_body)

	if json_body.Username == "" {
		json_body.Username = "Anonymous"
	}

	user, err = io.UpdateUname(user.ID, json_body.Username)

	if err != nil {
		HandleHTTPErr(response, "internal_err", 500)
		util.LogHandlerError(request, err)
		return
	}

	var response_map map[string]interface{} = map[string]interface{}{
		"username": user.Name,
	}

	SendHTTPJsonResponse(response, response_map)
	return
}

func postHandleUser(response http.ResponseWriter, request *http.Request) {
	var body []byte
	var err error
	body, err = ioutil.ReadAll(request.Body)

	if err != nil {
		HandleHTTPErr(response, "internal_err", 500)
		util.LogHandlerError(request, err)
		return
	}

	var json_body struct {
		Username string `json:"username"`
		Password string `json:"password"`
	}
	err = json.Unmarshal(body, &json_body)

	if err != nil {
		HandleHTTPErr(response, "malformed_json", 400)
		return
	}

	if json_body.Password == "" {
		HandleHTTPErr(response, "password_missing", 400)
		return
	}

	if json_body.Username == "" {
		json_body.Username = "Anonymous"
	}

	var user models.User = io.NewUser(json_body.Username, json_body.Password)
	var key string
	key, err = io.NewKey(user.ID, json_body.Password)

	if err != nil {
		HandleHTTPErr(response, "internal_err", 500)
		util.LogHandlerError(request, err)
		return
	}

	var response_map map[string]interface{} = map[string]interface{}{
		"user_id": user.ID,
		"key":     key,
	}

	SendHTTPJsonResponse(response, response_map)
}

func getHandleUserTree(response http.ResponseWriter, request *http.Request) {
	var id string = strings.Replace(request.URL.Path, "/user/", "", 1)
	var user models.User
	var exists bool
	var err error
	user, exists, err = io.UserFromID(id)

	if !exists {
		HandleHTTPErr(response, "no_such_user", 404)
		return
	}

	if err != nil {
		HandleHTTPErr(response, "internal_err", 500)
		util.LogHandlerError(request, err)
		return
	}

	var response_map map[string]interface{} = map[string]interface{}{
		"username": user.Name,
		"user_id":  user.ID,
	}

	SendHTTPJsonResponse(response, response_map)
	return
}

func HandleUser(response http.ResponseWriter, request *http.Request) {
	switch request.Method {
	case "POST":
		postHandleUser(response, request)
		return
	case "PUT":
		putHandleUser(response, request)
		return
	default:
		HandleHTTPErr(response, "bad_method", 405)
	}
}

func HandleUserTree(response http.ResponseWriter, request *http.Request) {
	switch request.Method {
	case "GET":
		getHandleUserTree(response, request)
		return
	default:
		HandleHTTPErr(response, "bad_method", 405)
	}
}
