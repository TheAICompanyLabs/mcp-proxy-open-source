# 📖 The Policy Cookbook

> **Secure, copy-paste guardrails for your Terminal Circuit Breaker.**

The Terminal Circuit Breaker protects your machine by enforcing strict, human-readable YAML policies before any Model Context Protocol (MCP) tool executes. 

This cookbook provides production-ready recipes for the most common attack vectors and MCP servers. 

---

## 🧭 How to Use This Cookbook

1. **Find your use case:** Browse the recipes below based on the MCP server you are running.
2. **Copy the YAML:** Copy the configuration block.
3. **Paste & Reload:** Paste it into your `~/.mcp-proxy/policy.yaml` file. The Circuit Breaker hot-reloads instantly.

---

## 📁 Recipe 1: The Local Filesystem Sandbox

**Target Server:** `@modelcontextprotocol/server-filesystem`

**Tags:** `[Beginner]` `[High-Utility]`

### 🛑 The Problem: Path Traversal

If you allow an agent to read files, a hallucination or prompt injection can trick it into reading sensitive SSH keys or environment variables outside your project folder.

### ✅ The Solution: Strict Scoping & Read-Only Fallbacks

This policy restricts all file reads strictly to your designated workspace and explicitly blocks dangerous payload paths.

```yaml
policies:
  read_file:
    action: ALLOW
    message: "Permitted read access within project bounds"
    validation:
      # Block directory traversal and sensitive OS directories
      reject_patterns:
        - "(?i)(\\.\\./)"
        - "^/etc/.*"
        - "^~/.ssh/.*"
        - "^.*\\.env$"
  
  write_file:
    action: REQUIRE_APPROVAL
    message: "Filesystem mutations always require human approval"
```

## 🗄️ Recipe 2: PostgreSQL Guardrails

**Target Server:** postgres-mcp-server

**Tags:** [Advanced] [High-Risk]

### 🛑 The Problem: Catastrophic Data Loss

When connecting an LLM to a database, you want it to query data, not delete it. Standard MCP implementations cannot differentiate between a SELECT query and a DROP TABLE command inside a raw execution tool.

### ✅ The Solution: Destructive SQL Interception

This policy allows standard schema exploration, but uses Layer 1 Argument Sanitization to intercept destructive SQL commands such as DELETE, DROP, and TRUNCATE before they reach your database.

```yaml
policies:
  execute_query:
    action: ALLOW
    validation:
      # Intercept and block destructive SQL operations
      reject_patterns:
        - "(?i)\\b(DROP|TRUNCATE|DELETE|ALTER|GRANT|REVOKE)\\b"
  
  list_tables:
    action: ALLOW
    message: "Schema exploration is safe"
```

## 🌐 Recipe 3: Cloud SSRF Protection

**Target Server:** @modelcontextprotocol/server-fetch

**Tags:** [Intermediate] [Security-Critical]

#### 🛑 The Problem: Server-Side Request Forgery (SSRF)

Agents equipped with fetch tools can be manipulated into pinging internal cloud metadata endpoints (like AWS 169.254.169.254) to exfiltrate temporary IAM credentials.

### ✅ The Solution: Internal IP Blacklisting

Block the agent from resolving any private, link-local, or loopback IP ranges.


```yaml
policies:
  fetch_url:
    action: ALLOW
    validation:
      reject_patterns:
        - "169\\.254\\.169\\.254"
        - "^https?://(10\\.|192\\.168\\.|172\\.(1[6-9]|2[0-9]|3[0-1])\\.)"
        - "^https?://(localhost|127\\.0\\.0\\.1)"
```

## 🏢 Enterprise Tier: MCP Sentinel

Are you managing a fleet of developers? Managing local .yaml files across hundreds of laptops does not scale.

MCP Sentinel (Phase 2) is our upcoming centralized control plane. It allows security teams to author these exact cookbook recipes in a centralized dashboard and sync them fleet-wide in real time. Learn more about Enterprise here.
