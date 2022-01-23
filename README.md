# Go PDF to Image Microservice Example with go-fitz
Example implementation of a REST microservice based on [`go-fitz`](https://github.com/gen2brain/go-fitz), a Go wrapper for the [MuPDF](https://mupdf.com/) fitz library.

## Build & Launch
Besides Go 1.16, MuPDF needs to be installed separately.

### Locally
Follow [build instructions](https://github.com/gen2brain/go-fitz) for `go-fitz` with MuPDF. When linking against shared libraries, one may want to set the `CGO_LDFLAGS` environment variable, e.g.:
```bash
export CGO_LDFLAGS="-lmupdf -lm -lmupdf-third -lfreetype -ljbig2dec -lharfbuzz -ljpeg -lopenjp2 -lz" \
    && go build -o /go-fitz-rest
```

Note, that `harfbuzz` [isn't listed as a dependency in MuPDF docs](https://mupdf.com/docs/building.html), but is indeed required.

Finally, run the application at port 8080:
```bash
$ go run .
```

### With Docker
```bash
$ docker build --pull --rm -t go-fitz-rest-example:latest .
$ docker run --rm -it  -p 8080:8080/tcp go-fitz-rest-example:latest
```

### With `docker-compose`
```bash
$ docker-compose up
```

## Usage
Send a `multipart/form-data` request with a single file named `pdf` to the `/scale` endpoint. Width and height of the target image is set via the URL query parameters `width` and `height` respectively. Optionally, image output `format` can be specified (`jpeg`, `png`, `gif`, `tiff` or `bmp`)

```bash
$ curl -F pdf=@/path/to/in.pdf -o /path/to/out.png http://localhost:8080/scale?width=200&height=200&format=png
```

Alternatively, there is also `test.html` with HTML `form` and `input`. While experimenting, just edit the `form action` URL.

### Parameters
| parameter | mandatory | value |
|---|:---:|---|
| `width` | yes | Min: `1`; Max: `65536` (theoretical). |
| `height` | yes | Min: `1`; Max: `65536` (theoretical). |
| `format` | yes | `jpg`, `jpeg`, `png`, `gif`, `tif`, `tiff`, `bmp`. |
| `pages` | no | Range `m-n` with `1-8388606` being the max theoretical range. If `n` is larger than the actual page count, an HTTP 400 is returned. Special case `pages=m`, equivalent to `pages=m-m`, indicating a single page is allowed. Using underscore for the last page is also possible, e.g. `pages=1-_` for all pages of the document. Default is the first page. |
| `archive` | yes/no | If a page range `m-n` is specified with `m`≠`n`, then providing an archive format is mandatory. `zip` and `tar` are possible values for a ZIP archive and a Tarball respectively. |
| `quality` | no | From `1` to `100` for JPEG (default `95`); `0` to `9` for PNG (default `6`). |
| `resize` | no | Image resizing mode for cases when the source document has a different aspect ratio than specified by the target dimensions. |
| `resample` | no | Resampling algorithm, affects both quality and performance. |

#### Image Resizing Modes
| value | description |
|---|---|
| `fit` | Image is scaled within the specified width-height box ("fit") preserving the aspect ratio of the image. Scaled images with an aspect ratio differing from the specified dimensions have therefore either less width or height. **This is the default mode.** |
| `fit-black` | Same as `fit`, but the left out pixels are filled with black color. Such images have horizontal or black bars on the opposite sides of the image, a practice also known as "letterboxing". |
| `fit-white` | Same as `fit-black`, but the bars are now white. |
| `fill` | Image is scaled to fill the entire width-height box preserving the aspect ratio of the image. If the image has a different aspect ratio, than specified by the target dimensions, then the image is cropped around the center of the image from both sides. Taller images are equally cropped at the top and bottom, accordingly wider ones are cropped left and right. |
| `fill-top-left` | Same as `fill`, but taller images are cropped from the top, hence the outsized bottom part is left out. Wider images are cropped from the left. This option may come in handy when processing both portrait and landscape oriented documents. |
| `fill-top` | Same as `fill`, but taller images are cropped from the top, hence the bottom part is left out. Wider images are still cropped at the center. |
| `fill-top-right` | Same as `fill-top-left`, but wider images are cropped from the right. |
| `fill-left` | Same as `fill`, but wider images are cropped from the right. Taller images are still cropped at the center. |
| `fill-right` | Same as `fill-left`, but wider images are cropped only from the right. |
| `fill-bottom-left` | The opposite of `fill-top-right`: Taller images are cropped from the bottom, wider images cropped from the left. |
| `fill-bottom` | The opposite of `fill-top`: Taller images are cropped from the bottom, wider images at the center. |
| `fill-bottom-right` | The opposite of `fill-top-left`: Taller images are cropped from the bottom, wider images cropped from the right. |
| `stretch` | The image is stretched to fill the target width and height without preserving the aspect ratio. |

#### Resampling Algorithms
The following resampling algorithms are provided by the Golang [Imaging library](https://github.com/disintegration/imaging), which has a decent basic overview of available resampling algorithms in its readme file. [Wikipedia](https://en.wikipedia.org/wiki/Image_scaling) and [ImageMagick 6 Docs](https://legacy.imagemagick.org/Usage/filter/) are also a good starting point for understanding the differences between particular algorithms and performance implications.

| value | description |
|---|---|
| `box` | Box sampling (**default**) |
| `nearest` | Nearest-neighbor interpolation |
| `linear` | Linear interpolation |
| `hermite` | Hermite cubic interpolation, i.e. BC-spline with B=0 and C=0 |
| `mitchell` | Mitchell-Netravali cubic interpolation, i.e. BC-spline with [ImageMagick](https://imagemagick.org) parameters B=⅓ and C=⅓ |
| `catmull` | Catmull–Rom spline with B=0 and C=½, also used in [GIMP](https://gimp.org) |
| `bspline` | B-spline, i.e. cubic interpolation with [Paint.NET](https://www.getpaint.net) parameters B=1 and C=0 |
| `bartlett` | Bartlett window sinc resampling (3 lobes) |
| `lanczos` | Lanczos sinc approximation (3 lobes) |
| `hann` | Hann window sinc resampling (3 lobes) |
| `hamming` | Hamming window sinc resampling (3 lobes) |
| `blackman` | Blackman window sinc resampling (3 lobes) |
| `welch` | Welch parabolic window sinc resampling (3 lobes) |
| `cosine` | Cosine window sinc resampling (3 lobes) |

## License
Inherits AGPLv3 from `go-fitz` and MuPDF.
