//go:build wireinject
// +build wireinject

package drivercampaign

import (
	"github.com/google/wire"
	"github.com/sirupsen/logrus"

	"github.com/BiTaksi/drivercampaign/internal/handler"
	"github.com/BiTaksi/drivercampaign/internal/route"
	"github.com/BiTaksi/drivercampaign/pkg/requestclient"
)

var handlerProviders = wire.NewSet(
	handler.NewHomeHandler,
)

var allProviders = wire.NewSet(
	handlerProviders,
)

func InitHealthCheckHandler() handler.IHealthCheckHandler {
	wire.Build(handler.NewHealthCheckHandler)
	return nil
}

func InitRoute(
	l *logrus.Logger,
	rqInstance requestclient.IRequestClient,
) route.IRoute {
	wire.Build(allProviders, route.NewRoute)
	return nil
}
