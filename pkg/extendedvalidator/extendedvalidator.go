package extendedvalidator

import (
	"fmt"
	"reflect"

	"github.com/go-playground/validator"
)

type Validator struct {
	Validator *validator.Validate
}

func MustNew() *Validator {
	v, err := New()
	if err != nil {
		panic(err)
	}

	return v
}

func New() (*Validator, error) {
	validate := validator.New()

	err := validate.RegisterValidation("is_valid_enum", func(fl validator.FieldLevel) bool {
		fieldValue := fl.Field()

		fieldType := fieldValue.Type()

		method, ok := fieldType.MethodByName("IsValid")
		if !ok {
			return false
		}

		output := method.Func.Call([]reflect.Value{fieldValue})
		if len(output) != 1 {
			return false
		}

		value := output[0].Interface()

		boolValue, ok := value.(bool)
		if !ok {
			return false
		}

		return boolValue
	})
	if err != nil {
		return nil, err
	}

	return &Validator{Validator: validate}, nil
}

func (vl *Validator) Validate(i any) error {
	const op = "app.Validate"

	err := vl.Validator.Struct(i)
	if err != nil {
		return fmt.Errorf("validation failed %s: %w", op, err)
	}

	return nil
}
