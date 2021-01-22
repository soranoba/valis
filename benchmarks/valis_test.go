package benchmarks

import (
	"github.com/soranoba/valis/is"
	valistag "github.com/soranoba/valis/tag"
	"testing"

	"github.com/soranoba/valis"
)

func BenchmarkValis_SimpleStruct_tag(b *testing.B) {
	v := valis.NewValidator()
	v.SetCommonRules()
	b.ResetTimer()

	for i := 0; i < b.N; i++ {
		s := &Simple{}
		if err := v.Validate(s, valistag.Validate); err == nil {
			panic("invalid results")
		}

		s = &Simple{FirstName: "123456789012345678901", LastName: "123456789012345678901"}
		if err := v.Validate(s, valistag.Validate); err == nil {
			panic("invalid results")
		}
	}
}

func BenchmarkValis_SimpleStruct_field(b *testing.B) {
	v := valis.NewValidator()
	v.SetCommonRules()
	b.ResetTimer()

	for i := 0; i < b.N; i++ {
		s := &Simple{}
		if err := v.Validate(
			s,
			valis.Field(&s.FirstName, is.Required, is.LengthBetween(0, 20)),
			valis.Field(&s.LastName, is.Required, is.LengthBetween(0, 20)),
		); err == nil {
			panic("invalid results")
		}

		s = &Simple{FirstName: "123456789012345678901", LastName: "123456789012345678901"}
		if err := v.Validate(
			s,
			valis.Field(&s.FirstName, is.Required, is.LengthBetween(0, 20)),
			valis.Field(&s.LastName, is.Required, is.LengthBetween(0, 20)),
		); err == nil {
			panic("invalid results")
		}
	}
}
