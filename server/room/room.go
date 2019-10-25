package room

import (
	"chat/models"
	"chat/storage"
	"chat/util"
	"encoding/json"
	"io/ioutil"
	"net/http"
)

func postHandleRoom(response http.ResponseWriter, request *http.Request) {
	if len(request.URL.Query()["key"]) == 0 {
		util.HandleHTTPErr(response, "no_key", 401)
		return
	}

	var key string = request.URL.Query()["key"][0]
	var user models.User
	var exists bool
	var err error
	user, exists, err = storage.UserFromKey(key)

	if !exists {
		util.HandleHTTPErr(response, "bad_key", 403)
		return
	}

	if err != nil {
		util.HandleHTTPErr(response, "internal_err", 500)
		util.LogHandlerError(request, err)
		return
	}

	var body []byte
	body, err = ioutil.ReadAll(request.Body)

	if err != nil {
		util.HandleHTTPErr(response, "internal_err", 500)
		util.LogHandlerError(request, err)
		return
	}

	var json_body struct {
		Open bool   `json:"open"`
		Name string `json:"name"`
	}

	err = json.Unmarshal(body, &json_body)

	if err != nil {
		util.HandleHTTPErr(response, "internal_err", 500)
		util.LogHandlerError(request, err)
		return
	}

	if json_body.Name == "" {
		json_body.Name = "Anonymous chat"
	}

	var room models.Room
	room, err = storage.NewRoom(json_body.Name, json_body.Open, user.ID)

	var response_map map[string]interface{} = map[string]interface{}{
		"name":       room.Name,
		"open":       room.Open,
		"owner":      room.OwnerId,
		"room_id":    room.ID,
		"user_count": room.UserCount,
	}
	util.SendHTTPJsonResponse(response, response_map)
}

func HandleRoom(response http.ResponseWriter, request *http.Request) {
	switch request.Method {
	case "POST":
		postHandleRoom(response, request)
	default:
		util.HandleHTTPErr(response, "bad_method", 405)
	}
}
