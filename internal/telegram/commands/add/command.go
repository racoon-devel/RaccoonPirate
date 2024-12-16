package add

import (
	"context"
	"errors"
	"strconv"

	"github.com/RacoonMediaServer/rms-bot-client/pkg/command"
	"github.com/RacoonMediaServer/rms-media-discovery/pkg/client/models"
	"github.com/RacoonMediaServer/rms-media-discovery/pkg/media"
	"github.com/RacoonMediaServer/rms-media-discovery/pkg/model"
	"github.com/RacoonMediaServer/rms-packages/pkg/communication"
	"github.com/racoon-devel/raccoon-pirate/internal/cache"
	"github.com/racoon-devel/raccoon-pirate/internal/frontend"
	internalModel "github.com/racoon-devel/raccoon-pirate/internal/model"
	"github.com/racoon-devel/raccoon-pirate/internal/selector"
	"go-micro.dev/v4/logger"
)

var Command command.Type = command.Type{
	ID:       "add",
	Title:    "Добавить",
	Help:     "",
	Internal: true,
	Factory:  New,
}

type state int

const (
	stateInitial state = iota
	stateChooseSeason
	stateChooseTorrent
	stateWaitFile
)

type addCommand struct {
	s *frontend.Setup
	c *cache.Cache
	l logger.Logger

	state    state
	stateMap map[state]command.Handler
	download command.Handler

	id       string
	season   *int64
	torrents []string
	all      bool

	torrentRecord internalModel.Torrent
	mov           *model.Movie
	mus           model.Music
	query         string
}

func (d *addCommand) Do(ctx command.Context) (done bool, messages []*communication.BotMessage) {
	return d.stateMap[d.state](ctx)
}

func (d *addCommand) doInitial(ctx command.Context) (bool, []*communication.BotMessage) {
	if len(ctx.Arguments) < 2 {
		return true, command.ReplyText(command.ParseArgumentsFailed)
	}
	switch ctx.Arguments[0] {
	case "auto":
		d.download = d.addAuto

	case "select":
		d.download = d.addSelect

	case "file":
		d.download = d.addFile

	default:
		return true, command.ReplyText(command.ParseArgumentsFailed)
	}

	idArgs := ctx.Arguments[1:]
	d.id = idArgs.String()

	item, ok := d.c.Load(d.id)
	if !ok {
		d.l.Logf(logger.ErrorLevel, "Fetch item %s from cache failed", d.id)
		return true, command.ReplyText(command.SomethingWentWrong)
	}

	switch item := item.(type) {
	case *model.Movie:
		d.mov = item
		d.torrentRecord.ExpandByMovie(item)

		if item.Type == model.MovieType_TvSeries && item.Seasons != 0 {
			d.state = stateChooseSeason
			return false, []*communication.BotMessage{formatSelectSeason(item)}
		}
	case model.Music:
		d.mus = item
		d.query = item.Title()
		d.torrentRecord.ExpandByMusic(item)
	case string:
		d.torrentRecord.Type = media.Other
		d.query = item
	default:
		d.l.Logf(logger.ErrorLevel, "Unknown type of media: %T", item)
		return true, command.ReplyText(command.SomethingWentWrong)
	}

	return d.download(ctx)
}

func (d *addCommand) addAuto(ctx command.Context) (bool, []*communication.BotMessage) {
	torrents, err := d.searchTorrents(ctx)
	if err != nil {
		d.l.Logf(logger.ErrorLevel, "Search torrents failed: %s", err)
		return true, command.ReplyText(command.SomethingWentWrong)
	}

	if len(torrents) == 0 {
		return true, command.ReplyText("Не удалось найти подходящую раздачу")
	}

	opts := selector.Options{
		MediaType:   d.torrentRecord.Type,
		Criteria:    d.s.SelectCriterion,
		Discography: d.torrentRecord.Type == media.Music && d.mus.IsArtist(),
	}

	if d.all {
		opts.Criteria = selector.CriteriaCompact
	}

	if d.torrentRecord.Type != media.Movies {
		opts.Query = d.query
	}

	picked := d.s.Selector.Select(torrents, opts)
	return d.addTorrent(ctx, *picked.Link)
}

func (d *addCommand) addSelect(ctx command.Context) (bool, []*communication.BotMessage) {
	torrents, err := d.searchTorrents(ctx)
	if err != nil {
		d.l.Logf(logger.ErrorLevel, "Search torrents failed: %s", err)
		return true, command.ReplyText(command.SomethingWentWrong)
	}

	if len(torrents) == 0 {
		return true, command.ReplyText("Не удалось найти подходящую раздачу")
	}

	d.torrents = make([]string, len(torrents))
	for i := range torrents {
		d.torrents[i] = *torrents[i].Link
	}

	d.state = stateChooseTorrent
	return false, []*communication.BotMessage{formatTorrents(torrents)}
}

func (d *addCommand) addFile(ctx command.Context) (bool, []*communication.BotMessage) {
	d.state = stateWaitFile
	return false, command.ReplyText("Необходимо прислать торрент-файл с содержимым выбранного фильма/сериала")
}

func (d *addCommand) doChooseSeason(ctx command.Context) (bool, []*communication.BotMessage) {
	if len(ctx.Arguments) != 1 {
		return false, command.ReplyText("Необходимо выбрать сезон")
	}
	if ctx.Arguments[0] == "Все" {
		d.all = true
		return d.download(ctx)
	}
	season, err := strconv.ParseUint(ctx.Arguments[0], 10, 8)
	if err != nil {
		return false, command.ReplyText("Неверно указан номер сезона")
	}
	s := int64(season)
	d.season = &s

	return d.download(ctx)
}

func (d *addCommand) doChooseTorrent(ctx command.Context) (bool, []*communication.BotMessage) {
	if len(ctx.Arguments) != 1 {
		return false, command.ReplyText("Необходимо выбрать раздачу")
	}
	no, err := strconv.ParseInt(ctx.Arguments[0], 10, 8)
	if err != nil || no <= 0 || no > int64(len(d.torrents)) {
		return false, command.ReplyText("Неверно указан номер раздачи")
	}

	link := d.torrents[no-1]
	return d.addTorrent(ctx, link)
}

func (d *addCommand) doWaitFile(ctx command.Context) (bool, []*communication.BotMessage) {
	if ctx.Attachment == nil {
		return false, command.ReplyText("Необходимо прислать торрент-файл")
	}
	if ctx.Attachment.MimeType != "application/x-bittorrent" {
		return false, command.ReplyText("Неверный формат файла")
	}

	if err := d.s.TorrentService.Add(&d.torrentRecord, ctx.Attachment.Content); err != nil {
		d.l.Logf(logger.ErrorLevel, "Add torrent failed: %s", err)
		return false, command.ReplyText(command.SomethingWentWrong)
	}

	return true, command.ReplyText("Добавлено")
}

func (d *addCommand) searchTorrents(ctx context.Context) (result []*models.SearchTorrentsResult, err error) {
	switch d.torrentRecord.Type {
	case media.Movies:
		result, err = d.s.DiscoveryService.SearchMovieTorrents(ctx, d.mov, d.season)
	case media.Music:
		result, err = d.s.DiscoveryService.SearchMusicTorrents(ctx, d.mus)
	case media.Other:
		result, err = d.s.DiscoveryService.SearchOtherTorrents(ctx, d.query)
	default:
		err = errors.New("unknown content type")
	}
	return
}

func (d *addCommand) addTorrent(ctx context.Context, link string) (bool, []*communication.BotMessage) {
	content, err := d.s.DiscoveryService.GetTorrent(ctx, link)
	if err != nil {
		d.l.Logf(logger.ErrorLevel, "Get torrent failed: %s", err)
		return false, command.ReplyText(command.SomethingWentWrong)
	}

	if err = d.s.TorrentService.Add(&d.torrentRecord, content); err != nil {
		d.l.Logf(logger.ErrorLevel, "Add torrent failed: %s", err)
		return false, command.ReplyText(command.SomethingWentWrong)
	}

	return true, command.ReplyText("Добавлено")
}

func New(interlayer command.Interlayer, l logger.Logger) command.Command {
	s, _ := command.InterlayerLoad[*frontend.Setup](&interlayer)
	c, _ := command.InterlayerLoad[*cache.Cache](&interlayer)
	d := &addCommand{
		s: s,
		c: c,
		l: l.Fields(map[string]interface{}{"command": "add"}),
	}

	d.stateMap = map[state]command.Handler{
		stateInitial:       d.doInitial,
		stateChooseSeason:  d.doChooseSeason,
		stateChooseTorrent: d.doChooseTorrent,
		stateWaitFile:      d.doWaitFile,
	}

	return d
}
