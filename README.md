# slam-cli

`slam-cli` 是一个面向 **SLAM 场景 / 平台能力** 的命令行工具骨架项目。

当前阶段目标不是一次性把所有业务能力做完，而是先基于 [`lark-cli`](https://github.com/larksuite/cli) 的整体设计思路，搭出一个 **可持续扩展、适合后续模块化演进** 的 CLI 基础架构，优先完成项目说明、目录规范和最小可运行骨架。

---

## 1. 项目目标

`slam-cli` 的目标是把未来分散的 SLAM 操作能力统一收敛到一个命令行入口中，支持：

- 面向开发者的本地命令执行
- 面向自动化脚本 / Agent 的稳定调用
- 面向后续平台能力扩展的模块化接入

参考 `lark-cli` 后，本项目优先借鉴以下思想：

1. **统一入口**：所有能力收口到一个二进制中。
2. **分层命令体系**：高频快捷命令、标准领域命令、原始接口命令并存。
3. **领域模块化**：按能力域拆分命令，避免所有逻辑堆在根命令下。
4. **基础设施抽离**：配置、输出、运行时上下文、接口访问逻辑独立沉淀。
5. **对人和自动化都友好**：既能交互使用，也便于脚本化编排。

---

## 2. 当前阶段范围

当前仓库先完成以下内容：

- 明确 `slam-cli` 的架构方向
- 建立最小可运行的 Go CLI 骨架
- 预留后续命令分层与模块扩展位置
- 在 README 中沉淀开发约定与演进路线

> 也就是说：**先把骨架搭对，再逐步填业务能力。**

---

## 3. 参考 lark-cli 提炼出的架构原则

虽然 `slam-cli` 的业务范围会比 `lark-cli` 更聚焦，但它的架构方式非常值得借鉴。

### 3.1 三层命令模型

`slam-cli` 建议采用下面的命令分层：

#### A. 快捷命令层（Shortcut Layer）
面向高频、短路径、结果导向的操作。

示例：

- `slam +doctor`
- `slam +login`
- `slam +init`

特点：

- 参数尽量少
- 默认行为尽量智能
- 输出更适合人直接阅读

#### B. 领域命令层（Domain Command Layer）
按业务域拆分标准命令，适合长期维护。

示例：

- `slam auth login`
- `slam config set`
- `slam workflow run`
- `slam project list`

特点：

- 结构清晰
- 参数语义稳定
- 更适合文档化和测试

#### C. 原始接口层（Raw API Layer）
作为能力兜底层，保留最细粒度的调用入口。

示例：

- `slam api call ...`

特点：

- 覆盖率最高
- 方便调试和快速接入新能力
- 不阻塞上层命令演进

---

## 4. slam-cli 的基础目录设计

结合 Go CLI 常见实践，以及 `lark-cli` 的“统一入口 + 模块化扩展”思路，当前建议目录结构如下：

```text
slam-cli/
├── cmd/
│   └── slam/
│       └── main.go              # 程序入口
├── internal/
│   ├── app/
│   │   ├── root.go              # 根命令装配
│   │   └── version.go           # 版本信息
│   ├── command/
│   │   ├── auth/
│   │   │   └── auth.go          # 认证相关命令
│   │   ├── config/
│   │   │   └── config.go        # 配置相关命令
│   │   ├── shortcut/
│   │   │   └── shortcut.go      # 快捷命令层
│   │   ├── workflow/
│   │   │   └── workflow.go      # 工作流/任务编排命令
│   │   └── api/
│   │       └── api.go           # 原始接口调用层
│   ├── config/
│   │   └── config.go            # 配置模型与默认值
│   └── output/
│       └── output.go            # 输出封装（plain/json 预留）
├── go.mod
├── go.sum
└── README.md
```

### 为什么这样拆？

- `cmd/slam`：只保留启动逻辑，避免入口过重。
- `internal/app`：负责命令树装配，是整个 CLI 的应用层。
- `internal/command/*`：每个命令域各自独立，便于后续继续拆分子命令。
- `internal/config`：后续可以接入本地配置文件、环境变量、远程配置。
- `internal/output`：统一人类可读输出与机器可解析输出。

---

## 5. 当前建议优先落地的模块

在业务能力尚未完全明确前，建议先落地这几类通用模块：

### 5.1 auth
负责认证、身份、token 管理等。

典型命令：

- `slam auth login`
- `slam auth logout`
- `slam auth status`

### 5.2 config
负责本地配置的读写和查看。

典型命令：

- `slam config get`
- `slam config set`
- `slam config list`

### 5.3 shortcut
负责高频场景封装。

典型命令：

- `slam +doctor`
- `slam +init`

### 5.4 workflow
负责较长链路任务编排。

典型命令：

- `slam workflow run`
- `slam workflow inspect`

### 5.5 api
负责调用底层原始接口，作为扩展兜底层。

典型命令：

- `slam api call`

---

## 6. 最小开发约定

为了保证后续扩展时不返工，建议从一开始就遵循下面的约定：

### 6.1 命令职责单一
- 一个目录只负责一个命令域。
- 根命令只做装配，不写具体业务逻辑。

### 6.2 输出统一
- 后续所有命令都尽量通过统一输出层返回结果。
- 先支持纯文本，后续扩展为 JSON / Table。

### 6.3 配置集中
- 本地配置读取逻辑统一放在 `internal/config`。
- 不在各个命令中重复拼接配置路径。

### 6.4 先骨架，后业务
- 先保证目录结构和命令树稳定。
- 业务细节逐步填充，不急于一次性做深。

---

## 7. 当前阶段的实现策略

本仓库建议按以下顺序推进：

1. **README 先行**：把项目定位、结构和路线讲清楚。
2. **最小命令树可运行**：先能执行 `slam --help`。
3. **基础模块占位**：把 `auth / config / shortcut / workflow / api` 这些目录先建立起来。
4. **逐步填充命令能力**：优先做通用能力，再做业务命令。

---

## 8. 本地运行

如果已经完成最小骨架，可以这样运行：

```bash
go run ./cmd/slam --help
```

后续示例：

```bash
go run ./cmd/slam version
go run ./cmd/slam +doctor
go run ./cmd/slam auth status
```

---

## 9. 近期 Roadmap

### Phase 1：CLI 基础骨架
- [x] README 架构说明
- [x] 根命令入口
- [x] 最小模块目录
- [ ] 基础输出封装
- [ ] 配置文件读写

### Phase 2：通用能力建设
- [ ] auth 登录态管理
- [ ] config 读写与校验
- [ ] version / env / doctor 命令

### Phase 3：业务能力接入
- [ ] 补充 SLAM 领域命令
- [ ] 接入 workflow / task / project 等能力
- [ ] 增加 JSON 输出与自动化调用支持

---

## 10. 总结

`slam-cli` 目前的重点不是“功能有多少”，而是“骨架是否足够稳”。

这次基础设计明确采用了从 `lark-cli` 提炼出来的几个核心思路：

- **统一 CLI 入口**
- **快捷命令 + 标准命令 + 原始接口 的三层结构**
- **按领域模块拆分命令目录**
- **配置 / 输出 / 运行时能力下沉为基础设施**

如果这个方向保持一致，后续无论接多少 SLAM 相关能力，都能在现有架构上平滑扩展。
