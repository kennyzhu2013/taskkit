package cmd

import (
	"fmt"
	"os"
	"os/exec"
	"path/filepath"

	"github.com/spec-kit/task-kit/internal/core"
	"github.com/spec-kit/task-kit/internal/ui"
	"github.com/spec-kit/task-kit/internal/util"
	"github.com/spf13/cobra"
)

type InitOptions struct {
	Repo             string // owner/repo 或 https://github.com/owner/repo
	ZipURL           string // 直接 ZIP 下载链接（优先于 Repo）
	Branch           string // 分支或 tag，默认 main（仅对 Repo 生效）
	Folder           string // 仓库内子目录，为空则使用仓库根
	TargetDir        string // 目标目录，默认为当前目录
	LatestRelease    bool   // 若为 true，则根据 AI+Script 匹配最新 release 资源
	AI               string // AI 助手名，空则交互选择
	IgnoreAgentTools bool   // 忽略 AI Agent 工具自动检测
	Script           string // 脚本类型（与 Python --script 对齐）
	KeepTemp         bool   // 保留 .task-kit-tmp 目录用于调试
	Debug            bool
	SkipTLSVerify    bool
	GithubToken      string
	// 新增：与 Python CLI 对齐的默认行为
	NoGit bool // 跳过 git 初始化
	Here  bool // 在当前目录初始化
	Force bool // 目标目录非空时强制继续
}

var initOpts InitOptions

var initCmd = &cobra.Command{
	Use:   "init",
	Short: "初始化项目：从 GitHub 下载模板并展开",
	RunE: func(cmd *cobra.Command, args []string) error {
		// 填充全局选项
		initOpts.Debug = debug
		initOpts.SkipTLSVerify = skipTLSVerify
		initOpts.GithubToken = core.ResolveGitHubToken(githubToken)

		// 处理 --here 与位置参数冲突
		if initOpts.Here {
			if len(args) > 0 && args[0] != "." {
				return fmt.Errorf("不能同时指定位置参数与 --here；若要在当前目录初始化，请使用 --here 或指定 '.'")
			}
			cwd, _ := os.Getwd()
			initOpts.TargetDir = cwd
		}

		// 当未显式提供来源参数时，默认使用最新 Release（与 Python CLI 对齐）
		if initOpts.ZipURL == "" && initOpts.Repo == "" && !initOpts.LatestRelease {
			initOpts.LatestRelease = true
			initOpts.Repo = "kennyzhu2013/taskkit" // initOpts.Repo = "github/spec-kit"
			if initOpts.Debug {
				util.Debugf("默认使用最新 Release: repo=%s\n", initOpts.Repo)
			}
		}
		// latest-release 默认仓库（保持原有兜底逻辑）
		if initOpts.LatestRelease && initOpts.Repo == "" {
			initOpts.Repo = "kennyzhu2013/taskkit" // initOpts.Repo = "github/spec-kit"
			if initOpts.Debug {
				util.Debugf("latest-release 默认 repo=%s\n", initOpts.Repo)
			}
		}

		// 将第一个位置参数作为目标目录（支持 `init test` 或 `init .`）
		if initOpts.TargetDir == "" && len(args) > 0 {
			if args[0] == "." {
				cwd, _ := os.Getwd()
				initOpts.TargetDir = cwd
			} else {
				initOpts.TargetDir = args[0]
			}
		}

		if initOpts.TargetDir == "" {
			cwd, _ := os.Getwd()
			initOpts.TargetDir = cwd
		}

		// 目标目录存在性与非空检查（与 Python 语义一致，需要 --force）
		if st, err := os.Stat(initOpts.TargetDir); err == nil && st.IsDir() {
			f, err := os.Open(initOpts.TargetDir)
			if err != nil {
				return err
			}
			names, _ := f.Readdirnames(1)
			_ = f.Close()
			if len(names) > 0 && !initOpts.Force {
				return fmt.Errorf("目标目录已存在且非空：%s。使用 --force 覆盖，或选择其他目录", initOpts.TargetDir)
			}
		}

		if initOpts.Branch == "" {
			initOpts.Branch = "main"
		}

		if initOpts.AI == "" {
			options := core.SupportedAgents()
			ai, err := core.SelectAI(options)
			if err != nil {
				return fmt.Errorf("需要交互选择 AI，请使用 --ai 指定或在交互终端选择: %w", err)
			}
			initOpts.AI = ai
		}
		if initOpts.Debug {
			util.Debugf("init with repo=%s branch=%s folder=%s dest=%s ai=%s skipTLS=%v\n", initOpts.Repo, initOpts.Branch, initOpts.Folder, initOpts.TargetDir, initOpts.AI, initOpts.SkipTLSVerify)
		}

		// 脚本类型：若未指定则按平台默认并可交互选择
		if initOpts.Script == "" {
			def := core.DetectDefaultScript()
			s, err := core.SelectScript(def)
			if err != nil {
				return fmt.Errorf("需要交互选择脚本类型，请使用 --script 指定或在交互终端选择: %w", err)
			}
			initOpts.Script = s
		}

		// 步骤 TUI（仅交互终端启用）
		var steps *ui.StepUI
		if true {
			var items []ui.StepItem
			if initOpts.LatestRelease {
				items = append(items, ui.StepItem{Key: "release", Label: "获取 Release 资产", Status: ui.StatusPending})
			}
			items = append(items,
				ui.StepItem{Key: "download", Label: "下载模板压缩包", Status: ui.StatusPending},
				ui.StepItem{Key: "extract", Label: "解压模板", Status: ui.StatusPending},
				ui.StepItem{Key: "flatten", Label: "扁平化目录", Status: ui.StatusPending},
				ui.StepItem{Key: "copy", Label: "复制到目标目录", Status: ui.StatusPending},
				ui.StepItem{Key: "perm", Label: "修复脚本权限", Status: ui.StatusPending},
				ui.StepItem{Key: "cleanup", Label: "清理临时文件", Status: ui.StatusPending},
				ui.StepItem{Key: "git", Label: "初始化 Git 仓库", Status: ui.StatusPending},
			)
			var err error
			steps, err = ui.StartStepUI("初始化模板", items)
			if err != nil {
				steps = nil
			}
		}
		start := func(k, d string) {
			if steps != nil {
				steps.Start(k, d)
			}
		}
		complete := func(k, d string) {
			if steps != nil {
				steps.Complete(k, d)
			}
		}
		errorStep := func(k, d string) {
			if steps != nil {
				steps.Error(k, d)
			}
		}
		skip := func(k, d string) {
			if steps != nil {
				steps.Skip(k, d)
			}
		}
		progressDL := func(done, total int64) {
			if steps != nil {
				steps.ProgressBytes("download", done, total)
			}
		}
		progressExtract := func(done, total int64) {
			if steps != nil {
				steps.ProgressBytes("extract", done, total)
			}
		}
		progressCopy := func(done, total int64) {
			if steps != nil {
				steps.ProgressBytes("copy", done, total)
			}
		}
		defer func() {
			if steps != nil {
				steps.Stop()
			}
		}()

		// 1) 下载归档：优先 ZipURL；其次 latest-release；否则仓库 zipball
		var archivePath string
		var err error
		if initOpts.ZipURL != "" {
			start("download", "直接链接")
			archivePath, err = core.DownloadZipFromURLWithProgress(initOpts.ZipURL, initOpts.GithubToken, initOpts.SkipTLSVerify, initOpts.Debug, progressDL)
			if err != nil {
				errorStep("download", err.Error())
				return err
			}
			complete("download", "完成")
		} else if initOpts.LatestRelease {
			start("release", initOpts.Repo)
			_, assetURL, selErr := core.GetLatestReleaseAssetURL(initOpts.Repo, initOpts.AI, initOpts.Script, initOpts.GithubToken, initOpts.SkipTLSVerify, initOpts.Debug)
			if selErr != nil {
				errorStep("release", selErr.Error())
				return selErr
			}
			complete("release", "已选择资产")
			start("download", "最新 Release")
			archivePath, err = core.DownloadZipFromURLWithProgress(assetURL, initOpts.GithubToken, initOpts.SkipTLSVerify, initOpts.Debug, progressDL)
			if err != nil {
				errorStep("download", err.Error())
				return err
			}
			complete("download", "完成")
		} else {
			start("download", fmt.Sprintf("%s@%s", initOpts.Repo, initOpts.Branch))
			archivePath, err = core.DownloadTemplateFromGitHubWithProgress(core.DownloadOptions{
				Repo:          initOpts.Repo,
				Branch:        initOpts.Branch,
				Token:         initOpts.GithubToken,
				SkipTLSVerify: initOpts.SkipTLSVerify,
				Debug:         initOpts.Debug,
			}, progressDL)
			if err != nil {
				errorStep("download", err.Error())
				return err
			}
			complete("download", "完成")
		}
		defer os.Remove(archivePath)

		// 2) 解压
		start("extract", ".task-kit-tmp")
		extractDir := filepath.Join(initOpts.TargetDir, ".task-kit-tmp")
		if err := os.MkdirAll(extractDir, 0o755); err != nil {
			errorStep("extract", err.Error())
			return err
		}
		if err := core.ExtractZipWithProgress(archivePath, extractDir, progressExtract); err != nil {
			errorStep("extract", err.Error())
			return err
		}
		complete("extract", "完成")

		// 3) 如果只有单个顶层目录则扁平化
		start("flatten", "单目录检测")
		actualRoot, err := core.FlattenSingleTopDir(extractDir)
		if err != nil {
			errorStep("flatten", err.Error())
			return err
		}
		complete("flatten", "完成")

		// 4) 如果指定子目录，仅迁移该子目录内容
		moveRoot := actualRoot
		if initOpts.Folder != "" {
			moveRoot = filepath.Join(actualRoot, filepath.FromSlash(initOpts.Folder))
			if _, err := os.Stat(moveRoot); err != nil {
				err := fmt.Errorf("仓库内子目录不存在: %s", initOpts.Folder)
				errorStep("copy", err.Error())
				return err
			}
		}

		start("copy", initOpts.TargetDir)
		if err := core.CopyTreeWithProgress(moveRoot, initOpts.TargetDir, progressCopy); err != nil {
			errorStep("copy", err.Error())
			return err
		}
		complete("copy", "完成")

		// 5) 权限修复（主要针对 *nix 脚本）
		start("perm", "修复可执行权限")
		_ = core.FixScriptPermissions(initOpts.TargetDir)
		complete("perm", "完成")

		// 6) 清理临时目录（除非用户要求保留）
		start("cleanup", "清理中间文件")
		if !initOpts.KeepTemp {
			_ = os.RemoveAll(extractDir)
			if initOpts.Debug {
				util.Debugf("临时目录已清理: %s\n", extractDir)
			}
		}
		complete("cleanup", "完成")

		// 7) Git 初始化（可跳过；失败不致命）
		start("git", "检测 git 环境")
		if initOpts.NoGit {
			skip("git", "--no-git 已设置，跳过")
		} else if _, err := exec.LookPath("git"); err != nil {
			skip("git", "未找到 git，已跳过")
		} else {
			// 是否已经处于 git 仓库中
			cmdCheck := exec.Command("git", "rev-parse", "--is-inside-work-tree")
			cmdCheck.Dir = initOpts.TargetDir
			if err := cmdCheck.Run(); err == nil {
				skip("git", "目标目录已是 git 仓库")
			} else {
				cmdInit := exec.Command("git", "init")
				cmdInit.Dir = initOpts.TargetDir
				if err := cmdInit.Run(); err != nil {
					skip("git", "git init 失败，已跳过")
				} else {
					cmdAdd := exec.Command("git", "add", ".")
					cmdAdd.Dir = initOpts.TargetDir
					_ = cmdAdd.Run()
					cmdCommit := exec.Command("git", "-c", "user.name=TaskKit", "-c", "user.email=taskkit@local", "commit", "-m", "chore: initialize with task-kit")
					cmdCommit.Dir = initOpts.TargetDir
					if err := cmdCommit.Run(); err != nil {
						skip("git", "git commit 失败，已跳过")
					} else {
						complete("git", "完成")
					}
				}
			}
		}

		fmt.Printf("模板已初始化到: %s\n", initOpts.TargetDir)
		if initOpts.AI != "" {
			fmt.Printf("已选择 AI: %s\n", initOpts.AI)
		}
		// 非强制校验模式：若选择了需 CLI 的 AI 且未检测到，给出提示（从 AgentConfig 获取）
		if !initOpts.IgnoreAgentTools {
			if info, ok := core.GetAgentInfo(initOpts.AI); ok && info.RequiresCLI && !core.IsAgentInstalled(initOpts.AI) {
				if info.InstallURL != "" {
					fmt.Fprintf(os.Stderr, "[warn] 未检测到 %s CLI。安装: %s\n", initOpts.AI, info.InstallURL)
				} else {
					fmt.Fprintf(os.Stderr, "[warn] 未检测到 %s CLI。\n", initOpts.AI)
				}
				fmt.Fprintln(os.Stderr, "提示：安装后再试，或使用 --ignore-agent-tools 跳过检查。")
			}
		}
		return nil
	},
}

func init() {
	initCmd.Flags().StringVar(&initOpts.Repo, "repo", "", "GitHub 仓库（owner/repo 或 https://github.com/owner/repo）")
	initCmd.Flags().StringVar(&initOpts.ZipURL, "zip-url", "", "直接 ZIP 下载链接（优先于 --repo）")
	initCmd.Flags().StringVar(&initOpts.Branch, "branch", "", "分支或 tag（默认为 main；在 --zip-url 或 --latest-release 下忽略）")
	initCmd.Flags().StringVar(&initOpts.Folder, "folder", "", "仓库内子目录，仅迁移该目录内容")
	initCmd.Flags().StringVar(&initOpts.TargetDir, "target-dir", "", "目标目录，默认当前目录")
	initCmd.Flags().StringVar(&initOpts.AI, "ai", "", "AI 助手名称（为空时交互选择）")
	initCmd.Flags().BoolVar(&initOpts.IgnoreAgentTools, "ignore-agent-tools", false, "忽略 AI Agent 工具自动检测")
	initCmd.Flags().StringVar(&initOpts.Script, "script", "", "脚本类型（与 Python --script 对齐）")
	initCmd.Flags().BoolVar(&initOpts.LatestRelease, "latest-release", false, "从最新 Release 中选择匹配 AI+Script 的模板资产（--repo 为空时默认 kennyzhu2013/taskkit）")
	initCmd.Flags().BoolVar(&initOpts.KeepTemp, "keep-temp", false, "保留中间 .task-kit-tmp 目录（调试用）")
	// 新增 flags
	initCmd.Flags().BoolVar(&initOpts.NoGit, "no-git", false, "跳过 Git 仓库初始化")
	initCmd.Flags().BoolVar(&initOpts.Here, "here", false, "在当前目录初始化项目")
	initCmd.Flags().BoolVar(&initOpts.Force, "force", false, "目标目录非空时强制继续")
}
