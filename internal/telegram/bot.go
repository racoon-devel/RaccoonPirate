package telegram

import (
	"context"
	"fmt"
	"strings"

	"github.com/RacoonMediaServer/rms-bot-client/pkg/bot"
	"github.com/RacoonMediaServer/rms-bot-client/pkg/command"
	"github.com/RacoonMediaServer/rms-bot-client/pkg/commands"
	"github.com/RacoonMediaServer/rms-bot-client/pkg/commands/unlink"
	rms_bot_client "github.com/RacoonMediaServer/rms-packages/pkg/service/rms-bot-client"
	"github.com/apex/log"
	"github.com/racoon-devel/raccoon-pirate/internal/frontend"
	"github.com/racoon-devel/raccoon-pirate/internal/telegram/commands/search"
	"google.golang.org/protobuf/types/known/emptypb"
)

const botNickName = "RaccoonPirateBot"

type Bot struct {
	frontend.Setup
	Transport bot.Transport

	b *bot.Bot

	ctx    context.Context
	cancel context.CancelFunc
}

func (b *Bot) Run() {
	pirateCommands := commands.MakeRegisteredCommands(search.Command, unlink.Command)

	settings := bot.Settings{
		Transport:  b.Transport,
		Interlayer: command.Interlayer{},
		CmdFactory: commands.NewDefaultFactory(pirateCommands),
	}

	command.InterlayerStore(&settings.Interlayer, &b.Setup)

	b.b = bot.New(settings)

	b.ctx, b.cancel = context.WithCancel(context.Background())

	go b.obtainIdentificationCode()
}

func (b *Bot) obtainIdentificationCode() {
	resp := rms_bot_client.GetIdentificationCodeResponse{}
	if err := b.b.GetIdentificationCode(b.ctx, &emptypb.Empty{}, &resp); err != nil {
		log.Errorf("Obtain Telegram identification code failed: %s", err)
		return
	}

	s := strings.Repeat("*", 30) + "\n"
	s += fmt.Sprintf("* %-27s*\n", "Bot: @"+botNickName)
	s += fmt.Sprintf("* %-27s*\n", "Code: "+resp.Code)
	s += strings.Repeat("*", 30) + "\n"
	log.Infof("Telegram connection info: \n%s", s)
}

func (b *Bot) Shutdown() {
	b.cancel()
	b.b.Shutdown()
}
