# Adversarial MCP Lab — Un-Gamed Zero-Config OSS Core Results

Run locally from a clean clone (0B cloud egress, no pre-tuning of `inspector.go`, no `policy.yaml` allowlist loaded):

```bash
RUN_ADV_LAB=1 go test -v -count=1 -run TestAdversarialMCPLab .
```
