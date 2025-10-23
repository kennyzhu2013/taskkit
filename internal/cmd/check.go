package cmd

import (
	"fmt"
	"os/exec"

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
		return nil
	},
}