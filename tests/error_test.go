package tests

import (
	"github.com/soranoba/valis"
	"testing"
)

func TestValidationError(t *testing.T) {
	var _ error = &valis.ValidationError{}
}

//func TestValidationError_Error(t *testing.T) {
//	assert := assert.New(t)
//
//	type User struct {
//	}
//	assert.Equal(
//		"(*is.RequiredRule) tests.User{} cannot be blank",
//		valis.NewErrorDetails(is.Required, User{}, errors.New("cannot be blank")).Error(),
//	)
//	assert.Equal(
//		"(*is.RequiredRule) \"\" cannot be blank",
//		valis.NewErrorDetails(is.Required, "", errors.New("cannot be blank")).Error(),
//	)
//}
