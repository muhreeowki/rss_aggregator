package main

import "net/http"

func handlerReadiness(writer http.ResponseWriter, r *http.Request) {
	respondWithJSON(writer, 200, struct{}{})
}
