package httphandlers

import (
	"context"
	services "ecommerceapi/Services"
	structs "ecommerceapi/Structs"
	"encoding/json"
	"errors"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/stretchr/testify/assert"
)

func setupCategoryHandlerTest() (*CategoryHandler, *mockCategoryRepo) {
	repo := &mockCategoryRepo{}
	svc := services.NewCategoryService(repo)
	handler := NewCategoryHandler(*svc)
	return handler, repo
}

func TestCategoryHandler_AddNewCategory(t *testing.T) {
	handler, repo := setupCategoryHandlerTest()

	t.Run("geçerli category_name ile başarıyla eklemeli", func(t *testing.T) {
		repo.createFn = func(ctx context.Context, category_name string) (structs.GetCategory, error) {
			assert.Equal(t, "Elektronik", category_name)
			return structs.GetCategory{
				Category_id:   "cat-1",
				Category_name: "Elektronik",
			}, nil
		}

		req := httptest.NewRequest(http.MethodPost, "/api/admin/createcategory?category_name=Elektronik", nil)
		rec := httptest.NewRecorder()

		handler.AddNewCategory(rec, req)

		assert.Equal(t, http.StatusOK, rec.Code)
		assert.Equal(t, "application/json", rec.Header().Get("Content-Type"))

		var resp structs.GetCategory
		err := json.Unmarshal(rec.Body.Bytes(), &resp)
		assert.NoError(t, err)
		assert.Equal(t, "cat-1", resp.Category_id)
		assert.Equal(t, "Elektronik", resp.Category_name)
	})

	t.Run("category_name parametresi eksikse 400 dönmeli", func(t *testing.T) {
		req := httptest.NewRequest(http.MethodPost, "/api/admin/createcategory", nil)
		rec := httptest.NewRecorder()

		handler.AddNewCategory(rec, req)

		assert.Equal(t, http.StatusBadRequest, rec.Code)
		assert.Contains(t, rec.Body.String(), "category_name is required")
	})

	t.Run("servis hata verdiğinde 500 dönmeli", func(t *testing.T) {
		repo.createFn = func(ctx context.Context, category_name string) (structs.GetCategory, error) {
			return structs.GetCategory{}, errors.New("insert failed")
		}

		req := httptest.NewRequest(http.MethodPost, "/api/admin/createcategory?category_name=Elektronik", nil)
		rec := httptest.NewRecorder()

		handler.AddNewCategory(rec, req)

		assert.Equal(t, http.StatusInternalServerError, rec.Code)
		assert.Contains(t, rec.Body.String(), "Failed to add new category")
	})
}

func TestCategoryHandler_DeleteCategory(t *testing.T) {
	handler, repo := setupCategoryHandlerTest()

	t.Run("geçerli category_id ile başarıyla silmeli", func(t *testing.T) {
		deleted := false
		repo.deleteFn = func(ctx context.Context, categoryId string) error {
			assert.Equal(t, "cat-1", categoryId)
			deleted = true
			return nil
		}

		req := httptest.NewRequest(http.MethodDelete, "/api/admin/deletecategory?category_id=cat-1", nil)
		rec := httptest.NewRecorder()

		handler.DeleteCategory(rec, req)

		assert.Equal(t, http.StatusOK, rec.Code)
		assert.True(t, deleted)
	})

	t.Run("category_id parametresi eksikse 400 dönmeli", func(t *testing.T) {
		req := httptest.NewRequest(http.MethodDelete, "/api/admin/deletecategory", nil)
		rec := httptest.NewRecorder()

		handler.DeleteCategory(rec, req)

		assert.Equal(t, http.StatusBadRequest, rec.Code)
		assert.Contains(t, rec.Body.String(), "category_id is required")
	})

	t.Run("servis hata verdiğinde 500 dönmeli", func(t *testing.T) {
		repo.deleteFn = func(ctx context.Context, categoryId string) error {
			return errors.New("delete failed")
		}

		req := httptest.NewRequest(http.MethodDelete, "/api/admin/deletecategory?category_id=cat-1", nil)
		rec := httptest.NewRecorder()

		handler.DeleteCategory(rec, req)

		assert.Equal(t, http.StatusInternalServerError, rec.Code)
		assert.Contains(t, rec.Body.String(), "Failed to delete category")
	})
}

func TestCategoryHandler_UpdateCategory(t *testing.T) {
	handler, repo := setupCategoryHandlerTest()

	t.Run("geçerli parametrelerle başarıyla güncellemeli", func(t *testing.T) {
		updated := false
		repo.updateFn = func(ctx context.Context, category_id string, category_name string) error {
			assert.Equal(t, "cat-1", category_id)
			assert.Equal(t, "GuncelKat", category_name)
			updated = true
			return nil
		}

		req := httptest.NewRequest(http.MethodPut, "/api/admin/updatecategory?category_id=cat-1&category_name=GuncelKat", nil)
		rec := httptest.NewRecorder()

		handler.UpdateCategory(rec, req)

		assert.Equal(t, http.StatusOK, rec.Code)
		assert.True(t, updated)
	})

	t.Run("parametrelerden biri eksikse 400 dönmeli", func(t *testing.T) {
		req := httptest.NewRequest(http.MethodPut, "/api/admin/updatecategory?category_id=cat-1", nil)
		rec := httptest.NewRecorder()

		handler.UpdateCategory(rec, req)

		assert.Equal(t, http.StatusBadRequest, rec.Code)
		assert.Contains(t, rec.Body.String(), "category_id and category_name are required")
	})

	t.Run("servis hata verdiğinde 500 dönmeli", func(t *testing.T) {
		repo.updateFn = func(ctx context.Context, category_id string, category_name string) error {
			return errors.New("update failed")
		}

		req := httptest.NewRequest(http.MethodPut, "/api/admin/updatecategory?category_id=cat-1&category_name=GuncelKat", nil)
		rec := httptest.NewRecorder()

		handler.UpdateCategory(rec, req)

		assert.Equal(t, http.StatusInternalServerError, rec.Code)
		assert.Contains(t, rec.Body.String(), "Failed to update category")
	})
}

func TestCategoryHandler_GetAll(t *testing.T) {
	handler, repo := setupCategoryHandlerTest()

	t.Run("başarıyla tüm kategorileri dönmeli", func(t *testing.T) {
		repo.getAllFn = func(ctx context.Context) ([]structs.GetCategory, error) {
			return []structs.GetCategory{
				{Category_id: "cat-1", Category_name: "Elektronik"},
				{Category_id: "cat-2", Category_name: "Moda"},
			}, nil
		}

		req := httptest.NewRequest(http.MethodGet, "/api/public/getcategories", nil)
		rec := httptest.NewRecorder()

		handler.GetAll(rec, req)

		assert.Equal(t, http.StatusOK, rec.Code)
		assert.Equal(t, "application/json", rec.Header().Get("Content-Type"))

		var resp []structs.GetCategory
		err := json.Unmarshal(rec.Body.Bytes(), &resp)
		assert.NoError(t, err)
		assert.Len(t, resp, 2)
	})

	t.Run("servis hata verdiğinde 500 dönmeli", func(t *testing.T) {
		repo.getAllFn = func(ctx context.Context) ([]structs.GetCategory, error) {
			return nil, errors.New("db failed")
		}

		req := httptest.NewRequest(http.MethodGet, "/api/public/getcategories", nil)
		rec := httptest.NewRecorder()

		handler.GetAll(rec, req)

		assert.Equal(t, http.StatusInternalServerError, rec.Code)
		assert.Contains(t, rec.Body.String(), "beklenmedik hataayla karşılaşıldı")
	})
}
