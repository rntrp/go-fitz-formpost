package rest

import (
	_ "embed"
	"net/http"
	"strconv"
)

//go:embed index.html
var html []byte

var htmlContentLength = strconv.Itoa(len(html))

func Index(w http.ResponseWriter, _ *http.Request) {
	w.Header().Set("Content-Type", "text/html; charset=utf-8")
	w.Header().Set("Content-Length", htmlContentLength)
	w.Write(html)
}
