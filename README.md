# go-strcase

[中文](README_zh.md)

[![Go Reference](https://pkg.go.dev/badge/github.com/MouXiaoJun/strcase.svg)](https://pkg.go.dev/github.com/MouXiaoJun/strcase)
[![Go Version](https://img.shields.io/badge/go-1.21+-00ADD8?style=flat-square&logo=go)](https://golang.org)
[![License](https://img.shields.io/badge/license-MulanPSL--2.0-green.svg?style=flat-square)](LICENSE)

A **zero-dependency** naming-convention conversion library for Go: camelCase, PascalCase, snake_case, kebab-case and UPPER_SNAKE_CASE, with well-defined word-splitting rules for acronyms, digits, consecutive capitals and mixed separators.

## Features

- ✅ **Zero dependencies** — stdlib only (`strings` / `unicode`); clean `go.mod`, no supply-chain risk
- ✅ **5 styles, 1 tokenizer** — every function splits the input into words by the same documented rules, then re-joins with its own casing
- ✅ **Acronym-aware** — `HTTPRequest` → `http_request`, not `httprequest`
- ✅ **Digit boundaries** — `version2` → `version_2`, `SHA256` → `sha_256`, `ID2Name` → `id_2_name`
- ✅ **Mixed separators** — `user--_name` and `USER_NAME` both normalize to `user_name`
- ✅ **Deterministic & idempotent** — separator styles are fixed points (`ToSnake(ToSnake(s)) == ToSnake(s)`) and mutually convertible (`ToSnake(ToKebab(s)) == ToSnake(s)`)
- ✅ **Fuzz-tested invariants** — `FuzzStrCase` verifies no-panic on arbitrary input plus the full invariant suite on ASCII input
- ✅ **Unicode-safe** — Chinese/Latin letters are preserved (`用户Name` stays readable); invalid UTF-8 never panics

## Quick start

```go
package main

import (
	"fmt"

	"github.com/MouXiaoJun/strcase"
)

func main() {
	fmt.Println(strcase.ToSnake("HTTPRequest"))      // http_request
	fmt.Println(strcase.ToCamel("user_name"))        // userName
	fmt.Println(strcase.ToPascal("user_name"))       // UserName
	fmt.Println(strcase.ToKebab("HTTPRequest"))      // http-request
	fmt.Println(strcase.ToUpperSnake("userName"))    // USER_NAME
}
```

## API

All functions take `string` and return `string`. The empty string maps to the empty string.

| Function | Example in | Example out | Notes |
| --- | --- | --- | --- |
| `ToCamel` | `user_name` | `userName` | first word fully lowercased; later words Title-cased |
| `ToPascal` | `user_name` | `UserName` | every word Title-cased (acronyms normalized too) |
| `ToSnake` | `HTTPRequest` | `http_request` | words lowercased, joined with `_` |
| `ToKebab` | `userName` | `user-name` | words lowercased, joined with `-` |
| `ToUpperSnake` | `userName` | `USER_NAME` | words uppercased, joined with `_` |

## Word-splitting rules

Every function splits the input into words first, then re-joins. The rules are deterministic and fully covered by table tests:

1. **Separators** — `_`, `-` and space are word boundaries; consecutive or leading/trailing separators are ignored and never produce empty words. `user--_name` → `user` `name`.
2. **Lower → upper** — a transition from a lowercase letter to an uppercase letter is a boundary: `userName` → `user` `Name`; `caféName` → `café` `Name`. CJK etc. `Lo` letters are *not* "lowercase", so `用户Name` stays **one word**.
3. **Consecutive capitals** — when an uppercase run is followed by a lowercase letter, the **last** capital starts a new word: `HTTPRequest` → `HTTP` `Request`; `XMLHTTPRequest` → `XMLHTTP` `Request`. A run at the end (or before a digit/separator) is one word: `HTTP` → `HTTP`.
4. **Digit boundaries** — every letter↔digit transition is a boundary; consecutive digits form one word: `version2` → `version` `2`; `SHA256` → `SHA` `256`; `ID2Name` → `ID` `2` `Name`; `2fa` → `2` `fa`; `00İ` → `00` `İ`. "Digit" means ASCII `0-9` only.
5. **Other characters** — emoji, punctuation etc. are opaque word content and never trigger boundaries (`😀A` → one word).

## Output rules & normalization

- **Separator styles** (`ToSnake` / `ToKebab` / `ToUpperSnake`) lowercase / lowercase / uppercase each whole word and join with `_` / `-` / `_`.
- **`ToCamel`** lowercases the **whole first word** (so `HTTPRequest` → `httpRequest`, not `hTTPRequest`), then Title-cases the rest.
- **`ToPascal`** Title-cases every word; acronyms are normalized too: `HTTPRequest` → `HttpRequest`, `ID` → `Id`.

## Idempotency & invariants

Proven invariants (asserted in `FuzzStrCase`; hold for any input unless noted):

| Invariant | Scope |
| --- | --- |
| `ToSnake(ToSnake(s)) == ToSnake(s)` | all input |
| `ToKebab(ToKebab(s)) == ToKebab(s)` | all input |
| `ToSnake(ToKebab(s)) == ToSnake(s)` | all input |
| `ToKebab(ToSnake(s)) == ToKebab(s)` | all input |
| `ToUpperSnake(ToUpperSnake(s)) == ToUpperSnake(s)` | ASCII |
| `ToSnake(ToUpperSnake(s)) == ToSnake(s)` | ASCII |
| `ToCamel(ToSnake(s)) == ToCamel(s)` | ASCII |
| `ToPascal(ToSnake(s)) == ToPascal(s)` | ASCII |

### Limitations (documented, not bugs)

- **Camel/Pascal are not idempotent.** They encode word boundaries with capitalization, which is lossy for consecutive single-letter words: `a_b_c` → Pascal `ABC` re-parses as one word. `ToCamel("aA_A")` → `"aAA"` → `"aAa"`. Use the separator styles when round-tripping matters.
- **Non-ASCII case-folding is asymmetric for a few characters** (Turkish `İ` lowercases to `i` which uppercases to `I`; `ϓ`/`ȷ` lack counterparts). For those inputs the fold-based invariants above are not guaranteed — `ToSnake`/`ToKebab` still behave sanely, and no input ever panics.

## Why another strcase?

- **Defined, tested, fuzz-verified rules** instead of regex magic: the tokenizer is ~60 lines and every rule above has table tests.
- **Deterministic normalization**: acronyms and digits get the same treatment in every direction (`HTTPRequest` and `HTTP_Request` both → `http_request`).
- **Zero dependencies**, go 1.21, Mulan PSL v2 — same family as [go-validator](https://github.com/MouXiaoJun/validator) and [go-copier](https://github.com/MouXiaoJun/copier).

## License

[Mulan PSL v2](LICENSE)
