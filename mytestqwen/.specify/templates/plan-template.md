# Implementation Plan: [FEATURE]

**Branch**: `[###-feature-name]` | **Date**: [DATE] | **Spec**: [link]
**Input**: 来自 `/specs/[###-feature-name]/spec.md` 的 Feature specification

**Note**: 此模板由 `/taskkit.plan` 命令填充。执行流程见 `.specify/templates/commands/plan.md`。

## Summary

[从 feature spec 中提取：主要需求 + 来自 research 的技术方案]

## Technical Context

<!--
  操作说明：请将本节内容替换为项目的技术细节。
  此结构仅为建议，用于指导迭代过程。
-->

**Language/Version**: [例如 Python 3.11, Swift 5.9, Rust 1.75 或 NEEDS CLARIFICATION]  
**Primary Dependencies**: [例如 FastAPI, UIKit, LLVM 或 NEEDS CLARIFICATION]  
**Storage**: [如适用，例如 PostgreSQL, CoreData, files 或 N/A]  
**Testing**: [例如 pytest, XCTest, cargo test 或 NEEDS CLARIFICATION]  
**Target Platform**: [例如 Linux server, iOS 15+, WASM 或 NEEDS CLARIFICATION]
**Project Type**: [single/web/mobile - 决定源码结构]  
**Performance Goals**: [领域相关，例如 1000 req/s, 10k lines/sec, 60 fps 或 NEEDS CLARIFICATION]  
**Constraints**: [领域相关，例如 <200ms p95, <100MB memory, offline-capable 或 NEEDS CLARIFICATION]  
**Scale/Scope**: [领域相关，例如 10k users, 1M LOC, 50 screens 或 NEEDS CLARIFICATION]

## Constitution Check

GATE：必须在 Phase 0 research 前通过。Phase 1 design 后复检。

[基于 constitution 文件确定的 Gates]

## Project Structure

### Documentation (this feature)

```
specs/[###-feature]/
├── plan.md              # 本文件（/taskkit.plan 命令输出）
├── research.md          # Phase 0 输出（/taskkit.plan 命令）
├── data-model.md        # Phase 1 输出（/taskkit.plan 命令）
├── quickstart.md        # Phase 1 输出（/taskkit.plan 命令）
├── contracts/           # Phase 1 输出（/taskkit.plan 命令）
└── tasks.md             # Phase 2 输出（/taskkit.tasks 命令 - 非 /taskkit.plan 生成）
```

### Source Code (repository root)
<!--
  操作说明：请将下方占位树替换为该 feature 的具体布局。
  删除未使用的选项，并用真实路径扩展所选结构（例如 apps/admin, packages/something）。
  交付的 plan 不得包含“Option”标签。
-->

```
# [REMOVE IF UNUSED] Option 1: Single project (DEFAULT)
src/
├── models/
├── services/
├── cli/
└── lib/

tests/
├── contract/
├── integration/
└── unit/

# [REMOVE IF UNUSED] Option 2: Web application (when "frontend" + "backend" detected)
backend/
├── src/
│   ├── models/
│   ├── services/
│   └── api/
└── tests/

frontend/
├── src/
│   ├── components/
│   ├── pages/
│   └── services/
└── tests/

# [REMOVE IF UNUSED] Option 3: Mobile + API (when "iOS/Android" detected)
api/
└── [same as backend above]

ios/ or android/
└── [platform-specific structure: feature modules, UI flows, platform tests]
```

**Structure Decision**: [记录所选结构，并引用上方捕获的真实目录]

## Complexity Tracking

仅当 Constitution Check 存在需说明的违反项时填写

| Violation | Why Needed | Simpler Alternative Rejected Because |
|-----------|------------|-------------------------------------|
| [例如，第 4 个 project] | [当前需求] | [为何 3 个 projects 不足] |
| [例如，Repository pattern] | [具体问题] | [为何直接 DB 访问不足] |