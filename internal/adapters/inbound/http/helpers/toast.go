package helpers

import (
	"uot-exam/web/templates/components/toast"

	"github.com/gin-gonic/gin"
)

func Toast(ctx *gin.Context, title string, description string, variant toast.Variant) {
	// 1. Start an OOB div. 'hx-swap-oob="beforeend:#toast-container"'
	// tells HTMX to append this HTML to your toast container,
	// ignoring the main swap rules.
	ctx.Writer.Write([]byte(`<div hx-swap-oob="beforeend:#toast-container">`))

	// 2. Render your existing component
	toast.Toast(toast.Props{
		Title:         title,
		Description:   description,
		Variant:       variant,
		Duration:      4000,
		ShowIndicator: true,
		Dismissible:   true,
		Icon:          true,
	}).Render(ctx.Request.Context(), ctx.Writer)

	// 3. Close the div
	ctx.Writer.Write([]byte(`</div>`))
}
