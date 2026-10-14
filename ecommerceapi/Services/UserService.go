package services

import (
	"context"
	e "ecommerceapi/CustomErrors"
	models "ecommerceapi/Database/Models"
	db "ecommerceapi/Database/Repository"
	st "ecommerceapi/Structs"
	"ecommerceapi/Types"
	"ecommerceapi/Utils"
	"fmt"
	"log"
	"strings"
	"time"

	"github.com/golang-jwt/jwt/v5"
	"github.com/google/uuid"
	"google.golang.org/api/idtoken"
)

type UserService struct {
	UserFunc           db.UserFunctions
	Token_service      TokenService
	CartFuncs          db.CartFunctions
	AddressFuncs       db.AddressFunctions
	GoogleAuthClientId string
}

func NewUserService(
	_userfunc db.UserFunctions,
	_token_service TokenService,
	_cartfunc db.CartFunctions,
	_addressFuncs db.AddressFunctions,
	_googleAuthClientId string) *UserService {
	return &UserService{
		UserFunc:           _userfunc,
		Token_service:      _token_service,
		CartFuncs:          _cartfunc,
		AddressFuncs:       _addressFuncs,
		GoogleAuthClientId: _googleAuthClientId,
	}
}

func (u *UserService) RegisterUser(ctx context.Context, request *st.UserRegisterReq) (*st.TokenWithName, error) {
	// 1. kontroller yapılsın : mail, telefon, parola doğrulama
	if request.Password != request.Password_check {
		return nil, e.ErrPasswordNotMatched
	}

	// 1.1 normalize email
	request.Mail = strings.TrimSpace(request.Mail)
	request.Mail = strings.ToLower(request.Mail)

	new_user := models.User{
		Password: request.Password,
		Name:     request.Name,
		Surname:  request.SecondName,
		Phone:    "",
		Email:    request.Mail,
	}

	// 2. veritabanına kayıt
	user_id, err1 := u.UserFunc.Create(ctx, &new_user)

	if err1 != nil {
		return nil, fmt.Errorf("user oluşturulamadi %w", err1)
	}

	_, err := u.CartFuncs.Create(ctx, user_id)
	if err != nil {
		log.Printf("Sepet oluşturulamadı %v", err)
		return nil, fmt.Errorf("user oluşturulamadi")
	}

	// 3. token üretimi
	family_id := uuid.New().String()
	reExpiresAt := time.Now().Add(7 * 24 * time.Minute)
	acExpiresAt := time.Now().Add(10 * time.Minute)
	access_token, err := u.Token_service.GenToken("user", acExpiresAt, family_id, user_id)
	refresh_token, err2 := u.Token_service.GenToken("user", reExpiresAt, family_id, user_id)

	if err != nil {
		return nil, fmt.Errorf("access token oluşturulamadı %w", err)
	}

	if err2 != nil {
		return nil, fmt.Errorf("refresh token oluşturulamadı %w", err2)
	}

	Token := st.TokenWithName{
		Access_token:  access_token.Token,
		Refresh_token: refresh_token.Token,
		Name:          request.Name,
		Surname:       request.SecondName,
	}

	//save refresh token to db
	errT := u.Token_service.SaveRefreshToken(ctx, refresh_token.Token, user_id, refresh_token.Expires_at, refresh_token.Jti, family_id)

	if errT != nil {
		return nil, fmt.Errorf("token veritabanına kaydedilemedi")
	}
	// 4. token gönderme ve yeniden yönlendirme
	return &Token, nil
}

func (u *UserService) Login(ctx context.Context, loginReq st.LoginReq, role string) (*st.TokenWithName, error) {

	data, err := u.UserFunc.GetUserByEmail(ctx, loginReq.Mail)

	if err != nil {
		log.Println("giriş sırasında hata oluştu line 93")
		return nil, fmt.Errorf("giriş sırasında hata oluştu %w", err)
	}

	if data.Role != role {
		log.Println("hatalı rol")
		return nil, fmt.Errorf("giriş sırasında hata oluştu %w", err)
	}

	if !Utils.CheckPasswordHash(loginReq.Password, data.Password) {
		return nil, fmt.Errorf("Parola yanlış")
	}

	// token things
	family_id := uuid.New().String()

	reExpiresAt := time.Now().Add((7 * 24) * time.Hour)
	//reExpiresAt := time.Now().Add(30 * time.Minute)
	acExpiresAt := time.Now().Add(10 * time.Minute)
	access_token, err1 := u.Token_service.GenToken(role, acExpiresAt, family_id, data.UserID)
	refresh_token, err2 := u.Token_service.GenToken(role, reExpiresAt, family_id, data.UserID)

	if err1 != nil {
		return nil, fmt.Errorf("access token oluşturulamadı %w", err1)
	}

	if err2 != nil {
		return nil, fmt.Errorf("refresh token oluşturulamadı %w", err2)
	}

	errT := u.Token_service.SaveRefreshToken(ctx, refresh_token.Token, data.UserID, refresh_token.Expires_at, refresh_token.Jti, family_id)

	if errT != nil {
		return nil, fmt.Errorf("token veritabanına kaydedilemedi")
	}
	return &st.TokenWithName{
		Access_token:  access_token.Token,
		Refresh_token: refresh_token.Token,
		Name:          data.Name,
		Surname:       data.Surname,
	}, nil
}

// Validating user role compare with db value
func (u *UserService) IsValidUser(ctx context.Context) (bool, error) {

	user_id, err_id := Utils.GetUserIdFromClaims(ctx)

	if err_id != nil || user_id == "" {
		log.Printf("user_id maybe zero length or invalid claims:%v", err_id)
		return false, fmt.Errorf("Beklenmeyen bir hatayla karşılaşıldı")
	}
	var role string = ctx.Value(Types.UserContextKey).(jwt.MapClaims)["role"].(string)
	userRole, err := u.UserFunc.GetUserById(ctx, user_id)

	if err != nil {
		log.Printf("hatayla karşılaşıldı %v", err)
		return false, fmt.Errorf("error")
	}

	if userRole != role {
		log.Printf("hatalı rol")
		return false, nil
	}

	return true, nil
}

func (u *UserService) GetUserCount(ctx context.Context) (int, error) {

	count, err := u.UserFunc.GetUserCount(ctx)

	if err != nil {
		return -1, err
	}

	return count, nil
}

func (u *UserService) VerifyGoogleIdToken(ctx context.Context, Token string) (*st.GoogleUser, error) {

	if Token == "" {
		return nil, e.ErrEmptyParameter
	}

	tokenPayload, err := idtoken.Validate(ctx, Token, u.GoogleAuthClientId)

	if err != nil {
		log.Printf("cannot parse idToken :%v", err)
		return nil, err
	}

	email, _ := tokenPayload.Claims["email"].(string)
	name, _ := tokenPayload.Claims["name"].(string)
	// Unique Google user ID
	sub, _ := tokenPayload.Claims["sub"].(string)
	// parse name
	namelist := strings.Fields(name)

	UserInfo := st.GoogleUser{
		UserId:  sub,
		Name:    namelist[0],
		Surname: namelist[1],
		Email:   email,
	}

	return &UserInfo, nil
}

func (u *UserService) LoginWithGoogle(ctx context.Context, IdToken string) (*st.TokenWithName, error) {

	userInfo, err := u.VerifyGoogleIdToken(ctx, IdToken)
	var userId string
	if err != nil {
		return nil, err
	}

	data, err1 := u.UserFunc.GetUserByOauthId(ctx, userInfo.UserId)

	if err1 != nil || data == nil {
		var err1 error
		userId, err1 = u.UserFunc.AddNewUserWithGoogle(ctx, userInfo)

		if err1 != nil || userId == "" {
			log.Printf("user_id is empty or error occurred: %v", err1)
			return nil, err1
		}

		_, err := u.CartFuncs.Create(ctx, userId)
		if err != nil {
			log.Printf("Sepet oluşturulamadı %v", err)
			return nil, fmt.Errorf("user oluşturulamadi")
		}

		Tokens, errT1 := u.InitToken(ctx, userId, userInfo.Name, userInfo.Surname)

		if errT1 != nil {
			return nil, errT1
		}

		return Tokens, nil

	}

	userId = data.UserID
	Token, errT2 := u.InitToken(ctx, userId, userInfo.Name, userInfo.Surname)

	if errT2 != nil {
		return nil, errT2
	}

	return Token, nil

}

func (u *UserService) InitToken(ctx context.Context, userId, userName, userSurname string) (*st.TokenWithName, error) {

	family_id := uuid.New().String()
	reExpiresAt := time.Now().Add(7 * 24 * time.Minute)
	acExpiresAt := time.Now().Add(10 * time.Minute)
	access_token, err := u.Token_service.GenToken("user", acExpiresAt, family_id, userId)
	refresh_token, err2 := u.Token_service.GenToken("user", reExpiresAt, family_id, userId)

	if err != nil {
		return nil, fmt.Errorf("access token oluşturulamadı %w", err)
	}

	if err2 != nil {
		return nil, fmt.Errorf("refresh token oluşturulamadı %w", err2)
	}

	Token := st.TokenWithName{
		Access_token:  access_token.Token,
		Refresh_token: refresh_token.Token,
		Name:          userName,
		Surname:       userSurname,
	}

	//save refresh token to db
	errT := u.Token_service.SaveRefreshToken(ctx, refresh_token.Token, userId, refresh_token.Expires_at, refresh_token.Jti, family_id)

	if errT != nil {
		return nil, fmt.Errorf("token veritabanına kaydedilemedi")
	}

	return &Token, nil
}
