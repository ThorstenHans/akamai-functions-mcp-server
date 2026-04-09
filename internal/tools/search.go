package tools

import (
	"context"
	"strings"

	"github.com/mark3labs/mcp-go/mcp"
)

type SearchAppArguments struct {
	Query string `json:"query" jsonschema:"A query to search apps for"`
}

type SearchResults struct {
	Results []SearchResult `json:"results"`
}

type SearchResult struct {
	AppName     string `json:"appName"`
	AppId       string `json:"appId"`
	AccountId   string `json:"accountId"`
	AccountName string `json:"accountName"`
}

func (a *AkamaiFunctionsTools) SearchAppByName(ctx context.Context, request mcp.CallToolRequest, args SearchAppArguments) (ToolResponse[SearchResults], error) {
	if len(args.Query) == 0 {
		a.logger.Println("Search tool called without query, will terminate")
		return NewToolErrorResponse[SearchResults]("You must provide a query"), nil
	}
	args.Query = strings.ToLower(args.Query)
	a.logger.Printf("Will find all your Akamai Functions accounts")
	accounts, err := a.backend.ListAccounts(ctx)
	if err != nil {
		a.logger.Printf("Error listing accounts: %v\n", err)
		return NewToolErrorResponse[SearchResults]("Error listing accounts"), err
	}
	a.logger.Printf("Found %d accounts\n", len(accounts))

	apps := make([]SearchResult, 0)
	for _, account := range accounts {
		a.logger.Printf("Will retrieve all apps for particular Akamai Functions account")
		existingApps, err := a.backend.ListApps(ctx, account.Id)
		if err != nil {
			a.logger.Printf("Error listing apps for Akamai Functions account: %v\n", err)
			return NewToolErrorResponse[SearchResults]("Error retrieving apps from Akamai Functions account"), err
		}
		a.logger.Printf("Found %d apps with authorized Akamai Functions account\n", len(existingApps))
		for _, app := range existingApps {
			if strings.Contains(strings.ToLower(app.Name), args.Query) {
				apps = append(apps, SearchResult{
					AppId:       app.Id,
					AppName:     app.Name,
					AccountId:   account.Id,
					AccountName: account.Name,
				})
			}
		}
	}
	a.logger.Printf("Search for query '%s' resulted in %d results\n", args.Query, len(apps))
	return NewToolSuccessResponse[SearchResults](SearchResults{
		Results: apps,
	}), nil
}
