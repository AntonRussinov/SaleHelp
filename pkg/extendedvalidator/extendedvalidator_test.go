package extendedvalidator_test

import (
	"salehelp/pkg/extendedvalidator"
	"testing"

	"github.com/stretchr/testify/require"
)

type (
	invalidField struct {
		Field string `validate:"is_valid_enum"`
	}

	noValidateFunc struct {
		Field enumWithNoValidateFunc `validate:"is_valid_enum"`
	}

	wrongValidateFuncParamsCount struct {
		Field enumWithWrongValidateFuncParamsCount `validate:"is_valid_enum"`
	}

	wrongValidateFuncParamType struct {
		Field enumWithWrongValidateFuncParamType `validate:"is_valid_enum"`
	}

	correctValidateFunc struct {
		Field enum `validate:"is_valid_enum"`
	}

	enum                                 string
	enumWithNoValidateFunc               string
	enumWithWrongValidateFuncParamsCount string
	enumWithWrongValidateFuncParamType   string
)

const (
	enumOne enum = "One"
	enumTwo enum = "Two"
)

func TestExtendedValidator(t *testing.T) {
	t.Parallel()

	validator, err := extendedvalidator.New()
	require.NoError(t, err)

	err = validator.Validate(invalidField{})
	require.Error(t, err)

	err = validator.Validate(noValidateFunc{})
	require.Error(t, err)

	err = validator.Validate(wrongValidateFuncParamsCount{})
	require.Error(t, err)

	err = validator.Validate(wrongValidateFuncParamType{})
	require.Error(t, err)

	err = validator.Validate(correctValidateFunc{
		Field: "Some invalid value",
	})
	require.Error(t, err)

	err = validator.Validate(correctValidateFunc{
		Field: enumOne,
	})
	require.NoError(t, err)
}

func (en enumWithWrongValidateFuncParamsCount) IsValid() (bool, string, int) {
	return true, "", 0
}

func (en enumWithWrongValidateFuncParamType) IsValid() string {
	return "invalid value"
}

func (en enum) IsValid() bool {
	return en == enumOne || en == enumTwo
}
