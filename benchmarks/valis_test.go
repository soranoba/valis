package benchmarks

import (
	"testing"

	"github.com/soranoba/valis"
	"github.com/soranoba/valis/is"
	"github.com/soranoba/valis/tagrule"
)

func BenchmarkValis_SimpleStruct_tag(b *testing.B) {
	type Simple struct {
		FirstName string `json:"first_name" validate:"min=1,max=20"`
		LastName  string `json:"last_name" validate:"min=1,max=20"`
	}

	v := valis.NewValidator()
	v.SetCommonRules()
	b.ResetTimer()

	for i := 0; i < b.N; i++ {
		s := &Simple{}
		if err := v.Validate(s, valis.EachFields(tagrule.Validate)); err == nil {
			panic("invalid results")
		}

		s = &Simple{FirstName: "123456789012345678901", LastName: "123456789012345678901"}
		if err := v.Validate(s, valis.EachFields(tagrule.Validate)); err == nil {
			panic("invalid results")
		}
	}
}

func BenchmarkValis_SimpleStruct_field(b *testing.B) {
	type Simple struct {
		FirstName string `json:"first_name"`
		LastName  string `json:"last_name"`
	}

	v := valis.NewValidator()
	v.SetCommonRules()
	b.ResetTimer()

	for i := 0; i < b.N; i++ {
		s := &Simple{}
		if err := v.Validate(
			s,
			valis.Field(&s.FirstName, is.LengthBetween(1, 20)),
			valis.Field(&s.LastName, is.LengthBetween(1, 20)),
		); err == nil {
			panic("invalid results")
		}

		s = &Simple{FirstName: "123456789012345678901", LastName: "123456789012345678901"}
		if err := v.Validate(
			s,
			valis.Field(&s.FirstName, is.LengthBetween(1, 20)),
			valis.Field(&s.LastName, is.LengthBetween(1, 20)),
		); err == nil {
			panic("invalid results")
		}
	}
}
