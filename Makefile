build:
	go build -o mcp-dnd-server .

run:
	go run .

clean:
	rm -f mcp-dnd-server

setup:
	go mod download
	npm install -g @modelcontextprotocol/inspector

run-inspector:
	mcp-inspector go run .

.PHONY: build run clean setup run-inspector
