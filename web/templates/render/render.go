package render

import (
	"bytes"
	"net/http"

	"github.com/a-h/templ"
	"github.com/gin-gonic/gin"
)

func Render(c *gin.Context, component templ.Component) {
	// If something was already written (like a Toast), we should just render directly
	// to avoid setting an incorrect Content-Length via ctx.Data.
	if c.Writer.Written() {
		_ = component.Render(c.Request.Context(), c.Writer)
		return
	}

	// Use a buffer to render the component to allow Gin to calculate the Content-Length
	// correctly for the response. This avoids HTTP/2 "wrote more than declared" errors.
	var buf bytes.Buffer
	if err := component.Render(c.Request.Context(), &buf); err != nil {
		c.Status(http.StatusInternalServerError)
		return
	}
	c.Data(http.StatusOK, "text/html; charset=utf-8", buf.Bytes())
}
