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
$ cd src
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

## License
Inherits AGPLv3 from `go-fitz` and MuPDF.
