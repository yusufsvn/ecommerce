package Customerrors

import "errors"

var (
	ErrInvalidEmail       = errors.New("geçersiz e-posta formatı")
	ErrWeakPassword       = errors.New("parola güvenlik standartlarını karşılamıyor")
	ErrUserExists         = errors.New("bu e-posta adresi zaten kayıtlı")
	ErrPasswordNotMatched = errors.New("parolalar eşleşmiyor")
	ErrInvalidPhoneNumber = errors.New("geçersiz telefon numarası")
	ErrUnauthorized       = errors.New("erişim reddedildi")
	ErrInvalidToken       = errors.New("geçersiz token")
	ZeroListResult        = errors.New("sepet boş")
	ErrInvalidItemType    = errors.New("Ürün tipini bu şekilde girin PHYSICAL veya VIRTUAL")
	ErrEmptyParameter     = errors.New("parametreler boş")
)
