package chat

import (
	"chat/server/key"
	"chat/server/user"
	"fmt"
	"net/http"
	"os"
)

func main() {
	var port string = ":5000"

	if len(os.Args) >= 2 {
		port = fmt.Sprintf(":%s", os.Args[1])
	}

	http.HandleFunc("/key", key.HandleKey)
	http.HandleFunc("/user", user.HandleUser)
	http.HandleFunc("/user/", user.HandleUserTree)

	http.ListenAndServe(port, nil)
}
