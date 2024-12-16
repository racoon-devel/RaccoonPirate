package search

import (
	"github.com/RacoonMediaServer/rms-bot-client/pkg/command"
	"github.com/RacoonMediaServer/rms-media-discovery/pkg/media"
	"github.com/RacoonMediaServer/rms-packages/pkg/communication"
	"github.com/racoon-devel/raccoon-pirate/internal/cache"
	"github.com/racoon-devel/raccoon-pirate/internal/frontend"
	"go-micro.dev/v4/logger"
)

var Command command.Type = command.Type{
	ID:      "search",
	Title:   "Поиск контента",
	Help:    "Позволяет искать информацию о фильмах/сериалах/музыке и перейти к скачиванию",
	Factory: New,
}

var contentTypeHelper = map[string]media.ContentType{
	"Фильмы": media.Movies,
	"Музыка": media.Music,
	"Другое": media.Other,
}

type searchCommand struct {
	s     *frontend.Setup
	c     *cache.Cache
	l     logger.Logger
	query string
}

func (s *searchCommand) Do(ctx command.Context) (bool, []*communication.BotMessage) {
	if len(ctx.Arguments) < 1 {
		return false, command.ReplyText("Что ищем?")
	}

	if s.query == "" {
		s.query = ctx.Arguments.String()
		msg := communication.BotMessage{
			Text:          "Тип контента?",
			KeyboardStyle: communication.KeyboardStyle_Chat,
		}
		for k := range contentTypeHelper {
			msg.Buttons = append(msg.Buttons, &communication.Button{Title: k, Command: k})
		}
		return false, []*communication.BotMessage{&msg}
	}

	contentType, ok := contentTypeHelper[ctx.Arguments.String()]
	if !ok {
		return false, command.ReplyText("Не удалось распознать тип контента")
	}

	switch contentType {
	case media.Movies:
		return true, s.searchMovies(ctx)
	case media.Music:
		return true, s.searchMusic(ctx)
	case media.Other:
		fallthrough
	default:
		return true, s.searchOther(ctx)
	}
}

func (s *searchCommand) searchMovies(ctx command.Context) []*communication.BotMessage {
	movies, err := s.s.DiscoveryService.SearchMovies(ctx, s.query)
	if err != nil {
		s.l.Logf(logger.ErrorLevel, "Search movies failed: %s", err)
		return command.ReplyText(command.SomethingWentWrong)
	}

	if len(movies) == 0 {
		return command.ReplyText(command.NothingFound)
	}

	// выводим в обратном порядке,чтобы не мотать ленту в тг
	result := make([]*communication.BotMessage, len(movies))
	for i, mov := range movies {
		s.c.Store(mov.ID, mov)
		result[len(result)-i-1] = s.formatMovieMessage(mov)
	}

	return result
}
func (s *searchCommand) searchMusic(ctx command.Context) []*communication.BotMessage {
	music, err := s.s.DiscoveryService.SearchMusic(ctx, s.query)
	if err != nil {
		s.l.Logf(logger.ErrorLevel, "Search music failed: %s", err)
		return command.ReplyText(command.SomethingWentWrong)
	}

	if len(music) == 0 {
		return command.ReplyText(command.NothingFound)
	}

	result := make([]*communication.BotMessage, len(music))
	for i, mu := range music {
		s.c.Store(mu.Title(), mu)
		result[len(result)-i-1] = s.formatMusicMessage(mu)
	}

	return result
}

func (s *searchCommand) searchOther(ctx command.Context) []*communication.BotMessage {
	// id := uuid.NewV4().String()
	// s.c.Store(id, s.query)
	// s.
	return nil
}

func New(interlayer command.Interlayer, l logger.Logger) command.Command {
	s, _ := command.InterlayerLoad[*frontend.Setup](&interlayer)
	c, _ := command.InterlayerLoad[*cache.Cache](&interlayer)
	return &searchCommand{
		s: s,
		c: c,
		l: l.Fields(map[string]interface{}{"command": "search"}),
	}
}
