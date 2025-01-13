package validator

type Validator struct {
	Errors map[string]string
}

// method to create a new instance of validator with empty error map
func New() *Validator {
	return &Validator{Errors: make(map[string]string)}
}

// method to check that error map is empty i.e. it does not contain any errors
func (validator *Validator) Valid() bool {
	return len(validator.Errors) == 0
}

// AddError adds an error message to the map (so long as no entry already exists for
// the given key).
func (validator *Validator) AddError(key, message string) {
	if _, exists := validator.Errors[key]; !exists {
		validator.Errors[key] = message
	}
}

// Check adds an error message to the map only if a validation check is not 'ok'.
func (validator *Validator) Check(ok bool, key, message string) {
	if !ok {
		validator.Errors[key] = message
	}
}
