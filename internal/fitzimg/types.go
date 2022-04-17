package fitzimg

import (
	"github.com/disintegration/imaging"
)

type Archive int

// Archive types
const (
	Raw Archive = iota
	Tar
	Zip
)

type Resize int

// Resize
const (
	Fit Resize = iota
	FitBlack
	FitWhite
	FitUpscale
	FitUpscaleBlack
	FitUpscaleWhite
	Fill
	FillTopLeft
	FillTop
	FillTopRight
	FillLeft
	FillRight
	FillBottomLeft
	FillBottom
	FillBottomRight
	Stretch
)

var resizeFillMap = map[Resize]imaging.Anchor{
	Fill:            imaging.Center,
	FillTopLeft:     imaging.TopLeft,
	FillTop:         imaging.Top,
	FillTopRight:    imaging.TopRight,
	FillLeft:        imaging.TopLeft,
	FillRight:       imaging.Right,
	FillBottomLeft:  imaging.BottomLeft,
	FillBottom:      imaging.Bottom,
	FillBottomRight: imaging.BottomRight,
}

type Params struct {
	Width, Height       int
	FirstPage, LastPage int
	Archive             Archive
	Format              imaging.Format
	Quality             int
	Resize              Resize
	Resample            imaging.ResampleFilter
}
