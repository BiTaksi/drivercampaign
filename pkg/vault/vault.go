package vault

import (
	"context"
	"fmt"
	"os"

	"github.com/hashicorp/vault/api"
	"github.com/hashicorp/vault/api/auth/approle"
	"github.com/mitchellh/mapstructure"
)

type Config struct {
	URL          string
	ENV          string
	AppRoleID    string
	Workspace    string
	SecretIDPath string
}

func InitVault(c Config, output interface{}) error {
	config := api.DefaultConfig()
	config.Address = c.URL

	client, err := api.NewClient(config)
	if err != nil {
		return fmt.Errorf("new client: %v", err)
	}

	ctx := context.Background()
	auth, authErr := approle.NewAppRoleAuth(c.AppRoleID, &approle.SecretID{FromFile: c.SecretIDPath})
	if authErr != nil {
		return fmt.Errorf("app role: %v", authErr)
	}

	login, loginErr := client.Auth().Login(ctx, auth)
	if loginErr != nil {
		return fmt.Errorf("auth: %v", loginErr)
	}
	if login == nil {
		return fmt.Errorf("no auth")
	}

	secret, secretErr := client.KVv2(c.ENV).Get(ctx, c.Workspace)
	if secretErr != nil {
		return fmt.Errorf("secret: %v", secretErr)
	}

	decoder, decoderErr := mapstructure.NewDecoder(&mapstructure.DecoderConfig{
		Metadata:         nil,
		Result:           output,
		WeaklyTypedInput: true,
		DecodeHook: mapstructure.ComposeDecodeHookFunc(
			mapstructure.StringToTimeDurationHookFunc(),
			mapstructure.StringToSliceHookFunc(","),
		),
	})
	if decoderErr != nil {
		return fmt.Errorf("decode: %v", decoderErr)
	}

	return decoder.Decode(secret.Data)
}

func GetConfigViaEnv() Config {
	return Config{
		URL:          os.Getenv("VAULT_URL"),
		ENV:          os.Getenv("VAULT_ENV"),
		AppRoleID:    os.Getenv("VAULT_APP_ROLE_ID"),
		Workspace:    os.Getenv("VAULT_WORKSPACE"),
		SecretIDPath: os.Getenv("VAULT_SECRET_ID_PATH"),
	}
}
