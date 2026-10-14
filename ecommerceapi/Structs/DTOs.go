package structs

import (
	"time"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5/pgtype"
)

// Structs/productwithimage.go
type ProductWithImage struct {
	ProductID           string         `db:"product_id" json:"product_id"`
	ProductName         string         `db:"product_name" json:"product_name"`
	Stock               int            `db:"stock" json:"stock"`
	ThumbnailImage      string         `db:"thumbnail_image" json:"thumbnail_image"`
	Description         string         `db:"description" json:"description"`
	Price               int64          `db:"price" json:"price"`
	CreatedAt           *time.Time     `db:"created_at" json:"created_at"`
	UpdatedAt           *time.Time     `db:"updated_at" json:"updated_at"`
	Images              []ProductImage `db:"images" json:"images"`
	TechnicalDetails    []byte         `db:"technical_details" json:"technical_details"`
	ProductCode         string         `db:"product_code" json:"product_code"`
	ProductPublishState string         `db:"product_publish_state" json:"product_publish_state"`
	ItemType            string         `db:"item_type" json:"item_type"`
}

type ProductImage struct {
	ImageID uuid.UUID `json:"image_id"`
	Path    string    `json:"path"`
}

// getproducts
type Getproducts struct {
	ProductID      string      `db:"product_id" json:"product_id"`
	ProductName    string      `db:"product_name" json:"product_name"`
	Stock          int         `db:"stock" json:"stock"`
	ThumbnailImage string      `db:"thumbnail_image" json:"thumbnail_image"`
	Price          int64       `db:"price" json:"price"`
	CreatedAt      pgtype.Date `db:"created_at"`
}

// getproducts for admin panel
type GetproductsAdmin struct {
	ProductID           string      `db:"product_id" json:"product_id"`
	ProductName         string      `db:"product_name" json:"product_name"`
	Stock               int         `db:"stock" json:"stock"`
	ThumbnailImage      string      `db:"thumbnail_image" json:"thumbnail_image"`
	Price               int64       `db:"price" json:"price"`
	CreatedAt           pgtype.Date `db:"created_at"`
	ProductCode         string      `db:"product_code" json:"product_code"`
	ProductPublishState string      `db:"product_publish_state" json:"product_publish_state"`
}

type ProductWithCategories struct {
	ProductID      uuid.UUID `db:"product_id" json:"product_id"`
	ProductName    string    `db:"product_name" json:"product_name"`
	Stock          int       `db:"stock" json:"stock"`
	ThumbnailImage string    `db:"thumbnail_image" json:"thumbnail_image"`
	Price          int64     `db:"price" json:"price"`
	Category       []string  `db:"category_name" json:"category_name"`
}

// Structs/Tokens.go
type TokenWithName struct {
	Access_token  string `json:"access_token"`
	Refresh_token string `json:"refresh_token"`
	Name          string `json:"name"`
	Surname       string `json:"surname"`
}

type TokenDTO struct {
	Token      string
	Expires_at time.Time
	Jti        string
}

type RefreshTokenSaveDTO struct {
	Jti        string
	Family_Id  string
	User_id    string
	Is_used    bool
	Is_revoked bool
	Expires_at time.Time
	Token_hash string
}

// Structs/DTOs.go
type CartProductDTO struct {
	Cart_id    string `json:"cart_id" db:"cart_id"`
	Product_id string `json:"product_id" db:"product_id"`
	Quantity   int    `json:"quantity" db:"quantity"`
}

type CartProduct struct {
	Product_id string `json:"product_id"`
	Quantity   int    `json:"quantity"`
}

type CartItems struct {
	Cart_id  string        `json:"cart_id"`
	Products []CartProduct `json:"products"`
}

type ProductForCart struct {
	ProductID      string  `json:"product_id" db:"product_id"`
	ProductName    string  `json:"product_name" db:"product_name"`
	Price          float64 `json:"price" db:"price"`
	Stock          int     `json:"stock" db:"stock"`
	ThumbnailImage string  `json:"thumbnail_image" db:"thumbnail_image"`
}

// Services/productService.go
type GetCartProducts struct {
	ProductID   string       `json:"product_id"`
	ProductName string       `json:"product_name"`
	Price       float64      `json:"price"`
	Stock       int          `json:"stock"`
	Images      ProductImage `json:"product_image"`
}

type UserCartId struct {
	User_id string
	Cart_id string
}

type CreatePaymentDto struct {
	Tckn            string
	CardData        Card
	BillingAddress  Address
	ShippingAddress Address
	User_ip         string
	SecurePay       bool
	SaveAddress     bool
	PhoneNumber     string
}

type ProductsFromCartDb struct {
	ProductID   string `db:"product_id" json:"product_id"`
	ProductName string `db:"product_name" json:"product_name"`
	Price       int64  `db:"price" json:"price"`
	ItemType    string `db:"item_type" json:"item_type"`
	Quantity    int    `db:"quantity" json:"quantity"`
}

type ProductsFromCartDTO struct {
	ProductID   string `db:"product_id" json:"product_id"`
	ProductName string `db:"product_name" json:"product_name"`
	Price       int64  `db:"price" json:"price"`
	ItemType    string `db:"item_type"`
	Quantity    int    `db:"quantity" json:"quantity"`
	Category    string `json:"category"`
}

type ProductsFromCartDTO2 struct {
	Product    ProductsFromCartDTO
	Categories []string
}

type UserWithAddress struct {
	Name       string `db:"name"`
	Surname    string `db:"surname"`
	Email      string `db:"email"`
	Phone      string `db:"phone"`
	Country    string `db:"country"`
	City       string `db:"city"`
	ZipCode    string `db:"post_code"`
	FullAdress string `db:"full_address"`
}

type CreateOrderDTO struct {
	Order_id   string
	TotalPrice string
	Cart_id    string
	CartItems  []ProductsFromCartDTO
	UserData   UserOrderDTO
	ReqData    CreatePaymentDto
}

type CartItemDto struct {
	ProductId string
	Name      string
	Price     string
	Category  string
	ItemType  string
}

type UserOrderDTO struct {
	User_id string
	Name    string
	Surname string
	Mail    string
	Phone   string
	Address Address
}

type OrderDTO struct {
	Order_id        string
	User_id         string
	Status          string
	TotalAmount     string
	ShippingAddress string
	BillingAddress  string
	OrderCode       string
}

type CreateAddressDTO struct {
	UserId      string `db:"user_id" json:"user_id"`
	Name        string `db:"name" json:"name"`
	PostCode    string `db:"post_code" json:"post_code"`
	Country     string `db:"country" json:"country"`
	District    string `db:"district" json:"district"`
	City        string `db:"city" json:"city"`
	FullAddress string `db:"full_address" json:"full_address"`
}

// get from order table -- first order request
type GetOrder struct {
	OrderId       string    `db:"order_code" json:"order_id"`
	UserId        string    `db:"user_id" json:"user_id"`
	TotalPrice    int64     `db:"total_amount" json:"total_amount"`
	Status        string    `db:"status" json:"status"`
	CreatedAt     time.Time `db:"created_at" json:"created_at"`
	TotalQuantity int       `db:"total_quantity" json:"total_quantity"`
	OrderCode     string    `db:"order_code" json:"order_code"`
}

// second order requests with id

// for first db result
type GetProductOrder struct {
	ProductId       string `db:"product_id"`
	Quantity        int    `db:"quantity"`
	PriceAtPurchase int64  `db:"price_at_purchase"`
}

// for seconde db result
type GetProducOrderDetails struct {
	ProductName string `db:"product_name"`
	Image       string `db:"thumbnail_image"`
}

type OrderResponse struct {
	Products        []OrderFullDetail `json:"products"`
	ShippingAddress string            `json:"shipping_address"`
	BillingAddress  string            `json:"billing_address"`
}

type ProductOrderDTO struct {
	ProductID string `db:"product_id" json:"product_id"`
	OrderId   string `db:"order_id" json:"order_id"`
	Price     int64  `db:"price" json:"price"`
	Quantity  int    `db:"quantity" json:"quantity"`
}

type OrderForAdmin struct {
	Name        string    `db:"name" json:"name"`
	Surname     string    `db:"surname" json:"surname"`
	OrderId     string    `db:"order_code" json:"order_id"`
	CreatedAt   time.Time `db:"created_at" json:"created_at"`
	Status      string    `db:"status" json:"status"`
	TotalAmount int64     `db:"total_amount" json:"total_amount"`
}

type OrderFullDetail struct {
	ProductID       string `db:"product_id" json:"product_id"`
	PriceAtPurchase int64  `db:"price_at_purchase" json:"price_at_purchase"`
	Quantity        int    `db:"quantity" json:"quantity"`
	ProductName     string `db:"product_name" json:"product_name"`
	ThumbnailImage  string `db:"thumbnail_image" json:"thumbnail_image"`
}

type GetAllCategory struct {
	ProductId  string           `db:"product_id"`
	Categories []CategoryDbJson `db:"categories_json"`
}

type CategoryDbJson struct {
	CategoryId   string `db:"category_id" json:"category_id"`
	CategoryName string `db:"category_name" json:"category_name"`
}

type GetAddressDTO struct {
	AddressId   string `db:"address_id" json:"AddressId"`
	Name        string `db:"name" json:"AddressName"`
	PostCode    string `db:"post_code" json:"ZipCode"`
	Country     string `db:"country" json:"Country"`
	District    string `db:"district" json:"District"`
	City        string `db:"city" json:"City"`
	FullAddress string `db:"full_address" json:"AddressStr"`
}

type GetAddressList struct {
	Addresses *[]GetAddressDTO `json:"Addresses"`
}

type GetOrderAddress struct {
	OrderId         string `db:"order_id"`
	ShippingAddress string `db:"shipping_address"`
	BillingAddress  string `db:"billing_address"`
}

type GoogleUser struct {
	UserId  string `db:"user_id"`
	Email   string `db:"email"`
	Name    string `db:"name"`
	Surname string `db:"surname"`
}

type UserWithoutAddress struct {
	Name    string `db:"name"`
	Surname string `db:"surname"`
	Email   string `db:"email"`
	Phone   string `db:"phone"`
}
