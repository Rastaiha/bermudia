package config

import "time"

type Config struct {
	DevMode                bool          `config:"dev_mode"`
	Postgres               Postgres      `config:"postgres"`
	TokenSigningKey        []byte        `config:"token_signing_key"`
	MockUsersPassword      string        `config:"mock_users_password"`
	BotToken               string        `config:"bot_token"`
	MinCorrectionDelay     time.Duration `config:"min_correction_delay"`
	CorrectionJobInterval  time.Duration `config:"correction_job_interval"`
	DefaultCorrectionGroup int64         `config:"default_correction_group"`
	CorrectionGroupsStr    string        `config:"correction_groups"`
	ContentFileID          string        `config:"content_file_id"`
	CorrectionGroups       map[string]int64
}

type Postgres struct {
	Enable  bool   `config:"enable"`
	Host    string `config:"host"`
	Port    int    `config:"port"`
	User    string `config:"user"`
	Pass    string `config:"pass"`
	DB      string `config:"db"`
	SSLMode string `config:"ssl_mode"`
}

func defaultConfig() *Config {
	return &Config{
		MinCorrectionDelay: 30 * time.Second,
		Postgres: Postgres{
			SSLMode: "disable",
		},
		CorrectionJobInterval: 10 * time.Second,
	}
}
