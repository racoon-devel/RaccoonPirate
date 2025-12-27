package frontend

type TelegramAccessData struct {
	BotUrl string
	IdCode string
}

type TelegramAccessProvider interface {
	GetTelegramAccessData() TelegramAccessData
}
