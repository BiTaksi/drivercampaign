package configs

import (
	"os"

	"github.com/BiTaksi/drivercampaign/pkg/vault"
	"github.com/BiTaksi/drivercampaign/pkg/viperconfig"
)

func InitConfig(output interface{}) error {
	if id := os.Getenv("VAULT_APP_ROLE_ID"); id != "" {
		return vault.InitVault(vault.GetConfigViaEnv(), &output)
	}

	path := "."
	if envConfigPath := os.Getenv("CONFIG_FILE_PATH"); envConfigPath != "" {
		path = envConfigPath
	}

	file := ".env"
	if envConfigFile := os.Getenv("CONFIG_FILE_NAME"); envConfigFile != "" {
		file = envConfigFile
	}

	vc := viperconfig.Config{
		Path:     path,
		FileName: file,
		Env:      os.Getenv("ENV"),
	}

	return viperconfig.Load(vc, output)
}
