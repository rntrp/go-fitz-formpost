package fitzimg

import (
	"io"

	"github.com/gen2brain/go-fitz"
	"github.com/rntrp/go-fitz-rest-example/internal/config"
)

func NumPage(src []byte) (int, error) {
	doc, err := fitz.NewFromMemory(src)
	if err != nil {
		return 0, err
	}
	defer doc.Close()
	return doc.NumPage(), nil
}

func Scale(src []byte, dst io.Writer, params *Params) error {
	switch config.GetProcessingMode() {
	case config.Interleaved:
		return interleave(src, dst, params)
	case config.InMemory:
		return inMemory(src, dst, params)
	default:
		return serial(src, dst, params)
	}
}
