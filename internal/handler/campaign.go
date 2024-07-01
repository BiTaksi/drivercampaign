package handler

import (
	"context"

	"github.com/BiTaksi/drivercampaign/internal/dto/request"
	"github.com/BiTaksi/drivercampaign/internal/service"
	"github.com/BiTaksi/drivercampaign/pkg/response"
	"github.com/BiTaksi/drivercampaign/pkg/utils"
	"github.com/gofiber/fiber/v2"
)

type ICampaignHandler interface {
	GetCampaignsHome(ctx *fiber.Ctx) error
}

type campaignHandler struct {
	campaignService service.ICampaignService
}

func NewCampaignHandler(cs service.ICampaignService) ICampaignHandler {
	return &campaignHandler{campaignService: cs}
}

func (ch *campaignHandler) GetCampaignsHome(c *fiber.Ctx) error {
	req := request.GetCampaignsHome{
		DriverID: c.Get("X-Driver-ID", ""),
	}

	ctx := c.Locals(utils.NewrelicContextKey).(context.Context)

	campaigns, err := ch.campaignService.GetCampaignsHome(ctx, req.ToEntity())
	if err != nil {
		return err
	}

	return c.JSON(response.NewSuccessResponse(campaigns))
}
