package responses

type LoginResponse struct {
	AccessToken string `json:"accessToken"`
	Error       string `json:"error"`
}

type RefreshResponse struct {
	AccessToken string `json:"accessToken"`
	Error       string `json:"error"`
}

type RegisterResponse struct {
	AccessToken string `json:"accessToken"`
	Error       string `json:"error"`
}

type ResetPasswordVerifyCodeResponse struct {
	AccessToken string `json:"accessToken"`
	Error       string `json:"error"`
}

type ResetPasswordChangeResponse struct {
	Error string `json:"error"`
}
