package main

import (
	"fmt"
	"time"

	"github.com/BiTaksi/drivercampaign/pkg/localizer"

	"github.com/nicksnyder/go-i18n/v2/i18n"
	"github.com/sirupsen/logrus"
	"golang.org/x/text/language"

	"github.com/BiTaksi/drivercampaign/configs"
	"github.com/BiTaksi/drivercampaign/pkg/logging"
	"github.com/BiTaksi/drivercampaign/pkg/nrclient"
	"github.com/BiTaksi/drivercampaign/pkg/requestclient"
	"github.com/BiTaksi/drivercampaign/pkg/tokenizer"
)

func boot(logger *logrus.Logger) (*application, error) {
	newRelicInstance, newRelicErr := initNewRelic()
	if newRelicErr != nil {
		return nil, fmt.Errorf("newrelic: %v", newRelicErr)
	}

	languageBundle := initLanguageBundler()
	tokenizerInstance := initTokenizer()
	requestClient := requestclient.NewClientRequest()

	return &application{
		logger:                logger,
		tokenizerInstance:     tokenizerInstance,
		newRelicInstance:      newRelicInstance,
		requestClientInstance: requestClient,
		languageBundle:        languageBundle,
	}, nil
}

func initConfig() error {
	var config configs.DriverCampaignScheme
	if err := configs.InitConfig(&config); err != nil {
		return err
	}

	configs.DriverCampaignApp = &config

	return nil
}

func initNewRelic() (nrclient.INewRelicInstance, error) {
	return nrclient.InitNewRelic(nrclient.Config{
		Key:     configs.DriverCampaignApp.NewRelic.Key,
		AppName: configs.DriverCampaignApp.NewRelic.AppName,
	})
}

func initTokenizer() tokenizer.ITokenizer {
	return tokenizer.NewTokenizer(tokenizer.JWTConfig{
		Secret: configs.DriverCampaignApp.MicroserviceJWT.Secret,
		Expr:   time.Second * time.Duration(configs.DriverCampaignApp.MicroserviceJWT.Expr),
	}, tokenizer.BasicAuthenticationConfig{
		Username: configs.DriverCampaignApp.BasicAuthenticationConfig.Username,
		Password: configs.DriverCampaignApp.BasicAuthenticationConfig.Password,
	})
}

func initLanguageBundler() *i18n.Bundle {
	return localizer.InitLocalizer(localizer.LanguageConfig{
		Default: language.English,
		Languages: []language.Tag{
			language.English,
		},
	})
}

func initLogger() *logrus.Logger {
	return logging.NewLogger(logging.Config{
		Service: logging.ServiceConfig{
			Env:     configs.DriverCampaignApp.Web.Env,
			AppName: configs.DriverCampaignApp.Web.AppName,
		},
	})
}
