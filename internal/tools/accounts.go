package tools

import (
	"context"
	"encoding/json"

	"github.com/ThorstenHans/akamai-functions-mcp/internal/spin"
	"github.com/mark3labs/mcp-go/mcp"
)

type spinAccountInfoResponse struct {
	AuthInfo spinAuthInfo `json:"auth_info"`
}

type spinAuthInfo struct {
	Accounts []AccountInfo `json:"accounts"`
}

type ListAccountResponse struct {
	Accounts []AccountInfo `json:"accounts" jsonschema:"List of Akamai Functions accounts you have access to"`
}

type AccountInfo struct {
	Id   string `json:"id" jsonschema:"Unique Akamai Functions account identifier"`
	Name string `json:"name" jsonschema:"The name of the Akamai Functions account"`
}

type ListAccountsArgs struct{}

func (a *AkamaiFunctionsTools) ListAccounts(ctx context.Context, request mcp.CallToolRequest, args ListAccountsArgs) (ToolResponse[ListAccountResponse], error) {
	command := []string{"aka", "info", "--format", "json"}
	a.logger.Printf("Will retrieve all your Akamai Functions accounts using the following spin command: %v", command)
	out, err := spin.RunCommand(command...)
	if err != nil {
		return NewToolErrorResponse[ListAccountResponse](err.Error()), err
	}
	var accountInfo spinAccountInfoResponse

	err = json.Unmarshal(out, &accountInfo)
	if err != nil {
		return NewToolErrorResponse[ListAccountResponse](err.Error()), err
	}
	return NewToolSuccessResponse(ListAccountResponse{
		Accounts: accountInfo.AuthInfo.Accounts,
	}), nil
}
