# Runtime Consequence Firewall (RCF) - WalmartLabs Architecture Sandbox

This air-gapped staging environment validates RCF's ability to act as a zero-trust `stdio` execution boundary for Spring Boot MCP architectures. It demonstrates sustained, high-volume agentic tool processing and sub-millisecond TOCTOU interception.

## 1. Boot the Sandbox
Spin up the isolated Debian environment. The container automatically compiles the Java JVM backend and authenticates the RCF FIPS-140-3 proxy binary in `block` mode:
```bash
docker compose up -d --build
docker compose exec sandbox /bin/bash
```

## 2. The TOCTOU Exploit (Air-Gapped Interception)
Traditional IAM protects the perimeter, but it misses post-auth SQL injections initiated by compromised internal agents. Execute the generated session through the proxy:

```
./build_session.sh | ./mcp-proxy-fips ./run_backend.sh

```

Observe the firewall successfully route the 20+ legitimate calls, dynamically intercept the SELECT anomaly on call 23, and sever the connection to protect the Aurora database.

## 3. Cryptographic Ledger Verification
Verify the Merkle root of the WORM audit log (SEC Rule 17a-4 / FINRA 4511 Compliant) to prove all events were cryptographically sealed:
```
node audit/verify-ledger.mjs --ledger /workspace/mcp_audit.log
```
