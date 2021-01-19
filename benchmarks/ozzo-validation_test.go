package benchmarks

import (
	validation "github.com/go-ozzo/ozzo-validation/v4"
	"testing"
)

func BenchmarkOzzo_SimpleStruct(b *testing.B) {
	for i := 0; i < b.N; i++ {
		s := &Simple{}
		if err := validation.ValidateStruct(
			s,
			validation.Field(&s.FirstName, validation.Required, validation.Length(0, 20)),
			validation.Field(&s.LastName, validation.Required, validation.Length(0, 20)),
		); err == nil {
			panic("invalid results")
		}

		s = &Simple{FirstName: "123456789012345678901", LastName: "123456789012345678901"}
		if err := validation.ValidateStruct(
			s,
			validation.Field(&s.FirstName, validation.Required, validation.Length(0, 20)),
			validation.Field(&s.LastName, validation.Required, validation.Length(0, 20)),
		); err == nil {
			panic("invalid results")
		}
	}
}
