FROM golang:1.18-alpine3.15 AS builder
RUN apk add --no-cache build-base upx
WORKDIR /app
COPY . ./
COPY internal ./internal
RUN go mod download \
    && go test -tags musl ./... \
    && go build -tags musl -ldflags="-s -w" -o /go-fitz-formpost

FROM alpine:3.15
COPY --from=builder /go-fitz-formpost ./
EXPOSE 8080
ENV FITZ_FORMPOST_ENV production
ENTRYPOINT [ "/go-fitz-formpost" ]
