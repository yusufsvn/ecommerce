package httphandlers

import (
	e "ecommerceapi/CustomErrors"
	services "ecommerceapi/Services"
	st "ecommerceapi/Structs"
	"ecommerceapi/Utils"
	"encoding/json"
	"errors"
	"log"
	"net/http"

	"golang.org/x/oauth2"
)

type UserHandler struct {
	UserService       services.UserService
	GoogleOauthConfig *oauth2.Config
}

func NewUserHandler(_userService services.UserService) *UserHandler {
	return &UserHandler{
		UserService: _userService,
	}
}

func (u *UserHandler) Login(w http.ResponseWriter, r *http.Request) {
	// Parse body as struct
	var loginReq st.LoginReq
	Utils.ReadJSON(w, r, &loginReq)

	// service
	token, err := u.UserService.Login(r.Context(), loginReq, "user")
	// token gen

	if err != nil {
		http.Error(w, "giriş sırasında hata oluştu", http.StatusInternalServerError)
		return
	}

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

	response := st.TokenWithName{
		Name:    token.Name,
		Surname: token.Surname,
	}
	json.NewEncoder(w).Encode(response)
	w.Header().Set("Content-Type", "application/json")

	w.WriteHeader(http.StatusOK)
}

func (u *UserHandler) Register(w http.ResponseWriter, r *http.Request) {

	ctx := r.Context()

	var RequestDTO st.UserRegisterReq

	Utils.ReadJSON(w, r, &RequestDTO)

	tokenObj, err := u.UserService.RegisterUser(ctx, &RequestDTO)

	if err != nil && tokenObj == nil {
		switch {
		// İstemci kaynaklı hataları (400) yakalama
		case errors.Is(err, e.ErrInvalidEmail), errors.Is(err, e.ErrWeakPassword), errors.Is(err, e.ErrInvalidPhoneNumber):
			http.Error(w, err.Error(), http.StatusBadRequest)

		// Çakışma hataları (409) yakalama
		case errors.Is(err, e.ErrUserExists):
			http.Error(w, err.Error(), http.StatusConflict)

		// Bilinmeyen tüm hatalar için loglama ve 500 döndürme
		default:
			log.Println("err %w", err)
			http.Error(w, "Sunucu içi hata", http.StatusInternalServerError)
		}
		return
	}
	accessCookie := http.Cookie{
		Name:     "access_token",
		Value:    tokenObj.Access_token, // Kendi değişken ismine göre uyarla
		HttpOnly: true,
		Secure:   true,
		SameSite: http.SameSiteLaxMode,
		Path:     "/",
		MaxAge:   600,
	}

	refreshCookie := http.Cookie{
		Name:     "refresh_token",
		Value:    tokenObj.Refresh_token, // Kendi değişken ismine göre uyarla
		HttpOnly: true,
		Secure:   true,
		SameSite: http.SameSiteLaxMode,
		Path:     "/", // Güvenlik için sadece refresh endpoint'i ile sınırlandırılabilir (örneğin: "/api/refresh")
		MaxAge:   604800,
	}

	http.SetCookie(w, &accessCookie)
	http.SetCookie(w, &refreshCookie)
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)

}

func (u *UserHandler) GetUserRole(w http.ResponseWriter, r *http.Request) {

	isValid, err := u.UserService.IsValidUser(r.Context())
	if err != nil {
		http.Error(w, "Beklenmedik bir hata oluştu", http.StatusInternalServerError)
		return
	}

	if !isValid {
		http.Error(w, "Kullanıcı bulunamadı", http.StatusUnauthorized)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
}

func (u *UserHandler) GetUserCount(w http.ResponseWriter, r *http.Request) {

	count, err := u.UserService.GetUserCount(r.Context())

	if err != nil {
		http.Error(w, "Kullanıcı sayısı Getirilirken hata oluştu", http.StatusInternalServerError)
		return
	}

	json.NewEncoder(w).Encode(count)
}

// google oauth login

func (u *UserHandler) LoginWithGoogle(w http.ResponseWriter, r *http.Request) {

	var IdToken st.LoginGoogle
	Utils.ReadJSON(w, r, &IdToken)

	// service call
	tokendata, err := u.UserService.LoginWithGoogle(r.Context(), IdToken.ReturnIdToken)

	if err != nil {
		http.Error(w, "Google ile giriş yapılırken hata oluştu", http.StatusInternalServerError)
		return
	}
	accessCookie := http.Cookie{
		Name:     "access_token",
		Value:    tokendata.Access_token,
		HttpOnly: true,
		Secure:   false,
		SameSite: http.SameSiteLaxMode,
		Path:     "/",
		MaxAge:   600,
	}

	refreshCookie := http.Cookie{
		Name:     "refresh_token",
		Value:    tokendata.Refresh_token,
		HttpOnly: true,
		Secure:   false,
		SameSite: http.SameSiteLaxMode,
		Path:     "/",
		MaxAge:   604800,
	}

	http.SetCookie(w, &accessCookie)
	http.SetCookie(w, &refreshCookie)

	response := st.TokenWithName{
		Name:    tokendata.Name,
		Surname: tokendata.Surname,
	}
	json.NewEncoder(w).Encode(response)
	w.Header().Set("Content-Type", "application/json")

	w.WriteHeader(http.StatusOK)
}
