FROM golang:1.16-alpine
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
        zlib \
    && apk add --no-cache --virtual .build-deps \
        build-base \
        freetype-dev \
        harfbuzz-dev \
        jbig2dec-dev \
        jpeg-dev \
        mupdf-dev \
        openjpeg-dev \
        zlib-dev
WORKDIR /app
COPY go.mod ./
COPY go.sum ./
COPY src/*.go ./
RUN export CGO_LDFLAGS="-lmupdf -lm -lmupdf-third -lfreetype -ljbig2dec -lharfbuzz -ljpeg -lopenjp2 -lz" \
    && go mod download \
    && go test \
    && go build -o /go-fitz-rest \
    && apk del .build-deps \
    && rm -rf /tmp/*
EXPOSE 8080
CMD [ "/go-fitz-rest" ]
