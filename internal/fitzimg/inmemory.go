package fitzimg

import (
	"fmt"
	"image/draw"
	"io"

	"github.com/gen2brain/go-fitz"
)

func inMemory(doc *fitz.Document, dst io.Writer, params *Params) error {
	out := initArchive(params.Archive, dst)
	bkg := background(params.Width, params.Height, params.Resize)
	from := params.FirstPage
	to := params.LastPage
	for page := from; page <= to; page++ {
		if err := entry(doc, bkg, out, page, params); err != nil {
			return fmt.Errorf("fitzimg.inMemory page=%d: %w", page, err)
		}
	}
	if err := out.Close(); err != nil {
		return fmt.Errorf("fitzimg.inMemory out.Close: %w", err)
	}
	return nil
}

func entry(doc *fitz.Document, bkg draw.Image, dst ArchiveWriter, page int, params *Params) error {
	if w, err := dst.StartEntry(name(page, params.Format)); err != nil {
		return fmt.Errorf("StartEntry: %w", err)
	} else if err := encode(doc, bkg, w, page, params); err != nil {
		return fmt.Errorf("encode: %w", err)
	} else if err := dst.FinishEntry(); err != nil {
		return fmt.Errorf("FinishEntry: %w", err)
	}
	return nil
}
