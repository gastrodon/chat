package key

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

func HTTPTestRequest(method string, endpoint string, post_data *string, handler_func func(http.ResponseWriter, *http.Request)) (*httptest.ResponseRecorder, error) {
	var request *http.Request
	var err error
	if post_data != nil {
		request, err = http.NewRequest(method, endpoint, strings.NewReader(*post_data))
	} else {
		request, err = http.NewRequest(method, endpoint, nil)
	}

	var recorder *httptest.ResponseRecorder = httptest.NewRecorder()
	if err != nil {
		return recorder, err
	}

	http.HandlerFunc(handler_func).ServeHTTP(recorder, request)
	return recorder, nil
}

func postKeyTest(test *testing.T, user_id string, password string) {
	var post_data string = fmt.Sprintf("{\"user_id\":\"%s\",\"password\":\"%s\"}", user_id, password)
	var recorder *httptest.ResponseRecorder
	var err error
	recorder, err = HTTPTestRequest("POST", "/key", &post_data, HandleKey)

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
	user, exists, err = storage.UserFromKey(response.Key)

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

func deleteKeyTest(test *testing.T, user_id string, key string) {
	var post_data string = fmt.Sprintf("{\"user_id\":\"%s\"}", user_id)
	var recorder *httptest.ResponseRecorder
	var err error
	recorder, err = HTTPTestRequest("DELETE", fmt.Sprintf("/key?key=%s", key), &post_data, HandleKey)

	if err != nil {
		test.Fatal(err)
	}

	var response struct {
		UserID string `json:"user_id"`
	}
	json.Unmarshal([]byte(recorder.Body.String()), &response)

	if response.UserID == "" {
		test.Error("No user_id returned")
	}

	if response.UserID != user_id {
		test.Errorf("user_id expected: %s, got: %s", user_id, response.UserID)
	}

	var exists bool
	_, exists, err = storage.UserFromKey(key)

	if err != nil {
		test.Fatal(err)
	}

	if exists {
		test.Error("UserFromKey key does exist")
	}

	_, exists, err = storage.UserFromID(user_id)

	if err != nil {
		test.Fatal(err)
	}

	if !exists {
		test.Error("UserFromID user_id does not exist")
	}
}

func errKeyTest(test *testing.T, method string, data string, key *string, error_desc string, code int) {
	var recorder *httptest.ResponseRecorder = httptest.NewRecorder()

	var request *http.Request
	var err error
	if key != nil {
		request, err = http.NewRequest(method, fmt.Sprintf("/key?key=%s", *key), strings.NewReader(data))
	} else {
		request, err = http.NewRequest(method, "/key?key=%s", strings.NewReader(data))

	}
	if err != nil {
		test.Fatal(err)
	}

	http.HandlerFunc(HandleKey).ServeHTTP(recorder, request)

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
		test.Errorf("response: %s", recorder.Body.String())
	}
}

func Test_postHandleKey(test *testing.T) {
	var uname string = "foobar"
	var password string = "foobar2000"
	var user models.User
	user, _ = storage.NewUser(uname, password)

	if user.Name != uname {
		test.Fatalf("user.Name expected: %s got: %s", uname, user.Name)
	}

	postKeyTest(test, user.ID, password)
}

func Test_postHandleKeyWrongPasswd(test *testing.T) {
	var passwd string = "foobar2000"
	var user models.User
	user, _ = storage.NewUser("foobar", passwd)

	if user.Name != "foobar" {
		test.Fatalf("user.Name expected: foobar, got: %s", user.Name)
	}

	var post_data string = fmt.Sprintf("{\"user_id\":\"%s\",\"password\":\"oof!\"}", user.ID)

	errKeyTest(test, "POST", post_data, nil, "bad_password", 403)
}

func Test_postHandleKeyMissingPasswd(test *testing.T) {
	var post_data string = "{\"user_id\":\"0\"}"

	errKeyTest(test, "POST", post_data, nil, "password_missing", 400)
}

func Test_postHandleKeyMissingUserID(test *testing.T) {
	var post_data string = "{\"password\":\"0\"}"

	errKeyTest(test, "POST", post_data, nil, "user_id_missing", 400)
}

func Test_postHandleKeyNoSuchUser(test *testing.T) {
	var post_data string = "{\"password\":\"0\",\"user_id\":\"0\"}"

	errKeyTest(test, "POST", post_data, nil, "no_such_user", 404)
}

func Test_postHandleKeyIntUserID(test *testing.T) {
	var post_data string = "{\"password\":\"0\",\"user_id\":0}"

	errKeyTest(test, "POST", post_data, nil, "malformed_json", 400)
}

func Test_postHandleKeyMalformed(test *testing.T) {
	var post_data string = "Wait, this isn't json!{{{{{{]]][:]}"

	errKeyTest(test, "POST", post_data, nil, "malformed_json", 400)
}

func Test_deleteHandleKey(test *testing.T) {
	var uname string = "foobar"
	var password string = "foobar2000"
	var user models.User
	user, _ = storage.NewUser(uname, password)

	if user.Name != uname {
		test.Fatalf("user.Name expected: %s got: %s", uname, user.Name)
	}

	var key string
	var err error
	key, err = storage.NewKey(user.ID, password)

	if err != nil {
		test.Fatal(err)
	}

	deleteKeyTest(test, user.ID, key)
}

func Test_deleteHandleKeyMissingKey(test *testing.T) {
	var post_data string = "{\"user_id\":\"0\"}"

	errKeyTest(test, "DELETE", post_data, nil, "no_key", 401)
}

func Test_deleteHandleKeyMissingUserID(test *testing.T) {
	var post_data string = "{}"

	var key string = "0"
	errKeyTest(test, "DELETE", post_data, &key, "user_id_missing", 400)
}

func Test_deleteHandleKeyNoSuchUser(test *testing.T) {
	var post_data string = "{\"key\":\"0\",\"user_id\":\"0\"}"

	var key string = "0"
	errKeyTest(test, "DELETE", post_data, &key, "bad_key", 401)
}

func Test_deleteHandleKeyBadKey(test *testing.T) {
	var uname string = "foobar"
	var password string = "foobar2000"
	var user models.User
	user, _ = storage.NewUser(uname, password)

	if user.Name != uname {
		test.Fatalf("user.Name expected: %s got: %s", uname, user.Name)
	}

	var post_data string = fmt.Sprintf("{\"user_id\":\"%s\"}", user.ID)

	var key string = "0"
	errKeyTest(test, "DELETE", post_data, &key, "bad_key", 401)
}

func Test_HandleKeyBadMethod(test *testing.T) {
	errKeyTest(test, "OOOO", "", nil, "bad_method", 405)
}

func Test_deleteHandleKeyIntUserID(test *testing.T) {
	var post_data string = "{\"password\":\"0\",\"user_id\":0}"

	var key string = ""
	errKeyTest(test, "DELETE", post_data, &key, "malformed_json", 400)
}
