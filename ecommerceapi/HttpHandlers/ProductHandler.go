package httphandlers

import (
	services "ecommerceapi/Services"
	structs "ecommerceapi/Structs"
	"ecommerceapi/Utils"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"strconv"
	"strings"
)

type ProductHandler struct {
	ProductService services.ProductService
}

func NewProductHandler(productService services.ProductService) *ProductHandler {
	return &ProductHandler{
		ProductService: productService,
	}
}

// Get Products with category for admin
func (h *ProductHandler) Getallproducts(w http.ResponseWriter, r *http.Request) {
	// fmt.Println("get all products http handlerına varış sağlandı")

	// Query Parameters
	queryparams := r.URL.Query()

	// limit
	limit, err := strconv.Atoi(queryparams.Get("limit"))
	if err != nil {
		http.Error(w, "tur donusumu sirasinda hatayla karsilasildi", http.StatusInternalServerError)
	}
	if limit > 500 {
		http.Error(w, "tek seferde max 500 urun istegi atabilirsiniz", http.StatusBadRequest)
	}

	// offset
	offset, err := strconv.Atoi(queryparams.Get("offset"))
	if err != nil {
		http.Error(w, "tur donusumu sirasinda hatayla karsilasildi", http.StatusInternalServerError)
	}

	product_list, err := h.ProductService.GetAllProducts(r.Context(), limit, offset)

	if err != nil {
		http.Error(w, "Ürünler getirilemedi: "+err.Error(), http.StatusInternalServerError)
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(product_list)
}

// Get Products for public
func (h *ProductHandler) PubGetAllProducts(w http.ResponseWriter, r *http.Request) {

	// Query Parameters
	queryparams := r.URL.Query()

	// limit
	limit, err := strconv.Atoi(queryparams.Get("limit"))
	if err != nil {
		http.Error(w, "tur donusumu sirasinda hatayla karsilasildi", http.StatusInternalServerError)
	}
	if limit > 500 {
		http.Error(w, "tek seferde max 500 urun istegi atabilirsiniz", http.StatusBadRequest)
	}

	// offset
	offset, err := strconv.Atoi(queryparams.Get("offset"))
	if err != nil {
		http.Error(w, "tur donusumu sirasinda hatayla karsilasildi", http.StatusInternalServerError)
	}

	data, err1 := h.ProductService.PubGetAllProducts(r.Context(), limit, offset)

	if err1 != nil {
		http.Error(w, "beklenmeyen hatayla karşılşıldı", http.StatusInternalServerError)
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(data)

}

func (h *ProductHandler) GetProductById(w http.ResponseWriter, r *http.Request) {
	// get product by id

	// get query parameters
	url := r.URL.Query()
	ProductID := url.Get("ProductID")
	// service call

	response, err := h.ProductService.GetProductById(r.Context(), ProductID)

	if err != nil {
		http.Error(w, "Ürün Getirilimedi", http.StatusInternalServerError)
		return
	}
	// response configuration

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(response)
}

func (h *ProductHandler) GetProductByIdForAdmin(w http.ResponseWriter, r *http.Request) {
	// get product by id

	// get query parameters
	url := r.URL.Query()
	ProductID := url.Get("ProductID")
	// service call

	response, err := h.ProductService.GetProductByIdForAdmin(r.Context(), ProductID)

	if err != nil {
		http.Error(w, "Ürün Getirilimedi", http.StatusInternalServerError)
		return
	}
	// response configuration

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(response)
}

// POST Function
func (h *ProductHandler) CreateProduct(w http.ResponseWriter, r *http.Request) {
	// fmt.Println("Create Product handlerıı")

	err := r.ParseMultipartForm(10 << 20)

	if err != nil {
		http.Error(w, "Form parse edilemedi", http.StatusBadRequest)
		return
	}

	var payload structs.CreateProductRequest

	var techDetailsJSON string = r.FormValue("TechnicalDetails")
	if techDetailsJSON != "" {
		errJSON := json.Unmarshal([]byte(techDetailsJSON), &payload.TechnicalDetails)
		if errJSON != nil {
			http.Error(w, "Teknik detaylar geçerli değil", http.StatusBadRequest)
			return
		}
	}

	payload.ProductName = r.FormValue("ProductName")
	payload.Description = r.FormValue("Description")
	payload.ProductCategory = r.Form["Categories"]

	fmt.Println("debug:::::createproduct")
	fmt.Println(payload.ProductCategory)

	payload.ItemType = r.FormValue("ItemType")
	payload.ProductCode = r.FormValue("ProductCode")
	payload.PublishState = r.FormValue("PublishState")
	boolean, errx := strconv.ParseBool(r.FormValue("IsChangedThumb"))

	if errx != nil {
		http.Error(w, "IsChangedThumb bulunamadı", http.StatusBadRequest)
		return
	}

	payload.IsChangedThumb = boolean

	payload.Price, err = strconv.ParseInt(r.FormValue("Price"), 10, 64)
	if err != nil {
		log.Println("price da hata", err)
		http.Error(w, "stok veya price tip dönüşümünde hata ", http.StatusBadRequest)
	}
	payload.Stock, err = strconv.Atoi(r.FormValue("Stock"))
	if err != nil {
		log.Println("stok da hata", err)
		http.Error(w, "stok tip dönüşümünde hata ", http.StatusBadRequest)
	}
	payload.Images = r.MultipartForm.File["Images"]

	err1 := h.ProductService.CreateProduct(r.Context(), &payload)

	if err1 != nil {
		log.Println("product oluşturmada hata %w", err)
		http.Error(w, "product serviste hata", http.StatusInternalServerError)
	}

	// fmt.Println(payload.ProductName)

}

// Delete Function

func (h *ProductHandler) DeleteProduct(w http.ResponseWriter, r *http.Request) {

	// query params
	urlquery := r.URL.Query()
	ProductID := urlquery.Get("ProductID")

	if ProductID == "" {
		http.Error(w, "ProductID parametresi zorunludur", http.StatusBadRequest)
		return // Hata durumunda akışı kesmek şarttır
	}
	// service call
	err := h.ProductService.DeleteProduct(r.Context(), ProductID)

	if err != nil {
		http.Error(w, "Ürün silinirken hatayla karşılaşıldı", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK) // veya içerik yoksa http.StatusNoContent (204)
	w.Write([]byte(`{"message": "Ürün başarıyla silindi"}`))

}

// Update Functions

func (h *ProductHandler) UpdateProduct(w http.ResponseWriter, r *http.Request) {

	var payload *structs.UpdateProductRequest

	if err := Utils.ReadJSON(w, r, &payload); err != nil {
		Utils.WriteError(w, err, http.StatusBadRequest)
		return
	}

	err := h.ProductService.UpdateProduct(r.Context(), payload)

	if err != nil {
		http.Error(w, "ürün güncellenirken hata çıktı ", http.StatusInternalServerError)
	}
}

func (h *ProductHandler) PatchProduct(w http.ResponseWriter, r *http.Request) {

	var payload *structs.PatchProductRequest

	if err := Utils.ReadJSON(w, r, &payload); err != nil {
		Utils.WriteError(w, err, http.StatusBadRequest)
		return
	}

	err := h.ProductService.PatchProduct(r.Context(), payload)

	if err != nil {
		http.Error(w, "ürün güncellenirken hata çıktı ", http.StatusInternalServerError)
	}
}

func (h *ProductHandler) AddProductImage(w http.ResponseWriter, r *http.Request) {

	err := r.ParseMultipartForm(10 << 20)

	if err != nil {
		http.Error(w, "Form Data Parse Edilemedi", http.StatusBadRequest)
		log.Println("parse hatası %w", err)
	}

	product_id := r.FormValue("product_id")
	isChangedThumbStr := r.FormValue("isChangedThumb")

	log.Println(isChangedThumbStr)
	isChangedThumb, err := strconv.ParseBool(isChangedThumbStr)
	log.Println(isChangedThumb)
	if err != nil {
		http.Error(w, "isChangedThumb değeri true veya false olmalıdır", http.StatusBadRequest)
		return
	}
	images := r.MultipartForm.File["Images"]

	err2 := h.ProductService.AddProductImage(r.Context(), images, product_id, isChangedThumb)

	if err2 != nil {
		http.Error(w, "Resimler Yüklenemedi", http.StatusInternalServerError)
		log.Println("Resim Yükleme Servisinde Hata %w", err)
	}

	//w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
}

func (h *ProductHandler) GetManyProducts(w http.ResponseWriter, r *http.Request) {

	ids := r.URL.Query().Get("ids")
	idslice := strings.Split(ids, ",")

	products, err := h.ProductService.GetManyProducts(r.Context(), idslice)

	if err != nil {
		http.Error(w, "Ürünler getirilemedi", http.StatusInternalServerError)
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(products)
	w.WriteHeader(http.StatusOK)
}

func (h *ProductHandler) GetProductsByCategoryId(w http.ResponseWriter, r *http.Request) {

	limit := r.URL.Query().Get("limit")
	offset := r.URL.Query().Get("offset")
	category_ids := r.URL.Query()["category"]

	if limit == "" || offset == "" || category_ids == nil {
		http.Error(w, "query parametleri boş", http.StatusBadRequest)
		return
	}

	data, err := h.ProductService.GetProductsByCategoryId(r.Context(), category_ids, limit, offset)

	if err != nil {
		http.Error(w, "Beklenmeyen bir hatayla karşılaşıldı", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(data)
}

func (h *ProductHandler) SearchProduct(w http.ResponseWriter, r *http.Request) {

	searchQuery := r.URL.Query().Get("searchQuery")
	limit := r.URL.Query().Get("limit")
	if searchQuery == "" || limit == "" {
		http.Error(w, "invalid search query", http.StatusBadRequest)
		return
	}

	data, err := h.ProductService.SearchProduct(r.Context(), limit, searchQuery)

	if err != nil {
		http.Error(w, "Beklenmeyen hata", http.StatusInternalServerError)
		return
	}

	json.NewEncoder(w).Encode(data)
}

func (h *ProductHandler) SearchProductAdmin(w http.ResponseWriter, r *http.Request) {

	searchQuery := r.URL.Query().Get("searchQuery")
	limit := r.URL.Query().Get("limit")
	if searchQuery == "" || limit == "" {
		http.Error(w, "invalid search query", http.StatusBadRequest)
		return
	}

	data, err := h.ProductService.SearchProductAdmin(r.Context(), limit, searchQuery)

	if err != nil {
		http.Error(w, "Beklenmeyen hata", http.StatusInternalServerError)
		return
	}

	json.NewEncoder(w).Encode(data)
}

// sıralamaya göre ürünler
func (h *ProductHandler) GetFilteredProducts(w http.ResponseWriter, r *http.Request) {

	limit := r.URL.Query().Get("limit")
	offset := r.URL.Query().Get("offset")
	category_ids := r.URL.Query()["category"]
	sortBy := r.URL.Query().Get("sortBy")

	if limit == "" || offset == "" {
		http.Error(w, "query parametleri boş", http.StatusBadRequest)
		return
	}

	log.Println(category_ids)

	data, err := h.ProductService.GetFilteredAndSortedProducts(r.Context(), category_ids, limit, offset, sortBy)

	if err != nil {
		http.Error(w, "beklenmeyen hata", http.StatusInternalServerError)
		return
	}

	json.NewEncoder(w).Encode(data)

}

func (h *ProductHandler) GetProductCount(w http.ResponseWriter, r *http.Request) {

	count, err := h.ProductService.ProductStore.GetProductCount(r.Context())

	if err != nil {
		http.Error(w, "Ürün sayısı Getirilirken hata oluştu", http.StatusInternalServerError)
		return
	}

	json.NewEncoder(w).Encode(count)

}

func (h *ProductHandler) GetProductCountAdmin(w http.ResponseWriter, r *http.Request) {

	count, err := h.ProductService.ProductStore.GetProductCountAdmin(r.Context())

	if err != nil {
		http.Error(w, "Ürün sayısı Getirilirken hata oluştu", http.StatusInternalServerError)
		return
	}

	json.NewEncoder(w).Encode(count)

}

func (h *ProductHandler) UpdatePublishState(w http.ResponseWriter, r *http.Request) {
	type UpdatePublishStateReq struct {
		ProductID string `json:"product_id"`
		NewState  string `json:"new_state"`
	}

	var req UpdatePublishStateReq
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "Geçersiz istek gövdesi", http.StatusBadRequest)
		return
	}

	if req.ProductID == "" || req.NewState == "" {
		http.Error(w, "product_id ve new_state alanları zorunludur", http.StatusBadRequest)
		return
	}

	err := h.ProductService.UpdateProductPublishState(r.Context(), req.ProductID, req.NewState)
	if err != nil {
		http.Error(w, "Yayın durumu güncellenirken hata oluştu: "+err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	w.Write([]byte(`{"message": "Ürün yayın durumu başarıyla güncellendi"}`))
}
