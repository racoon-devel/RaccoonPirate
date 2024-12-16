package library

import (
	"fmt"

	"github.com/RacoonMediaServer/rms-packages/pkg/communication"
	"github.com/racoon-devel/raccoon-pirate/internal/model"
)

func formatTorrent(t *model.Torrent) *communication.BotMessage {
	msg := communication.BotMessage{}

	msg.Text = "<b>"
	if t.BelongsTo != "" {
		msg.Text += t.BelongsTo + " "
	}

	msg.Text += fmt.Sprintf("{ %s }", t.Title)

	if t.Year != 0 {
		msg.Text += fmt.Sprintf(" (%d)", t.Year)
	}
	msg.Text += "</b>"

	msg.KeyboardStyle = communication.KeyboardStyle_Message
	msg.Buttons = append(msg.Buttons, &communication.Button{
		Title:   "Удалить",
		Command: "/remove " + t.ID,
	})

	return &msg
}
