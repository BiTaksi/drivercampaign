package localizer

import (
	"fmt"

	"github.com/BurntSushi/toml"
	"github.com/nicksnyder/go-i18n/v2/i18n"
	"golang.org/x/text/language"
)

type LanguageConfig struct {
	Default   language.Tag
	Languages []language.Tag
}

func InitLocalizer(l LanguageConfig) *i18n.Bundle {
	bundle := i18n.NewBundle(l.Default)
	bundle.RegisterUnmarshalFunc("toml", toml.Unmarshal)

	for _, language := range l.Languages {
		bundle.MustLoadMessageFile(fmt.Sprintf("locale/active.%s.toml", language.String()))
	}

	return bundle
}
