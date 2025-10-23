package core

import (
	"runtime"

	"github.com/spec-kit/task-kit/internal/ui"
)

func DetectDefaultScript() string {
	if runtime.GOOS == "windows" {
		return "powershell"
	}
	return "bash"
}

func SelectScript(defaultScript string) (string, error) {
	options := []string{"powershell", "bash", "zsh"}
	// ensure default is first (用于非交互回退)
	start := 0
	for i, s := range options {
		if s == defaultScript {
			start = i
			break
		}
	}
	if start != 0 {
		options[0], options[start] = options[start], options[0]
	}
	return ui.SelectFromList(options, "选择脚本类型", defaultScript)
}