package repository

import (
	"context"
	st "ecommerceapi/Structs"
	"fmt"
	"log"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
)

type TokenFunctions interface {
	AddToken(ctx context.Context, token st.RefreshTokenSaveDTO) error
	RemoveToken(ctx context.Context, token_id string) error
	GetTokenById(ctx context.Context, token_id string) (*st.TokenQueryRes, error)
	UpdateToken(ctx context.Context, token_id string, is_used bool) error
	PanicUpdate(ctx context.Context, family_id string, is_revoked bool) error
}

type TokenTable struct {
	db *pgxpool.Pool
}

func NewTokenStore(_db *pgxpool.Pool) *TokenTable {
	return &TokenTable{
		db: _db,
	}
}

func (tk *TokenTable) AddToken(ctx context.Context, token st.RefreshTokenSaveDTO) error {
	var query string = "insert into refreshtokens (token_id,user_id,token_hash,family_id,is_used,is_revoked,expires_at) values ($1,$2,$3,$4,$5,$6,$7)"

	log.Println("tokenRepo line 34")
	log.Println(token.Token_hash)
	_, err := tk.db.Exec(ctx, query, token.Jti, token.User_id, token.Token_hash, token.Family_Id, token.Is_used, token.Is_revoked, token.Expires_at)

	if err != nil {
		log.Printf("tokenRepo--AddToken--error:%v", err)
		return fmt.Errorf("token db'ye eklenemedi %w", err)
	}

	return nil
}

func (tk *TokenTable) RemoveToken(ctx context.Context, token_id string) error {

	var query string = "Delete from refreshtokens where token_id = $1"

	_, err := tk.db.Exec(ctx, query, token_id)

	if err != nil {
		log.Printf("tokenRepo--RemoveToken--error:%v", err)
		return fmt.Errorf("token silinirken hatayla karsilasildi %w", err)
	}
	return nil
}

func (tk *TokenTable) GetTokenById(ctx context.Context, token_id string) (*st.TokenQueryRes, error) {
	var query string = "select user_id,family_id,is_used,is_revoked from refreshtokens where token_id=$1"

	rows, err := tk.db.Query(ctx, query, token_id)

	if err != nil {
		log.Printf("tokenRepo--GetTokenById--error:%v", err)
		return nil, fmt.Errorf("token cannot get : db error %w", err)
	}

	query_response, err1 := pgx.CollectOneRow(rows, pgx.RowToStructByName[st.TokenQueryRes])

	if err1 != nil {
		log.Printf("tokenRepo--GetTokenById--error:%v", err1)
		return nil, fmt.Errorf("struct a dönüştürmede hata %w", err1)
	}

	return &query_response, nil

}

func (tk *TokenTable) UpdateToken(ctx context.Context, token_id string, is_used bool) error {

	var query string = "update refreshtokens set is_used=$1 where token_id=$2"

	_, err := tk.db.Exec(ctx, query, is_used, token_id)

	if err != nil {
		log.Printf("tokenRepo--UpdateToken--error:%v", err)
		return fmt.Errorf("update işleminde hata %w", err)
	}
	return nil
}

func (tk *TokenTable) PanicUpdate(ctx context.Context, family_id string, is_revoked bool) error {

	var query string = "update refreshtokens set is_revoked=$1 where family_id=$2"

	_, err := tk.db.Exec(ctx, query, is_revoked, family_id)

	if err != nil {
		log.Printf("tokenRepo--PanicUpdate--error:%v", err)
		return fmt.Errorf("update işleminde hata %w", err)
	}
	return nil

}
