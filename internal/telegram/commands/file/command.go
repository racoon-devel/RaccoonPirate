package file

import (
	"github.com/RacoonMediaServer/rms-bot-client/pkg/command"
	"github.com/RacoonMediaServer/rms-packages/pkg/communication"
	"github.com/racoon-devel/raccoon-pirate/internal/frontend"
	"github.com/racoon-devel/raccoon-pirate/internal/model"
	"go-micro.dev/v4/logger"
)

var Command command.Type = command.Type{
	ID:       "file",
	Title:    "Загрузить",
	Help:     "Загрузка файлов на сервер",
	Factory:  New,
	Internal: true,
}

type fileCommand struct {
	s          *frontend.Setup
	l          logger.Logger
	content    []byte
	chooseType bool
}

func (c *fileCommand) Do(ctx command.Context) (bool, []*communication.BotMessage) {
	if c.chooseType {
		contentType, ok := frontend.DetermineContentType(ctx.Arguments.String())
		if !ok {
			return false, command.ReplyText("Неизвестный тип медиа")
		}

		record := model.Torrent{Type: contentType, Content: c.content}
		if err := c.s.TorrentService.Add(ctx, &record); err != nil {
			c.l.Logf(logger.ErrorLevel, "Add torrent failed: %s", err)
			return true, command.ReplyText(command.SomethingWentWrong)
		}

		return true, command.ReplyText("Добавлено")
	}

	if ctx.Attachment == nil {
		return true, command.ReplyText("Необходимо прислать торрент-файл")
	}
	if ctx.Attachment.MimeType != "application/x-bittorrent" {
		return true, command.ReplyText("Неверный формат файла. Поддерживаются только торрент-файлы")
	}

	c.content = ctx.Attachment.Content
	c.chooseType = true

	msg := communication.BotMessage{
		Text:          "Выберите тип загружаемого медиа",
		KeyboardStyle: communication.KeyboardStyle_Chat,
		Buttons:       frontend.GetContentTypesButtonsRu(),
	}

	return false, []*communication.BotMessage{&msg}
}

func New(interlayer command.Interlayer, l logger.Logger) command.Command {
	s, _ := command.InterlayerLoad[*frontend.Setup](&interlayer)
	return &fileCommand{
		s: s,
		l: l.Fields(map[string]interface{}{"command": "file"}),
	}
}
