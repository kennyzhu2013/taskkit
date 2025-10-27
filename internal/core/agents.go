package core

import (
	"os"
	"os/exec"
	"path/filepath"
)

// AgentInfo captures unified metadata for AI assistants as per design
// Name: display name; Folder: local dir convention; InstallURL: official install doc; RequiresCLI: whether a CLI must be present
// ExecName: optional CLI executable name override (defaults to key when RequiresCLI)
type AgentInfo struct {
	Name        string
	Folder      string
	InstallURL  string
	RequiresCLI bool
	ExecName    string
}

// AgentConfig mirrors the design document with 13 assistants supported
var AgentConfig = map[string]AgentInfo{
	"copilot": {
		Name:        "GitHub Copilot",
		Folder:      ".github/",
		RequiresCLI: false,
	},
	"claude": {
		Name:        "Claude Code",
		Folder:      ".claude/",
		InstallURL:  "https://docs.anthropic.com/en/docs/claude-code/setup",
		RequiresCLI: true,
		ExecName:    "claude",
	},
	"gemini": {
		Name:        "Gemini CLI",
		Folder:      ".gemini/",
		InstallURL:  "https://github.com/google-gemini/gemini-cli",
		RequiresCLI: true,
		ExecName:    "gemini",
	},
	"cursor-agent": {
		Name:        "Cursor",
		Folder:      ".cursor/",
		RequiresCLI: false,
	},
	"qwen": {
		Name:        "Qwen Code",
		Folder:      ".qwen/",
		InstallURL:  "https://github.com/QwenLM/qwen-code",
		RequiresCLI: true,
		ExecName:    "qwen",
	},
	"opencode": {
		Name:        "opencode",
		Folder:      ".opencode/",
		InstallURL:  "https://opencode.ai",
		RequiresCLI: true,
		ExecName:    "opencode",
	},
	"codex": {
		Name:        "Codex CLI",
		Folder:      ".codex/",
		InstallURL:  "https://github.com/openai/codex",
		RequiresCLI: true,
		ExecName:    "codex",
	},
	"codebuddy": {
		Name:        "CodeBuddy",
		Folder:      ".codebuddy/",
		InstallURL:  "https://www.codebuddy.ai",
		RequiresCLI: true,
		ExecName:    "codebuddy",
	},
}

// Presentation order for selection UI
var orderedAgents = []string{
	"copilot", "claude", "gemini", "cursor-agent", "qwen", "opencode", "codex", "codebuddy",
}

// SupportedAgents returns ordered agent keys for selection
func SupportedAgents() []string {
	ret := make([]string, 0, len(orderedAgents))
	for _, k := range orderedAgents {
		if _, ok := AgentConfig[k]; ok {
			ret = append(ret, k)
		}
	}
	return ret
}

// GetAgentInfo returns configuration by key
func GetAgentInfo(key string) (AgentInfo, bool) {
	info, ok := AgentConfig[key]
	return info, ok
}

// GetAgentExecName resolves the CLI executable name for detection
func GetAgentExecName(key string) string {
	info, ok := AgentConfig[key]
	if !ok {
		return key
	}
	if info.ExecName != "" {
		return info.ExecName
	}
	return key
}

// DetectAgentTools scans PATH (and known locations) and returns installed AI agent tool names
func DetectAgentTools() []string {
	installed := make([]string, 0, len(AgentConfig))
	for key, info := range AgentConfig {
		if !info.RequiresCLI {
			continue
		}
		// Special-case Claude local migration path
		if key == "claude" {
			if hasClaudeLocal() {
				installed = append(installed, key)
				continue
			}
		}
		execName := GetAgentExecName(key)
		if _, err := exec.LookPath(execName); err == nil {
			installed = append(installed, key)
		}
	}
	return installed
}

// IsAgentInstalled reports whether the given agent CLI is available.
// Special-cases Claude's migrated local path besides PATH lookup.
func IsAgentInstalled(key string) bool {
	if key == "claude" {
		if hasClaudeLocal() {
			return true
		}
	}
	execName := GetAgentExecName(key)
	_, err := exec.LookPath(execName)
	return err == nil
}

// hasClaudeLocal detects migrated Claude CLI at ~/.claude/local/claude (and .exe on Windows)
func hasClaudeLocal() bool {
	home, err := os.UserHomeDir()
	if err != nil || home == "" {
		return false
	}
	p1 := filepath.Join(home, ".claude", "local", "claude")
	if st, err := os.Stat(p1); err == nil && !st.IsDir() {
		return true
	}
	p2 := filepath.Join(home, ".claude", "local", "claude.exe")
	if st, err := os.Stat(p2); err == nil && !st.IsDir() {
		return true
	}
	return false
}
