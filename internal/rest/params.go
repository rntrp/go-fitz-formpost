package rest

import (
	"fmt"
	"math"
	"strconv"
	"strings"

	"github.com/disintegration/imaging"
	"github.com/rntrp/go-fitz-formpost/internal/fitzimg"
)

// Netpbm formates have the smallest valid file size of 8 bytes
// See https://github.com/mathiasbynens/small
const minValidFileSize = 8

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
const maxImageDim = 65_500

func coerceWidth(dim string) (int, error) {
	errFormat := "supported width range is [%d;%d], got %d"
	return parseInt(1, maxImageDim, 256, dim, errFormat)
}

func coerceHeight(dim string) (int, error) {
	errFormat := "supported height range is [%d;%d], got %d"
	return parseInt(1, maxImageDim, 256, dim, errFormat)
}

// Theoretical max number of pages presentable by a 32-bit PDF reader app
// https://community.adobe.com/t5/acrobat-discussions/is-there-a-pdf-size-limit/m-p/4387327#M12286
const maxPage = 8_388_606

const lastPage = -1

func coerceContentLength(contentLength string) (int64, error) {
	return strconv.ParseInt(contentLength, 10, 64)
}

const maxMemoryBufferSize = int64(math.MaxInt64) - 1

func coerceMemoryBufferSize(memoryBufferSize int64) int64 {
	if memoryBufferSize < 0 || memoryBufferSize > maxMemoryBufferSize {
		return maxMemoryBufferSize
	}
	return memoryBufferSize
}

func coerceFrom(dim string) (int, error) {
	errFormat := "start page index range is [%d;%d], got %d"
	return parseInt(1, maxPage, 1, dim, errFormat)
}

func coerceTo(dim string) (int, error) {
	errFormat := "end page index range is [%d;%d], got %d"
	return parseInt(1, maxPage, lastPage, dim, errFormat)
}

func coerceQuality(format imaging.Format, quality string) (int, error) {
	switch format {
	case imaging.JPEG:
		errFormat := "supported JPEG quality range is [%d;%d], got %d"
		return parseInt(1, 100, 95, quality, errFormat)
	case imaging.PNG:
		errFormat := "supported PNG compression level range is [%d;%d], got %d"
		return parseInt(0, 9, 6, quality, errFormat)
	default:
		return 0, nil
	}
}

func parseInt(min, max, def int, num, errFormat string) (int, error) {
	if len(num) == 0 {
		return def, nil
	}
	n, err := strconv.Atoi(num)
	if err != nil {
		return n, err
	} else if min > n || n > max {
		return n, fmt.Errorf(errFormat, min, max, n)
	}
	return n, nil
}

func coerceArchive(archive string) (fitzimg.Archive, error) {
	switch strings.ToLower(archive) {
	case "":
		return fitzimg.Raw, nil
	case "tar":
		return fitzimg.Tar, nil
	case "zip":
		return fitzimg.Zip, nil
	default:
		return fitzimg.Raw, fmt.Errorf("unknown archive format: %s", archive)
	}
}

func coerceFormat(format string) (imaging.Format, error) {
	return imaging.FormatFromExtension(format)
}

func coerceResize(resize string) (fitzimg.Resize, error) {
	switch strings.ToLower(resize) {
	case "", "fit":
		return fitzimg.Fit, nil
	case "fit-black":
		return fitzimg.FitBlack, nil
	case "fit-white":
		return fitzimg.FitWhite, nil
	case "fit-upscale":
		return fitzimg.FitUpscale, nil
	case "fit-upscale-black":
		return fitzimg.FitUpscaleBlack, nil
	case "fit-upscale-white":
		return fitzimg.FitUpscaleWhite, nil
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
		return fitzimg.Fit, fmt.Errorf("unknown resize mode: %s", resize)
	}
}

func coerceResample(resample string) (imaging.ResampleFilter, error) {
	switch strings.ToLower(resample) {
	case "", "box":
		return imaging.Box, nil
	case "nearest":
		return imaging.NearestNeighbor, nil
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
		return imaging.Box, fmt.Errorf("unknown resampling algorithm: %s", resample)
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
