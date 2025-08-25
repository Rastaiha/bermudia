package config

import "time"

type Config struct {
	TokenSigningKey       []byte        `config:"token_signing_key"`
	MockUsersPassword     string        `config:"mock_users_password"`
	BotToken              string        `config:"bot_token"`
	MinCorrectionDelay    time.Duration `config:"min_correction_delay"`
	CorrectionJobInterval time.Duration `config:"min_correction_delay"`
}

func defaultConfig() *Config {
	return &Config{
		MinCorrectionDelay:    30 * time.Second,
		CorrectionJobInterval: 10 * time.Second,
	}
}
