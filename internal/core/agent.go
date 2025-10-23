package core

import (
	"os"
	"os/exec"
	"path/filepath"
)

// DetectAgentTools scans PATH (and known locations) and returns installed AI agent tool names
func DetectAgentTools() []string {
	candidates := []string{"claude", "qwen", "gemini", "codebuddy", "q"}
	installed := make([]string, 0, len(candidates))

	for _, name := range candidates {
		if name == "claude" {
			if hasClaudeLocal() {
				installed = append(installed, name)
				continue
			}
		}
		if _, err := exec.LookPath(name); err == nil {
			installed = append(installed, name)
		}
	}
	return installed
}

// IsAgentInstalled reports whether the given agent CLI is available.
// Special-cases Claude's migrated local path besides PATH lookup.
func IsAgentInstalled(name string) bool {
	if name == "claude" {
		if hasClaudeLocal() {
			return true
		}
	}
	_, err := exec.LookPath(name)
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