package app

import "errors"

var (
	ErrBotSecretToken = errors.New("bot: failed to generate webhook secret token")
	ErrBotTokenNotSet = errors.New("bot: token is not set")
	ErrBotCreate      = errors.New("bot: failed to create bot")
)
