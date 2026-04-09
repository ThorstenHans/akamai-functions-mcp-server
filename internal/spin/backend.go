package spin

import "context"

type Account struct {
	Id   string `json:"id" jsonschema:"Unique Akamai Functions account identifier"`
	Name string `json:"name" jsonschema:"The name of the Akamai Functions account"`
}

type App struct {
	Id   string `json:"id" jsonschema:"Identifier of the application deployed to Akamai Functions"`
	Name string `json:"name" jsonschema:"Name of the application deployed to Akamai Functions"`
}

type AppStatus struct {
	Id          string   `json:"id"`
	Name        string   `json:"name"`
	Urls        []string `json:"urls"`
	CreatedAt   string   `json:"created_at"`
	Invocations int      `json:"invocations"`
}

type AppHistory struct {
	EventType string `json:"event_type"`
	Version   int    `json:"version"`
	Timestamp string `json:"timestamp"`
}

type Backend interface {
	ListAccounts(ctx context.Context) ([]Account, error)
	ListApps(ctx context.Context, accountId string) ([]App, error)
	GetAppStatus(ctx context.Context, accountId string, appId string, appName string) (*AppStatus, error)
	GetAppHistory(ctx context.Context, accountId string, appId string, appName string) ([]AppHistory, error)
	GetAppLogs(ctx context.Context, maxLines int, accountId string, appId string, appName string) ([]string, error)
	DeployApp(ctx context.Context, variables []string, isFirstTimeDeployment bool, accountId string, appId string, appName string) ([]string, error)
	LinkApp(ctx context.Context, accountId string, appId string, appName string) error
	UnlinkApp(ctx context.Context, accountId string) error
}

func getExtraArgs(accountId string, appId string, appName string) []string {
	result := []string{}
	if len(accountId) > 0 {
		result = append(result, "--account-id", accountId)
	}
	if len(appId) > 0 {
		result = append(result, "--app-id", appId)
	} else if len(appName) > 0 {
		result = append(result, "--app-name", appName)
	}
	return result
}
