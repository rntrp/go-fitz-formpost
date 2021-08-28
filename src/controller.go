package main

import (
	"fmt"
	"image"
	"image/draw"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"runtime"
	"strconv"
	"strings"

	"github.com/disintegration/imaging"
	"github.com/gen2brain/go-fitz"
)

const maxMemory int64 = 1024 * 64
const maxFileSize int64 = 1024 * 1024 * 512

func Welcome(w http.ResponseWriter, r *http.Request) {
	w.Write([]byte("welcome"))
}

func Scale(w http.ResponseWriter, r *http.Request) {
	printMemUsage()
	query := r.URL.Query()
	width, _ := strconv.Atoi(query.Get("width"))
	height, _ := strconv.Atoi(query.Get("height"))
	format := getFormat(query.Get("format"))
	if err := r.ParseMultipartForm(maxMemory); err != nil {
		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte(err.Error()))
		return
	}
	f, fh, err := r.FormFile("pdf")
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte("File 'pdf' could not be processed."))
		return
	}
	defer f.Close()
	if fh.Size > maxFileSize {
		w.WriteHeader(http.StatusRequestEntityTooLarge)
		fmt.Fprintf(w, "Max file size is %d", maxFileSize)
		return
	}
	doc, err := fitz.NewFromReader(f)
	if err != nil {
		log.Println(err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	defer doc.Close()
	tmp, err := ioutil.TempFile("", "fitz.*.png")
	if err != nil {
		log.Println(err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	defer os.Remove(tmp.Name())
	page1, err := doc.ImageDPI(0, 72.0)
	if err != nil {
		log.Println(err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	resized := imaging.Fit(page1, width, height, imaging.Box)
	background := image.NewNRGBA(image.Rect(0, 0, width, height))
	draw.Draw(background, background.Bounds(), image.White, image.ZP, draw.Src)
	out := imaging.OverlayCenter(background, resized, 1.0)
	imaging.Encode(w, out, format,
		imaging.JPEGQuality(99), imaging.PNGCompressionLevel(9))
	printMemUsage()
}

func getFormat(format string) imaging.Format {
	switch (strings.ToLower(format)) {
	case "png": return imaging.PNG
	case "gif": return imaging.GIF
	case "tif", "tiff": return imaging.TIFF
	case "bmp", "dib": return imaging.BMP
	default: return imaging.JPEG
	}
}

func printMemUsage() {
	var m runtime.MemStats
	runtime.ReadMemStats(&m)
	bToMiB := func (b uint64) string {
		return strconv.FormatFloat(float64(b) / 1_048_576, 'f', 3, 64)
	}
	fmt.Printf("Alloc = %v MiB" +
		"\tTotalAlloc = %v MiB" +
		"\tSys = %v MiB" +
		"\tNumGC = %v\n",
		bToMiB(m.Alloc),
		bToMiB(m.TotalAlloc),
		bToMiB(m.Sys),
		m.NumGC)
}
