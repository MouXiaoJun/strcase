# go-strcase

[English](README.md)

[![Go Reference](https://pkg.go.dev/badge/github.com/MouXiaoJun/strcase.svg)](https://pkg.go.dev/github.com/MouXiaoJun/strcase)
[![Go Version](https://img.shields.io/badge/go-1.21+-00ADD8?style=flat-square&logo=go)](https://golang.org)
[![License](https://img.shields.io/badge/license-MIT-green.svg?style=flat-square)](LICENSE)

一个**零依赖**的 Go 命名风格转换库:小驼峰 / 大驼峰 / 蛇形 / 短横线 / 大写蛇形,内置一套文档化、可测试的单词切分规则,从容处理缩写、数字、连续大写与混合分隔符。

## 特性

- ✅ **零依赖** — 仅标准库(`strings` / `unicode`),`go.mod` 干净,无供应链风险
- ✅ **5 种风格,1 套分词器** — 所有函数按同一套文档化规则切词,再按各自风格拼装
- ✅ **缩写感知** — `HTTPRequest` → `http_request`(而不是 `httprequest`)
- ✅ **数字边界** — `version2` → `version_2`;`SHA256` → `sha_256`;`ID2Name` → `id_2_name`
- ✅ **混合分隔符** — `user--_name`、`USER_NAME` 都归一化为 `user_name`
- ✅ **确定性 + ASCII 幂等** — 对纯 ASCII 输入,分隔符风格是不动点(`ToSnake(ToSnake(s)) == ToSnake(s)`),且互相可转(`ToSnake(ToKebab(s)) == ToSnake(s)`)
- ✅ **Fuzz 验证不变式** — `FuzzStrCase` 验证任意输入不 panic,并在 ASCII 输入上验证整套不变式
- ✅ **Unicode 安全** — 中文/拉丁字母原样保留(`用户Name` 可读);非法 UTF-8 也不 panic

## 快速开始

```go
package main

import (
	"fmt"

	"github.com/MouXiaoJun/strcase"
)

func main() {
	fmt.Println(strcase.ToSnake("HTTPRequest"))   // http_request
	fmt.Println(strcase.ToCamel("user_name"))     // userName
	fmt.Println(strcase.ToPascal("user_name"))    // UserName
	fmt.Println(strcase.ToKebab("HTTPRequest"))   // http-request
	fmt.Println(strcase.ToUpperSnake("userName")) // USER_NAME
}
```

## API

全部函数输入输出 `string`;空串映射为空串。

| 函数 | 示例输入 | 示例输出 | 说明 |
| --- | --- | --- | --- |
| `ToCamel` | `user_name` | `userName` | 首词整词小写,后续词首字母大写 |
| `ToPascal` | `user_name` | `UserName` | 每词首字母大写(缩写同样归一化) |
| `ToSnake` | `HTTPRequest` | `http_request` | 每词小写,`_` 连接 |
| `ToKebab` | `userName` | `user-name` | 每词小写,`-` 连接 |
| `ToUpperSnake` | `userName` | `USER_NAME` | 每词大写,`_` 连接 |

## 单词切分规则

所有函数先按同一套规则把输入切成词,再重新拼装。规则确定、完全由表驱动测试覆盖:

1. **分隔符** — `_`、`-`、空格都是单词边界;连续或首尾的分隔符被忽略,不产生空词。`user--_name` → `user` `name`。
2. **小写→大写** — 小写字母到大写字母的转换是边界:`userName` → `user` `Name`;`caféName` → `café` `Name`。中文等「其它字母」(Lo 类)不算小写,**不触发切分**:`用户Name` 是一个词。
3. **连续大写** — 大写串后紧跟小写字母时,最后一个大写字母开新词:`HTTPRequest` → `HTTP` `Request`;`XMLHTTPRequest` → `XMLHTTP` `Request`。大写串在末尾(或后跟数字/分隔符)保持为一个词:`HTTP` → `HTTP`。
4. **数字边界** — 字母↔数字 的转换都是边界,连续数字是一个词:`version2` → `version` `2`;`SHA256` → `SHA` `256`;`ID2Name` → `ID` `2` `Name`;`2fa` → `2` `fa`;`00İ` → `00` `İ`。「数字」仅指 ASCII `0-9`。
5. **其它字符** — emoji、标点等作为词内容保留,不触发边界(`😀A` 是一个词)。

## 输出规则与归一化

- **分隔符风格**(`ToSnake` / `ToKebab` / `ToUpperSnake`)把每词整词小写 / 小写 / 大写,以 `_` / `-` / `_` 连接。
- **`ToCamel`** 把**整个首词**小写(`HTTPRequest` → `httpRequest`,而非 `hTTPRequest`),再首字母大写其余词。
- **`ToPascal`** 每词首字母大写、其余小写,缩写同样归一化:`HTTPRequest` → `HttpRequest`,`ID` → `Id`。

## 幂等性与不变式

`FuzzStrCase` 对纯 ASCII 输入检查下列不变式;非 ASCII 输入只检查五个函数不 panic。Fuzz 测试不等于对所有字符串的数学证明。

| 不变式 | 适用范围 |
| --- | --- |
| `ToSnake(ToSnake(s)) == ToSnake(s)` | ASCII |
| `ToKebab(ToKebab(s)) == ToKebab(s)` | ASCII |
| `ToSnake(ToKebab(s)) == ToSnake(s)` | ASCII |
| `ToKebab(ToSnake(s)) == ToKebab(s)` | ASCII |
| `ToUpperSnake(ToUpperSnake(s)) == ToUpperSnake(s)` | ASCII |
| `ToSnake(ToUpperSnake(s)) == ToSnake(s)` | ASCII |
| `ToCamel(ToSnake(s)) == ToCamel(s)` | ASCII |
| `ToPascal(ToSnake(s)) == ToPascal(s)` | ASCII |

### 限制(文档化,非缺陷)

- **驼峰/帕斯卡不保证幂等。** 它们用大小写表达词边界,对连续单字母词有信息损失:`a_b_c` → 帕斯卡 `ABC` 会被重切分为一个词;`ToCamel("aA_A")` → `"aAA"` → `"aAa"`。对纯 ASCII 输入,分隔符风格互转可保持上述分词结果。
- **非 ASCII 大小写映射可能改变边界。** 例如 `ToSnake("Aϓ")` 得到 `"aϓ"`,再次转换则得到 `"a_ϓ"`;上表不变式不保证对 Unicode 输入成立。
- **不做 Unicode 规范化或标识符校验。** `éName` 得到 `é_name`,而 `e\u0301Name` 得到 `e\u0301name`。只有普通空格、`_`、`-` 是分隔符;制表符、换行和标点保留,也不修复开头的数字。非法 UTF-8 会替换为 U+FFFD,不会逐字节保留。

## 为什么再造一个 strcase?

- **规则定义清晰、有测试、有 fuzz 背书** — 分词器约 60 行,每条规则都有表驱动用例,而不是正则魔法。
- **归一化确定** — 缩写与数字在任何方向都得到一致处理(`HTTPRequest` 与 `HTTP_Request` 都 → `http_request`)。
- **零依赖、go 1.21、MIT** — 与 [go-validator](https://github.com/MouXiaoJun/validator)、[go-copier](https://github.com/MouXiaoJun/copier) 同家族。

## 许可证

[MIT](LICENSE)
