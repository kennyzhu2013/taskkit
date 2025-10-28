---
description: 基于可用的设计制品为该功能生成可执行的、按依赖排序的 tasks.md。
scripts:
  sh: scripts/bash/check-prerequisites.sh --json
  ps: scripts/powershell/check-prerequisites.ps1 -Json
---

## 用户输入

```text
$ARGUMENTS
```

在继续之前（且用户输入不为空时），你必须先充分考虑用户输入。

## 大纲

1. **Setup**：从仓库根目录运行 `{SCRIPT}`，解析 FEATURE_DIR 和 AVAILABLE_DOCS 列表。所有路径必须为绝对路径。对于参数中包含单引号的内容（如 "I'm Groot"），请使用转义语法：例如 `I'\''m Groot`（或尽可能改用双引号：`"I'm Groot"`）。

2. **加载设计文档**：从 FEATURE_DIR 读取：
   - **必需**：`plan.md`（技术栈、库、结构）、`spec.md`（用户故事及其优先级）
   - **可选**：`data-model.md`（实体）、`contracts/`（API 端点）、`research.md`（决策）、`quickstart.md`（测试场景）
   - 注意：并非所有项目都具备全部文档。应根据“可用文档”生成任务。

3. **执行任务生成工作流**：
   - 加载 `plan.md`，提取技术栈、依赖库、项目结构
   - 加载 `spec.md`，提取用户故事及其优先级（P1、P2、P3 等）
   - 若存在 `data-model.md`：提取实体并映射到用户故事
   - 若存在 `contracts/`：将端点映射到用户故事
   - 若存在 `research.md`：提取与初始化相关的决策，用于 Setup 任务
   - 按用户故事组织生成任务（见下方“任务生成规则”）
   - 生成显示用户故事完成顺序的依赖关系图
   - 为每个用户故事创建并行执行示例
   - 校验任务完备性（每个用户故事都有完成所需的全部任务，并可独立测试）

4. **生成 tasks.md**：使用 `.specify/templates/tasks-template.md` 作为结构模板，填充：
   - 来自 `plan.md` 的正确功能名称
   - Phase 1：Setup 任务（项目初始化）
   - Phase 2：Foundational 任务（所有用户故事的前置阻塞项）
   - Phase 3+：按 `spec.md` 中的优先级（P1、P2、P3…）为每个用户故事建立独立 Phase
   - 每个 Phase 包含：故事目标、独立测试标准、（如请求）测试任务、实施任务
   - Final Phase：打磨与跨切面关注点
   - 所有任务必须严格遵循清单格式（见下方“任务生成规则”）
   - 为每个任务提供清晰的文件路径
   - 依赖关系部分给出用户故事的完成顺序
   - 每个故事提供并行执行示例
   - 实施策略部分（MVP 优先、增量交付）

5. **报告**：输出生成的 `tasks.md` 路径与摘要：
   - 总任务数
   - 每个用户故事的任务数
   - 识别到的并行机会
   - 每个故事的独立测试标准
   - 建议的 MVP 范围（通常为 User Story 1）
   - 格式校验：确认所有任务均遵循清单格式（复选框、ID、标签、文件路径）

任务生成上下文：{ARGS}

`tasks.md` 必须“可立即执行”——每个任务都应足够具体，使得一个 LLM 无需额外上下文即可完成。

## 任务生成规则

**关键要求**：任务必须以用户故事为主组织维度，以支持独立实施与测试。

**测试为可选**：仅当功能规范明确要求或用户请求采用 TDD 方法时，才生成测试相关任务。

### 清单格式（必填）

每个任务必须严格遵循以下格式：

```text
- [ ] [TaskID] [P?] [Story?] 描述（含文件路径）
```

**格式组件**：

1. **复选框**：始终以 `- [ ]` 开始（Markdown 复选框）
2. **任务 ID**：按执行顺序的序号（T001、T002、T003…）
3. **[P] 标记**：仅在任务可并行（涉及不同文件，且不依赖未完成任务）时加入
4. **[Story] 标签**：仅在“用户故事”阶段的任务中必须包含
   - 形式：`[US1]`、`[US2]`、`[US3]` 等（映射自 `spec.md` 的用户故事）
   - Setup 阶段：不含故事标签
   - Foundational 阶段：不含故事标签  
   - 用户故事阶段：必须含故事标签
   - Polish 阶段：不含故事标签
5. **描述**：清晰的操作说明，包含精确文件路径

**示例**：

- ✅ 正确：`- [ ] T001 Create project structure per implementation plan`
- ✅ 正确：`- [ ] T005 [P] Implement authentication middleware in src/middleware/auth.py`
- ✅ 正确：`- [ ] T012 [P] [US1] Create User model in src/models/user.py`
- ✅ 正确：`- [ ] T014 [US1] Implement UserService in src/services/user_service.py`
- ❌ 错误：`- [ ] Create User model`（缺少 ID 与故事标签）
- ❌ 错误：`T001 [US1] Create model`（缺少复选框）
- ❌ 错误：`- [ ] [US1] Create User model`（缺少任务 ID）
- ❌ 错误：`- [ ] T001 [US1] Create model`（缺少文件路径）

### 任务组织

1. **来源于用户故事（spec.md）——主组织维度**：
   - 每个用户故事（P1、P2、P3…）单独成 Phase
   - 将相关组件映射到所属的用户故事：
     - 该故事所需的模型
     - 该故事所需的服务
     - 该故事所需的端点 / UI
     - 若请求测试：该故事的测试任务
   - 标注故事间的依赖关系（大多数故事应保持相互独立）

2. **来源于 Contracts**：
   - 将每个 contract / endpoint → 映射到其服务的用户故事
   - 若请求测试：每个 contract → 在该故事的 Phase 中，于实现前添加可并行的“contract 测试任务”

3. **来源于数据模型**：
   - 将每个实体映射到需要它的用户故事（们）
   - 若某实体服务多个故事：将其放入最早的相关故事或 Setup 阶段
   - 关系 → 在相应故事阶段的服务层任务中体现

4. **来源于 Setup / 基础设施**：
   - 共享基础设施 → Setup 阶段（Phase 1）
   - 统一的阻塞前置项 → Foundational 阶段（Phase 2）
   - 故事专属的初始化工作 → 归入该故事的阶段

### 阶段结构

- **Phase 1**：Setup（项目初始化）
- **Phase 2**：Foundational（阻塞前置项——必须先于用户故事完成）
- **Phase 3+**：按优先级顺序的用户故事（P1、P2、P3…）
  - 在每个故事内：测试（如请求）→ 模型 → 服务 → 端点 → 集成
  - 每个阶段应构成一个可独立测试的完整增量
- **Final Phase**：Polish & Cross-Cutting Concerns（打磨与跨切面关注点）
