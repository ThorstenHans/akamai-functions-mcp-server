package tools

import (
	"context"

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
	Apps []spin.App `json:"apps"`
}

func (a *AkamaiFunctionsTools) ListApps(ctx context.Context, request mcp.CallToolRequest, args MaybeByAccountArgs) (ToolResponse[*ListAppsResponse], error) {

	a.logger.Printf("Will find all your Akamai Functions account")
	apps, err := a.backend.ListApps(ctx, args.Account.Id)
	if err != nil {
		a.logger.Printf("Error listing apps for account %s: %v\n", args.Account.Id, err)
		return NewToolErrorResponse[*ListAppsResponse](err.Error()), err
	}
	a.logger.Printf("Found %d apps for account %s\n", len(apps), args.Account.Id)
	return NewToolSuccessResponse(&ListAppsResponse{
		Apps: apps,
	}), nil
}
