[![Release](https://img.shields.io/github/v/release/rntrp/go-fitz-formpost)](https://github.com/rntrp/go-fitz-formpost/releases)
[![Go Version](https://img.shields.io/github/go-mod/go-version/rntrp/go-fitz-formpost)](https://go.dev/dl/)
[![Go Report Card](https://goreportcard.com/badge/github.com/rntrp/go-fitz-formpost)](https://goreportcard.com/report/github.com/rntrp/go-fitz-formpost)
[![Docker Image](https://img.shields.io/docker/image-size/rntrp/go-fitz-formpost/latest?logo=docker)](https://hub.docker.com/r/rntrp/go-fitz-formpost)
[![Tests](https://github.com/rntrp/go-fitz-formpost/actions/workflows/tests.yml/badge.svg)](https://github.com/rntrp/go-fitz-formpost/actions/workflows/tests.yml)

# Document to Image Formpost Microservice
A microservice based on [`go-fitz`](https://github.com/gen2brain/go-fitz), a Go wrapper for the [MuPDF](https://mupdf.com/) Fitz library. Accepts PDF and EPUB documents via `multipart/form-data` POST requests, delivers images in response.

## Build & Launch
Besides Go 1.21, MuPDF needs to be installed separately.

### Locally
Follow [build instructions](https://github.com/gen2brain/go-fitz) for `go-fitz` with MuPDF. When linking against shared libraries, one may want to set the `CGO_LDFLAGS` environment variable.

`go-fitz` provides [prebuilt libraries](https://github.com/gen2brain/go-fitz/tree/master/libs) for common operating systems and acrhitectures. Using those on Ubuntu, macOS or Windows makes building straightforward:
```bash
$ go mod download
$ go build -o /go-fitz-formpost
```
On Alpine Linux this would be:
```bash
$ go build -tags musl -o /go-fitz-formpost
```
Or you can also build and run the application at port 8080:
```bash
$ go run .
```
Alpine still requires the tag:
```bash
$ go run -tags musl .
```
If you are updating the dependencies locally, you may want to update `go.sum`:
```bash
$ go mod download
$ go mod tidy
```
> [!IMPORTANT]
> **Issues with vendoring**: Unfortunately `go mod vendor` does not pull the header files and binaries from the `go-fitz` repository. This issue is discussed at golang/go/issues/26366. Please stay with `go mod download` instead.

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
Send a `multipart/form-data` request with a single file named `doc` to the `/convert` endpoint. Width and height of the target image is set via the URL query parameters `width` and `height` respectively. Optionally, image output `format` can be specified (`jpeg`, `png`, `gif`, `tiff` or `bmp`)

```bash
$ curl -F doc=@/path/to/in.pdf -o /path/to/out.png http://localhost:8080/convert?width=200&height=200&format=png
```

Alternatively, there is also an endpoint serving basic HTML form upload page. It is available under http://localhost:8080 or http://localhost:8080/index.html.

### Endpoint Overview
| path | method | description |
|---|---|---|
| `/` | `GET` | Returns a simple HTML page with form uploads for testing `/convert` and `/pages` endpoints. |
| `/index.html` | `GET` | Same as `/`. |
| `/live` | `GET` | Liveness probe endpoint. Returns HTTP 200 while the application is running. |
| `/metrics` | `GET` | Returns application metrics in Prometheus text-based format collected by the [Prometheus Go client library](https://github.com/prometheus/client_golang). The endpoint is disabled by default; enabled with `FITZ_FORMPOST_ENABLE_PROMETHEUS=true` (see [Environment Variables](#environment-variables) below). |
| `/convert` | `POST` | Pivotal endpoint within the whole application. Accepts a single document per `multipart/form-data` request. Converts the document to the requested image format. Target dimensions are specified as part of the URL query string. See [Parameters](#parameters) below for further details. |
| `/pages` | `POST` | Accepts a single document per `multipart/form-data` request. Returns page count of the document as `text/plain; charset=utf-8`. |
| `/shutdown` | `POST` | Initiates graceful shutdown of the application when triggered by a POST request with an arbitrary payload. Returns HTTP 204 after the shutdown process has been started. The endpoint is disabled by default; enabled with `FITZ_FORMPOST_ENABLE_SHUTDOWN_ENDPOINT=true` (see [Environment Variables](#environment-variables) below). |

### Environment Variables
The application supports configuration via environment variables or a `.env` file. Environment variables have higher priority.

| variable | default | description |
|---|---|---|
| `FITZ_FORMPOST_ENV` | `development` | Currently, this setting only affects the [`.env` file precedence](https://github.com/bkeepers/dotenv#what-other-env-files-can-i-use), no actual distinction between execution environments is made. Possible values are `development`, `test` and `production`. |
| `FITZ_FORMPOST_ENV_DIR` | _empty_ | Path to directory containing the `.env` file. Absolute paths or paths relative to the application folder are possible. If the variable is left _empty_, `.env` file is read from the application folder. |
| `FITZ_FORMPOST_TCP_ADDRESS` | `:8080` | Application TCP address as described by the Golang's `http.Server.Addr` field, most prominently in form of `host:port`. See also [`net` package docs](https://pkg.go.dev/net). |
| `FITZ_FORMPOST_TEMP_DIR` | _OS temp folder_ | Path to directory, where applications's temporary _output_ files are managed. Absolute paths or paths relative to the application folder are possible. If the variable is left _empty_, operating system's default temporary directory is used. Please note that the application always writes _input_ multipart data to the OS' temporary directory, regardless of the `FITZ_FORMPOST_TEMP_DIR` value, unless `FITZ_FORMPOST_MEMORY_BUFFER_SIZE` is negative. This limitation is inflicted by the `mime/multipart` Go API and cannot be feasibly altered (yet). |
| `FITZ_FORMPOST_MAX_REQUEST_SIZE` | `-1` | Maximum size of a multipart request in bytes, which can be processed by the application. Note that the request size amounts to the entire HTTP request body, including multipart boundary delimiters, content disposition headers, line breaks and, indeed, the actual payload. Decent clients which send the `Content-Length` header also enjoy fail-fast behavior if that value exceeds the provided maximum. Either way, the application counts bytes during upload and returns `413 Request Entity Too Large` as soon as the limit is exceeded. By default no request size limit is set. |
| `FITZ_FORMPOST_MEMORY_BUFFER_SIZE` | `10485760` | Number of bytes stored in memory when uploading multipart data. If the payload size is exceeding this number, then the remaining bytes are dumped onto the filesystem. Accordingly, the size of `0` prompts the application to always write all bytes to a temporary file, whereas any negative value such as `-1` will prevent the application from hitting the filesystem and retain the whole request payload in memory. Keep in mind that Go `mime/multipart` always adds 10 MiB on top of this value for the "non-file parts", i.e. boundaries etc., hence the actual minimum is 10 MiB plus 1 byte. The default value is therefore effectively 20 MiB. |
| `FITZ_FORMPOST_ENABLE_PROMETHEUS` | `false` | Expose application metrics via [Prometheus](https://prometheus.io) endpoint `/metrics`. |
| `FITZ_FORMPOST_ENABLE_SHUTDOWN_ENDPOINT` | `false` | Enable shutdown endpoint under `/shutdown`. A single POST request with arbitrary payload to this endpoint will cause the application to shutdown gracefully. |
| `FITZ_FORMPOST_SHUTDOWN_TIMEOUT` | `0s` | Specifies amount of time to wait before ongoing requests are forcefully cancelled in order to perform a graceful shutdown. A zero value lets the application wait indefinetely for all requests to complete. At least one time unit must be specified, e.g. `45s` or `5m15s123ms`. See [Go `time.ParseDuration` format](https://pkg.go.dev/time#ParseDuration) for further details. |
| `FITZ_FORMPOST_PROCESSING_MODE` | `serialized` | Choose between `serialized`, `interleaved` and `inmemory` processing modes. See [Processing Modes](#processing-modes) for further details. |
| `FITZ_FORMPOST_RENDERING_DPI` | `300` | Sets the default DPI resolution at which the source PDF document is rasterized. Higher value improves font and image quality (both vector and large enough raster images) but also increases memory and CPU usage. Lower DPI values may result in upscaling of the intermediate raster image, hence nullifying potential quality surplus of higher target dimensions. Any `float64` value between `1` and `math.MaxFloat64` is allowed. Traditional values are `72`, `150`, `200`, `300`, `600`, `1200` and so on. |

#### Processing Modes
| value | description |
|---|---|
| `serialized` | Standard processing mode. First, image output is written to a local temporary file. Contents of the file are then transferred as part of the response. The process is repeated for every following page requested. This allows for a single temporary file to be reused, hence minimizing required free disk space. |
| `interleaved` | Slightly enhanced version of the serialized mode. Two temporary files are created, with two goroutines writing images and transferring the contents interchangeably. While the first goroutine is encoding the document page, the second one is busy with serving the image of the previous page. This helps to increase throughput on such systems where the network speed is slower or merely faster than the encoding performance. This comes at a cost of two temporary files needes to be present on disk during processing. If network speed clearly dominates over encoding, then the advantage is insignificant. |
| `inmemory` | Processing is done completely in memory, hence no temporary files are created. May be susceptible to DoS attacks, if the client issues too many requests while artificially cutting down the network throughput, so that many files remain in memory at the same point in time. |

#### Pure In Memory Configuration
In some deployments such as AWS one may want to avoid hitting the filesystem with temporary files. This is possible with the following environment variables set:

* `FITZ_FORMPOST_MEMORY_BUFFER_SIZE=-1`
* `FITZ_FORMPOST_PROCESSING_MODE=inmemory`

Please also consider setting a reasonable request payload size limit to mitigate excessive memory usage:

* `FITZ_FORMPOST_MAX_REQUEST_SIZE=134217728` (128 MiB)

### Parameters
| parameter | mandatory | value |
|---|:---:|---|
| `width` | yes | Min: `1`; Max: `65500` (theoretical). |
| `height` | yes | Min: `1`; Max: `65500` (theoretical). |
| `format` | yes | `jpg`, `jpeg`, `png`, `gif`, `tif`, `tiff`, `bmp`. |
| `from` | no | Min: `1`; Max: `8388606` (theoretical). This parameter specifies either the sole page of the document to generate image from, or the first page of a page range if the parameter `archive` is provided. Default is `1`. |
| `to` | no | Min: `1`; Max: `8388606` (theoretical). The last page of a page range. Must be larger or equal to `from` if provided. Setting this parameter also makes `archive` mandatory. Default is the last page of the document. |
| `archive` | yes/no | Specifies how the images of the specified page range should be bundled together. This parameter becomes mandatory, if `from` is provided. `zip` and `tar` are possible values for a ZIP archive and a tarball respectively. Images are placed directly within the root directory having filename pattern `imgNNNNNNN.ext`.  |
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
