package structs

import (
	"net/url"
	"time"

	"github.com/google/uuid"
)

// HttpRequests/GetAllProductsResponse.go
type ProductResponse struct {
	ProductID uuid.UUID
	//ItemType    string
	ProductName string
	Stock       int
	PhotosURL   string
	Description string
	Price       float64
	Images      []ProductImage
}

type FullProductResponse struct {
	ProductID           string            `json:"product_id"`
	ProductName         string            `json:"product_name"`
	Stock               int               `json:"stock"`
	Description         string            `json:"description"`
	Price               int64             `json:"price"`
	TechnicalDetails    map[string]string `json:"technical_details"`
	Categories          []GetCategory     `json:"categories"`
	Images              []ProductImage    `json:"images"`
	ProductCode         string            `json:"product_code"`
	ProductPublishState string            `json:"product_publish_state"`
	ItemType            string            `json:"item_type"`
}

// Structs/UserStructs.go
type UserLoginRes struct {
	UserID   string `json:"user_id"`
	Name     string `json:"name"`
	Surname  string `json:"surname"`
	Password string `json:"password"`
	Email    string `json:"email"`
	Role     string `json:"role"`
}

type AddToCartRes struct {
	CartID    string `json:"cart_id"`
	ProductID string `json:"product_id"`
	Quantity  int    `json:"quantity"`
}

// Structs/Tokens.go
type TokenQueryRes struct {
	UserId     string `db:"user_id"`
	FamilyId   string `db:"family_id"`
	IsUsed     bool   `db:"is_used"`
	Is_revoked bool   `db:"is_revoked"`
}

// Services/CartService.go
type CartResponse struct {
	CartID   string        `json:"cart_id"`
	Products []CartProduct `json:"products"`
}

// Services/FileManagerService.go
type UrlResult struct {
	CreatedDate  time.Time
	ExpireDate   time.Duration
	PresignedURL *url.URL
}

type ProductResult struct {
	Count int                     `json:"count"`
	Data  []GetAllProductResponse `json:"data"`
}

type ProductResultAdmin struct {
	Count int                          `json:"count"`
	Data  []GetAllProductResponseAdmin `json:"data"`
}

type OrderResult struct {
	Count int             `json:"count"`
	Data  []OrderForAdmin `json:"data"`
}

type GetAllProductResponse struct {
	Product    Getproducts      `json:"Product"`
	Categories []CategoryDbJson `json:"Categories"`
}

type GetAllProductResponseAdmin struct {
	Product    GetproductsAdmin `json:"Product"`
	Categories []CategoryDbJson `json:"Categories"`
}

type CategoryFilteredProducts struct {
	Count int           `json:"Count"`
	Data  []Getproducts `json:"Data"`
}

type GetCategory struct {
	Category_id   string `db:"category_id" json:"category_id"`
	Category_name string `db:"category_name" json:"category_name"`
}
