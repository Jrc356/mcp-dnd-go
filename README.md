# D&D 5e MCP Server (Go)

This project is a Go implementation of an MCP server for D&D 5e API integration, using [github.com/mark3labs/mcp-go](https://github.com/mark3labs/mcp-go).

## Features

- List D&D 5e API categories
- List items in a category
- Get details for an item

## Getting Started

1. Ensure you have Docker and Docker Compose installed.
2. (Optional) Run `make setup` to install Go dependencies locally (for development).
3. Build and start the server using Docker Compose:

   ```sh
   docker compose build
   docker compose up -d
   ```

   Or use the Makefile for convenience:

   ```sh
   make build-and-run
   ```

## Project Structure

- `main.go`: Entry point, MCP server setup
- `config.go`: Configuration constants
- `api.go`: D&D 5e API client logic
- `tools.go`: MCP tool registration and handlers
- `compose.yaml`: Docker Compose configuration
- `Dockerfile`: Multi-stage build for Go MCP server
- `Makefile`: Common build/run/clean commands

## Usage

- To build and run the server:

  ```sh
  docker compose build
  docker compose up -d
  # or
  make build-and-run
  ```

- To stop and remove the server:

  ```sh
  docker compose stop
  docker compose rm -f
  # or
  make clean
  ```

## License

MIT
