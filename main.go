package main

import (
	"chat/server"
	"net/http"
)

func main() {
	http.HandleFunc("/user", server.HandleUser)

	http.ListenAndServe(":5000", nil)
}
