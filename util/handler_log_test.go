package util

import (
	"errors"
	"fmt"
	"net/http"
	"net/http/httptest"
	"testing"
)

func dummyHandler(response http.ResponseWriter, request *http.Request) {
	response.Header().Set("Content-Type", "application/json")
	response.WriteHeader(200)
	fmt.Fprint(response, "{\"test\":true}")
}

func Test_LogHandlerErr(test *testing.T) {
	var request *http.Request
	var recorder *httptest.ResponseRecorder = httptest.NewRecorder()
	var handler http.HandlerFunc = http.HandlerFunc(dummyHandler)
	var err error
	request, err = http.NewRequest("GET", "/", nil)

	if err != nil {
		test.Fatal(err)
	}

	handler.ServeHTTP(recorder, request)

	var test_err error = errors.New("LogHandlerError")
	LogHandlerError(request, test_err)
}

func Test_LogInternalError(test *testing.T) {
	LogInternalError(errors.New("LogInternalError"))
}
