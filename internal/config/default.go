package config

import (
	"time"
)

// Default is default configuration
var Default = Config{
	Bot: Bot{
		LongPollTimeout: 10 * time.Second,
		UseWebhook:      false,
		WebhookListen:   ":8080",
	},
	DB: DB{
		RequiredVersion: 1,
	},
}
