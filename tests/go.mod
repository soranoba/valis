module github.com/soranoba/valis/tests

go 1.15

require (
	github.com/soranoba/valis v0.0.0
	github.com/soranoba/valis/is v0.0.0
	github.com/soranoba/valis/to v0.0.0-00010101000000-000000000000
	github.com/stretchr/testify v1.7.0
)

replace github.com/soranoba/valis => ../

replace github.com/soranoba/valis/is => ../is

replace github.com/soranoba/valis/to => ../to
