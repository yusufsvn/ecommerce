package services

import (
	structs "ecommerceapi/Structs"
	"fmt"
	"log"

	"github.com/JspBack/iyzipay-go"
)

type PaymentService struct {
	client *iyzipay.IyzipayClient
}

func NewPaymentService(client *iyzipay.IyzipayClient) *PaymentService {
	return &PaymentService{client: client}
}

func (s *PaymentService) CreatePaymentNonTDS(req structs.CreateOrderDTO) bool {

	ntdsReq := s.InitializePayment(req)

	response, err := s.client.NonTDSPaymentRequest(&ntdsReq)

	if err != nil {
		log.Printf("Error creating payment request: %v", err)
		return false
	}

	if response.Status == "success" {
		log.Printf("Payment: %v", response)
		return true
	}

	log.Println("Payment failed")
	return false
}

func (s *PaymentService) CreatePaymentWith3DS(data structs.CreateOrderDTO) (string, error) {

	ntdsReq := s.InitializePayment(data)

	//conv_id := ntdsReq.ConversationID
	req := &iyzipay.InitTDSRequest{
		PaymentRequest: ntdsReq,
		CallbackUrl:    "https://www.merchant.com/callback",
	}

	res, decodedHtmlContent, err := s.client.InitializeTDSPayment(
		req,
		iyzipay.WithBinRequest(false),       // Bin kontrolü yapmak istemiyorsanız false yapabilirsiniz.
		iyzipay.WithHtmlDecodeRequest(true), // Html içeriğini decode etmek istemiyorsanız false yapabilirsiniz.
	)

	if err != nil {
		log.Printf("Error creating payment request: %v", err)
		return "", fmt.Errorf("ödeme sırasında hatayla karşılaşıldı")
	}

	if res.Status == "success" {
		log.Printf("PaymentId: %v", res.PaymentID)
		log.Printf("DecodedHtmlContent: %v", decodedHtmlContent)
	}

	return decodedHtmlContent, nil
}

func (s *PaymentService) InitializePayment(req structs.CreateOrderDTO) iyzipay.PaymentRequest {

	basketItems := []iyzipay.BasketItem{}

	for i := range req.CartItems {
		quantity := req.CartItems[i].Quantity

		if quantity > 1 {
			for range quantity {
				basketItems = append(basketItems, iyzipay.BasketItem{
					ID:        req.CartItems[i].ProductID,
					Name:      req.CartItems[i].ProductName,
					Category1: req.CartItems[i].Category,
					ItemType:  req.CartItems[i].ItemType,
					Price:     fmt.Sprintf("%d", req.CartItems[i].Price),
				})
			}
		} else {
			basketItems = append(basketItems, iyzipay.BasketItem{
				ID:        req.CartItems[i].ProductID,
				Name:      req.CartItems[i].ProductName,
				Category1: req.CartItems[i].Category,
				ItemType:  req.CartItems[i].ItemType,
				Price:     fmt.Sprintf("%d", req.CartItems[i].Price),
			})
		}
	}

	ntdsReq := &iyzipay.PaymentRequest{
		Locale:         "tr",
		ConversationID: req.Order_id,
		Price:          req.TotalPrice,
		PaidPrice:      req.TotalPrice, // şimdilik ayrı masraf veya indirim yok Price a eşit olacak
		Installment:    1,
		PaymentChannel: "WEB",
		BasketID:       req.Cart_id,
		PaymentGroup:   "PRODUCT",
		PaymentCard: iyzipay.PaymentCard{
			CardHolderName: req.ReqData.CardData.CardName,
			CardNumber:     req.ReqData.CardData.CardNumber,
			ExpireYear:     req.ReqData.CardData.ExpireYear,
			ExpireMonth:    req.ReqData.CardData.ExpireMonth,
			CVC:            req.ReqData.CardData.Cvc,
			RegisterCard:   0,
		},
		// if you have a saved card
		// PaymentCard: iyzipay.PaymentCard{
		// 	CardUserKey: "card user key",
		// 	CardToken:   "card token",
		// },
		Buyer: iyzipay.Buyer{
			ID:                  req.UserData.User_id,
			Name:                req.UserData.Name,
			Surname:             req.UserData.Surname,
			IdentityNumber:      req.ReqData.Tckn,
			Email:               req.UserData.Mail,
			GSMNumber:           req.UserData.Phone,
			RegistrationAddress: req.UserData.Address.AddressStr,
			City:                req.UserData.Address.City,
			Country:             req.UserData.Address.Country,
			ZipCode:             req.UserData.Address.ZipCode,
			IP:                  req.ReqData.User_ip,
		},
		ShippingAddress: iyzipay.ShippingAddress{
			Address:     req.ReqData.ShippingAddress.AddressStr,
			ContactName: req.UserData.Name + " " + req.UserData.Surname,
			City:        req.ReqData.ShippingAddress.City,
			Country:     req.ReqData.ShippingAddress.Country,
			ZipCode:     req.ReqData.ShippingAddress.ZipCode,
		},
		BillingAddress: iyzipay.BillingAddress{
			Address:     req.ReqData.BillingAddress.AddressStr,
			ContactName: req.UserData.Name + " " + req.UserData.Surname,
			City:        req.ReqData.BillingAddress.City,
			Country:     req.ReqData.BillingAddress.Country,
			ZipCode:     req.ReqData.BillingAddress.ZipCode,
		},
		BasketItems: basketItems,
		Currency:    "TRY",
	}

	return *ntdsReq
}
