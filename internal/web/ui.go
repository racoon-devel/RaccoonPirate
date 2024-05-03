package web

import "github.com/gin-gonic/gin"

type uiPage struct {
	Redirect string
}

type errorPage struct {
	uiPage
	Error string
}

func displayError(ctx *gin.Context, status int, err string) {
	ctx.HTML(status, "error.tmpl", &errorPage{
		Error: err,
	})
}
