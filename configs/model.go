package configs

type WebConfig struct {
	AppName string `mapstructure:"APP_NAME"`
	Port    string `mapstructure:"PORT"`
	Env     string `mapstructure:"ENV"`
}

type MicroserviceJWTConfig struct {
	Secret string `mapstructure:"MICROSERVICE_JWT_SECRET"`
	Expr   int    `mapstructure:"MICROSERVICE_ACCESS_TOKEN_EXPIRATION"`
}

type KafkaConnectionConfig struct {
	Brokers []string `mapstructure:"KAFKA_BROKERS"`
	TLS     bool     `mapstructure:"KAFKA_TLS"`
}

type NewRelicConfig struct {
	AppName string `mapstructure:"NEWRELIC_APP_NAME"`
	Key     string `mapstructure:"NEWRELIC_LICENSE_KEY"`
}

type BasicAuthenticationConfig struct {
	Username string `mapstructure:"BASIC_AUTH_USERNAME"`
	Password string `mapstructure:"BASIC_AUTH_PASSWORD"`
}
