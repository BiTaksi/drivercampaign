package main

import (
	"fmt"
	"log"
	"os"
	"os/signal"
	"syscall"

	"github.com/BiTaksi/drivercampaign/configs"
	constants "github.com/BiTaksi/drivercampaign/pkg/constans"
	"github.com/BiTaksi/drivercampaign/pkg/healthcheck"
)

func main() {
	if configErr := initConfig(); configErr != nil {
		log.Fatalf("initialization config: %v", configErr)
	}

	logger := initLogger()
	app, appErr := boot(logger)

	if appErr != nil {
		logger.Fatalf("initialization: %v", appErr)
	}
	defer shutdown(app)

	server := initApplication(app)

	go func() {
		healthcheck.InitHealthCheck()

		if serveErr := server.Listen(fmt.Sprintf(":%s", configs.DriverCampaignApp.Web.Port)); serveErr != nil {
			logger.Fatalf("connection: web server %v", serveErr)
		}
	}()

	// Wait for gracefully shutdown (Interrupt)
	c := make(chan os.Signal, 1)
	signal.Notify(c, syscall.SIGINT, syscall.SIGTERM)

	<-c

	healthcheck.ServerShutdown()
	if shutdownErr := server.Shutdown(); shutdownErr != nil {
		logger.Error(shutdownErr)
	}
}

func shutdown(app *application) {
	app.newRelicInstance.Application().Shutdown(constants.NewRelicShutdownTimeout)

	app.logger.Info("shutdown: completed")
}
