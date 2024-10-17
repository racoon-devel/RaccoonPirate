package web

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/racoon-devel/raccoon-pirate/internal/model"
)

type uiPage struct {
	Redirect string
}

type errorPage struct {
	uiPage
	Error string
}

type okPage struct {
	uiPage
	Text string
}

func displayError(ctx *gin.Context, status int, err string) {
	ctx.HTML(status, "error.tmpl", &errorPage{
		Error: err,
	})
}

func displayOK(ctx *gin.Context, text, redirect string) {
	page := okPage{
		Text: text,
	}
	page.Redirect = redirect
	ctx.HTML(http.StatusOK, "success.tmpl", &page)
}

func iotaSeasons(count uint) []uint {
	result := make([]uint, count)
	for i := uint(1); i <= count; i++ {
		result[i-1] = i
	}
	return result
}

func decodeMediaType(t string) model.MediaType {
	switch t {
	case "movies":
		return model.MediaTypeMovie
	case "music":
		return model.MediaTypeArtist
	case "others":
		return model.MediaTypeOther
	default:
		return model.MediaTypeMovie
	}
}
