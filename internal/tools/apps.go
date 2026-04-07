package tools

import (
	"context"
	"encoding/json"

	"github.com/ThorstenHans/akamai-functions-mcp/internal/spin"
	"github.com/mark3labs/mcp-go/mcp"
)

type MaybeByAccountArgs struct {
	Account `json:"account,omitempty" jsonschema:"Specify the desired Akamai Functions account by either providing the account name or its identifier"`
}

type App struct {
	Name string `json:"name,omitempty" jsonschema:"Name of the application running on Akamai Functions"`
	Id   string `json:"id,omitempty" jsonschema:"Identifier of the application running on Akamai Functions"`
}

type Account struct {
	Id string `json:"id,omitempty" jsonschema:"Akamai Functions Account Id"`
}

type ListAppsItem struct {
	Id   string `json:"id" jsonschema:"Identifier of the application deployed to Akamai Functions"`
	Name string `json:"name" jsonschema:"Name of the application deployed to Akamai Functions"`
}

type ListAppsResponse struct {
	Apps []ListAppsItem `json:"apps"`
}

func (a *AkamaiFunctionsTools) ListApps(ctx context.Context, request mcp.CallToolRequest, args MaybeByAccountArgs) (*ListAppsResponse, error) {
	command := []string{"aka", "apps", "list", "--format", "json"}
	if len(args.Account.Id) > 0 {
		command = append(command, "--account-id", args.Account.Id)
	}
	a.logger.Printf("Will retrieve all apps for the specified account using the following spin command: %v\n", command)
	out, err := spin.RunCommand(command...)
	if err != nil {
		a.logger.Printf("Error running command %v: %v\nOutput was: %s\n", command, err, string(out))
		return nil, err
	}
	var apps []ListAppsItem
	err = json.Unmarshal(out, &apps)
	if err != nil {
		a.logger.Printf("Error unmarshalling output of command %v: %v\nOutput was: %s\n", command, err, string(out))
		return nil, err
	}
	a.logger.Printf("Found %d apps for account %s\n", len(apps), args.Account.Id)
	return &ListAppsResponse{
		Apps: apps,
	}, nil
}
