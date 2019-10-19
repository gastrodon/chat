package server

import (
    "net/http"
    "net/http/httptest"
    "chat/io"
    "chat/models"
    "fmt"
    "strings"
    "encoding/json"
    "testing"
)

func postKeyTest(test *testing.T, user_id string, password string) {
    var post_data string = fmt.Sprintf("{\"user_id\":\"%s\",\"password\":\"%s\"}", user_id, password)

    var handler http.HandlerFunc = http.HandlerFunc(HandleKey)
    var recorder *httptest.ResponseRecorder = httptest.NewRecorder()

    var request *http.Request
    var err error
    request, err = http.NewRequest("POST", "/key", strings.NewReader(post_data))
    if err != nil {
        test.Fatal(err)
    }

    handler.ServeHTTP(recorder, request)
    var response struct {
        Key string `json:"key"`
    }
    json.Unmarshal([]byte(recorder.Body.String()), &response)

    if response.Key == "" {
        test.Errorf("key got: %s", response.Key)
    }

    var key_user models.User
    var key_exists bool
    key_user, key_exists, err = io.UserFromKey(response.Key)

    if err != nil {
        test.Fatal(err)
    }

    if !key_exists {
        test.Errorf("User not found for key %s", response.Key)
    }

    if key_user.ID != user_id {
        test.Errorf("user.ID expected: %s, got: %s", user_id, key_user.ID)
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
    var response struct{
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

func Test_postHandleKey(test *testing.T)  {
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
