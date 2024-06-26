package handler

import (
	"github.com/gofiber/fiber/v2"

	"github.com/BiTaksi/drivercampaign/pkg/healthcheck"
)

const (
	healthLivenessStatusOk       = "UP"
	healthLivenessStatusShutdown = "SHUTDOWN"
)

type IHealthCheckHandler interface {
	Liveness(c *fiber.Ctx) error
}

type healthCheckHandler struct{}

func NewHealthCheckHandler() IHealthCheckHandler {
	return &healthCheckHandler{}
}

func (h *healthCheckHandler) Liveness(c *fiber.Ctx) error {
	if !healthcheck.Liveness() {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"status": healthLivenessStatusShutdown})
	}

	return c.JSON(fiber.Map{"status": healthLivenessStatusOk})
}
