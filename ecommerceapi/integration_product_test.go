package main

import (
	"bytes"
	"context"
	models "ecommerceapi/Database/Models"
	structs "ecommerceapi/Structs"
	"encoding/json"
	"fmt"
	"mime/multipart"
	"net/http"
	"net/url"
	"testing"
	"time"

	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
)

func TestIntegration_Category_Lifecycle(t *testing.T) {
	app := setupIntegrationApp(t)

	_, adminCookie, _ := app.createTestUser(t, fmt.Sprintf("admin-cat-%s@example.com", uuid.NewString()[:8]), "Pass123!", "admin")

	var categoryID string
	catName := fmt.Sprintf("Kategori-%s", uuid.NewString()[:6])

	t.Run("admin yeni kategori ekleyebilmeli", func(t *testing.T) {
		path := fmt.Sprintf("/api/admin/createcategory?category_name=%s", url.QueryEscape(catName))
		resp := app.doRequest(t, http.MethodPost, path, nil, adminCookie)
		assert.Equal(t, http.StatusOK, resp.StatusCode)

		var createdCat structs.GetCategory
		err := json.NewDecoder(resp.Body).Decode(&createdCat)
		assert.NoError(t, err)
		assert.NotEmpty(t, createdCat.Category_id)
		assert.Equal(t, catName, createdCat.Category_name)
		categoryID = createdCat.Category_id
	})

	t.Run("public kullanıcı tüm kategorileri listeleyebilmeli", func(t *testing.T) {
		resp := app.doRequest(t, http.MethodGet, "/api/public/getcategories", nil)
		assert.Equal(t, http.StatusOK, resp.StatusCode)

		var cats []structs.GetCategory
		err := json.NewDecoder(resp.Body).Decode(&cats)
		assert.NoError(t, err)

		found := false
		for _, c := range cats {
			if c.Category_id == categoryID {
				found = true
				break
			}
		}
		assert.True(t, found, "Eklenen kategori listede bulunmalı")
	})

	t.Run("admin kategoriyi güncelleyebilmeli", func(t *testing.T) {
		updatedName := catName + "-Guncel"
		path := fmt.Sprintf("/api/admin/updatecategory?category_id=%s&category_name=%s", categoryID, url.QueryEscape(updatedName))
		resp := app.doRequest(t, http.MethodPut, path, nil, adminCookie)
		assert.Equal(t, http.StatusOK, resp.StatusCode)

		// DB'den kontrol et
		ctx := context.Background()
		cats, err := app.CategoryStore.GetAll(ctx)
		assert.NoError(t, err)
		var match bool
		for _, c := range cats {
			if c.Category_id == categoryID && c.Category_name == updatedName {
				match = true
				break
			}
		}
		assert.True(t, match, "Kategori adı güncellenmiş olmalı")
	})

	t.Run("admin kategoriyi silebilmeli", func(t *testing.T) {
		path := fmt.Sprintf("/api/admin/deletecategory?category_id=%s", categoryID)
		resp := app.doRequest(t, http.MethodDelete, path, nil, adminCookie)
		assert.Equal(t, http.StatusOK, resp.StatusCode)
	})
}

func TestIntegration_Product_Lifecycle(t *testing.T) {
	app := setupIntegrationApp(t)

	_, adminCookie, _ := app.createTestUser(t, fmt.Sprintf("admin-prod-%s@example.com", uuid.NewString()[:8]), "Pass123!", "admin")

	// Önce bir kategori oluşturalım
	ctx := context.Background()
	cat, err := app.CategoryStore.Create(ctx, fmt.Sprintf("Cat-%s", uuid.NewString()[:6]))
	assert.NoError(t, err)

	var productID string
	prodName := fmt.Sprintf("Test Telefon %s", uuid.NewString()[:6])

	t.Run("admin multipart form ile ürün oluşturabilmeli", func(t *testing.T) {
		var body bytes.Buffer
		writer := multipart.NewWriter(&body)

		_ = writer.WriteField("ProductName", prodName)
		_ = writer.WriteField("Description", "Kapsamlı entegrasyon test telefonu")
		_ = writer.WriteField("Price", "15000")
		_ = writer.WriteField("Stock", "25")
		_ = writer.WriteField("ItemType", "PHYSICAL")
		_ = writer.WriteField("ProductCode", fmt.Sprintf("CODE%s", uuid.NewString()[:4]))
		_ = writer.WriteField("ProductPublishState", "published")
		_ = writer.WriteField("IsChangedThumb", "false")
		_ = writer.WriteField("TechnicalDetails", `{"renk":"uzay grisi","hafiza":"256gb"}`)
		_ = writer.WriteField("Categories", cat.Category_name)

		writer.Close()

		req, err := http.NewRequest(http.MethodPost, app.Server.URL+"/api/admin/createproduct", &body)
		assert.NoError(t, err)
		req.Header.Set("Content-Type", writer.FormDataContentType())
		req.AddCookie(adminCookie)

		client := &http.Client{}
		resp, err := client.Do(req)
		assert.NoError(t, err)
		assert.Equal(t, http.StatusOK, resp.StatusCode)

		// DB'den ürünü bulalım
		prods, err := app.ProductStore.GetAllAdmin(10, 0, ctx)
		assert.NoError(t, err)
		for _, p := range prods {
			if p.ProductName == prodName {
				productID = p.ProductID
				break
			}
		}
		assert.NotEmpty(t, productID, "Ürün DB'ye kaydedilmiş olmalı")
	})

	t.Run("public kullanıcı getproductbyid ile ürün detaylarını alabilmeli", func(t *testing.T) {
		resp := app.doRequest(t, http.MethodGet, fmt.Sprintf("/api/public/getproductbyid?ProductID=%s", productID), nil)
		assert.Equal(t, http.StatusOK, resp.StatusCode)

		var fullProd structs.FullProductResponse
		err := json.NewDecoder(resp.Body).Decode(&fullProd)
		assert.NoError(t, err)
		assert.Equal(t, prodName, fullProd.ProductName)
		assert.Equal(t, int64(15000), fullProd.Price)
		assert.Equal(t, 25, fullProd.Stock)
		assert.Equal(t, "uzay grisi", fullProd.TechnicalDetails["renk"])
	})

	t.Run("admin patchproduct ile stok ve fiyat güncelleyebilmeli", func(t *testing.T) {
		patchPayload := structs.PatchProductRequest{
			ProductID:   productID,
			ProductName: prodName,
			Stock:       50,
			Price:       17500,
		}

		resp := app.doRequest(t, http.MethodPatch, "/api/admin/patchproduct", patchPayload, adminCookie)
		assert.Equal(t, http.StatusOK, resp.StatusCode)

		// Kontrol et
		parsedUUID, _ := uuid.Parse(productID)
		prod, err := app.ProductStore.GetProductById(ctx, parsedUUID)
		assert.NoError(t, err)
		assert.Equal(t, 50, prod.Stock)
		assert.Equal(t, int64(17500), prod.Price)
	})

	t.Run("admin ürünü silebilmeli", func(t *testing.T) {
		resp := app.doRequest(t, http.MethodDelete, fmt.Sprintf("/api/admin/deleteproduct?ProductID=%s", productID), nil, adminCookie)
		assert.Equal(t, http.StatusOK, resp.StatusCode)

		// DB'de silindiğini doğrula
		parsedUUID, _ := uuid.Parse(productID)
		_, err := app.ProductStore.GetProductById(ctx, parsedUUID)
		assert.Error(t, err, "Silinen ürün artık sorgulanamamalı")
	})
}

func TestIntegration_Product_PublishState(t *testing.T) {
	app := setupIntegrationApp(t)

	_, adminCookie, _ := app.createTestUser(t, fmt.Sprintf("admin-state-%s@example.com", uuid.NewString()[:8]), "Pass123!", "admin")
	ctx := context.Background()

	// 1. Published ürün oluşturalım
	prodID := uuid.NewString()
	now := time.Now()
	pCode := fmt.Sprintf("PUB%s", uuid.NewString()[:4])
	prod := models.Product{
		ProductID:           prodID,
		ProductName:         fmt.Sprintf("Yayın Test Ürünü %s", uuid.NewString()[:6]),
		Stock:               10,
		Price:               2000,
		ThumbnailImage:      uuid.NewString(),
		Description:         "Yayın durumu entegrasyon testi",
		CreatedAt:           &now,
		UpdatedAt:           &now,
		ItemType:            "PHYSICAL",
		ProductCode:         pCode,
		ProductPublishState: "published",
		TechnicalDetails:    []byte(`{}`),
	}
	cat, err := app.CategoryStore.Create(ctx, fmt.Sprintf("StateCat-%s", uuid.NewString()[:6]))
	assert.NoError(t, err)

	err = app.ProductStore.Create(ctx, &prod)
	assert.NoError(t, err)

	err = app.PPCStore.AddMany(ctx, []string{cat.Category_id}, prodID)
	assert.NoError(t, err)

	t.Run("published ürün public listede görünmeli", func(t *testing.T) {
		resp := app.doRequest(t, http.MethodGet, "/api/public/getallproducts?limit=50&offset=0", nil)
		assert.Equal(t, http.StatusOK, resp.StatusCode)

		var body structs.CategoryFilteredProducts
		_ = json.NewDecoder(resp.Body).Decode(&body)

		var found bool
		for _, p := range body.Data {
			if p.ProductID == prodID {
				found = true
				break
			}
		}
		assert.True(t, found, "Published ürün public listede yer almalı")
	})

	t.Run("yayın durumu draft'a çekildiğinde public listeden gizlenmeli", func(t *testing.T) {
		updateStatePayload := map[string]string{
			"product_id": prodID,
			"new_state":  "draft",
		}
		resp := app.doRequest(t, http.MethodPut, "/api/admin/product/publishstate", updateStatePayload, adminCookie)
		assert.Equal(t, http.StatusOK, resp.StatusCode)

		// Public listede aranır - bulunmamalı
		pubResp := app.doRequest(t, http.MethodGet, "/api/public/getallproducts?limit=50&offset=0", nil)
		var pubBody structs.CategoryFilteredProducts
		_ = json.NewDecoder(pubResp.Body).Decode(&pubBody)

		var foundInPublic bool
		for _, p := range pubBody.Data {
			if p.ProductID == prodID {
				foundInPublic = true
				break
			}
		}
		assert.False(t, foundInPublic, "Draft ürün public listede OLMAMALI")

		// Admin listesinde aranır - bulunmalı
		adminResp := app.doRequest(t, http.MethodGet, "/api/admin/getallproducts?limit=50&offset=0", nil, adminCookie)
		assert.Equal(t, http.StatusOK, adminResp.StatusCode)
		var adminResult structs.ProductResultAdmin
		_ = json.NewDecoder(adminResp.Body).Decode(&adminResult)

		var foundInAdmin bool
		for _, p := range adminResult.Data {
			if p.Product.ProductID == prodID {
				foundInAdmin = true
				break
			}
		}
		assert.True(t, foundInAdmin, "Draft ürün admin listesinde görünmeli")
	})
}

func TestIntegration_Product_SearchAndFilter(t *testing.T) {
	app := setupIntegrationApp(t)
	ctx := context.Background()

	// Kategori ve ürün oluşturalım
	cat, err := app.CategoryStore.Create(ctx, fmt.Sprintf("Kulaklıklar-%s", uuid.NewString()[:6]))
	assert.NoError(t, err)

	prodID := uuid.NewString()
	now := time.Now()
	pCode := uuid.NewString()[:8]
	prodName := fmt.Sprintf("Kablosuz Bluetooth Kulaklık %s", uuid.NewString()[:4])
	prod := models.Product{
		ProductID:           prodID,
		ProductName:         prodName,
		Stock:               15,
		Price:               3500,
		ThumbnailImage:      uuid.NewString(),
		Description:         "Gürültü engelleyici kablosuz bluetooth kulaklık",
		CreatedAt:           &now,
		UpdatedAt:           &now,
		ItemType:            "PHYSICAL",
		ProductCode:         pCode,
		ProductPublishState: "published",
		TechnicalDetails:    []byte(`{}`),
	}
	err = app.ProductStore.Create(ctx, &prod)
	assert.NoError(t, err)

	// Kategori ilişkisi ekle
	err = app.PPCStore.AddMany(ctx, []string{cat.Category_id}, prodID)
	assert.NoError(t, err)

	t.Run("full-text search ile ürün bulunabilmeli", func(t *testing.T) {
		resp := app.doRequest(t, http.MethodGet, "/api/user/searchproduct?searchQuery=Bluetooth&limit=10", nil)
		assert.Equal(t, http.StatusOK, resp.StatusCode)

		var results []structs.Getproducts
		err := json.NewDecoder(resp.Body).Decode(&results)
		assert.NoError(t, err)
		assert.GreaterOrEqual(t, len(results), 1)

		var found bool
		for _, p := range results {
			if p.ProductID == prodID {
				found = true
				break
			}
		}
		assert.True(t, found, "Arama sonucu eklenen ürünü içermeli")
	})

	t.Run("kategoriye göre ürünler filtrelenebilmeli", func(t *testing.T) {
		path := fmt.Sprintf("/api/public/getproductsbycategory?category=%s&limit=10&offset=0", cat.Category_id)
		resp := app.doRequest(t, http.MethodGet, path, nil)
		assert.Equal(t, http.StatusOK, resp.StatusCode)

		var body structs.CategoryFilteredProducts
		err := json.NewDecoder(resp.Body).Decode(&body)
		assert.NoError(t, err)
		assert.GreaterOrEqual(t, body.Count, 1)

		var found bool
		for _, p := range body.Data {
			if p.ProductID == prodID {
				found = true
				break
			}
		}
		assert.True(t, found, "Kategori listesinde ürün yer almalı")
	})
}
