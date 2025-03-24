package main

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"time"

	"github.com/mark3labs/mcp-go/client"
	"github.com/mark3labs/mcp-go/mcp"
)

func main() {
	// apiKey := os.Getenv("ELEVENLABS_API_KEY")
	c, err := client.NewStdioMCPClient(
		"go",
		[]string{
			// fmt.Sprintf("ELEVENLABS_API_KEY=%s", apiKey),
		},
		"run",
		"main.go",
	)
	if err != nil {
		log.Fatalf("Failed to create client: %v", err)
	}
	defer c.Close()

	// Create context with timeout
	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	// Initialize the client
	fmt.Println("Initializing client...")
	initRequest := mcp.InitializeRequest{}
	initRequest.Params.ProtocolVersion = mcp.LATEST_PROTOCOL_VERSION
	initRequest.Params.ClientInfo = mcp.Implementation{
		Name:    "example-client",
		Version: "1.0.0",
	}

	initResult, err := c.Initialize(ctx, initRequest)
	if err != nil {
		log.Fatalf("Failed to initialize: %v", err)
	}
	fmt.Printf(
		"Initialized with server: %s %s\n\n",
		initResult.ServerInfo.Name,
		initResult.ServerInfo.Version,
	)

	// List Tools
	fmt.Println("Listing available tools...")
	toolsRequest := mcp.ListToolsRequest{}
	tools, err := c.ListTools(ctx, toolsRequest)
	if err != nil {
		log.Fatalf("Failed to list tools: %v", err)
	}
	for _, tool := range tools.Tools {
		fmt.Printf("- %s: %s\n", tool.Name, tool.Description)
	}
	fmt.Println()

	// Say
	fmt.Println("Say...")
	sayRequest := mcp.CallToolRequest{}
	sayRequest.Params.Name = "say"
	sayRequest.Params.Arguments = map[string]any{
		"text": "Hello, world!",
		// "voice": "Daniel",
	}

	result, err := c.CallTool(ctx, sayRequest)
	if err != nil {
		log.Fatalf("Failed to run say: %v", err)
	}
	printToolResult(result)
	fmt.Println()

	// ElevenLabs
	fmt.Println("ElevenLabs...")
	elevenLabsRequest := mcp.CallToolRequest{}
	elevenLabsRequest.Params.Name = "elevenlabs"
	elevenLabsRequest.Params.Arguments = map[string]any{
		"text":  "Hello, world!",
		"voice": "V9fdGZs6AiHI4uyiAiza",
	}

	result, err = c.CallTool(ctx, elevenLabsRequest)
	if err != nil {
		log.Fatalf("Failed to run elevenlabs: %v", err)
	}
	printToolResult(result)
	fmt.Println()
}

// Helper function to print tool results
func printToolResult(result *mcp.CallToolResult) {
	for _, content := range result.Content {
		if textContent, ok := content.(mcp.TextContent); ok {
			fmt.Println(textContent.Text)
		} else {
			jsonBytes, _ := json.MarshalIndent(content, "", "  ")
			fmt.Println(string(jsonBytes))
		}
	}
}
