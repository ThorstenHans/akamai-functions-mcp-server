package tools

import (
	"context"

	"github.com/ThorstenHans/akamai-functions-mcp/internal/spin"
	"github.com/mark3labs/mcp-go/mcp"
)

type ListAccountResponse struct {
	Accounts []spin.Account `json:"accounts" jsonschema:"List of Akamai Functions accounts you have access to"`
}

type ListAccountsArgs struct{}

func (a *AkamaiFunctionsTools) ListAccounts(ctx context.Context, request mcp.CallToolRequest, args ListAccountsArgs) (ToolResponse[ListAccountResponse], error) {

	accounts, err := a.backend.ListAccounts(ctx)
	if err != nil {
		a.logger.Panic(err)
		return NewToolErrorResponse[ListAccountResponse](err.Error()), err
	}
	return NewToolSuccessResponse(ListAccountResponse{
		Accounts: accounts,
	}), nil
}
