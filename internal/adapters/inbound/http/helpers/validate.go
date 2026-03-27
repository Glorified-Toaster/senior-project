package helpers

import (
	"mime/multipart"
	"net/http"
	"uot-exam/web/templates/components/toast"

	"github.com/gin-gonic/gin"
)

func ValidateImage(ctx *gin.Context, questionImage *multipart.FileHeader) {
	if questionImage.Size > 8<<20 {
		ctx.Header("HX-Reswap", "none")
		Toast(ctx, "Create Question Failed", "Question image size exceeds 8MB limit", toast.VariantError)
		return
	}

	buffer := make([]byte, 512)
	file, err := questionImage.Open()
	if err != nil {
		ctx.Header("HX-Reswap", "none")
		Toast(ctx, "Create Question Failed", "Failed to open question image", toast.VariantError)
		return
	}
	defer file.Close()
	_, err = file.Read(buffer)
	if err != nil {
		ctx.Header("HX-Reswap", "none")
		Toast(ctx, "Create Question Failed", "Failed to read question image", toast.VariantError)
		return
	}
	contentType := http.DetectContentType(buffer)

	if contentType != "image/jpeg" && contentType != "image/png" {
		ctx.Header("HX-Reswap", "none")
		Toast(ctx, "Create Question Failed", "Question image must be a JPEG or PNG file", toast.VariantError)
		return
	}

	if _, err := file.Seek(0, 0); err != nil {
		ctx.Header("HX-Reswap", "none")
		Toast(ctx, "Error", "Failed to reset file pointer", toast.VariantError)
		return
	}

}
