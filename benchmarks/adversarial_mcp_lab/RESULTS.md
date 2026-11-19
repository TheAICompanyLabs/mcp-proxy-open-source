# Adversarial MCP Lab — Un-Gamed Zero-Config OSS Core Results

Run locally from a clean clone (0B cloud egress, no pre-tuning of `inspector.go`, no `policy.yaml` allowlist loaded):

```bash
RUN_ADV_LAB=1 go test -v -count=1 -run TestAdversarialMCPLab .
```

## Summary (Apple M3, Go 1.23+, Zero-Config OSS Core)
* **Red Canary Atomic Red Team (MITRE ATT&CK via MCP):** 8 / 10 deterministically blocked
* **HexStrike-MCP Style Autonomous Offensive Bursts:** 4 / 8 deterministically blocked
* **Total Burst Inspection Time (18 calls):** 374.791µs (~20.8µs avg/call)

## Explicit Misses & Remediation Boundary
The 6 vectors allowed under zero-config OSS defaults demonstrate why deterministic argument scanning must pair with least-privilege tool allowlists (`policy.yaml`) or Wasm payload detonation:

* **Out-of-Scope Offensive Tools (`nmap_scan`, `sqlmap_run`, `ffuf_fuzz`):** Syntactically valid tool arguments—blocked immediately by enabling default-deny tool allowlisting (`* -> REQUIRE_APPROVAL` or `DENY`) in `policy.yaml`.
* **Non-Chained Shell Recon / Cleanup (`history -c && unset HISTFILE`, `grep -rn 'API_KEY'`, `echo ... | base64 -d`):** Requires commercial Wasm pre-flight detonation / entropy inspection or OS-level container confinement (seccomp + read-only rootfs). See `LIMITATIONS.md`.
