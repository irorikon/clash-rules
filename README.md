# Domain List Custom

基于 [v2fly/domain-list-community](https://github.com/v2fly/domain-list-community) 的数据构建域名规则列表，并提供 V2Ray/Xray `geosite.dat`、纯文本规则、GFWList 和 Clash Meta MRS 格式。

## 规则列表调整

生成 `geosite.dat` 时，默认按以下规则过滤属性：

- `cn`：排除 `@!cn` 和 `@ads`
- `geolocation-cn`：排除 `@!cn` 和 `@ads`
- `geolocation-!cn`：排除 `@cn` 和 `@ads`

这些过滤用于避免不需要的规则进入相应分类。其余规则继承上游列表及 include 关系。

## 下载

- [geosite.dat](https://github.com/Loyalsoldier/domain-list-custom/releases/latest/download/geosite.dat)
- [geosite.dat SHA-256](https://github.com/Loyalsoldier/domain-list-custom/releases/latest/download/geosite.dat.sha256sum)
- [所有发布文件](https://github.com/Loyalsoldier/domain-list-custom/releases/latest)

### Clash Meta MRS

MRS 文件按分类单独生成，文件名与对应的纯文本列表一致：

| 分类 | MRS 文件 | 纯文本文件 |
| --- | --- | --- |
| 广告域名 | [category-ads-all.mrs](https://github.com/Loyalsoldier/domain-list-custom/releases/latest/download/category-ads-all.mrs) | `category-ads-all.txt` |
| 中国顶级域名 | [tld-cn.mrs](https://github.com/Loyalsoldier/domain-list-custom/releases/latest/download/tld-cn.mrs) | `tld-cn.txt` |
| 中国域名 | [cn.mrs](https://github.com/Loyalsoldier/domain-list-custom/releases/latest/download/cn.mrs) | `cn.txt` |
| 中国大陆域名地理分类 | [geolocation-cn.mrs](https://github.com/Loyalsoldier/domain-list-custom/releases/latest/download/geolocation-cn.mrs) | `geolocation-cn.txt` |
| 非中国顶级域名 | [tld-!cn.mrs](https://github.com/Loyalsoldier/domain-list-custom/releases/latest/download/tld-!cn.mrs) | `tld-!cn.txt` |
| 非中国域名 | [geolocation-!cn.mrs](https://github.com/Loyalsoldier/domain-list-custom/releases/latest/download/geolocation-!cn.mrs) | `geolocation-!cn.txt` |
| 私有域名 | [private.mrs](https://github.com/Loyalsoldier/domain-list-custom/releases/latest/download/private.mrs) | `private.txt` |
| Apple | [apple.mrs](https://github.com/Loyalsoldier/domain-list-custom/releases/latest/download/apple.mrs) | `apple.txt` |
| iCloud | [icloud.mrs](https://github.com/Loyalsoldier/domain-list-custom/releases/latest/download/icloud.mrs) | `icloud.txt` |
| Google | [google.mrs](https://github.com/Loyalsoldier/domain-list-custom/releases/latest/download/google.mrs) | `google.txt` |
| Steam | [steam.mrs](https://github.com/Loyalsoldier/domain-list-custom/releases/latest/download/steam.mrs) | `steam.txt` |

MRS 采用 Clash Meta 的 `domain` 行为，只支持完整域名和根域名规则。源列表中的 `keyword`、`regexp` 规则不会写入 MRS 文件；纯文本列表和 `geosite.dat` 仍保留这些规则。

Rule Provider 配置示例：

```yaml
rule-providers:
  private:
    type: http
    behavior: domain
    format: mrs
    url: https://github.com/Loyalsoldier/domain-list-custom/releases/latest/download/private.mrs
    path: ./ruleset/private.mrs
    interval: 86400

rules:
  - RULE-SET,private,DIRECT
```

## 构建

需要 Go 环境以及 `v2fly/domain-list-community` 的 `data` 目录：

```sh
go run ./ \
  --datapath=./domain-list-community/data \
  --outputpath=./publish
```

程序默认输出：

- `geosite.dat`
- 默认分类对应的 `.txt` 和 `.mrs` 文件
- `gfwlist.txt`

可通过命令行参数调整输入、输出和导出分类：

| 参数 | 默认值 | 说明 |
| --- | --- | --- |
| `datapath` | `./data` | 上游域名列表目录 |
| `datname` | `geosite.dat` | 生成的 GeoSite 文件名 |
| `outputpath` | `./publish` | 输出目录 |
| `exportlists` | `category-ads-all,tld-cn,cn,geolocation-cn,tld-!cn,geolocation-!cn,private,apple,icloud,google,steam` | 以逗号分隔的纯文本和 MRS 分类 |
| `excludeattrs` | `cn@!cn@ads,geolocation-cn@!cn@ads,geolocation-!cn@cn@ads` | `geosite.dat` 中按列表排除的属性 |
| `togfwlist` | `geolocation-!cn` | 生成 GFWList 的列表 |
