package library

import (
	"github.com/RacoonMediaServer/rms-bot-client/pkg/command"
	"github.com/RacoonMediaServer/rms-media-discovery/pkg/media"
	"github.com/RacoonMediaServer/rms-packages/pkg/communication"
	"github.com/racoon-devel/raccoon-pirate/internal/frontend"
	"go-micro.dev/v4/logger"
)

var Command command.Type = command.Type{
	ID:      "library",
	Title:   "Библиотека",
	Help:    "Можно посмотреть, что было скачано и добавлено",
	Factory: New,
}

type libraryCommand struct {
	s *frontend.Setup
	l logger.Logger
}

var contentTypeHelper = map[string]media.ContentType{
	"Фильмы": media.Movies,
	"Музыка": media.Music,
	"Другое": media.Other,
}

func (s *libraryCommand) Do(ctx command.Context) (bool, []*communication.BotMessage) {
	if len(ctx.Arguments) == 0 {
		msg := communication.BotMessage{Text: "Что ищем?"}
		msg.KeyboardStyle = communication.KeyboardStyle_Chat
		for k := range contentTypeHelper {
			msg.Buttons = append(msg.Buttons, &communication.Button{Title: k, Command: k})
		}
		return false, []*communication.BotMessage{&msg}
	}

	contentType, ok := contentTypeHelper[ctx.Arguments[0]]
	if !ok {
		return false, command.ReplyText("Неизвестный тип медиа")
	}

	result, err := s.s.TorrentService.GetTorrentsList(contentType)
	if err != nil {
		s.l.Logf(logger.ErrorLevel, "Fetch torrent list failed: %s", err)
		return true, command.ReplyText(command.SomethingWentWrong)
	}

	if len(result) == 0 {
		return false, command.ReplyText(command.NothingFound)
	}

	messages := make([]*communication.BotMessage, len(result))
	for i, r := range result {
		messages[len(messages)-i-1] = formatTorrent(r)
	}
	return false, messages
}

func New(interlayer command.Interlayer, l logger.Logger) command.Command {
	s, _ := command.InterlayerLoad[*frontend.Setup](&interlayer)
	return &libraryCommand{
		s: s,
		l: l.Fields(map[string]interface{}{"command": "library"}),
	}
}
