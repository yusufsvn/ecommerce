package services

import (
	"context"
	Customerrors "ecommerceapi/CustomErrors"
	db "ecommerceapi/Database/Repository"
	st "ecommerceapi/Structs"
	"ecommerceapi/Types"
	"ecommerceapi/Utils"
	"fmt"
	"log"
	"time"

	"github.com/golang-jwt/jwt/v5"
	"github.com/google/uuid"
)

type TokenService struct {
	secret_key string
	tokenFuncs db.TokenFunctions
}

func NewTokenService(tokenFuncs db.TokenFunctions, _secret_key string) *TokenService {

	return &TokenService{
		secret_key: _secret_key,
		tokenFuncs: tokenFuncs,
	}
}

func (ts *TokenService) GenToken(role string, expiresAt time.Time, family_id string, user_id string) (*st.TokenDTO, error) {

	jti := uuid.New().String()

	claims := jwt.MapClaims{
		"role":      role,
		"jti":       jti,
		"user_id":   user_id,
		"family_id": family_id,
		"iss":       "my-ecommerce",
		"aud":       []string{"core-backend"},
		"exp":       jwt.NewNumericDate(expiresAt),
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	token_string, err := token.SignedString([]byte(ts.secret_key))

	if err != nil {
		return nil, fmt.Errorf("token üretilemedi %w", err)
	}

	token_dto := st.TokenDTO{
		Token:      token_string,
		Expires_at: expiresAt,
		Jti:        jti,
	}
	return &token_dto, nil
}

func (ts *TokenService) SaveRefreshToken(ctx context.Context, token string, user_id string, expires_at time.Time, jti string, family_id string) error {

	token_hash, err := Utils.HashRefreshToken(token)

	log.Println("saverefreshtoken line 63")
	log.Println(token_hash)
	if err != nil {
		return fmt.Errorf("hashlemede sorunla karşılaşıldı %w", err)
	}
	token_model := st.RefreshTokenSaveDTO{
		Jti:        jti,
		User_id:    (user_id),
		Family_Id:  family_id,
		Is_used:    false,
		Is_revoked: false,
		Expires_at: expires_at,
		Token_hash: token_hash,
	}

	err1 := ts.tokenFuncs.AddToken(ctx, token_model)

	if err1 != nil {
		log.Printf("satır no 82 %v", err1)
		return fmt.Errorf("refresh token db ye kaydedilemedi %w", err1)
	}

	return nil
}

func (ts *TokenService) RefreshingToken(ctx context.Context) (*st.TokenWithName, error) {
	//db check refresh token
	token_props, ok := ctx.Value(Types.UserContextKey).(jwt.MapClaims)

	if !ok {
		return nil, fmt.Errorf("cannot fetch token claims: nil error or invalid type")
	}

	token_id, ok := token_props["jti"].(string)
	user_id, ok1 := token_props["user_id"].(string)
	role, ok2 := token_props["role"].(string)
	family_id, ok3 := token_props["family_id"].(string)
	if !ok {
		return nil, fmt.Errorf("jti is not a string or missing")
	}
	if !ok1 {
		return nil, fmt.Errorf("user_id is not a string or missing")
	}
	if !ok2 {
		return nil, fmt.Errorf("role is not a string or missing")
	}
	if !ok3 {
		return nil, fmt.Errorf("family id is not a string or missing")
	}

	// get token from db if has a matched result
	query_response, err := ts.tokenFuncs.GetTokenById(ctx, token_id)

	if err != nil || query_response == nil {
		log.Printf("satır no 114 %v", err)
		return nil, fmt.Errorf("token db den getirelemedi veya veri yok %w", err)
	}

	// Token kontrolleri
	if query_response.IsUsed != false || query_response.Is_revoked != false {
		ts.PanicToken(ctx, family_id, query_response.Is_revoked)
		// remove all token that are family id
		// not implemented
		return nil, Customerrors.ErrUnauthorized
	}

	// token oluştururken user_id alanı kodda eklenmedi
	if query_response.UserId != user_id {
		return nil, Customerrors.ErrUnauthorized
	}

	// Create New Tokens
	acExpiresAt := time.Now().Add(10 * time.Minute)
	reExpiresAt, errdate := token_props.GetExpirationTime()

	if errdate != nil {
		log.Printf("refresh tokenın expire time ı alınamadı:%v", errdate)
		return nil, errdate
	}

	new_access_token, err1 := ts.GenToken(role, acExpiresAt, "", user_id)
	new_refresh_token, err2 := ts.GenToken(role, reExpiresAt.Time, query_response.FamilyId, user_id)

	if err1 != nil {
		log.Printf("satır no 146 %v", err1)
		return nil, fmt.Errorf("AcToken Generate edilemedi %w", err1)
	}
	if err2 != nil {
		log.Printf("satır no 149 %v", err2)
		return nil, fmt.Errorf("ReToken Generate edilemedi %w", err2)
	}

	// Save the Db New Refresh Token
	ts.SaveRefreshToken(
		ctx,
		new_refresh_token.Token,
		query_response.UserId,
		new_refresh_token.Expires_at,
		new_access_token.Jti,
		family_id)

	// update old refresh token
	err4 := ts.tokenFuncs.UpdateToken(ctx, token_id, true)
	if err4 != nil {
		log.Printf("satır no 161 %v", err4)
		return nil, fmt.Errorf("update sorgusunda hatayla karşılaşıldı %w", err4)
	}
	// return the tokens

	newToken := st.TokenWithName{
		Access_token:  new_access_token.Token,
		Refresh_token: new_refresh_token.Token,
	}

	return &newToken, nil
}

func (ts *TokenService) PanicToken(ctx context.Context, family_id string, is_revoked bool) error {
	err := ts.tokenFuncs.PanicUpdate(ctx, family_id, is_revoked)

	if err != nil {
		return fmt.Errorf("Panic token update sorgusunda hata %w", err)
	}
	return nil
}
