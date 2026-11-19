.PHONY: bench build verify

bench:
	@go test -v -count=1 -run TestReproducibleMCPBench .
	@echo ""
	@echo "Running standard Go allocation & throughput benchmark..."
	@go test -bench=BenchmarkInspectJSONRPC -benchmem -run=^$$ .

build:
	@CGO_ENABLED=0 go build -trimpath -ldflags="-s -w" -o rcf-proxy .

verify:
	@shasum -a 256 --ignore-missing -c checksums.txt
