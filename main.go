package main

import (
	"chat/server"
	"fmt"
	"net/http"
	"os"
)

func main() {
	var port string = ":5000"

	if len(os.Args) >= 2 {
		port = fmt.Sprintf(":%s", os.Args[1])
	}

	http.HandleFunc("/key", server.HandleKey)

	http.HandleFunc("/user", server.HandleUser)
	http.HandleFunc("/user/", server.HandleUserTree)

	http.ListenAndServe(port, nil)
}
