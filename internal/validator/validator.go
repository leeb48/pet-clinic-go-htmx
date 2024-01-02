package validator

type Validator struct {
	FieldErrors map[string]string
}

func (validator *Validator) Valid() bool {
	return len(validator.FieldErrors) == 0
}

func (validator *Validator) AddFieldError(key, message string) {
	if validator.FieldErrors == nil {
		validator.FieldErrors = map[string]string{}
	}

	if _, exists := validator.FieldErrors[key]; !exists {
		validator.FieldErrors[key] = message
	}
}

func (validator *Validator) CheckField(ok bool, key, message string) {
	if !ok {
		validator.AddFieldError(key, message)
	}
}
