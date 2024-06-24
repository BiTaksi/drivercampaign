package main

import (
	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/requestid"
	"github.com/nicksnyder/go-i18n/v2/i18n"
	"github.com/sirupsen/logrus"

	dc "github.com/BiTaksi/drivercampaign"
	"github.com/BiTaksi/drivercampaign/internal/middleware"
	"github.com/BiTaksi/drivercampaign/pkg/response"
	"github.com/BiTaksi/drivercampaign/pkg/utils"
	"github.com/BiTaksi/drivercampaign/pkg/validation"

	"github.com/BiTaksi/drivercampaign/pkg/nrclient"
	"github.com/BiTaksi/drivercampaign/pkg/requestclient"
	"github.com/BiTaksi/drivercampaign/pkg/tokenizer"
)

type application struct {
	logger                *logrus.Logger
	languageBundle        *i18n.Bundle
	tokenizerInstance     tokenizer.ITokenizer
	newRelicInstance      nrclient.INewRelicInstance
	requestClientInstance requestclient.IRequestClient
}

func initApplication(a *application) *fiber.App {
	app := fiber.New(fiber.Config{
		// Override default error handler - Internal server err
		ErrorHandler: func(c *fiber.Ctx, err error) error {
			code := fiber.StatusInternalServerError
			if e, ok := err.(*fiber.Error); ok {
				code = e.Code
			}

			errBag := utils.ErrorBag{Code: utils.UnexpectedErrCode}

			return c.Status(code).JSON(response.NewErrorResponse(c.Context(), errBag, utils.UnexpectedMsg))
		},
	})

	// Health check routes
	a.addHealthCheckRoutes(app)

	// Common middleware
	a.addCommonMiddleware(app)

	// 404 handler
	app.Use(func(c *fiber.Ctx) error {
		errBag := utils.ErrorBag{Code: utils.NotFoundErrCode}

		return c.Status(fiber.StatusNotFound).JSON(response.NewErrorResponse(c.Context(), errBag, utils.NotFoundMsg))
	})

	return app
}

func (a *application) addCommonMiddleware(app *fiber.App) {
	app.Use(middleware.RecoverMiddleware(a.logger))
	app.Use(requestid.New())
	app.Use(middleware.LoggerMiddleware(a.logger))
	app.Use(middleware.LocalizerMiddleware(a.languageBundle))
	app.Use(middleware.NewRelicMiddleware(a.newRelicInstance))

	// Validator
	validator := validation.InitValidator()
	app.Use(func(c *fiber.Ctx) error {
		c.Locals(utils.ValidatorKey, validator)

		return c.Next()
	})

	// Tokenizer
	app.Use(func(c *fiber.Ctx) error {
		c.Locals(utils.TokenizerKey, a.tokenizerInstance)

		return c.Next()
	})
}

func (a *application) addHealthCheckRoutes(app *fiber.App) {
	healthCheckHandler := dc.InitHealthCheckHandler()
	app.Get("/liveness", healthCheckHandler.Liveness)
}
