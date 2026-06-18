# xflow create 示例

## 命令用法

```bash
slam-cli xflow create --file create-xflow.json
slam-cli xflow create --file create-xflow.json --action-user alice
```

## publish 场景一：仅上源，不发起变更

常见于源码包，仅填写 `reasonNoChange`，不传 `changeOrder`。

```json
{
  "type": "publish",
  "package": "bdaa-sdk",
  "publishOpts": {
    "docsV2": {
      "reviewDoc": "https://example.com/review",
      "releaseNote": "https://example.com/release",
      "testDoc": "https://example.com/test"
    },
    "versionNote": "填写版本信息",
    "reasonNoChange": "仅上源，不涉及线上机器变更",
    "branch": "main",
    "region": ["cn"],
    "dist": ["Debian13:riscv64"],
    "files": [],
    "hwKeys": []
  }
}
```

## publish 场景二：上源并发起变更

常见于二进制包，传 `changeOrder`，bin 包必须提供 `files`。

```json
{
  "type": "publish",
  "package": "ip-address-env",
  "publishOpts": {
    "docsV2": {
      "reviewDoc": "https://example.com/review",
      "releaseNote": "https://example.com/release",
      "testDoc": "https://example.com/test"
    },
    "versionNote": "Debian13 riscv64 上源，同时发起变更",
    "smokeTestParams": {
      "linuxKernel": ["linux-6.6"]
    },
    "changeOrder": {
      "name": "ip-address-env 升级变更",
      "source": "system-package",
      "changeType": "upgrade",
      "changeTemplate": {
        "targetPackage": "ip-address-env",
        "packageVersion": "tbd",
        "versionSource": "SLAM"
      },
      "changeLevel": "L2",
      "changeEnv": "Online",
      "changeObject": {
        "type": "PSM",
        "psmList": ["data.example.psm"],
        "ips": {
          "ipSelector": "ipFile",
          "ipList": [],
          "ipFile": ""
        },
        "customFilterConfigs": [
          {
            "filterType": "idc",
            "items": ["lf"]
          }
        ]
      },
      "healthCheckItems": ["host_alive", "ssh_connectivity"],
      "creator": "alice",
      "executors": ["alice"],
      "approvers": {
        "sreList": ["bob"],
        "packageOwnerList": ["carol"]
      },
      "changeWindow": [1750000000000, 1750086400000],
      "changeScripts": {
        "pre": "",
        "deploy": "apt install -y ip-address-env",
        "post": "",
        "rollback": "apt install -y ip-address-env=previous-version"
      },
      "notification": {
        "type": ["larkUser"],
        "larkUserList": ["alice"],
        "larkGroupList": []
      },
      "grayConfig": {
        "concurrency": 1,
        "grayObject": {
          "selector": "percentage",
          "percentage": 10,
          "filter": []
        },
        "taskConfig": [
          { "percentage": 50, "cd": 300 },
          { "percentage": 100, "cd": 600 }
        ]
      },
      "prodConfig": {
        "concurrency": 5,
        "taskConfig": [
          { "percentage": 50, "cd": 300 },
          { "percentage": 100, "cd": 600 }
        ]
      },
      "autoCreateXflow": false
    },
    "reasonNoChange": "",
    "region": ["cn"],
    "dist": ["Debian13:riscv64"],
    "files": [
      {
        "name": "ip-address-env_3.9.1+byted_all.deb",
        "key": "packages/ip-address-env/ip-address-env_3.9.1+byted_all.deb"
      }
    ],
    "hwKeys": []
  }
}
```

## changeOrder 字段说明

| 字段 | 类型 | 含义 |
|------|------|------|
| `name` | string | 变更名称 |
| `source` | enum | 变更来源：`system-package` / `kernel` / `idc-infra` |
| `changeType` | enum | 变更类型：`upgrade` / `downgrade` / `install` / `uninstall` 等 |
| `changeTemplate` | object | `targetPackage` / `packageVersion` / `versionSource`（`SLAM` 或 `Custom`） |
| `changeLevel` | enum | `L0`~`L3` |
| `changeEnv` | enum | `Online` / `BOE` |
| `changeObject` | object | 变更对象，包含 `type`(`PSM`/`IP`/`PSMandIP`)、`psmList`、`ips`、`customFilterConfigs` |
| `healthCheckItems` | string[] | 主机健康检查项 |
| `creator` / `executors` / `approvers` | - | 发起人 / 操作人 / 审批人（含 `sreList`、`packageOwnerList`） |
| `changeWindow` | [number, number] | 变更窗口起止时间戳（毫秒） |
| `changeScripts` | object | `pre` / `deploy` / `post` / `rollback` 脚本 |
| `notification` | object | 通知配置：`type`(`larkUser`/`larkGroup`)、`larkUserList`、`larkGroupList` |
| `grayConfig` | object | 灰度阶段配置：`concurrency`、`grayObject`、`taskConfig` |
| `prodConfig` | object | 正式阶段配置：`concurrency`、`taskConfig` |
| `autoCreateXflow` | bool | 是否自动创建 xflow |

## register 类型示例

```json
{
  "type": "register",
  "package": "new-package",
  "registerOpts": {
    "packageType": "src",
    "owner": ["alice"],
    "repo": "git@example.com/repo.git",
    "branch": "main",
    "goVersion": "1.22",
    "files": [],
    "dist": ["Debian13:riscv64"],
    "region": ["cn"],
    "dkms": false,
    "desc": "package description",
    "agent": false,
    "agentType": "",
    "resourceLimit": "",
    "publishNotifyChatID": "",
    "openSource": false,
    "codeLanguage": ["golang"],
    "listenPorts": [],
    "involvedHardtype": "",
    "hasNDA": false
  }
}
```

## sync 类型示例

```json
{
  "type": "sync",
  "package": "bdaa-sdk",
  "syncOpts": {
    "versionID": 12345,
    "region": ["cn", "sg"]
  }
}
```
