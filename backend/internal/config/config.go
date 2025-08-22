package config

type Config struct {
	TokenSigningKey   []byte `config:"token_signing_key"`
	MockUsersPassword string `config:"mock_users_password"`
	BotToken          string `config:"bot_token"`
}

func defaultConfig() *Config {
	return &Config{}
}
