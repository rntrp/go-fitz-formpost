FROM golang:1.17-alpine3.15 AS builder
RUN apk add --no-cache build-base
WORKDIR /app
COPY go.mod ./
COPY go.sum ./
COPY main.go ./
COPY internal ./internal
RUN go mod download \
    && go test -tags musl ./... \
    && go build -tags musl -o /go-fitz-formpost

FROM alpine:3.15
RUN apk add --no-cache ghostscript-fonts
COPY --from=builder /go-fitz-formpost ./
EXPOSE 8080
CMD [ "/go-fitz-formpost" ]
