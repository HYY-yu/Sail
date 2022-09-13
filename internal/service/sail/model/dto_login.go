package model

type LoginParams struct {
	UserName string `json:"user_name" v:"required|length:3,20"`
	Password string `json:"password" v:"required|length:6,30"`
}

type LoginResponse struct {
	AccessToken  string `json:"access_token"`
	RefreshToken string `json:"refresh_token"`
	InitPassword bool   `json:"init_password"`
}
