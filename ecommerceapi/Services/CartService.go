package services

import (
	"context"
	Customerrors "ecommerceapi/CustomErrors"
	db "ecommerceapi/Database/Repository"
	st "ecommerceapi/Structs"
	"ecommerceapi/Utils"
	"errors"
	"fmt"
	"log"
)

type CartService struct {
	CartFuncs    db.CartFunctions
	CartProducts db.CartProductFunctions
}

func NewCartService(_cartFuncs db.CartFunctions, _cartProducts db.CartProductFunctions) *CartService {
	return &CartService{
		CartFuncs:    _cartFuncs,
		CartProducts: _cartProducts,
	}
}

func (c *CartService) AddToCart(ctx context.Context, req *st.AddToCartReq) error {

	// eğer sepet id yoksa db ye sor yoksa oluştur
	ids, err := c.GetUserCartId(ctx)
	if err != nil {
		log.Printf("getUserCartId de hata %v", err)
		return fmt.Errorf("sepet bulunamadı %w", err)
	}

	user_id := ids.User_id
	cart_id := ids.Cart_id

	if cart_id == "" {
		cartId, err2 := c.CartFuncs.Create(ctx, user_id)
		if err2 != nil {
			log.Printf("sepet oluşturulamadı %v", err2)
			return fmt.Errorf("sepet oluşturulamadı %w", err2)
		}
		cart_id = cartId
	}

	// sepete ekle
	cartDTO := st.CartProductDTO{
		Cart_id:    cart_id,
		Product_id: req.ProductID,
		Quantity:   req.Quantity,
	}

	err1 := c.CartProducts.Add(ctx, &cartDTO)
	if err1 != nil {
		log.Printf("could not add to cart %v", err1)
		return fmt.Errorf("sepete eklenemedi")
	}
	return nil
}

func (c *CartService) RemoveFromCart(ctx context.Context, product_id string) error {

	ids, err := c.GetUserCartId(ctx)
	if err != nil {
		log.Printf("getUserCartId de hata %v", err)
		return fmt.Errorf("sepet bulunamadı %w", err)
	}
	log.Printf("idler userid %s cartid: %s", ids.User_id, ids.Cart_id)

	err2 := c.CartProducts.Remove(ctx, ids.Cart_id, product_id)

	if err2 != nil {
		log.Printf("sepetten silinemedi db error %v", err2)
		return fmt.Errorf("sepetten silinemedi from RemoveFromCart")
	}

	return nil
}

func (c *CartService) UpdateQuantity(ctx context.Context, data *st.UpdateCartReq) error {

	ids, err := c.GetUserCartId(ctx)
	if err != nil {
		log.Printf("getUserCartId de hata %v", err)
		return fmt.Errorf("sepet bulunamadı %w", err)
	}

	updateData := st.CartProductDTO{
		Cart_id:    ids.Cart_id,
		Product_id: data.Product_id,
		Quantity:   data.Quantity,
	}
	err2 := c.CartProducts.UpdateQuantity(ctx, &updateData)

	if err2 != nil {
		log.Printf("sepet güncellenemedi db error %v", err2)
		return fmt.Errorf("sepet güncellenemedi")
	}
	return nil
}

func (c *CartService) GetCart(ctx context.Context) (*st.CartItems, error) {

	user_id, err_id := Utils.GetUserIdFromClaims(ctx)

	if err_id != nil || user_id == "" {
		log.Printf("user_id maybe zero length or invalid claims:%v", err_id)
		return nil, fmt.Errorf("Beklenmeyen bir hatayla karşılaşıldı")
	}

	cartProducts, err := c.CartFuncs.GetCartProductsByUserID(ctx, user_id)

	if err != nil {
		if errors.Is(err, Customerrors.ZeroListResult) {
			log.Println("zero list")
			return nil, err
		}
		log.Printf("sepet bulunamadı line 115 CartService--GetCard %v", err)
		return nil, fmt.Errorf("sepet bulunamadı %w", err)
	}

	if len(cartProducts) == 0 {
		return nil, Customerrors.ZeroListResult
	}

	var cartProductsTemp []st.CartProduct = []st.CartProduct{}

	for i := range cartProducts {
		cartProductsTemp = append(cartProductsTemp, st.CartProduct{
			Product_id: cartProducts[i].Product_id,
			Quantity:   cartProducts[i].Quantity,
		})
	}

	cartResponse := st.CartItems{
		Cart_id:  cartProducts[0].Cart_id,
		Products: cartProductsTemp,
	}

	if cartResponse.Cart_id == "" {
		log.Printf("cart_id nil :%s", cartResponse.Cart_id)
		return nil, fmt.Errorf("sepet bulunamadı")
	}

	return &cartResponse, nil
}

// helper function
func (c *CartService) GetUserCartId(ctx context.Context) (*st.UserCartId, error) {

	user_id, err_id := Utils.GetUserIdFromClaims(ctx)

	if err_id != nil || user_id == "" {
		log.Printf("user_id maybe zero length or invalid claims:%v", err_id)
		return nil, fmt.Errorf("Beklenmeyen bir hatayla karşılaşıldı")
	}

	cart, err := c.CartFuncs.GetCartByUserID(ctx, user_id)

	if err != nil {
		log.Printf("sepet bulunamadı %v", err)
		return nil, fmt.Errorf("sepet bulunamadı %w", err)
	}

	ids := st.UserCartId{
		User_id: user_id,
		Cart_id: cart,
	}
	return &ids, nil
}
