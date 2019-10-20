package server

import (
	"chat/io"
	"chat/models"
	"chat/util"
	"encoding/json"
	"fmt"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
)

func postKeyTest(test *testing.T, user_id string, password string) {
	var post_data string = fmt.Sprintf("{\"user_id\":\"%s\",\"password\":\"%s\"}", user_id, password)
	var recorder *httptest.ResponseRecorder
	var err error
	recorder, err = util.HTTPTestRequest("POST", "/key", &post_data, HandleKey)

	if err != nil {
		test.Fatal(err)
	}

	var response struct {
		Key string `json:"key"`
	}
	json.Unmarshal([]byte(recorder.Body.String()), &response)

	if response.Key == "" {
		test.Error("No key returned")
	}

	var user models.User
	var exists bool
	user, exists, err = io.UserFromKey(response.Key)

	if err != nil {
		test.Fatal(err)
	}

	if !exists {
		test.Errorf("UserFromKey %s does not exist", response.Key)
	}

	if user.ID != user_id {
		test.Errorf("user.ID expected: %s, got: %s", user_id, user.ID)
	}

	if recorder.Code != 200 {
		test.Errorf("response.Code expected: 200, got: %d", recorder.Code)
		test.Errorf("response: %s", recorder.Body.String())
	}
}

func errKeyTest(test *testing.T, method string, data string, error_desc string, code int) {
	var post_data string = data
	var handler http.HandlerFunc = http.HandlerFunc(HandleKey)
	var recorder *httptest.ResponseRecorder = httptest.NewRecorder()

	var request *http.Request
	var err error
	request, err = http.NewRequest(method, "/key", strings.NewReader(post_data))
	if err != nil {
		test.Fatal(err)
	}

	handler.ServeHTTP(recorder, request)
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

func Test_postHandleKey(test *testing.T) {
	var password string = "foobar2000"
	var user models.User = io.NewUser("foobar", password)

	if user.Name != "foobar" {
		test.Fatalf("user.Name expected: foobar, got: %s", user.Name)
	}

	postKeyTest(test, user.ID, password)
}

func Test_postHandleKeyWrongPasswd(test *testing.T) {
	var passwd string = "foobar2000"
	var user models.User = io.NewUser("foobar", passwd)

	if user.Name != "foobar" {
		test.Fatalf("user.Name expected: foobar, got: %s", user.Name)
	}

	var post_data string = fmt.Sprintf("{\"user_id\":\"%s\",\"password\":\"oof!\"}", user.ID)

	errKeyTest(test, "POST", post_data, "bad_password", 403)
}

func Test_postHandleKeyMissingPasswd(test *testing.T) {
	var post_data string = "{\"user_id\":\"0\"}"

	errKeyTest(test, "POST", post_data, "password_missing", 400)
}

func Test_postHandleKeyMissingUserID(test *testing.T) {
	var post_data string = "{\"password\":\"0\"}"

	errKeyTest(test, "POST", post_data, "user_id_missing", 400)
}

func Test_postHandleKeyNoSuchUser(test *testing.T) {
	var post_data string = "{\"password\":\"0\",\"user_id\":\"0\"}"

	errKeyTest(test, "POST", post_data, "no_such_user", 404)
}

func Test_postHandleKeyIntUserID(test *testing.T) {
	var post_data string = "{\"password\":\"0\",\"user_id\":0}"

	errKeyTest(test, "POST", post_data, "malformed_json", 400)
}

func Test_postHandleKeyMalformed(test *testing.T) {
	var post_data string = "Wait, this isn't json!{{{{{{]]][:]}"

	errKeyTest(test, "POST", post_data, "malformed_json", 400)
}

func Test_HandleKeyBadMethod(test *testing.T) {
	errKeyTest(test, "OOOO", "", "bad_method", 405)
}
