# slam-cli

`slam-cli` 是一个面向 **SLAM 软件包查询与认证** 场景的命令行工具。

当前项目已经具备以下可用能力：

- 飞书扫码登录并保存本地认证信息
- 查看当前登录状态
- 查看本地配置
- 按条件分页查询软件包列表
- 查询一个或多个软件包详情
- 提供基础快捷命令 `+doctor` 与 `+init`

当前版本：`0.1.0`，定义见 `internal/app/version.go:3`。

---

## 功能概览

根命令注册在 `internal/app/root.go:14`，当前 CLI 暴露的命令包括：

- `slam-cli help`
- `slam-cli version`
- `slam-cli auth login`
- `slam-cli auth status`
- `slam-cli config`
- `slam-cli config list`
- `slam-cli config get`
- `slam-cli package list`
- `slam-cli package info`
- `slam-cli +doctor`
- `slam-cli +init`

对应实现位置：

- 根命令装配：`internal/app/root.go:14`
- 认证命令：`internal/command/auth/auth.go:14`
- 配置命令：`internal/command/config/config.go:12`
- 软件包命令：`internal/command/package/package.go:17`
- 快捷命令：`internal/command/shortcut/shortcut.go:9`

> 说明：`internal/pkg/caller/caller.go:101` 中已经实现了 `GetAllPackages` 底层调用，但当前版本还没有对应的 CLI 子命令对外暴露。

---

## 项目结构

当前仓库核心结构如下：

```text
slam-cli/
├── cmd/
│   └── slam/
│       └── main.go
├── internal/
│   ├── app/
│   │   ├── root.go
│   │   └── version.go
│   ├── command/
│   │   ├── auth/
│   │   │   └── auth.go
│   │   ├── config/
│   │   │   └── config.go
│   │   ├── package/
│   │   │   └── package.go
│   │   └── shortcut/
│   │       └── shortcut.go
│   └── pkg/
│       ├── authflow/
│       │   ├── manager.go
│       │   └── qr.go
│       ├── caller/
│       │   ├── DTO/
│       │   │   ├── req.go
│       │   │   └── resp.go
│       │   ├── caller.go
│       │   └── common.go
│       ├── config/
│       │   └── config.go
│       └── output/
│           └── output.go
├── go.mod
├── go.sum
└── README.md
```

职责划分：

- `cmd/slam/main.go:9`：程序入口
- `internal/app`：根命令装配、帮助与版本信息
- `internal/command/*`：CLI 命令层
- `internal/pkg/authflow`：飞书扫码登录、JWT 获取与用户信息拉取
- `internal/pkg/caller`：SLAM 后端 API 调用与响应解析
- `internal/pkg/config`：本地配置与登录态管理
- `internal/pkg/output/output.go:8`：基础输出封装

---

## 快速开始

### 本地运行

```bash
go run ./cmd/slam help
go run ./cmd/slam version
```

### 构建二进制

```bash
go build -o slam-cli ./cmd/slam
./slam-cli help
```

### 运行环境

- Go 版本：`1.26.1`，定义见 `go.mod:3`
- 主要依赖：`cobra`、`resty`、`persistent-cookiejar`、`qrterminal`，见 `go.mod:5`

---

## 认证与本地配置

### 登录

登录命令定义在 `internal/command/auth/auth.go:23`。

```bash
slam-cli auth login
```

支持区域参数：

```bash
slam-cli auth login --region cn
slam-cli auth login --region us
slam-cli auth login --region sg
```

区域 flag 定义见 `internal/command/auth/auth.go:21`。

登录流程说明：

1. 创建认证管理器，见 `internal/pkg/authflow/manager.go:75`
2. 拉起飞书扫码登录流程，见 `internal/pkg/authflow/manager.go:191`
3. 获取 JWT token，见 `internal/pkg/authflow/manager.go:312`
4. 获取用户信息，见 `internal/pkg/authflow/manager.go:161`
5. 保存本地认证信息，见 `internal/command/auth/auth.go:44`

登录成功后会在本地写入：

- 配置目录：`~/.slam-cli`，见 `internal/pkg/config/config.go:65`
- 配置文件：`~/.slam-cli/config.json`，见 `internal/pkg/config/config.go:79`
- token 文件目录：`~/.slam-cli/auth/`，见 `internal/pkg/authflow/manager.go:86`
- cookie jar 文件：`~/.slam-cli/cookies`，见 `internal/pkg/authflow/manager.go:91`

### 查看登录状态

```bash
slam-cli auth status
```

状态命令定义在 `internal/command/auth/auth.go:80`，会输出：

- region
- 是否已登录
- 更新时间
- token 文件路径
- token 摘要
- 用户名
- 邮箱

登录状态判断逻辑见 `internal/pkg/config/config.go:38`。

### 查看配置

以下三个命令当前行为一致，都会输出配置文件路径和完整配置 JSON：

```bash
slam-cli config
slam-cli config list
slam-cli config get
```

实现位置：`internal/command/config/config.go:12`。

默认配置定义在 `internal/pkg/config/config.go:54`，默认值包括：

- `profile: default`
- `output: plain`
- `auth.logged_in: false`
- `auth.region: cn`

---

## 软件包查询

软件包命令定义在 `internal/command/package/package.go:17`。

软件包查询依赖已登录状态：`caller.New()` 会先调用 `configpkg.RequireLogin()`，见 `internal/pkg/caller/caller.go:19` 与 `internal/pkg/config/config.go:158`。

### 查询软件包列表

默认查询第一页：

```bash
slam-cli package list
```

指定页码与页面大小：

```bash
slam-cli package list --index 1 --size 20
slam-cli package list -i 2 -s 50
```

分页 flag 定义见：

- `internal/command/package/package.go:33`
- `internal/command/package/package.go:34`

按条件过滤时，参数必须以 **字段 / 值** 成对传入，校验逻辑见 `internal/command/package/package.go:99`。

示例：

```bash
slam-cli package list name bdaa
slam-cli package list owner liubing.065
slam-cli package list type src
slam-cli package list dkms true
slam-cli package list agent false
slam-cli package list name bdaa owner liubing.065 type src
```

当前支持的过滤字段：

- `name`
- `keyword`
- `keywords`
- `owner`
- `dkms`
- `agent`
- `type`

参数规则：

- `type` 仅支持 `src` 或 `bin`，见 `internal/command/package/package.go:135`
- `dkms` / `agent` 支持 `true/false/yes/no/1/0`，见 `internal/command/package/package.go:149`

列表查询最终调用 `/package` 接口，见 `internal/pkg/caller/caller.go:68`。

### 查询软件包详情

查询单个软件包：

```bash
slam-cli package info bdaa-sdk
```

查询多个软件包：

```bash
slam-cli package info bdaa-sdk bdaa-drivers
```

实现位置：`internal/command/package/package.go:48`。

注意事项：

- 至少需要一个包名参数，见 `internal/command/package/package.go:53`
- 当前输出始终为 JSON 数组，即使只查一个包也会返回数组结构，见 `internal/command/package/package.go:56`

详情查询最终调用 `/package/{name}` 接口，见 `internal/pkg/caller/caller.go:52`。

### 返回数据结构

请求与响应 DTO 定义在：

- 请求结构：`internal/pkg/caller/DTO/req.go:3`
- 列表响应：`internal/pkg/caller/DTO/resp.go:16`
- 全量响应：`internal/pkg/caller/DTO/resp.go:21`
- 详情响应：`internal/pkg/caller/DTO/resp.go:43`
- 版本信息：`internal/pkg/caller/DTO/resp.go:26`
- 简版软件包：`internal/pkg/caller/DTO/resp.go:79`

统一响应解码逻辑见 `internal/pkg/caller/caller.go:164`，格式为：

```json
{
  "code": 0,
  "data": {},
  "error": ""
}
```

---

## 快捷命令

快捷命令定义在 `internal/command/shortcut/shortcut.go:9`。

### `+doctor`

```bash
slam-cli +doctor
```

当前输出：

```text
slam-cli doctor: ok (skeleton mode)
```

实现位置：`internal/command/shortcut/shortcut.go:12`。

### `+init`

```bash
slam-cli +init
```

当前输出：

```text
slam-cli init: not implemented yet
```

实现位置：`internal/command/shortcut/shortcut.go:18`。

---

## 环境变量

SLAM 后端地址支持通过环境变量覆盖，逻辑见 `internal/pkg/caller/caller.go:124`。

### 直接指定完整 Base URL

```bash
export SLAM_API_BASE_URL="https://slam.byted.org/api/slam/v2"
```

### 只指定 Host

```bash
export SLAM_API_HOST="https://slam.byted.org"
```

程序会自动拼接为：

```text
https://slam.byted.org/api/slam/v2
```

默认地址常量定义在 `internal/pkg/caller/common.go:10`。

---

## 开发与校验

### 运行测试

```bash
go test ./...
```

### 常用本地调试命令

```bash
go run ./cmd/slam help
go run ./cmd/slam version
go run ./cmd/slam auth --help
go run ./cmd/slam config --help
go run ./cmd/slam package list --help
```

---

## 当前限制

截至当前版本，需要注意以下几点：

1. **软件包查询依赖本地登录态**
   - 未登录时会直接报错，见 `internal/pkg/config/config.go:164`。

2. **`package list` 过滤参数必须成对传入**
   - 否则会报错，见 `internal/command/package/package.go:100`。

3. **`package info` 至少需要一个包名**
   - 参数校验见 `internal/command/package/package.go:53`。

4. **`GetAllPackages` 仅在调用层实现，CLI 尚未暴露对应命令**
   - 见 `internal/pkg/caller/caller.go:101`。

5. **`+doctor` 与 `+init` 仍是占位实现**
   - 当前主要用于预留命令入口，见 `internal/command/shortcut/shortcut.go:9`。

---

## 总结

当前 `slam-cli` 已经具备基础可用性：

- 能登录并保存认证信息
- 能查看当前登录状态
- 能查看本地配置
- 能分页查询软件包列表
- 能查询一个或多个软件包详情

整体结构也已经比较清晰：

- CLI 装配集中在 `internal/app/root.go:14`
- 命令按领域放在 `internal/command/*`
- 基础设施能力沉淀在 `internal/pkg/*`

后续如果继续扩展，比较自然的方向包括：

- 增加 `auth logout`
- 增加 `config set`
- 暴露 `package all` 或类似全量查询命令
- 补充更完善的 `doctor` 检查项
- 增加更友好的输出格式与错误提示
