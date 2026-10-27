# Runtime Consequence Firewall (RCF) - WalmartLabs Staging Sandbox

This repository provides a fully localized staging environment to validate RCF's <1ms payload interception against a simulated Walmart Cart & Checkout backend.

**Note:** This repository uses the compiled Linux binary `mcp-proxy-fips`. It requires zero dependencies other than Docker.

## 1. Start the Environment
```bash
docker compose up --build
