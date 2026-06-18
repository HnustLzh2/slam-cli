# xflow create

创建 xflow 工单，对应 SLAM API：`POST /api/slam/v2/xflow`。

这是写接口，执行前必须确认已登录：`slam-cli auth status`。

- 请求体通过 `--file` 传入，JSON 字段与 SLAM API `CreateXflowReq` 对齐。
- `--action-user` 会写入 `X-Action-User` header，主要用于服务账号代真实用户发起的场景。

## publish 场景

publish 有两种常见场景：

1. **仅上源不发起变更**：不传 `changeOrder`，填写 `reasonNoChange`。
2. **上源并发起变更**：传 `changeOrder`，`reasonNoChange` 置空。

### 后端行为说明

- `publish` + `changeOrder`：后端会基于 `changeOrder` 创建变更。
- `publish` + 无 `changeOrder`：评审消息会展示"仅上源未变更，理由：reasonNoChange"。
- src 包 publish：`branch` 必填。
- bin 包 publish：`files` 必填。
- dkms 包 publish：`smokeTestParams` 必填。
- `dist` 必须使用后端认可的发行版字符串，例如 `Debian13:riscv64`。

## create-xflow.json 必含字段

| 字段 | 说明 |
|------|------|
| `type` | 工单类型：`publish`、`register`、`sync` |
| `package` | 软件包名 |

## 不同 type 对应的 opts

| type | opts | 说明 |
|------|------|------|
| `publish` | `publishOpts` | 发布已有 package 的新版本 |
| `register` | `registerOpts` | 注册新 package |
| `sync` | `syncOpts` | 同步已有版本到指定 region |

## RegisterOpts：注册软件包参数

| 字段 | 说明 |
|------|------|
| `packageType` | 软件包类型。可选值通常为 `src` / `bin`，分别表示源码包 / 二进制包 |
| `owner` | 软件包负责人列表 |
| `repo` | 源码仓库地址或仓库标识。源码包通常必填 |
| `branch` | 源码分支。源码包注册或发布时使用 |
| `goVersion` | Go 版本。仅源码包使用，用于后续构建任务选择 Go 环境 |
| `files` | 文件列表。二进制包注册时通常需要上传文件，元素为 `rpc.File` |
| `dist` | 目标发行版 / 架构 / 分发目标列表 |
| `region` | 生效区域列表 |
| `dkms` | 是否是 DKMS 包。DKMS 包发布时通常要求提供冒烟测试参数 |
| `agent` | 是否是 Agent 类软件包 |
| `desc` | 软件包描述信息 |
| `publishNotifyChatID` | 发布通知群 ID。用于发布完成后的通知 |
| `agentType` | Agent 类型。`agent=true` 时需要提供 |
| `resourceLimit` | 资源限制信息。Agent 包或版本元数据中会使用 |
| `openSource` | 是否为开源软件 |
| `codeLanguage` | 代码语言列表。源码包使用，例如 `c`、`cpp`、`java`、`python`、`golang` 等 |
| `listenPorts` | 监听端口列表，用于描述服务暴露或监听的端口信息 |
| `involvedHardtype` | 涉及的硬件类型 |
| `hasNDA` | 是否涉及 NDA |

## PublishOpts：发布版本参数

| 字段 | 说明 |
|------|------|
| `docs` | 旧版文档字段。当前主要使用 `docsV2` |
| `docsV2` | 新版交付文档信息，类型为 `rpc.PublishDocs`，通常包含评审文档、发布说明、测试文档等 |
| `versionNote` | 版本说明 / 发布说明 |
| `smokeTestParams` | 冒烟测试参数。DKMS 包发布时通常必填 |
| `changeOrder` | 变更单数据。如果提供，会用于创建变更；这是"上源+变更"场景 |
| `reasonNoChange` | 仅上源不发起变更的原因。`changeOrder` 为空时应填写 |
| `dist` | 本次发布覆盖的发行版 / 目标列表 |
| `region` | 本次发布覆盖的区域列表 |
| `branch` | 本次发布使用的源码分支。源码包发布时通常必填 |
| `files` | 本次发布上传的产物文件列表。二进制包发布时通常必填 |
| `hwKeys` | 硬件 key 列表，用于标识硬件相关适配或约束 |

## SyncOpts：同步版本参数

| 字段 | 说明 |
|------|------|
| `versionID` | 要同步的版本 ID |
| `region` | 要同步到的区域列表 |
