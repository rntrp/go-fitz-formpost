package rest

import (
	"errors"
	"fmt"
	"math"
	"strconv"
	"strings"

	"github.com/disintegration/imaging"
	"github.com/rntrp/go-fitz-rest-example/internal/fitzimg"
)

// https://stackoverflow.com/q/17279712
const MinValidFileSize = 64

// Theoretical max dimensions according to its respective specs:
//
// BMP (uint32 in go/x/image/bmp): 4,294,967,295
//
// GIF: 65,535
//
// JPEG: 65,535 (65,500 for the libjpeg-turbo based software)
//
// PNG: 4,294,967,295
//
// TIFF (uint32 in go/x/image/tiff): 4,294,967,295
//
// libjpeg-turbo max of 65,500 pixels appears to be a good limit for the other image
// formats, since all known image viewers are getting problems with image dimensions
// higher than this value, or even at values far smaller than this.
const MaxImageDim = 65_500

func coerceWidth(dim string) (int, error) {
	errFormat := "Supported width range is [%d;%d], got %d"
	return parseInt(1, MaxImageDim, 256, dim, errFormat)
}

func coerceHeight(dim string) (int, error) {
	errFormat := "Supported height range is [%d;%d], got %d"
	return parseInt(1, MaxImageDim, 256, dim, errFormat)
}

// First page index
const MinPage = 1

// Theoretical max number of pages presentable by a 32-bit PDF reader app
// https://community.adobe.com/t5/acrobat-discussions/is-there-a-pdf-size-limit/m-p/4387327#M12286
const MaxPage = 8_388_606

// Default page index, namely the first page
const DefaultPage = 1

const LastPage = -1

func coerceContentLength(contentLength string) (int64, error) {
	return strconv.ParseInt(contentLength, 10, 64)
}

const MaxMemoryBufferSize = int64(math.MaxInt64) - 1

func coerceMemoryBufferSize(memoryBufferSize int64) int64 {
	if memoryBufferSize < 0 || memoryBufferSize > MaxMemoryBufferSize {
		return MaxMemoryBufferSize
	}
	return memoryBufferSize
}

func coercePages(dim string) (int, int, error) {
	if len(dim) == 0 {
		return DefaultPage, DefaultPage, nil
	}
	s := strings.SplitN(dim, "-", 2)
	pageErrFormat := "Supported page number range is [%d;%d], got %d"
	start, err := parseInt(MinPage, MaxPage, DefaultPage, s[0], pageErrFormat)
	if err != nil {
		return DefaultPage, DefaultPage, err
	} else if len(s) == 1 {
		return start, start, nil
	} else if s[1] == "_" {
		return start, LastPage, nil
	}
	end, err := parseInt(MinPage, MaxPage, DefaultPage, s[1], pageErrFormat)
	if err != nil {
		return DefaultPage, DefaultPage, err
	}
	if end < start {
		pageRangeErrFormat := "Invalid page range [%d;%d]"
		return start, end, errors.New(fmt.Sprintf(pageRangeErrFormat, start, end))
	}
	return start, end, nil
}

func coerceQuality(format imaging.Format, quality string) (int, error) {
	switch format {
	case imaging.JPEG:
		errFormat := "Supported JPEG quality range is [%d;%d], got %d"
		return parseInt(1, 100, 95, quality, errFormat)
	case imaging.PNG:
		errFormat := "Supported PNG compression level range is [%d;%d], got %d"
		return parseInt(0, 9, 6, quality, errFormat)
	default:
		return 0, nil
	}
}

func parseInt(min, max, def int, num, errFormat string) (int, error) {
	if num == "" {
		return def, nil
	}
	n, err := strconv.Atoi(num)
	if err != nil {
		return n, err
	}
	if min > n || n > max {
		return n, errors.New(fmt.Sprintf(errFormat, min, max, n))
	}
	return n, nil
}

func checkPageRange(first, last, numPages int) error {
	if first < 1 || last > numPages || first > last {
		msg := "Page range [%d,%d] is beyond [1,%d]"
		return errors.New(fmt.Sprintf(msg, first, last, numPages))
	}
	return nil
}

func coerceArchive(archive string) (fitzimg.Archive, error) {
	if len(archive) == 0 {
		return fitzimg.Raw, nil
	}
	s := strings.ToLower(archive)
	switch s {
	case "tar":
		return fitzimg.Tar, nil
	case "zip":
		return fitzimg.Zip, nil
	default:
		unknownArchiveFormat := "Unknown archive format: %s"
		return fitzimg.Raw, errors.New(fmt.Sprintf(unknownArchiveFormat, archive))
	}
}

func coerceFormat(format string) (imaging.Format, error) {
	return imaging.FormatFromExtension(format)
}

func coerceResize(resize string) (fitzimg.Resize, error) {
	if len(resize) == 0 {
		return fitzimg.Fit, nil
	}
	s := strings.ToLower(resize)
	switch s {
	case "fit":
		return fitzimg.Fit, nil
	case "fit-black":
		return fitzimg.FitBlack, nil
	case "fit-white":
		return fitzimg.FitWhite, nil
	case "fill":
		return fitzimg.Fill, nil
	case "fill-top-left":
		return fitzimg.FillTopLeft, nil
	case "fill-top":
		return fitzimg.FillTop, nil
	case "fill-top-right":
		return fitzimg.FillTopRight, nil
	case "fill-left":
		return fitzimg.FillLeft, nil
	case "fill-right":
		return fitzimg.FillRight, nil
	case "fill-bottom-left":
		return fitzimg.FillBottomLeft, nil
	case "fill-bottom":
		return fitzimg.FillBottom, nil
	case "fill-bottom-right":
		return fitzimg.FillBottomRight, nil
	case "stretch":
		return fitzimg.Stretch, nil
	default:
		unknownResizeFormat := "Unknown resize mode: %s"
		return fitzimg.Fit, errors.New(fmt.Sprintf(unknownResizeFormat, resize))
	}
}

func coerceResample(resample string) (imaging.ResampleFilter, error) {
	if len(resample) == 0 {
		return imaging.Box, nil
	}
	s := strings.ToLower(resample)
	switch s {
	case "nearest":
		return imaging.NearestNeighbor, nil
	case "box":
		return imaging.Box, nil
	case "linear":
		return imaging.Linear, nil
	case "hermite":
		return imaging.Hermite, nil
	case "mitchell":
		return imaging.MitchellNetravali, nil
	case "catmull":
		return imaging.CatmullRom, nil
	case "bspline":
		return imaging.BSpline, nil
	case "bartlett":
		return imaging.Bartlett, nil
	case "lanczos":
		return imaging.Lanczos, nil
	case "hann":
		return imaging.Hann, nil
	case "hamming":
		return imaging.Hamming, nil
	case "blackman":
		return imaging.Blackman, nil
	case "welch":
		return imaging.Welch, nil
	case "cosine":
		return imaging.Cosine, nil
	default:
		unknownResampleFormat := "Unknown resampling algorithm: %s"
		return imaging.Box, errors.New(fmt.Sprintf(unknownResampleFormat, resample))
	}
}

func getOutputFileMeta(output fitzimg.Archive, format imaging.Format) (string, string) {
	switch output {
	case fitzimg.Tar:
		return ".tar", "application/x-tar"
	case fitzimg.Zip:
		return ".zip", "application/zip"
	default:
		switch format {
		case imaging.JPEG:
			return ".jpg", "image/jpeg"
		case imaging.PNG:
			return ".png", "image/png"
		case imaging.GIF:
			return ".gif", "image/gif"
		case imaging.TIFF:
			return ".tif", "image/tiff"
		case imaging.BMP:
			return ".bmp", "image/bmp"
		default:
			return "", "application/octet-stream"
		}
	}
}
