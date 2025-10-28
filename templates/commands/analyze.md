---
description: 在生成 `tasks.md` 后，对 `spec.md`、`plan.md`、`tasks.md` 进行非破坏性的跨制品一致性与质量分析。
scripts:
  sh: scripts/bash/check-prerequisites.sh --json --require-tasks --include-tasks
  ps: scripts/powershell/check-prerequisites.ps1 -Json -RequireTasks -IncludeTasks
---

## 用户输入

```text
$ARGUMENTS
```

在继续之前，若用户输入不为空，你必须考虑它。

## 目标

在实施之前，识别三大核心制品（`spec.md`、`plan.md`、`tasks.md`）中的不一致、重复、模糊和未充分规定项。该命令必须在 `/tasks` 成功生成完整的 `tasks.md` 之后运行。

## 运行约束

**STRICTLY READ-ONLY**：严格只读，禁止修改任何文件。输出结构化的分析报告，可提供可选的修复方案（用户需明确批准后，才会由后续编辑类命令手动执行）。

**Constitution Authority**：项目宪章（`/memory/constitution.md`）在本次分析范围内具有不可协商的权威。与宪章冲突的内容自动判定为 CRITICAL，必须通过调整 `spec`、`plan` 或 `tasks` 解决，而不是削弱、曲解或忽略原则。若原则本身需要变更，必须在 `/analyze` 之外以单独的显式宪章更新进行。

## 执行步骤

### 1. 初始化分析上下文

从仓库根目录运行 `{SCRIPT}` 一次，并解析 JSON 获取 FEATURE_DIR 与 AVAILABLE_DOCS。派生绝对路径：

- SPEC = FEATURE_DIR/spec.md
- PLAN = FEATURE_DIR/plan.md
- TASKS = FEATURE_DIR/tasks.md

若任一必需文件缺失则中止并提示错误（指导用户运行缺失的前置命令）。对于单引号参数（如 "I'm Groot"），使用转义：例如 `'I'\''m Groot'`（或尽可能使用双引号："I'm Groot"）。

### 2. 渐进式加载制品

仅加载每个制品的必要上下文：

**来自 spec.md：**
- Overview/Context
- Functional Requirements
- Non-Functional Requirements
- User Stories
- Edge Cases（如存在）

**来自 plan.md：**
- Architecture/stack choices
- Data Model 引用
- Phases
- Technical constraints

**来自 tasks.md：**
- Task IDs
- Descriptions
- Phase 分组
- 并行标记 [P]
- 引用的文件路径

**来自宪章：**
- 加载 `/memory/constitution.md` 以进行原则校验

### 3. 构建语义模型

创建内部表示（不要在输出中包含原始制品）：
- Requirements inventory：为每条 Functional 与 Non-Functional Requirement 派生稳定键（基于祈使短语生成 slug，例如 "User can upload file" → `user-can-upload-file`）
- User story/action inventory：离散的用户动作及其 Acceptance Criteria
- Task coverage mapping：将每个任务映射到一个或多个 Requirements 或 User Stories（通过关键字/显式引用模式如 ID 或关键短语推断）
- Constitution rule set：抽取原则名及 MUST/SHOULD 规范性陈述

### 4. 检测通道（Token 高效）

聚焦高信号发现。限定最多 50 条发现，其余汇总到溢出摘要。

#### A. Duplication Detection（重复）
- 识别近似重复的 Requirements
- 标记表述质量较低者以供合并

#### B. Ambiguity Detection（模糊）
- 标记缺乏可度量标准的形容词（fast、scalable、secure、intuitive、robust）
- 标记未解决占位符（TODO、TKTK、???、`<placeholder>` 等）

#### C. Underspecification（未充分规定）
- 仅有动词而缺少对象或可度量结果的 Requirements
- 缺少与 Acceptance Criteria 对齐的 User Stories
- 引用 `spec/plan` 未定义文件或组件的 Tasks

#### D. Constitution Alignment（宪章对齐）
- 任何与宪章 MUST 原则冲突的 Requirement 或 Plan 元素
- 宪章要求但缺失的强制性章节或质量闸口

#### E. Coverage Gaps（覆盖缺口）
- 无任何关联任务的 Requirements
- 无映射 Requirement/Story 的 Tasks
- 未在 Tasks 体现的 Non-Functional Requirements（例如性能、安全）

#### F. Inconsistency（不一致）
- 术语漂移（同一概念跨文档命名不一致）
- Plan 中引用的数据实体在 Spec 缺失（或反之）
- 任务排序矛盾（如在未声明依赖的情况下，集成任务早于基础设置任务）
- 冲突的 Requirements（例如一处要求 Next.js，另一处指定 Vue）

### 5. 赋予严重级别

使用以下启发式进行优先级：
- CRITICAL：违反宪章 MUST、缺失核心 Spec 制品、阻断基础功能的零覆盖 Requirement
- HIGH：重复或冲突 Requirement、关于安全/性能的模糊属性、不可测试的 Acceptance Criterion
- MEDIUM：术语漂移、Non-Functional 任务缺失、未充分规定的 Edge Case
- LOW：风格/措辞改进、对执行顺序无影响的轻微冗余

### 6. 生成精简分析报告

输出 Markdown 报告（不写文件），结构如下：

## 规范分析报告

| ID | 类别 | 严重级别 | 位置 | 摘要 | 建议 |
|----|----------|----------|-------------|---------|----------------|
| A1 | Duplication | HIGH | spec.md:L120-134 | Two similar requirements ... | Merge phrasing; keep clearer version |

覆盖摘要表：

| Requirement Key | Has Task? | Task IDs | Notes |
|-----------------|-----------|----------|-------|

Constitution Alignment Issues：（如有）

Unmapped Tasks：（如有）

度量：
- Total Requirements
- Total Tasks
- Coverage %（requirements with >=1 task）
- Ambiguity Count
- Duplication Count
- Critical Issues Count

### 7. 提供下一步行动

在报告末尾提供精炼的 Next Actions：
- 若存在 CRITICAL：建议在 `/implement` 前先解决
- 若仅 LOW/MEDIUM：可继续，但给出改进建议
- 提供明确的命令建议：如 “Run /specify with refinement”、“Run /plan to adjust architecture”、“手动编辑 tasks.md 为 'performance-metrics' 添加覆盖”

### 8. 提供修复建议

询问用户：“是否需要我为前 N 个问题提出具体的修复编辑建议？”（不要自动应用）。

## 运行原则

### Context Efficiency（上下文效率）
- Minimal high-signal tokens：聚焦可执行发现，不做穷尽文档
- Progressive disclosure：按需增量加载制品，避免一次性倾倒
- Token-efficient output：发现表限制 50 行，超出部分摘要
- Deterministic results：无变更重复运行应产生一致的 ID 与计数

### Analysis Guidelines（分析指南）
- NEVER 修改文件（严格只读）
- NEVER 虚构缺失章节（如确实缺失，准确报告）
- 优先处理宪章违规（总是 CRITICAL）
- 使用具体实例而非泛化规则（引用明确位置）
- 零问题也要优雅报告（输出成功报告与覆盖统计）

## Context

{ARGS}