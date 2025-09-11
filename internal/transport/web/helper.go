package web

import (
	"log"
	"net/http"
)

func handleError(w http.ResponseWriter, msg string, err error, code int) {
	if err != nil {
		log.Printf("[ERR] %s: %v", msg, err)
	}
	http.Error(w, msg, code)
}
