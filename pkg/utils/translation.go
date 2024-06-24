package utils

import (
	"context"

	"github.com/gofiber/fiber/v2"
	"github.com/nicksnyder/go-i18n/v2/i18n"

	constants "github.com/BiTaksi/drivercampaign/pkg/constans"
)

func TranslateByIDWithContext(ctx context.Context, msgID string) string {
	if l, ok := ctx.Value(LocalizerKey).(*i18n.Localizer); ok {
		msg, _ := l.LocalizeMessage(&i18n.Message{
			ID: msgID,
		})
		return msg
	}
	return ""
}

func GetLanguageWithContext(_ *fiber.Ctx) string {
	return constants.EnglishLanguage
}
