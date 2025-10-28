**目录作用概览**
- `templates/commands/`：为斜杠命令提供标准化提示与执行准则（微提示/操作手册），对应 “/constitution, /specify, /plan, /tasks, /implement, /clarify, /analyze”。
- `spec-template.md`：产品/技术规格说明的骨架文档，配合 `/specify` 生成规范初稿。
- `plan-template.md`：实施计划的骨架文档，配合 `/plan` 生成包含里程碑、风险与测试策略的计划文档。
- `tasks-template.md`：任务分解的骨架文档，配合 `/tasks` 产出可执行任务列表（含依赖与验收标准）。
- `checklist-template.md`：检查清单骨架，用于质量保障与发布前自检；常与 `/analyze` 或手动 QA 流结合。
- `agent-file-template.md`：新增/维护 AI 助手配置的模板（名称、键、是否需要 CLI、安装链接等），与内部 `core.AgentConfig` 生态衔接。
- `vscode-settings.json`：VS Code 项目级设置，提升 Markdown/代码编辑体验（如格式化、换行、文件关联等）。

**与斜杠命令的关系**
- 命令到模板的映射：
  - `/constitution` → 指导原则与工作方式（来自 `commands/constitution.md`），通常落地为项目“宪章”文档（如 `memory/constitution.md`）。
  - `/specify` → 读取 `spec-template.md` 作为骨架，结合 `commands/specify.md` 的方法提示，生成规范文档。
  - `/plan` → 读取 `plan-template.md`，配合 `commands/plan.md` 指南，形成实施计划（里程碑/风险/测试策略）。
  - `/tasks` → 读取 `tasks-template.md`，参照 `commands/tasks.md` 组织任务分解与依赖/验收标准。
  - `/implement` → 参照 `commands/implement.md` 的执行说明，结合脚本类型（PowerShell/Bash）触发对应的项目脚本与落地动作。
  - `/clarify`（可选前置）→ 使用 `commands/clarify.md` 结构化澄清问题，在 `/plan` 前降低歧义。
  - `/analyze`（可选中后段）→ 使用 `commands/analyze.md` 进行跨工件一致性分析，常在 `/tasks` 后、`/implement` 前进行。
- 关系链路（典型流）：`/constitution` → `/specify` → `/plan` → `/tasks` → `/implement`，中间可插入 `/clarify` 与 `/analyze` 增强质量与一致性。

**模板之间的协同**
- `commands/*.md` 提供操作方法与内容期望（指南），`spec/plan/tasks-*template.md` 提供文档结构（骨架）。前者指导如何填充，后者规定放什么、放在哪。
- `checklist-template.md` 独立于三大文档，可在 `/analyze` 产出“检查清单”，或在发布/验收阶段手工使用。
- `agent-file-template.md` 不直接参与文档产出，但影响 `/init`、工具检测与 AI 斜杠命令可用性（如是否需要安装 CLI）。
- `vscode-settings.json` 保证在 VS Code 中编辑这些 Markdown/脚本文件时体验一致，利于团队协作与规范统一。

**关键字段与示例说明**
- `plan-template.md` 的测试字段示例：
  - `**Testing**: [e.g., pytest, XCTest, cargo test or NEEDS CLARIFICATION]`
  - 这里的 “NEEDS CLARIFICATION” 是刻意的占位标识，提示当技术栈不明时先用 `/clarify` 明确测试框架与策略，再回写到计划文档。
- `spec-template.md` 通常包含“背景/范围/需求/非功能/约束/开放问题”等结构，结合 `/specify` 的方法提示来填充。
- `tasks-template.md` 约定任务条目包含“目标/步骤/依赖/验收标准/负责人”等，方便后续映射到脚本与实现。

**落地与存放位置（建议）**
- 初始成果通常生成在项目根或 `memory/`、`.specify/` 等目录下，便于版本控制与复盘。
- 在你的仓库中已存在 `memory/constitution.md`，说明宪章文档已被生成/维护；其他文档可沿用同样的存放策略（例如 `docs/` 或 `memory/`）。

**使用与定制建议**
- 先根据你的技术栈定制 `spec/plan/tasks-*template.md` 的章节与术语，让 `/specify`、`/plan`、`/tasks` 输出更贴近团队语言。
- 将 `plan-template.md` 的测试字段具体化为你的栈（如 `pytest`、`vitest`、`go test`、`XCTest`），去掉模糊占位。
- 若你引入新的 AI 助手或本地 CLI，按 `agent-file-template.md` 要求补齐配置后，再使用 `/init` 与斜杠命令确保工具检测与指导面板一致。
- 若团队需要更严格的规范检查，把 `checklist-template.md` 与 `/analyze` 流程结合，形成 PR 审查清单或发布门禁。

**与 Init 指导面板的联动**
- “Next Steps” 面板引导你按斜杠命令执行的顺序推进；这些命令对应 `templates/commands/`。
- “Script Commands” 面板给出 PowerShell/Bash 的执行格式；与 `/implement` 的执行落地相配合。
- “Project Commands” 提示 `task-kit check`、`code .` 与 `quickstart.md`，完善环境与文档入口，帮助你在模板基础上快速迭代。

如果你希望，我可以按你的技术栈（如 Go/Node/Python/iOS）为 `spec/plan/tasks-*template.md` 给出一版定制化骨架，并映射到你当前的脚本目录结构。
        