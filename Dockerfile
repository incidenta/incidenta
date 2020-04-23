FROM golang:1.14.2 as builder
WORKDIR /go/src/github.com/incidenta/incidenta
COPY . .
RUN go build ./cmd/incidenta

FROM debian:10-slim
COPY --from=builder /go/src/github.com/incidenta/incidenta/incidenta /usr/local/bin/
COPY --from=builder /go/src/github.com/incidenta/incidenta/static /shared/app
