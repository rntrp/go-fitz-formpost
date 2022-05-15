package rest

import (
	"fmt"
	"log"
	"net/http"

	"github.com/rntrp/go-fitz-formpost/internal/fitzimg"
)

func NumPage(w http.ResponseWriter, r *http.Request) {
	input := handleFileUpload(w, r)
	total, totalOk := handleNumPage(w, input)
	if !totalOk {
		return
	}
	w.Header().Set("Content-Type", "text/plain; charset=utf-8")
	fmt.Fprintf(w, "%d", total)
}

func Convert(w http.ResponseWriter, r *http.Request) {
	query := r.URL.Query()
	width, errW := coerceWidth(query.Get("width"))
	height, errH := coerceHeight(query.Get("height"))
	from, errFr := coerceFrom(query.Get("from"))
	to, errTo := coerceTo(query.Get("to"))
	format, errF := coerceFormat(query.Get("format"))
	archive, errO := coerceArchive(query.Get("archive"))
	quality, errQ := coerceQuality(format, query.Get("quality"))
	resize, errRz := coerceResize(query.Get("resize"))
	resample, errRs := coerceResample(query.Get("resample"))
	for _, err := range [...]error{errW, errH, errFr, errTo, errF, errO, errQ, errRz, errRs} {
		if err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}
	}
	input := handleFileUpload(w, r)
	if input == nil {
		return
	}
	total, totalOk := handleNumPage(w, input)
	if !totalOk || !handlePageRange(w, from, to, total, archive) {
		return
	}
	to = resolveLastPage(from, to, total, archive)
	params := &fitzimg.Params{
		Width:     width,
		Height:    height,
		FirstPage: from - 1,
		LastPage:  to - 1,
		Archive:   archive,
		Format:    format,
		Quality:   quality,
		Resize:    resize,
		Resample:  resample,
	}
	ext, mime := getOutputFileMeta(archive, format)
	w.Header().Set("Content-Disposition", "attachment; filename=result"+ext)
	w.Header().Set("Content-Type", mime)
	if err := fitzimg.Convert(input, w, params); err != nil {
		log.Println(err)
	}
}
