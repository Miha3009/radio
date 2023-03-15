package requests

type LoginRequest struct {
	Email    string `json:"email"`
	Password string `json:"password"`
}

type RegisterRequest struct {
	Name     string `json:"name"`
	Email    string `json:"email"`
	Password string `json:"password"`
}

type RefreshRequest struct {
	RefreshToken string
	IP           string
}

type LogoutRequest struct {
	UserID string
}

type DeleteUserRequset struct {
	UserID string
}

type ResetPasswordSendCodeRequest struct {
	Email string `json:"email"`
}

type ResetPasswordVerifyCodeRequest struct {
	Email string `json:"email"`
	Code  string `json:"code"`
}

type ResetPasswordChangeRequest struct {
	NewPassword string `json:"newPassword"`
	UserID      string
}

type UploadAvatarRequest struct {
	UserID string
	Avatar string
}
