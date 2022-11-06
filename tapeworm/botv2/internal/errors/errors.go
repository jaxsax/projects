package errors

import "fmt"

var (
	ErrInvalidDomain = UserFacingError{
		err:          fmt.Errorf("domain is not allowed"),
		userResponse: "Domain is not allowed",
	}
)

type UserFacingError struct {
	err          error
	userResponse string
}

func (e UserFacingError) Error() string {
	return fmt.Sprintf("%s: %s", e.userResponse, e.err)
}

func (e UserFacingError) UserResponse() string {
	return e.userResponse
}
