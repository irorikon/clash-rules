# Domain List Custom

基于 [v2fly/domain-list-community](https://github.com/v2fly/domain-list-community) 数据构建 `geosite.dat`，并按分流用途发布 Clash 规则集。

## 分类规则

GitHub Actions 按以下分类生成独立文件：

| 分类 | 类型 | 文件 |
| --- | --- | --- |
| `apple` | 域名 | [apple.txt](https://github.com/Loyalsoldier/domain-list-custom/releases/latest/download/apple.txt) · [apple.mrs](https://github.com/Loyalsoldier/domain-list-custom/releases/latest/download/apple.mrs) |
| `applications` | classical | [applications.txt](https://github.com/Loyalsoldier/domain-list-custom/releases/latest/download/applications.txt) |
| `cncidr` | IP-CIDR | [cncidr.txt](https://github.com/Loyalsoldier/domain-list-custom/releases/latest/download/cncidr.txt) · [cncidr.mrs](https://github.com/Loyalsoldier/domain-list-custom/releases/latest/download/cncidr.mrs) |
| `direct` | 域名 | [direct.txt](https://github.com/Loyalsoldier/domain-list-custom/releases/latest/download/direct.txt) · [direct.mrs](https://github.com/Loyalsoldier/domain-list-custom/releases/latest/download/direct.mrs) |
| `gfw` | 域名 | [gfw.txt](https://github.com/Loyalsoldier/domain-list-custom/releases/latest/download/gfw.txt) · [gfw.mrs](https://github.com/Loyalsoldier/domain-list-custom/releases/latest/download/gfw.mrs) |
| `google` | 域名 | [google.txt](https://github.com/Loyalsoldier/domain-list-custom/releases/latest/download/google.txt) · [google.mrs](https://github.com/Loyalsoldier/domain-list-custom/releases/latest/download/google.mrs) |
| `greatfire` | 域名 | [greatfire.txt](https://github.com/Loyalsoldier/domain-list-custom/releases/latest/download/greatfire.txt) · [greatfire.mrs](https://github.com/Loyalsoldier/domain-list-custom/releases/latest/download/greatfire.mrs) |
| `icloud` | 域名 | [icloud.txt](https://github.com/Loyalsoldier/domain-list-custom/releases/latest/download/icloud.txt) · [icloud.mrs](https://github.com/Loyalsoldier/domain-list-custom/releases/latest/download/icloud.mrs) |
| `lancidr` | IP-CIDR | [lancidr.txt](https://github.com/Loyalsoldier/domain-list-custom/releases/latest/download/lancidr.txt) · [lancidr.mrs](https://github.com/Loyalsoldier/domain-list-custom/releases/latest/download/lancidr.mrs) |
| `private` | 域名 | [private.txt](https://github.com/Loyalsoldier/domain-list-custom/releases/latest/download/private.txt) · [private.mrs](https://github.com/Loyalsoldier/domain-list-custom/releases/latest/download/private.mrs) |
| `proxy` | 域名 | [proxy.txt](https://github.com/Loyalsoldier/domain-list-custom/releases/latest/download/proxy.txt) · [proxy.mrs](https://github.com/Loyalsoldier/domain-list-custom/releases/latest/download/proxy.mrs) |
| `reject` | 域名 | [reject.txt](https://github.com/Loyalsoldier/domain-list-custom/releases/latest/download/reject.txt) · [reject.mrs](https://github.com/Loyalsoldier/domain-list-custom/releases/latest/download/reject.mrs) |
| `telegramcidr` | IP-CIDR | [telegramcidr.txt](https://github.com/Loyalsoldier/domain-list-custom/releases/latest/download/telegramcidr.txt) · [telegramcidr.mrs](https://github.com/Loyalsoldier/domain-list-custom/releases/latest/download/telegramcidr.mrs) |
| `tld-not-cn` | 域名 | [tld-not-cn.txt](https://github.com/Loyalsoldier/domain-list-custom/releases/latest/download/tld-not-cn.txt) · [tld-not-cn.mrs](https://github.com/Loyalsoldier/domain-list-custom/releases/latest/download/tld-not-cn.mrs) |

`.txt` 文件为 Clash Rule Provider 格式。`applications` 包含进程匹配规则，使用 `behavior: classical`，不生成 MRS；其余分类按 `domain` 或 `ipcidr` 行为生成 MRS。

`geosite.dat` 仍会生成，并沿用以下属性过滤：

- `cn`：排除 `@!cn` 和 `@ads`
- `geolocation-cn`：排除 `@!cn` 和 `@ads`
- `geolocation-!cn`：排除 `@cn` 和 `@ads`

## 使用示例

通过 MRS 加载域名分类：

```yaml
rule-providers:
  proxy:
    type: http
    behavior: domain
    format: mrs
    url: https://github.com/Loyalsoldier/domain-list-custom/releases/latest/download/proxy.mrs
    path: ./ruleset/proxy.mrs
    interval: 86400

rules:
  - RULE-SET,proxy,PROXY
```

CIDR 分类使用 `behavior: ipcidr`。`applications` 是 classical 格式，应使用 `format: yaml`，不能作为 MRS 加载。

## 下载

- [geosite.dat](https://github.com/Loyalsoldier/domain-list-custom/releases/latest/download/geosite.dat)
- [geosite.dat SHA-256](https://github.com/Loyalsoldier/domain-list-custom/releases/latest/download/geosite.dat.sha256sum)
- [全部发布文件](https://github.com/Loyalsoldier/domain-list-custom/releases/latest)

## 本地构建

需要 Go、Bash、curl、awk，以及 `v2fly/domain-list-community` 的 `data` 目录：

```sh
go run ./ \
  --datapath=./domain-list-community/data \
  --outputpath=./publish \
  --exportlists=icloud,private,tld-!cn

bash ./.github/scripts/generate-rules.sh ./publish
```

分类生成脚本会从上游规则源拉取并整理其余分类，再生成对应的 MRS 文件。`applications.txt` 为 classical 规则，单独发布。
