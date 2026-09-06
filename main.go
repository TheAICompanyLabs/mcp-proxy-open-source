package main

import (
	"bufio"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"os"
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
	// 1. Run silently in the background
	go func() {
		home, _ := os.UserHomeDir()
		cacheFile := filepath.Join(home, ".mcp-proxy", "last_update_check")

		// 2. Only check once every 24 hours
		if info, err := os.Stat(cacheFile); err == nil {
			if time.Since(info.ModTime()) < 24*time.Hour {
				return
			}
		}

		// 3. Ping GitHub Releases API
		resp, err := http.Get("https://api.github.com/repos/TheAICompanyLabs/mcp-proxy-open-source/releases/latest")
		if err != nil || resp.StatusCode != 200 {
			return
		}
		defer resp.Body.Close()

		var release struct {
			TagName string `json:"tag_name"`
		}
		if err := json.NewDecoder(resp.Body).Decode(&release); err == nil {
			// 4. Update the cache file timestamp
			os.MkdirAll(filepath.Dir(cacheFile), 0755)
			os.WriteFile(cacheFile, []byte(release.TagName), 0644)

			// 5. Notify if the version is newer
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
	// 1. Check environment variable (highest priority for CI/CD & automation)
	if envKey := strings.TrimSpace(os.Getenv("MCP_LICENSE_KEY")); envKey != "" {
		return envKey
	}

	// 2. Check local saved configuration
	cfg, err := LoadConfig()
	if err != nil {
		log.Printf("⚠️ Warning: Could not read config file: %v", err)
	} else if cfg != nil && strings.TrimSpace(cfg.LicenseKey) != "" {
		return strings.TrimSpace(cfg.LicenseKey)
	}

	return ""
}

func main() {
	// 1. Check if user is invoking the login subcommand
	if len(os.Args) > 1 && os.Args[1] == "login" {
		handleLogin()
		return
	}

	// 2. Resolve the active license key (Free vs. Enterprise)
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

	// 3. Verify target MCP server command is passed
	if len(os.Args) < 2 {
		fmt.Println("\nError: No target MCP server command provided.")
		fmt.Println("Usage: go run . <your_server_command>")
		os.Exit(1)
	}

	// 4. Initialize arguments and package-scoped logger
	args := os.Args[1:]
	serverCmdString := strings.Join(args, " ")

	logFile, err := os.OpenFile("mcp_audit.log", os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0644)
	if err != nil {
		log.Fatalf("Failed to open log file: %v", err)
	}
	defer logFile.Close()

	// Assign to the package-level appLogger variable
	appLogger = log.New(logFile, "[MCP-PROXY] ", log.LstdFlags)
	loadPolicy("policy.yaml", appLogger)
	InitLicense(appLogger)
	pendingDecisionChan = make(chan bool)

	// 1. Detect if the process is running headless
	fileInfo, _ := os.Stdout.Stat()
	isHeadless := (fileInfo.Mode() & os.ModeCharDevice) == 0

	var prog *tea.Program
	var tty *os.File
	var ttyErr error

	// 2. Only boot the Terminal UI if a TTY is attached
	if !isHeadless {
		ttyPath := "/dev/tty"
		if runtime.GOOS == "windows" {
			ttyPath = "CONIN$" // Native Windows console input
		}

		tty, ttyErr = os.OpenFile(ttyPath, os.O_RDWR, 0)
		if ttyErr == nil {
			prog = tea.NewProgram(initialUIModel(serverCmdString), tea.WithInput(tty), tea.WithOutput(tty))
		} else {
			appLogger.Printf("Failed to open keyboard TTY: %v\n", ttyErr)
		}
	} else {
		appLogger.Println("[SYSTEM] Running in Headless Mode. Interactive TUI disabled.")
	}

	// 3. Pass the `isHeadless` flag into the pipeline
	cmd, err := StartProxyPipeline(args, prog, appLogger, isHeadless)
	if err != nil {
		appLogger.Fatalf("Pipeline startup failed: %v", err)
	}

	// 4. Handle Execution Halting
	if prog != nil {
		// If TUI is active, run it. It blocks until the user quits.
		if _, err := prog.Run(); err != nil {
			appLogger.Printf("TUI runtime error: %v", err)
		}
		if tty != nil {
			tty.Close()
		}
	} else {
		// If Headless, block on the underlying server process
		cmd.Wait()
	}
}
