package fitzimg

import (
	"io"

	"github.com/gen2brain/go-fitz"
)

func serial(doc *fitz.Document, dst io.Writer, params *Params) error {
	tmp, err := initTmp()
	if err != nil {
		return err
	}
	defer removeTmp(tmp)
	out, closer := initArchive(params.Archive, dst)
	bkg := background(params.Width, params.Height, params.Resize)
	from := params.FirstPage
	to := params.LastPage
	for page := from; page <= to; page++ {
		if err := process(doc, bkg, tmp, out, page, params); err != nil {
			return err
		}
	}
	if closer != nil {
		return closer.Close()
	}
	return nil
}
