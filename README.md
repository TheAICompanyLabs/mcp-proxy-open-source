# 🛡️ Universal MCP Proxy 
**An open-source, zero-latency security firewall for the Model Context Protocol (MCP).**

[![Go Version](https://img.shields.io/badge/Go-1.21+-00ADD8?style=flat&logo=go)](https://golang.org)
[![License](https://img.shields.io/badge/License-Apache%202.0-blue.svg)](https://opensource.org/licenses/Apache-2.0)
[![Platform](https://img.shields.io/badge/Platform-macOS%20%7C%20Linux%20%7C%20Windows-lightgray)](#installation)

Local AI agents executing arbitrary tool commands is a critical security vulnerability. Running a standard MCP server gives language models unchecked access to your filesystem, databases, and internal APIs. 

**MCP Proxy sits exactly in the middle.** It intercepts JSON-RPC payloads in microseconds, enforcing strict, human-readable YAML policies before any command reaches the server. 

![Demo GIF Placeholder: Show a red terminal alert blocking a 'rm -rf /' command]
*(Add a 5-second GIF here showing the terminal UI blocking a malicious command)*

## ⚡ Core Capabilities

*   **Runtime Guardrails:** Intercepts and validates every tool execution request in real-time.
*   **Vulnerability Protection:** Natively blocks SSRF attempts, stops directory traversal payloads, and detects prompt injection signatures.
*   **Interactive TUI:** Built with a high-performance terminal interface that acts as a physical circuit breaker, requiring explicit human-in-the-loop approvals for sensitive actions.
*   **True Air-Gapped Execution:** Zero outbound network calls. 100% of Layer 1 and Layer 2 security runs locally on your machine.

## 🚀 One-Line Installation

**macOS & Linux**
```bash
curl -sSfL https://raw.githubusercontent.com/TheAICompanyLabs/mcp-proxy-open-source/main/scripts/install.sh | bash
```
**Windows (PowerShell)**
```
irm https://raw.githubusercontent.com/TheAICompanyLabs/mcp-proxy-open-source/main/scripts/install.ps1 | iex
```
