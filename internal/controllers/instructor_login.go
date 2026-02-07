package controllers

import (
	"net/http"

	"github.com/Glorified-Toaster/senior-project/internal/templates"
	"github.com/Glorified-Toaster/senior-project/internal/utils"
	"github.com/gin-gonic/gin"
	csrf "github.com/utrack/gin-csrf"
)

func (ctrl *Controllers) InstructorLoginPageRender() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		csrfToken := csrf.GetToken(ctx)

		render := utils.NewRender(ctx, http.StatusOK, templates.InstructorLoginPage(csrfToken))
		ctx.Render(http.StatusOK, render)
	}
}
