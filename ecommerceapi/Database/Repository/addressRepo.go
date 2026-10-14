package repository

import (
	"context"
	st "ecommerceapi/Structs"
	"fmt"
	"log"
	"time"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
)

type AddressFunctions interface {
	Create(ctx context.Context, address *st.CreateAddressDTO) error
	Delete(ctx context.Context, addressId string) error
	Update(ctx context.Context, address st.Address, user_id string) error
	GetByUserId(ctx context.Context, user_id string) (*[]st.GetAddressDTO, error)
}

type AddressTable struct {
	db *pgxpool.Pool
}

func NewAddressStore(dbpool *pgxpool.Pool) *AddressTable {
	return &AddressTable{
		db: dbpool,
	}
}

func (s *AddressTable) Create(ctx context.Context, address *st.CreateAddressDTO) error {
	now := time.Now()

	var Query string = `Insert into Addresses(user_id, name, post_code, country, city, district, full_address, created_at) 
	Values ($1,$2,$3,$4,$5,$6,$7,$8)`
	_, err := s.db.Exec(ctx, Query,
		address.UserId,
		address.Name,
		address.PostCode,
		address.Country,
		address.City,
		address.District,
		address.FullAddress,
		now,
	)
	if err != nil {
		log.Printf("addressRepo--Create--error:%v", err)
		return fmt.Errorf("adres oluşturma hatası: %w", err)
	}
	return nil
}

func (s *AddressTable) Delete(ctx context.Context, addressId string) error {
	query := `Delete From Addresses where address_id=$1`
	_, err := s.db.Exec(ctx, query, addressId)
	if err != nil {
		log.Printf("addressRepo--Delete--error:%v", err)
		return fmt.Errorf("adres silme hatası: %w", err)
	}
	return nil
}

func (s *AddressTable) Update(ctx context.Context, address st.Address, user_id string) error {

	now := time.Now()

	var query string = `Update Addresses set 
	full_address=$1,post_code=$2,city=$3,district=$4,country=$5,updated_at=$6 where user_id=$7`

	_, err := s.db.Exec(ctx, query,
		address.AddressStr,
		address.ZipCode,
		address.City,
		address.District,
		address.Country,
		now,
		user_id)

	if err != nil {
		log.Printf("addressRepo--Update--error:%v", err)
		return fmt.Errorf("adress update edilirken hata olustu %w", err)
	}

	return nil
}

func (s *AddressTable) GetByUserId(ctx context.Context, user_id string) (*[]st.GetAddressDTO, error) {

	query := `select address_id,name,post_code,country,district,city,full_address from addresses where user_id=$1`

	rows, err := s.db.Query(ctx, query, user_id)

	if err != nil {
		log.Printf("db error--addressRepo--GetByUserId--%v", err)
		return nil, fmt.Errorf("db error")
	}

	data, err2 := pgx.CollectRows(rows, pgx.RowToStructByName[st.GetAddressDTO])

	if err2 != nil {
		log.Printf("collect rows error--addressRepo--GetByUserId--%v", err2)
		return nil, fmt.Errorf("collect rows error")
	}

	return &data, nil
}
