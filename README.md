# Task Kit 项目说明
**语言与技术栈**: Go CLI (Cobra + Bubbletea + Resty)

## 概览
  - `templates/` 规范、计划、任务与命令模板，驱动“/specify、/plan、/tasks”等流程。
  - `scripts/` Bash/PowerShell 脚本，与命令模板配套，生成分支与规范文件、更新 Agent 上下文。
  - `specs/` 样例规范目录结构（如 `specs/main/`）。

## 仓库结构
- `task-kit/`：Go 项目根目录，包含 CLI 入口、内部模块、脚本与模板。
- `templates/`：跨工具的规范、计划、任务与命令模板。
- `scripts/`：与 `templates/commands/*.md` 中的脚本占位相配的实现。
- `specs/`：功能级文档输出目录（feature 级 plan/spec 等）。

## CLI 入口与命令
- 入口文件：`task-kit/cmd/task-kit/main.go`
  - 调用 `internal/cmd.Execute()` 启动 CLI。
- 根命令：`task-kit/internal/cmd/root.go`
  - Flags（持久化）：`--debug`、`--skip-tls-verify`、`--github-token`、`--interactive`。
  - 子命令：
    - `init`：项目初始化（模板下载、解压、复制、权限修复、Git 初始化、后续提示）。
    - `check`：本机工具链自检（标准工具与 AI Agent CLI 检测）。

## 依赖与环境
- Go：`spf13/cobra`（CLI）、`charmbracelet/bubbletea`（TUI）、`go-resty/resty`（HTTP）。
- 可选：`--skip-tls-verify` 仅用于受信本地调试环境；生产需启用 TLS 校验。

## 使用与下一步
- 初始化项目：
  - `task-kit init <PROJECT_NAME>`
  - 指定脚本：`task-kit init <PROJECT_NAME> --script ps|sh`（Windows 推荐 `ps`）。
  - 指定来源：`--zip-url <URL>` 或 `--repo owner/repo --branch <ref>`；未指定时默认 `--latest-release`。
- 进入项目目录后，按面板提示使用斜杠命令：
  - `/constitution` → 建立项目原则
  - `/specify` → 创建基线规范（调用 `create-new-feature.*`）
  - `/plan` → 制定实施计划（生成 `plan.md` 等）
  - `/tasks` → 生成可执行任务清单
  - `/implement` → 执行实现与集成
