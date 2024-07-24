package domain

type (
	ServiceError struct {
		Fn  string
		Err error
		Msg string
	}

	NotFoundError struct {
		Entity     string
		ParamName  string
		ParamValue string
		Err        error
	}

	NotFoundErrorNew struct {
		ServiceError
		Entity     string
		ParamName  string
		ParamValue string
	}

	ForbiddenError struct {
		Msg string
		Err error
	}

	DuplicateError struct {
		Msg   string
		Field string
		Err   error
	}

	ApplicationLogicError struct {
		Msg string
	}

	ValidationErrors []FieldValidationError

	FieldValidationError struct {
		Field   string `json:"field"`
		Message string `json:"message"`
		Details string `json:"details"`
	}
)

func (e *NotFoundError) Error() string {
	errMsg := ""
	if e.Err != nil {
		errMsg = e.Err.Error()
	}

	return e.Entity + ": " + "not found by " + e.ParamName + " = " + e.ParamValue + ": " + errMsg
}

func (e *NotFoundError) Unwrap() error { return e.Err }

func (e *ServiceError) Error() string {
	errMsg := ""
	if e.Err != nil {
		errMsg = e.Err.Error()
	}

	return e.Fn + ": " + errMsg
}

func (e *ServiceError) Unwrap() error { return e.Err }

func (e *ForbiddenError) Error() string {
	errMsg := ""
	if e.Err != nil {
		errMsg = e.Err.Error()
	}

	return e.Msg + ": " + errMsg
}

func (e *ForbiddenError) Unwrap() error { return e.Err }

func NewNotFoundError(fn, entity, paramName, paramValue string, err error) *ServiceError {
	return &ServiceError{
		Fn: fn,
		Err: &NotFoundError{
			Entity:     entity,
			ParamName:  paramName,
			ParamValue: paramValue,
			Err:        err,
		},
	}
}

func NewNotFoundErrorNew(fn, entity, paramName, paramValue string, err error) *NotFoundErrorNew {
	return &NotFoundErrorNew{
		ServiceError: ServiceError{
			Fn:  fn,
			Err: err,
			Msg: "not found",
		},
		Entity:     entity,
		ParamName:  paramName,
		ParamValue: paramValue,
	}
}

func (ValidationErrors) Error() string {
	return "validation error"
}

func (e *DuplicateError) Error() string {
	errMsg := ""
	if e.Err != nil {
		errMsg = e.Err.Error()
	}

	return e.Msg + ": " + errMsg
}

func (e *DuplicateError) Unwrap() error { return e.Err }

func NewForbiddenError(fn, msg string, err error) *ServiceError {
	return &ServiceError{
		Fn: fn,
		Err: &ForbiddenError{
			Msg: msg,
			Err: err,
		},
	}
}

func NewValidationErrors(fn, msg string, errs []FieldValidationError) *ServiceError {
	return &ServiceError{
		Fn:  fn,
		Msg: msg,
		Err: ValidationErrors(errs),
	}
}

func NewDuplicateError(fn, msg, field string, err error) *ServiceError {
	return &ServiceError{
		Fn: fn,
		Err: &DuplicateError{
			Msg:   msg,
			Field: field,
			Err:   err,
		},
	}
}

func NewApplicationLogicError(fn, msg string) *ServiceError {
	return &ServiceError{
		Fn: fn,
		Err: &ApplicationLogicError{
			Msg: msg,
		},
	}
}

func (e *ApplicationLogicError) Error() string {
	return e.Msg
}
