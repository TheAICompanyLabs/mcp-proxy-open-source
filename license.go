// license.go
package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
	"os"
	"time"
)

type VerifyRequest struct {
	LicenseKey string `json:"licenseKey"`
}

type VerifyResponse struct {
	Valid    bool     `json:"valid"`
	Status   string   `json:"status"`
	Tier     string   `json:"tier"`
	Features []string `json:"features"`
	Error    string   `json:"error,omitempty"`
}

type EnterpriseContext struct {
	IsEnterprise bool
	LicenseKey   string
	Features     map[string]bool
}

var EnterpriseState = &EnterpriseContext{
	IsEnterprise: false,
	Features:     make(map[string]bool),
}

func VerifyLicense(apiURL string, key string) (*VerifyResponse, error) {
	if key == "" {
		return nil, fmt.Errorf("no license key provided")
	}

	payload, err := json.Marshal(VerifyRequest{LicenseKey: key})
	if err != nil {
		return nil, err
	}

	client := &http.Client{Timeout: 5 * time.Second}
	resp, err := client.Post(apiURL+"/api/license/verify", "application/json", bytes.NewBuffer(payload))
	if err != nil {
		return nil, fmt.Errorf("verification network error: %w", err)
	}
	defer resp.Body.Close()

	var result VerifyResponse
	if err := json.NewDecoder(resp.Body).Decode(&result); err != nil {
		return nil, fmt.Errorf("failed to decode verification response: %w", err)
	}

	return &result, nil
}

func InitLicense(logger interface{ Printf(string, ...interface{}) }) {
	key := os.Getenv("MCP_ENTERPRISE_KEY")
	apiURL := os.Getenv("MCP_API_URL")

	if apiURL == "" {
		apiURL = "https://mcp-proxy-pink.vercel.app"
	}

	if key == "" {
		logger.Printf("[LICENSE] Running in Free / Community Open Core mode (air-gapped).\n")
		return
	}

	resp, err := VerifyLicense(apiURL, key)
	if err != nil || !resp.Valid {
		logger.Printf("[LICENSE] Enterprise validation failed: %v. Falling back to Free Tier.\n", err)
		return
	}

	EnterpriseState.IsEnterprise = true
	EnterpriseState.LicenseKey = key
	for _, f := range resp.Features {
		EnterpriseState.Features[f] = true
	}

	logger.Printf("[LICENSE] Enterprise Tier unlocked for key: %s...\n", key[:12])
}
