package utils

// System errors
const (
	NotFoundErrCode     = "not_found"
	ValidationErrCode   = "validation_failed"
	UnexpectedErrCode   = "unexpected_error"
	UnauthorizedErrCode = "unauthorized"
	BodyParserErrCode   = "body_parser_failed"
	QueryParserErrCode  = "query_parser_failed"

	NotFoundMsg     = "Not found!"
	UnexpectedMsg   = "An unexpected error has occurred."
	ValidationMsg   = "the given data was invalid"
	UnauthorizedMsg = "Authentication failed."
	BodyParserMsg   = "The given values could not be parsed."
	QueryParserMsg  = "The given values on query could not be parsed."
)

type ErrorBag struct {
	Code  string `json:"code"`
	Cause error  `json:"cause"`
}

func (e ErrorBag) GetCode() string {
	return e.Code
}

func (e ErrorBag) Error() string {
	return e.Cause.Error()
}
