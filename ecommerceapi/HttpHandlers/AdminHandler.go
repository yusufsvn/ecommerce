package httphandlers

import (
	services "ecommerceapi/Services"
	st "ecommerceapi/Structs"
	"ecommerceapi/Utils"
	"log"
	"net/http"
)

type AdminHandler struct {
	userService services.UserService
}

func NewAdminHandler(_userService services.UserService) *AdminHandler {
	return &AdminHandler{
		userService: _userService,
	}
}

// admin giriş yapar token generate edilir. access token ve refresh token gönderilir

func (Ad *AdminHandler) Login(w http.ResponseWriter, r *http.Request) {

	// Parse body as struct
	var loginReq st.LoginReq
	Utils.ReadJSON(w, r, &loginReq)

	log.Printf("mail %s", loginReq.Mail)
	log.Printf("password %s", loginReq.Password)
	// service
	token, err := Ad.userService.Login(r.Context(), loginReq, "admin")
	// token gen

	if err != nil {
		http.Error(w, "giriş sırasında hata oluştu", http.StatusInternalServerError)
		return
	}

	log.Printf("token %s", token.Access_token)

	accessCookie := http.Cookie{
		Name:     "access_token",
		Value:    token.Access_token,
		HttpOnly: true,
		Secure:   false,
		SameSite: http.SameSiteLaxMode,
		Path:     "/",
		MaxAge:   600,
	}

	refreshCookie := http.Cookie{
		Name:     "refresh_token",
		Value:    token.Refresh_token,
		HttpOnly: true,
		Secure:   false,
		SameSite: http.SameSiteLaxMode,
		Path:     "/",
		MaxAge:   604800,
	}

	http.SetCookie(w, &accessCookie)
	http.SetCookie(w, &refreshCookie)
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
}
