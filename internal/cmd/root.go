package cmd

import (
	"fmt"
	"os"

	"github.com/spf13/cobra"
)

var (
	debug            bool
	githubToken      string
	skipTLSVerify    bool
	forceInteractive bool
)

// rootCmd is the base command for task-kit
var rootCmd = &cobra.Command{
	Use:   "task-kit",
	Short: "Task Kit - initialize and scaffold from templates",
	PersistentPreRun: func(cmd *cobra.Command, args []string) {
		if debug {
			fmt.Fprintln(os.Stderr, "[debug] debug mode enabled")
		}
	},
}

// Execute runs the root command
func Execute() {
	// persistent flags on root
	rootCmd.PersistentFlags().BoolVar(&debug, "debug", false, "启用调试输出")
	rootCmd.PersistentFlags().BoolVar(&skipTLSVerify, "skip-tls-verify", false, "跳过 TLS 证书校验（仅用于受信环境）")
	rootCmd.PersistentFlags().StringVar(&githubToken, "github-token", "", "GitHub 令牌；为空时依次回退 GH_TOKEN, GITHUB_TOKEN")
	rootCmd.PersistentFlags().BoolVar(&forceInteractive, "interactive", false, "强制启用交互式选择（在某些终端下 TUI 判定失败时使用）")

	// add subcommands
	rootCmd.AddCommand(initCmd)
	rootCmd.AddCommand(checkCmd)

	if err := rootCmd.Execute(); err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
}
