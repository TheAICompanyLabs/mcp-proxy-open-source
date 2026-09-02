package main

import (
	"log"
	"os"
	"regexp"
	"sync"

	"go.yaml.in/yaml/v4"
)

type PolicyConfig struct {
	Version  string   `yaml:"version"`
	Policies []Policy `yaml:"policies"`
}

type Policy struct {
	Tool       string      `yaml:"tool"`
	Action     string      `yaml:"action"`
	Message    string      `yaml:"message"`
	Conditions []Condition `yaml:"conditions,omitempty"`
}

type Condition struct {
	Argument string         `yaml:"argument"`
	Operator string         `yaml:"operator"`
	Value    string         `yaml:"value"`
	Compiled *regexp.Regexp `yaml:"-"`
}

var (
	activePolicies map[string]Policy
	policyMu       sync.RWMutex
)

func initPolicies() {
	policyMu.Lock()
	defer policyMu.Unlock()
	if activePolicies == nil {
		activePolicies = make(map[string]Policy)
	}
}

func loadPolicy(filename string, logger *log.Logger) {
	initPolicies()

	data, err := os.ReadFile(filename)
	if err != nil {
		logger.Printf("Policy file %s not found. Initializing Zero-Trust fallback.\n", filename)
		policyMu.Lock()
		activePolicies["*"] = Policy{
			Tool:    "*",
			Action:  "REQUIRE_APPROVAL",
			Message: "Unrecognized tool invoked. Defaulting to Human-in-the-Loop review.",
		}
		policyMu.Unlock()
		savePolicyFile(filename, logger)
		return
	}

	var config PolicyConfig
	if err := yaml.Unmarshal(data, &config); err != nil {
		logger.Printf("Failed to parse YAML: %v\n", err)
		return
	}

	policyMu.Lock()
	defer policyMu.Unlock()
	activePolicies = make(map[string]Policy)

	for _, p := range config.Policies {
		for i, cond := range p.Conditions {
			if cond.Operator == "regex" {
				compiled, err := regexp.Compile(cond.Value)
				if err != nil {
					logger.Printf("Warning: Invalid regex '%s' on %s: %v\n", cond.Value, p.Tool, err)
					continue
				}
				p.Conditions[i].Compiled = compiled
			}
		}
		activePolicies[p.Tool] = p
	}
	logger.Printf("Loaded %d policies from %s\n", len(activePolicies), filename)
}

func savePolicyFile(filename string, logger *log.Logger) {
	policyMu.RLock()
	var config PolicyConfig
	config.Version = "1.2"
	config.Policies = make([]Policy, 0, len(activePolicies))

	for _, p := range activePolicies {
		config.Policies = append(config.Policies, p)
	}
	policyMu.RUnlock()

	data, err := yaml.Marshal(&config)
	if err != nil {
		logger.Printf("Failed to marshal YAML: %v\n", err)
		return
	}

	if err := os.WriteFile(filename, data, 0644); err != nil {
		logger.Printf("Failed to write %s: %v\n", filename, err)
		return
	}
	logger.Printf("Persisted %d policies to %s\n", len(config.Policies), filename)
}
