package httphandlers

import (
	"bytes"
	"context"
	Entities "ecommerceapi/Database/Models"
	models "ecommerceapi/Database/Models"
	services "ecommerceapi/Services"
	st "ecommerceapi/Structs"
	structs "ecommerceapi/Structs"
	"encoding/json"
	"errors"
	"mime/multipart"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
)

func setupProductHandlerTest() (
	*ProductHandler,
	*mockProductRepo,
	*mockCategoryRepo,
	*mockProductProductCategoryRepo,
	*mockImageRepo,
	*mockFileManager,
) {
	prodRepo := &mockProductRepo{}
	catRepo := &mockCategoryRepo{}
	ppcRepo := &mockProductProductCategoryRepo{}
	imgRepo := &mockImageRepo{}
	fm := &mockFileManager{}

	imageSvc := services.NewImageService(imgRepo, fm)
	prodSvc := services.NewProductService(
		prodRepo,
		catRepo,
		ppcRepo,
		"http://storage.api",
		"http://public.storage",
		"test-bucket",
		*imageSvc,
	)
	handler := NewProductHandler(*prodSvc)
	return handler, prodRepo, catRepo, ppcRepo, imgRepo, fm
}

func TestProductHandler_Getallproducts(t *testing.T) {
	t.Run("geçerli limit ve offset ile admin ürün listesini dönmeli", func(t *testing.T) {
		handler, prodRepo, _, ppcRepo, _, _ := setupProductHandlerTest()

		prodRepo.getProductCountAdminFn = func(ctx context.Context) (int, error) {
			return 1, nil
		}
		prodRepo.getAllAdminFn = func(limit, offset int, ctx context.Context) ([]structs.GetproductsAdmin, error) {
			return []structs.GetproductsAdmin{
				{ProductID: "p-1", ProductName: "Admin Urun"},
			}, nil
		}
		ppcRepo.getAllProductCategoriesFn = func(ctx context.Context, ids []string) ([]*structs.GetAllCategory, error) {
			return []*structs.GetAllCategory{
				{ProductId: "p-1", Categories: []structs.CategoryDbJson{{CategoryName: "Kat"}}},
			}, nil
		}

		req := httptest.NewRequest(http.MethodGet, "/api/admin/getallproducts?limit=10&offset=0", nil)
		rec := httptest.NewRecorder()

		handler.Getallproducts(rec, req)

		assert.Equal(t, http.StatusOK, rec.Code)
		assert.Equal(t, "application/json", rec.Header().Get("Content-Type"))

		var resp structs.ProductResultAdmin
		err := json.Unmarshal(rec.Body.Bytes(), &resp)
		assert.NoError(t, err)
		assert.Equal(t, 1, resp.Count)
	})

	t.Run("limit 500'den büyükse 400 dönmeli", func(t *testing.T) {
		handler, _, _, _, _, _ := setupProductHandlerTest()

		req := httptest.NewRequest(http.MethodGet, "/api/admin/getallproducts?limit=501&offset=0", nil)
		rec := httptest.NewRecorder()

		handler.Getallproducts(rec, req)

		assert.Equal(t, http.StatusBadRequest, rec.Code)
		assert.Contains(t, rec.Body.String(), "tek seferde max 500 urun istegi atabilirsiniz")
	})

	t.Run("limit sayıya dönüştürülemezse 500 dönmeli", func(t *testing.T) {
		handler, _, _, _, _, _ := setupProductHandlerTest()

		req := httptest.NewRequest(http.MethodGet, "/api/admin/getallproducts?limit=abc&offset=0", nil)
		rec := httptest.NewRecorder()

		handler.Getallproducts(rec, req)

		assert.Equal(t, http.StatusInternalServerError, rec.Code)
		assert.Contains(t, rec.Body.String(), "tur donusumu sirasinda hatayla karsilasildi")
	})
}

func TestProductHandler_PubGetAllProducts(t *testing.T) {
	t.Run("herkese açık ürün listesini başarıyla dönmeli", func(t *testing.T) {
		handler, prodRepo, _, _, _, _ := setupProductHandlerTest()

		prodRepo.getProductCountFn = func(ctx context.Context) (int, error) {
			return 1, nil
		}
		prodRepo.getAllFn = func(limit, offset int, ctx context.Context) ([]structs.Getproducts, error) {
			return []structs.Getproducts{
				{ProductID: "p-1", ProductName: "Public Urun"},
			}, nil
		}

		req := httptest.NewRequest(http.MethodGet, "/api/public/getallproducts?limit=10&offset=0", nil)
		rec := httptest.NewRecorder()

		handler.PubGetAllProducts(rec, req)

		assert.Equal(t, http.StatusOK, rec.Code)
		assert.Equal(t, "application/json", rec.Header().Get("Content-Type"))

		var resp structs.ProductResult
		err := json.Unmarshal(rec.Body.Bytes(), &resp)
		assert.NoError(t, err)
		assert.Equal(t, 1, resp.Count)
	})

	t.Run("limit 500'den büyükse 400 dönmeli", func(t *testing.T) {
		handler, _, _, _, _, _ := setupProductHandlerTest()

		req := httptest.NewRequest(http.MethodGet, "/api/public/getallproducts?limit=600&offset=0", nil)
		rec := httptest.NewRecorder()

		handler.PubGetAllProducts(rec, req)

		assert.Equal(t, http.StatusBadRequest, rec.Code)
		assert.Contains(t, rec.Body.String(), "tek seferde max 500 urun istegi atabilirsiniz")
	})
}

func TestProductHandler_GetProductById(t *testing.T) {
	t.Run("geçerli ProductID ile ürün detayını başarıyla dönmeli", func(t *testing.T) {
		handler, prodRepo, catRepo, ppcRepo, imgRepo, _ := setupProductHandlerTest()

		prodID := uuid.New()
		thumbID := uuid.New()
		prodRepo.getProductByIdFn = func(ctx context.Context, ProductID uuid.UUID) (*structs.ProductWithImage, error) {
			return &structs.ProductWithImage{
				ProductID:        prodID.String(),
				ProductName:      "Telefon",
				Price:            5000,
				ThumbnailImage:   thumbID.String(),
				TechnicalDetails: []byte(`{"renk":"siyah"}`),
				Images:           []structs.ProductImage{},
			}, nil
		}
		catRepo.getCategoriesByProductIdFn = func(ctx context.Context, product_id string) ([]structs.GetCategory, error) {
			return []structs.GetCategory{}, nil
		}
		ppcRepo.getByProductIDFn = func(ctx context.Context, productID uuid.UUID) ([]Entities.ProductProductCategory, error) {
			return []Entities.ProductProductCategory{}, nil
		}
		imgRepo.getProductImagesFn = func(ctx context.Context, product_id string) ([]string, error) {
			return []string{}, nil
		}

		req := httptest.NewRequest(http.MethodGet, "/api/public/getproductbyid?ProductID="+prodID.String(), nil)
		rec := httptest.NewRecorder()

		handler.GetProductById(rec, req)

		assert.Equal(t, http.StatusOK, rec.Code)
		assert.Equal(t, "application/json", rec.Header().Get("Content-Type"))

		var resp structs.FullProductResponse
		err := json.Unmarshal(rec.Body.Bytes(), &resp)
		assert.NoError(t, err)
		assert.Equal(t, "Telefon", resp.ProductName)
	})

	t.Run("servis hata verdiğinde 500 dönmeli", func(t *testing.T) {
		handler, _, _, _, _, _ := setupProductHandlerTest()

		req := httptest.NewRequest(http.MethodGet, "/api/public/getproductbyid?ProductID=not-a-uuid", nil)
		rec := httptest.NewRecorder()

		handler.GetProductById(rec, req)

		assert.Equal(t, http.StatusInternalServerError, rec.Code)
		assert.Contains(t, rec.Body.String(), "Ürün Getirilimedi")
	})
}

func TestProductHandler_CreateProduct(t *testing.T) {
	t.Run("başarıyla multipart form ile ürün oluşturmalı", func(t *testing.T) {
		handler, prodRepo, catRepo, ppcRepo, _, _ := setupProductHandlerTest()

		prodRepo.createFn = func(ctx context.Context, product *models.Product) error {
			assert.Equal(t, "Yeni Telefon", product.ProductName)
			return nil
		}
		catRepo.getByCategoryNameFn = func(ctx context.Context, categoryName string) (uuid.UUID, error) {
			return uuid.New(), nil
		}
		ppcRepo.addFn = func(ctx context.Context, item *Entities.ProductProductCategory) error {
			return nil
		}

		body := &bytes.Buffer{}
		writer := multipart.NewWriter(body)
		writer.WriteField("ProductName", "Yeni Telefon")
		writer.WriteField("Description", "Aciklama")
		writer.WriteField("ItemType", "PHYSICAL")
		writer.WriteField("ProductCode", "PRD1")
		writer.WriteField("ProductPublishState", "PUBLISHED")
		writer.WriteField("IsChangedThumb", "false")
		writer.WriteField("Price", "1500")
		writer.WriteField("Stock", "10")
		writer.WriteField("TechnicalDetails", `{"renk":"siyah"}`)
		writer.WriteField("Categories", "Elektronik")
		writer.Close()

		req := httptest.NewRequest(http.MethodPost, "/api/admin/createproduct", body)
		req.Header.Set("Content-Type", writer.FormDataContentType())
		rec := httptest.NewRecorder()

		handler.CreateProduct(rec, req)

		assert.Equal(t, http.StatusOK, rec.Code)
	})

	t.Run("geçersiz JSON technical details verildiğinde 400 dönmeli", func(t *testing.T) {
		handler, _, _, _, _, _ := setupProductHandlerTest()

		body := &bytes.Buffer{}
		writer := multipart.NewWriter(body)
		writer.WriteField("ProductName", "Yeni Telefon")
		writer.WriteField("TechnicalDetails", `invalid-json`)
		writer.Close()

		req := httptest.NewRequest(http.MethodPost, "/api/admin/createproduct", body)
		req.Header.Set("Content-Type", writer.FormDataContentType())
		rec := httptest.NewRecorder()

		handler.CreateProduct(rec, req)

		assert.Equal(t, http.StatusBadRequest, rec.Code)
		assert.Contains(t, rec.Body.String(), "Teknik detaylar geçerli değil")
	})

	t.Run("IsChangedThumb eksik veya geçersizse 400 dönmeli", func(t *testing.T) {
		handler, _, _, _, _, _ := setupProductHandlerTest()

		body := &bytes.Buffer{}
		writer := multipart.NewWriter(body)
		writer.WriteField("ProductName", "Yeni Telefon")
		writer.WriteField("IsChangedThumb", "not-a-bool")
		writer.Close()

		req := httptest.NewRequest(http.MethodPost, "/api/admin/createproduct", body)
		req.Header.Set("Content-Type", writer.FormDataContentType())
		rec := httptest.NewRecorder()

		handler.CreateProduct(rec, req)

		assert.Equal(t, http.StatusBadRequest, rec.Code)
		assert.Contains(t, rec.Body.String(), "IsChangedThumb bulunamadı")
	})
}

func TestProductHandler_DeleteProduct(t *testing.T) {
	t.Run("geçerli ProductID ile ürünü silmeli", func(t *testing.T) {
		handler, prodRepo, _, ppcRepo, imgRepo, fm := setupProductHandlerTest()

		prodID := uuid.New()
		imgRepo.getProductImagesFn = func(ctx context.Context, product_id string) ([]string, error) {
			return []string{}, nil
		}
		fm.deleteMultipleFileFn = func(ctx context.Context, images []string) error { return nil }
		imgRepo.deleteImageFn = func(ctx context.Context, images []string) error { return nil }
		ppcRepo.getByProductIDOneRowFn = func(ctx context.Context, productID uuid.UUID) (uuid.UUID, error) {
			return uuid.New(), nil
		}
		ppcRepo.removeFn = func(ctx context.Context, productID, categoryID uuid.UUID) error { return nil }
		prodRepo.deleteFn = func(ctx context.Context, productId *uuid.UUID) error {
			assert.Equal(t, prodID, *productId)
			return nil
		}

		req := httptest.NewRequest(http.MethodDelete, "/api/admin/deleteproduct?ProductID="+prodID.String(), nil)
		rec := httptest.NewRecorder()

		handler.DeleteProduct(rec, req)

		assert.Equal(t, http.StatusOK, rec.Code)
		assert.Contains(t, rec.Body.String(), "Ürün başarıyla silindi")
	})

	t.Run("ProductID parametresi eksikse 400 dönmeli", func(t *testing.T) {
		handler, _, _, _, _, _ := setupProductHandlerTest()

		req := httptest.NewRequest(http.MethodDelete, "/api/admin/deleteproduct", nil)
		rec := httptest.NewRecorder()

		handler.DeleteProduct(rec, req)

		assert.Equal(t, http.StatusBadRequest, rec.Code)
		assert.Contains(t, rec.Body.String(), "ProductID parametresi zorunludur")
	})

	t.Run("silme sırasında hata çıkarsa 500 dönmeli", func(t *testing.T) {
		handler, _, _, _, _, _ := setupProductHandlerTest()

		req := httptest.NewRequest(http.MethodDelete, "/api/admin/deleteproduct?ProductID=not-a-uuid", nil)
		rec := httptest.NewRecorder()

		handler.DeleteProduct(rec, req)

		assert.Equal(t, http.StatusInternalServerError, rec.Code)
		assert.Contains(t, rec.Body.String(), "Ürün silinirken hatayla karşılaşıldı")
	})
}

func TestProductHandler_UpdateProduct(t *testing.T) {
	t.Run("başarıyla ürünü güncellemeli", func(t *testing.T) {
		handler, prodRepo, _, _, _, _ := setupProductHandlerTest()

		prodRepo.updateFn = func(ctx context.Context, product *models.Product) error {
			assert.Equal(t, "Guncel Saat", product.ProductName)
			return nil
		}

		body := map[string]interface{}{
			"ProductID":   "prod-1",
			"ProductName": "Guncel Saat",
			"Price":       250,
			"Stock":       5,
		}
		bodyBytes, _ := json.Marshal(body)

		req := httptest.NewRequest(http.MethodPut, "/api/admin/updateproduct", bytes.NewReader(bodyBytes))
		rec := httptest.NewRecorder()

		handler.UpdateProduct(rec, req)

		assert.Equal(t, http.StatusOK, rec.Code)
	})

	t.Run("geçersiz JSON verildiğinde 400 dönmeli", func(t *testing.T) {
		handler, _, _, _, _, _ := setupProductHandlerTest()

		req := httptest.NewRequest(http.MethodPut, "/api/admin/updateproduct", bytes.NewReader([]byte("invalid-json")))
		rec := httptest.NewRecorder()

		handler.UpdateProduct(rec, req)

		assert.Equal(t, http.StatusBadRequest, rec.Code)
	})
}

func TestProductHandler_PatchProduct(t *testing.T) {
	t.Run("başarıyla patch uygulamalı", func(t *testing.T) {
		handler, prodRepo, _, _, _, _ := setupProductHandlerTest()

		prodRepo.patchFn = func(ctx context.Context, product *models.Product) error {
			assert.Equal(t, "Patch Saat", product.ProductName)
			return nil
		}

		body := map[string]interface{}{
			"ProductID":   "prod-1",
			"ProductName": "Patch Saat",
		}
		bodyBytes, _ := json.Marshal(body)

		req := httptest.NewRequest(http.MethodPatch, "/api/admin/patchproduct", bytes.NewReader(bodyBytes))
		rec := httptest.NewRecorder()

		handler.PatchProduct(rec, req)

		assert.Equal(t, http.StatusOK, rec.Code)
	})

	t.Run("geçersiz JSON verildiğinde 400 dönmeli", func(t *testing.T) {
		handler, _, _, _, _, _ := setupProductHandlerTest()

		req := httptest.NewRequest(http.MethodPatch, "/api/admin/patchproduct", bytes.NewReader([]byte("invalid-json")))
		rec := httptest.NewRecorder()

		handler.PatchProduct(rec, req)

		assert.Equal(t, http.StatusBadRequest, rec.Code)
	})
}

func TestProductHandler_AddProductImage(t *testing.T) {
	t.Run("geçersiz isChangedThumb verildiğinde 400 dönmeli", func(t *testing.T) {
		handler, _, _, _, _, _ := setupProductHandlerTest()

		body := &bytes.Buffer{}
		writer := multipart.NewWriter(body)
		writer.WriteField("product_id", "prod-1")
		writer.WriteField("isChangedThumb", "not-a-bool")
		writer.Close()

		req := httptest.NewRequest(http.MethodPost, "/api/admin/addimage", body)
		req.Header.Set("Content-Type", writer.FormDataContentType())
		rec := httptest.NewRecorder()

		handler.AddProductImage(rec, req)

		assert.Equal(t, http.StatusBadRequest, rec.Code)
		assert.Contains(t, rec.Body.String(), "isChangedThumb değeri true veya false olmalıdır")
	})

	t.Run("başarıyla görsel yüklemeli", func(t *testing.T) {
		handler, _, _, _, imgRepo, _ := setupProductHandlerTest()

		imgRepo.addOneImageFn = func(ctx context.Context, image models.Image) error {
			return nil
		}

		body := &bytes.Buffer{}
		writer := multipart.NewWriter(body)
		writer.WriteField("product_id", "prod-1")
		writer.WriteField("isChangedThumb", "false")
		writer.Close()

		req := httptest.NewRequest(http.MethodPost, "/api/admin/addimage", body)
		req.Header.Set("Content-Type", writer.FormDataContentType())
		rec := httptest.NewRecorder()

		handler.AddProductImage(rec, req)

		assert.Equal(t, http.StatusOK, rec.Code)
	})
}

func TestProductHandler_GetManyProducts(t *testing.T) {
	t.Run("virgülle ayrılmış ID'ler ile ürünleri başarıyla dönmeli", func(t *testing.T) {
		handler, prodRepo, _, _, _, _ := setupProductHandlerTest()

		prodRepo.getManyFn = func(ctx context.Context, ids []string) ([]*structs.ProductForCart, error) {
			assert.Equal(t, []string{"p1", "p2"}, ids)
			return []*structs.ProductForCart{
				{ProductID: "p1", ProductName: "Urun 1"},
				{ProductID: "p2", ProductName: "Urun 2"},
			}, nil
		}

		req := httptest.NewRequest(http.MethodGet, "/api/public/getmanyproducts?ids=p1,p2", nil)
		rec := httptest.NewRecorder()

		handler.GetManyProducts(rec, req)

		assert.Equal(t, http.StatusOK, rec.Code)
		assert.Equal(t, "application/json", rec.Header().Get("Content-Type"))

		var resp []*structs.ProductForCart
		err := json.Unmarshal(rec.Body.Bytes(), &resp)
		assert.NoError(t, err)
		assert.Len(t, resp, 2)
	})
}

func TestProductHandler_GetProductsByCategoryId(t *testing.T) {
	t.Run("başarıyla kategori ürünlerini dönmeli", func(t *testing.T) {
		handler, prodRepo, _, _, _, _ := setupProductHandlerTest()

		prodRepo.getCountProductsByCategoryFn = func(ctx context.Context, category_ids []string) (int, error) {
			return 1, nil
		}
		prodRepo.getProductsByCategoryIdFn = func(ctx context.Context, category_ids []string, limit, offset string) ([]structs.Getproducts, error) {
			return []structs.Getproducts{
				{ProductID: "p-1", ProductName: "Kategori Urunu"},
			}, nil
		}

		req := httptest.NewRequest(http.MethodGet, "/api/public/getproductsbycategory?category=cat-1&limit=10&offset=0", nil)
		rec := httptest.NewRecorder()

		handler.GetProductsByCategoryId(rec, req)

		assert.Equal(t, http.StatusOK, rec.Code)
		assert.Equal(t, "application/json", rec.Header().Get("Content-Type"))

		var resp structs.CategoryFilteredProducts
		err := json.Unmarshal(rec.Body.Bytes(), &resp)
		assert.NoError(t, err)
		assert.Equal(t, 1, resp.Count)
	})

	t.Run("parametreler eksikse 400 dönmeli", func(t *testing.T) {
		handler, _, _, _, _, _ := setupProductHandlerTest()

		req := httptest.NewRequest(http.MethodGet, "/api/public/getproductsbycategory?category=cat-1", nil)
		rec := httptest.NewRecorder()

		handler.GetProductsByCategoryId(rec, req)

		assert.Equal(t, http.StatusBadRequest, rec.Code)
		assert.Contains(t, rec.Body.String(), "query parametleri boş")
	})
}

func TestProductHandler_SearchProduct(t *testing.T) {
	t.Run("başarıyla arama sonuçlarını dönmeli", func(t *testing.T) {
		handler, prodRepo, _, _, _, _ := setupProductHandlerTest()

		prodRepo.searchProductCountFn = func(ctx context.Context, searchQuery string) int {
			return 1
		}
		prodRepo.searchProductFn = func(ctx context.Context, searchQuery, limit string) ([]st.Getproducts, error) {
			return []st.Getproducts{
				{ProductID: "p-1", ProductName: "Laptop"},
			}, nil
		}

		req := httptest.NewRequest(http.MethodGet, "/api/public/search?searchQuery=Laptop&limit=10", nil)
		rec := httptest.NewRecorder()

		handler.SearchProduct(rec, req)

		assert.Equal(t, http.StatusOK, rec.Code)
		var resp []st.Getproducts
		err := json.Unmarshal(rec.Body.Bytes(), &resp)
		assert.NoError(t, err)
		assert.Len(t, resp, 1)
	})

	t.Run("arama parametreleri eksikse 400 dönmeli", func(t *testing.T) {
		handler, _, _, _, _, _ := setupProductHandlerTest()

		req := httptest.NewRequest(http.MethodGet, "/api/public/search", nil)
		rec := httptest.NewRecorder()

		handler.SearchProduct(rec, req)

		assert.Equal(t, http.StatusBadRequest, rec.Code)
		assert.Contains(t, rec.Body.String(), "invalid search query")
	})
}

func TestProductHandler_SearchProductAdmin(t *testing.T) {
	t.Run("başarıyla admin arama sonuçlarını dönmeli", func(t *testing.T) {
		handler, prodRepo, _, ppcRepo, _, _ := setupProductHandlerTest()

		prodRepo.searchProductCountAdminFn = func(ctx context.Context, searchQuery string) int {
			return 1
		}
		prodRepo.searchProductAdminFn = func(ctx context.Context, searchQuery, limit string) ([]st.GetproductsAdmin, error) {
			return []st.GetproductsAdmin{
				{ProductID: "p-1", ProductName: "Admin Laptop"},
			}, nil
		}
		ppcRepo.getAllProductCategoriesFn = func(ctx context.Context, ids []string) ([]*structs.GetAllCategory, error) {
			return []*structs.GetAllCategory{
				{ProductId: "p-1", Categories: []structs.CategoryDbJson{{CategoryName: "Bilgisayar"}}},
			}, nil
		}

		req := httptest.NewRequest(http.MethodGet, "/api/admin/search?searchQuery=Laptop&limit=10", nil)
		rec := httptest.NewRecorder()

		handler.SearchProductAdmin(rec, req)

		assert.Equal(t, http.StatusOK, rec.Code)
		var resp []structs.GetAllProductResponseAdmin
		err := json.Unmarshal(rec.Body.Bytes(), &resp)
		assert.NoError(t, err)
		assert.Len(t, resp, 1)
	})

	t.Run("arama parametreleri eksikse 400 dönmeli", func(t *testing.T) {
		handler, _, _, _, _, _ := setupProductHandlerTest()

		req := httptest.NewRequest(http.MethodGet, "/api/admin/search", nil)
		rec := httptest.NewRecorder()

		handler.SearchProductAdmin(rec, req)

		assert.Equal(t, http.StatusBadRequest, rec.Code)
		assert.Contains(t, rec.Body.String(), "invalid search query")
	})
}

func TestProductHandler_GetFilteredProducts(t *testing.T) {
	t.Run("başarıyla filtrelenmiş ürünleri dönmeli", func(t *testing.T) {
		handler, prodRepo, _, _, _, _ := setupProductHandlerTest()

		prodRepo.getCountProductsByCategoryFn = func(ctx context.Context, category_ids []string) (int, error) {
			return 1, nil
		}
		prodRepo.filteredAndSortedProductsFn = func(ctx context.Context, category_ids []string, limit, offset, sortBy, sortOrder string) ([]structs.Getproducts, error) {
			return []structs.Getproducts{
				{ProductID: "p-1", ProductName: "Filtrelenmis Urun"},
			}, nil
		}

		req := httptest.NewRequest(http.MethodGet, "/api/public/getfilteredproducts?limit=10&offset=0&sortBy=cheapest", nil)
		rec := httptest.NewRecorder()

		handler.GetFilteredProducts(rec, req)

		assert.Equal(t, http.StatusOK, rec.Code)
		var resp structs.CategoryFilteredProducts
		err := json.Unmarshal(rec.Body.Bytes(), &resp)
		assert.NoError(t, err)
		assert.Equal(t, 1, resp.Count)
	})

	t.Run("limit veya offset eksikse 400 dönmeli", func(t *testing.T) {
		handler, _, _, _, _, _ := setupProductHandlerTest()

		req := httptest.NewRequest(http.MethodGet, "/api/public/getfilteredproducts", nil)
		rec := httptest.NewRecorder()

		handler.GetFilteredProducts(rec, req)

		assert.Equal(t, http.StatusBadRequest, rec.Code)
		assert.Contains(t, rec.Body.String(), "query parametleri boş")
	})
}

func TestProductHandler_GetProductCount(t *testing.T) {
	t.Run("başarıyla genel ürün sayısını dönmeli", func(t *testing.T) {
		handler, prodRepo, _, _, _, _ := setupProductHandlerTest()

		prodRepo.getProductCountFn = func(ctx context.Context) (int, error) {
			return 100, nil
		}

		req := httptest.NewRequest(http.MethodGet, "/api/public/productcount", nil)
		rec := httptest.NewRecorder()

		handler.GetProductCount(rec, req)

		assert.Equal(t, http.StatusOK, rec.Code)
		var count int
		err := json.Unmarshal(rec.Body.Bytes(), &count)
		assert.NoError(t, err)
		assert.Equal(t, 100, count)
	})

	t.Run("sayı sorgusu hata verirse 500 dönmeli", func(t *testing.T) {
		handler, prodRepo, _, _, _, _ := setupProductHandlerTest()

		prodRepo.getProductCountFn = func(ctx context.Context) (int, error) {
			return 0, errors.New("db count failed")
		}

		req := httptest.NewRequest(http.MethodGet, "/api/public/productcount", nil)
		rec := httptest.NewRecorder()

		handler.GetProductCount(rec, req)

		assert.Equal(t, http.StatusInternalServerError, rec.Code)
		assert.Contains(t, rec.Body.String(), "Ürün sayısı Getirilirken hata oluştu")
	})
}

func TestProductHandler_GetProductCountAdmin(t *testing.T) {
	t.Run("başarıyla admin ürün sayısını dönmeli", func(t *testing.T) {
		handler, prodRepo, _, _, _, _ := setupProductHandlerTest()

		prodRepo.getProductCountAdminFn = func(ctx context.Context) (int, error) {
			return 150, nil
		}

		req := httptest.NewRequest(http.MethodGet, "/api/admin/productcount", nil)
		rec := httptest.NewRecorder()

		handler.GetProductCountAdmin(rec, req)

		assert.Equal(t, http.StatusOK, rec.Code)
		var count int
		err := json.Unmarshal(rec.Body.Bytes(), &count)
		assert.NoError(t, err)
		assert.Equal(t, 150, count)
	})

	t.Run("sayı sorgusu hata verirse 500 dönmeli", func(t *testing.T) {
		handler, prodRepo, _, _, _, _ := setupProductHandlerTest()

		prodRepo.getProductCountAdminFn = func(ctx context.Context) (int, error) {
			return 0, errors.New("admin count failed")
		}

		req := httptest.NewRequest(http.MethodGet, "/api/admin/productcount", nil)
		rec := httptest.NewRecorder()

		handler.GetProductCountAdmin(rec, req)

		assert.Equal(t, http.StatusInternalServerError, rec.Code)
		assert.Contains(t, rec.Body.String(), "Ürün sayısı Getirilirken hata oluştu")
	})
}

func TestProductHandler_UpdatePublishState(t *testing.T) {
	t.Run("başarıyla yayın durumunu güncellemeli", func(t *testing.T) {
		handler, prodRepo, _, _, _, _ := setupProductHandlerTest()

		prodRepo.changePublishStateFn = func(ctx context.Context, product_id, newState string) error {
			assert.Equal(t, "prod-1", product_id)
			assert.Equal(t, "published", newState)
			return nil
		}

		body := map[string]string{
			"product_id": "prod-1",
			"new_state":  "published",
		}
		bodyBytes, _ := json.Marshal(body)

		req := httptest.NewRequest(http.MethodPut, "/api/admin/product/publishstate", bytes.NewReader(bodyBytes))
		rec := httptest.NewRecorder()

		handler.UpdatePublishState(rec, req)

		assert.Equal(t, http.StatusOK, rec.Code)
		assert.Equal(t, "application/json", rec.Header().Get("Content-Type"))
		assert.Contains(t, rec.Body.String(), "Ürün yayın durumu başarıyla güncellendi")
	})

	t.Run("alanlar eksikse 400 dönmeli", func(t *testing.T) {
		handler, _, _, _, _, _ := setupProductHandlerTest()

		body := map[string]string{
			"product_id": "prod-1",
		}
		bodyBytes, _ := json.Marshal(body)

		req := httptest.NewRequest(http.MethodPut, "/api/admin/product/publishstate", bytes.NewReader(bodyBytes))
		rec := httptest.NewRecorder()

		handler.UpdatePublishState(rec, req)

		assert.Equal(t, http.StatusBadRequest, rec.Code)
		assert.Contains(t, rec.Body.String(), "product_id ve new_state alanları zorunludur")
	})

	t.Run("geçersiz JSON gövdesinde 400 dönmeli", func(t *testing.T) {
		handler, _, _, _, _, _ := setupProductHandlerTest()

		req := httptest.NewRequest(http.MethodPut, "/api/admin/product/publishstate", bytes.NewReader([]byte("invalid-json")))
		rec := httptest.NewRecorder()

		handler.UpdatePublishState(rec, req)

		assert.Equal(t, http.StatusBadRequest, rec.Code)
		assert.Contains(t, rec.Body.String(), "Geçersiz istek gövdesi")
	})

	t.Run("geçersiz durum değeri verildiğinde 500 dönmeli", func(t *testing.T) {
		handler, _, _, _, _, _ := setupProductHandlerTest()

		body := map[string]string{
			"product_id": "prod-1",
			"new_state":  "UNKNOWN_STATE",
		}
		bodyBytes, _ := json.Marshal(body)

		req := httptest.NewRequest(http.MethodPut, "/api/admin/product/publishstate", bytes.NewReader(bodyBytes))
		rec := httptest.NewRecorder()

		handler.UpdatePublishState(rec, req)

		assert.Equal(t, http.StatusInternalServerError, rec.Code)
		assert.Contains(t, rec.Body.String(), "Yayın durumu güncellenirken hata oluştu")
	})
}
