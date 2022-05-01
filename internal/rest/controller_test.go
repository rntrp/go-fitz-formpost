package rest

import (
	"bytes"
	"encoding/base64"
	"encoding/hex"
	"math"
	"net/http/httptest"
	"testing"

	"github.com/rntrp/go-fitz-formpost/internal/config"
)

func init() {
	config.Load()
}

const pdfB64 = "JVBERi0xLjAKMSAwIG9iajw8L1BhZ2VzIDIgMCBSPj5lbmRvYmogMiAwIG9iajw8L0tpZHNbMyAwIFJdL0" +
	"NvdW50IDE+PmVuZG9iaiAzIDAgb2JqPDwvTWVkaWFCb3hbMCAwIDMgM10+PmVuZG9iagp0cmFpbG" +
	"VyPDwvUm9vdCAxIDAgUj4+Cg=="
const pdf3PagesB64 = "JVBERi0xLjAKJeLjz9MKMiAwIG9iago8PCAvVHlwZSAvQ2F0YWxvZyAvUGFnZXMgMyAwIFI+Pgpl" +
	"bmRvYmoKMyAwIG9iago8PCAvVHlwZSAvUGFnZXMgL0tpZHMgWyA0IDAgUiA1IDAgUiA2IDAgUiBd" +
	"IC9Db3VudCAzID4+CmVuZG9iago0IDAgb2JqCjw8IC9NZWRpYUJveCBbMCAwIDMgMyBdICAvVHJp" +
	"bUJveCBbMCAwIDMgMyBdICAvUm90YXRlIDAgID4+CmVuZG9iago1IDAgb2JqCjw8IC9NZWRpYUJv" +
	"eCBbMCAwIDMgMyBdICAvVHJpbUJveCBbMCAwIDMgMyBdICAvUm90YXRlIDAgID4+CmVuZG9iago2" +
	"IDAgb2JqCjw8IC9NZWRpYUJveCBbMCAwIDMgMyBdICAvVHJpbUJveCBbMCAwIDMgMyBdICAvUm90" +
	"YXRlIDAgID4+CmVuZG9iagp4cmVmDQowIDcNCjAwMDAwMDAwMDEgNjU1MzUgZg0KMDAwMDAwMDAw" +
	"MCAwMDAwMCBmDQowMDAwMDAwMDE1IDAwMDAwIG4NCjAwMDAwMDAwNjMgMDAwMDAgbg0KMDAwMDAw" +
	"MDEzNCAwMDAwMCBuDQowMDAwMDAwMjA5IDAwMDAwIG4NCjAwMDAwMDAyODQgMDAwMDAgbg0KdHJh" +
	"aWxlcg0KPDwvU2l6ZSA1IC9JRCBbKJI0YKT5scABPrNkOEg8/MMpICiSNGCk+bHAAT6zZDhIPPzD" +
	"KSBdIC9Sb290IDIgMCBSID4+IA0Kc3RhcnR4cmVmDQozNTkNCiUlRU9GDQo=+PmVuZG9iaiAzIDA" +
	"gb2JqPDwvTWVkaWFCb3hbMCAwIDMgM10+PmVuZG9iagp0cmFpbGVyPDwvUm9vdCAxIDAgUj4+Cg=="

func TestScale(t *testing.T) {
	rec := httptest.NewRecorder()

	pdfBytes, _ := base64.StdEncoding.DecodeString(pdfB64)
	body := []byte("--boundary\r\nContent-Disposition: form-data; name=\"pdf\"; filename=\"pdf\"\r\n\r\n")
	body = append(body, pdfBytes...)
	body = append(body, []byte("\r\n--boundary--")...)
	req := httptest.NewRequest("POST", "/scale?width=2&height=2&page=1&format=png", bytes.NewReader(body))
	req.Header.Add("Content-Type", "multipart/form-data; boundary=\"boundary\"")

	Convert(rec, req)

	if rec.Code != 200 {
		t.Errorf("rec.Code = %v; want 200", rec.Code)
	}

	out := rec.Body.Bytes()
	outMagic := out[:int(math.Min(8, float64(len(out))))]
	pngMagic, _ := hex.DecodeString("89504e470d0a1a0a")
	if !bytes.Equal(outMagic, pngMagic) {
		got := hex.EncodeToString(outMagic)
		t.Errorf("rec.Body PNG Magic Number = %v; want 89504e470d0a1a0a", got)
	}
}

func TestArchiveScale(t *testing.T) {
	rec := httptest.NewRecorder()

	pdfBytes, _ := base64.StdEncoding.DecodeString(pdf3PagesB64)
	body := []byte("--boundary\r\nContent-Disposition: form-data; name=\"pdf\"; filename=\"pdf\"\r\n\r\n")
	body = append(body, pdfBytes...)
	body = append(body, []byte("\r\n--boundary--")...)
	req := httptest.NewRequest("POST", "/scale?width=2&height=2&page=1&format=png", bytes.NewReader(body))
	req.Header.Add("Content-Type", "multipart/form-data; boundary=\"boundary\"")

	Convert(rec, req)

	if rec.Code != 200 {
		t.Errorf("rec.Code = %v; want 200", rec.Code)
	}

	out := rec.Body.Bytes()
	outMagic := out[:int(math.Min(8, float64(len(out))))]
	pngMagic, _ := hex.DecodeString("89504e470d0a1a0a")
	if !bytes.Equal(outMagic, pngMagic) {
		got := hex.EncodeToString(outMagic)
		t.Errorf("rec.Body PNG Magic Number = %v; want 89504e470d0a1a0a", got)
	}
}
