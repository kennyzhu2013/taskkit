package main

import (
	"fmt"

	"github.com/spec-kit/task-kit/internal/core"
)

func main() {
	options := []string{"copilot", "claude", "qwen", "gemini", "codebuddy", "q"}
	fmt.Println("[TEST] 开始 AI 选择 (TUI 若不可用将回退 stdin)")
	ai, err := core.SelectAI(options)
	if err != nil {
		fmt.Println("[TEST] ERR:", err)
		return
	}
	fmt.Println("[TEST] AI:", ai)
}
