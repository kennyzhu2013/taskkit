package cmd

import (
	"fmt"
	"os/exec"
	"strings"

	"github.com/spec-kit/task-kit/internal/core"
	"github.com/spf13/cobra"
)

var checkCmd = &cobra.Command{
	Use:   "check",
	Short: "检查本机工具链是否就绪（Go 草案）",
	RunE: func(cmd *cobra.Command, args []string) error {
		bins := []string{"git", "curl", "unzip", "tar", "code", "code-insiders"}
		for _, b := range bins {
			if _, err := exec.LookPath(b); err != nil {
				fmt.Printf("[x] 未找到: %s\n", b)
			} else {
				fmt.Printf("[✓] 已安装: %s\n", b)
			}
		}

		// AI Agent 工具检测：基于 DetectAgentTools 结果
		agents := core.DetectAgentTools()
		if len(agents) == 0 {
			fmt.Println("[i] 未检测到 AI Agent 工具")
		} else {
			fmt.Printf("[✓] 检测到 AI Agent 工具: %s\n", strings.Join(agents, ", "))
		}
		return nil
	},
}