# Build stage
FROM golang:1.24-bookworm AS builder
WORKDIR /app
COPY . .
RUN go mod download && go build -o mcp-dnd-server .

# Runtime stage (distroless)
FROM gcr.io/distroless/base-debian12:latest
WORKDIR /app
COPY --from=builder /app/mcp-dnd-server .
EXPOSE 8000
ENTRYPOINT ["/app/mcp-dnd-server"]
