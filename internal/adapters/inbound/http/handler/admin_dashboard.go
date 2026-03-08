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

func (h *UserHandler) AdminDashboardMainRender() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		users, err := h.userApp.ListAllUsers(ctx.Request.Context(), ports.ListAllUsersParams{Limit: 4, Offset: 0})
		if err != nil {
			ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}
		usernameVal, exists := ctx.Get("username")
		if !exists {
			ctx.Redirect(http.StatusSeeOther, "/admin/login")
			return
		}
		username := usernameVal.(string)
		fullnameVal, exists := ctx.Get("fullname")
		if !exists {
			ctx.Redirect(http.StatusSeeOther, "/admin/login")
			fmt.Println("fullname not found")
			return
		}
		fullname := fullnameVal.(string)

		exams := []domain.Exam{}

		userCount, err := h.userApp.CountUsers(ctx)
		if err != nil {
			return
		}

		params := page.AdminDashboardParam{
			Users:      users,
			Username:   username,
			FullName:   fullname,
			Exams:      exams,
			TotalUsers: fmt.Sprintf("%d", userCount),
		}
		render.Render(ctx, pages.BasePage("Admin Dashboard", page.AdminMainPage(params)))
	}
}

func (h *UserHandler) AdminLogin() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		render.Render(ctx, pages.BasePage("Admin Login", page.LoginPage("")))
	}
}

func (h *UserHandler) UserPageRender() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		users, err := h.userApp.ListAllUsers(ctx.Request.Context(), ports.ListAllUsersParams{Limit: 10, Offset: 0})
		if err != nil {
			ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}
		userCount, err := h.userApp.CountUsers(ctx)
		if err != nil {
			return
		}
		params := page.AdminDashboardParam{
			Users:      users,
			TotalUsers: fmt.Sprintf("%d", userCount),
		}
		render.Render(ctx, pages.BasePage("Admin Dashboard", page.AllUsers(params)))
	}
}
