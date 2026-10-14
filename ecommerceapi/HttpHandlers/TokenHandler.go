package httphandlers

import (
	s "ecommerceapi/Services"
	"log"
	"net/http"
)

type TokenHandler struct {
	tokenService s.TokenService
}

func NewTokenHandler(_tokenService s.TokenService) *TokenHandler {
	return &TokenHandler{
		tokenService: _tokenService,
	}
}

func (th *TokenHandler) GetNewAccessToken(w http.ResponseWriter, r *http.Request) {

	log.Println("refresh token endpoint i ------------------")
	tokens, err1 := th.tokenService.RefreshingToken(r.Context())
	if err1 != nil {
		log.Printf("token refresh hatası %v", err1)
		http.Error(w, "geçersiz refresh token", http.StatusUnauthorized)
		return
	}

	accessCookie := http.Cookie{
		Name:     "access_token",
		Value:    tokens.Access_token,
		HttpOnly: true,
		Secure:   false,
		SameSite: http.SameSiteLaxMode,
		Path:     "/",
		MaxAge:   600,
	}

	refreshCookie := http.Cookie{
		Name:     "refresh_token",
		Value:    tokens.Refresh_token,
		HttpOnly: true,
		Secure:   false,
		SameSite: http.SameSiteLaxMode,
		Path:     "/",
		MaxAge:   604800,
	}

	http.SetCookie(w, &accessCookie)
	http.SetCookie(w, &refreshCookie)
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(200)

}
