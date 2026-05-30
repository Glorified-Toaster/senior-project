package middleware

import (
	"fmt"
	"net/http"
	"strings"
	"uot-exam/internal/adapters/inbound/http/helpers"
	"uot-exam/internal/domain"
	"uot-exam/web/templates/components/toast"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

func (m *AuthMiddleware) handleError(c *gin.Context, title, msg string, status int) {

	if status == http.StatusUnauthorized {
		c.Redirect(http.StatusSeeOther, "/login")
		c.Abort()
		return
	}

	if status == http.StatusForbidden {
		role, _ := c.Get("role")
		isInstructorPath := strings.HasPrefix(c.Request.URL.Path, "/instructor") || role == string(domain.RoleInstructor)

		if isInstructorPath {
			if c.GetHeader("HX-Request") == "true" {
				c.Header("HX-Reswap", "none")
				helpers.Toast(c, title, msg, toast.VariantError)
				c.Status(http.StatusOK)
				c.Abort()
				return
			}

			// For full page load, redirect to dashboard with error query parameter
			redirectURL := "/instructor/dashboard"
			c.Redirect(http.StatusSeeOther, fmt.Sprintf("%s?error=%s", redirectURL, msg))
			c.Abort()
			return
		}

		c.HTML(http.StatusForbidden, "403.html", gin.H{
			"Title":   title,
			"Message": msg,
		})
		c.Abort()
		return
	}

	// Fallback for other errors
	role, _ := c.Get("role")
	redirectURL := "/login"
	if role == string(domain.RoleStudent) {
		redirectURL = "/student/dashboard"
	} else if role == string(domain.RoleInstructor) {
		redirectURL = "/instructor/dashboard"
	} else if role == string(domain.RoleAdmin) {
		redirectURL = "/admin/dashboard"
	}

	c.Redirect(http.StatusSeeOther, fmt.Sprintf("%s?error=%s", redirectURL, msg))
	c.Abort()
}

// SubjectAccessMiddleware ensures the student is enrolled or the instructor is assigned to the subject.
// This middleware expects a ":id" parameter that represents the Subject ID.
func (m *AuthMiddleware) SubjectAccessMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		role, _ := c.Get("role")
		userIDStr, _ := c.Get("userID")

		userID, err := uuid.Parse(userIDStr.(string))
		if err != nil {
			m.handleError(c, "Authentication Error", "Invalid user ID", http.StatusUnauthorized)
			return
		}

		// Admins have full access
		if role == string(domain.RoleAdmin) {
			c.Next()
			return
		}

		subjectIDStr := c.Param("id")
		subjectID, err := uuid.Parse(subjectIDStr)
		if err != nil {
			m.handleError(c, "Invalid Request", "Invalid subject ID", http.StatusBadRequest)
			return
		}

		if role == string(domain.RoleStudent) {
			enrolled, err := m.app.IsStudentEnrolled(c.Request.Context(), subjectID, userID)
			if err != nil || !enrolled {
				m.logger.LogErrorWithLevel("warn", "AUTH_ERROR", "UNAUTHORIZED_SUBJECT_ACCESS", "Student attempted to access non-enrolled subject", nil)
				m.handleError(c, "Access Denied", "You are not enrolled in this subject", http.StatusForbidden)
				return
			}
		} else if role == string(domain.RoleInstructor) {
			assigned, err := m.app.IsInstructorAssigned(c.Request.Context(), subjectID, userID)
			if err != nil || !assigned {
				m.logger.LogErrorWithLevel("warn", "AUTH_ERROR", "UNAUTHORIZED_SUBJECT_ACCESS", "Instructor attempted to access non-assigned subject", nil)
				m.handleError(c, "Access Denied", "You are not assigned to this subject", http.StatusForbidden)
				return
			}

			if c.Request.Method != http.MethodGet && !strings.HasSuffix(c.Request.URL.Path, "/search") {
				subject, err := m.app.GetSubjectByID(c.Request.Context(), subjectID)
				if err == nil && subject.Status == domain.SubjectStatusPublished {
					m.handleError(c, "Forbidden", "Cannot modify a published subject", http.StatusForbidden)
					return
				}
			}
		}

		c.Next()
	}
}

// ExamAccessMiddleware ensures the user has access to the subject the exam belongs to.
// This middleware expects a ":id" parameter that represents the Exam ID.
func (m *AuthMiddleware) ExamAccessMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		role, _ := c.Get("role")
		userIDStr, _ := c.Get("userID")

		userID, err := uuid.Parse(userIDStr.(string))
		if err != nil {
			m.handleError(c, "Authentication Error", "Invalid user ID", http.StatusUnauthorized)
			return
		}

		// Admins have full access
		if role == string(domain.RoleAdmin) {
			c.Next()
			return
		}

		examIDStr := c.Param("id")
		examID, err := uuid.Parse(examIDStr)
		if err != nil {
			m.handleError(c, "Invalid Request", "Invalid exam ID", http.StatusBadRequest)
			return
		}

		// Get the exam to find its subject ID
		exam, err := m.app.GetExamByID(c.Request.Context(), examID)
		if err != nil {
			m.handleError(c, "Not Found", "Exam not found", http.StatusNotFound)
			return
		}

		if role == string(domain.RoleStudent) {
			enrolled, err := m.app.IsStudentEnrolled(c.Request.Context(), exam.SubjectID, userID)
			if err != nil || !enrolled {
				m.logger.LogErrorWithLevel("warn", "AUTH_ERROR", "UNAUTHORIZED_EXAM_ACCESS", "Student attempted to access exam in non-enrolled subject", nil)
				m.handleError(c, "Access Denied", "You are not enrolled in the subject for this exam", http.StatusForbidden)
				return
			}
		} else if role == string(domain.RoleInstructor) {
			// Check if the instructor is the creator of the exam
			if exam.CreatedBy != userID {
				m.logger.LogErrorWithLevel("warn", "AUTH_ERROR", "UNAUTHORIZED_EXAM_ACCESS", "Instructor attempted to access exam they did not create", nil)
				m.handleError(c, "Access Denied", "You can only access exams you have created", http.StatusForbidden)
				return
			}

			// Also ensure they are still assigned to the subject (extra security)
			assigned, err := m.app.IsInstructorAssigned(c.Request.Context(), exam.SubjectID, userID)
			if err != nil || !assigned {
				m.logger.LogErrorWithLevel("warn", "AUTH_ERROR", "UNAUTHORIZED_EXAM_ACCESS", "Instructor attempted to access exam in non-assigned subject", nil)
				m.handleError(c, "Access Denied", "You are no longer assigned to the subject for this exam", http.StatusForbidden)
				return
			}

			if c.Request.Method != http.MethodGet && !strings.HasSuffix(c.Request.URL.Path, "/search") {
				subject, err := m.app.GetSubjectByID(c.Request.Context(), exam.SubjectID)
				if err == nil && subject.Status == domain.SubjectStatusPublished {
					m.handleError(c, "Forbidden", "Cannot modify an exam in a published subject", http.StatusForbidden)
					return
				}
			}
		}

		c.Next()
	}
}
