package main

import (
	"context"
	"encoding/json"
	"flag"
	"fmt"
	"log"
	"net/http"

	"github.com/labstack/echo/v4"
	"github.com/mark3labs/mcp-go/mcp"
	"github.com/mark3labs/mcp-go/server"
)

func main() {
	var addr string
	flag.StringVar(&addr, "addr", ":8080", "address to listen on")
	flag.Parse()

	mcpServer := server.NewMCPServer("weather-mcp", "1.0.0")

	// Add a weather forecast tool
	weatherTool := mcp.NewTool("weather",
		mcp.WithDescription("Get weather forecast"),
		mcp.WithString("Latitude",
			mcp.Required(),
			mcp.Description("Latitude of the location"),
		),
		mcp.WithString("Longitude",
			mcp.Required(),
			mcp.Description("Longitude of the location"),
		),
		mcp.WithString("Time",
			mcp.Required(),
			mcp.Description("Time of the forecast"),
		),
	)

	// Add the weather handler (uses open-meteo API)
	mcpServer.AddTool(weatherTool, weatherHandler)

	sseServer := server.NewSSEServer(
		mcpServer,
		server.WithStaticBasePath("/api/"),
		server.WithBaseURL(fmt.Sprintf("http://localhost%s", addr)),
		server.WithUseFullURLForMessageEndpoint(true),
	)

	e := echo.New()
	e.GET("/api/sse", echo.WrapHandler(sseServer.SSEHandler()))
	e.POST("/api/message", echo.WrapHandler(sseServer.MessageHandler()))

	log.Printf("Static SSE server listening on %s", addr)
	if err := e.Start(addr); err != nil {
		log.Fatalf("Server error: %v", err)
	}
}

func weatherHandler(ctx context.Context, request mcp.CallToolRequest) (*mcp.CallToolResult, error) {
	latitude := request.Params.Arguments["Latitude"].(string)
	longitude := request.Params.Arguments["Longitude"].(string)
	time := request.Params.Arguments["Time"].(string)
	url := fmt.Sprintf("https://api.open-meteo.com/v1/forecast?latitude=%s&longitude=%s&hourly=temperature_2m,relative_humidity_2m,wind_speed_10m&timezone=auto", latitude, longitude)
	resp, err := http.Get(url)
	if err != nil {
		return mcp.NewToolResultError("failed to fetch weather data"), nil
	}
	defer resp.Body.Close()
	if resp.StatusCode != http.StatusOK {
		return mcp.NewToolResultError("failed to fetch weather data"), nil
	}
	var data struct {
		Hourly struct {
			Temperature     []float64 `json:"temperature_2m"`
			RelativeHumidity []float64 `json:"relative_humidity_2m"`
			WindSpeed       []float64 `json:"wind_speed_10m"`
		} `json:"hourly"`
	}
	if err := json.NewDecoder(resp.Body).Decode(&data); err != nil {
		return mcp.NewToolResultError("failed to parse weather data"), nil
	}
	if len(data.Hourly.Temperature) == 0 {
		return mcp.NewToolResultError("no temperature data available"), nil
	}
	temperature := data.Hourly.Temperature[0]
	humidity := data.Hourly.RelativeHumidity[0]
	windSpeed := data.Hourly.WindSpeed[0]
	return mcp.NewToolResultText(fmt.Sprintf("The temperature at %s, %s on %s is %.2f°C with %.2f%% humidity and a wind speed of %.2f m/s", latitude, longitude, time, temperature, humidity, windSpeed)), nil
}
