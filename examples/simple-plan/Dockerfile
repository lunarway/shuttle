FROM golang:1.10.3-alpine as builder

COPY . /go

RUN go build -v -o app

FROM alpine

COPY --from=builder /go/app /app

ENTRYPOINT [ "/app" ]