package tools

import (
	"context"
	"fmt"
	"strings"

	"github.com/ThorstenHans/akamai-functions-mcp/internal/spin"
	"github.com/mark3labs/mcp-go/mcp"
)

type ByAppArgs struct {
	App     App     `json:"app" jsonschema:"The Akamai Functions App (ID or Name). If unknown, check the 'local://app-context' resource first."`
	Account Account `json:"account,omitempty" jsonschema:"Optionally specify the desired Akamai Functions account by either specifying the account name or its identifier"`
}

func (a ByAppArgs) Validate() error {
	if len(a.App.Name) == 0 && len(a.App.Id) == 0 {
		return fmt.Errorf("application Name or ID is required. Hint: Check 'local://app-context' if you are in a project folder")
	}
	if len(a.App.Name) > 0 && len(a.App.Id) > 0 {
		return fmt.Errorf("specify either Name or ID, not both")
	}
	return nil
}

// Get Application History tool

type AppDeploymentHistoryResponse struct {
	History []spin.AppHistory `json:"history"`
}

func (a *AkamaiFunctionsTools) GetAppDeploymentHistory(ctx context.Context, request mcp.CallToolRequest, args ByAppArgs) (ToolResponse[AppDeploymentHistoryResponse], error) {
	if err := args.Validate(); err != nil {
		a.logger.Printf("Invalid arguments for GetAppDeploymentHistory: %v\n", err)
		return NewToolErrorResponse[AppDeploymentHistoryResponse](err.Error()), nil
	}
	history, err := a.backend.GetAppHistory(ctx, args.Account.Id, args.App.Id, args.App.Name)
	if err != nil {
		a.logger.Printf("Error running command to get app deployment history: %v\n", err)
		return NewToolErrorResponse[AppDeploymentHistoryResponse](err.Error()), err
	}

	return NewToolSuccessResponse(AppDeploymentHistoryResponse{History: history}), nil
}

// Get Application Logs tool

type GetAppLogArguments struct {
	App      App     `json:"app" jsonschema:"The Akamai Functions App (ID or Name). If unknown, check the 'local://app-context' resource first."`
	Account  Account `json:"account,omitempty" jsonschema:"Optionally specify the desired Akamai Functions account by either specifying the account name or its identifier"`
	MaxLines int     `json:"maxLines,omitempty" jsonschema:"Maximum number of log lines to retrieve,default=10"`
}

func (a GetAppLogArguments) Validate() error {
	if len(a.App.Name) == 0 && len(a.App.Id) == 0 {
		return fmt.Errorf("application Name or ID is required. Hint: Check 'local://app-context' if you are in a project folder")
	}
	if len(a.App.Name) > 0 && len(a.App.Id) > 0 {
		return fmt.Errorf("specify either Name or ID, not both")
	}
	if a.MaxLines < 0 {
		return fmt.Errorf("MaxLines cannot be negative")
	}
	return nil
}

func (a *AkamaiFunctionsTools) GetAppLogs(ctx context.Context, request mcp.CallToolRequest, args GetAppLogArguments) (ToolResponse[[]string], error) {
	if err := args.Validate(); err != nil {
		a.logger.Printf("Invalid arguments for GetAppLogs: %v\n", err)
		return NewToolErrorResponse[[]string](err.Error()), err
	}

	logs, err := a.backend.GetAppLogs(ctx, args.MaxLines, args.Account.Id, args.App.Id, args.App.Name)
	if err != nil {
		a.logger.Printf("Error running command: %v\n", err)
		return NewToolErrorResponse[[]string](err.Error()), err
	}
	return NewToolSuccessResponse(logs), nil

}

// Get Application Status & Get Application URL tools share common logic

type AppStatusResponse struct {
	Id          string   `json:"id"`
	Name        string   `json:"name"`
	Urls        []string `json:"urls"`
	CreatedAt   string   `json:"created_at"`
	Invocations int      `json:"invocations"`
}

func (a *AkamaiFunctionsTools) GetAppStatus(ctx context.Context, request mcp.CallToolRequest, args ByAppArgs) (ToolResponse[spin.AppStatus], error) {
	if err := args.Validate(); err != nil {
		a.logger.Printf("Invalid arguments for GetAppLogs: %v\n", err)
		return NewToolErrorResponse[spin.AppStatus](err.Error()), err
	}
	status, err := a.backend.GetAppStatus(ctx, args.Account.Id, args.App.Id, args.App.Name)
	if err != nil {
		a.logger.Printf("Error getting app status: %v\n", err)
		return NewToolErrorResponse[spin.AppStatus](err.Error()), nil
	}
	return NewToolSuccessResponse(*status), nil
}

const deprecatedAkamaiFunctionsDomain = "aka.fermyon.tech"

func (a *AkamaiFunctionsTools) GetAppUrl(ctx context.Context, request mcp.CallToolRequest, args ByAppArgs) (ToolResponse[string], error) {
	if err := args.Validate(); err != nil {
		a.logger.Printf("Invalid arguments for GetAppUrl: %v\n", err)
		return NewToolErrorResponse[string](err.Error()), err
	}
	status, err := a.backend.GetAppStatus(ctx, args.Account.Id, args.App.Id, args.App.Name)
	if err != nil {
		a.logger.Printf("Error getting app url: %v\n", err)
		return NewToolErrorResponse[string](err.Error()), err
	}

	for _, url := range status.Urls {
		if !strings.HasSuffix(url, deprecatedAkamaiFunctionsDomain) {
			return NewToolSuccessResponse(url), nil
		}
	}
	return NewToolErrorResponse[string](fmt.Sprintf("Could not determine URL for app %s", status.Name)), fmt.Errorf("Could not determine URL for app %s", status.Name)
}

// Deploy application

type DeployAppArgs struct {
	App                   App      `json:"app,omitempty" jsonschema:"description=The target app. If empty, the server checks local context."`
	Account               Account  `json:"account,omitempty" jsonschema:"description=Optionally specify the target Akamai Functions account by ID. If omitted, the user's default account is used."`
	IsFirstTimeDeployment bool     `json:"isFirstTimeDeployment" jsonschema:"description=CRITICAL: Set to true ONLY if you intend to create a brand new application. If you want to deploy to an existing app that isn't linked locally, you must provide the App.Id to link it."`
	Variables             []string `json:"variables,omitempty" jsonschema:"description=List of variables in order. Use 'KEY=VALUE' for inline or '@file.json' for files. E.g. ['@base.json', 'ENV=prod']. The last specified key wins."`
}

func (a DeployAppArgs) Validate() error {

	if len(a.App.Name) > 0 && len(a.App.Id) > 0 {
		return fmt.Errorf("specify either app Name or ID, not both")
	}

	return nil
}

func (a *AkamaiFunctionsTools) DeployApp(ctx context.Context, request mcp.CallToolRequest, args DeployAppArgs) (ToolResponse[[]string], error) {

	logs, err := a.backend.DeployApp(ctx, args.Variables, args.IsFirstTimeDeployment, args.Account.Id, args.App.Id, args.App.Name)
	if err != nil {
		return NewToolErrorResponse[[]string](fmt.Sprintf("Deployment failed: %v\n", err)), nil
	}
	return NewToolSuccessResponse(logs), nil
}

// Link Application
type LinkAppArgs struct {
	App     App     `json:"app" jsonschema:"description=REQUIRED: The existing Akamai application to link to the current workspace. You must provide either the ID or the Name."`
	Account Account `json:"account,omitempty" jsonschema:"description=Optionally specify the target Akamai account by name or ID. Defaults to the current account context."`
}

func (a LinkAppArgs) Validate() error {
	if len(a.App.Name) == 0 && len(a.App.Id) == 0 {
		return fmt.Errorf("an application Name or ID is strictly required to link a workspace")
	}
	if len(a.App.Name) > 0 && len(a.App.Id) > 0 {
		return fmt.Errorf("specify either app Name or ID, not both")
	}
	return nil
}

func (a *AkamaiFunctionsTools) LinkApp(ctx context.Context, request mcp.CallToolRequest, args LinkAppArgs) (ToolResponse[[]string], error) {
	if err := args.Validate(); err != nil {
		a.logger.Printf("Validation failed for LinkApp: %v\n", err)
		// Return conversational error so the LLM knows to ask the user or look up the ID
		return NewToolErrorResponse[[]string](err.Error()), nil
	}
	err := a.backend.LinkApp(ctx, args.Account.Id, args.App.Id, args.App.Name)
	if err != nil {
		a.logger.Printf("CLI Error during link: %v\n", err)
		return NewToolErrorResponse[[]string](fmt.Sprintf("Failed to link app: %v", err)), nil
	}
	return NewToolSuccessResponse([]string{"Successfully linked workspace to app"}), nil
}
