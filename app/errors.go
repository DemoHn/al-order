package app

// Error - customized AppError for displaying error response
type Error struct {
	Name       string
	StatusCode int
	Code       int
	Details    string
}

// ErrIDNotFound - if input ID is empty
func ErrIDNotFound(details string) Error {
	return Error{
		Name:       "ErrIDNotFound",
		StatusCode: 400,
		Code:       10000,
		Details:    details,
	}
}
