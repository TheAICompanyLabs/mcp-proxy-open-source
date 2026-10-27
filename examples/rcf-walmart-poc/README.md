## 🚨 The TOCTOU Exploit Payload (Test This Yourself)

Traditional WSO2/API Gateways validate the `Authorization` header and pass the payload downstream. If an ADK Agent is compromised, it can execute the following payload to bypass traditional RBAC. 

Run this `curl` against the RCF proxy to watch the wazero micro-VM intercept the `override_query` parameter in <1ms:

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
