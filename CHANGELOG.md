# Changelog

All notable changes to this project will be documented in this file.

The format is based on [Keep a Changelog](https://keepachangelog.com/en/1.0.0/),
and this project adheres to [Semantic Versioning](https://semver.org/spec/v2.0.0.html).

## [1.0.1] - 2026-08-30

### Fixed / Changed

- Document the distinction between ASCII invariants and Unicode casing behavior.
- Add fixed Unicode examples and fuzz seeds without changing the conversion algorithm.
- Align LICENSE and current documentation with MIT, as confirmed by the maintainer; preserve existing copyright notices.

## [1.0.0] - 2026-02-XX

### Added
- 初始版本发布
- 5 个转换函数:ToCamel / ToPascal / ToSnake / ToKebab / ToUpperSnake
- 文档化的单词切分规则:分隔符混合、小写→大写、连续大写缩写、数字边界、非 ASCII
- 表驱动测试覆盖全部规则 + FuzzStrCase(任意输入不 panic,ASCII 输入验证幂等与规范形不变式)
- 零依赖,go 1.21,Mulan PSL v2
