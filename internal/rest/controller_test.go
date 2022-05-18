package rest

import (
	"archive/zip"
	"bytes"
	_ "embed"
	"encoding/hex"
	"fmt"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/rntrp/go-fitz-formpost/internal/config"
)

func init() {
	config.Load()
}

//go:embed test.pdf
var pdf []byte

func testInitPostReq(query string) *http.Request {
	body := []byte("--boundary\r\nContent-Disposition: form-data; name=\"pdf\"; filename=\"pdf\"\r\n\r\n")
	body = append(body, pdf...)
	body = append(body, []byte("\r\n--boundary--")...)
	req := httptest.NewRequest("POST", query, bytes.NewReader(body))
	req.Header.Add("Content-Type", "multipart/form-data; boundary=\"boundary\"")
	return req
}

func TestPages(t *testing.T) {
	rec := httptest.NewRecorder()
	req := testInitPostReq("/pages")

	NumPage(rec, req)
	if rec.Code != 200 {
		t.Errorf("rec.Code = %v; want 200", rec.Code)
	}

	pages := string(rec.Body.Bytes())
	if pages != "3" {
		t.Errorf("Number of pages = %v; want 3", pages)
	}
}

func TestConvert(t *testing.T) {
	rec := httptest.NewRecorder()
	req := testInitPostReq("/convert?width=2&height=2&from=1&format=png")

	Convert(rec, req)
	if rec.Code != 200 {
		t.Errorf("rec.Code = %v; want 200", rec.Code)
	}

	outMagic := make([]byte, 8)
	rec.Body.Read(outMagic)
	pngMagic, _ := hex.DecodeString("89504e470d0a1a0a")
	if !bytes.Equal(outMagic, pngMagic) {
		got := hex.EncodeToString(outMagic)
		t.Errorf("rec.Body PNG Magic Number = %v; want 89504e470d0a1a0a", got)
	}
}

func TestConvertArchive(t *testing.T) {
	rec := httptest.NewRecorder()
	req := testInitPostReq("/convert?width=2&height=2&from=1&to=2&archive=zip&format=png")

	Convert(rec, req)

	if rec.Code != 200 {
		t.Errorf("rec.Code = %v; want 200", rec.Code)
	}

	out := rec.Body.Bytes()

	outMagic := out[:4]
	zipMagic, _ := hex.DecodeString("504b0304")
	if !bytes.Equal(outMagic, zipMagic) {
		got := hex.EncodeToString(outMagic)
		t.Errorf("rec.Body ZIP Magic Number = %v; want 504b0304", got)
	}

	zipReader, _ := zip.NewReader(bytes.NewReader(out), int64(len(out)))
	numFiles := len(zipReader.File)
	if numFiles != 2 {
		t.Errorf("Number of zipped files = %v; want 2", numFiles)
	}

	pngMagic, _ := hex.DecodeString("89504e470d0a1a0a")
	for i, zipFile := range zipReader.File {
		gotName := zipFile.Name
		wantName := fmt.Sprintf("img%07d.png", i+1)
		if gotName != wantName {
			t.Errorf("PNG file name = %v; want %v", gotName, wantName)
		}

		gotMagic := make([]byte, 8)
		file, _ := zipFile.Open()
		file.Read(gotMagic)
		file.Close()
		if !bytes.Equal(gotMagic, pngMagic) {
			got := hex.EncodeToString(gotMagic)
			t.Errorf("rec.Body PNG Magic Number = %v; want 89504e470d0a1a0a", got)
		}
	}
}
