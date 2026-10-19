# RCF Standardized Benchmark Harness (v2.0)

This directory contains the reproducible test harness and cryptographic WORM (Write Once, Read Many) ledgers proving the 100% containment rate of the Runtime Consequence Firewall (RCF).

## Files Included
* `swarm_simulator.py`: The asynchronous Python harness simulating a 10-agent swarm executing BFCL, InjecAgent, and AgentDojo threat vectors over local `stdio` pipes.
* `audit-ledger.json`: The raw SEC Rule 17a-4 compliant cryptographic ledger generated during the benchmark run.
* `verify-ledger.mjs`: The local verification script to mathematically prove the integrity of the Merkle root hash.

## How to Reproduce the Benchmark
1. Download the compiled, SHA-256 verified `mcp-proxy-fips` binary from the [GitHub Releases](../../releases) page.
2. Export your staging license key: `export MCP_LICENSE_KEY="<YOUR_KEY>"`
3. Run the complete multi-suite benchmark:
   ```bash
   python3 swarm_simulator.py --suite=all

## How to Verify the WORM Ledger

The proxy records all execution boundaries without sending raw prompts or codebase data to the cloud. You can verify the integrity of the attached ledger offline:

```
node verify-ledger.mjs audit-ledger.json
```
