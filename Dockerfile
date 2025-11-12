FROM golang:1.25-alpine3.22 AS builder
RUN apk add --no-cache build-base upx
WORKDIR /app
COPY . ./
COPY internal ./internal
RUN go mod download \
    && go test -tags musl ./... \
    && go build -tags musl -ldflags="-s -w" -o /go-fitz-formpost \
    && upx --best --lzma /go-fitz-formpost

FROM alpine:3.22
COPY --from=builder /go-fitz-formpost ./
EXPOSE 8080
ENV FITZ_FORMPOST_ENV=production
ENTRYPOINT [ "/go-fitz-formpost" ]
