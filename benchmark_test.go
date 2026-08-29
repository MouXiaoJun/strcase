package strcase

import "testing"

// 转换吞吐基准:核心场景是长标识符的 snake 转换(如字段名映射)。
func BenchmarkToSnake(b *testing.B) {
	for i := 0; i < b.N; i++ {
		_ = ToSnake("HTTPRequestWithLongNameAndAPIKey2FA")
	}
}

func BenchmarkToCamel(b *testing.B) {
	for i := 0; i < b.N; i++ {
		_ = ToCamel("http_request_with_long_name_and_api_key_2fa")
	}
}
