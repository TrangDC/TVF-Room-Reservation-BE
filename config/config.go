package config

type Configurations struct {
	App          AppConfig          `mapstructure:"app"`
	Postgres     PostgresConfig     `mapstructure:"postgres"`
	AzureADOAuth AzureADOAuthConfig `mapstructure:"azureAD"`
}

type AppConfig struct {
	Name        string `mapstructure:"name"`
	Environment string `mapstructure:"environment"`
	Version     string `mapstructure:"version"`
	ServiceName string `mapstructure:"serviceName"`
	SeverName   string `mapstructure:"serverName"`
	Debug       bool   `mapstructure:"debug"`
}

type PostgresConfig struct {
	ConnectionString string `mapstructure:"connectionString"`
}

type AzureADOAuthConfig struct {
	Enabled           bool   `mapstructure:"enabled"`
	ClientID          string `mapstructure:"clientID"`
	ClientSecret      string `mapstructure:"clientSecret"`
	TenantID          string `mapstructure:"tenantID"`
	RedirectUrl       string `mapstructure:"redirectUrl"`
	Scopes            string `mapstructure:"scopes"`
	ClientRedirectUrl string `mapstructure:"clientRedirectUrl"`
}
