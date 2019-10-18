package server

import (
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
    "encoding/json"
)

func postUserTest(test *testing.T, data string) {
	var post_data string = data
	var handler http.HandlerFunc = http.HandlerFunc(HandleUser)
	var recorder *httptest.ResponseRecorder = httptest.NewRecorder()

	var request *http.Request
	var err error
	request, err = http.NewRequest("POST", "/user", strings.NewReader(post_data))
	if err != nil {
		test.Fatal(err)
	}

	handler.ServeHTTP(recorder, request)
    var response struct{
        Key string `json:"key"`
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

func Test_postHandleUser(test *testing.T) {
    postUserTest(test, "{\"username\":\"foobar\",\"password\":\"foobar2000\"}")
}

func Test_postHandleUserNoUname(test *testing.T) {
    postUserTest(test, "{\"password\":\"foobar2000\"}")
}

func Test_postHandleUserBadUname(test *testing.T) {
    postUserTest(test, "{\"username\":42,\"password\":\"foobar2000\"}")
}
