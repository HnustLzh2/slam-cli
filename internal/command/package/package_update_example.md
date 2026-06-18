# package update 示例

## 命令用法

```bash
slam-cli package update bdaa-sdk --file update-package.json
```

## update-package.json 示例

```json
{
  "owner": ["alice"],
  "type": "src",
  "repo": "git@example.com/repo.git",
  "branch": "main",
  "goVersion": "1.22",
  "region": ["cn"],
  "dist": ["Debian13:riscv64"],
  "dkms": false,
  "agent": false,
  "desc": "package description",
  "publishNotifyChatID": "",
  "agentType": "",
  "resourceLimit": "",
  "openSource": false,
  "codeLanguage": ["golang"],
  "listenPorts": [
    {"protocol": "tcp", "port": 8080, "desc": "http service"}
  ],
  "involvedHardtype": "",
  "hasNDA": false
}
```

## bin 类型 files 示例

```json
"files": [{"name": "demo.deb", "key": "tmp/uploaded-key"}]
```

## dist 字段说明

`dist` 建议使用后端认可的完整发行版字符串，例如：

```json
"dist": ["Debian13:riscv64"]
```
