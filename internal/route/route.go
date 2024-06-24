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
	homeHandler handler.IHomeHandler
}

func NewRoute(
	hHandler handler.IHomeHandler,
) IRoute {
	return &route{
		homeHandler: hHandler,
	}
}

func (r *route) SetupRoutes(ac *AppContext) {
	// v1 routes
	v1Group := ac.App.Group("/v1")
	v1Group.Get("/", r.homeHandler.Home, middleware.JWTAuthMiddleware())
}
