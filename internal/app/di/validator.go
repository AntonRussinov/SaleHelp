package di

import (
	"errors"
	"salehelp/pkg/extendedvalidator"
)

var ErrInitValidator = errors.New("failed to initialize validator")

func (prov *Provider) InitValidator() error {
	validator, err := extendedvalidator.New()
	if err != nil {
		return err
	}

	prov.Validator = validator

	return nil
}
