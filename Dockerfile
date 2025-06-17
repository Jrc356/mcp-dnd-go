FROM golang:1.24-bookworm AS builder
WORKDIR /app
COPY go.mod go.sum ./
RUN go mod download
COPY . .
RUN go build -o mcp-dnd-server .

FROM gcr.io/distroless/base-debian12:latest
WORKDIR /app
COPY --from=builder /app/mcp-dnd-server .
EXPOSE 8000
ENTRYPOINT ["/app/mcp-dnd-server"]
