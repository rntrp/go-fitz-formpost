package rest

import "net/http"

func Live(w http.ResponseWriter, _ *http.Request) {
	w.WriteHeader(http.StatusOK)
}
