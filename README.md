# D&D 5e MCP Server (Go)

> **Note:** This project is a work in progress and currently provides limited functionality. Only basic spell and monster queries are supported at this time. Additional features and categories may be added in the future.

This project is a Go implementation of a Model Context Protocol (MCP) server for D&D 5e API integration, using [github.com/mark3labs/mcp-go](https://github.com/mark3labs/mcp-go). It exposes D&D 5e data (spells, monsters, and more) as MCP tools for use in LLM-powered applications.

## Features

- List D&D 5e API categories (spells, monsters, etc.)
- List items in a category (e.g., all spells, all monsters, with filtering)
- Get details for a specific item (e.g., spell by name, monster by index)
- Implements MCP server using Go's standard library for HTTP requests
- Includes unit tests for core logic

## Getting Started

1. Ensure you have Go installed.
2. Run `make setup` to install Go dependencies and the MCP Inspector (for development).
3. Use the provided Makefile commands to build, run, or clean the project.

## Makefile Commands

- `make build` — Build the server binary (`mcp-dnd-server`).
- `make run` — Run the server locally.
- `make clean` — Remove the built server binary.
- `make setup` — Install Go dependencies and the MCP Inspector tool.
- `make run-inspector` — Launch the MCP Inspector with the server for local testing.

## Usage

To build and run the server:

```sh
make build
make run
```

Or run directly (without building):

```sh
make run
```

To stop and clean up:

```sh
make clean
```

To run the MCP Inspector (for local development):

```sh
make run-inspector
```

## MCP Tools

### Spells Tool

The Spells tool allows you to:

- List all D&D 5e spells
- Filter spells by level and/or school of magic
- Retrieve detailed information for a specific spell by name

**Filtering options:**

- `level` (integer): Only return spells of a specific level (e.g., 1 for Magic Missile, 3 for Fireball)
- `school` (string): Only return spells from a specific school of magic (e.g., "evocation", "illusion")

**Input fields:**

- `name` (string, optional): The name of the spell to retrieve details for. If provided, returns only that spell's details.
- `filter` (object, optional): Filtering options for listing spells. If `name` is not provided, returns a list of spells matching the filter.

#### Example: Get a spell by name

```json
{
  "name": "Fireball"
}
```

#### Example: List all 1st-level evocation spells

```json
{
  "filter": {
    "level": 1,
    "school": "evocation"
  }
}
```

#### Spells Tool Example Output

```json
{
  "results": [
    { "name": "Magic Missile", "level": 1, ... }
  ],
  "count": 1
}
```

### Monsters Tool

The Monsters tool allows you to:

- List all D&D 5e monsters
- Filter monsters by challenge rating (CR)
- Retrieve detailed information for a specific monster by index (e.g., "goblin", "orc")

**Filtering options:**

- `challenge_rating` (array of numbers): Only return monsters matching one or more challenge ratings (e.g., `[1, 2.5]`)

**Input fields:**

- `name` (string, optional): The index of the monster to retrieve details for. If provided, returns only that monster's details.
- `filter` (object, optional): Filtering options for listing monsters. If `name` is not provided, returns a list of monsters matching the filter.

#### Example: Get a monster by index

```json
{
  "name": "goblin"
}
```

#### Example: List all monsters with CR 1 or 2.5

```json
{
  "filter": {
    "challenge_rating": [1, 2.5]
  }
}
```

#### Monsters Tool Example Output

```json
{
  "results": [
    { "name": "Goblin", "challenge_rating": 1, ... },
    { "name": "Orc", "challenge_rating": 2.5, ... }
  ],
  "count": 2
}
```

## Development & Testing

Run all unit tests:

```sh
  go test ./...
```

Test data is in `testdata/`.
Use the Makefile for common tasks.

## License

MIT
