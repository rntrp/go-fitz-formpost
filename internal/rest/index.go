package rest

import (
	_ "embed"
	"encoding/base64"
	"hash/fnv"
	"net/http"
	"strconv"
)

//go:embed index.html
var html []byte

var htmlContentLength = strconv.Itoa(len(html))

var htmlETag = b64ETag(html)

func b64ETag(b []byte) string {
	h := fnv.New64a()
	if _, err := h.Write(b); err != nil {
		panic(err)
	}
	b64 := base64.URLEncoding.EncodeToString(h.Sum(nil))
	return "\"" + b64 + "\""
}

func Index(w http.ResponseWriter, r *http.Request) {
	wh := w.Header()
	wh.Set("ETag", htmlETag)
	if r.Header.Get("If-None-Match") == htmlETag {
		w.WriteHeader(http.StatusNotModified)
		return
	}
	wh.Set("Content-Type", "text/html; charset=utf-8")
	wh.Set("Content-Length", htmlContentLength)
	w.Write(html)
}
