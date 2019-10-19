package util

import (
	"bytes"
	"log"
	"net/http"
)

var handlerErr *log.Logger
var internalErr *log.Logger

func init() {
	var buf bytes.Buffer
	handlerErr = log.New(&buf, "[Handler err]\n", log.Llongfile)
	internalErr = log.New(&buf, "[Internal err]\n", log.Llongfile)
}

func LogHandlerError(request *http.Request, err error) {
	handlerErr.Printf("> [Request Body]\n%s\n> [Err]\n%s\n", request.Body, err.Error())
}

func LogInternalError(err error) {
	internalErr.Print(err.Error())
}
