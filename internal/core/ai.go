package core

import (
    "fmt"

    "github.com/spec-kit/task-kit/internal/ui"
)

// SelectAI prompts user to select an AI when none provided
func SelectAI(options []string) (string, error) {
    if len(options) == 0 {
        return "", fmt.Errorf("没有可选 AI 选项")
    }
    // 默认优先项：copilot 若存在
    defaultVal := ""
    for _, o := range options {
        if o == "copilot" {
            defaultVal = o
            break
        }
    }
    // 构建显示名映射：key -> AgentInfo.Name（若存在）
    display := make(map[string]string, len(options))
    for _, k := range options {
        if info, ok := GetAgentInfo(k); ok && info.Name != "" {
            display[k] = info.Name
        } else {
            display[k] = k
        }
    }
    return ui.SelectFromListWithDisplay(options, "选择 AI 助手", defaultVal, display)
}
