package resources

import (
	"context"
	_ "embed"
	"encoding/json"
	"fmt"
	"log"
	"os"
	"path/filepath"

	"github.com/mark3labs/mcp-go/mcp"
	"github.com/mark3labs/mcp-go/server"
	"github.com/pelletier/go-toml/v2"
)

//go:embed aka.md
var spinAkaCommandReference string

type AkamaiFunctionsResources struct {
	logger *log.Logger
}

func NewAkamaiFunctionsResources(logger *log.Logger) *AkamaiFunctionsResources {
	return &AkamaiFunctionsResources{
		logger: logger,
	}
}

const (
	mimeTypeMarkdown        = "text/markdown"
	resourceIdAkaCommandRef = "akamai-functions://docs/reference/spin-aka"
)

func (a *AkamaiFunctionsResources) RegisterAllWith(s *server.MCPServer) {
	a.registerLocalAppContext(s)
	a.registerAkaCommandReference(s)
}

type AkaConfig struct {
	Application struct {
		ID string `toml:"id" json:"id"`
	} `toml:"application" json:"application"`
}

func getLocalAkaConfig() (*AkaConfig, error) {
	cwd, err := os.Getwd()
	if err != nil {
		return nil, fmt.Errorf("failed to get working directory: %w", err)
	}

	configPath := filepath.Join(cwd, ".spin-aka", "config.toml")

	// Read the file content first
	data, err := os.ReadFile(configPath)
	if err != nil {
		return nil, fmt.Errorf("failed to read config file: %w", err)
	}

	var config AkaConfig
	// Use Unmarshal (standard for v2)
	if err := toml.Unmarshal(data, &config); err != nil {
		return nil, fmt.Errorf("failed to parse TOML: %w", err)
	}

	return &config, nil
}

func (a *AkamaiFunctionsResources) registerLocalAppContext(s *server.MCPServer) {
	isLocalApp := mcp.NewResource("local://app-context", "Current Application Context",
		mcp.WithResourceDescription("Contains the Application ID from the local .spin-aka/config.toml if present."),
		mcp.WithMIMEType("application/json"),
	)
	s.AddResource(isLocalApp, func(ctx context.Context, request mcp.ReadResourceRequest) ([]mcp.ResourceContents, error) {

		config, err := getLocalAkaConfig()
		if err != nil {
			// Return a helpful message so the LLM knows why it can't see the context
			return []mcp.ResourceContents{
				mcp.TextResourceContents{MIMEType: "text/plain", Text: "No .spin-aka/config.toml found in the current directory."},
			}, nil
		}

		// Convert config to JSON for the LLM to read easily
		jsonData, _ := json.MarshalIndent(config, "", "  ")

		return []mcp.ResourceContents{
			mcp.TextResourceContents{MIMEType: "application/json", Text: string(jsonData)},
		}, nil

	})
}

func (a *AkamaiFunctionsResources) registerAkaCommandReference(s *server.MCPServer) {
	ref := mcp.NewResource(resourceIdAkaCommandRef, "The 'spin aka' command reference",
		mcp.WithResourceDescription("Reference documentation for the `spin aka` command"),
		mcp.WithMIMEType(mimeTypeMarkdown),
	)

	s.AddResource(ref, func(ctx context.Context, request mcp.ReadResourceRequest) ([]mcp.ResourceContents, error) {
		return []mcp.ResourceContents{
			mcp.TextResourceContents{MIMEType: mimeTypeMarkdown, Text: spinAkaCommandReference},
		}, nil
	})
}
