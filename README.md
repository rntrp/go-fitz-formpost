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
| `pages` | no | Range `m-n` with `1-8388606` being the max theoretical range. If `n` is larger than the actual page count, an HTTP 400 is returned. Special case `pages=m`, equivalent to `pages=m-m`, indicating a single page is allowed. Using underscore for the last page is also possible, e.g. `pages=1-_` for all pages of the document. |
| `archive` | yes/no | If a page range `m-n` is specified with `m`≠`n`, then providing an archive format is mandatory. `zip` and `tar` are possible values for a ZIP archive and a Tarball respectively. |
| `quality` | no | From `1` to `100` for JPEG (default `95`); `0` to `9` for PNG (default `6`). |
| `resize` | no | Image resizing mode for cases when the source document has a different aspect ratio than specified by the target dimensions. |
| `resample` | no | Resamling algorithm, affection both quality and performance. |

## License
Inherits AGPLv3 from `go-fitz` and MuPDF.
