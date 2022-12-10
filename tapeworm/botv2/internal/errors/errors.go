package errors

import "fmt"

var (
	ErrBlocklistedDomain = UserFacingError{
		err:          fmt.Errorf("domain is not allowed"),
		userResponse: "Domain is not allowed",
	}
	ErrInvalidDomain = UserFacingError{
		err:          fmt.Errorf("invalid domain"),
		userResponse: "Domain is not valid",
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
