package main

import (
	"flag"
	"io"
	"log"
	"os"

	"github.com/ThorstenHans/akamai-functions-mcp/internal/resources"
	"github.com/ThorstenHans/akamai-functions-mcp/internal/spin"
	"github.com/ThorstenHans/akamai-functions-mcp/internal/tools"

	"github.com/mark3labs/mcp-go/server"
)

func main() {
	debug := flag.Bool("debug", false, "enable debug logging")
	flag.Parse()

	var logger *log.Logger
	if *debug {
		logFile, err := os.OpenFile("akamai-functions-mcp.log", os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0o666)
		if err != nil {
			log.Fatal(err)
		}
		defer logFile.Close()

		logger = log.New(logFile, "[MCP] ", log.LstdFlags|log.Lshortfile)
	} else {
		logger = log.New(io.Discard, "", 0)
	}

	var backend = spin.NewSpinBackend(logger)

	afTools := tools.NewAkamaiFunctionsTools(backend, logger)
	afResources := resources.NewAkamaiFunctionsResources(logger)

	serverOptions := []server.ServerOption{
		server.WithToolCapabilities(false),
		server.WithPromptCapabilities(false),
		server.WithRecovery(),
	}

	if *debug {
		serverOptions = append(serverOptions, server.WithLogging())
	}

	mcpServer := server.NewMCPServer(
		"Akamai Functions MCP Server",
		"0.1.0",
		serverOptions...,
	)

	afTools.RegisterAllWith(mcpServer)
	afResources.RegisterAllWith(mcpServer)
	if err := server.ServeStdio(mcpServer); err != nil {
		logger.Printf("Stdio server failed: %v\n", err)
	}

}
