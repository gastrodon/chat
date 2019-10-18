package server

import (
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
    "encoding/json"
)

func postTestErr(test *testing.T, endpoint string, data string, error_desc string, code int) {
    var post_data string = data
	var handler http.HandlerFunc = http.HandlerFunc(HandleUser)
	var recorder *httptest.ResponseRecorder = httptest.NewRecorder()

	var request *http.Request
	var err error
	request, err = http.NewRequest("POST", endpoint, strings.NewReader(post_data))
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
    }

    if recorder.Code != code {
        test.Errorf("response.Code expected: %d, got: %d", code, recorder.Code)
    }
}

func Test_postHandleUserNoPasswd(test *testing.T) {
    postTestErr(test, "/user", "{\"username\":\"foobar\"}", "password_missing", 400)
}

func Test_postHandleUserBadPasswd(test *testing.T) {
    postTestErr(test, "/user", "{\"password\":42}", "password_missing", 400)
}
