package key

import (
	"chat/models"
	"chat/storage"
	"chat/util"
	"encoding/json"
	"io/ioutil"
	"net/http"
)

func postHandleKey(response http.ResponseWriter, request *http.Request) {
	var body []byte
	var err error

	body, err = ioutil.ReadAll(request.Body)

	if err != nil {
		util.HandleHTTPErr(response, "internal_err", 500)
		util.LogHandlerError(request, err)
		return
	}

	var json_body struct {
		UserID   string `json:"user_id"`
		Password string `json:"password"`
	}
	err = json.Unmarshal(body, &json_body)

	if err != nil {
		util.HandleHTTPErr(response, "malformed_json", 400)
		return
	}

	if json_body.Password == "" {
		util.HandleHTTPErr(response, "password_missing", 400)
		return
	}

	if json_body.UserID == "" {
		util.HandleHTTPErr(response, "user_id_missing", 400)
		return
	}

	var user models.User
	var user_exists bool
	user, user_exists, err = storage.UserFromID(json_body.UserID)

	if err != nil {
		util.HandleHTTPErr(response, "internal_err", 500)
		util.LogHandlerError(request, err)
		return
	}

	if !user_exists {
		util.HandleHTTPErr(response, "no_such_user", 404)
		return
	}

	var passwd_match bool
	passwd_match, err = storage.CheckPasswd(user.ID, json_body.Password)

	if err != nil {
		util.HandleHTTPErr(response, "internal_err", 500)
		util.LogHandlerError(request, err)
		return
	}

	if !passwd_match {
		util.HandleHTTPErr(response, "bad_password", 403)
		return
	}

	var key string
	key, err = storage.NewKey(user.ID, json_body.Password)

	if err != nil {
		util.HandleHTTPErr(response, "internal_err", 500)
		util.LogHandlerError(request, err)
		return
	}

	var response_map map[string]interface{} = map[string]interface{}{
		"key": key,
	}

	util.SendHTTPJsonResponse(response, response_map)
}

func deleteHandleKey(response http.ResponseWriter, request *http.Request) {
	if len(request.URL.Query()["key"]) == 0 {
		util.HandleHTTPErr(response, "no_key", 401)
		return
	}

	var key string = request.URL.Query()["key"][0]

	var body []byte
	var err error
	body, err = ioutil.ReadAll(request.Body)

	if err != nil {
		util.HandleHTTPErr(response, "internal_err", 500)
		util.LogHandlerError(request, err)
		return
	}

	var json_body struct {
		UserID string `json:"user_id"`
	}
	err = json.Unmarshal(body, &json_body)

	if err != nil {
		util.HandleHTTPErr(response, "malformed_json", 400)
		return
	}

	if json_body.UserID == "" {
		util.HandleHTTPErr(response, "user_id_missing", 400)
		return
	}

	var exists bool
	_, exists, err = storage.UserFromKey(key)

	if err != nil {
		util.HandleHTTPErr(response, "internal_err", 500)
		util.LogHandlerError(request, err)
		return
	}

	if !exists {
		util.HandleHTTPErr(response, "bad_key", 401)
		return
	}

	err = storage.DeleteKey(key)

	if err != nil {
		util.HandleHTTPErr(response, "internal_err", 500)
		util.LogHandlerError(request, err)
		return
	}

	var response_map map[string]interface{} = map[string]interface{}{
		"user_id": json_body.UserID,
	}
	util.SendHTTPJsonResponse(response, response_map)

}

func HandleKey(response http.ResponseWriter, request *http.Request) {
	switch request.Method {
	case "POST":
		postHandleKey(response, request)
	case "DELETE":
		deleteHandleKey(response, request)
	default:
		util.HandleHTTPErr(response, "bad_method", 405)
	}
}
