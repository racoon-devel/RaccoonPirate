package web

import (
	"net/http"

	"github.com/gin-gonic/gin"
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
