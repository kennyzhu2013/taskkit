package main

import (
	"fmt"

	"github.com/spec-kit/task-kit/internal/core"
)

func main() {
	def := core.DetectDefaultScript()
	fmt.Println("[TEST] 开始 脚本类型 选择 (TUI 若不可用将回退 stdin)")
	s, err := core.SelectScript(def)
	if err != nil {
		fmt.Println("[TEST] ERR:", err)
		return
	}
	fmt.Println("[TEST] SCRIPT:", s)
}