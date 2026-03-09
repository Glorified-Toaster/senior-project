package handler

import (
	"fmt"
	"net/http"
	"strconv"
	"uot-exam/internal/domain"
	"uot-exam/internal/ports"
	"uot-exam/web/templates/pages"
	"uot-exam/web/templates/pages/admin_dashboard/components"
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
			Users:           users,
			Username:        username,
			FullName:        fullname,
			Exams:           exams,
			TotalUsers:      fmt.Sprintf("%d", userCount),
			TotalUsersCount: userCount,
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
		limitStr := ctx.DefaultQuery("limit", "10")
		offsetStr := ctx.DefaultQuery("offset", "0")

		limit, _ := strconv.Atoi(limitStr)
		offset, _ := strconv.Atoi(offsetStr)

		users, err := h.userApp.ListAllUsers(ctx.Request.Context(), ports.ListAllUsersParams{
			Limit:  int32(limit),
			Offset: int32(offset),
		})
		if err != nil {
			ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}

		userCount, err := h.userApp.CountUsers(ctx)
		if err != nil {
			return
		}

		props := components.UserTableProps{
			Title:      "All Users",
			Users:      users,
			TotalCount: userCount,
			Limit:      int32(limit),
			Offset:     int32(offset),
			BaseURL:    "/admin/dashboard/users",
		}

		// Check if it's an HTMX request
		if ctx.GetHeader("HX-Request") != "" {
			render.Render(ctx, components.UserTableContainer(props))
			return
		}

		params := page.AdminDashboardParam{
			Users:           users,
			TotalUsers:      fmt.Sprintf("%d", userCount),
			TotalUsersCount: userCount,
		}
		render.Render(ctx, pages.BasePage("Admin Dashboard", page.AllUsers(params)))
	}
}
