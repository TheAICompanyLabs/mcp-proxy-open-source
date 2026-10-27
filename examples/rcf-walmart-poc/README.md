# Runtime Consequence Firewall (RCF) - WalmartLabs Architecture Sandbox

This air-gapped staging environment validates RCF's ability to act as a zero-trust `stdio` execution boundary for Spring Boot MCP architectures.

## 1. Secure Authentication & Boot
Export your provided Enterprise Staging License Key to authenticate the FIPS-140-3 proxy binary, then boot the sandbox.
```bash
export MCP_LICENSE_KEY="<PASTE_YOUR_PROVIDED_KEY_HERE>"
docker compose up -d --build
docker compose exec sandbox /bin/bash

```

## 2. Execute the Autonomous Agent Simulation
Simulate a compromised agent processing 20+ legitimate cart validations before attempting a post-auth SQL injection.

```
./build_session.sh | ./mcp-proxy-fips ./run_backend.sh

```

Observe the firewall successfully route the 20+ legitimate calls, dynamically intercept the SELECT anomaly on call 23, and sever the connection to protect the Aurora database.

## 3. Cryptographic Ledger Verification
Verify the Merkle root of the WORM audit log (SEC Rule 17a-4 / FINRA 4511 Compliant) to prove all events were cryptographically sealed:
```
node audit/verify-ledger.mjs --ledger /workspace/mcp_audit.log
```
