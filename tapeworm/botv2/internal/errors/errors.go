package errors

import "fmt"

var (
	ErrBlocklistedDomain = UserFacingError{
		err:             fmt.Errorf("domain is not allowed"),
		stickerResponse: "CAACAgQAAxkBAAIBlmOZmV1PDFKylJ9aWowPvB_wXKFmAAL7CwACWkWRU7PGcGTSdnIuLAQ",
	}
	ErrInvalidDomain = UserFacingError{
		err:          fmt.Errorf("invalid domain"),
		userResponse: "Domain is not valid",
	}
)

type UserFacingError struct {
	err          error
	userResponse string

	// file id of the sticker
	stickerResponse string
}

func (e UserFacingError) Error() string {
	return fmt.Sprintf("%s: %s", e.userResponse, e.err)
}

func (e UserFacingError) UserResponse() string {
	return e.userResponse
}

func (e UserFacingError) StickerResponse() string {
	return e.stickerResponse
}
