package nrclient

import (
	"time"

	"github.com/newrelic/go-agent/v3/newrelic"
)

const (
	slowQueryThreshold = 100 * time.Millisecond
)

type INewRelicInstance interface {
	Application() *newrelic.Application
}

type newRelicInstance struct {
	application *newrelic.Application
}

type Config struct {
	Key     string
	AppName string
}

func InitNewRelic(config Config) (INewRelicInstance, error) {
	app, err := newrelic.NewApplication(
		newrelic.ConfigAppName(config.AppName),
		newrelic.ConfigLicense(config.Key),
		newrelic.ConfigDistributedTracerEnabled(true),
		newrelic.ConfigAppLogDecoratingEnabled(true),
		func(c *newrelic.Config) {
			c.DatastoreTracer.SlowQuery.Threshold = slowQueryThreshold
		},
	)
	if err != nil {
		return nil, err
	}

	return &newRelicInstance{application: app}, nil
}

func (n *newRelicInstance) Application() *newrelic.Application {
	return n.application
}
