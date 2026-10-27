# Runtime Consequence Firewall (RCF) - Enterprise Air-Gapped Sandbox

This repository provides a fully localized, zero-dependency staging environment to validate RCF's <1ms payload interception and SEC 17a-4 WORM audit logging against a simulated enterprise backend.

## 1. Start the Environment
Run the entire architecture (Backend + RCF Proxy + Load Tester) locally:
```bash
docker compose up --build

```

## 2. Run this `curl` against the RCF proxy to watch the wazero micro-VM intercept the `override_query` parameter in <1ms:

```
bash
curl -X POST http://localhost:8080/api/v3/mcp/orchestration/execute \
-H "Authorization: Bearer valid_wso2_oauth_token" \
-H "Content-Type: application/json" \
-d '{
  "agent_id": "adk_checkout_agent_09",
  "tool_name": "process_payment_gateway",
  "arguments": {
    "cart_id": "99812",
    "override_query": "SELECT full_pan, cvv FROM aurora_pci_vault WHERE tx_date = CURRENT_DATE;"
  }
}'
```
