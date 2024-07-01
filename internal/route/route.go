package route

import (
	"github.com/gofiber/fiber/v2"

	"github.com/BiTaksi/drivercampaign/internal/handler"
	"github.com/BiTaksi/drivercampaign/internal/middleware"
)

type AppContext struct {
	App *fiber.App
}

type IRoute interface {
	SetupRoutes(ac *AppContext)
}

type route struct {
	homeHandler     handler.IHomeHandler
	campaignHandler handler.ICampaignHandler
}

func NewRoute(
	hHandler handler.IHomeHandler,
	cHandler handler.ICampaignHandler,
) IRoute {
	return &route{
		homeHandler:     hHandler,
		campaignHandler: cHandler,
	}
}

func (r *route) SetupRoutes(ac *AppContext) {
	// v1 routes
	v1Group := ac.App.Group("/v1")
	v1Group.Get("/", r.homeHandler.Home, middleware.JWTAuthMiddleware())

	r.campaignsRoutes(v1Group)
}

func (r *route) campaignsRoutes(fr fiber.Router) {
	fr.Get("/campaigns-home", r.campaignHandler.GetCampaignsHome)
}
