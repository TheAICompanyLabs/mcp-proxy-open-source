# 📖 The Policy Cookbook

> **Production-grade, copy-paste guardrails for your Terminal Circuit Breaker.**

The Terminal Circuit Breaker protects your environment by evaluating every Model Context Protocol (MCP) JSON-RPC payload against strict, human-readable YAML policies before any execution occurs.

This cookbook contains audited, modular policy suites categorized by domain and threat profile.

---

## ⚡ Quickstart: The Universal Baseline

If you just want an immediate, multi-domain starter pack covering the top 6 attack vectors in one file, use our root starter template:

👉 **[View Root Unified Baseline (`./policy.yaml`)](./policy.yaml)**

```bash
# Copy the unified baseline into your active proxy directory
cp cookbook/policy.yaml ~/.mcp-proxy/policy.yaml
```

---

## 🗂️ Domain-Specific Policy Recipes

Each folder contains a targeted `policy.yaml` fine-tuned for specific MCP servers, threat vectors, and risk profiles.

| Domain Recipe | Target Servers | Threat Vectors Neutralized | Risk Tier |
|---|---|---|---|
| 1. Filesystem & OS | `@modelcontextprotocol/server-filesystem`, `bash-mcp`, `os-mcp` | Directory traversal (`../../`), credential harvesting (`.env`, `id_rsa`), arbitrary shell execution, and disk wipes. | **CRITICAL** |
| 2. Database & Data | `@modelcontextprotocol/server-postgres`, `server-sqlite`, `mysql-mcp` | Unchecked `DROP TABLE`, unindexed mass `DELETE`/`UPDATE` without `WHERE` clauses, and SQL injection chaining via `;`. | **CRITICAL** |
| 3. Web & SSRF | `@modelcontextprotocol/server-fetch`, `server-puppeteer`, `playwright-mcp` | Cloud metadata exfiltration (`169.254.169.254`), private subnet port scanning, and browser token extraction (`document.cookie`). | **HIGH** |
| 4. Git & Code | `@modelcontextprotocol/server-github`, `git-mcp`, `gitlab-mcp` | Force-pushes (`--force`), pushing directly to `main`/`master`, unauthorized repo deletions, and autonomous self-merging. | **HIGH** |
| 5. Cloud & DevOps | `docker_mcp`, `kubernetes-mcp`, `aws-mcp`, `terraform-mcp` | Privileged container creation (`--privileged`), root volume mounting (`-v /:`), crypto-mining instance sizing, and `terraform destroy`. | **CRITICAL** |
| 6. Productivity & Comms | `slack-mcp`, `gmail-mcp`, `google-drive-mcp`, `notion-mcp`, `cal-mcp` | Unauthorized `#all-hands` or `@everyone` Slack blasts, competitor/investor email leaks, and mass deletion of Drive or Notion tables. | **MEDIUM** |

---

## 🛠️ Deep Dive: Recipe Breakdown

### 1. Filesystem & OS Guardrails
**Target Servers:** `@modelcontextprotocol/server-filesystem`, `bash-mcp`, `os-mcp`
**Default Posture:** `deny`

**Key Controls:**
- **ALLOW:** Read-only queries scoped strictly inside `/opt/safe_workspace/*`.
- **BLOCK:** Auto-rejects any path containing `../`, `.env`, `.aws`, `.ssh`, `/etc/shadow`, `/etc/passwd`, or `id_rsa`.
- **REQUIRE_APPROVAL:** File creation stripped of executable permissions (`.sh`, `.exe`, `.bat`, `.ps1`).
- **INTERCEPT:** Blocks `curl`, `wget`, `nc`, `rm -rf`, `sudo`, and `chmod` in bash tools.

### 2. Database & Data Guardrails
**Target Servers:** `@modelcontextprotocol/server-postgres`, `server-sqlite`, `mysql-mcp`
**Default Posture:** `deny`

**Key Controls:**
- **ALLOW:** Schema inspection (`list_tables`, `describe_schema`) and pure read operations matching `^(SELECT|EXPLAIN)`.
- **BLOCK:** Chained SQL injections attempting to append `; DROP`, `; DELETE`, or `; ALTER`.
- **REQUIRE_APPROVAL:** Single data mutations; hard-blocks unconstrained queries (e.g., `DELETE FROM users;` without a `WHERE` clause).

### 3. Web & SSRF Guardrails
**Target Servers:** `@modelcontextprotocol/server-fetch`, `server-puppeteer`, `playwright-mcp`, `firecrawl-mcp`
**Default Posture:** `deny`

**Key Controls:**
- **ALLOW:** Web fetches matching `^https://.*`.
- **BLOCK:** Cloud metadata endpoints (`169.254.169.254`, `metadata.google.internal`), loopback addresses (`127.0.0.1`, `localhost`), and RFC 1918 internal subnets (`10.0.0.0/8`, `172.16.0.0/12`, `192.168.0.0/16`).
- **REQUIRE_APPROVAL:** In-browser JavaScript execution attempts (`evaluate`), blocking access to `document.cookie` or `localStorage`.

### 4. Git & Code Guardrails
**Target Servers:** `@modelcontextprotocol/server-github`, `git-mcp`, `gitlab-mcp`
**Default Posture:** `deny`

**Key Controls:**
- **ALLOW:** Read-only inspect tools (`git_status`, `git_diff`, `github_search_repositories`) and pull request creation.
- **REQUIRE_APPROVAL:** Commit and push operations.
- **HARD DENY:** Pushes containing `main`, `master`, `prod`, or `--force`/`-f`. AI tools cannot merge PRs or delete repositories.

### 5. Cloud & DevOps Guardrails
**Target Servers:** `docker_mcp`, `kubernetes-mcp`, `aws-mcp`, `terraform-mcp`
**Default Posture:** `deny`

**Key Controls:**
- **ALLOW:** Observability commands (`list_containers`, `kubectl_get`, `terraform_plan`).
- **REQUIRE_APPROVAL:** Container launches, blocking `--privileged`, root mounts (`-v /:`), and expensive compute tiers (`p3.8xlarge`).
- **HARD DENY:** Destructive calls (`aws_delete_*`, `kubectl_delete`, `terraform_destroy`, `remove_volume`).

### 6. Productivity & Comms Guardrails
**Target Servers:** `slack-mcp`, `gmail-mcp`, `google-drive-mcp`, `notion-mcp`, `cal-mcp`
**Default Posture:** `deny`

**Key Controls:**
- **ALLOW:** Context retrieval (`slack_get_channel_history`, `drive_search_files`, `cal_get_events`).
- **REQUIRE_APPROVAL:** Outbound messages. Auto-blocks posts directed at `#general`, `#all-hands`, `@channel`, or `@everyone`.
- **HARD DENY:** Irreversible data destruction (`drive_delete_file`, `gmail_delete_message`, `notion_delete_database`).

---

## 🏢 Fleet-Wide Sync: MCP Sentinel (Enterprise)

Running independent `policy.yaml` files across an engineering team leads to config drift.

**MCP Sentinel** is our centralized enterprise plane that lets security teams distribute, version, and cryptographically audit these exact policies across thousands of developer laptops.
