package http_test

type request struct {
	AccessToken  string `json:"access_token"`
	RefreshToken string `json:"refresh_token"`
}

type response struct {
	Password     string `json:"password"`
	ClientSecret string `json:"client_secret"`
	Username     string `json:"username"`
}
