package tools

import (
	"context"
	"encoding/json"
	"fmt"
	"strconv"
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

func (a ByAppArgs) getCommandArgs() []string {
	result := []string{}
	if len(a.Account.Id) > 0 {
		result = append(result, "--account-id", a.Account.Id)
	}
	if len(a.App.Id) > 0 {
		result = append(result, "--app-id", a.App.Id)
	} else if len(a.App.Name) > 0 {
		result = append(result, "--app-name", a.App.Name)
	}
	return result
}

type ByAppOrUrlArgs struct {
	App     App     `json:"app,omitempty" jsonschema:"The Akamai Functions App (ID or Name). If unknown, check the 'local://app-context' resource first."`
	Account Account `json:"account,omitempty" jsonschema:"Optionally specify the desired Akamai Functions account by either specifying the account name or its identifier"`
}

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

func (a GetAppLogArguments) getCommandArgs() []string {
	result := []string{}
	if len(a.Account.Id) > 0 {
		result = append(result, "--account-id", a.Account.Id)
	}
	if len(a.App.Id) > 0 {
		result = append(result, "--app-id", a.App.Id)
	} else if len(a.App.Name) > 0 {
		result = append(result, "--app-name", a.App.Name)
	}
	if a.MaxLines > 0 && a.MaxLines != 10 {
		result = append(result, "--max-lines", strconv.Itoa(a.MaxLines))
	}
	return result
}

type AppHistory struct {
	EventType string `json:"event_type"`
	Version   int    `json:"version"`
	Timestamp string `json:"timestamp"`
}

type AppDeploymentHistoryResponse struct {
	History []AppHistory `json:"history"`
}

func (a *AkamaiFunctionsTools) GetAppDeploymentHistory(ctx context.Context, request mcp.CallToolRequest, args ByAppArgs) (ToolResponse[AppDeploymentHistoryResponse], error) {
	if err := args.Validate(); err != nil {
		a.logger.Printf("Invalid arguments for GetAppDeploymentHistory: %v\n", err)
		return NewToolErrorResponse[AppDeploymentHistoryResponse](err.Error()), nil
	}
	command := []string{"aka", "app", "history", "--format", "json"}
	extraArgs := args.getCommandArgs()
	if len(extraArgs) > 0 {
		command = append(command, extraArgs...)
	}
	a.logger.Printf("Will retrieve deployment history using the following spin command: %v\n", command)
	out, err := spin.RunCommand(command...)
	if err != nil {
		a.logger.Printf("Error running command to get app deployment history: %v\nOutput was: %s\n", err, string(out))
		return NewToolErrorResponse[AppDeploymentHistoryResponse](err.Error()), err
	}
	var history []AppHistory
	err = json.Unmarshal(out, &history)
	if err != nil {
		a.logger.Printf("Error unmarshalling output of command to get app deployment history: %v\nOutput was: %s\n", err, string(out))
		return NewToolErrorResponse[AppDeploymentHistoryResponse](err.Error()), err
	}

	return NewToolSuccessResponse(AppDeploymentHistoryResponse{History: history}), nil
}

func (a *AkamaiFunctionsTools) GetAppLogs(ctx context.Context, request mcp.CallToolRequest, args GetAppLogArguments) (ToolResponse[[]string], error) {
	if err := args.Validate(); err != nil {
		a.logger.Printf("Invalid arguments for GetAppLogs: %v\n", err)
		return NewToolErrorResponse[[]string](err.Error()), err
	}
	command := []string{"aka", "logs"}
	extraArgs := args.getCommandArgs()
	if len(extraArgs) > 0 {
		command = append(command, extraArgs...)
	}
	a.logger.Printf("Will retrieve logs using the following spin command: %v\n", command)
	out, err := spin.RunCommand(command...)
	if err != nil {
		a.logger.Printf("Error running command: %v\nOutput: %s\n", err, string(out))
		return NewToolErrorResponse[[]string](err.Error()), err
	}
	logs := strings.Split(string(out), "\n")
	logs = append([]string{}, logs...)
	for i := len(logs) - 1; i >= 0; i-- {
		if strings.TrimSpace(logs[i]) == "" {
			logs = append(logs[:i], logs[i+1:]...)
		}
	}
	return NewToolSuccessResponse(logs), nil

}

type AppStatusResponse struct {
	Id          string   `json:"id"`
	Name        string   `json:"name"`
	Urls        []string `json:"urls"`
	CreatedAt   string   `json:"created_at"`
	Invocations int      `json:"invocations"`
}

func (a *AkamaiFunctionsTools) GetAppStatus(ctx context.Context, request mcp.CallToolRequest, args ByAppArgs) (ToolResponse[AppStatusResponse], error) {
	status, err := a.getAppStatus(args)
	if err != nil {
		a.logger.Printf("Error getting app status: %v\n", err)
		return NewToolErrorResponse[AppStatusResponse](err.Error()), err
	}
	return NewToolSuccessResponse(*status), nil
}

func (a *AkamaiFunctionsTools) GetAppUrl(ctx context.Context, request mcp.CallToolRequest, args ByAppArgs) (ToolResponse[string], error) {
	status, err := a.getAppStatus(args)
	if err != nil {
		a.logger.Printf("Error getting app status: %v\n", err)
		return NewToolErrorResponse[string](err.Error()), err
	}
	for _, url := range status.Urls {
		if !strings.HasSuffix(url, deprecatedAkamaiFunctionsDomain) {
			return NewToolSuccessResponse(url), nil
		}
	}
	return NewToolErrorResponse[string](fmt.Sprintf("Could not determine URL for app %s", status.Name)), fmt.Errorf("Could not determine URL for app %s", status.Name)
}

const deprecatedAkamaiFunctionsDomain = "aka.fermyon.tech"

func (a *AkamaiFunctionsTools) getAppStatus(args ByAppArgs) (*AppStatusResponse, error) {
	if err := args.Validate(); err != nil {

		return nil, err
	}
	command := []string{"aka", "app", "status", "--format", "json"}
	extraArgs := args.getCommandArgs()
	if len(extraArgs) > 0 {
		command = append(command, extraArgs...)
	}
	a.logger.Printf("Will retrieve status using the following spin command: %v\n", command)
	out, err := spin.RunCommand(command...)
	if err != nil {
		return nil, err
	}
	var status AppStatusResponse
	err = json.Unmarshal(out, &status)
	if err != nil {
		return nil, err
	}

	n := 0
	for _, url := range status.Urls {
		if !strings.HasSuffix(url, deprecatedAkamaiFunctionsDomain) {
			status.Urls[n] = url
			n++
		}
	}

	status.Urls = status.Urls[:n]
	return &status, nil
}
