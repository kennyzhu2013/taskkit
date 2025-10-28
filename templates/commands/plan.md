---
description: 使用 plan 模板执行实施规划工作流以生成设计制品。
scripts:
  sh: scripts/bash/setup-plan.sh --json
  ps: scripts/powershell/setup-plan.ps1 -Json
agent_scripts:
  sh: scripts/bash/update-agent-context.sh __AGENT__
  ps: scripts/powershell/update-agent-context.ps1 -AgentType __AGENT__
---

## 用户输入

```text
$ARGUMENTS
```

在继续之前，若用户输入不为空，你必须考虑它。

## 大纲（Outline）

1. **Setup**：从仓库根目录运行 `{SCRIPT}` 并解析 JSON：FEATURE_SPEC、IMPL_PLAN、SPECS_DIR、BRANCH。对于单引号参数（如 "I'm Groot"），使用转义：例如 `'I'\''m Groot'`（或尽可能使用双引号："I'm Groot"）。

2. **Load context**：读取 FEATURE_SPEC 与 `/memory/constitution.md`。加载已复制的 IMPL_PLAN 模板。

3. **Execute plan workflow**：按照 IMPL_PLAN 模板结构：
   - 填写 Technical Context（未知项标注为 "NEEDS CLARIFICATION"）
   - 从宪章填充 Constitution Check 章节
   - 评估闸口（若存在未合理化的违规则 ERROR）
   - Phase 0：生成 `research.md`（解决所有 NEEDS CLARIFICATION）
   - Phase 1：生成 `data-model.md`、`contracts/`、`quickstart.md`
   - Phase 1：运行 Agent 脚本更新 Agent 上下文
   - 设计完成后重新评估 Constitution Check

4. **Stop and report**：命令在 Phase 2 规划结束后终止。报告分支、IMPL_PLAN 路径与已生成制品。

## Phases

### Phase 0: Outline & Research

1. **从 Technical Context 提取未知项**：
   - 每个 NEEDS CLARIFICATION → research 任务
   - 每个 dependency → best practices 任务
   - 每个 integration → patterns 任务

2. **生成并派发 research agents**：
   ```
   For each unknown in Technical Context:
     Task: "Research {unknown} for {feature context}"
   For each technology choice:
     Task: "Find best practices for {tech} in {domain}"
   ```

3. **整合发现**至 `research.md`，格式：
   - Decision: [what was chosen]
   - Rationale: [why chosen]
   - Alternatives considered: [what else evaluated]

**Output**：包含所有 NEEDS CLARIFICATION 的 `research.md`

### Phase 1: Design & Contracts

**Prerequisites**：`research.md` 完成

1. **从 Feature Spec 提取实体** → `data-model.md`：
   - Entity name、fields、relationships
   - 来自 Requirements 的校验规则
   - 如适用的状态转换

2. **从 Functional Requirements 生成 API contracts**：
   - 每个用户动作 → endpoint
   - 使用标准 REST/GraphQL 模式
   - 将 OpenAPI/GraphQL schema 输出到 `/contracts/`

3. **Agent context update**：
   - 运行 `{AGENT_SCRIPT}`
   - 这些脚本检测当前使用的 AI agent
   - 更新相应的 Agent 专用上下文文件
   - 仅添加当前计划中的新技术
   - 在标记之间保留人工添加内容

**Output**：`data-model.md`、`/contracts/*`、`quickstart.md`、Agent 专用文件

## Key rules

- 使用绝对路径
- 在闸口失败或澄清未解决时 ERROR