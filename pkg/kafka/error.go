package kafka

import (
	"fmt"
)

type ErrorBag struct {
	EventType  string `json:"eventType"`
	Reportable bool   `json:"isReportable"`
	Retryable  bool   `json:"retryable"`
	Cause      error  `json:"cause"`
}

func (eb ErrorBag) Error() string {
	return fmt.Sprintf("err: %s", eb.Cause.Error())
}

func (eb ErrorBag) IsReportable() bool {
	return eb.Reportable
}

func (eb ErrorBag) FormattedText() string {
	return fmt.Sprintf("event: type: %s - message: unhandled err", eb.EventType)
}
