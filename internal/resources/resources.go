package resources

import (
	"context"
	_ "embed"
	"encoding/json"
	"fmt"
	"log"
	"net/url"
	"os"
	"path/filepath"
	"strings"

	"github.com/ThorstenHans/akamai-functions-mcp/internal/spin"
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

type CmdHelp struct {
	Command string `json:"command" jsonschema:"The 'spin aka' sub-command you want help with. For example: 'apps list' or 'apps deploy'"`
}

type AkaConfig struct {
	Application struct {
		ID string `toml:"id" json:"id"`
	} `toml:"application" json:"application"`
}

func (a *AkamaiFunctionsResources) RegisterAllWith(s *server.MCPServer) {
	a.registerLocalAppContext(s)
	a.registerAkaCommandReference(s)
	a.registerAkaCommandHelp(s)
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

func (a *AkamaiFunctionsResources) registerAkaCommandHelp(s *server.MCPServer) {
	cmdHelp := mcp.NewResourceTemplate("aka-help://{command}", "Help for 'spin aka' sub-commands",
		mcp.WithTemplateDescription("Help documentation for the `spin aka` sub-commands"),
		mcp.WithTemplateMIMEType("text/plain"),
	)

	s.AddResourceTemplate(cmdHelp, func(ctx context.Context, request mcp.ReadResourceRequest) ([]mcp.ResourceContents, error) {
		cmd := extractCommandFromUri(request.Params.URI)
		cmd, err := url.QueryUnescape(cmd)
		if err != nil {
			a.logger.Printf("Error unescaping command from URI %s: %v\n", request.Params.URI, err)
			return []mcp.ResourceContents{
				mcp.TextResourceContents{MIMEType: "text/plain", Text: "Error parsing command from URI"},
			}, err
		}
		cmdParts := strings.Split(cmd, " ")
		arguments := append([]string{"aka"}, cmdParts...)
		arguments = append(arguments, "--help")
		a.logger.Printf("Running command for help resource: %v\n", arguments)
		out, err := spin.RunCommand(arguments...)
		if err != nil {
			a.logger.Printf("Error running command %s: %v\nOutput was: %s\n", arguments, err, string(out))
			return []mcp.ResourceContents{
				mcp.TextResourceContents{MIMEType: "text/plain", Text: "Error retrieving help for command " + cmd},
			}, err
		}
		return []mcp.ResourceContents{
			mcp.TextResourceContents{MIMEType: "text/plain", Text: string(out)},
		}, nil
	})
}

func extractCommandFromUri(uri string) string {
	// Extract ID from "users://123" format
	parts := strings.Split(uri, "://")
	if len(parts) == 2 {
		return parts[1]
	}
	return ""
}
