package repository

import (
	"context"
	e "ecommerceapi/CustomErrors"
	Entities "ecommerceapi/Database/Models"
	st "ecommerceapi/Structs"
	"ecommerceapi/Utils"
	"errors"
	"fmt"
	"log"
	"time"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgconn"
	"github.com/jackc/pgx/v5/pgxpool"
)

// user entity

// fonksiyonlar için interface

type UserFunctions interface {
	GetAll(limit int, offset int, ctx context.Context) ([]Entities.User, error)
	Create(ctx context.Context, user *Entities.User) (string, error)
	Delete(ctx context.Context, userId *uuid.UUID) error
	Update(ctx context.Context, user *Entities.User) error
	GetUserByEmail(ctx context.Context, email string) (*st.UserLoginRes, error)
	GetUserById(ctx context.Context, userId string) (string, error)
	GetUser(ctx context.Context, user_id string) (*st.UserWithoutAddress, error)
	GetUserCount(ctx context.Context) (int, error)
	AddNewUserWithGoogle(ctx context.Context, user *st.GoogleUser) (string, error)
	GetUserByOauthId(ctx context.Context, oauthId string) (*st.UserLoginRes, error)
}

// db struct ı

type UserTable struct {
	db *pgxpool.Pool
}

func NewUserStore(dbpool *pgxpool.Pool) *UserTable {

	return &UserTable{
		db: dbpool,
	}
}

// Concrete Functions

func (s *UserTable) GetAll(limit int, offset int, ctx context.Context) ([]Entities.User, error) {
	query := "SELECT user_id, coalesce(password, '') as password, name, surname, coalesce(phone, '') as phone, email, user_role, created_at, updated_at, oauth_provider, oauth_id FROM users"
	if limit > 0 {
		query += fmt.Sprintf(" LIMIT %d OFFSET %d", limit, offset)
	}

	rows, err := s.db.Query(ctx, query)
	if err != nil {
		log.Printf("userRepo--GetAll--error:%v", err)
		return nil, fmt.Errorf("sorgu sırasında hata %w", err)
	}

	users, err := pgx.CollectRows(rows, pgx.RowToStructByName[Entities.User])
	if err != nil {
		log.Printf("userRepo--GetAll--error:%v", err)
		return nil, fmt.Errorf("satır dönüşüm hatası: %w", err)
	}
	return users, nil
}

func (s *UserTable) Create(ctx context.Context, user *Entities.User) (string, error) {

	var user_id string
	var Query string = "Insert into Users(password,name,surname,phone,email) Values ($1,$2,$3,$4,$5) Returning user_id"
	hashedpassw, err := Utils.HashPassword(user.Password)

	if err != nil {
		log.Printf("userRepo--Create--error:%v", err)
		return "", fmt.Errorf("parola hashlenemedi %w", err)
	}

	errExec := s.db.QueryRow(ctx, Query,
		hashedpassw,
		user.Name,
		user.Surname,
		user.Phone,
		user.Email,
	).Scan(&user_id)

	if errExec != nil {
		var pgErr *pgconn.PgError

		// unique contraint error
		if errors.As(errExec, &pgErr) {
			if pgErr.Code == "23505" {
				return "", e.ErrUserExists
			}
		}
		log.Printf("userRepo--Create--error:%v", errExec)
		return "", fmt.Errorf("user oluşturulamadı %w", errExec)
	}

	return user_id, nil
}

func (s *UserTable) Delete(ctx context.Context, userid *uuid.UUID) error {

	query := `Delete From Users where user_id=$1`
	_, err := s.db.Exec(ctx, query, userid)

	if err != nil {
		log.Printf("userRepo--Delete--error:%v", err)
		return fmt.Errorf("kullanıcı silme hatası: %w", err)
	}
	return nil
}

func (s *UserTable) Update(ctx context.Context, user *Entities.User) error {

	var Query string = `Update Users set password=$2,name=$3,surname=$4,phone=$5,email=$6,updated_at=$7 where user_id=$1`

	hashedpw, err := Utils.HashPassword(user.Password)

	if err != nil {
		log.Printf("userRepo--Update--error:%v", err)
		return fmt.Errorf("parola hashleme hatası: %w", err)
	}

	now := time.Now()
	user.UpdatedAt = &now

	_, queryerr := s.db.Exec(ctx, Query, user.UserID, hashedpw, user.Name, user.Surname, user.Phone, user.Email, user.UpdatedAt)

	if queryerr != nil {
		log.Printf("userRepo--Update--error:%v", queryerr)
		return fmt.Errorf("güncelleme sorgusu hatası: %w", queryerr)
	}
	return nil
}

func (s *UserTable) GetUserByEmail(ctx context.Context, email string) (*st.UserLoginRes, error) {
	var usrreq st.UserLoginRes
	query := `SELECT user_id, coalesce(password, ''), name, surname, user_role FROM users WHERE email = $1`
	err := s.db.QueryRow(ctx, query, email).Scan(&usrreq.UserID, &usrreq.Password, &usrreq.Name, &usrreq.Surname, &usrreq.Role)
	if err != nil {
		log.Printf("userRepo--GetUserByEmail--error:%v", err)
		return nil, fmt.Errorf("db error %w", err)
	}
	return &usrreq, nil
}

// just only get user role
func (s *UserTable) GetUserById(ctx context.Context, userId string) (string, error) {

	var userRole string

	query := "SELECT user_role FROM users WHERE user_id = $1"
	err := s.db.QueryRow(ctx, query, userId).Scan(&userRole)
	if err != nil {
		log.Printf("userRepo--GetUserById--error:%v", err)
		return "", fmt.Errorf("db error %w", err)
	}
	return userRole, nil
}

func (s *UserTable) GetUser(ctx context.Context, user_id string) (*st.UserWithoutAddress, error) {

	query := `select u.name, u.surname, u.email, coalesce(u.phone, '') as phone from users u where u.user_id=$1`

	rows, err := s.db.Query(ctx, query, user_id)

	if err != nil {
		log.Printf("db error--UserRepo--GetUser--%v", err)
		return nil, fmt.Errorf("db error %w", err)
	}

	data, err1 := pgx.CollectRows(rows, pgx.RowToAddrOfStructByName[st.UserWithoutAddress])

	if err1 != nil {
		log.Printf("collect rows error: %v--userRepo-GetUser", err1)
		return nil, fmt.Errorf("collect rows error: %w", err1)
	}

	if len(data) == 0 {
		return nil, fmt.Errorf("kullanıcı bulunamadı")
	}

	return data[0], nil
}

func (s *UserTable) GetUserCount(ctx context.Context) (int, error) {

	var count int

	query := `select count(*) from users`

	err := s.db.QueryRow(ctx, query).Scan(&count)

	if err != nil {
		return -1, err
	}

	// without admin
	return count - 1, nil

}

func (s *UserTable) AddNewUserWithGoogle(ctx context.Context, user *st.GoogleUser) (string, error) {

	query := `Insert into users(name,surname,email,oauth_provider,oauth_id) values($1,$2,$3,$4,$5)
	returning user_id`

	var user_id string
	err := s.db.QueryRow(ctx, query, user.Name, user.Surname, user.Email, "google", user.UserId).Scan(&user_id)

	if err != nil {
		log.Printf("userRepo--AddNewUserWithGoogle--error:%v", err)
		return "", fmt.Errorf("db error %w", err)
	}

	return user_id, nil
}

func (s *UserTable) GetUserByOauthId(ctx context.Context, oauthId string) (*st.UserLoginRes, error) {

	query := `select user_id::text,name,surname,email,user_role from users where oauth_id=$1 and oauth_provider='google'`

	var user st.UserLoginRes

	err := s.db.QueryRow(ctx, query, oauthId).Scan(&user.UserID, &user.Name, &user.Surname, &user.Email, &user.Role)

	if err != nil {
		log.Printf("userRepo--GetUserByOauthId--error:%v", err)
		return nil, fmt.Errorf("db error %w", err)
	}

	log.Print(user.UserID)
	log.Print(user.Name)

	return &user, nil
}
