package model

type MockError struct {
	Message string `json:"message"`
}

func ErrorFrom(err error) *MockError {
	return &MockError{
		Message: err.Error(),
	}
}