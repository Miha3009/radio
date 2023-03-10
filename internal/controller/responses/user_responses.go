package responses

type LoginResponse struct {
	AccessToken string `json:"accessToken"`
	Error       string `json:"error"`
	UserID      int    `json:"-"`
}

type RefreshResponse struct {
	AccessToken string `json:"accessToken"`
	Error       string `json:"error"`
	UserID      int    `json:"-"`
}

type RegisterResponse struct {
	AccessToken string `json:"accessToken"`
	Error       string `json:"error"`
	UserID      int    `json:"-"`
}

type ResetPasswordVerifyCodeResponse struct {
	AccessToken string `json:"accessToken"`
	Error       string `json:"error"`
	UserID      int    `json:"-"`
}

type ResetPasswordChangeResponse struct {
	Error string `json:"error"`
}
