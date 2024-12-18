package config

type MessagingConfig struct {
	XSMS          XSMSConfig    `mapstructure:"xsms"`
	Every8DConfig Every8DConfig `mapstructure:"every8d"`
	Mailgun       MailgunConfig
}

type XSMSConfig struct {
	MDN      string `mapstructure:"mdn"`
	Username string
	Password string
	BaseURL  string `mapstructure:"base_url"`
}

type MailgunConfig struct {
	BaseURL string `mapstructure:"base_url"`
	APIKey  string `mapstructure:"api_key"`
}

type Every8DConfig struct {
	Username string `mapstructure:"username"`
	Password string `mapstructure:"password"`
	BaseURL  string `mapstructure:"base_url"`
}

type PushNotifConfig struct {
	Onesignal map[string]OnesignalConfig `mapstructure:"onesignal"`
}

type OnesignalConfig struct {
	BaseURL          string `mapstructure:"base_url"`
	AppID            string `mapstructure:"app_id"`
	APIKey           string `mapstructure:"api_key"`
	AndroidChannelID string `mapstructure:"android_channel_id"`
}
