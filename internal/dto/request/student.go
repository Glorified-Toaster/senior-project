package request

type CreateStudentRequest struct {
	FirstName  string `json:"first_name" validate:"required,min=2,max=32"`
	LastName   string `json:"last_name" validate:"required,min=2,max=32"`
	Department string `json:"department,omitempty"`
	StudentID  string `json:"student_id" validate:"required"`
	Email      string `json:"email" validate:"email,required"`
	Password   string `json:"password" validate:"required,min=8"`
}

type StudentLoginRequest struct {
	Email     string `json:"email" validate:"required,email"`
	StudentID string `json:"student_id" validate:"required"`
	Password  string `json:"password" validate:"required"`
}

type AdminResetPasswordRequest struct {
	StudentID   string `json:"student_id" validate:"required"`
	NewPassword string `json:"new_password" validate:"required,min=8"`
}
