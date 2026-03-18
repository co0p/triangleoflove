FROM golang:latest AS builder

WORKDIR /src

COPY go.mod ./
RUN go mod download

COPY . .
RUN CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -o /out/backend ./cmd/server

FROM gcr.io/distroless/static-debian12

WORKDIR /
COPY --from=builder /out/backend /backend

EXPOSE 8080
ENTRYPOINT ["/backend"]
