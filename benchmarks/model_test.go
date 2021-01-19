package benchmarks

type Simple struct {
	FirstName string `json:"first_name" validate:"required,max=20"`
	LastName  string `json:"last_name" validate:"required,max=20"`
}
