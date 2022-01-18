package fitzimg

import (
	"archive/tar"
	"archive/zip"
	"bufio"
	"bytes"
	"image/draw"
	"io"

	"github.com/gen2brain/go-fitz"
)

func inMemory(src []byte, dst io.Writer, params *Params) error {
	out, closer := initArchive(params.Archive, dst)
	buf := new(bytes.Buffer)
	bkg := background(params.Width, params.Height, params.Resize)
	from := params.FirstPage
	to := params.LastPage
	for page := from; page <= to; page++ {
		buf.Reset()
		if err := direct(src, buf, bkg, out, page, params); err != nil {
			return err
		}
	}
	if closer != nil {
		return closer.Close()
	}
	return nil
}

func direct(src []byte, buf *bytes.Buffer, bkg draw.Image, dst interface{}, page int, params *Params) error {
	doc, err := fitz.NewFromMemory(src)
	if err != nil {
		return err
	}
	defer doc.Close()
	switch params.Archive {
	case Tar:
		return directTar(doc, buf, bkg, dst.(*tar.Writer), page, params)
	case Zip:
		return directZip(doc, bkg, dst.(*zip.Writer), page, params)
	default:
		return encode(doc, bkg, dst.(io.Writer), page, params)
	}
}

func directTar(doc *fitz.Document, buf *bytes.Buffer, bkg draw.Image, dst *tar.Writer, page int, params *Params) error {
	w := bufio.NewWriter(buf)
	if err := encode(doc, bkg, w, page, params); err != nil {
		return err
	}
	hdr := &tar.Header{
		Name: name(page, params.Format),
		Mode: 0600,
		Size: int64(buf.Len()),
	}
	if err := dst.WriteHeader(hdr); err != nil {
		return err
	} else if _, err := buf.WriteTo(dst); err != nil {
		return err
	}
	return nil
}

func directZip(doc *fitz.Document, bkg draw.Image, dst *zip.Writer, page int, params *Params) error {
	w, err := dst.Create(name(page, params.Format))
	if err != nil {
		return err
	} else if err := encode(doc, bkg, w, page, params); err != nil {
		return err
	}
	return nil
}
