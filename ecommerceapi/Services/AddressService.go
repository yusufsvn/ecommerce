package services

import (
	"context"
	db "ecommerceapi/Database/Repository"
	st "ecommerceapi/Structs"
	"ecommerceapi/Utils"
	"fmt"
	"log"
)

type AddressService struct {
	AddressFuncs db.AddressFunctions
}

func NewAddressService(_addressFuncs db.AddressFunctions) *AddressService {
	return &AddressService{
		AddressFuncs: _addressFuncs,
	}
}

func (as *AddressService) GetAddress(ctx context.Context) (*st.GetAddressList, error) {

	user_id, err1 := Utils.GetUserIdFromClaims(ctx)

	if err1 != nil {
		return nil, err1
	}

	data, err := as.AddressFuncs.GetByUserId(ctx, user_id)

	if err != nil {
		log.Printf("UserService--GetAddress--%v", err)
		return nil, fmt.Errorf("adres alınamadı")
	}

	datalist := st.GetAddressList{
		Addresses: data,
	}

	return &datalist, nil

}

func (as *AddressService) UpdateAddress(ctx context.Context, addressreq st.Address) error {

	user_id, err1 := Utils.GetUserIdFromClaims(ctx)

	if err1 != nil {
		return err1
	}

	err := as.AddressFuncs.Update(ctx, addressreq, user_id)

	if err != nil {
		log.Printf("db error-- updateaddress--%v", err)
		return fmt.Errorf("db error")
	}

	return nil
}

func (as *AddressService) CreateAddress(ctx context.Context, formData st.Address) error {

	user_id, err := Utils.GetUserIdFromClaims(ctx)

	if err != nil {
		log.Printf("CreateAddress : %v", err)
		return err
	}

	dto := st.CreateAddressDTO{
		UserId:      user_id,
		Name:        formData.AddressName,
		PostCode:    formData.ZipCode,
		Country:     formData.Country,
		District:    formData.District,
		City:        formData.City,
		FullAddress: formData.AddressStr,
	}

	err1 := as.AddressFuncs.Create(ctx, &dto)

	if err1 != nil {
		log.Printf("CreateAddress DB: %v", err1)
		return err1
	}

	return nil
}

func (as *AddressService) DeleteAddress(ctx context.Context, id string) error {

	err := as.AddressFuncs.Delete(ctx, id)

	if err != nil {
		return err
	}

	return nil
}
