package main

import (
    "net/http"
    "chat/server"
)

func main() {
    http.HandleFunc("/user", server.HandleUser)

    http.ListenAndServe(":5000", nil)
}
