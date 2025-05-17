# Go Weather MCP Server

This project implements a simple weather forecast server using the `mcp-go` library.
It exposes an MCP tool that fetches weather data from the Open-Meteo API.

The server provides:
- An MCP tool named "weather" to get forecasts based on latitude, longitude, and time.
- HTTP endpoints for SSE (`/api/sse`) and messaging (`/api/message`).

## Running the server

To run the server, execute:

```sh
go run main.go
```

By default, it listens on port `:8080`. You can change the address using the `-addr` flag:

```sh
go run main.go -addr :your_port
```

## Using the server on a LLM client

To use the server with a LLM client, you first need to configure the client to connect to the server. An example config is:
```json
{
    "servers": {
        "calculator": {
            "type": "sse",
            "url": "http://localhost:8080/api/sse",
        }
    }
}
``` 
You can use it in clients like Claude Desktop or VSCode Agent Chat. See the followign video for a demo: 
![Demo](demo.gif)

## Why this project?
Just wanted to play around with mcp servers and found `mcp-go` to be a great library for building them. It abstracts away all the jsonrpc protocol stuff and lets you focus on the actual functionality.

Why echo? Because I like it. Very overkill for this project, but I wanted to test it with mcp because I will be building bigger projects with it.