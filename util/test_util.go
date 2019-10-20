package util

import (
    "net/http"
    "net/http/httptest"
    "strings"
)

func HTTPTestRequest(method string, endpoint string, post_data *string, handler_func func(http.ResponseWriter, *http.Request)) (*httptest.ResponseRecorder, error) {
	var recorder *httptest.ResponseRecorder = httptest.NewRecorder()
	var request *http.Request
	var err error
	if post_data != nil {
		request, err = http.NewRequest(method, endpoint, strings.NewReader(*post_data))
	} else {
		request, err = http.NewRequest(method, endpoint, nil)
	}

	if err != nil {
		return recorder, err
	}

	var handler http.HandlerFunc = http.HandlerFunc(handler_func)
	handler.ServeHTTP(recorder, request)
	return recorder, nil
}
