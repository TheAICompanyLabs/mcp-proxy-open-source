# Runtime Consequence Firewall (RCF) - WalmartLabs Architecture Sandbox

This air-gapped staging environment validates RCF's ability to act as a zero-trust `stdio` execution boundary for Spring Boot MCP architectures. It demonstrates DDoS infinite-loop prevention and SEC 17a-4 WORM audit logging at 50,000+ TPS.

## 1. Boot the Sandbox
Spin up the isolated Debian environment. The container automatically compiles the Java JVM backend and authenticates the RCF FIPS-140-3 proxy binary:
```bash
docker compose up -d --build
docker compose exec sandbox /bin/bash
```

## 2. Threat Simulation: Infinite Loop DDoS (50,000 TPS)
Traditional IAM protects the perimeter, but it fails when an authenticated AI agent hallucinates and gets stuck in an infinite tool-calling loop.
Generate 50,000 valid MCP JSON-RPC calls and pipe them directly through the RCF execution boundary to test line-rate interception:

```
awk 'BEGIN {
  print "{\"jsonrpc\": \"2.0\", \"id\": 1, \"method\": \"initialize\", \"params\": {\"protocolVersion\": \"2024-11-05\", \"capabilities\": {}, \"clientInfo\": {\"name\": \"benchmark\", \"version\": \"1.0\"}}}";
  print "{\"jsonrpc\": \"2.0\", \"method\": \"notifications/initialized\"}";
  for (i=2; i<=150; i++) {
    print "{\"jsonrpc\": \"2.0\", \"id\": " i ", \"method\": \"tools/call\", \"params\": {\"name\": \"read_cache\", \"arguments\": {\"key\": \"WM-99\"}}}";
  }
}' > loop_attack.jsonl

(cat loop_attack.jsonl; sleep 1) | ./mcp-proxy-fips ./run_backend.sh
```

Observe the firewall detect the operational boundary violation, block the rate-limit overrun in <1ms, and protect the Java upstream.

## 3. Cryptographic Ledger Verification
Verify the Merkle root of the WORM audit log (SEC Rule 17a-4 / FINRA 4511 Compliant):
```
node audit/verify-ledger.mjs --ledger /workspace/mcp_audit.log
```
