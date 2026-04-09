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

### Install the latest version of the plugin

The latest stable release of the `akamai-functions-mcp` plugin can be installed like so:

```bash
spin plugins update
spin plugins install akamai-functions-mcp
```

### Install the canary version of the plugin

The canary release of the `akamai-functions-mcp` represents the most recent commits on main and may not be stable, with some features still in progress.

```bash
spin plugins install --url https://github.com/ThorstenHans/akamai-functions-mcp-server/releases/download/canary/akamai-functions-mcp.json
```

### Install from a local build

Alternatively, use the `spin pluginify` plugin to install from a fresh build. This will use the `pluginify` manifest (`spin-pluginify.toml`) to package the plugin and proceed to install it:

```bash
spin plugins install pluginify
go build -o akamai-functions-mcp main.go
spin pluginify --install
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

## Available Resources 📄

- The `spin aka` command Reference - Returns the full command reference as markdown