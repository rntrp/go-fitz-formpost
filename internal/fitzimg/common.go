package fitzimg

import (
	"fmt"
	"image"
	"image/draw"
	"image/png"
	"io"
	"io/ioutil"
	"log"
	"os"
	"strings"

	"github.com/disintegration/imaging"
	"github.com/gen2brain/go-fitz"
	"github.com/rntrp/go-fitz-formpost/internal/config"
)

func initTmp() (string, error) {
	dir := config.GetTempDir()
	f, err := ioutil.TempFile(dir, "fitz.*.image")
	if err != nil {
		return "", err
	}
	tmp := f.Name()
	if err := f.Close(); err != nil {
		os.Remove(tmp)
		return "", err
	}
	return tmp, nil
}

func removeTmp(tmp string) {
	if err := os.Remove(tmp); err != nil {
		log.Println(err)
	}
}

func background(width, height int, resize Resize) *image.NRGBA {
	var color *image.Uniform
	switch resize {
	case FitBlack, FitUpscaleBlack:
		color = image.Black
	case FitWhite, FitUpscaleWhite:
		color = image.White
	default:
		return nil
	}
	bkg := image.NewNRGBA(image.Rect(0, 0, width, height))
	draw.Draw(bkg, bkg.Bounds(), color, image.ZP, draw.Src)
	return bkg
}

func process(src []byte, bkg draw.Image, dst string, out interface{}, page int, params *Params) error {
	if err := dump(src, bkg, dst, page, params); err != nil {
		return err
	}
	n := name(page, params.Format)
	return transfer(dst, out, n, params)
}

func name(page int, format imaging.Format) string {
	return fmt.Sprintf("img%07d.%s", page+1, strings.ToLower(format.String()))
}

func dump(src []byte, bkg draw.Image, dst string, page int, params *Params) error {
	doc, err := fitz.NewFromMemory(src)
	if err != nil {
		return err
	}
	defer doc.Close()
	img, err := os.OpenFile(dst, os.O_WRONLY|os.O_TRUNC|os.O_EXCL, 0600)
	if err != nil {
		return err
	}
	defer img.Close()
	return encode(doc, bkg, img, page, params)
}

func transfer(dst string, out interface{}, name string, params *Params) error {
	img, err := os.OpenFile(dst, os.O_RDONLY, 0400)
	if err != nil {
		return err
	}
	defer img.Close()
	return write(params.Archive, out, img, name)
}

func encode(doc *fitz.Document, bkg draw.Image, dst io.Writer, page int, params *Params) error {
	img, err := doc.ImageDPI(page, 72.0)
	if err != nil {
		return err
	}
	out := resize(img, bkg, params)
	q := quality(params.Format, params.Quality)
	return imaging.Encode(dst, out, params.Format, q...)
}

func quality(format imaging.Format, value int) []imaging.EncodeOption {
	switch format {
	case imaging.JPEG:
		return []imaging.EncodeOption{imaging.JPEGQuality(value)}
	case imaging.PNG:
		var level png.CompressionLevel
		switch value {
		case 0:
			level = png.NoCompression
		case 1, 2, 3:
			level = png.BestSpeed
		case 7, 8, 9:
			level = png.BestCompression
		default:
			level = png.DefaultCompression
		}
		return []imaging.EncodeOption{imaging.PNGCompressionLevel(level)}
	default:
		return nil
	}
}

func resize(img image.Image, bkg draw.Image, params *Params) *image.NRGBA {
	switch params.Resize {
	case FitUpscale:
		return fitUpscale(img, params.Width, params.Height, params.Resample)
	case FitBlack, FitWhite:
		resized := imaging.Fit(img, params.Width, params.Height, params.Resample)
		return imaging.OverlayCenter(bkg, resized, 1.0)
	case FitUpscaleBlack, FitUpscaleWhite:
		resized := fitUpscale(img, params.Width, params.Height, params.Resample)
		return imaging.OverlayCenter(bkg, resized, 1.0)
	case FillTopLeft, FillTop, FillTopRight,
		FillLeft, Fill, FillRight,
		FillBottomLeft, FillBottom, FillBottomRight:
		anchor := resizeFillMap[params.Resize]
		return imaging.Fill(img, params.Width, params.Height, anchor, params.Resample)
	case Stretch:
		return imaging.Resize(img, params.Width, params.Height, params.Resample)
	default:
		return imaging.Fit(img, params.Width, params.Height, params.Resample)
	}
}

func fitUpscale(img image.Image, w, h int, r imaging.ResampleFilter) *image.NRGBA {
	b := img.Bounds()
	if b.Dx() < b.Dy() {
		// portrait images are bound by height
		return imaging.Resize(img, 0, h, r)
	}
	return imaging.Resize(img, w, 0, r)
}
