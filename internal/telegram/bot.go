package telegram

import (
	"context"
	"fmt"
	"strings"
	"sync"
	"time"

	"github.com/RacoonMediaServer/rms-bot-client/pkg/bot"
	"github.com/RacoonMediaServer/rms-bot-client/pkg/command"
	"github.com/RacoonMediaServer/rms-bot-client/pkg/commands"
	"github.com/RacoonMediaServer/rms-bot-client/pkg/commands/unlink"
	rms_bot_client "github.com/RacoonMediaServer/rms-packages/pkg/service/rms-bot-client"
	"github.com/apex/log"
	"github.com/racoon-devel/raccoon-pirate/internal/cache"
	"github.com/racoon-devel/raccoon-pirate/internal/frontend"
	"github.com/racoon-devel/raccoon-pirate/internal/telegram/commands/add"
	"github.com/racoon-devel/raccoon-pirate/internal/telegram/commands/file"
	"github.com/racoon-devel/raccoon-pirate/internal/telegram/commands/library"
	"github.com/racoon-devel/raccoon-pirate/internal/telegram/commands/remove"
	"github.com/racoon-devel/raccoon-pirate/internal/telegram/commands/search"
	"google.golang.org/protobuf/types/known/emptypb"
)

const botURL = "https://t.me/RaccoonPirateBot"

const cacheItemTTL = 1 * time.Hour

const refreshCodeInterval = 1 * time.Minute

type Bot struct {
	frontend.Setup
	Transport bot.Transport

	b *bot.Bot
	c *cache.Cache

	ctx    context.Context
	cancel context.CancelFunc

	mu     sync.Mutex
	tgData frontend.TelegramAccessData
}

func (b *Bot) Run() {
	pirateCommands := commands.MakeRegisteredCommands(search.Command, add.Command, file.Command, library.Command, remove.Command, unlink.Command)

	settings := bot.Settings{
		Transport:  b.Transport,
		Interlayer: command.Interlayer{},
		CmdFactory: commands.NewDefaultFactory(pirateCommands),
	}

	b.c = cache.New(cacheItemTTL)
	command.InterlayerStore(&settings.Interlayer, b.c)

	command.InterlayerStore(&settings.Interlayer, &b.Setup)

	b.b = bot.New(settings)

	b.ctx, b.cancel = context.WithCancel(context.Background())

	go b.refreshIdentificationCode()
}

func (b *Bot) refreshIdentificationCode() {
	timer := time.NewTicker(refreshCodeInterval)
	defer timer.Stop()

	b.obtainIdentificationCode()

	for {
		select {
		case <-timer.C:
			b.obtainIdentificationCode()
		case <-b.ctx.Done():
			return
		}
	}
}

func (b *Bot) obtainIdentificationCode() {
	resp := rms_bot_client.GetIdentificationCodeResponse{}
	if err := b.b.GetIdentificationCode(b.ctx, &emptypb.Empty{}, &resp); err != nil {
		log.Errorf("Obtain Telegram identification code failed: %s", err)
		return
	}

	b.mu.Lock()
	changed := b.tgData.IdCode != resp.Code
	b.tgData = frontend.TelegramAccessData{
		BotUrl: botURL,
		IdCode: resp.Code,
	}
	b.mu.Unlock()

	if changed {
		s := strings.Repeat("*", 50) + "\n"
		s += fmt.Sprintf("* %-47s*\n", "Bot: "+botURL)
		s += fmt.Sprintf("* %-47s*\n", "Code: "+resp.Code)
		s += strings.Repeat("*", 50) + "\n"
		log.Infof("Telegram connection info: \n%s", s)
	}
}

func (b *Bot) GetTelegramAccessData() frontend.TelegramAccessData {
	b.mu.Lock()
	defer b.mu.Unlock()

	return b.tgData
}

func (b *Bot) Shutdown() {
	b.cancel()
	b.b.Shutdown()
}
