# Runtime Consequence Firewall (RCF) - WalmartLabs Architecture Sandbox

This air-gapped staging environment validates RCF's ability to act as a zero-trust `stdio` execution boundary for Spring Boot MCP architectures. It demonstrates sub-millisecond TOCTOU payload interception and SEC 17a-4 WORM audit logging at 50,000+ TPS.

## 1. Boot the Sandbox
Spin up the isolated Debian environment containing the Java JVM and the RCF FIPS-140-3 proxy binary:
```bash
docker compose up -d --build
docker compose exec sandbox /bin/bash
```
## 2. Compile the Target Upstream (Spring Boot)
Ensure the Java backend is ready to accept the proxy's stdio routing:
```
javac backend/src/main/java/com/walmartlabs/checkout/McpExecutionController.java
chmod +x run_backend.sh
```

## 3. High-Throughput Load Test (50,000 TPS)
Generate 50,000 valid MCP JSON-RPC calls and pipe them directly through the RCF execution boundary to test line-rate throughput latency:
```
awk 'BEGIN {
  print "{\"jsonrpc\": \"2.0\", \"id\": 1, \"method\": \"initialize\", \"params\": {\"protocolVersion\": \"2024-11-05\", \"capabilities\": {}, \"clientInfo\": {\"name\": \"benchmark\", \"version\": \"1.0\"}}}";
  print "{\"jsonrpc\": \"2.0\", \"method\": \"notifications/initialized\"}";
  for (i=2; i<=50000; i++) {
    print "{\"jsonrpc\": \"2.0\", \"id\": " i ", \"method\": \"tools/call\", \"params\": {\"name\": \"read_cache\", \"arguments\": {\"key\": \"WM-99\"}}}";
  }
}' > benchmark.jsonl

time (cat benchmark.jsonl; sleep 1) | ./mcp-proxy-fips ./run_backend.sh
```
## 4. The TOCTOU Exploit (Air-Gapped Interception)
Traditional API gateways validate auth headers but miss post-auth parameter mutations. Execute this 3-part handshake sequence to simulate a compromised ADK agent attempting to bypass RBAC and exfiltrate PCI data via override_query:
```
(cat attack_payload.jsonl; sleep 2) | ./mcp-proxy-fips ./run_backend.sh
```
Observe the proxy intercept the schema anomaly in <1ms, drop the connection, and prevent the Java backend from receiving the SQL injection.

## 5. Cryptographic Ledger Verification
Verify the Merkle root of the WORM audit log (SEC Rule 17a-4 / FINRA 4511 Compliant):
```
node audit/verify-ledger.mjs /tmp/rcf_worm_audit.log
```
