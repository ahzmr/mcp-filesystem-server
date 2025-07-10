package main

import (
	"encoding/json"
	"flag"
	"fmt"
	"log"
	"net/http"
	"os"
	"strings"

	"github.com/mark3labs/mcp-filesystem-server/filesystemserver"
	"github.com/mark3labs/mcp-go/server"
)

// parseToolConfig parses the tool configuration string
func parseToolConfig(toolsStr string) *filesystemserver.ToolConfig {
	toolsStr = strings.TrimSpace(toolsStr)
	if toolsStr == "" || toolsStr == "all" {
		return &filesystemserver.ToolConfig{EnableAll: true}
	}

	tools := strings.Split(toolsStr, ",")
	enabledTools := make([]string, 0, len(tools))
	for _, tool := range tools {
		tool = strings.TrimSpace(tool)
		if tool != "" {
			enabledTools = append(enabledTools, tool)
		}
	}

	return &filesystemserver.ToolConfig{
		EnabledTools: enabledTools,
		EnableAll:    false,
	}
}

func main() {
	// Define command line flags
	var (
		transport = flag.String("transport", "stdio", "Transport type: stdio or http")
		port      = flag.Int("port", 8080, "Port to listen on (http transport only)")
		host      = flag.String("host", "localhost", "Host to bind to (http transport only)")
		tools     = flag.String("tools", "all", "Comma-separated list of tools to enable (default: all). Supports wildcards like 'read_*,write_*'")
		help      = flag.Bool("help", false, "Show help message")
	)
	flag.Parse()

	// Show help if requested
	if *help {
		fmt.Fprintf(os.Stderr, "Usage: %s [options] <allowed-directory> [additional-directories...]\n", os.Args[0])
		fmt.Fprintf(os.Stderr, "\nOptions:\n")
		flag.PrintDefaults()
		fmt.Fprintf(os.Stderr, "\nTransport types:\n")
		fmt.Fprintf(os.Stderr, "  stdio: Standard input/output (default, for local MCP clients)\n")
		fmt.Fprintf(os.Stderr, "  http:  Streamable HTTP transport (for remote MCP clients)\n")
		fmt.Fprintf(os.Stderr, "\nTool configuration:\n")
		fmt.Fprintf(os.Stderr, "  all:                    Enable all available tools (default)\n")
		fmt.Fprintf(os.Stderr, "  tool1,tool2:           Enable specific tools\n")
		fmt.Fprintf(os.Stderr, "  read_*,write_*:        Enable tools matching wildcards\n")
		fmt.Fprintf(os.Stderr, "\nExamples:\n")
		fmt.Fprintf(os.Stderr, "  %s /home/user/documents                                    # stdio transport, all tools\n", os.Args[0])
		fmt.Fprintf(os.Stderr, "  %s --transport http --port 8080 /home/user/documents      # http transport, all tools\n", os.Args[0])
		fmt.Fprintf(os.Stderr, "  %s --tools read_file,write_file /home/user/documents       # only read and write tools\n", os.Args[0])
		fmt.Fprintf(os.Stderr, "  %s --tools 'read_*,list_*' /home/user/documents            # tools matching wildcards\n", os.Args[0])
		os.Exit(0)
	}

	// Parse remaining arguments as allowed directories
	allowedDirs := flag.Args()
	if len(allowedDirs) < 1 {
		fmt.Fprintf(os.Stderr, "Error: At least one allowed directory must be specified\n")
		fmt.Fprintf(os.Stderr, "Usage: %s [options] <allowed-directory> [additional-directories...]\n", os.Args[0])
		fmt.Fprintf(os.Stderr, "Use --help for more information\n")
		os.Exit(1)
	}

	// Parse tool configuration
	toolConfig := parseToolConfig(*tools)

	// Create the filesystem server
	fss, err := filesystemserver.NewFilesystemServer(allowedDirs, toolConfig)
	if err != nil {
		log.Fatalf("Failed to create server: %v", err)
	}

	// Start server based on transport type
	switch *transport {
	case "stdio":
		// Serve requests via stdio (original behavior)
		log.Printf("Starting MCP Filesystem Server with stdio transport")
		log.Printf("Allowed directories: %v", allowedDirs)
		if err := server.ServeStdio(fss); err != nil {
			log.Fatalf("Server error: %v", err)
		}

	case "http":
		// Serve requests via HTTP (streamable-http transport)
		if err := serveHTTP(fss, *host, *port); err != nil {
			log.Fatalf("HTTP server error: %v", err)
		}

	default:
		log.Fatalf("Unknown transport type: %s. Use 'stdio' or 'http'", *transport)
	}
}

func serveHTTP(mcpServer *server.MCPServer, host string, port int) error {
	// Create StreamableHTTPServer
	streamableServer := server.NewStreamableHTTPServer(mcpServer)

	// Create HTTP server with CORS support
	mux := http.NewServeMux()

	// Add CORS middleware
	corsHandler := func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			w.Header().Set("Access-Control-Allow-Origin", "*")
			w.Header().Set("Access-Control-Allow-Methods", "GET, POST, PUT, DELETE, OPTIONS")
			w.Header().Set("Access-Control-Allow-Headers", "Content-Type, Authorization, Mcp-Session-Id")
			w.Header().Set("Access-Control-Expose-Headers", "Mcp-Session-Id")

			if r.Method == "OPTIONS" {
				w.WriteHeader(http.StatusOK)
				return
			}

			next.ServeHTTP(w, r)
		})
	}

	// Register MCP server at /mcp endpoint using StreamableHTTPServer
	mux.Handle("/mcp", corsHandler(streamableServer))

	// Health check endpoint
	mux.HandleFunc("/health", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		w.Write([]byte(`{"status":"ok","transport":"streamable-http"}`))
	})

	// Root endpoint with server information
	mux.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		info := map[string]interface{}{
			"name":        "MCP Filesystem Server",
			"version":     "1.0.0",
			"transport":   "streamable-http",
			"endpoints": map[string]string{
				"mcp":    "/mcp",
				"health": "/health",
			},
			"description": "MCP server providing filesystem operations with streamable-HTTP transport",
		}
		json.NewEncoder(w).Encode(info)
	})

	addr := fmt.Sprintf("%s:%d", host, port)
	fmt.Printf("Starting MCP filesystem server with HTTP transport on %s\n", addr)
	fmt.Printf("MCP endpoint: http://%s/mcp\n", addr)
	fmt.Printf("Health check: http://%s/health\n", addr)

	return http.ListenAndServe(addr, mux)
}
