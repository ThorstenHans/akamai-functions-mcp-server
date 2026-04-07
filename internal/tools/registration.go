package tools

import (
	"github.com/mark3labs/mcp-go/mcp"
	"github.com/mark3labs/mcp-go/server"
)

func (a *AkamaiFunctionsTools) RegisterAllWith(s *server.MCPServer) {
	searchAppTool := mcp.NewTool("search_app",
		mcp.WithDescription("Find an app using the specified query in any of my Akamai Functions accounts"),
		mcp.WithReadOnlyHintAnnotation(true),
		mcp.WithDestructiveHintAnnotation(false),
		mcp.WithInputSchema[SearchAppArguments](),
		mcp.WithOutputSchema[ToolResponse[SearchResults]]())

	s.AddTool(searchAppTool, mcp.NewStructuredToolHandler(a.SearchAppByName))

	listAppsTool := mcp.NewTool("list_apps",
		mcp.WithDescription("Get all Spin apps deployed to your Akamai Functions account"),
		mcp.WithReadOnlyHintAnnotation(true),
		mcp.WithDestructiveHintAnnotation(false),
		mcp.WithInputSchema[MaybeByAccountArgs](),
		mcp.WithOutputSchema[ToolResponse[ListAppsResponse]]())
	s.AddTool(listAppsTool, mcp.NewStructuredToolHandler(a.ListApps))

	listAccountsTool := mcp.NewTool("list_accounts",
		mcp.WithDescription("List all Akamai Functions accounts I have access to"),
		mcp.WithReadOnlyHintAnnotation(true),
		mcp.WithDestructiveHintAnnotation(false),
		mcp.WithInputSchema[ListAccountsArgs](),
		mcp.WithOutputSchema[ToolResponse[ListAccountResponse]]())
	s.AddTool(listAccountsTool, mcp.NewStructuredToolHandler(a.ListAccounts))

	getAppStatusTool := mcp.NewTool("get_app_status",
		mcp.WithDescription("Retrieve the status of an Akamai Function. You can omit the ID or Name if you are running this in a project directory containing a .spin-aka/config.toml file; the server will auto-detect the application."),
		mcp.WithReadOnlyHintAnnotation(true),
		mcp.WithDestructiveHintAnnotation(false),
		mcp.WithInputSchema[ByAppArgs](),
		mcp.WithOutputSchema[ToolResponse[AppStatusResponse]]())

	s.AddTool(getAppStatusTool, mcp.NewStructuredToolHandler(a.GetAppStatus))

	getAppUrlTool := mcp.NewTool("get_app_url",
		mcp.WithDescription("Retrieve the public endpoint for an Akamai Function. You can omit the ID or Name if you are running this in a project directory containing a .spin-aka/config.toml file; the server will auto-detect the application."),
		mcp.WithReadOnlyHintAnnotation(true),
		mcp.WithDestructiveHintAnnotation(false),
		mcp.WithInputSchema[ByAppArgs](),
		mcp.WithOutputSchema[ToolResponse[string]]())
	s.AddTool(getAppUrlTool, mcp.NewStructuredToolHandler(a.GetAppUrl))

	getAppLogsTool := mcp.NewTool("get_app_logs",
		mcp.WithDescription("Retrieves logs for an Akamai Function. You can omit the ID or Name if you are running this in a project directory containing a .spin-aka/config.toml file; the server will auto-detect the application."),
		mcp.WithReadOnlyHintAnnotation(true),
		mcp.WithDestructiveHintAnnotation(false),
		mcp.WithInputSchema[GetAppLogArguments](),
		mcp.WithOutputSchema[ToolResponse[[]string]]())
	s.AddTool(getAppLogsTool, mcp.NewStructuredToolHandler(a.GetAppLogs))

	getAppHistoryTool := mcp.NewTool("get_app_history",
		mcp.WithDescription("Retrieve the history of an Akamai Function. You can omit the ID or Name if you are running this in a project directory containing a .spin-aka/config.toml file; the server will auto-detect the application."),
		mcp.WithReadOnlyHintAnnotation(true),
		mcp.WithDestructiveHintAnnotation(false),
		mcp.WithInputSchema[ByAppArgs](),
		mcp.WithOutputSchema[ToolResponse[AppDeploymentHistoryResponse]]())

	s.AddTool(getAppHistoryTool, mcp.NewStructuredToolHandler(a.GetAppDeploymentHistory))
}
