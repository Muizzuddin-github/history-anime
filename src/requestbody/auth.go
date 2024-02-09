package requestbody

type Register struct {
	Username string `json:"username"`
	Email    string `json:"email"`
	Password string `json:"password"`
}

type Login struct {
	Email    string `json:"email" validate:"required"`
	Password string `json:"password" validate:"required"`
}

type ForgotPassword struct {
	Email string `json:"email" validate:"required"`
}

type ResetPassword struct {
	NewPassword string `json:"newPassword" validate:"required"`
	Token       string `json:"token" validate:"required"`
}
