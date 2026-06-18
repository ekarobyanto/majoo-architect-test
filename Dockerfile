FROM golang:1.25-alpine AS builder

RUN apk add --no-cache git ca-certificates
WORKDIR /src

# Cache deps
COPY go.mod go.sum ./
RUN go mod download

# Copy sources and build
COPY . .
RUN CGO_ENABLED=0 GOOS=linux GOARCH=amd64 \
    go build -trimpath -ldflags "-s -w" -o /app/api ./cmd/api

FROM alpine:3.18 AS runtime
RUN apk add --no-cache ca-certificates
RUN addgroup -S app && adduser -S -G app app
WORKDIR /app

# Copy binary from builder
COPY --from=builder /app/api /app/api
RUN chown app:app /app/api

USER app
EXPOSE 8080
ENV PORT=8080
ENTRYPOINT ["/app/api"]
