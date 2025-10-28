# Feature Specification: [FEATURE NAME]

**Feature Branch**: `[###-feature-name]`  
**Created**: [DATE]  
**Status**: Draft  
**Input**: 用户描述: "$ARGUMENTS"

## User Scenarios & Testing *(mandatory)*

<!--
  重要：User stories 应按重要性排序为用户旅程。
  每个 User story/用户旅程必须可独立测试——即仅实现其中一个，
  也应形成可为用户提供价值的 MVP（Minimum Viable Product）。
  
  请为每个故事标注优先级（P1、P2、P3 等），其中 P1 最关键。
  将每个故事视为一个独立的功能切片，能够：
  - 独立开发
  - 独立测试
  - 独立部署
  - 独立向用户演示
-->

### User Story 1 - [Brief Title] (Priority: P1)

[用通俗语言描述该用户旅程]

**Why this priority**: [解释其价值及为何具有该优先级]

**Independent Test**: [说明如何独立测试，例如：“可通过 [具体动作] 完整测试并交付 [具体价值]”]

**Acceptance Scenarios**:

1. **Given** [初始状态], **When** [动作], **Then** [预期结果]
2. **Given** [初始状态], **When** [动作], **Then** [预期结果]

---

### User Story 2 - [Brief Title] (Priority: P2)

[用通俗语言描述该用户旅程]

**Why this priority**: [解释其价值及为何具有该优先级]

**Independent Test**: [说明如何独立测试]

**Acceptance Scenarios**:

1. **Given** [初始状态], **When** [动作], **Then** [预期结果]

---

### User Story 3 - [Brief Title] (Priority: P3)

[用通俗语言描述该用户旅程]

**Why this priority**: [解释其价值及为何具有该优先级]

**Independent Test**: [说明如何独立测试]

**Acceptance Scenarios**:

1. **Given** [初始状态], **When** [动作], **Then** [预期结果]

---

[根据需要添加更多 User stories，并为每个故事分配优先级]

### Edge Cases

<!--
  操作说明：本节内容为占位符。
  请用实际的边界情况填充。
-->

- 当 [边界条件] 时会发生什么？
- 系统如何处理 [错误场景]？

## Requirements *(mandatory)*

<!--
  操作说明：本节内容为占位符。
  请用实际的功能需求填充。
-->

### Functional Requirements

- **FR-001**: 系统必须 [具体能力，例如“允许用户创建账户”]
- **FR-002**: 系统必须 [具体能力，例如“校验电子邮箱地址”]  
- **FR-003**: 用户必须能够 [关键交互，例如“重置密码”]
- **FR-004**: 系统必须 [数据要求，例如“持久化用户偏好”]
- **FR-005**: 系统必须 [行为要求，例如“记录所有安全事件”]

*标注不清需求的示例：*

- **FR-006**: 系统必须通过 [NEEDS CLARIFICATION：认证方式未指定 - email/password、SSO、OAuth？] 进行身份认证
- **FR-007**: 系统必须保留用户数据 [NEEDS CLARIFICATION：保留周期未指定]

### Key Entities *(include if feature involves data)*

- **[Entity 1]**: [其代表的含义、关键属性（不涉及实现）]
- **[Entity 2]**: [其代表的含义、与其他实体的关系]

## Success Criteria *(mandatory)*

<!--
  操作说明：定义可度量的成功标准。
  这些标准必须与技术无关且可度量。
-->

### Measurable Outcomes

- **SC-001**: [可度量指标，例如“用户在 2 分钟内完成账户创建”]
- **SC-002**: [可度量指标，例如“系统在 1000 并发下无性能下降”]
- **SC-003**: [用户满意度指标，例如“90% 的用户首次成功完成主要操作”]
- **SC-004**: [业务指标，例如“将与 [X] 相关的支持工单减少 50%”]