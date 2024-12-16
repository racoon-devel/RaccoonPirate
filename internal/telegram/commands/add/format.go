package add

import (
	"fmt"
	"strconv"

	"github.com/RacoonMediaServer/rms-media-discovery/pkg/client/models"
	"github.com/RacoonMediaServer/rms-media-discovery/pkg/model"
	"github.com/RacoonMediaServer/rms-packages/pkg/communication"
)

func formatSelectSeason(mov *model.Movie) *communication.BotMessage {
	msg := communication.BotMessage{Text: "Выберите сезон"}
	msg.KeyboardStyle = communication.KeyboardStyle_Chat
	msg.Buttons = append(msg.Buttons, &communication.Button{
		Title:   "Все",
		Command: "Все",
	})

	for i := uint(1); i <= mov.Seasons; i++ {
		no := strconv.FormatUint(uint64(i), 10)
		msg.Buttons = append(msg.Buttons, &communication.Button{Title: no, Command: no})
	}

	return &msg
}

func formatTorrents(torrents []*models.SearchTorrentsResult) *communication.BotMessage {
	msg := communication.BotMessage{}
	msg.KeyboardStyle = communication.KeyboardStyle_Chat
	for i, t := range torrents {
		msg.Text += fmt.Sprintf("%d. %s [ %.2f Gb, %d seeds]\n", i+1, *t.Title, float32(*t.Size)/1024., *t.Seeders)
		no := fmt.Sprintf("%d", i+1)
		msg.Buttons = append(msg.Buttons, &communication.Button{Title: no, Command: no})
	}
	return &msg
}
