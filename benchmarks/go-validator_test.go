package benchmarks

import (
	"testing"

	"github.com/go-playground/validator/v10"
)

func BenchmarkGoPlayground_SimpleStruct(b *testing.B) {
	validate := validator.New()
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		if err := validate.Struct(&Simple{}); err == nil {
			panic("invalid results")
		}
		if err := validate.Struct(&Simple{FirstName: "123456789012345678901", LastName: "12345678901234568901"}); err == nil {
			panic("invalid results")
		}
	}
}
