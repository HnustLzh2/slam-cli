# slam-cli

`slam-cli` 是面向 agent 使用的 SLAM 命令行工具。当前能力覆盖本地认证与配置、软件包查询与更新、livepatch/third-repo/task/xflow 查询、xflow 创建，以及本地辅助命令。

当前版本：`0.1.0`，定义见 `internal/app/version.go`。

## 当前能力

根命令注册在 `internal/app/root.go`，当前工作区入口为 `cmd/slam-cli/main.go`。核心命令包括：

```text
slam-cli help
slam-cli version
slam-cli auth login
slam-cli auth status
slam-cli config
slam-cli config list
slam-cli config get
slam-cli package list
slam-cli package info
slam-cli package update
slam-cli livepatch list
slam-cli livepatch <product-id>
slam-cli third-repo list
slam-cli task list
slam-cli xflow list
slam-cli xflow create
slam-cli completion bash|zsh|fish|powershell|ps
slam-cli +doctor
slam-cli +init
```

后端查询和写入命令都会通过 `internal/pkg/caller.New()` 读取本地登录态；未登录时会要求先执行 `slam-cli auth login`。

## 项目结构

```text
slam-cli/
├── cmd/slam-cli/main.go
├── internal/app/                 # 根命令、help、version、completion
├── internal/command/auth/         # 登录与状态
├── internal/command/config/       # 配置展示
├── internal/command/package/      # package list/info/update
├── internal/command/livepatch/    # livepatch list/detail
├── internal/command/thirdrepo/    # third-repo list
├── internal/command/task/         # task list
├── internal/command/xflow/        # xflow list/create
├── internal/command/shortcut/     # +doctor、+init
├── internal/pkg/authflow/         # 飞书扫码登录、JWT、用户信息
├── internal/pkg/caller/           # SLAM API client
├── internal/pkg/caller/DTO/       # 请求与响应结构
├── internal/pkg/config/           # 本地配置与登录态
└── skills/slam-cli/               # slam-cli agent skills
```

## 快速开始

项目要求 Go `1.26.1`，定义见 `go.mod`。

```bash
go run ./cmd/slam-cli help
go run ./cmd/slam-cli version
go build -o slam-cli ./cmd/slam-cli
```

主要依赖包括 `cobra`、`resty`、`persistent-cookiejar` 和 `qrterminal`。

## 认证与配置

登录：

```bash
slam-cli auth login --region cn
slam-cli auth login --region us
slam-cli auth login --region sg
```

`--region` 支持短参数 `-r`，默认是 `cn`。登录会执行飞书扫码 SSO，获取 JWT token 和用户信息，并写入本地配置。

查看登录状态：

```bash
slam-cli auth status
```

查看配置：

```bash
slam-cli config
slam-cli config list
slam-cli config get
```

本地文件：

- 配置目录：`~/.slam-cli`
- 配置文件：`~/.slam-cli/config.json`
- token 文件：`~/.slam-cli/auth/<region>.token`
- cookie jar：`~/.slam-cli/cookies`

## 软件包

列表查询：

```bash
slam-cli package list
slam-cli package list --index 1 --size 20
slam-cli package list -i 2 -s 50
slam-cli package list --my
```

过滤参数必须按 `字段 值` 成对传入：

```bash
slam-cli package list name bdaa
slam-cli package list owner alice
slam-cli package list type src
slam-cli package list dkms true
slam-cli package list agent false
slam-cli package list name bdaa owner alice type src
```

支持字段：

- `name` / `keyword` / `keywords`
- `owner`
- `dkms`
- `agent`
- `type`

规则：

- `type` 仅支持 `src` 或 `bin`
- `dkms` / `agent` 支持 `true/false/yes/no/1/0`，代码也接受 `t/f/y/n`
- `--my` / `-m` 会使用当前登录用户作为 owner

详情查询：

```bash
slam-cli package info bdaa-sdk
slam-cli package info bdaa-sdk bdaa-drivers
```

输出始终是 JSON 数组，即使只查一个包。

更新软件包：

```bash
slam-cli package update bdaa-sdk --file update-package.json
```

这是写接口，对应 `PUT /api/slam/v2/package/:name`。`<name>` 是路径里的权威包名，CLI 会写入请求体 `name` 字段。请求 JSON 字段与 `internal/pkg/caller/DTO/req.go` 的 `UpdatePackageReq` 对齐。详细字段说明可查看：

```bash
slam-cli package update --help
```

## Work 查询

### Livepatch

```bash
slam-cli livepatch 1001
slam-cli livepatch list
slam-cli livepatch list --index 1 --size 20
slam-cli livepatch list --id 1001 --my
slam-cli livepatch list --status testing --creator alice
slam-cli livepatch list --distribution debian --os-version 12 --arch amd64
```

`livepatch list` 支持 `--id`、`--name`、`--status`、`--version`、`--distribution`、`--os-version`、`--kernel-version`、`--arch`、`--creator`、`--start-time`、`--end-time`、`--my`。

### Third Repo

```bash
slam-cli third-repo list
slam-cli third-repo list --index 1 --size 20
slam-cli third-repo list --my
```

当前 CLI 只暴露分页和 `--my`，没有暴露 name/type/owner flags。

### Task

```bash
slam-cli task list
slam-cli task list --index 1 --size 20
slam-cli task list --id 123
slam-cli task list --package kernel-agent
slam-cli task list --type build --state running
slam-cli task list --source xflow --creator alice
slam-cli task list --my
```

支持 `--id`、`--package`、`--type`、`--state`、`--source`、`--creator`、`--my`。

### Xflow

查询：

```bash
slam-cli xflow list
slam-cli xflow list --index 1 --size 20
slam-cli xflow list --id 123456
slam-cli xflow list --package kernel-agent
slam-cli xflow list --type publish --state running
slam-cli xflow list --creator alice
slam-cli xflow list --my
```

支持 `--id`、`--package`、`--type`、`--state`、`--creator`、`--my`。

创建：

```bash
slam-cli xflow create --file create-xflow.json
slam-cli xflow create --file create-xflow.json --action-user alice
```

这是写接口，对应 `POST /api/slam/v2/xflow`。请求 JSON 字段与 `CreateXflowReq` 对齐，`type` 支持 `publish`、`register`、`sync`，分别使用 `publishOpts`、`registerOpts`、`syncOpts`。详细 JSON 说明可查看：

```bash
slam-cli xflow create --help
```

其中 `publish` 要区分两种场景：

- 仅上源不变更：不传 `changeOrder`，填写 `reasonNoChange`
- 上源加变更：传 `changeOrder`，`reasonNoChange` 置空

并且要注意：

- src 包发布通常需要 `branch`
- bin 包发布通常需要 `files`
- dkms 包发布通常需要 `smokeTestParams`
- `dist` 建议使用后端认可的完整字符串，例如 `Debian13:riscv64`

## 本地辅助命令

```bash
slam-cli help
slam-cli version
slam-cli +doctor
slam-cli completion bash
slam-cli completion zsh
slam-cli completion fish
slam-cli completion powershell
slam-cli completion ps
```

`+doctor` 当前输出：

```text
slam-cli doctor: ok (skeleton mode)
```

`+init` 会修改用户 shell 配置：

```bash
slam-cli +init
```

它会向 `~/.zshrc` 追加 `source <(slam-cli completion zsh)` 配置块，向 `~/.bashrc` 追加 `source <(slam-cli completion bash)` 配置块，并执行 `source ~/.zshrc` 与 `source ~/.bashrc`。agent 执行前必须确认用户允许。

## Skills 使用方法

仓库里的 skills 目录是：`skills/slam-cli/`。

- 路由入口：`skills/slam-cli/SKILL.md`
- 项目开发：`skills/slam-cli/slam-cli-project/SKILL.md`
- 认证配置：`skills/slam-cli/slam-cli-auth-config/SKILL.md`
- 包查询更新：`skills/slam-cli/slam-cli-package-query/SKILL.md`
- work 查询与 xflow：`skills/slam-cli/slam-cli-work-query/SKILL.md`
- 本地辅助工具：`skills/slam-cli/slam-cli-local-tools/SKILL.md`

使用原则：

1. 优先选最小子 skill，不要一次加载所有 skill。
2. 涉及后端查询或写接口前，先确认登录态：`slam-cli auth status`。
3. `xflow create` 统一通过 JSON 文件驱动，先看：`slam-cli xflow create --help`。
4. `publish` 有两种场景：
   - 仅上源不变更：不传 `changeOrder`，填写 `reasonNoChange`
   - 上源加变更：传 `changeOrder`，`reasonNoChange` 置空
5. `dist` 建议使用后端认可的完整字符串，例如 `Debian13:riscv64`。

## 环境变量

默认后端地址定义在 `internal/pkg/caller/common.go`：

```text
https://slam.byted.org/api/slam/v2
```

可用环境变量：

```bash
export SLAM_API_BASE_URL="https://slam.byted.org/api/slam/v2"
export SLAM_API_HOST="https://slam.byted.org"
```

`SLAM_API_BASE_URL` 会完整覆盖 base URL；`SLAM_API_HOST` 会拼接为 `<host>/api/slam/v2`。

## 开发与校验

执行测试时必须使用项目要求的 Go 版本，并从 `/Users/bytedance/sdk/` 下查找对应 Go：

```bash
find /Users/bytedance/sdk -maxdepth 4 -type f -name go
/Users/bytedance/sdk/<go1.26.1>/bin/go test ./...
```

如果 `/Users/bytedance/sdk/` 下没有 Go `1.26.1`，跳过测试并说明原因。

常用本地调试命令：

```bash
go run ./cmd/slam-cli help
go run ./cmd/slam-cli package list --help
go run ./cmd/slam-cli package update --help
go run ./cmd/slam-cli xflow create --help
```

## 当前边界

- `GetAllPackages` 在 `internal/pkg/caller` 中已有底层方法，但当前没有 CLI 子命令。
- `third-repo list` 底层 DTO 有 `name`、`type`、`owner` 字段，但当前 CLI 没有注册对应 flags。
- 当前支持的写接口只有 `package update` 和 `xflow create`；没有删除、审批、发布执行、任意 API 调用等命令。
- `config`、`config list`、`config get` 当前行为一致，都会打印配置文件路径和完整配置 JSON。
