package helpers

import (
	"fmt"
	"net/http"
	"runtime/debug"

	"github.com/Orololuwa/go-backend-boilerplate/src/config"
)

var app *config.AppConfig

func NewHelper(a *config.AppConfig){
	app = a
}

func ClientError(w http.ResponseWriter, status int, message string) {
	errorMessage := message
	if errorMessage == "" {
		errorMessage = "Client error with status of "
	}

	app.InfoLog.Println(errorMessage, status)
	http.Error(w, http.StatusText(status), status)
}

func ServerError(w http.ResponseWriter, err error) {
	trace := fmt.Sprintf("%s\n%s", err.Error(), debug.Stack())
	app.ErrorLog.Println(trace)
	http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
}