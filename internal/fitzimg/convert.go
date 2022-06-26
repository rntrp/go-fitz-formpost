package fitzimg

import (
	"errors"
	"io"
	"io/ioutil"

	"github.com/gen2brain/go-fitz"
	"github.com/rntrp/go-fitz-formpost/internal/config"
)

func IsErrorFormatIssue(err error) bool {
	return errors.Is(err, fitz.ErrOpenDocument) ||
		errors.Is(err, fitz.ErrNeedsPassword)
}

func NumPage(src []byte) (int, error) {
	defer preventGC(src)
	doc, err := fitz.NewFromMemory(src)
	if err != nil {
		return 0, err
	}
	defer doc.Close()
	return doc.NumPage(), nil
}

func Convert(src []byte, dst io.Writer, params *Params) error {
	defer preventGC(src)
	doc, err := fitz.NewFromMemory(src)
	if err != nil {
		return err
	}
	defer doc.Close()
	switch config.GetProcessingMode() {
	case config.Interleaved:
		return interleave(doc, dst, params)
	case config.InMemory:
		return inMemory(doc, dst, params)
	default:
		return serial(doc, dst, params)
	}
}

// Prevent premature GC of the underlying byte array.
// Workaround for https://github.com/gen2brain/go-fitz/issues/55
func preventGC(b []byte) {
	// Try to avoid optimization by performing some bogus logic
	if len(b) > 0 {
		// Write the first byte to /dev/null
		ioutil.Discard.Write(b[:1])
	}
}
