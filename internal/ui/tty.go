package ui

import (
	"os"
	"strings"

	"golang.org/x/term"
)

// disableTUI 通过环境变量控制是否禁用 TUI，支持：1/true/yes（不区分大小写）
func disableTUI() bool {
	v := strings.TrimSpace(strings.ToLower(os.Getenv("TASKKIT_NO_TUI")))
	return v == "1" || v == "true" || v == "yes"
}

// isInteractiveTerminal 同时要求 stdin 与 stdout 均为 TTY
func isInteractiveTerminal() bool {
	inTTY := term.IsTerminal(int(os.Stdin.Fd()))
	outTTY := term.IsTerminal(int(os.Stdout.Fd()))
	return inTTY && outTTY
}