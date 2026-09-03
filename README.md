# 🛡️ Universal MCP Proxy 
**An open-source, zero-latency security firewall for the Model Context Protocol (MCP).**

[![Go Version](https://img.shields.io/badge/Go-1.21+-00ADD8?style=flat&logo=go)](https://golang.org)
[![License](https://img.shields.io/badge/License-Apache%202.0-blue.svg)](https://opensource.org/licenses/Apache-2.0)
[![Platform](https://img.shields.io/badge/Platform-macOS%20%7C%20Linux%20%7C%20Windows-lightgray)](#installation)

Local AI agents executing arbitrary tool commands is a critical security vulnerability. Running a standard MCP server gives language models unchecked access to your filesystem, databases, and internal APIs. 

**MCP Proxy sits exactly in the middle.** It intercepts JSON-RPC payloads via standard input/output (stdio) in microseconds, enforcing strict, human-readable YAML policies before any command executes.

![MCP Proxy blocking a malicious command](./assets/demo.gif)

## ⚡ Core Architecture & Capabilities

*   **Runtime Guardrails:** Parses JSON-RPC payloads in real-time, verifying arguments against your defined rules.
*   **Universal Protection:** Natively blocks SSRF (Server-Side Request Forgery), directory traversal (`../../`), and known prompt injection signatures.
*   **Fail-Closed Security:** If a policy is malformed or missing, the proxy defaults to zero-trust, completely isolating the underlying server.
*   **True Air-Gapped Execution:** 100% of the Open Core security engine runs locally. Zero outbound telemetry or network calls.

## 🚀 Installation

**macOS & Linux**
```bash
curl -sSfL https://raw.githubusercontent.com/TheAICompanyLabs/mcp-proxy-open-source/main/scripts/install.sh | bash
```


**Windows (PowerShell)**
```powershell
irm https://raw.githubusercontent.com/TheAICompanyLabs/mcp-proxy-open-source/main/scripts/install.ps1 | iex
```

## 💻 Universal Quick Start
MCP Proxy acts as a wrapper around any standard MCP server.
Syntax:
```
mcp-proxy <your-standard-mcp-server-startup-command>
```
# Examples: 

# 1. Securing a local SQLite database server
```
mcp-proxy uvx mcp-server-sqlite --db-path ./local.db
```
# 2. Securing a filesystem server
```
mcp-proxy npx -y @modelcontextprotocol/server-filesystem /path/to/safe/dir
```
# 3. Securing a custom Python server
```
mcp-proxy python3 main.py
```

## 📖 The Policy Cookbook
Security configurations should live alongside your code. When you run mcp-proxy for the first time, it generates a default policy.yaml in your working directory.
👉 View the complete Policy Cookbook here for copy-paste templates covering Databases, File Systems, and Web APIs.
Sample policy.yaml:
```
version: "1.0"
default_action: deny # Enforce Zero-Trust by default

rules:
  # Example: Allow the agent to read files, but strictly block directory traversal
  - tool: "read_file"
    action: allow
    conditions:
      - path_matches: "^/safe/workspace/.*"
      - block_traversal: true 

  # Example: Require explicit human approval (via TUI) for any database writes
  - tool: "query_database"
    action: require_approval
    conditions:
      - contains_regex: "(?i)(INSERT|UPDATE|DELETE|DROP)"
```

## ⚠️ Troubleshooting OS Warnings
Because this is a newly compiled security binary, your operating system may flag it initially.
macOS "Unidentified Developer" Error:
macOS Gatekeeper may block the binary from running. To explicitly trust the proxy, run:
```
sudo xattr -d com.apple.quarantine /usr/local/bin/mcp-proxy
```

Windows SmartScreen Warning:
If Windows Defender prompts "Windows protected your PC", click More info -> Run anyway. Alternatively, unblock the downloaded .exe via PowerShell:
```
Unblock-File -Path "C:\mcp-proxy\mcp-proxy.exe"
```


## 🏢 Enterprise Tier
Need centralized policy distribution, cryptographic compliance logs, and team-wide telemetry across your MLOps pipeline? Visit TheAICompanyLabs for our upcoming Enterprise Dashboard.

