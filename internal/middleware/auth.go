package middleware

import (
	"github.com/gofiber/fiber/v2"

	"github.com/BiTaksi/drivercampaign/pkg/response"
	"github.com/BiTaksi/drivercampaign/pkg/tokenizer"
	"github.com/BiTaksi/drivercampaign/pkg/utils"
)

func JWTAuthMiddleware() func(c *fiber.Ctx) error {
	return func(c *fiber.Ctx) error {
		t := c.Context().Value(utils.TokenizerKey).(tokenizer.ITokenizer)
		token := c.Get(utils.MicroserviceHeaderKey)
		if token == "" {
			return c.Status(fiber.StatusForbidden).JSON(response.NewAuthorizationError())
		}

		decoded, err := t.VerifyJWTToken(token)
		if err != nil || !decoded.Valid {
			return c.Status(fiber.StatusForbidden).JSON(response.NewAuthorizationError())
		}

		return c.Next()
	}
}

func BasicAuthMiddleware() func(c *fiber.Ctx) error {
	return func(c *fiber.Ctx) error {
		t := c.Context().Value(utils.TokenizerKey).(tokenizer.ITokenizer)
		token := c.Get(utils.XAPIHeaderKey)
		if token == "" {
			return c.Status(fiber.StatusUnauthorized).JSON(response.NewAuthorizationError())
		}

		if !t.IsBasicAuthorized(token) {
			return c.Status(fiber.StatusUnauthorized).JSON(response.NewAuthorizationError())
		}

		return c.Next()
	}
}
