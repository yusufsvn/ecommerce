package structs

import "mime/multipart"

// HttpRequests/CreateProduct.go
type CreateProductRequest struct {
	ProductName         string                  `json:"ProductName"`
	Stock               int                     `json:"Stock"`
	Description         string                  `json:"Description"`
	Price               int64                   `json:"Price"`
	ProductCategory     []string                `json:"ProductCategory"`
	ItemType            string                  `json:"ItemType"`
	Images              []*multipart.FileHeader `json:"Images"`
	IsChangedThumb      bool                    `json:"IsChangedThub"`
	TechnicalDetails    map[string]string       `json:"TechnicalDetails"`
	ProductCode         string                  `json:"ProductCode"`
	PublishState string                  `json:"PublishState"`
}

type FilterAndSortInfo struct {
	NewestItems bool `json:"newest_items"`
	Cheapest    bool `json:"cheapest"`
	Expensive   bool `json:"expensive"`
	AZ          bool `json:"az"`
	ZA          bool `json:"za"`
}

// HttpRequests/RefreshTokenReq.go
type RefrefshTokenReq struct {
	RefreshToken string
}

// HttpRequests/UpdateRequest.go
type UpdateProductRequest struct {
	ProductID           string            `json:"ProductID"`
	Stock               int               `json:"Stock"`
	Price               int64             `json:"Price"`
	Description         string            `json:"Description"`
	ProductName         string            `json:"ProductName"`
	DeletedImages       []string          `json:"DeletedImages"`
	Categories          UpdateCategories  `json:"Categories"`
	TechnicalDetails    map[string]string `json:"TechnicalDetails"`
	ItemType            string            `json:"ItemType"`
	ProductCode         string            `json:"ProductCode"`
	ProductPublishState string            `json:"ProductPublishState"`
}

type UpdateCategories struct {
	NewCategories     []string `json:"NewCategories"`
	DeletedCategories []string `json:"DeletedCategories"`
}

type PatchProductRequest struct {
	ProductID   string `json:"ProductID"`
	ProductName string `json:"ProductName"`
	Stock       int    `json:"Stock"`
	Price       int64  `json:"Price"`
}

// Structs/UserStructs.go
type LoginReq struct {
	Mail     string `json:"mail"`
	Password string `json:"password"`
}

type UserRegisterReq struct {
	Name           string `json:"name"`
	SecondName     string `json:"second_name"`
	Mail           string `json:"mail"`
	Password       string `json:"password"`
	Password_check string `json:"password_check"`
	Phone_number   string `json:"phone_number"`
	Address        string `json:"address"`
	City           string `json:"city"`
	District       string `json:"district"`
	Country        string `json:"country"`
	ZipCode        string `json:"zip_code"`
}

type AddToCartReq struct {
	ProductID string `json:"product_id"`
	Quantity  int    `json:"quantity"`
}

// Structs/Requests.go
type RemoveFromCartReq struct {
	Product_id string `json:"product_id"`
}

type UpdateCartReq struct {
	Product_id string `json:"product_id"`
	Quantity   int    `json:"quantity"`
}

// Payment Request
type PaymentRequest struct {
	Tckn            string
	CardData        Card
	BillingAddress  Address
	ShippingAddress Address
}

// Payment Request sub struct 1
type Card struct {
	CardName    string // kart sahibinin isim soyismi
	CardNumber  string
	ExpireYear  string
	ExpireMonth string
	Cvc         string
}

// Payment Request sub struct 2
type Address struct {
	AddressName string `json:"AddressName"`
	AddressStr  string `json:"AddressStr"`
	District    string `json:"District"`
	City        string `json:"City"`
	Country     string `json:"Country"`
	ZipCode     string `json:"ZipCode"`
}

type UpdateOrderStateReq struct {
	NewState string `json:"new_state"`
	Order_id string `json:"order_id"`
}

type UpdateMultiOrderReq struct {
	OrderIds []string `json:"OrderIds"`
	NewState string   `json:"New_state"`
}

type LoginGoogle struct {
	ReturnIdToken string `json:"returnIdToken"`
}
