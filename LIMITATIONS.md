# Security Scope, Threat Model & Explicit Limitations

RCF (`mcp-proxy-open-source`) is a deterministic, local Layer-7 JSON-RPC execution boundary for Model Context Protocol (MCP) servers. Security engineering teams evaluating any inline proxy should know both what it enforces and where its boundary ends.

## What This Proxy Deterministically Enforces
1. **Pre-Execution JSON-RPC `tools/call` Parameter Inspection:** Evaluates tool names and structured arguments locally before forwarding to the upstream `stdio` or Streamable HTTP MCP server.
2. **Deterministic Injection & Traversal Blocking:** Flags and blocks path traversal (`../`, URL/double-encoded `%252e`), shell metacharacter injection, destructive SQL DDL patterns, and unauthorized filesystem/command patterns defined in `policy.yaml`.
3. **Tool Schema Pinning (Rug-Pull Detection):** Hashes tool definitions observed during `tools/list` so silent post-approval tool description mutations (OWASP MCP03) are detected locally.
4. **Zero Cloud Egress:** Runs locally inside your host, container, or VPC with 0 bytes of prompt, argument, or credential egress to external SaaS APIs.

## What This Proxy Does NOT Solve (Explicit Boundaries)
1. **Semantic / Reasoning Drift Inside Policy-Allowed Parameters:** Deterministic AST and regex/policy evaluation does **not** detect subtle natural-language manipulation when the resulting tool call is syntactically valid and permitted by your configured `policy.yaml` (e.g., an agent tricked into summarizing a permitted document in a misleading way or emailing an allowlisted recipient with fabricated business text).
2. **Upstream Server Vulnerabilities & Direct Host Bypasses:** If an agent or compromised process invokes OS binaries, sockets, or databases directly—bypassing the MCP boundary—or exploits a memory-safety bug inside the upstream MCP server itself, an MCP-layer proxy cannot intercept that out-of-band path. Pair RCF with OS/container isolation (seccomp, non-root UID, network egress rules).
3. **Upstream Identity & Token Issuance:** RCF validates and enforces tool-call boundaries locally; it does not replace your OAuth 2.1 identity provider, short-lived credential broker, or secret rotation pipeline.
4. **Open-Source Core vs. Commercial Rulepack:** This repository contains the auditable Go proxy core, baseline policy engine, and reproducible local benchmark harness. Maintained enterprise OWASP MCP Top 10 rulepacks, Wasm detonation modules, and AIUC-1 / FINOS AIGF v3 / DORA cryptographic WORM exporters are licensed commercially via [getrcf.dev](https://getrcf.dev/trust).

## Reproducing Our Telemetry
Never trust vendor benchmark cards without running them yourself. Run `make bench` from a clean clone to reproduce P50/P99 latency and deterministic block/allow results on your own hardware.
