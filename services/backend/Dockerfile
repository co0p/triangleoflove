FROM golang:latest AS builder

WORKDIR /src

COPY go.mod ./
RUN go mod download

COPY . .
RUN CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -o /out/backend ./cmd/server

FROM busybox:musl AS busybox

FROM gcr.io/distroless/static-debian12

WORKDIR /
COPY --from=builder /out/backend /backend
COPY --from=busybox /bin/wget /wget

EXPOSE 8080
HEALTHCHECK --interval=5s --timeout=3s --start-period=10s --retries=3 \
  CMD ["/wget", "-q", "-O", "-", "http://localhost:8080/api/v1/health"]
ENTRYPOINT ["/backend"]
