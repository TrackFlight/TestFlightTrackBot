package types

const (
	ErrFloodWait            = "FLOODWAIT_X"
	ErrUnauthorized         = "UNAUTHORIZED"
	ErrBadRequest           = "BAD_REQUEST"
	ErrInternalServer       = "INTERNAL_SERVER_ERROR"
	ErrLinkAlreadyFollowing = "LINK_ALREADY_FOLLOWING"
	ErrInvalidLanguageCode  = "INVALID_LANGUAGE_CODE"
)

type ErrorResponse struct {
	Code    string `json:"code"`
	Message string `json:"message"`
	Seconds int    `json:"seconds,omitempty"`
}
