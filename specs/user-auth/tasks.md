# 用户认证功能任务拆分

> **Branch**: `feature/user-auth`  
> **Spec Path**: `specs/user-auth/`  
> **Created**: 2026-01-01  
> **Input**: spec.md, plan.md

---

## Task Format

```
[ID] [P?] [Story] Description
```

| 标记 | 含义 |
|------|------|
| `T001` | 任务 ID |
| `[P]` | 可并行执行（不同文件，无依赖） |
| `[US1]` | 关联 User Story 1 |

---

## Task Overview

| ID | Task | Story | Status | Parallel | Est. Lines |
|----|------|-------|--------|----------|------------|
| T001 | 项目基础设置检查 | Setup | ⏸️ | - | - |
| T002 | UUID工具实现 | Foundation | ⏸️ | [P] | 30 |
| T003 | 密码验证工具实现 | Foundation | ⏸️ | [P] | 40 |
| T004 | SMS服务接口定义 | Foundation | ⏸️ | [P] | 25 |
| T005 | SMS服务Redis实现 | Foundation | ⏸️ | - | 45 |
| T006 | JWT服务接口定义 | Foundation | ⏸️ | [P] | 30 |
| T007 | JWT服务实现 | Foundation | ⏸️ | - | 45 |
| T008 | 设备信息解析工具 | Foundation | ⏸️ | [P] | 35 |
| T009 | 用户表DDL定义 | Data Model | ⏸️ | [P] | 20 |
| T010 | 登录历史表DDL定义 | Data Model | ⏸️ | [P] | 20 |
| T011 | User Model接口定义 | Data Model | ⏸️ | - | 30 |
| T012 | User Model类型定义 | Data Model | ⏸️ | [P] | 25 |
| T013 | User Model常量错误 | Data Model | ⏸️ | [P] | 20 |
| T014 | User Model GORM实现 | Data Model | ⏸️ | - | 50 |
| T015 | LoginHistory Model接口定义 | Data Model | ⏸️ | [P] | 25 |
| T016 | LoginHistory Model类型定义 | Data Model | ⏸️ | [P] | 20 |
| T017 | LoginHistory Model常量错误 | Data Model | ⏸️ | [P] | 20 |
| T018 | LoginHistory Model GORM实现 | Data Model | ⏸️ | - | 45 |
| T019 | API文件定义 | API | ⏸️ | - | 50 |
| T020 | goctl生成Handler/Types | API | ⏸️ | - | - |
| T021 | 发送注册验证码Logic | US1 | ⏸️ | - | 45 |
| T022 | 用户注册Logic | US1 | ⏸️ | - | 50 |
| T023 | 用户登录Logic | US2 | ⏸️ | - | 50 |
| T024 | 发送密码重置验证码Logic | US3 | ⏸️ | - | 40 |
| T025 | 请求重置密码Logic | US3 | ⏸️ | - | 35 |
| T026 | 确认重置密码Logic | US3 | ⏸️ | - | 50 |
| T027 | 查询登录历史Logic | US4 | ⏸️ | - | 45 |
| T028 | ServiceContext配置 | Integration | ⏸️ | - | 30 |
| T029 | User Model单元测试 | Test | ⏸️ | [P] | 40 |
| T030 | LoginHistory Model单元测试 | Test | ⏸️ | [P] | 35 |
| T031 | Logic层单元测试 | Test | ⏸️ | - | 50 |
| T032 | 代码清理和格式化 | Polish | ⏸️ | - | - |

---

## Phase 1: Setup

**目的**: 项目初始化和基础配置

- [ ] T001 确认 Go-Zero 项目结构已就绪
- [ ] T001.1 确认 goctl 工具已安装
- [ ] T001.2 确认数据库连接配置已就绪
- [ ] T001.3 确认 Redis 连接配置已就绪

**Checkpoint**: ✅ 开发环境就绪

---

## Phase 2: Foundation (基础工具和服务)

**目的**: 实现基础工具和服务，为业务逻辑提供支撑

### Step 1: UUID工具

- [ ] T002 [P] 创建 `pkg/uuid/uuid.go`
  - 实现 `GenerateUUID()` 函数，使用 `github.com/google/uuid` 生成 UUID v7
  - 返回 string 格式的 UUID
  - **验收标准**: 函数可生成有效的 UUID v7 字符串

### Step 2: 密码验证工具

- [ ] T003 [P] 创建 `pkg/validator/password.go`
  - 实现 `ValidatePasswordStrength(password string) error`
  - 验证密码长度（8-32字符）
  - 验证密码复杂度（数字+大小写字母+特殊字符）
  - **验收标准**: 能正确验证密码强度，返回明确的错误信息

### Step 3: SMS服务

- [ ] T004 [P] 创建 `pkg/sms/sms.go`
  - 定义 `Service` 接口（SendCode, VerifyCode, CheckRateLimit）
  - **验收标准**: 接口定义完整，符合 plan.md 规范

- [ ] T005 创建 `pkg/sms/redis_store.go`
  - 实现 Redis 存储验证码
  - 实现发送频率限制（60秒）
  - 实现验证码验证（一次性使用）
  - **验收标准**: 验证码可正确存储、验证和过期清理

### Step 4: JWT服务

- [ ] T006 [P] 创建 `pkg/jwt/jwt.go` 和 `pkg/jwt/claims.go`
  - 定义 `Service` 接口（GenerateToken, VerifyToken, InvalidateUserTokens）
  - 定义 `Claims` 结构体（包含 UserID string）
  - **验收标准**: 接口定义完整，Claims 包含必要字段

- [ ] T007 实现 `pkg/jwt/jwt.go`
  - 实现 Token 生成（7天有效期）
  - 实现 Token 验证
  - 实现 Token 失效机制（Redis存储失效列表）
  - **验收标准**: 可生成和验证 JWT Token，密码重置后Token失效

### Step 5: 设备信息解析

- [ ] T008 [P] 创建 `pkg/device/parser.go`
  - 实现 `ParseDeviceInfo(userAgent string) (deviceType, deviceID string)`
  - 解析设备类型（Web/Android/iOS）
  - 生成设备标识（基于User-Agent哈希）
  - **验收标准**: 能正确解析User-Agent并生成设备信息

**Checkpoint**: ✅ 基础工具和服务就绪，可开始数据模型实现

---

## Phase 3: Data Model (数据模型)

**目的**: 实现数据访问层，包括DDL和Model层

### Step 1: DDL定义

- [ ] T009 [P] 创建 `migrations/auth/users.sql`
  - 定义 users 表结构（UUID v7主键）
  - 包含所有必要字段和索引
  - **验收标准**: DDL可成功执行，表结构符合设计

- [ ] T010 [P] 创建 `migrations/auth/login_history.sql`
  - 定义 login_history 表结构（UUID v7主键）
  - 包含外键关联和索引
  - **验收标准**: DDL可成功执行，外键关系正确

### Step 2: User Model实现

- [ ] T011 创建 `model/auth/users/interface.go`
  - 定义 `Model` 接口（Insert, FindOne, FindByMobile, Update, UpdatePassword, UpdateStatus, WithTx, Trans）
  - **验收标准**: 接口定义完整，符合 plan.md 规范

- [ ] T012 [P] 创建 `model/auth/users/types.go`
  - 定义 `User` 结构体（UUID v7 ID）
  - 定义 TableName 方法
  - **验收标准**: 结构体字段与DDL一致，GORM标签正确

- [ ] T013 [P] 创建 `model/auth/users/vars.go`
  - 定义错误常量（ErrUserNotFound, ErrUserAlreadyExists等）
  - **验收标准**: 错误定义完整，符合业务规则

- [ ] T014 创建 `model/auth/users/factory.go` 和 `model/auth/users/gorm_dao.go`
  - 实现工厂函数
  - 实现 GORM DAO（所有接口方法）
  - 实现 UUID v7 主键生成
  - **验收标准**: 所有接口方法实现完整，支持事务

### Step 3: LoginHistory Model实现

- [ ] T015 [P] 创建 `model/auth/login_history/interface.go`
  - 定义 `Model` 接口（Insert, FindOne, FindByUserID, DeleteOldRecords, CountByUserID）
  - **验收标准**: 接口定义完整，符合 plan.md 规范

- [ ] T016 [P] 创建 `model/auth/login_history/types.go`
  - 定义 `LoginHistory` 结构体（UUID v7 ID）
  - 定义 TableName 方法
  - **验收标准**: 结构体字段与DDL一致

- [ ] T017 [P] 创建 `model/auth/login_history/vars.go`
  - 定义错误常量
  - **验收标准**: 错误定义完整

- [ ] T018 创建 `model/auth/login_history/factory.go` 和 `model/auth/login_history/gorm_dao.go`
  - 实现工厂函数
  - 实现 GORM DAO（所有接口方法）
  - 实现分页查询
  - **验收标准**: 所有接口方法实现完整，分页查询正确

**Checkpoint**: ✅ 数据模型就绪，可开始API定义

---

## Phase 4: API Definition (API定义)

**目的**: 定义API接口并生成代码框架

### Step 1: API文件定义

- [ ] T019 创建 `api/doc/auth/user_auth.api`
  - 定义所有请求/响应类型
  - 接口定义的group使用 `auth`（与 plan.md 中的 @server group 一致）
  - 定义7个API端点（发送注册验证码、注册、登录、发送重置验证码、请求重置、确认重置、查询登录历史）
  - **验收标准**: API定义完整，类型定义正确，符合 plan.md 规范

- [ ] T019.1 在入口文件中导入模块
  - 在 `api/doc/api.api` 中添加 `import "auth/user_auth.api"`
  - ⚠️ **注意**: 导入路径是相对于 `api.doc/` 目录的，文件路径为 `api/doc/auth/user_auth.api`，所以导入路径为 `auth/user_auth.api`
  - **验收标准**: 入口文件正确导入新模块，导入路径格式正确

### Step 2: 生成代码框架

- [ ] T020 运行 goctl 生成 Handler/Types
  ```bash
  goctl api go -api api/doc/api.api -dir api/ --style=go_zero --type-group
  ```
  - ⚠️ **重要**: 必须使用入口文件 `api.doc/api.api`，不能直接使用模块文件
  - ⚠️ **原因**: 使用模块文件会覆盖 `routes.go`，导致其他模块路由丢失
  - **验收标准**: Handler 和 Types 文件成功生成，无编译错误，routes.go 包含所有模块路由

**Checkpoint**: ✅ API框架就绪，可开始Logic层实现

---

## Phase 5: User Story 1 - 用户注册 (P1) 🎯 MVP

**目标**: 实现用户注册功能，包括发送验证码和注册接口

**独立测试**: 用户成功提交手机号、密码和短信验证码后，系统创建账户并返回用户ID

### Step 1: 发送注册验证码Logic

- [ ] T021 [US1] 实现 `api/internal/logic/auth/sendregistercode_logic.go`
  - 验证手机号格式
  - 检查发送频率限制（SMS服务）
  - 生成6位验证码
  - 发送验证码（SMS服务）
  - 返回成功响应
  - **验收标准**: 可成功发送验证码，频率限制生效，验证码存储在Redis

### Step 2: 用户注册Logic

- [ ] T022 [US1] 实现 `api/internal/logic/auth/register_logic.go`
  - 验证手机号格式
  - 验证验证码（SMS服务）
  - 验证密码强度（validator工具）
  - 检查手机号是否已注册（User Model）
  - 加密密码（bcrypt）
  - 生成UUID v7作为用户ID
  - 创建用户（User Model）
  - 返回用户ID
  - **验收标准**: 可成功注册用户，密码加密存储，返回UUID v7格式的用户ID

**Checkpoint**: ✅ User Story 1 可独立测试和验证

---

## Phase 6: User Story 2 - 用户登录 (P1) 🎯 MVP

**目标**: 实现用户登录功能，返回JWT Token并记录登录历史

**独立测试**: 用户使用正确的手机号和密码登录后，系统返回JWT Token

### Step 1: 用户登录Logic

- [ ] T023 [US2] 实现 `api/internal/logic/auth/login_logic.go`
  - 验证手机号格式
  - 查询用户（User Model）
  - 检查账户状态（启用/禁用/锁定）
  - 验证密码（bcrypt）
  - 生成JWT Token（JWT服务）
  - 更新最后登录时间（User Model）
  - 解析设备信息（device工具）
  - 记录登录历史（LoginHistory Model）
  - 返回Token
  - **验收标准**: 可成功登录，返回有效Token，登录历史正确记录

**Checkpoint**: ✅ User Story 2 可独立测试和验证

---

## Phase 7: User Story 3 - 忘记密码重置 (P2)

**目标**: 实现密码重置功能，包括发送验证码、请求重置和确认重置

**独立测试**: 用户提交手机号、短信验证码和新密码后，系统更新密码并允许登录

### Step 1: 发送密码重置验证码Logic

- [ ] T024 [US3] 实现 `api/internal/logic/auth/sendresetcode_logic.go`
  - 验证手机号格式
  - 检查用户是否存在（User Model）
  - 检查发送频率限制（SMS服务）
  - 生成并发送验证码（SMS服务）
  - 返回成功响应
  - **验收标准**: 可成功发送重置验证码，仅已注册用户可请求

### Step 2: 请求重置密码Logic

- [ ] T025 [US3] 实现 `api/internal/logic/auth/resetpasswordrequest_logic.go`
  - 验证手机号格式
  - 检查用户是否存在（User Model）
  - 发送验证码（复用T024逻辑或调用SMS服务）
  - 返回成功响应
  - **验收标准**: 可成功请求重置密码

### Step 3: 确认重置密码Logic

- [ ] T026 [US3] 实现 `api/internal/logic/auth/resetpasswordconfirm_logic.go`
  - 验证手机号格式
  - 验证验证码（SMS服务）
  - 验证密码强度（validator工具）
  - 查询用户（User Model）
  - 检查新密码是否与旧密码相同
  - 加密新密码（bcrypt）
  - 更新密码（User Model）
  - 使旧Token失效（JWT服务）
  - 返回成功响应
  - **验收标准**: 可成功重置密码，旧Token失效，新密码可登录

**Checkpoint**: ✅ User Story 3 可独立测试和验证

---

## Phase 8: User Story 4 - 登录历史记录 (P2)

**目标**: 实现登录历史查询功能

**独立测试**: 用户或管理员可以查询登录历史记录，包括登录时间、IP地址、设备信息等

### Step 1: 查询登录历史Logic

- [ ] T027 [US4] 实现 `api/internal/logic/auth/loginhistory_logic.go`
  - 从Token获取用户ID（JWT服务）
  - 验证权限（用户只能查询自己的历史）
  - 解析分页参数
  - 解析筛选条件（时间范围、IP地址）
  - 查询登录历史（LoginHistory Model）
  - 格式化返回数据
  - 返回分页结果
  - **验收标准**: 可成功查询登录历史，支持分页和筛选，权限控制正确

**Checkpoint**: ✅ User Story 4 可独立测试和验证

---

## Phase 9: Integration (集成配置)

**目的**: 配置ServiceContext，集成所有服务

- [ ] T028 更新 `api/internal/svc/servicecontext.go`
  - 添加 User Model
  - 添加 LoginHistory Model
  - 添加 SMS Service
  - 添加 JWT Service
  - **验收标准**: ServiceContext包含所有必要的服务，可正常初始化

**Checkpoint**: ✅ 服务集成完成，可进行端到端测试

---

## Phase 10: Testing (测试)

**目的**: 编写单元测试，确保代码质量

- [ ] T029 [P] 创建 `model/auth/users/*_test.go`
  - 测试所有Model方法
  - 测试事务功能
  - **验收标准**: 测试覆盖率 > 80%，所有测试通过

- [ ] T030 [P] 创建 `model/auth/login_history/*_test.go`
  - 测试所有Model方法
  - 测试分页查询
  - **验收标准**: 测试覆盖率 > 80%，所有测试通过

- [ ] T031 创建 `api/internal/logic/auth/*_test.go`
  - 测试所有Logic方法（Mock Model和Service）
  - 测试业务规则验证
  - **验收标准**: 测试覆盖率 > 80%，所有测试通过

**Checkpoint**: ✅ 测试覆盖充分，代码质量达标

---

## Phase 11: Polish (收尾工作)

**目的**: 代码清理和质量检查

- [ ] T032 代码清理和格式化
  - 运行 `go fmt ./...`
  - 运行 `golangci-lint run`
  - 补充中文注释
  - 检查函数行数（≤50行）
  - **验收标准**: 无lint错误，所有公开接口有中文注释

- [ ] T033 确认测试覆盖率
  - 运行 `go test -cover ./...`
  - 确认覆盖率 > 80%
  - **验收标准**: 测试覆盖率达标

- [ ] T034 编译检查
  - 运行 `go build ./...`
  - **验收标准**: 编译无错误

**Checkpoint**: ✅ 代码质量达标，可提交代码

---

## Dependencies

```
Phase 1 (Setup)
    ↓
Phase 2 (Foundation)
    ↓
Phase 3 (Data Model)
    ↓
Phase 4 (API Definition)
    ↓
Phase 5 (US1: 注册) → Phase 6 (US2: 登录) → Phase 7 (US3: 重置) → Phase 8 (US4: 历史)
    ↓
Phase 9 (Integration)
    ↓
Phase 10 (Testing)
    ↓
Phase 11 (Polish)
```

### 并行执行说明

- `[P]` 标记的任务可与同 Phase 内其他 `[P]` 任务并行
- Phase 2 中的基础工具可以并行开发
- Phase 3 中的两个Model可以并行开发
- Phase 10 中的测试可以并行编写

### 关键依赖

1. **Foundation → Data Model**: 基础工具必须在Model层之前完成
2. **Data Model → API**: Model层必须在API定义之前完成
3. **API → Logic**: API框架必须在Logic实现之前生成
4. **US1 → US2**: 注册功能必须在登录功能之前完成（需要用户数据）
5. **US2 → US4**: 登录功能必须在登录历史查询之前完成（需要登录记录）

---

## Notes

- 每个 Task 完成后提交代码
- 每个 Checkpoint 进行验证
- 遇到问题及时记录到 Open Questions
- 函数代码量必须 ≤ 50 行
- 所有公开接口必须有中文注释
- 使用统一的错误处理（pkg/errorx）
- 使用统一的响应格式（pkg/response）

---

## Revision History

| Version | Date | Author | Changes |
|---------|------|--------|---------|
| 1.0 | 2026-01-01 | - | 初始版本 |

