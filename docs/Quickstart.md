# 🚀 Step-by-Step Quickstart Guide

This guide will walk you through installing the Terminal Circuit Breaker, verifying your configuration, and triggering your first secure, zero-prompt autonomous agent workflow.

## Step 1: Verify Prerequisites

Terminal Circuit Breaker runs entirely locally and requires zero dependencies other than a standard terminal.

Ensure you are operating under a standard user account (do not run these commands as `root` or `Administrator`).

## Step 2: Install the Circuit Breaker

Download and install the pre-compiled binary for your operating system.

**macOS & Linux**

Open your terminal and run the installation script:

```bash
curl -sSfL https://raw.githubusercontent.com/TheAICompanyLabs/mcp-proxy-open-source/main/scripts/install.sh | bash
```

**Windows (PowerShell)**

Open a non-administrator PowerShell window and run:

```powershell
irm https://raw.githubusercontent.com/TheAICompanyLabs/mcp-proxy-open-source/main/scripts/install.ps1 | iex
```

## Step 3: Verify Installation

Confirm the proxy was added to your system path successfully:

```bash
mcp-proxy --version
```

*(If the command is not found, restart your terminal session to refresh your path variables.)*

## Step 4: Configure Your AI Client

You must route your AI client's MCP traffic through the proxy. We support all major clients:

1. Open your client's MCP configuration file (e.g., `claude_desktop_config.json` for Claude, or `mcp.json` for Cursor).
2. Set the `command` to `mcp-proxy`.
3. Pass your actual server command (like `npx` or `python`) into the `args` array.

👉 [View exact configuration file paths for Claude, Cursor, Windsurf, Cline, and LM Studio here](https://github.com/TheAICompanyLabs/mcp-proxy-open-source/blob/main/docs/integration.md).

**CRITICAL:** Fully quit and restart your AI client after saving the configuration file.

## Step 5: Trigger the First-Run Scaffolding

The Circuit Breaker builds your security profile dynamically.

1. Open your AI client and prompt it to use the tool (e.g., "List the files in my workspace directory").
2. The proxy will intercept this unknown request and trigger a high-contrast Bubble Tea warning in your terminal (or a native OS popup).
3. Approve the request (`Y` or "Always Allow").

**What just happened?** The proxy instantly created a `~/.mcp-proxy/policy.yaml` file in your home directory and permanently whitelisted that specific tool.

## Step 6: Test Zero-Prompt Autonomy

Send a follow-up prompt to your AI client (e.g., "Now read the contents of the README file in that directory").

Because the agent is operating within the bounds of your newly generated policy, the request will complete end-to-end with zero popups or terminal interruptions. Your agent is now secure and autonomous.

**Next Steps:** Ready to lock down databases, cloud infrastructure, and git repositories? Explore our [Policy Cookbook](https://github.com/TheAICompanyLabs/mcp-proxy-open-source/blob/main/cookbook/README.md).
