package strcase

import (
	"strings"
	"testing"
)

// FuzzStrCase 验证:
//
// 任意 Unicode 输入(含非法 UTF-8 字节):
//   - 五个转换函数都不 panic。
//
// 纯 ASCII 输入(大小写折叠无信息损失,且所有大写字母都有小写形式):
//   - 分隔符风格(snake/kebab/大写蛇形)幂等。
//   - snake 是词切分的规范形:ToSnake(ToKebab(s)) == ToSnake(s)、
//     ToSnake(ToUpperSnake(s)) == ToSnake(s)、
//     ToCamel(ToSnake(s)) == ToCamel(s)、ToPascal(ToSnake(s)) == ToPascal(s)
//     等交叉不变式。
//
// 不在 fuzz 中断言的限制(见 README「限制」):
//   - 驼峰/帕斯卡用大小写表达词边界,连续单字母词会合并
//     (a_b_c → 帕斯卡 ABC),故二者不保证幂等,
//     也不保证 ToSnake(ToCamel(s)) == ToSnake(s)。
//   - 个别非 ASCII 字符大小写映射不对称(土耳其 İ 小写为 i 再大写为 I;
//     ϓ、ȷ 等无对应大小写形式),折叠会改变重切分的边界,
//     折叠类不变式对非 ASCII 输入不成立。
func FuzzStrCase(f *testing.F) {
	f.Add("userName")
	f.Add("HTTPRequest")
	f.Add("XMLHTTPRequest")
	f.Add("user_name")
	f.Add("USER_NAME")
	f.Add("user--_name")
	f.Add("version2")
	f.Add("SHA256")
	f.Add("ID2Name")
	f.Add("2fa")
	f.Add("用户Name")
	f.Add("a_b_c")
	f.Add("aA")
	f.Add("00İ")
	f.Add("用a00000")
	f.Add("caféName")
	f.Add("İstanbul")
	f.Add("Aϓ")
	f.Add("e\u0301Name")
	f.Add("a\xffB")
	f.Add("")
	f.Add("___")
	f.Add(strings.Repeat("ABC_", 20))
	f.Add(strings.Repeat("a", 10000))
	f.Fuzz(func(t *testing.T, s string) {
		// 任意输入:不 panic
		_ = ToCamel(s)
		_ = ToPascal(s)
		_ = ToSnake(s)
		_ = ToKebab(s)
		_ = ToUpperSnake(s)

		// 折叠类不变式仅对纯 ASCII 输入保证
		if !isASCIIOnly(s) {
			return
		}
		snake := ToSnake(s)
		kebab := ToKebab(s)
		upper := ToUpperSnake(s)

		// 分隔符风格幂等
		if ToSnake(snake) != snake {
			t.Fatalf("ToSnake not idempotent for %q: %q", s, snake)
		}
		if ToKebab(kebab) != kebab {
			t.Fatalf("ToKebab not idempotent for %q: %q", s, kebab)
		}
		if ToUpperSnake(upper) != upper {
			t.Fatalf("ToUpperSnake not idempotent for %q: %q", s, upper)
		}

		// snake 是规范形
		if ToSnake(kebab) != snake {
			t.Fatalf("snake(kebab(%q)) = %q != %q", s, ToSnake(kebab), snake)
		}
		if ToKebab(snake) != kebab {
			t.Fatalf("kebab(snake(%q)) = %q != %q", s, ToKebab(snake), kebab)
		}
		if ToSnake(upper) != snake {
			t.Fatalf("snake(upper(%q)) = %q != %q", s, ToSnake(upper), snake)
		}
		if ToUpperSnake(snake) != upper {
			t.Fatalf("upper(snake(%q)) = %q != %q", s, ToUpperSnake(snake), upper)
		}
		if ToCamel(snake) != ToCamel(s) {
			t.Fatalf("camel(snake(%q)) mismatch", s)
		}
		if ToPascal(snake) != ToPascal(s) {
			t.Fatalf("pascal(snake(%q)) mismatch", s)
		}
	})
}

// isASCIIOnly 报告 s 是否只含 ASCII 字符(其大小写折叠无信息损失)。
func isASCIIOnly(s string) bool {
	for i := 0; i < len(s); i++ {
		if s[i] >= 0x80 {
			return false
		}
	}
	return true
}
