FROM golang:1.17-alpine3.14 AS builder
RUN apk add --no-cache \
        build-base \
        freetype-dev \
        harfbuzz-dev \
        jbig2dec-dev \
        jpeg-dev \
        mupdf mupdf-dev \
        openjpeg-dev \
        zlib-dev
WORKDIR /app
COPY go.mod ./
COPY go.sum ./
COPY main.go ./
COPY internal ./internal
RUN export CGO_LDFLAGS="-lmupdf -lm -lmupdf-third -lfreetype -ljbig2dec -lharfbuzz -ljpeg -lopenjp2 -lz" \
    && go mod download \
    && go test ./... \
    && go build -o /go-fitz-formpost

FROM alpine:3.14
RUN apk add --no-cache \
        fontconfig \
        ghostscript-fonts \
        msttcorefonts-installer \
    && update-ms-fonts && fc-cache -f \
    && apk add --no-cache \
        freetype \
        harfbuzz \
        jbig2dec \
        jpeg \
        mupdf \
        openjpeg \
        zlib
COPY --from=builder /go-fitz-formpost ./
EXPOSE 8080
CMD [ "/go-fitz-formpost" ]
