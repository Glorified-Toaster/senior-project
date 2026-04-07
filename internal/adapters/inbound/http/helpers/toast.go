package helpers

import (
	"uot-exam/web/templates/components/toast"

	"github.com/gin-gonic/gin"
)

func Toast(ctx *gin.Context, title string, description string, variant toast.Variant) {
	ctx.Writer.Write([]byte(`<div hx-swap-oob="beforeend:#toast-container">`))

	toast.Toast(toast.Props{
		Title:         title,
		Description:   description,
		Variant:       variant,
		Duration:      4000,
		ShowIndicator: true,
		Dismissible:   true,
		Icon:          true,
	}).Render(ctx.Request.Context(), ctx.Writer)

	ctx.Writer.Write([]byte(`</div>`))
}
