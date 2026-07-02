package platformuser

type ResetPasswordResponse struct {
	UserID             string `json:"user_id"`
	Login              string `json:"login"`
	TemporaryPassword  string `json:"temporary_password"`
	MustChangePassword bool   `json:"must_change_password"`
}
