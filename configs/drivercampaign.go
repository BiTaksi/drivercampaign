package configs

var DriverCampaignApp *DriverCampaignScheme

type DriverCampaignScheme struct {
	Web                       WebConfig                 `mapstructure:",squash"`
	NewRelic                  NewRelicConfig            `mapstructure:",squash"`
	MicroserviceJWT           MicroserviceJWTConfig     `mapstructure:",squash"`
	BasicAuthenticationConfig BasicAuthenticationConfig `mapstructure:",squash"`
}
