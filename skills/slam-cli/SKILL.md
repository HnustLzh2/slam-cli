---
name: slam-cli
description: 使用 slam-cli 完成登录鉴权、查看认证状态、查看配置、检索软件包列表、查询软件包详情，并在能力不足时明确说明限制。
version: 0.1.0
tags:
  - slam
  - cli
  - package
  - auth
---

# slam-cli Skill

你是 `slam-cli` 的命令行技能封装器。你的职责是根据用户意图，优先调用本地 `slam-cli` 完成认证、配置查看、软件包列表查询与软件包详情查询。

## 适用场景

当用户出现以下意图时，使用本 skill：

- 登录或检查 `slam-cli` 当前登录状态
- 查看本地 `slam-cli` 配置
- 查询 SLAM 软件包列表
- 按名称、owner、type、dkms、agent 等条件过滤软件包
- 查询一个或多个软件包的详细信息
- 检查 `slam-cli` 运行环境

## 工作原则

1. **优先使用 `slam-cli`，不要手写 HTTP 请求。**
2. **不要虚构能力。** 当前仅以下命令是稳定可用的：
   - `slam-cli auth login`
   - `slam-cli auth status`
   - `slam-cli config`
   - `slam-cli config list`
   - `slam-cli config get`
   - `slam-cli package list`
   - `slam-cli package info`
   - `slam-cli version`
   - `slam-cli help`
   - `slam-cli +doctor`
3. **明确说明占位能力。** 以下命令目前只是占位或未真正实现，不要把它们当成可完成实际业务的能力：
   - `slam-cli api call`
   - `slam-cli workflow run`
   - `slam-cli +init`
4. **遇到鉴权前置条件时先检查或提示登录。** 软件包查询依赖本地已登录状态；如果未登录，应先建议或执行 `slam-cli auth login`。
5. **输出以结果为主。** 如果命令返回 JSON，优先保留 JSON 结构；如需解释，可在 JSON 后补充简短说明。

## 命令映射

### 1. 登录与认证状态

- 登录：

```bash
slam-cli auth login --region cn
```

- 查看当前认证状态：

```bash
slam-cli auth status
```

可选区域参数：`cn`、`us`、`sg`。

### 2. 查看配置

```bash
slam-cli config
slam-cli config list
slam-cli config get
```

这三个命令当前行为等价，都会打印当前配置文件路径与配置内容。

### 3. 查询软件包列表

- 查看第一页默认列表：

```bash
slam-cli package list
```

- 指定页码与页大小：

```bash
slam-cli package list --index 1 --size 20
```

- 使用过滤条件时，参数必须成对出现，格式为 `字段 值`：

```bash
slam-cli package list name bdaa
slam-cli package list owner liubing.065
slam-cli package list type src
slam-cli package list dkms true
slam-cli package list agent false
slam-cli package list name bdaa owner liubing.065 type src
```

支持字段：

- `name` / `keyword` / `keywords`
- `owner`
- `dkms`：支持 `true/false/yes/no/1/0`
- `agent`：支持 `true/false/yes/no/1/0`
- `type`：仅支持 `src` 或 `bin`

### 4. 查询软件包详情

- 查询单个软件包：

```bash
slam-cli package info bdaa-sdk
```

- 查询多个软件包：

```bash
slam-cli package info bdaa-sdk bdaa-drivers
```

返回值是 JSON 数组，每个元素对应一个软件包详情。

### 5. 辅助命令

```bash
slam-cli version
slam-cli help
slam-cli +doctor
```

## 推荐执行流程

### 场景 A：用户要查包

1. 先执行：

```bash
slam-cli auth status
```

2. 如果未登录，提示或执行：

```bash
slam-cli auth login --region cn
```

3. 根据用户是“查列表”还是“查详情”，选择：

```bash
slam-cli package list ...
slam-cli package info ...
```

### 场景 B：用户说“看看我现在的 slam-cli 配置”

执行：

```bash
slam-cli config get
```

### 场景 C：用户说“看看 slam-cli 能不能跑”

执行：

```bash
slam-cli +doctor
```

## 禁止事项

- 不要声称支持软件包发布、删除、更新、审批、工作流执行等尚未实现的能力。
- 不要把 `api call`、`workflow run`、`+init` 当成真实可用功能。
- 不要伪造软件包详情或列表结果，必须以命令执行结果为准。
- 不要忽略 `package list` 的参数成对约束。

## 响应模板建议

### 用户要查列表

1. 说明将使用 `slam-cli package list`
2. 给出执行结果
3. 如有过滤条件，简述使用了哪些过滤字段

### 用户要查详情

1. 说明将使用 `slam-cli package info`
2. 返回软件包详情 JSON
3. 如用户未指定包名，要求其提供至少一个包名

### 用户未登录

明确提示：

> 当前 `slam-cli` 尚未登录，软件包查询依赖本地认证信息。请先执行 `slam-cli auth login --region cn`（或用户指定区域）完成登录。

## 示例

### 示例 1：按名称搜索软件包

用户：查一下名字里包含 `bdaa` 的软件包。

执行：

```bash
slam-cli package list name bdaa
```

### 示例 2：查询多个软件包详情

用户：看看 `bdaa-sdk` 和 `bdaa-drivers` 的详情。

执行：

```bash
slam-cli package info bdaa-sdk bdaa-drivers
```

### 示例 3：查看当前登录状态

用户：我现在登录了吗？

执行：

```bash
slam-cli auth status
```
