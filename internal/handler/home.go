package handler

import (
	"time"

	"github.com/gofiber/fiber/v2"

	"github.com/BiTaksi/drivercampaign/configs"
	"github.com/BiTaksi/drivercampaign/internal/dto/resource"
	"github.com/BiTaksi/drivercampaign/pkg/response"
)

type IHomeHandler interface {
	Home(c *fiber.Ctx) error
}

type homeHandler struct{}

func NewHomeHandler() IHomeHandler {
	return &homeHandler{}
}

func (h *homeHandler) Home(c *fiber.Ctx) error {
	return c.JSON(response.NewSuccessResponse(&resource.HomeResource{
		App:  configs.DriverCampaignApp.Web.AppName,
		Env:  configs.DriverCampaignApp.Web.Env,
		Time: time.Now(),
	}))
}
