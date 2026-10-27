# Runtime Consequence Firewall (RCF) - WalmartLabs Architecture Sandbox

This air-gapped staging environment validates RCF's ability to act as a zero-trust `stdio` execution boundary for Spring Boot MCP architectures. It demonstrates sub-millisecond TOCTOU payload interception and SEC 17a-4 WORM audit logging at 50,000+ TPS.

## 1. Boot the Sandbox
Spin up the isolated Debian environment containing the Java JVM and the RCF FIPS-140-3 proxy binary:
```bash
docker compose up -d --build
docker compose exec sandbox /bin/bash
