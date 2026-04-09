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

	linkAppTool := mcp.NewTool("link_app",
		mcp.WithDescription(`Links the current local workspace to an existing Akamai Functions application. 
USE THIS WHEN:
1. The user explicitly asks to link this folder to an existing app.
2. A deployment or log request fails because the workspace is not linked (workspace has no link to an application deployed to Akamai Functions), and the user wants to attach it to an existing remote app.
IMPORTANT: You must provide the exact App ID or Name. Use 'list_apps' to see all apps and ask the user which remote app should be linked to.`),
		mcp.WithReadOnlyHintAnnotation(false),
		mcp.WithDestructiveHintAnnotation(false),
		mcp.WithInputSchema[LinkAppArgs](),
		mcp.WithOutputSchema[ToolResponse[[]string]](),
	)

	s.AddTool(linkAppTool, mcp.NewStructuredToolHandler(a.LinkApp))

	unlinkAppTool := mcp.NewTool("unlink_app",
		mcp.WithDescription(`Unlinks the current local workspace from the application deployed to Akamai Functions. 
USE THIS WHEN:
1. The user explicitly asks to unlink this workspace from an existing app.
IMPORTANT: This does NOT delete any remote applications, it only removes the local link (the .spin-aka/config.toml file).`),
		mcp.WithReadOnlyHintAnnotation(false),
		mcp.WithDestructiveHintAnnotation(false),
		mcp.WithInputSchema[MaybeByAccountArgs](),
		mcp.WithOutputSchema[ToolResponse[[]string]](),
	)

	s.AddTool(unlinkAppTool, mcp.NewStructuredToolHandler(a.UnlinkApp))

	deployAppTool := mcp.NewTool("deploy_app",
		mcp.WithDescription(`Deploys the application to Akamai Functions. 
SAFETY RULES:
1. You MUST ALWAYS ask the user for confirmation before calling this tool.
2. If the tool fails because no local app is detected, ask the user if they want to create a brand new application, or link to an existing one.
3. Set 'isFirstTimeDeployment' to true ONLY if deploying an app for the first time.`),
		mcp.WithReadOnlyHintAnnotation(false),    // Modifies remote state
		mcp.WithDestructiveHintAnnotation(false), // It overrides, but we protect against accidental creation
		mcp.WithInputSchema[DeployAppArgs](),
		mcp.WithOutputSchema[ToolResponse[[]string]](),
	)

	s.AddTool(deployAppTool, mcp.NewStructuredToolHandler(a.DeployApp))
}
