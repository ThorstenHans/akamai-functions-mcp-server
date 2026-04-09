# 🚀 Akamai Functions MCP Server

[![Build Status](https://github.com/ThorstenHans/akamai-functions-mcp-server/actions/workflows/release.yml/badge.svg)](https://github.com/ThorstenHans/akamai-functions-mcp-server/actions)
[![Go Report Card](https://goreportcard.com/badge/github.com/ThorstenHans/akamai-functions-mcp-server)](https://goreportcard.com/report/github.com/ThorstenHans/akamai-functions-mcp-server)
[![Spin Plugin](https://img.shields.io/badge/Spin-Plugin-9d31f0?style=flat-square&logo=webassembly)](https://github.com/spinframework/spin)
![Go](https://img.shields.io/badge/Go-00ADD8?style=flat-square&logo=go&logoColor=white)

Unlock the power of **Akamai Functions** within your AI Assistant. This MCP (Model Context Protocol) server allows tools like Claude Desktop or Cursor to interact with your edge applications using natural language.

Built in **Go** and delivered as a **Spin Plugin**.

---

## 🛑 Prerequisites

Before adding the server to your AI client, ensure your local environment is ready:

1. **Plugin Installed**: You need the `aka` plugin for Spin:

    ```bash
    spin plugins update
    spin plugins install aka --yes
    ```

2. **Login Required**: You must be authenticated via the `aka` plugin for `spin:

    ```bash
    spin aka login
    ```

---

## 📦 Installation

### Install with Homebrew (recommended)

Install via Homebrew by adding the Akamai Developers tap and installing the formula:

```bash
brew tap akamai-developers/tap
brew install akamai-developers/tap/akamai-functions-mcp
```

### Download a precompiled binary

If Homebrew is not available (for example on Windows), download a precompiled binary [from the latest release](https://github.com/akamai-developers/akamai-functions-mcp-server/releases/latest).

### Build from source

Build the binary locally using Go:

```bash
go build -o akamai-functions-mcp main.go
```

## ⚙️ Configuration

To let your AI Assistant use the server, add the `spin` command to your configuration. Optionally, enable debug logging by appending the `--debug` flag to the `args` array.

### Claude Desktop

Update your `claude_desktop_config.json` (found in `~/Library/Application Support/Claude/` on macOS or `%APPDATA%\Claude\` on Windows):

```json
{
  "mcpServers": {
    "akamai-functions": {
      "command": "spin",
      "args": ["akamai-functions-mcp"]
    }
  }
}
```

### Cursor / VS Code

- Navigate to Settings > MCP
- Click + Add Server
- Name: `Akamai Functions`
- Type: `command`
- Command: `spin`
- Args: `akamai-functions-mcp`

## 🤖 Example Interactions

Your AI now acts as a technical assistant for your edge infrastructure. Try asking:

- "Show me all applications in my 'staging' account."
- "What's the current status and URL for the 'edge-api' app?"
- "Check the logs for 'image-processor' look for any 404 errors."
- "List my available Akamai Functions accounts."

## 🛠 Available Tools

| Tool | Description |
| ---- | ----------- |
| `list_accounts` | Fetch all Akamai Functions accounts associated with your login. |
| `list_apps` | Inventory of deployed applications (filterable by account/name). |
| `get_app_status` | Detailed health, deployment history, and active URLs. |
| `get_app_url` | Retrieve the public URL of an application deployed to Akamai Functions |
| `get_app_logs` | Real-time execution log retrieval. |
| `get_app_history` | Retrieve deployment history of an application deployed to Akamai Functions |
| `search_app` | Find an application by a search term (`query`). This command iterates over all Akamai Functions accounts you've access to |
| `deploy_app` | Deploys the application to Akamai Functions. |
| `link_app` | Links the current local workspace to an existing Akamai Functions application. |
| `unlink_app`| |

## Available Resources 📄

- The `spin aka` command Reference - Returns the full command reference as markdown