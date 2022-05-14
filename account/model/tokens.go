package model

// Token used for returning pairs of id and refresh tokens
type Token struct {
	IDToken      string `json:"idToken"`
	RefreshToken string `json:"refreshToken"`
}