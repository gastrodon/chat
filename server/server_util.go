package server

import (
	"encoding/json"
	"fmt"
	"net/http"
)

func HandleHTTPErr(response http.ResponseWriter, err error, code int) {
	var parse_error error
	var response_data []byte
	var response_map map[string]string = map[string]string{
		"error": err.Error(),
	}

	response_data, parse_error = json.Marshal(response_map)

	if parse_error != nil {
		http.Error(response, parse_error.Error(), 500)
		return
	}

	http.Error(response, string(response_data), code)
}

func SendHTTPJsonResponse(response http.ResponseWriter, response_map map[string]interface{}) {
	var parse_error error
	var response_data []byte
	response_data, parse_error = json.Marshal(response_map)

	if parse_error != nil {
		HandleHTTPErr(response, parse_error, 500)
		return
	}

	response.Header().Set("Content-Type", "application/json")
	response.WriteHeader(200)
	fmt.Fprint(response, string(response_data))
}
