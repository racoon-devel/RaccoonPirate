package search

import (
	"bytes"
	"embed"
	"fmt"
	"strings"
	"text/template"
	"unicode/utf8"

	"github.com/RacoonMediaServer/rms-media-discovery/pkg/model"
	"github.com/RacoonMediaServer/rms-packages/pkg/communication"
	"go-micro.dev/v4/logger"
)

//go:embed templates
var templates embed.FS

var parsedTemplates *template.Template

func init() {
	parsedTemplates = template.Must(template.ParseFS(templates, "templates/*.txt"))
}

func formatGenres(genres []string) string {
	result := ""
	for _, g := range genres {
		result += strings.TrimSpace(strings.ToLower(g)) + ", "
	}
	if len(result) > 2 {
		result = result[0 : len(result)-2]
	}
	return result
}

func formatDescription(d string) string {
	const maxLength = 350
	if utf8.RuneCountInString(d) <= maxLength {
		return d
	}

	cnt := 0
	found := false
	split := strings.FieldsFunc(d, func(r rune) bool {
		cnt++
		if cnt > maxLength && r == ' ' && !found {
			found = true
			return true
		}
		return false
	})
	return split[0] + "..."
}

func (s *searchCommand) formatMovieMessage(mov *model.Movie) *communication.BotMessage {
	m := &communication.BotMessage{}
	if mov.Poster != "" {
		m.Attachment = &communication.Attachment{
			Type:     communication.Attachment_PhotoURL,
			MimeType: "",
			Content:  []byte(mov.Poster),
		}
	}

	m.Buttons = append(m.Buttons, &communication.Button{Title: "Добавить", Command: "/add " + mov.ID})
	m.Buttons = append(m.Buttons, &communication.Button{Title: "Выбрать раздачу", Command: "/add select " + mov.ID})
	m.Buttons = append(m.Buttons, &communication.Button{Title: "Файл", Command: "/add file " + mov.ID})

	m.KeyboardStyle = communication.KeyboardStyle_Message

	var ui struct {
		Title       string
		Year        uint
		Rating      string
		Genres      string
		Description string
	}
	ui.Title = mov.Title
	ui.Year = mov.Year
	ui.Rating = fmt.Sprintf("%.1f", mov.Rating)
	ui.Genres = formatGenres(mov.Genres)
	ui.Description = formatDescription(mov.Description)

	var buf bytes.Buffer
	if err := parsedTemplates.ExecuteTemplate(&buf, "movie", &ui); err != nil {
		s.l.Logf(logger.ErrorLevel, "execute template failed: %s", err)
	}
	m.Text = buf.String()
	return m
}
