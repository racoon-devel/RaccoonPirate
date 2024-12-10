package telegram

import "github.com/RacoonMediaServer/rms-bot-client/pkg/bot"

func New(transport bot.Transport) *bot.Bot {
	settings := bot.Settings{
		Transport: transport,
	}
	return bot.New(settings)
}
