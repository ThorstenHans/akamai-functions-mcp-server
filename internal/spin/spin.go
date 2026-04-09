package spin

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"os"
	"os/exec"
	"path/filepath"
	"strings"

	"github.com/pelletier/go-toml/v2"
)

const spinBinary = "spin"

type spinBackend struct {
	logger *log.Logger
}

func NewSpinBackend(logger *log.Logger) Backend {
	return &spinBackend{
		logger: logger,
	}
}

type spinAccountInfoResponse struct {
	AuthInfo spinAuthInfo `json:"auth_info"`
}

type spinAuthInfo struct {
	Accounts []Account `json:"accounts"`
}

func (b *spinBackend) ListAccounts(ctx context.Context) ([]Account, error) {
	command := []string{"aka", "info", "--format", "json"}
	b.logger.Printf("Will retrieve all your Akamai Functions accounts using the following spin command: %v", command)
	out, err := b.run(command...)
	if err != nil {
		b.logger.Printf("Error listing accounts: %v\n", err)
		return nil, err
	}
	var result spinAccountInfoResponse
	err = json.Unmarshal(out, &result)
	if err != nil {
		return nil, err
	}
	return result.AuthInfo.Accounts, nil
}

func (b *spinBackend) ListApps(ctx context.Context, accountId string) ([]App, error) {
	command := []string{"aka", "apps", "list", "--format", "json"}
	if len(accountId) > 0 {
		command = append(command, "--account-id", accountId)
	}
	b.logger.Printf("Will retrieve all apps for the specified account using the following spin command: %v\n", command)
	out, err := b.run(command...)
	if err != nil {
		b.logger.Printf("Error listing apps for account %s: %v\n", accountId, err)
		return nil, err
	}
	var apps []App
	err = json.Unmarshal(out, &apps)
	if err != nil {
		b.logger.Printf("Error unmarshalling output of command to list apps for account %s: %v\nOutput was: %s\n", accountId, err, string(out))
		return nil, err
	}
	return apps, nil
}

func (b *spinBackend) GetAppHistory(ctx context.Context, accountId string, appId string, appName string) ([]AppHistory, error) {
	command := []string{"aka", "app", "history", "--format", "json"}
	extraArgs := getExtraArgs(accountId, appId, appName)
	if len(extraArgs) > 0 {
		command = append(command, extraArgs...)
	}
	b.logger.Printf("Will retrieve deployment history using the following spin command: %v\n", command)
	out, err := b.run(command...)
	if err != nil {
		b.logger.Printf("Error running command to get app deployment history: %v\nOutput was: %s\n", err, string(out))
		return nil, err
	}
	var history []AppHistory
	err = json.Unmarshal(out, &history)
	if err != nil {
		return nil, err
	}
	return history, nil
}

func (b *spinBackend) GetAppLogs(ctx context.Context, maxLines int, accountId string, appId string, appName string) ([]string, error) {
	command := []string{"aka", "logs"}
	if maxLines < 0 {
		return nil, fmt.Errorf("maxLines cannot be negative")
	}
	extraArgs := getExtraArgs(accountId, appId, appName)
	if len(extraArgs) > 0 {
		command = append(command, extraArgs...)
	}
	command = append(command, "--max-lines", fmt.Sprintf("%d", maxLines))
	out, err := b.run(command...)
	if err != nil {
		b.logger.Printf("Error running command to get app logs: %v\nOutput was: %s\n", err, string(out))
		return nil, err
	}
	logs := strings.Split(string(out), "\n")
	logs = append([]string{}, logs...)
	for i := len(logs) - 1; i >= 0; i-- {
		if strings.TrimSpace(logs[i]) == "" {
			logs = append(logs[:i], logs[i+1:]...)
		}
	}
	return logs, nil
}

func (b *spinBackend) GetAppStatus(ctx context.Context, accountId string, appId string, appName string) (*AppStatus, error) {
	command := []string{"aka", "app", "status", "--format", "json"}
	extraArgs := getExtraArgs(accountId, appId, appName)
	if len(extraArgs) > 0 {
		command = append(command, extraArgs...)
	}
	b.logger.Printf("Will retrieve app status using the following spin command: %v\n", command)
	out, err := b.run(command...)
	if err != nil {
		b.logger.Printf("Error running command to get app status: %v\nOutput was: %s\n", err, string(out))
		return nil, err
	}
	var status AppStatus
	err = json.Unmarshal(out, &status)
	if err != nil {
		b.logger.Printf("Error unmarshalling output of command to get app status: %v\nOutput was: %s\n", err, string(out))
		return nil, err
	}
	return &status, nil
}
func (b *spinBackend) DeployApp(ctx context.Context, variables []string, isFirstTimeDeployment bool, accountId string, appId string, appName string) ([]string, error) {
	// 1. Build base command
	command := []string{"aka", "deploy", "--build", "--no-confirm"}
	if len(accountId) > 0 {
		command = append(command, "--account-id", accountId)
	}
	// if it is a first time deployment, we must provide the --create-name flag
	if isFirstTimeDeployment {
		command = append(command, "--create-name", appName)
	} else {
		// check if there is either an appId provided or if  .spin-aka/config.toml file exists, if not return an error and tell the LLM that
		// it must link the app first
		// otherwise add an app-id flag
		if appId != "" {
			command = append(command, "--app-id", appId)
		} else {
			localConfig, err := getLocalAkaConfig()
			if err != nil {
				return nil, fmt.Errorf("safety abort: no app ID provided and no local config detected. If you want to create a new app, set isFirstTimeDeployment to true. If you want to deploy to an existing unlinked app, provide the App ID.")
			}
			command = append(command, "--app-id", localConfig.Application.ID)
		}
	}

	if len(variables) > 0 {
		for _, val := range variables {
			command = append(command, "--variable", val)
		}
	}

	b.logger.Printf("Executing deploy: spin %v\n", command)

	// 5. Execute Command
	out, err := b.run(command...)
	if err != nil {
		b.logger.Printf("CLI Error: %s\n", string(out))
		return nil, fmt.Errorf("Deployment failed: %v\nOutput: %s", err, string(out))
	}

	// 6. Format Output
	logs := strings.Split(strings.TrimSpace(string(out)), "\n")
	return logs, nil
}

func (b *spinBackend) LinkApp(ctx context.Context, accountId string, appId string, appName string) error {

	command := []string{"aka", "app", "link"}

	// 1. Append the required target app
	if appId != "" {
		command = append(command, "--app-id", appId)
	} else if appName != "" {
		command = append(command, "--app-name", appName)
	}

	// 2. Append optional account context
	if accountId != "" {
		command = append(command, "--account-id", accountId)
	}

	b.logger.Printf("Executing link: spin %v\n", command)
	out, err := b.run(command...)
	if err != nil {
		b.logger.Printf("Error linking app: %v\nOutput was: %s\n", err, string(out))
		return err
	}

	return nil
}

func (b *spinBackend) UnlinkApp(ctx context.Context, accountId string) error {
	command := []string{"aka", "app", "unlink"}

	if accountId != "" {
		command = append(command, "--account-id", accountId)
	}

	b.logger.Printf("Executing unlink: spin %v\n", command)
	out, err := b.run(command...)
	if err != nil {
		b.logger.Printf("Error unlinking app: %v\nOutput was: %s\n", err, string(out))
		return err
	}

	return nil
}

func (b *spinBackend) run(command ...string) ([]byte, error) {

	cmd := exec.Command(spinBinary, command...)
	spinOut, err := cmd.Output()
	if err != nil {
		return nil, err
	}
	return spinOut, nil

}

type AkaConfig struct {
	Application struct {
		ID string `toml:"id" json:"id"`
	} `toml:"application" json:"application"`
}

func getLocalAkaConfig() (*AkaConfig, error) {
	cwd, err := os.Getwd()
	if err != nil {
		return nil, fmt.Errorf("failed to get working directory: %w", err)
	}

	configPath := filepath.Join(cwd, ".spin-aka", "config.toml")

	// Read the file content first
	data, err := os.ReadFile(configPath)
	if err != nil {
		return nil, fmt.Errorf("failed to read config file: %w", err)
	}

	var config AkaConfig
	// Use Unmarshal (standard for v2)
	if err := toml.Unmarshal(data, &config); err != nil {
		return nil, fmt.Errorf("failed to parse TOML: %w", err)
	}

	return &config, nil
}
