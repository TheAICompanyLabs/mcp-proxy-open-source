# RCF Official Pre-Compiled FIPS Binaries (v2.0 Antimatter)

Pre-compiled, zero-dependency Go executables for the Runtime Consequence Firewall (RCF). Built under strict FIPS 140-3 validation mode (`GOFIPS140=v1.0.0 GODEBUG=fips140=only`).

## Available Binaries

| Operating System | Architecture | Binary File |
| :--- | :--- | :--- |
| macOS (Apple Silicon) | arm64 (M1/M2/M3/M4) | [mcp-proxy-fips-darwin-arm64](./mcp-proxy-fips-darwin-arm64) |
| macOS (Intel) | amd64 (x86_64) | [mcp-proxy-fips-darwin-amd64](./mcp-proxy-fips-darwin-amd64) |
| Linux | arm64 (aarch64) | [mcp-proxy-fips-linux-arm64](./mcp-proxy-fips-linux-arm64) |
| Windows | amd64 (x86_64) | [mcp-proxy-fips-windows-amd64.exe](./mcp-proxy-fips-windows-amd64.exe) |

---

## Verifiable SHA-256 Checksums

Run `shasum -a 256 <binary>` locally to verify integrity before execution:

* **mcp-proxy-fips-darwin-arm64:**
  `0e35264563d694c2bebeda48580a2959f5bccb13244280163c12f91d8e139361`
* **mcp-proxy-fips-darwin-amd64:**
  `47229a815e15d4b4458df98370c3239ddfae93939ab244d91502fc9b2595ae0`
* **mcp-proxy-fips-linux-arm64:**
  `665d7034b17fa2b316d6f2145ab071b14e72708a37b27ba4120bb24aa2573891`
* **mcp-proxy-fips-windows-amd64.exe:**
  `b2f6afdb9445c21750aa1ccbdc7b186bde22331f3cf195ac03864dcead0a51cb`
* **mcp-proxy-fips-linux-amd64:**
  `ba40bcb04d8fd6c507dae6f3ea957424f19d52efe6c630911add123d8a723d09`

---

## macOS Gatekeeper Bypass (CLI Quickstart)

Corporate MacBooks attach the `com.apple.quarantine` attribute to files downloaded via browser. Run these commands to execute without quarantine prompts:

```bash
# 1. Remove quarantine flag
xattr -d com.apple.quarantine mcp-proxy-fips-darwin-arm64

# 2. Add execution permission
chmod +x mcp-proxy-fips-darwin-arm64

# 3. Verify execution
./mcp-proxy-fips-darwin-arm64 --version
```
