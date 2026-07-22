FROM golang:1.22-alpine AS builder

WORKDIR /build
COPY go.mod ./
RUN go mod download
COPY . .
RUN CGO_ENABLED=0 go build -o /yaml-validator -ldflags="-s -w" .

FROM alpine:3.20
RUN apk add --no-cache ca-certificates
COPY --from=builder /yaml-validator /usr/local/bin/yaml-validator
ENTRYPOINT ["/usr/local/bin/yaml-validator"]
CMD []
