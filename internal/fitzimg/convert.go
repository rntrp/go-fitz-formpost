package fitzimg

import (
	"errors"
	"fmt"
	"io"

	"github.com/gen2brain/go-fitz"
	"github.com/rntrp/go-fitz-formpost/internal/config"
)

func IsErrorFormatIssue(err error) bool {
	return errors.Is(err, fitz.ErrOpenDocument) ||
		errors.Is(err, fitz.ErrNeedsPassword)
}

func NumPage(src []byte) (int, error) {
	doc, err := fitz.NewFromMemory(src)
	if err != nil {
		return 0, fmt.Errorf("fitzimg.NumPage: %w", err)
	}
	defer doc.Close()
	return doc.NumPage(), nil
}

func Convert(src []byte, dst io.Writer, params *Params) error {
	doc, err := fitz.NewFromMemory(src)
	if err != nil {
		return fmt.Errorf("fitzimg.Convert: %w", err)
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
