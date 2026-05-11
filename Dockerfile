FROM golang:1.26-alpine AS builder

WORKDIR /app

RUN apk add --no-cache git

COPY go.mod go.sum ./
RUN go mod download

COPY cmd/ cmd/
COPY internal/ internal/
COPY docs/ docs/

# Build static binary
RUN CGO_ENABLED=0 GOOS=linux GOARCH=amd64 \
    go build -ldflags="-s -w" -o app ./cmd/server

FROM gcr.io/distroless/base-debian12 AS production

WORKDIR /app

COPY --from=builder /app/app .
COPY migration migration/

USER nonroot:nonroot

EXPOSE 8080

ENTRYPOINT ["/app/app"]
