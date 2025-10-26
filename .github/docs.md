## 

## 上游 github/spec-kit 仓库的发布工作流或构建脚本对应的目录 .github 是如何执行生成模板ZIP的？

**整体流程**
- 触发条件：`push` 到 `main` 且改动 `memory/**`、`scripts/**`、`templates/**` 或 `.github/workflows/**`，或手动触发。
- 版本计算：运行 `get-next-version.sh`，读取最新 tag（无则 `v0.0.0`），补丁位 +1，输出到 `GITHUB_OUTPUT`。
- 释放检查：运行 `check-release-exists.sh <version>`，用 `gh release view` 判断是否已有该版本，已存在则整条流程跳过。
- 资产打包：运行 `create-release-packages.sh <version>`，按 AI 助手 × 脚本类型生成 ZIP 到 `.genreleases/`。
- 生成说明：运行 `generate-release-notes.sh <new_version> <last_tag>`，从 `git log` 生成 `release_notes.md`。
- 创建发布：运行 `create-github-release.sh <version>`，使用 `gh release create` 上传所有 ZIP，并附上 `release_notes.md`。
- 版本写回：运行 `update-version.sh <version>`，将 `pyproject.toml` 的 `version` 更新为不带 `v` 前缀的版本号（仅用于发布构件）。

**打包脚本要点（`.github/workflows/scripts/create-release-packages.sh`）**
- 构建产物目录：所有中间与最终 ZIP 放在 `.genreleases/`。
- 变体选择：
  - 代理列表：默认全量，`ALL_AGENTS=(claude gemini copilot cursor-agent qwen opencode windsurf codex kilocode auggie roo codebuddy amp q)`。
  - 脚本类型：`ALL_SCRIPTS=(sh ps)`。
  - 可通过环境变量限制：`AGENTS="copilot,gemini"` 或 `SCRIPTS=ps` 等。
- 每个变体调用 `build_variant <agent> <script>`：
  - 目录搭建：在 `.genreleases/sdd-{agent}-package-{script}` 下构建最终结构。
  - 复制内存：`memory/**` → `.specify/memory/`（存在时）。
  - 复制脚本：
    - `script=sh`：复制 `scripts/bash/**` → `.specify/scripts/`，并复制 `scripts` 根目录下的非变体脚本（`find scripts -maxdepth 1 -type f`）。
    - `script=ps`：复制 `scripts/powershell/**` → `.specify/scripts/`，并复制 `scripts` 根目录下的非变体脚本。
  - 复制模板：除命令模板与 `vscode-settings.json` 外的 `templates/**` → `.specify/templates/`，用 `cp --parents` 保持相对路径。
  - 生成命令模板：从 `templates/commands/*.md` 生成各代理所需的命令文件（见下节）。
  - 特定代理补充：
    - `copilot`：生成到 `.github/prompts`，并复制 `templates/vscode-settings.json` → `.vscode/settings.json`。
    - `gemini`、`qwen`：生成到各自的 `.{agent}/commands`（TOML 格式），并可额外复制 `agent_templates/gemini/GEMINI.md` 或 `agent_templates/qwen/QWEN.md`。
    - 其他代理：分别生成到各自目录（如 `.claude/commands`、`.cursor/commands`、`.windsurf/workflows` 等）。
  - 压缩归档：`zip -r "../spec-kit-template-{agent}-{script}-{VERSION}.zip" .`，输出到 `.genreleases/`。

**命令模板生成（`generate_commands`）**
- 输入源：`templates/commands/*.md`，带 YAML frontmatter（含 `description`、`scripts:`、可选 `agent_scripts:`）。
- 脚本命令注入：
  - 根据变体读取 `scripts:` 下的 `sh:` 或 `ps:` 值，替换正文中的 `{SCRIPT}`。
  - 若存在 `agent_scripts:`（按变体），替换正文中的 `{AGENT_SCRIPT}`。
- 占位符替换：
  - `{ARGS}`：按代理选择格式（Markdown/Prompt 为 `"$ARGUMENTS"`，TOML 为 `"{{args}}"`）。
  - `__AGENT__`：替换为代理名。
  - 路径重写：`memory/` → `.specify/memory/`，`scripts/` → `.specify/scripts/`，`templates/` → `.specify/templates/`。
- frontmatter清理：移除 `scripts:` 与 `agent_scripts:` 段，保留其余 YAML。
- 输出位置与扩展名：
  - `copilot`：`.github/prompts/speckit.<cmd>.prompt.md`
  - `claude`、`cursor-agent` 等：`.md`
  - `gemini`、`qwen`：`.toml`

**发布创建（`.github/workflows/scripts/create-github-release.sh`）**
- 统一上传 `.genreleases/` 下所有 `spec-kit-template-{agent}-{script}-{version}.zip`。
- 标题为 `"Spec Kit Templates - {version_no_v}"`，备注取自 `release_notes.md`。

**版本/发布辅助脚本**
- `get-next-version.sh`：读取最新 tag，补丁位 +1，输出 `latest_tag` 与 `new_version`。
- `check-release-exists.sh`：用 `gh release view <version>` 判断是否存在，存在则将 `exists=true` 输出到 `GITHUB_OUTPUT`。
- `generate-release-notes.sh`：取上次 tag 到当前之间的提交摘要生成 `release_notes.md`。
- `update-version.sh`：将 `pyproject.toml` 的 `version` 改为不带 `v` 的版本号。

**命名规范与产物位置**
- ZIP 命名：`spec-kit-template-{agent}-{script}-{vX.Y.Z}.zip`，如 `spec-kit-template-copilot-ps-v0.0.72.zip`。
- 构建目录：`.genreleases/sdd-{agent}-package-{script}` 为打包前目录，ZIP 汇总在 `.genreleases/`。

**本地复刻（可选）**
- 在支持 Bash 的环境中（WSL/Git Bash）运行：
  - `bash .github/workflows/scripts/get-next-version.sh`
  - `bash .github/workflows/scripts/create-release-packages.sh v0.0.78`
  - ZIP 将生成到 `.genreleases/`；可用 `bash .github/workflows/scripts/create-github-release.sh v0.0.78` 上传（需 `gh` 登录）。

这套 `.github` 工作流和脚本把本地 `templates/`、`scripts/`、`memory/` 等源，统一映射到 `.specify/` 结构，按代理与脚本类型生成命令文件，再打成标准化的模板 ZIP 并发布到 GitHub Releases。

## bash .github/workflows/scripts/get-next-version.sh运行显示.github/workflows/scripts/get-next-version.sh: line 10: GITHUB_OUTPUT: unbound variable

**原因**
- 这个脚本是为 GitHub Actions 设计的，使用 `GITHUB_OUTPUT` 写 Step 的输出值。
- 本地直接运行时没有该环境变量；在 `set -euo pipefail` 下，未定义变量会触发 “unbound variable”。

**解决方案**
- 在本地运行前显式设置 `GITHUB_OUTPUT` 指向一个可写文件。
- 两种方式都可以：
  - Bash（Git Bash/WSL）: `export GITHUB_OUTPUT=.github_output.txt`
  - PowerShell: `$env:GITHUB_OUTPUT = "$PWD\.github_output.txt"`

**本地运行示例**
- PowerShell（推荐在你当前终端环境）：
  - 设置输出文件: ``$env:GITHUB_OUTPUT = "$PWD\.github_output.txt"``
  - 运行脚本: ``bash .github/workflows/scripts/get-next-version.sh``
  - 查看结果: ``Get-Content .github_output.txt``
- Bash（Git Bash/WSL）：
  - 设置输出文件: ``export GITHUB_OUTPUT="$(mktemp)"`` 或 ``export GITHUB_OUTPUT=.github_output.txt``
  - 运行脚本: ``bash .github/workflows/scripts/get-next-version.sh``
  - 查看结果: ``cat "$GITHUB_OUTPUT"``

你会看到两行输出键值：
- `latest_tag=<最新标签或 v0.0.0>`
- `new_version=<自动+1后的版本，如 v0.0.78>`

**可选的脚本健壮性改造**
- 若你希望脚本在本地无环境变量也能运行，可在 `get-next-version.sh` 开头添加回退逻辑（在 `set -u` 下安全的写法）：
  - 在 `set -euo pipefail` 之后加：
    - `if [[ -z "${GITHUB_OUTPUT:-}" ]]; then export GITHUB_OUTPUT="$(pwd)/.github_output.txt"; fi`
- 同类依赖 `GITHUB_OUTPUT` 的脚本（如 `.github/workflows/scripts/check-release-exists.sh`）也建议加同样的回退。

**后续打包与发布（本地）**
- 设定版本并构建 ZIP（仅打包，不上传）：
  - ``export GITHUB_OUTPUT=.github_output.txt``
  - 生成下一个版本: ``bash .github/workflows/scripts/get-next-version.sh``
  - 从 `.github_output.txt` 读取 `new_version`，然后执行:
  - ``bash .github/workflows/scripts/create-release-packages.sh v0.0.78``（替换为上一步得到的版本）
- 如果要本地上传到 GitHub Release，需要安装并登录 `gh` CLI：
  - 检查是否已存在: ``bash .github/workflows/scripts/check-release-exists.sh v0.0.78``
  - 创建发布并上传: ``bash .github/workflows/scripts/create-github-release.sh v0.0.78``

这样在本地即可完整复刻工作流中的版本计算与打包过程，同时避免 `GITHUB_OUTPUT` 未定义的错误。
        

## 执行 bash .github/workflows/scripts/create-github-release.sh报错：
.github/workflows/scripts/create-github-release.sh: line 18: gh: command not found
     
**错误原因**
- `gh` 是 GitHub CLI。脚本用它调用 GitHub Releases（`gh release create`/`gh release view`）。
- 你的本地环境没有安装 `gh`；因此 Bash 报错 `gh: command not found`。

**安装 GitHub CLI（Windows）**
- winget（推荐）:
  - `winget install --id GitHub.cli -e`
- scoop（若已装 scoop）:
  - `scoop install gh`
- Chocolatey:
  - `choco install gh -y`
- 验证安装:
  - `gh --version`

**认证与令牌**
- 交互登录（最简单）:
  - `gh auth login`
  - 选择 `GitHub.com` → `HTTPS` → Paste token（或浏览器登录）
- 或设置环境变量（无交互）:
  - PowerShell: `$env:GH_TOKEN = "<你的PAT>"` 或 `$env:GITHUB_TOKEN = "<你的PAT>"`
  - PAT 至少需要 `repo` 权限（创建 release 和上传资产）。

**确保目标仓库正确**
- 默认会在当前仓库的 `origin` 对应的 GitHub 仓库创建发布：
  - 检查：`git remote -v`
  - 如需指定仓库，可改脚本为：`gh release create "$VERSION" --repo <owner>/<repo> ...`
- 建议先用：`bash .github/workflows/scripts/check-release-exists.sh vX.Y.Z` 验证 CLI 连通性与仓库指向。

**本地完整打包与发布流程**
- 设置 `GITHUB_OUTPUT`（供版本脚本写出值）:
  - PowerShell: `$env:GITHUB_OUTPUT = "$PWD\.github_output.txt"`
- 计算版本：
  - `bash .github/workflows/scripts/get-next-version.sh`
  - 查看：`Get-Content .github_output.txt`（包含 `latest_tag=` 与 `new_version=`）
- 生成 ZIP：
  - `bash .github/workflows/scripts/create-release-packages.sh vX.Y.Z`
  - 产物在 `.genreleases/`
- 生成说明：
  - `bash .github/workflows/scripts/generate-release-notes.sh vX.Y.Z <latest_tag>`
- 检查发布是否存在：
  - `bash .github/workflows/scripts/check-release-exists.sh vX.Y.Z`
- 创建发布并上传资产：
  - `bash .github/workflows/scripts/create-github-release.sh vX.Y.Z`

**不安装 gh 的替代方案**
- Web UI：直接在 GitHub 仓库的 Releases 页面创建新版本，手动上传 `.genreleases/*.zip`，将 `release_notes.md` 内容复制为备注。
- 纯 API（高级）：用 `curl` 调 GitHub Releases API 创建 release 并上传资产，但流程更复杂，不如安装 `gh` 简单。

**常见坑位与提示**
- PATH 更新：安装后如 `gh --version` 仍失败，重启终端或确保 PATH 包含 GitHub CLI 安装路径。
- Bash 环境：你在 PowerShell 调用 `bash`；确保 `gh.exe` 在系统 PATH 中，Git Bash 也能找到它。
- 令牌优先级：设置了 `GH_TOKEN`/`GITHUB_TOKEN` 后，可不需要 `gh auth login`；否则会进入交互登录。
- 仓库来源：如果当前项目不是你要发布到的 GitHub 仓库，务必在脚本里加 `--repo <owner>/<repo>`，避免发布到错误仓库。

## release_notes.md是在哪生成的？open release_notes.md: The system cannot find the file specified.

在哪里生成

- release_notes.md 在仓库根目录生成（ d:\GoWorks\src\require-gen\release_notes.md ）。
- 在 GitHub Actions 中由步骤 “Generate release notes” 调用 .github/workflows/scripts/generate-release-notes.sh <new_version> <last_tag> 生成。
- 我已为你在本地复刻生成，文件已创建。
为什么你之前找不到

- 你还未运行“生成发布说明”的脚本，或在本地环境没有 Bash（导致脚本未执行），所以文件不存在。
- 该脚本不会自动运行，需显式传入 new_version 与 last_tag 参数。
本地生成方式

- 有 Bash（Git Bash/WSL）时：
  - 设置版本输出文件： export GITHUB_OUTPUT=.github_output.txt
  - 计算版本： bash .github/workflows/scripts/get-next-version.sh
  - 读取值： cat .github_output.txt （得到 latest_tag=... 、 new_version=... ）
  - 生成说明： bash .github/workflows/scripts/generate-release-notes.sh vX.Y.Z vA.B.C
- 无 Bash（你当前的 PowerShell 环境）：
  - 我已用 PowerShell生成了 release_notes.md ，内容与脚本逻辑一致（取最近 10 条或从 latest_tag 到 HEAD 的提交信息）。
确认文件

- 路径： d:\GoWorks\src\require-gen\release_notes.md
- 打开查看： Get-Content release_notes.md
后续发布

- 安装 GitHub CLI（ gh ）后，可运行：
  - bash .github/workflows/scripts/check-release-exists.sh vX.Y.Z
  - bash .github/workflows/scripts/create-github-release.sh vX.Y.Z
- 或直接到 GitHub Releases 页面，手动创建发布并上传 .genreleases/*.zip ，将 release_notes.md 内容粘贴为说明。

## 脚本清单




          
**脚本清单**
- `d:\GoWorks\src\require-gen\.github\workflows\scripts\get-next-version.sh`
  - 作用：读取最近的 Git 标签（无则视为 `v0.0.0`），把补丁位 +1，生成下一个版本号。
  - 输出：向 `GITHUB_OUTPUT` 写入两行键值：`latest_tag=...`、`new_version=...`。
  - 依赖：`git`；在本地运行需先设置 `GITHUB_OUTPUT`（例如 `export GITHUB_OUTPUT=.github_output.txt`）。
  - 用法示例：`bash .github/workflows/scripts/get-next-version.sh` 后查看 `.github_output.txt`。

- `d:\GoWorks\src\require-gen\.github\workflows\scripts\check-release-exists.sh`
  - 作用：检查指定版本的 GitHub Release 是否已存在，决定后续步骤是否跳过。
  - 输出：向 `GITHUB_OUTPUT` 写入 `exists=true` 或 `exists=false`。
  - 依赖：`gh`（GitHub CLI）；需要登录或设置令牌（`GH_TOKEN`/`GITHUB_TOKEN`）。
  - 用法示例：`bash .github/workflows/scripts/check-release-exists.sh v0.0.78`。

- `d:\GoWorks\src\require-gen\.github\workflows\scripts\create-release-packages.sh`
  - 作用：为所有支持的 AI 助手（`claude`、`gemini`、`copilot` 等）和脚本类型（`sh`/`ps`）生成模板 ZIP 包。
  - 输入：位置参数 `<version>`（形如 `v0.0.78`）；可选环境变量 `AGENTS`（限制代理集合）、`SCRIPTS`（限制脚本集合）。
  - 行为：
    - 在 `.genreleases/` 下为每个组合构建目录 `sdd-{agent}-package-{script}`。
    - 映射资源：`memory/**` → `.specify/memory/`；`scripts/bash/**` 或 `scripts/powershell/**` → `.specify/scripts/`；`templates/**`（除命令与 `vscode-settings.json`）→ `.specify/templates/`。
    - 生成命令文件：从 `templates/commands/*.md` 读取 frontmatter 的 `scripts:` 和可选 `agent_scripts:`，替换 `{SCRIPT}`、`{AGENT_SCRIPT}`、`{ARGS}`，并重写路径到 `.specify/...`；根据代理输出到不同目录（如 `copilot` → `.github/prompts`、`gemini/qwen` → TOML 到各自 `.gemini/.qwen`）。
    - 生成 ZIP：`spec-kit-template-{agent}-{script}-{version}.zip` 写入 `.genreleases/`。
  - 依赖：`zip`、`sed`、`awk`、`find`、`cp` 等常用 *nix 工具。
  - 用法示例：`bash .github/workflows/scripts/create-release-packages.sh v0.0.78`。

- `d:\GoWorks\src\require-gen\.github\workflows\scripts\generate-release-notes.sh`
  - 作用：根据提交历史生成 `release_notes.md`。
  - 输入：位置参数 `<new_version> <last_tag>`。
  - 行为：
    - 若 `last_tag` 为 `v0.0.0`，取最近 10 条提交（或全部提交）；否则取 `last_tag..HEAD` 区间。
    - 在仓库根生成 `release_notes.md` 并打印预览。
  - 依赖：`git`。
  - 用法示例：`bash .github/workflows/scripts/generate-release-notes.sh v0.0.78 v0.0.77`。

- `d:\GoWorks\src\require-gen\.github\workflows\scripts\create-github-release.sh`
  - 作用：创建指定版本的 GitHub Release，并一次性上传所有模板 ZIP。
  - 输入：位置参数 `<version>`（例如 `v0.0.78`）。
  - 行为：
    - 通过 `gh release create` 上传 `.genreleases/spec-kit-template-*-{version}.zip` 的全套组合。
    - 标题：`Spec Kit Templates - {version_no_v}`；备注从仓库根的 `release_notes.md` 读取。
  - 依赖：`gh` CLI；需要本地已存在 `.genreleases/*.zip` 和 `release_notes.md`。
  - 用法示例：`bash .github/workflows/scripts/create-github-release.sh v0.0.78`。

- `d:\GoWorks\src\require-gen\.github\workflows\scripts\update-version.sh`
  - 作用：将 `pyproject.toml` 的 `version` 替换为不带 `v` 前缀的版本号，以便发布工件标识。
  - 输入：位置参数 `<version>`；示例 `v0.0.78` → 写入 `0.0.78`。
  - 行为：使用 `sed -i` 就地替换；若找不到 `pyproject.toml` 则输出警告。
  - 依赖：`sed`。
  - 用法示例：`bash .github/workflows/scripts/update-version.sh v0.0.78`。

**工作流关联**
- `d:\GoWorks\src\require-gen\.github\workflows\release.yml` 按顺序调用：
  - 获取新版本（`get-next-version.sh`）→ 检查是否已存在（`check-release-exists.sh`）→ 打包（`create-release-packages.sh`）→ 生成说明（`generate-release-notes.sh`）→ 创建发布（`create-github-release.sh`）→ 更新版本号（`update-version.sh`）。

**本地运行提示（Windows）**
- 在 PowerShell 环境中运行 Bash 脚本需要 Git Bash 或 WSL；否则可用 PowerShell复刻逻辑（如我已为你在仓库根生成了 `release_notes.md`）。
- 使用 `gh` 时先安装并登录：
  - 安装：`winget install --id GitHub.cli -e`
  - 登录：`gh auth login` 或设置 `GH_TOKEN`/`GITHUB_TOKEN`。
- 本地计算版本需设置 `GITHUB_OUTPUT`：
  - 例如：`$env:GITHUB_OUTPUT = \"$PWD\\.github_output.txt\"` 后运行 `bash .github/workflows/scripts/get-next-version.sh`。
        