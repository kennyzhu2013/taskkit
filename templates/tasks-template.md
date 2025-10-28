---
description: "用于功能实现的任务清单模板"
---

# Tasks: [FEATURE NAME]

**Input**: 来自 `/specs/[###-feature-name]/` 的设计文档
**Prerequisites**: plan.md（必需）、spec.md（用户故事必需）、research.md、data-model.md、contracts/

**Tests**: 下方示例包含测试任务。Tests 为可选——只有在 feature specification 明确要求时才包含。

**Organization**: Tasks 按 User story 分组，以便每个故事可独立实现与测试。

## Format: `[ID] [P?] [Story] Description`
- **[P]**: 可并行（不同文件、无依赖）
- **[Story]**: 任务所属的用户故事（如 US1、US2、US3）
- 在描述中包含精确的文件路径

## Path Conventions
- **Single project**: 仓库根目录下使用 `src/`、`tests/`
- **Web app**: 使用 `backend/src/`、`frontend/src/`
- **Mobile**: 使用 `api/src/`、`ios/src/` 或 `android/src/`
- 下方路径示例基于单项目结构——请根据 plan.md 的结构调整

<!-- 
  ============================================================================
  重要说明：下方任务仅为示例，用于说明格式。
  
  `/taskkit.tasks` 命令必须依据以下内容生成实际任务：
  - 来自 spec.md 的 User stories（及其优先级 P1、P2、P3...）
  - 来自 plan.md 的 Feature requirements
  - 来自 data-model.md 的 Entities
  - 来自 contracts/ 的 Endpoints
  
  任务必须根据 User story 组织，使每个故事均可：
  - 独立实施
  - 独立测试
  - 作为 MVP 增量交付
  
  请勿在生成的 tasks.md 文件中保留这些示例任务。
  ============================================================================
-->

## Phase 1: Setup (Shared Infrastructure)

**Purpose**: 项目初始化与基础结构

- [ ] T001 按实施计划创建项目结构
- [ ] T002 使用 [framework] 初始化 [language] 项目并添加依赖
- [ ] T003 [P] 配置代码规范与格式化工具

---

## Phase 2: Foundational (Blocking Prerequisites)

**Purpose**: 所有 User story 开始前必须完成的核心基础设施

**⚠️ CRITICAL**: 在本阶段完成之前，禁止开始任何 User story 工作

基础任务示例（请根据项目调整）：

- [ ] T004 设置数据库 schema 与 migrations 框架
- [ ] T005 [P] 实现 authentication/authorization 框架
- [ ] T006 [P] 搭建 API routing 与 middleware 结构
- [ ] T007 创建所有故事依赖的基础 models/entities
- [ ] T008 配置 error handling 与 logging 基础设施
- [ ] T009 设置环境配置管理

**Checkpoint**: 基础完备——可以开始并行实施 User stories

---

## Phase 3: User Story 1 - [Title] (Priority: P1) 🎯 MVP

**Goal**: [该故事将交付的内容简述]

**Independent Test**: [如何在独立状态下验证该故事]

### Tests for User Story 1 (OPTIONAL - only if tests requested) ⚠️

**NOTE: 请先编写这些 Tests，并确保在实现前它们 FAIL**

- [ ] T010 [P] [US1] 为 [endpoint] 编写 contract test 于 `tests/contract/test_[name].py`
- [ ] T011 [P] [US1] 为 [user journey] 编写 integration test 于 `tests/integration/test_[name].py`

### Implementation for User Story 1

- [ ] T012 [P] [US1] 在 `src/models/[entity1].py` 创建 [Entity1] model
- [ ] T013 [P] [US1] 在 `src/models/[entity2].py` 创建 [Entity2] model
- [ ] T014 [US1] 在 `src/services/[service].py` 实现 [Service]（依赖 T012、T013）
- [ ] T015 [US1] 在 `src/[location]/[file].py` 实现 [endpoint/feature]
- [ ] T016 [US1] 添加 validation 与 error handling
- [ ] T017 [US1] 为 User story 1 的操作添加 logging

**Checkpoint**: 此时，User Story 1 应可完全独立运行并测试

---

## Phase 4: User Story 2 - [Title] (Priority: P2)

**Goal**: [该故事将交付的内容简述]

**Independent Test**: [如何在独立状态下验证该故事]

### Tests for User Story 2 (OPTIONAL - only if tests requested) ⚠️

- [ ] T018 [P] [US2] 为 [endpoint] 编写 contract test 于 `tests/contract/test_[name].py`
- [ ] T019 [P] [US2] 为 [user journey] 编写 integration test 于 `tests/integration/test_[name].py`

### Implementation for User Story 2

- [ ] T020 [P] [US2] 在 `src/models/[entity].py` 创建 [Entity] model
- [ ] T021 [US2] 在 `src/services/[service].py` 实现 [Service]
- [ ] T022 [US2] 在 `src/[location]/[file].py` 实现 [endpoint/feature]
- [ ] T023 [US2] 如需，与 User Story 1 组件集成

**Checkpoint**: 此时，User Story 1 与 2 应均可独立运行

---

## Phase 5: User Story 3 - [Title] (Priority: P3)

**Goal**: [该故事将交付的内容简述]

**Independent Test**: [如何在独立状态下验证该故事]

### Tests for User Story 3 (OPTIONAL - only if tests requested) ⚠️

- [ ] T024 [P] [US3] 为 [endpoint] 编写 contract test 于 `tests/contract/test_[name].py`
- [ ] T025 [P] [US3] 为 [user journey] 编写 integration test 于 `tests/integration/test_[name].py`

### Implementation for User Story 3

- [ ] T026 [P] [US3] 在 `src/models/[entity].py` 创建 [Entity] model
- [ ] T027 [US3] 在 `src/services/[service].py` 实现 [Service]
- [ ] T028 [US3] 在 `src/[location]/[file].py` 实现 [endpoint/feature]

**Checkpoint**: 所有 User stories 现均应可独立运行

---

[可根据需要添加更多 User story 阶段，遵循相同模式]

---

## Phase N: Polish & Cross-Cutting Concerns

**Purpose**: 影响多个 User stories 的改进

- [ ] TXXX [P] 在 `docs/` 中更新文档
- [ ] TXXX 代码清理与重构
- [ ] TXXX 全局性能优化
- [ ] TXXX [P] 追加 unit tests（如有请求）于 `tests/unit/`
- [ ] TXXX 安全加固
- [ ] TXXX 运行 `quickstart.md` 验证

---

## Dependencies & Execution Order

### Phase Dependencies

- **Setup (Phase 1)**: 无依赖——可立即开始
- **Foundational (Phase 2)**: 依赖 Setup 完成——阻塞所有 User stories
- **User Stories (Phase 3+)**: 均依赖 Foundational 完成
  - 然后可并行推进（视团队容量）
  - 或按优先级顺序依次推进（P1 → P2 → P3）
- **Polish (Final Phase)**: 依赖所有目标 User stories 完成

### User Story Dependencies

- **User Story 1 (P1)**: Foundational（Phase 2）完成后即可开始——不依赖其他 stories
- **User Story 2 (P2)**: Foundational（Phase 2）完成后即可开始——可能与 US1 集成，但应可独立测试
- **User Story 3 (P3)**: Foundational（Phase 2）完成后即可开始——可能与 US1/US2 集成，但应可独立测试

### Within Each User Story

- 如包含 Tests，必须先编写并确保在实现前 FAIL
- 先 models，再 services
- 先 services，再 endpoints
- 先核心实现，再集成
- 完成当前 story 再进入下一个优先级

### Parallel Opportunities

- 所有标记为 [P] 的 Setup 任务可并行
- Foundational 阶段内所有标记为 [P] 的任务可并行
- Foundational 完成后，所有 User stories 可并行启动（视团队容量）
- 某个 User story 中所有标记为 [P] 的 Tests 可并行
- 同一 story 内标记为 [P] 的 models 可并行
- 不同 User stories 可由不同成员并行推进

---

## Parallel Example: User Story 1

```bash
# 如请求 Tests，可一起启动 User Story 1 的所有 Tests：
Task: "Contract test for [endpoint] in tests/contract/test_[name].py"
Task: "Integration test for [user journey] in tests/integration/test_[name].py"

# 一起启动 User Story 1 的所有 models：
Task: "Create [Entity1] model in src/models/[entity1].py"
Task: "Create [Entity2] model in src/models/[entity2].py"
```

---

## Implementation Strategy

### MVP First (User Story 1 Only)

1. 完成 Phase 1: Setup
2. 完成 Phase 2: Foundational（关键——阻塞所有 stories）
3. 完成 Phase 3: User Story 1
4. 停止并验证：独立测试 User Story 1
5. 如已准备好，部署/演示

### Incremental Delivery

1. 完成 Setup + Foundational → 基础就绪
2. 添加 User Story 1 → 独立测试 → 部署/演示（MVP!）
3. 添加 User Story 2 → 独立测试 → 部署/演示
4. 添加 User Story 3 → 独立测试 → 部署/演示
5. 每个故事在不破坏之前故事的情况下增添价值

### Parallel Team Strategy

多人协作：

1. 团队共同完成 Setup + Foundational
2. Foundational 完成后：
   - 开发者 A：User Story 1
   - 开发者 B：User Story 2
   - 开发者 C：User Story 3
3. 各故事独立完成并集成

---

## Notes

- 标记为 [P] 的任务 = 不同文件、无依赖
- [Story] 标签将任务映射到具体 User story，便于追踪
- 每个 User story 应可独立完成与测试
- 在实现前，先验证 Tests 失败
- 每个任务或逻辑组完成后进行一次提交
- 在任一检查点暂停以独立验证当前 story
- 避免：含糊任务、同一文件冲突、破坏独立性的跨故事依赖