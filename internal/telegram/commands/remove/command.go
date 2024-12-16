package remove

import (
	"github.com/RacoonMediaServer/rms-bot-client/pkg/command"
	"github.com/RacoonMediaServer/rms-packages/pkg/communication"
	"github.com/racoon-devel/raccoon-pirate/internal/frontend"
	"go-micro.dev/v4/logger"
)

var Command command.Type = command.Type{
	ID:       "remove",
	Title:    "Удалить",
	Help:     "Удаление фильмов и сериалов",
	Factory:  New,
	Internal: true,
}

type removeCommand struct {
	s *frontend.Setup
	l logger.Logger
}

func (r *removeCommand) Do(ctx command.Context) (bool, []*communication.BotMessage) {
	if len(ctx.Arguments) != 1 {
		return true, command.ReplyText(command.ParseArgumentsFailed)
	}

	if err := r.s.TorrentService.Remove(ctx.Arguments[0]); err != nil {
		r.l.Logf(logger.ErrorLevel, "Remove torrent failed: %s", err)
		return true, command.ReplyText(command.SomethingWentWrong)
	}

	return true, command.ReplyText(command.Removed)
}

func New(interlayer command.Interlayer, l logger.Logger) command.Command {
	s, _ := command.InterlayerLoad[*frontend.Setup](&interlayer)
	return &removeCommand{
		s: s,
		l: l.Fields(map[string]interface{}{"command": "remove"}),
	}
}
