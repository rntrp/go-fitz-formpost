package rest

import (
	"io/ioutil"
	"log"
	"net/http"

	"github.com/rntrp/go-fitz-rest-example/internal/config"
	"github.com/rntrp/go-fitz-rest-example/internal/fitzimg"
)

func handleError(w http.ResponseWriter, err error) bool {
	if err == nil {
		return false
	} else if fitzimg.IsErrorFormatIssue(err) {
		http.Error(w, err.Error(), http.StatusUnprocessableEntity)
	} else {
		log.Println(err)
		http.Error(w, http.StatusText(http.StatusInternalServerError),
			http.StatusInternalServerError)
	}
	return true
}

func handleFileUpload(w http.ResponseWriter, r *http.Request) []byte {
	if !setupFileSizeChecks(w, r) {
		return nil
	}
	memBufSize := coerceMemoryBufferSize(config.GetMemoryBufferSize())
	if err := r.ParseMultipartForm(memBufSize); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return nil
	} else if r.MultipartForm != nil {
		defer r.MultipartForm.RemoveAll()
	}
	f, fh, err := r.FormFile("pdf")
	if err != nil {
		http.Error(w, "File 'pdf' is missing", http.StatusBadRequest)
		return nil
	}
	defer f.Close()
	if fh.Size < MinValidFileSize {
		http.Error(w, "File size too small for a valid document",
			http.StatusBadRequest)
		return nil
	}
	input, err := ioutil.ReadAll(f)
	if err != nil {
		log.Println(err)
		http.Error(w, http.StatusText(http.StatusInternalServerError),
			http.StatusInternalServerError)
		return nil
	}
	return input
}

func setupFileSizeChecks(w http.ResponseWriter, r *http.Request) bool {
	clen, err := coerceContentLength(r.Header.Get("Content-Length"))
	if err == nil && clen < MinValidFileSize {
		http.Error(w, "Content-Length too short for a valid document",
			http.StatusBadRequest)
		return false
	}
	maxReqSize := config.GetMaxRequestSize()
	if maxReqSize >= 0 {
		if err == nil && clen > maxReqSize {
			http.Error(w, "http: Content-Length too large",
				http.StatusRequestEntityTooLarge)
			return false
		}
		r.Body = http.MaxBytesReader(w, r.Body, maxReqSize)
	}
	return true
}
