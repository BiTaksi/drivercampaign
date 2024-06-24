package httpclient

import (
	"fmt"
)

type ResponseErrorBag struct {
	Code    int    `json:"returnCode"`
	Message string `json:"returnMessage"`
	Cause   error  `json:"cause"`

	Response
}

type ResponseInternalErrorBag struct {
	Err   Error `json:"error"`
	Cause error `json:"cause"`

	Response
}

type Error struct {
	Code    string `json:"code"`
	Message string `json:"message"`
}

func (re ResponseErrorBag) Error() string {
	return fmt.Sprintf("status: %d - body: %s - err: %v", re.StatusCode, string(re.Response.Body), re.Cause.Error())
}

func (re ResponseInternalErrorBag) Error() string {
	return fmt.Sprintf("status: %d - body: %s - err: %v", re.StatusCode, string(re.Response.Body), re.Cause.Error())
}
