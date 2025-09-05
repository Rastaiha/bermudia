package config

import "time"

type Config struct {
	DevMode                bool          `config:"dev_mode"`
	TokenSigningKey        []byte        `config:"token_signing_key"`
	MockUsersPassword      string        `config:"mock_users_password"`
	BotToken               string        `config:"bot_token"`
	MinCorrectionDelay     time.Duration `config:"min_correction_delay"`
	CorrectionJobInterval  time.Duration `config:"correction_job_interval"`
	DefaultCorrectionGroup int64         `config:"default_correction_group"`
	CorrectionGroupsStr    string        `config:"correction_groups"`
	CorrectionGroups       map[string]int64
}

func defaultConfig() *Config {
	return &Config{
		MinCorrectionDelay:    30 * time.Second,
		CorrectionJobInterval: 10 * time.Second,
	}
}
