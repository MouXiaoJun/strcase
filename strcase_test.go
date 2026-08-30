package strcase

import "testing"

func TestUnicodeSeparatorBoundary(t *testing.T) {
	first := ToSnake("Aϓ")
	if second := ToSnake(first); first != "aϓ" || second != "a_ϓ" {
		t.Fatalf("Unicode mapping boundary: first=%q second=%q", first, second)
	}
}

func TestToSnakeInputBoundaries(t *testing.T) {
	for _, tc := range []struct{ in, want string }{
		{"éName", "é_name"},
		{"e\u0301Name", "e\u0301name"},
		{"a b", "a_b"},
		{"a\tb", "a\tb"},
		{"a\nb", "a\nb"},
		{"a.b", "a.b"},
		{"2fa", "2_fa"},
		{"a\xffB", "a\ufffdb"},
	} {
		if got := ToSnake(tc.in); got != tc.want {
			t.Errorf("ToSnake(%q) = %q, want %q", tc.in, got, tc.want)
		}
	}
}

// 表驱动:每个函数独立用例,覆盖文档化的全部切分规则:
// 分隔符混合/连续、小写→大写、连续大写缩写、数字边界、非 ASCII、空串。
func TestToSnake(t *testing.T) {
	cases := []struct {
		in, want string
	}{
		// 规格示例
		{"userName", "user_name"},
		{"UserName", "user_name"},
		{"HTTPRequest", "http_request"},
		{"version2", "version_2"},
		// 输入风格:蛇形/短横线/大写蛇形/混合分隔符
		{"user_name", "user_name"},
		{"user-name", "user_name"},
		{"user name", "user_name"},
		{"USER_NAME", "user_name"},
		{"user--_name", "user_name"},
		{"-user_Name-", "user_name"},
		{"__user____name__", "user_name"},
		// 连续大写缩写
		{"HTTPServer", "http_server"},
		{"XMLHTTPRequest", "xmlhttp_request"},
		{"ID", "id"},
		{"APIKey", "api_key"},
		{"UserID", "user_id"},
		{"A", "a"},
		{"ABc", "a_bc"},
		{"aBC", "a_bc"},
		{"AbC", "ab_c"},
		{"UsEr", "us_er"},
		// 数字边界
		{"Version2", "version_2"},
		{"SHA256", "sha_256"},
		{"UTF8", "utf_8"},
		{"ID2Name", "id_2_name"},
		{"2fa", "2_fa"},
		{"v2ray", "v_2_ray"},
		{"A2B", "a_2_b"},
		{"user_name_2", "user_name_2"},
		{"123", "123"},
		// 组合
		{"snake_case2HTTP", "snake_case_2_http"},
		{"Foo-Bar_2baz", "foo_bar_2_baz"},
		// 连续单字母词:分隔符保真,snake 不丢信息
		{"a_b_c", "a_b_c"},
		{"aA", "a_a"},
		// 非 ASCII:中英混排按边界切分,emoji 是词内容
		{"用户Name", "用户name"}, // Lo 字母不算小写,不切分
		{"用户_Name", "用户_name"},
		{"中文_English", "中文_english"},
		{"😀A", "😀a"},
		// 大小写折叠不破坏幂等(İ → i)
		{"00İ", "00_i"},
		{"İstanbul", "istanbul"},  // 单个大写开头不切分(Request 同理)
		{"caféName", "café_name"}, // 变音符小写字母触发边界
		{"用a00000", "用a_00000"},   // 折叠不产生新边界,幂等
		{"İ_stanbul", "i_stanbul"},
		// 空串 / 纯分隔符
		{"", ""},
		{"___", ""},
		{"---", ""},
		{" - _ ", ""},
	}
	for _, c := range cases {
		if got := ToSnake(c.in); got != c.want {
			t.Errorf("ToSnake(%q) = %q, want %q", c.in, got, c.want)
		}
	}
}

func TestToKebab(t *testing.T) {
	cases := []struct {
		in, want string
	}{
		{"userName", "user-name"},
		{"HTTPRequest", "http-request"},
		{"user_name", "user-name"},
		{"USER_NAME", "user-name"},
		{"user--_name", "user-name"},
		{"version2", "version-2"},
		{"SHA256", "sha-256"},
		{"用户Name", "用户name"},
		{"00İ", "00-i"},
		{"", ""},
		{"_", ""},
	}
	for _, c := range cases {
		if got := ToKebab(c.in); got != c.want {
			t.Errorf("ToKebab(%q) = %q, want %q", c.in, got, c.want)
		}
	}
}

func TestToUpperSnake(t *testing.T) {
	cases := []struct {
		in, want string
	}{
		{"userName", "USER_NAME"},
		{"HTTPRequest", "HTTP_REQUEST"},
		{"user_name", "USER_NAME"},
		{"user-name", "USER_NAME"},
		{"version2", "VERSION_2"},
		{"SHA256", "SHA_256"},
		{"ID2Name", "ID_2_NAME"},
		{"2fa", "2_FA"},
		{"ID", "ID"},
		{"aBC", "A_BC"},
		{"用户Name", "用户NAME"},
		{"用a00000", "用A_00000"}, // 折叠后重切分稳定
		{"", ""},
		{"_", ""},
	}
	for _, c := range cases {
		if got := ToUpperSnake(c.in); got != c.want {
			t.Errorf("ToUpperSnake(%q) = %q, want %q", c.in, got, c.want)
		}
	}
}

func TestToCamel(t *testing.T) {
	cases := []struct {
		in, want string
	}{
		// 规格示例:保留首字母小写
		{"user_name", "userName"},
		{"UserName", "userName"},
		{"USER_NAME", "userName"},
		{"user-name", "userName"},
		{"user name", "userName"},
		{"user--_name", "userName"},
		// 连续大写缩写:首词整词小写
		{"HTTPRequest", "httpRequest"},
		{"HTTPServer", "httpServer"},
		{"APIKey", "apiKey"},
		{"ID", "id"},
		// 数字边界
		{"version2", "version2"},
		{"user_name_2", "userName2"},
		{"ID2Name", "id2Name"},
		{"SHA256", "sha256"},
		{"2fa", "2Fa"},
		// 奇形输入
		{"ABc", "aBc"},
		{"aBC", "aBc"},
		{"UsEr", "usEr"},
		// 连续单字母词:驼峰输出合并(信息损失,见 README 限制)
		{"a_b_c", "aBC"},
		{"aA", "aA"},
		// 非 ASCII
		{"用户Name", "用户name"},
		{"用户_Name", "用户Name"},
		// 空串
		{"", ""},
		{"_", ""},
	}
	for _, c := range cases {
		if got := ToCamel(c.in); got != c.want {
			t.Errorf("ToCamel(%q) = %q, want %q", c.in, got, c.want)
		}
	}
}

func TestToPascal(t *testing.T) {
	cases := []struct {
		in, want string
	}{
		// 规格示例
		{"user_name", "UserName"},
		{"userName", "UserName"},
		{"USER_NAME", "UserName"},
		// 缩写归一化:每词首字母大写其余小写
		{"HTTPRequest", "HttpRequest"},
		{"HTTPServer", "HttpServer"},
		{"APIKey", "ApiKey"},
		{"ID", "Id"},
		{"XMLHTTPRequest", "XmlhttpRequest"},
		// 数字边界
		{"version2", "Version2"},
		{"user_name_2", "UserName2"},
		{"ID2Name", "Id2Name"},
		{"SHA256", "Sha256"},
		{"2fa", "2Fa"},
		// 奇形输入
		{"ABc", "ABc"},
		{"aBC", "ABc"},
		{"AbC", "AbC"},
		// 连续单字母词:帕斯卡输出合并为一个词(见 README 限制)
		{"a_b_c", "ABC"},
		{"aA", "AA"},
		// 非 ASCII
		{"用户Name", "用户name"}, // 用户Name 一个词
		{"用户_Name", "用户Name"},
		// 空串
		{"", ""},
		{"_", ""},
	}
	for _, c := range cases {
		if got := ToPascal(c.in); got != c.want {
			t.Errorf("ToPascal(%q) = %q, want %q", c.in, got, c.want)
		}
	}
}

// 交叉一致性:同一输入在五个函数间的关系固定(也由 fuzz 全量验证)。
func TestCrossFunctionConsistency(t *testing.T) {
	cases := []struct {
		in, snake, kebab, upperSnake, camel, pascal string
	}{
		{"userName", "user_name", "user-name", "USER_NAME", "userName", "UserName"},
		{"HTTPRequest", "http_request", "http-request", "HTTP_REQUEST", "httpRequest", "HttpRequest"},
		{"user_name", "user_name", "user-name", "USER_NAME", "userName", "UserName"},
		{"version2", "version_2", "version-2", "VERSION_2", "version2", "Version2"},
		{"ID2Name", "id_2_name", "id-2-name", "ID_2_NAME", "id2Name", "Id2Name"},
		{"2fa", "2_fa", "2-fa", "2_FA", "2Fa", "2Fa"},
		{"", "", "", "", "", ""},
	}
	for _, c := range cases {
		if got := ToSnake(c.in); got != c.snake {
			t.Errorf("ToSnake(%q) = %q, want %q", c.in, got, c.snake)
		}
		if got := ToKebab(c.in); got != c.kebab {
			t.Errorf("ToKebab(%q) = %q, want %q", c.in, got, c.kebab)
		}
		if got := ToUpperSnake(c.in); got != c.upperSnake {
			t.Errorf("ToUpperSnake(%q) = %q, want %q", c.in, got, c.upperSnake)
		}
		if got := ToCamel(c.in); got != c.camel {
			t.Errorf("ToCamel(%q) = %q, want %q", c.in, got, c.camel)
		}
		if got := ToPascal(c.in); got != c.pascal {
			t.Errorf("ToPascal(%q) = %q, want %q", c.in, got, c.pascal)
		}
	}
}
