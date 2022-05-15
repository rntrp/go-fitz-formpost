package rest

import (
	"fmt"
	"io/ioutil"
	"log"
	"net/http"

	"github.com/rntrp/go-fitz-formpost/internal/config"
	"github.com/rntrp/go-fitz-formpost/internal/fitzimg"
)

func handleNumPage(w http.ResponseWriter, input []byte) (int, bool) {
	total, err := fitzimg.NumPage(input)
	if err == nil {
		return total, true
	} else if fitzimg.IsErrorFormatIssue(err) {
		http.Error(w, err.Error(), http.StatusUnprocessableEntity)
	} else {
		log.Println(err)
		http.Error(w, http.StatusText(http.StatusInternalServerError),
			http.StatusInternalServerError)
	}
	return total, false
}

func handlePageRange(w http.ResponseWriter, from, to, total int, archive fitzimg.Archive) bool {
	var msg string
	switch archive {
	case fitzimg.Raw:
		if from != to && to != LastPage {
			msg = fmt.Sprintf("Missing 'archive' for page range [%d;%d]", from, to)
		} else if from > total {
			msg = fmt.Sprintf("Requested page %d exceeds page count of %d", from, total)
		}
	default:
		if to > total || from > to {
			msg = fmt.Sprintf("Page range [%d;%d] is beyond [1;%d]", from, to, total)
		}
	}
	if len(msg) > 0 {
		http.Error(w, msg, http.StatusBadRequest)
		return false
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

func resolveLastPage(from, to, total int, archive fitzimg.Archive) int {
	switch {
	case to != LastPage:
		return to
	case archive == fitzimg.Raw:
		return from
	default:
		return total
	}
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
