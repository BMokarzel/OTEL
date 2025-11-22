package errors

type ProviderError struct {
	GenericError
	Provider string `json:"provider"`
}

/*
func (e *ProviderError) Error() string {
	return fmt.Sprintf("provider error: %s", e.Message)
}

func NewProviderError(provider string) *ProviderError {
	return &ProviderError{
		GenericError{
			Message: fmt.Sprintf("[%s] -  Provider", provider),
		},
		provider,
	}
}

func NewProviderErrorMsg(code, msg string) *ProviderError {
	return &ProviderError{

	}
}

func AsProviderError(err error) (bool, ProviderError) {
	e, ok := err.(*ProviderError)
	if !ok {
		return false, ProviderError{}
	}
	return ok, *e
}

func IsProviderError(err error) bool {
	_, ok := err.(*ProviderError)
	if !ok {
		return false
	}
	return ok
}
*/
