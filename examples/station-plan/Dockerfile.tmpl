FROM golang:{{ get "GO_VERSION" .Args }}-alpine as builder

COPY . /go

RUN go build -v -o app

FROM alpine

COPY --from=builder /go/app /app

ENTRYPOINT [ "/app" ]