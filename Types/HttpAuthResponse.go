package Types

type HttpAuthResponse struct {
	AccessToken  string `json:"accessToken"`
	RefreshToken string `json:"refreshToken"`
	HttpResponse
}
