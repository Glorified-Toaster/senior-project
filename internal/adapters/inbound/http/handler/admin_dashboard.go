package handler

import (
	"fmt"
	"net/http"
	"uot-exam/internal/domain"
	"uot-exam/internal/ports"
	"uot-exam/web/templates/pages"
	"uot-exam/web/templates/pages/admin_dashboard/page"
	"uot-exam/web/templates/render"

	"github.com/gin-gonic/gin"
)

func (h *UserHandler) AdminDashboardMain() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		users, err := h.userApp.ListAllUsers(ctx.Request.Context(), ports.ListAllUsersParams{Limit: 4, Offset: 0})
		if err != nil {
			ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}
		exams := []domain.Exam{}

		userCount, err := h.userApp.CountUsers(ctx)
		if err != nil {
			return
		}

		params := page.AdminMainDashboardParam{
			Users:      users,
			Exams:      exams,
			TotalUsers: fmt.Sprintf("%d", userCount),
		}
		render.Render(ctx, pages.BasePage("Admin Dashboard", page.AdminMainPage(params)))
	}
}
