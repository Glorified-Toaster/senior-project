package request

type CreateStudentRequest struct {
	FirstName  string `json:"first_name" validate:"required,min=2,max=32"`
	LastName   string `json:"last_name" validate:"required,min=2,max=32"`
	Department string `json:"department,omitempty"`
	StudentID  string `json:"student_id" validate:"required"`
	Email      string `json:"email" validate:"email,required"`
	Password   string `json:"password" validate:"required,min=8"`
}

type AdminResetPasswordRequest struct {
	StudentID   string `json:"student_id" validate:"required"`
	NewPassword string `json:"new_password" validate:"required,min=8"`
}

type LoginRequest struct {
	UserID   string `form:"user_id" json:"user_id" binding:"required"`
	Password string `form:"password" json:"password" binding:"required"`
	Role     string `form:"role" json:"role" binding:"required"`
}
