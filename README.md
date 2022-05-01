# PDF to Image Formpost Microservice
A microservice based on [`go-fitz`](https://github.com/gen2brain/go-fitz), a Go wrapper for the [MuPDF](https://mupdf.com/) Fitz library. Accepts PDF documents via `multipart/form-data` POST requests, delivers images in response.

## Build & Launch
Besides Go 1.14, MuPDF needs to be installed separately.

### Locally
Follow [build instructions](https://github.com/gen2brain/go-fitz) for `go-fitz` with MuPDF. When linking against shared libraries, one may want to set the `CGO_LDFLAGS` environment variable, e.g.:
```bash
export CGO_LDFLAGS="-lmupdf -lm -lmupdf-third -lfreetype -ljbig2dec -lharfbuzz -ljpeg -lopenjp2 -lz" \
    && go build -o /go-fitz-rest
```

Note that `harfbuzz` [isn't listed as a dependency in MuPDF docs](https://mupdf.com/docs/building.html), but is indeed required.

Finally, run the application at port 8080:
```bash
$ go run .
```

### With Docker
```bash
$ docker build --pull --rm -t go-fitz-formpost:latest .
$ docker run --rm -it  -p 8080:8080/tcp go-fitz-formpost:latest
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

Alternatively, there is also an endpoint serving basic HTML form upload page. It is available under http://localhost:8080 or http://localhost:8080/index.html.

### Environment Variables
The application supports configuration via environment variables or a `.env` file. Environment variables have higher priority.

| variable | default | description |
|---|---|---|
| `FITZREST_ENV` | `development` | Currently, this setting only affects the [`.env` file precedence](https://github.com/bkeepers/dotenv#what-other-env-files-can-i-use), no actual distinction between execution environments is made. Possible values are `development`, `test` and `production`. |
| `FITZREST_ENV_DIR` | _empty_ | Path to directory containing the `.env` file. Absolute paths or paths relative to the application folder are possible. If the variable is left _empty_, `.env` file is read from the application folder. |
| `FITZREST_TCP_ADDRESS` | `:8080` | Application TCP address as described by the Golang's `http.Server.Addr` field, most prominently in form of `host:port`. See also [`net` package docs](https://pkg.go.dev/net). |
| `FITZREST_TEMP_DIR` | _empty_ | Path to directory, where applications's temporary _output_ files are managed. Absolute paths or paths relative to the application folder are possible. If the variable is left _empty_, operating system's default temporary directory is used. Please note that the application always writes _input_ multipart data to the OS' temporary directory, regardless of the `FITZREST_TEMP_DIR` value, unless `FITZREST_MEMORY_BUFFER_SIZE` is negative. This limitation is inflicted by the `mime/multipart` Go API and cannot be feasibly altered (yet). |
| `FITZREST_MAX_REQUEST_SIZE` | `-1` | Maximum size of a multipart request in bytes, which can be processed by the application. Note that the request size amounts to the entire HTTP request body, including multipart boundary delimiters, content disposition headers, line breaks and, indeed, the actual payload. Decent clients which send the `Content-Length` header also enjoy fail-fast behavior if that value exceeds the provided maximum. Either way, the application counts bytes during upload and returns `413 Request Entity Too Large` as soon as the limit is exceeded. By default no request size limit is set. |
| `FITZREST_MEMORY_BUFFER_SIZE` | `10485760` | Number of bytes stored in memory when uploading multipart data. If the payload size is exceeding this number, then the remaining bytes are dumped onto the filesystem. Accordingly, the size of `0` prompts the application to always write all bytes to a temporary file, whereas any negative value such as `-1` will prevent the application from hitting the filesystem and retain the whole request payload in memory. Keep in mind that Go `mime/multipart` always adds 10 MiB on top of this value for the "non-file parts", i.e. boundaries etc., hence the actual minimum is 10 MiB plus 1 byte. The default value is therefore effectively 20 MiB. |
| `FITZREST_ENABLE_PROMETHEUS` | `false` | Expose application metrics via [Prometheus](https://prometheus.io) endpoint `/metrics`. |
| `FITZREST_ENABLE_SHUTDOWN_ENDPOINT` | `false` | Enable shutdown endpoint under `/shutdown`. A single POST request with arbitrary payload to this endpoint will cause the application to shutdown gracefully. |
| `FITZREST_SHUTDOWN_TIMEOUT_SECONDS` | `0` | Specifies amount of seconds to wait before ongoing requests are forcefully cancelled in order to perform a graceful shutdown. A zero value lets the application wait indefinetely for all requests to complete. |
| `FITZREST_PROCESSING_MODE` | `serialized` | Choose between `serialized`, `interleaved` and `inmemory` processing modes. |

#### Processing Modes
| value | description |
|---|---|
| `serialized` | Standard processing mode. First, image output is written to a local temporary file. Contents of the file are then transferred as part of the response. The process is repeated for every following page requested. This allows for a single temporary file to be reused, hence minimizing required free disk space. |
| `interleaved` | Slightly enhanced version of the serialized mode. Two temporary files are created, with two goroutines writing images and transferring the contents interchangeably. While the first goroutine is encoding the document page, the second one is busy with serving the image of the previous page. This helps to increase throughput on such systems where the network speed is slower or merely faster than the encoding performance. This comes at a cost of two temporary files needes to be present on disk during processing. If network speed clearly dominates over encoding, then the advantage is insignificant. |
| `inmemory` | Processing is done completely in memory, hence no temporary files are created. May be susceptible to DoS attacks, if the client issues too many requests while artificially cutting down the network throughput, so that many files remain in memory at the same point in time. |

#### Pure In Memory Configuration
In some deployments such as AWS one may want to avoid hitting the filesystem with temporary files. This is possible with the following environment variables set:

* `FITZREST_MEMORY_BUFFER_SIZE=-1`
* `FITZREST_PROCESSING_MODE=inmemory`

Please also consider setting a reasonable request payload size limit to mitigate excessive memory usage:

* `FITZREST_MAX_REQUEST_SIZE=134217728` (128 MiB)

### Parameters
| parameter | mandatory | value |
|---|:---:|---|
| `width` | yes | Min: `1`; Max: `65500` (theoretical). |
| `height` | yes | Min: `1`; Max: `65500` (theoretical). |
| `format` | yes | `jpg`, `jpeg`, `png`, `gif`, `tif`, `tiff`, `bmp`. |
| `pages` | no | Range `m-n` with `1-8388606` being the max theoretical range. If `n` is larger than the actual page count, an HTTP 400 is returned. Special case `pages=m`, equivalent to `pages=m-m`, indicating a single page is allowed. Using underscore for the last page is also possible, e.g. `pages=1-_` for all pages of the document. Default is the first page. |
| `archive` | yes/no | If a page range `m-n` is specified with `m`≠`n`, then providing an archive format is mandatory. `zip` and `tar` are possible values for a ZIP archive and a Tarball respectively. |
| `quality` | no | From `1` to `100` for JPEG (default `95`); `0` to `9` for PNG (default `6`). |
| `resize` | no | Image resizing mode for cases when the source document has a different aspect ratio than specified by the target dimensions. |
| `resample` | no | Resampling algorithm, affects both quality and performance. |

#### Image Resizing Modes
| value | description |
|---|---|
| `fit` | Image is scaled within the specified width-height box ("fit") preserving the aspect ratio of the image. Scaled images with an aspect ratio differing from the specified dimensions have therefore either less width or height. **Note:** if the source image is smaller than the targeted bounding box, the resulting image will retain its source dimensions. Use `fit-upscale` in case the source image should be upscaled. **`fit` is the default mode.** |
| `fit-black` | Same as `fit`, but the left out pixels are filled with black color. Such images have horizontal or vertical black bars on the opposite sides of the image, a practice also known as "letterboxing" or "pillarboxing", respectively. **Note:** if the source image is smaller than the targeted bounding box, the resulting image will be surrounded by black bars, i.e. "windowboxed". Use `fit-upscale-black` if the source image should be upscaled. |
| `fit-white` | Same as `fit-black`, but the bars are now white. |
| `fit-upscale` | Same as `fit`, but upscales the source image so either image width or height fits the targeted bounding box. |
| `fit-upscale-black` | Same as `fit-black`, but upscales the source image avoiding windowboxing. |
| `fit-upscale-white` | Same as `fit-upscale-black`, but the bars are now white. |
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
| `mitchell` | Mitchell–Netravali cubic interpolation, i.e. BC-spline with [ImageMagick](https://imagemagick.org) parameters B=⅓ and C=⅓ |
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
