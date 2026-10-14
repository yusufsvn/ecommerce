package services

import (
	structs "ecommerceapi/Structs"
	"testing"

	"github.com/JspBack/iyzipay-go"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestPaymentService_InitializePayment(t *testing.T) {
	client, err := iyzipay.New("test-api-key", "test-api-secret")
	require.NoError(t, err)
	svc := NewPaymentService(client)

	t.Run("tekil ve çoklu miktardaki sepet ürünlerini doğru şekilde eşlemeli", func(t *testing.T) {
		req := structs.CreateOrderDTO{
			Order_id:   "ORD-12345",
			TotalPrice: "1500",
			Cart_id:    "CART-999",
			ReqData: structs.CreatePaymentDto{
				Tckn:    "11111111110",
				User_ip: "192.168.1.1",
				CardData: structs.Card{
					CardName:    "Ahmet Yılmaz",
					CardNumber:  "5528790000000001",
					ExpireYear:  "2028",
					ExpireMonth: "12",
					Cvc:         "123",
				},
				ShippingAddress: structs.Address{
					AddressStr: "Bağdat Cad. No:10",
					City:       "İstanbul",
					Country:    "Türkiye",
					ZipCode:    "34710",
				},
				BillingAddress: structs.Address{
					AddressStr: "Atatürk Bulv. No:5",
					City:       "Ankara",
					Country:    "Türkiye",
					ZipCode:    "06000",
				},
			},
			UserData: structs.UserOrderDTO{
				User_id: "USR-001",
				Name:    "Ahmet",
				Surname: "Yılmaz",
				Mail:    "ahmet@example.com",
				Phone:   "05321112233",
				Address: structs.Address{
					AddressStr: "Kayıtlı Adres",
					City:       "İzmir",
					Country:    "Türkiye",
					ZipCode:    "35000",
				},
			},
			CartItems: []structs.ProductsFromCartDTO{
				{
					ProductID:   "prod-1",
					ProductName: "Kitap",
					Price:       100,
					Quantity:    1,
					ItemType:    "PHYSICAL",
					Category:    "Yayın",
				},
				{
					ProductID:   "prod-2",
					ProductName: "Kalem",
					Price:       50,
					Quantity:    3, // 3 adet -> 3 ayrı basket item oluşturulmalı
					ItemType:    "PHYSICAL",
					Category:    "Kırtasiye",
				},
			},
		}

		paymentReq := svc.InitializePayment(req)

		assert.Equal(t, "tr", paymentReq.Locale)
		assert.Equal(t, "ORD-12345", paymentReq.ConversationID)
		assert.Equal(t, "1500", paymentReq.Price)
		assert.Equal(t, "1500", paymentReq.PaidPrice)
		assert.Equal(t, 1, paymentReq.Installment)
		assert.Equal(t, "WEB", paymentReq.PaymentChannel)
		assert.Equal(t, "CART-999", paymentReq.BasketID)
		assert.Equal(t, "PRODUCT", paymentReq.PaymentGroup)
		assert.Equal(t, "TRY", paymentReq.Currency)

		// Kart bilgileri
		assert.Equal(t, "Ahmet Yılmaz", paymentReq.PaymentCard.CardHolderName)
		assert.Equal(t, "5528790000000001", paymentReq.PaymentCard.CardNumber)
		assert.Equal(t, "2028", paymentReq.PaymentCard.ExpireYear)
		assert.Equal(t, "12", paymentReq.PaymentCard.ExpireMonth)
		assert.Equal(t, "123", paymentReq.PaymentCard.CVC)
		assert.Equal(t, 0, paymentReq.PaymentCard.RegisterCard)

		// Alıcı (Buyer) bilgileri
		assert.Equal(t, "USR-001", paymentReq.Buyer.ID)
		assert.Equal(t, "Ahmet", paymentReq.Buyer.Name)
		assert.Equal(t, "Yılmaz", paymentReq.Buyer.Surname)
		assert.Equal(t, "11111111110", paymentReq.Buyer.IdentityNumber)
		assert.Equal(t, "ahmet@example.com", paymentReq.Buyer.Email)
		assert.Equal(t, "05321112233", paymentReq.Buyer.GSMNumber)
		assert.Equal(t, "Kayıtlı Adres", paymentReq.Buyer.RegistrationAddress)
		assert.Equal(t, "İzmir", paymentReq.Buyer.City)
		assert.Equal(t, "Türkiye", paymentReq.Buyer.Country)
		assert.Equal(t, "35000", paymentReq.Buyer.ZipCode)
		assert.Equal(t, "192.168.1.1", paymentReq.Buyer.IP)

		// Teslimat Adresi
		assert.Equal(t, "Bağdat Cad. No:10", paymentReq.ShippingAddress.Address)
		assert.Equal(t, "Ahmet Yılmaz", paymentReq.ShippingAddress.ContactName)
		assert.Equal(t, "İstanbul", paymentReq.ShippingAddress.City)
		assert.Equal(t, "Türkiye", paymentReq.ShippingAddress.Country)
		assert.Equal(t, "34710", paymentReq.ShippingAddress.ZipCode)

		// Fatura Adresi
		assert.Equal(t, "Atatürk Bulv. No:5", paymentReq.BillingAddress.Address)
		assert.Equal(t, "Ahmet Yılmaz", paymentReq.BillingAddress.ContactName)
		assert.Equal(t, "Ankara", paymentReq.BillingAddress.City)
		assert.Equal(t, "Türkiye", paymentReq.BillingAddress.Country)
		assert.Equal(t, "06000", paymentReq.BillingAddress.ZipCode)

		// Sepet Ürünleri (1 adet prod-1 + 3 adet prod-2 = toplam 4 adet BasketItem)
		assert.Len(t, paymentReq.BasketItems, 4)
		assert.Equal(t, "prod-1", paymentReq.BasketItems[0].ID)
		assert.Equal(t, "Kitap", paymentReq.BasketItems[0].Name)
		assert.Equal(t, "100", paymentReq.BasketItems[0].Price)
		assert.Equal(t, "prod-2", paymentReq.BasketItems[1].ID)
		assert.Equal(t, "prod-2", paymentReq.BasketItems[2].ID)
		assert.Equal(t, "prod-2", paymentReq.BasketItems[3].ID)
	})
}

func TestPaymentService_CreatePaymentNonTDS(t *testing.T) {
	client, err := iyzipay.New("invalid-key", "invalid-secret")
	require.NoError(t, err)
	svc := NewPaymentService(client)

	t.Run("geçersiz kimlik bilgileri veya ağ hatasında false dönmeli", func(t *testing.T) {
		req := structs.CreateOrderDTO{
			Order_id:   "ORD-FAIL",
			TotalPrice: "100",
			ReqData: structs.CreatePaymentDto{
				CardData: structs.Card{CardNumber: "5528790000000001"},
			},
		}

		success := svc.CreatePaymentNonTDS(req)
		assert.False(t, success)
	})
}

func TestPaymentService_CreatePaymentWith3DS(t *testing.T) {
	client, err := iyzipay.New("invalid-key", "invalid-secret")
	require.NoError(t, err)
	svc := NewPaymentService(client)

	t.Run("geçersiz kimlik bilgileri veya ağ hatasında hata dönmeli", func(t *testing.T) {
		req := structs.CreateOrderDTO{
			Order_id:   "ORD-3DS-FAIL",
			TotalPrice: "100",
			ReqData: structs.CreatePaymentDto{
				CardData: structs.Card{CardNumber: "5528790000000001"},
			},
		}

		html, err := svc.CreatePaymentWith3DS(req)
		assert.Error(t, err)
		assert.Empty(t, html)
		assert.Contains(t, err.Error(), "ödeme sırasında hatayla karşılaşıldı")
	})
}
