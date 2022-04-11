FROM golang:1.18 AS builder

WORKDIR /app
COPY go.mod .
COPY go.sum .
RUN go mod download

COPY main.go .
COPY ./podchaosmonkey/ /app/podchaosmonkey/
COPY ./util/ /app/util/

RUN go test -v ./...
RUN CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -o /podchaosmonkey

FROM alpine:latest
WORKDIR /
COPY --from=builder /podchaosmonkey /podchaosmonkey
ENTRYPOINT ["/podchaosmonkey"]
