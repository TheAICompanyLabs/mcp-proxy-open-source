package main

import (
	"bufio"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"os"
	"os/exec"
	"path/filepath"
	"runtime"
	"strings"
	"time"

	tea "github.com/charmbracelet/bubbletea"
)

// Package-level variables accessible across ui.go, proxy.go, etc.
var (
	pendingDecisionChan chan bool
	appLogger           *log.Logger
)

const currentVersion = "v1.0.1"

func checkUpdateAsync() {
	go func() {
		home, _ := os.UserHomeDir()
		cacheFile := filepath.Join(home, ".mcp-proxy", "last_update_check")

		if info, err := os.Stat(cacheFile); err == nil {
			if time.Since(info.ModTime()) < 24*time.Hour {
				return
			}
		}

		resp, err := http.Get("https://api.github.com/repos/TheAICompanyLabs/mcp-proxy-open-source/releases/latest")
		if err != nil || resp.StatusCode != 200 {
			return
		}
		defer resp.Body.Close()

		var release struct {
			TagName string `json:"tag_name"`
		}
		if err := json.NewDecoder(resp.Body).Decode(&release); err == nil {
			os.MkdirAll(filepath.Dir(cacheFile), 0755)
			os.WriteFile(cacheFile, []byte(release.TagName), 0644)

			if release.TagName != currentVersion && release.TagName != "" {
				fmt.Printf("\n\033[33m🚀 A new version of Universal MCP Proxy is available! (%s -> %s)\033[0m\n", currentVersion, release.TagName)
				fmt.Println("\033[33mRun your installation script to update.\033[0m")
			}
		}
	}()
}

func handleLogin() {
	reader := bufio.NewReader(os.Stdin)
	fmt.Print("Enter your MCP Proxy License Key: ")
	key, err := reader.ReadString('\n')
	if err != nil {
		log.Fatalf("❌ Failed to read input: %v", err)
	}

	key = strings.TrimSpace(key)
	if key == "" {
		log.Fatal("❌ License key cannot be empty.")
	}

	apiURL := os.Getenv("MCP_API_URL")
	if apiURL == "" {
		apiURL = "https://mcp-proxy-pink.vercel.app"
	}

	fmt.Println("Authenticating with licensing server...")
	resp, err := VerifyLicense(apiURL, key)
	if err != nil {
		log.Fatalf("❌ Network error connecting to licensing server: %v", err)
	}
	if !resp.Valid {
		log.Fatalf("❌ Authentication failed: Invalid or inactive license key")
	}

	cfg := &Config{LicenseKey: key}
	if err := SaveConfig(cfg); err != nil {
		log.Fatalf("❌ Failed to save credentials: %v", err)
	}

	configPath, _ := getConfigPath()
	fmt.Printf("🔒 Credentials successfully stored in %s\n", configPath)
	fmt.Println("🚀 You're all set! You can now run the proxy directly.")
	os.Exit(0)
}

func resolveLicenseKey() string {
	if envKey := strings.TrimSpace(os.Getenv("MCP_LICENSE_KEY")); envKey != "" {
		return envKey
	}

	cfg, err := LoadConfig()
	if err != nil {
		log.Printf("⚠️ Warning: Could not read config file: %v", err)
	} else if cfg != nil && strings.TrimSpace(cfg.LicenseKey) != "" {
		return strings.TrimSpace(cfg.LicenseKey)
	}

	return ""
}

// Helper to safely execute node/batch scripts on Windows
func execCommandHelper(name string, arg ...string) *exec.Cmd {
	if runtime.GOOS == "windows" {
		winArgs := append([]string{"/c", name}, arg...)
		return exec.Command("cmd", winArgs...)
	}
	return exec.Command(name, arg...)
}

func main() {
	if len(os.Args) > 1 && os.Args[1] == "login" {
		handleLogin()
		return
	}

	licenseKey := resolveLicenseKey()

	if licenseKey == "" {
		fmt.Println("🚀 Booting Universal MCP Proxy in Free Open Core Mode...")
		fmt.Println("   [Air-gapped Layer 1 & 2 Security Active. Enterprise modules locked.]")
	} else {
		apiURL := os.Getenv("MCP_API_URL")
		if apiURL == "" {
			apiURL = "https://mcp-proxy-pink.vercel.app"
		}

		fmt.Println("Verifying enterprise license with central server...")
		resp, err := VerifyLicense(apiURL, licenseKey)
		if err != nil {
			log.Fatalf("❌ Network error connecting to licensing server: %v", err)
		}
		if !resp.Valid {
			log.Fatalf("❌ Access Denied: Invalid or inactive license key")
		}
	}

	if len(os.Args) < 2 {
		fmt.Println("\nError: No target MCP server command provided.")
		fmt.Println("Usage: go run . <your_server_command>")
		os.Exit(1)
	}

	args := os.Args[1:]
	serverCmdString := strings.Join(args, " ")

	// --- PATH FIX: Resolve home directory to bypass Claude's Read-Only Sandbox ---
	homeDir, err := os.UserHomeDir()
	if err != nil {
		log.Fatalf("Failed to get home directory: %v", err)
	}

	proxyDir := filepath.Join(homeDir, ".mcp-proxy")
	if err := os.MkdirAll(proxyDir, 0755); err != nil {
		log.Fatalf("Failed to create proxy config directory: %v", err)
	}

	logPath := filepath.Join(proxyDir, "mcp_audit.log")
	logFile, err := os.OpenFile(logPath, os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0644)
	if err != nil {
		log.Fatalf("Failed to open log file: %v", err)
	}
	defer logFile.Close()

	appLogger = log.New(logFile, "[MCP-PROXY] ", log.LstdFlags)

	policyPath := filepath.Join(proxyDir, "policy.yaml")
	loadPolicy(policyPath, appLogger)
	// ----------------------------------------------------------------------------

	InitLicense(appLogger)
	pendingDecisionChan = make(chan bool)

	fileInfo, _ := os.Stdout.Stat()
	isHeadless := (fileInfo.Mode() & os.ModeCharDevice) == 0

	var prog *tea.Program

	if !isHeadless {
		// Bubble Tea automatically handles Windows/Mac terminal attachment natively
		prog = tea.NewProgram(initialUIModel(serverCmdString))
	} else {
		appLogger.Println("[SYSTEM] Running in Headless Mode. Interactive TUI disabled.")
	}

	cmd, err := StartProxyPipeline(args, prog, appLogger, isHeadless)
	if err != nil {
		appLogger.Fatalf("Pipeline startup failed: %v", err)
	}

	if prog != nil {
		if _, err := prog.Run(); err != nil {
			appLogger.Printf("TUI runtime error: %v", err)
		}
	} else {
		cmd.Wait()
	}
}
