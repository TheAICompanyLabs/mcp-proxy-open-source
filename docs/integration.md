# 🔌 Client Integration Guide

Terminal Circuit Breaker is universally compatible with any Model Context Protocol (MCP) client. Below are the exact setup instructions for the top 5 developer AI platforms.

---

## 1. Claude Desktop
Claude uses a global configuration file to mount servers.

*   **macOS Path:** `~/Library/Application Support/Claude/claude_desktop_config.json`
*   **Windows Path:** `%APPDATA%\Claude\claude_desktop_config.json`
*   **Configuration:**
    ```json
    {
      "mcpServers": {
        "secure-filesystem": {
          "command": "mcp-proxy",
          "args": ["npx", "-y", "@modelcontextprotocol/server-filesystem", "/path/to/workspace"]
        }
      }
    }
    ```
    > **Note:** You must completely quit (Cmd+Q) and restart Claude Desktop for changes to take effect.

---

## 2. Cursor IDE
Cursor allows configuring MCP servers either globally or per-project.

*   **Global Path:** `~/.cursor/mcp.json` (Mac/Linux) or `%USERPROFILE%\.cursor\mcp.json` (Windows).
*   **Project Path:** `.cursor/mcp.json` inside your project directory.
*   **Configuration:** 
    ```json
    {
      "mcpServers": {
        "secure-filesystem": {
          "command": "mcp-proxy",
          "args": ["npx", "-y", "@modelcontextprotocol/server-filesystem", "/path/to/workspace"]
        }
      }
    }
    ```
    > **Alternative:** You can also configure this via the UI by navigating to **Settings > Tools & MCP** and clicking **Add new MCP server**.

---

## 3. Windsurf
Windsurf (by Codeium) reads from its own dedicated config file.

*   **macOS / Linux Path:** `~/.codeium/windsurf/mcp_config.json`.
*   **Windows Path:** `%USERPROFILE%\.codeium\windsurf\mcp_config.json`.
*   **Configuration:**
    ```json
    {
      "mcpServers": {
        "secure-filesystem": {
          "command": "mcp-proxy",
          "args": ["npx", "-y", "@modelcontextprotocol/server-filesystem", "/path/to/workspace"]
        }
      }
    }
    ```
    > **Note:** Check the Cascade MCP settings panel (via the hammer icon) to confirm the server status is active.

---

## 4. Cline (VS Code Extension)
Cline stores its MCP settings deeply within the VS Code global storage directory.

*   **macOS Path:** `~/Library/Application Support/Code/User/globalStorage/saoudrizwan.claude-dev/settings/cline_mcp_settings.json`.
*   **CLI Path:** If using the Cline CLI, use `~/.cline/mcp.json`.
*   **Configuration:**
    ```json
    {
      "mcpServers": {
        "secure-filesystem": {
          "command": "mcp-proxy",
          "args": ["npx", "-y", "@modelcontextprotocol/server-filesystem", "/path/to/workspace"]
        }
      }
    }
    ```
    > **Alternative:** In the Cline VS Code panel, click the **MCP Servers icon** (stacked server icon in the top toolbar) to edit the config directly.

---

## 5. LM Studio
LM Studio allows you to mount servers for local, air-gapped LLM usage.

*   **Path:** Edit `mcp.json` via the LM Studio Developer tab.
*   **Configuration:**
    ```json
    {
      "mcpServers": {
        "secure-filesystem": {
          "command": "mcp-proxy",
          "args": ["npx", "-y", "@modelcontextprotocol/server-filesystem", "/path/to/workspace"]
        }
      }
    }
    ```
    > **Windows Users:** Ensure you specify `.exe` and `.cmd` extensions if required by your environment (e.g., `"command": "mcp-proxy.exe", "args": ["npx.cmd", ...]`).
