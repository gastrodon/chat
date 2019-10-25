package user

import (
	"chat/storage"
	"chat/models"
	"encoding/json"
	"fmt"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
)

func postUserTest(test *testing.T, data string) {
	var request *http.Request
	var err error
	request, err = http.NewRequest("POST", "/user", strings.NewReader(data))
	if err != nil {
		test.Fatal(err)
	}

	var recorder *httptest.ResponseRecorder = httptest.NewRecorder()
	http.HandlerFunc(HandleUser).ServeHTTP(recorder, request)

	var response struct {
		Key    string `json:"key"`
		UserID string `json:"user_id"`
	}
	json.Unmarshal([]byte(recorder.Body.String()), &response)

	if response.UserID == "" {
		test.Errorf("user_id got: %s", response.UserID)
	}

	if response.Key == "" {
		test.Errorf("key got: %s", response.Key)
	}

	if recorder.Code != 200 {
		test.Errorf("response.Code expected: 200, got: %d", recorder.Code)
		test.Errorf("response: %s", recorder.Body.String())
	}
}

func putUserTest(test *testing.T, data string, key string) {
	var request *http.Request
	var err error
	request, err = http.NewRequest("PUT", fmt.Sprintf("/user?key=%s", key), strings.NewReader(data))
	if err != nil {
		test.Fatal(err)
	}

	var recorder *httptest.ResponseRecorder = httptest.NewRecorder()
	http.HandlerFunc(HandleUser).ServeHTTP(recorder, request)

	var response struct {
		Username string `json:"username"`
	}
	json.Unmarshal([]byte(recorder.Body.String()), &response)

	if response.Username == "" {
		test.Errorf("username got: %s", response.Username)
	}

	if recorder.Code != 200 {
		test.Errorf("response.Code expected: 200, got: %d", recorder.Code)
		test.Errorf("response: %s", recorder.Body.String())
	}
}

func getUserTest(test *testing.T, id string) {
	var request *http.Request
	var err error
	request, err = http.NewRequest("GET", fmt.Sprintf("/user/%s", id), nil)

	if err != nil {
		test.Fatal(err)
	}

	var recorder *httptest.ResponseRecorder = httptest.NewRecorder()
	http.HandlerFunc(HandleUserTree).ServeHTTP(recorder, request)

	var response struct {
		Username string `json:"username"`
		UserID   string `json:"user_id"`
	}
	json.Unmarshal([]byte(recorder.Body.String()), &response)

	if response.Username != "foobar" {
		test.Errorf("username expected: foobar, got: %s", response.Username)
	}

	if response.UserID != id {
		test.Errorf("user_id expected: %s, got: %s", id, response.UserID)
	}

	if recorder.Code != 200 {
		test.Errorf("response.Code expected: 200, got: %d", recorder.Code)
		test.Errorf("response: %s", recorder.Body.String())
	}
}

func errUserTest(test *testing.T, method string, data string, error_desc string, code int) {
	var request *http.Request
	var err error
	request, err = http.NewRequest(method, "/user", strings.NewReader(data))
	if err != nil {
		test.Fatal(err)
	}

	var recorder *httptest.ResponseRecorder = httptest.NewRecorder()
	http.HandlerFunc(HandleUser).ServeHTTP(recorder, request)

	var response struct {
		Error string `json:"error"`
	}
	json.Unmarshal([]byte(recorder.Body.String()), &response)

	if response.Error != error_desc {
		test.Errorf("error expected: %s, got: %s", error_desc, response.Error)
		test.Errorf("response: %s", recorder.Body.String())
	}

	if recorder.Code != code {
		test.Errorf("response.Code expected: %d, got: %d", code, recorder.Code)
	}
}

func errUserTreeTest(test *testing.T, user_id string, method string, data string, error_desc string, code int) {
	var request *http.Request
	var err error
	request, err = http.NewRequest(method, fmt.Sprintf("/user/%s", user_id), strings.NewReader(data))
	if err != nil {
		test.Fatal(err)
	}

	var recorder *httptest.ResponseRecorder = httptest.NewRecorder()
	http.HandlerFunc(HandleUserTree).ServeHTTP(recorder, request)

	var response struct {
		Error string `json:"error"`
	}
	json.Unmarshal([]byte(recorder.Body.String()), &response)

	if response.Error != error_desc {
		test.Errorf("error expected: %s, got: %s", error_desc, response.Error)
		test.Errorf("response: %s", recorder.Body.String())
	}

	if recorder.Code != code {
		test.Errorf("response.Code expected: %d, got: %d", code, recorder.Code)
	}
}

func Test_postHandleUser(test *testing.T) {
	postUserTest(test, "{\"username\":\"foobar\",\"password\":\"foobar2000\"}")
}

func Test_postHandleUserNoUname(test *testing.T) {
	postUserTest(test, "{\"password\":\"foobar2000\"}")
}

func Test_putNewUsername(test *testing.T) {
	var user models.User
	user, _ = storage.NewUser("foobar", "foobar2000")

	if user.Name != "foobar" {
		test.Fatalf("user.Name expected: foobar, got: %s", user.Name)
	}

	var key string
	var err error
	key, err = storage.NewKey(user.ID, "foobar2000")

	if err != nil {
		test.Fatal(err)
	}

	putUserTest(test, "{\"username\":\"foobar3000\"}", key)

	var exists bool
	user, exists, err = storage.UserFromKey(key)

	if !exists {
		test.Fatalf("User with key %s does not exist", key)
	}

	if err != nil {
		test.Fatal(err)
	}

	if user.Name != "foobar3000" {
		test.Errorf("user.Name expected: foobar3000, got: %s", user.Name)
	}
}

func Test_putNewUsernameBlank(test *testing.T) {
	var user models.User
	user, _ = storage.NewUser("foobar", "foobar2000")

	if user.Name != "foobar" {
		test.Fatalf("user.Name expected: foobar, got: %s", user.Name)
	}

	var key string
	var err error
	key, err = storage.NewKey(user.ID, "foobar2000")

	if err != nil {
		test.Fatal(err)
	}

	putUserTest(test, "", key)

	var exists bool
	user, exists, err = storage.UserFromKey(key)

	if !exists {
		test.Fatalf("User with key %s does not exist", key)
	}

	if err != nil {
		test.Fatal(err)
	}

	if user.Name != "Anonymous" {
		test.Errorf("user.Name expected: Anonymous, got: %s", user.Name)
	}
}

func Test_putHandleUserNoKey(test *testing.T) {
	errUserTest(test, "PUT", "", "no_key", 401)
}

func Test_putHandleUserBadKey(test *testing.T) {
	var request *http.Request
	var err error
	request, err = http.NewRequest("PUT", "/user?key=foobar", nil)

	if err != nil {
		test.Fatal(err)
	}

	var recorder *httptest.ResponseRecorder = httptest.NewRecorder()
	http.HandlerFunc(HandleUser).ServeHTTP(recorder, request)

	var response struct {
		Error string `json:"error"`
	}
	json.Unmarshal([]byte(recorder.Body.String()), &response)

	if response.Error != "bad_key" {
		test.Errorf("error expected: %s, got: %s", "bad_key", response.Error)
		test.Errorf("response: %s", recorder.Body.String())
	}

	if recorder.Code != 401 {
		test.Errorf("response.Code expected: %d, got: %d", 401, recorder.Code)
	}
}

func Test_getHandleUserTree(test *testing.T) {
	var user models.User
	user, _ = storage.NewUser("foobar", "foobar2000")

	if user.Name != "foobar" {
		test.Fatalf("user.Name expected: foobar, got: %s", user.Name)
	}

	getUserTest(test, user.ID)
}

func Test_getHandleUserTreeNoUser(test *testing.T) {
	errUserTreeTest(test, "foobar", "GET", "", "no_such_user", 404)
}

func Test_postHandleUserNoPasswd(test *testing.T) {
	errUserTest(test, "POST", "{\"username\":\"foobar\"}", "password_missing", 400)
}

func Test_postHandleUserBadPasswd(test *testing.T) {
	errUserTest(test, "POST", "{\"password\":42}", "malformed_json", 400)
}

func Test_postHandleUserBadUname(test *testing.T) {
	errUserTest(test, "POST", "{\"username\":42,\"password\":\"foobar2000\"}", "malformed_json", 400)
}

func Test_HandleUserBadMethod(test *testing.T) {
	errUserTest(test, "OOOO", "", "bad_method", 405)
}

func Test_HandleUserTreeBadMethod(test *testing.T) {
	errUserTreeTest(test, "", "OOOO", "", "bad_method", 405)
}
