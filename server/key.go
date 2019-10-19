package server

import (
	"chat/io"
	"chat/models"
	"chat/util"
	"encoding/json"
	"io/ioutil"
	"net/http"
)

// TODO: reorganize variable declaration
func postHandleKey(response http.ResponseWriter, request *http.Request) {
	var response_map map[string]interface{}
	var body []byte
	var err error
	var json_body struct {
		UserID   string `json:"user_id"`
		Password string `json:"password"`
	}

	body, err = ioutil.ReadAll(request.Body)

	if err != nil {
		HandleHTTPErr(response, "internal_err", 500)
		util.LogHandlerError(request, err)
		return
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

	if json_body.UserID == "" {
		HandleHTTPErr(response, "user_id_missing", 400)
		return
	}

	var user models.User
	var user_exists bool
	user, user_exists, err = io.UserFromID(json_body.UserID)

	if err != nil {
		HandleHTTPErr(response, "internal_err", 500)
		util.LogHandlerError(request, err)
		return
	}

	if !user_exists {
		HandleHTTPErr(response, "no_such_user", 404)
		return
	}

	var passwd_match bool
	passwd_match, err = io.CheckPasswd(user.ID, json_body.Password)

	if err != nil {
		HandleHTTPErr(response, "internal_err", 500)
		util.LogHandlerError(request, err)
		return
	}

	if !passwd_match {
		HandleHTTPErr(response, "bad_password", 403)
		return
	}

	var key string
	key, err = io.NewKey(user.ID, json_body.Password)

	if err != nil {
		HandleHTTPErr(response, "internal_err", 500)
		util.LogHandlerError(request, err)
		return
	}

	response_map = map[string]interface{}{
		"key": key,
	}

	SendHTTPJsonResponse(response, response_map)
}

func HandleKey(response http.ResponseWriter, request *http.Request) {
	switch request.Method {
	case "POST":
		postHandleKey(response, request)
	default:
		HandleHTTPErr(response, "bad_method", 405)
	}
}
