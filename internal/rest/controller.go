package rest

import (
	"fmt"
	"io/ioutil"
	"log"
	"net/http"

	"github.com/rntrp/go-fitz-rest-example/internal/config"
	"github.com/rntrp/go-fitz-rest-example/internal/fitzimg"
)

func NumPage(w http.ResponseWriter, r *http.Request) {
	input := handleFileUpload(w, r)
	if input == nil {
		return
	}
	numPage, err := fitzimg.NumPage(input)
	if err != nil {
		log.Println(err)
		http.Error(w, http.StatusText(http.StatusInternalServerError),
			http.StatusInternalServerError)
		return
	}
	w.Header().Set("Content-Type", "text/plain; charset=utf-8")
	fmt.Fprintf(w, "%d", numPage)
}

func Scale(w http.ResponseWriter, r *http.Request) {
	query := r.URL.Query()
	width, errW := coerceWidth(query.Get("width"))
	height, errH := coerceHeight(query.Get("height"))
	pageStart, pageEnd, errP := coercePages(query.Get("pages"))
	format, errF := coerceFormat(query.Get("format"))
	archive, errO := coerceArchive(query.Get("archive"))
	quality, errQ := coerceQuality(format, query.Get("quality"))
	resize, errRZ := coerceResize(query.Get("resize"))
	resample, errRF := coerceResample(query.Get("resample"))
	for _, err := range []error{errW, errH, errP, errF, errO, errQ, errRZ, errRF} {
		if err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}
	}
	if pageStart != pageEnd && archive == fitzimg.Raw {
		http.Error(w, "Invalid archive format", http.StatusBadRequest)
		return
	}
	input := handleFileUpload(w, r)
	if input == nil {
		return
	}
	numPage, err := fitzimg.NumPage(input)
	if err != nil {
		log.Println(err)
		http.Error(w, http.StatusText(http.StatusInternalServerError),
			http.StatusInternalServerError)
		return
	} else if pageEnd == LastPage {
		pageEnd = numPage
	}
	if err := checkPageRange(pageStart, pageEnd, numPage); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	params := &fitzimg.Params{
		Width:     width,
		Height:    height,
		FirstPage: pageStart - 1,
		LastPage:  pageEnd - 1,
		Archive:   archive,
		Format:    format,
		Quality:   quality,
		Resize:    resize,
		Resample:  resample,
	}
	ext, mime := getOutputFileMeta(archive, format)
	w.Header().Set("Content-Disposition", "attachment; filename=result"+ext)
	w.Header().Set("Content-Type", mime)
	if err := fitzimg.Scale(input, w, params); err != nil {
		log.Println(err)
	}
}

func handleFileUpload(w http.ResponseWriter, r *http.Request) []byte {
	if err := r.ParseMultipartForm(config.GetMemoryBufferSize()); err != nil {
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
	if fh.Size > config.GetMaxFileSize() {
		http.Error(w, fmt.Sprintf("Max file size is %d", config.GetMaxFileSize()),
			http.StatusRequestEntityTooLarge)
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
