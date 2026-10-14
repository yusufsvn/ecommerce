package httphandlers

import (
	"encoding/json"
	"errors"
	"log"
	"net/http"

	Customerrors "ecommerceapi/CustomErrors"
	services "ecommerceapi/Services"
	st "ecommerceapi/Structs"
	"ecommerceapi/Utils"
)

type CartHandler struct {
	CartService services.CartService
}

func NewCartHandler(_cartService services.CartService) *CartHandler {
	return &CartHandler{
		CartService: _cartService,
	}
}

func (c *CartHandler) AddToCart(w http.ResponseWriter, r *http.Request) {

	var RequestDTO st.AddToCartReq
	Utils.ReadJSON(w, r, &RequestDTO)

	err := c.CartService.AddToCart(r.Context(), &RequestDTO)

	if err != nil {
		log.Printf("Sepete ürün eklenirken hata oluştu %v", err)
		http.Error(w, "Sepete ürün eklenirken hata oluştu", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)

}

func (c *CartHandler) GetCart(w http.ResponseWriter, r *http.Request) {

	cartProducts, err := c.CartService.GetCart(r.Context())

	if errors.Is(err, Customerrors.ZeroListResult) {
		log.Printf("sepet boş %v", err)
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		json.NewEncoder(w).Encode(cartProducts)
		return
	}

	if err != nil {
		log.Printf("Sepette ürün bulunamadı %v", err)
		http.Error(w, "Sepette ürün bulunamadı", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(cartProducts)

}

func (c *CartHandler) RemoveFromCart(w http.ResponseWriter, r *http.Request) {

	var RequestDTO st.RemoveFromCartReq
	Utils.ReadJSON(w, r, &RequestDTO)

	err := c.CartService.RemoveFromCart(r.Context(), RequestDTO.Product_id)

	if err != nil {
		log.Printf("Sepetten ürün çıkarılırken hata oluştu %v", err)
		http.Error(w, "Sepetten ürün çıkarılırken hata oluştu", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)

}

func (c *CartHandler) UpdateCart(w http.ResponseWriter, r *http.Request) {

	var RequestDTO st.UpdateCartReq
	Utils.ReadJSON(w, r, &RequestDTO)

	err := c.CartService.UpdateQuantity(r.Context(), &RequestDTO)

	if err != nil {
		log.Printf("Sepet güncellenirken hata oluştu %v", err)
		http.Error(w, "Sepet güncellenirken hata oluştu", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)

}
