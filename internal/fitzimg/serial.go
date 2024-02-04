package fitzimg

import (
	"fmt"
	"io"

	"github.com/gen2brain/go-fitz"
)

func serial(doc *fitz.Document, dst io.Writer, params *Params) error {
	tmp, err := initTmp()
	if err != nil {
		return fmt.Errorf("fitzimg.serial initTmp: %w", err)
	}
	defer removeTmp(tmp)
	out := initArchive(params.Archive, dst)
	bkg := background(params.Width, params.Height, params.Resize)
	from := params.FirstPage
	to := params.LastPage
	for page := from; page <= to; page++ {
		if err := process(doc, bkg, tmp, out, page, params); err != nil {
			return fmt.Errorf("fitzimg.serial page=%d: %w", page, err)
		}
	}
	if err := out.Close(); err != nil {
		return fmt.Errorf("fitzimg.serial out.Close: %w", err)
	}
	return nil
}
