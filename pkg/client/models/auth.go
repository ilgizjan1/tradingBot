package models

type SignUpInput struct {
	Name          string `json:"name"`
	Username      string `json:"username"`
	Password      string `json:"password"`
	PublicAPIKey  string `json:"public_api_key"`
	PrivateAPIKey string `json:"private_api_key"`
}

type SignUpResponse struct {
	ID      int    `json:"id"`
	Message string `json:"message"`
}

type SignInInput struct {
	Username string `json:"username"`
	Password string `json:"password"`
}

type SignInResponse struct {
	AccessToken string `json:"access_token"`
	Message     string `json:"message"`
}

type LogoutInput struct {
	JWTToken string
}

type LogoutResponse struct {
	Message string `json:"message"`
}
