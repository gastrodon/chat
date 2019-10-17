package server

import (
    "log"
    "errors"
    "chat/io"
    "net/http"
    "io/ioutil"
    "chat/models"
)

func postHandleUser(response http.ResponseWriter, request *http.Request) {
    var uname, passwd []string
    var uname_supplied, passwd_supplied bool

    var response_map map[string]interface{}
    var err error

    request.ParseForm()
    log.Printf("request: %s", request)

    passwd, passwd_supplied = request.Form["password"]
    log.Printf("passwd: %s", passwd)

    if !passwd_supplied {
        HandleHTTPErr(response, errors.New("passwd_missing"), 400)
        return
    }

    uname, uname_supplied = request.Form["username"]
    log.Printf("uname: %s", uname)

    if !uname_supplied {
        uname[0] = "Anonymous"
    }

    var user models.User = io.NewUser(uname[0], passwd[0])
    var key string
    key, err = io.NewKey(user.ID, passwd[0])

    if err != nil {
        HandleHTTPErr(response, err, 500)
    }

    response_map = map[string]interface{} {
        "user_id":  user.ID,
        "key":      key,
    }

    SendHTTPJsonResponse(response, response_map)
}

func HandleUser(response http.ResponseWriter, request *http.Request) {
    switch request.Method {
        case "POST":
            postHandleUser(response, request)
    }
}
