


          
**目标说明**
- 为 Go 和 Java 提供定制化的规范、计划、任务与检查清单骨架。
- 将执行动作映射到你当前脚本目录结构：`scripts/powershell` 与 `scripts/bash`。
- 命令同时给出 PowerShell 与 Bash 两种格式，按你的脚本类型选择使用。

**Go 骨架**
- 规范（Spec）
  - 背景与目标：业务背景、用户场景、成功指标。
  - 范围与边界：包含/不包含的能力，外部系统接口。
  - 架构概览：分层结构（`cmd/`, `internal/`, `pkg/`），依赖关系，关键数据流。
  - 模块与接口：核心模块、公共接口、DTO/Model、错误处理策略。
  - API 设计：端点、方法、请求/响应、状态码、鉴权与速率限制。
  - 数据与存储：表结构/索引、迁移策略、事务/一致性（NEEDS CLARIFICATION 如未定）。
  - 非功能需求：性能目标（P95/P99）、可用性、容错、观测（日志/指标/追踪）。
  - 安全合规：鉴权、加密、密钥管理策略。
  - 依赖与版本：Go 版本、第三方库、工具链（`go`, `golangci-lint`, `gofmt`, `docker` 可选）。
  - 开放问题：需澄清的技术栈或业务细节（保留 NEEDS CLARIFICATION 标记）。
- 计划（Plan）
  - 里程碑：M1 原型/接口对齐，M2 核心功能，M3 测试与优化，M4 发布。
  - 风险与缓解：性能、依赖版本、并发数据竞争、跨平台行为。
  - 测试策略：`go test ./...`，`go vet ./...`，`golangci-lint run` 或 NEEDS CLARIFICATION。
  - 发布策略：构建产物、配置/密钥策略、回滚方案、灰度与监控。
  - 验收标准：功能验收、指标达成、故障演练、文档完整度。
- 任务（Tasks）
  - 任务条目格式：目标、步骤、依赖、脚本映射、验收标准、负责人。
  - 示例条目：
    - 目标：实现用户服务 API
    - 步骤：定义接口/模型 → 编写 handler → 接入存储 → 单元测试
    - 依赖：数据库初始化、配置就绪
    - 脚本（选其一）：
      - PowerShell: `pwsh -NoProfile -ExecutionPolicy Bypass -File scripts\powershell\build.ps1`
      - Bash: `bash ./scripts/bash/build.sh`
    - 验收：接口通过集成测试、`go test` 通过、`lint` 无严重问题
- 检查清单（Checklist）
  - 代码质量：`go vet`, `golangci-lint`, `gofmt`。
  - 安全：无明文密钥、遵循最小权限、TLS 与证书校验。
  - 文档：README、API 文档、运维手册、变更日志。
  - 发布门禁：构建产物校验、回滚验证、监控告警配置。

- 脚本映射（Go）
  - 构建
    - `ps`: `pwsh -NoProfile -ExecutionPolicy Bypass -File scripts\powershell\build.ps1`
    - `sh`: `bash ./scripts/bash/build.sh`
  - 测试
    - `ps`: `pwsh -NoProfile -ExecutionPolicy Bypass -File scripts\powershell\test.ps1`
    - `sh`: `bash ./scripts/bash/test.sh`
  - 运行
    - `ps`: `pwsh -NoProfile -ExecutionPolicy Bypass -File scripts\powershell\run.ps1`
    - `sh`: `bash ./scripts/bash/run.sh`
  - Lint/格式化（可选）
    - `ps`: `pwsh -NoProfile -ExecutionPolicy Bypass -File scripts\powershell\lint.ps1`
    - `sh`: `bash ./scripts/bash/lint.sh`

**Java 骨架**
- 规范（Spec）
  - 背景与目标：业务目标与成功指标，平台与环境（JDK 版本、框架）。
  - 范围与边界：微服务/单体，前后端接口范围，外部系统。
  - 架构概览：分层（Controller/Service/Repository）、依赖管理（Maven/Gradle）、打包与运行。
  - 模块与接口：核心服务、DTO/实体、异常与事务处理策略。
  - API 设计：Rest/GraphQL，认证（JWT/OAuth2）、错误码与幂等。
  - 数据与存储：ORM（JPA/MyBatis）、迁移（Flyway/Liquibase）、索引策略。
  - 非功能需求：性能目标、线程池与资源管理、稳定性与观测。
  - 安全合规：依赖漏洞扫描（OWASP）、输入校验、密钥管理。
  - 依赖与版本：JDK/框架版本、构建插件、测试库（JUnit/Mockito）。
  - 开放问题：打包格式（jar/war/docker）、部署环境（K8s/VM）等 NEEDS CLARIFICATION。
- 计划（Plan）
  - 里程碑：M1 项目骨架与依赖，M2 功能与DAO层，M3 集成测试与性能优化，M4 发布与监控。
  - 风险与缓解：依赖冲突、GC 性能、数据库连接池、事务边界。
  - 测试策略：单元测试（JUnit/Mockito）、集成测试（Testcontainers 可选）、端到端。
  - 发布策略：`mvn package`/`gradle build`，Docker 镜像，环境配置与密钥注入。
  - 验收标准：覆盖率阈值、端点契约测试、性能指标达标、文档完整。
- 任务（Tasks）
  - 示例条目：
    - 目标：订单服务实现与集成测试
    - 步骤：定义实体/仓储 → 服务实现 → 控制器 → JUnit/Mockito 测试 → 打包
    - 依赖：数据库就绪、配置文件
    - 脚本（选其一）：
      - PowerShell: `pwsh -NoProfile -ExecutionPolicy Bypass -File scripts\powershell\build.ps1`
      - Bash: `bash ./scripts/bash/build.sh`
    - 验收：`mvn test`/`gradle test` 通过，集成测试通过，打包产物可运行
- 检查清单（Checklist）
  - 代码质量：`spotbugs`/`checkstyle`/`pmd`（可在构建脚本中集成）。
  - 安全：依赖漏洞扫描（如 `owasp-dependency-check`），密钥不入库。
  - 文档：API 文档（Swagger/OpenAPI）、运行与运维手册。
  - 发布门禁：构建成功、镜像扫描、回滚与健康检查。

- 脚本映射（Java）
  - 构建
    - `ps`: `pwsh -NoProfile -ExecutionPolicy Bypass -File scripts\powershell\build.ps1`（内部调用 `mvn package` 或 `gradle build`）
    - `sh`: `bash ./scripts/bash/build.sh`
  - 测试
    - `ps`: `pwsh -NoProfile -ExecutionPolicy Bypass -File scripts\powershell\test.ps1`（内部调用 `mvn test` 或 `gradle test`）
    - `sh`: `bash ./scripts/bash/test.sh`
  - 运行
    - `ps`: `pwsh -NoProfile -ExecutionPolicy Bypass -File scripts\powershell\run.ps1`（例如 `java -jar target/app.jar`）
    - `sh`: `bash ./scripts/bash/run.sh`
  - 质量检查（可选）
    - `ps`: `pwsh -NoProfile -ExecutionPolicy Bypass -File scripts\powershell\lint.ps1`（内部调用 `spotbugs/checkstyle/pmd`）
    - `sh`: `bash ./scripts/bash/lint.sh`

**与现有模板的结合**
- 将上述骨架用于填充：
  - 规范 → `templates/spec-template.md`
  - 计划 → `templates/plan-template.md`（把 Testing 字段替换为实际栈：Go `go test`, Java `JUnit`，并移除 NEEDS CLARIFICATION）
  - 任务 → `templates/tasks-template.md`
  - 检查清单 → `templates/checklist-template.md`
- 若你希望区分语言版本，可新增：
  - `templates/spec-template-go.md` / `templates/spec-template-java.md`
  - `templates/plan-template-go.md` / `templates/plan-template-java.md`
  - `templates/tasks-template-go.md` / `templates/tasks-template-java.md`

**脚本命名建议**
- PowerShell：`build.ps1`, `test.ps1`, `run.ps1`, `lint.ps1` 存放于 `scripts\powershell\`
- Bash：`build.sh`, `test.sh`, `run.sh`, `lint.sh` 存放于 `scripts/bash/`
- 在 `/implement` 或你的执行流程中直接调用上述路径，与你当前“Script Commands”面板一致。

**小提示**
- 若目前脚本文件还未创建，可以先以上述命名建立空脚本并逐步充实内容。
- 在计划模板中保留“开放问题”段落，先用 `/clarify` 完成澄清，再替换计划模板中的占位标识。
- 有需要我可以按你的具体框架（比如 Go: gin/fiber/echo；Java: Spring Boot/Quarkus/Micronaut）进一步细化骨架章节与示例命令。
        