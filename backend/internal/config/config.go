package config

import "time"

type Config struct {
	DevMode                bool          `config:"dev_mode"`
	Postgres               Postgres      `config:"postgres"`
	TokenSigningKey        []byte        `config:"token_signing_key"`
	MockUsersPassword      string        `config:"mock_users_password"`
	BotToken               string        `config:"bot_token"`
	CorrectionJobInterval  time.Duration `config:"correction_job_interval"`
	DefaultCorrectionGroup int64         `config:"default_correction_group"`
	CorrectionGroupsStr    string        `config:"correction_groups"`
	ContentFileID          string        `config:"content_file_id"`
	CorrectionGroups       map[string]int64
	CorrectionRevertWindow time.Duration `config:"correction_revert_window"`
	CreateMock             bool          `config:"create_mock"`
	AdminsGroup            int64         `config:"admins_group"`
}

func (c Config) MinCorrectionDelay() time.Duration {
	if c.DevMode {
		return 10 * time.Second
	}
	return 3 * time.Minute
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
		Postgres: Postgres{
			SSLMode: "disable",
		},
		CorrectionJobInterval: 10 * time.Second,
	}
}
