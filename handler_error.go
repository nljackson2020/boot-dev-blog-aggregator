package main

import (
	"net/http"
)

func handlerError(w http.ResponseWriter, r *http.Request) {
	msg := "Internal Server Error"
	respondWithError(w, http.StatusInternalServerError, msg)
}
