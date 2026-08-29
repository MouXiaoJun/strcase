// Package strcase 提供零依赖的命名风格转换:小驼峰 / 大驼峰 / 蛇形 /
// 短横线 / 大写蛇形。
//
// 输入支持任意组合:驼峰(userName)、帕斯卡(UserName)、蛇形(user_name)、
// 短横线(user-name)、大写蛇形(USER_NAME)、混合分隔符(user--_name)、
// 连续大写缩写(HTTPRequest)、数字边界(version2)、非 ASCII(中文/emoji)。
//
// # 单词切分规则
//
//  1. 分隔符:下划线 `_`、连字符 `-`、空格都是单词边界;连续或首尾的
//     分隔符被忽略,不产生空词(user--_name → user|name)。
//  2. 小写→大写:小写字母(Ll,含变音符如 é)到大写字母的转换是边界:
//     userName → user|Name;caféName → café|Name。
//     中文等「其它字母」(Lo)不算小写,不触发切分:用户Name → 一个词。
//  3. 连续大写:一段大写串后紧跟小写字母时,最后一个大写字母开新词:
//     HTTPRequest → HTTP|Request;XMLHTTPRequest → XMLHTTP|Request;
//     大写串在末尾(或后面是数字/分隔符)则保持为一个词:HTTP → HTTP。
//  4. 数字边界:字母↔数字 的转换都是边界(任意书写系统的字母),连续数字
//     是一个词:version2 → version|2;SHA256 → SHA|256;ID2Name → ID|2|Name;
//     2fa → 2|fa;用户2 → 用户|2;00İ → 00|İ。数字仅指 ASCII 数字(0-9)。
//  5. 其它字符(中文、emoji、标点等)作为词内容保留;emoji/标点不参与
//     边界判断(😀A → 一个词)。
//
// # 输出规则
//
//   - ToSnake / ToKebab / ToUpperSnake:每词整词转小写 / 小写 / 大写,
//     以 `_` / `-` / `_` 连接。分隔符风格是保真的:幂等
//     (ToSnake(ToSnake(s)) == ToSnake(s)),且互相可转
//     (ToSnake(ToKebab(s)) == ToSnake(s))。
//   - ToCamel:首词整词小写,后续词首字母大写、其余小写(缩写同样归一化):
//     user_name → userName;HTTPRequest → httpRequest。
//   - ToPascal:每词首字母大写、其余小写(缩写同样归一化):
//     user_name → UserName;HTTPRequest → HttpRequest。
//
// 驼峰 / 帕斯卡用大小写表达词边界,存在固有限制:连续单字母词会合并
// (a_b_c → 帕斯卡 ABC,重切分为一个词 ABC),因此二者不保证幂等;
// 词切分的规范形是 snake:对 ASCII 输入 ToCamel(ToSnake(s)) == ToCamel(s)、
// ToPascal(ToSnake(s)) == ToPascal(s) 恒成立(见 fuzz 测试)。
package strcase

import (
	"strings"
	"unicode"
)

// isSeparator 报告 r 是否为单词分隔符(下划线 / 连字符 / 空格)。
func isSeparator(r rune) bool {
	return r == '_' || r == '-' || r == ' '
}

// isASCIIDigit 仅把 ASCII 数字(0-9)当作数字边界参与者(规则 4)。
func isASCIIDigit(r rune) bool { return r >= '0' && r <= '9' }

// splitWords 按文档规则把输入切成单词列表,不含空词。
func splitWords(s string) [][]rune {
	runes := []rune(s)
	var words [][]rune
	var cur []rune
	flush := func() {
		if len(cur) > 0 {
			words = append(words, cur)
			cur = nil
		}
	}
	for i := 0; i < len(runes); i++ {
		r := runes[i]
		switch {
		case isSeparator(r):
			flush() // 规则 1:分隔符是边界,连续分隔符被忽略
		case isASCIIDigit(r):
			// 规则 4:字母→数字 是边界。
			if n := len(cur); n > 0 && unicode.IsLetter(cur[n-1]) {
				flush()
			}
			cur = append(cur, r)
		case unicode.IsLetter(r):
			// 规则 4:数字→字母 是边界。
			if n := len(cur); n > 0 && isASCIIDigit(cur[n-1]) {
				flush()
			}
			if !unicode.IsUpper(r) {
				cur = append(cur, r)
				continue
			}
			// 规则 2:小写字母→大写字母 是边界(userName、caféName)。
			// 中文等 Lo 字母不算小写,不切分(用户Name → 一个词)。
			if n := len(cur); n > 0 && unicode.IsLower(cur[n-1]) {
				flush()
			}
			// 规则 3:连续大写后跟小写时,最后一个大写字母开新词。
			if n := len(cur); n > 0 && unicode.IsUpper(cur[n-1]) &&
				i+1 < len(runes) && unicode.IsLower(runes[i+1]) {
				flush()
			}
			cur = append(cur, r)
		default:
			// 规则 5:emoji、标点等非字母字符作为词内容保留。
			cur = append(cur, r)
		}
	}
	flush()
	return words
}

// ToSnake 转蛇形:userName → user_name;HTTPRequest → http_request。
func ToSnake(s string) string {
	return joinWords(splitWords(s), "_", strings.ToLower)
}

// ToKebab 转短横线:userName → user-name;HTTPRequest → http-request。
func ToKebab(s string) string {
	return joinWords(splitWords(s), "-", strings.ToLower)
}

// ToUpperSnake 转大写蛇形:userName → USER_NAME;HTTPRequest → HTTP_REQUEST。
func ToUpperSnake(s string) string {
	return joinWords(splitWords(s), "_", strings.ToUpper)
}

// ToCamel 转小驼峰:首词整词小写,后续词首字母大写、其余小写。
// user_name → userName;User_Name → userName;HTTPRequest → httpRequest。
func ToCamel(s string) string {
	words := splitWords(s)
	if len(words) == 0 {
		return ""
	}
	var b strings.Builder
	b.Grow(len(s))
	b.WriteString(strings.ToLower(string(words[0])))
	for _, w := range words[1:] {
		b.WriteString(titleWord(w))
	}
	return b.String()
}

// ToPascal 转大驼峰:每词首字母大写、其余小写(缩写同样归一化)。
// user_name → UserName;HTTPRequest → HttpRequest。
func ToPascal(s string) string {
	words := splitWords(s)
	var b strings.Builder
	b.Grow(len(s))
	for _, w := range words {
		b.WriteString(titleWord(w))
	}
	return b.String()
}

// titleWord 首字母大写、其余小写。
func titleWord(w []rune) string {
	if len(w) == 0 {
		return ""
	}
	return string(unicode.ToUpper(w[0])) + strings.ToLower(string(w[1:]))
}

// joinWords 把词表以 sep 连接,每词经 mapFn 变换后输出。
func joinWords(words [][]rune, sep string, mapFn func(string) string) string {
	if len(words) == 0 {
		return ""
	}
	parts := make([]string, 0, len(words))
	for _, w := range words {
		parts = append(parts, mapFn(string(w)))
	}
	return strings.Join(parts, sep)
}
