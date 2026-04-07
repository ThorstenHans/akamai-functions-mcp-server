package tools

import (
	"context"
	"encoding/json"
	"strings"

	"github.com/ThorstenHans/akamai-functions-mcp/internal/spin"
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
	command := []string{"aka", "info", "--format", "json"}
	a.logger.Printf("Will find all your Akamai Functions accounts using the following spin command: %v\n", command)
	out, err := spin.RunCommand(command...)
	if err != nil {
		a.logger.Printf("Error running command %v: %v\n", command, err)
		return NewToolErrorResponse[SearchResults]("Error running spin command"), err
	}
	var accountInfo spinAccountInfoResponse

	err = json.Unmarshal(out, &accountInfo)
	if err != nil {
		a.logger.Printf("Error unmarshalling output of command %v: %v\nOutput was: %s\n", command, err, string(out))
		return NewToolErrorResponse[SearchResults]("Error unmarshalling spin command output"), err
	}
	a.logger.Printf("Found %d accounts\n", len(accountInfo.AuthInfo.Accounts))

	apps := make([]SearchResult, 0)
	for _, account := range accountInfo.AuthInfo.Accounts {
		a.logger.Printf("Will retrieve all apps for particular account using the following spin command: %v\n", []string{"aka", "apps", "list", "--format", "json", "--account-id", account.Id})
		out, err := spin.RunCommand("aka", "apps", "list", "--format", "json", "--account-id", account.Id)
		if err != nil {
			a.logger.Printf("Error running command to get apps for account %s: %v\nOutput was: %s\n", account.Name, err, string(out))
			return NewToolErrorResponse[SearchResults]("Error running spin command to get apps for account " + account.Name), err
		}
		var appsPerAccount []ListAppsItem
		err = json.Unmarshal(out, &appsPerAccount)
		if err != nil {
			a.logger.Printf("Error unmarshalling output of command to get apps for account %s: %v\nOutput was: %s\n", account.Name, err, string(out))
			return NewToolErrorResponse[SearchResults]("Error unmarshalling spin command output for account " + account.Name), err
		}
		a.logger.Printf("Found %d apps for account %s\n", len(appsPerAccount), account.Name)
		for _, app := range appsPerAccount {
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
