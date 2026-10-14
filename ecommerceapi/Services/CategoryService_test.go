package services

import (
	"context"
	structs "ecommerceapi/Structs"
	"errors"
	"testing"

	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

type mockCategoryRepoForCategoryService struct {
	createFn                   func(ctx context.Context, categoryName string) (structs.GetCategory, error)
	deleteFn                   func(ctx context.Context, categoryId string) error
	updateFn                   func(ctx context.Context, categoryId string, categoryName string) error
	getAllFn                   func(ctx context.Context) ([]structs.GetCategory, error)
	getByCategoryNameFn        func(ctx context.Context, categoryName string) (uuid.UUID, error)
	getCategoriesByProductIdFn func(ctx context.Context, product_id string) ([]structs.GetCategory, error)
	getByCategoryIdsFn         func(ctx context.Context, categoryNames []string) ([]string, error)
}

func (m *mockCategoryRepoForCategoryService) Create(ctx context.Context, categoryName string) (structs.GetCategory, error) {
	if m.createFn != nil {
		return m.createFn(ctx, categoryName)
	}
	return structs.GetCategory{}, nil
}

func (m *mockCategoryRepoForCategoryService) Delete(ctx context.Context, categoryId string) error {
	if m.deleteFn != nil {
		return m.deleteFn(ctx, categoryId)
	}
	return nil
}

func (m *mockCategoryRepoForCategoryService) Update(ctx context.Context, categoryId string, categoryName string) error {
	if m.updateFn != nil {
		return m.updateFn(ctx, categoryId, categoryName)
	}
	return nil
}

func (m *mockCategoryRepoForCategoryService) GetAll(ctx context.Context) ([]structs.GetCategory, error) {
	if m.getAllFn != nil {
		return m.getAllFn(ctx)
	}
	return nil, nil
}

func (m *mockCategoryRepoForCategoryService) GetByCategoryName(ctx context.Context, categoryName string) (uuid.UUID, error) {
	if m.getByCategoryNameFn != nil {
		return m.getByCategoryNameFn(ctx, categoryName)
	}
	return uuid.Nil, nil
}

func (m *mockCategoryRepoForCategoryService) GetCategoriesByProductId(ctx context.Context, product_id string) ([]structs.GetCategory, error) {
	if m.getCategoriesByProductIdFn != nil {
		return m.getCategoriesByProductIdFn(ctx, product_id)
	}
	return nil, nil
}

func (m *mockCategoryRepoForCategoryService) GetByCategoryIds(ctx context.Context, categoryNames []string) ([]string, error) {
	if m.getByCategoryIdsFn != nil {
		return m.getByCategoryIdsFn(ctx, categoryNames)
	}
	return nil, nil
}

func TestCategoryService_AddNewCategory(t *testing.T) {
	mockRepo := &mockCategoryRepoForCategoryService{}
	svc := NewCategoryService(mockRepo)
	ctx := context.Background()

	t.Run("veritabanı Create hata verirse hatayı dönmeli", func(t *testing.T) {
		mockRepo.createFn = func(ctx context.Context, name string) (structs.GetCategory, error) {
			return structs.GetCategory{}, errors.New("db create error")
		}

		res, err := svc.AddNewCategory(ctx, "Elektronik")
		assert.Error(t, err)
		assert.Equal(t, structs.GetCategory{}, res)
	})

	t.Run("başarıyla yeni kategori oluşturup dönmeli", func(t *testing.T) {
		expectedCategory := structs.GetCategory{
			Category_id:   "cat-123",
			Category_name: "Elektronik",
		}
		mockRepo.createFn = func(ctx context.Context, name string) (structs.GetCategory, error) {
			assert.Equal(t, "Elektronik", name)
			return expectedCategory, nil
		}

		res, err := svc.AddNewCategory(ctx, "Elektronik")
		assert.NoError(t, err)
		assert.Equal(t, expectedCategory, res)
	})
}

func TestCategoryService_DeleteCategory(t *testing.T) {
	mockRepo := &mockCategoryRepoForCategoryService{}
	svc := NewCategoryService(mockRepo)
	ctx := context.Background()

	t.Run("veritabanı Delete hata verirse hatayı dönmeli", func(t *testing.T) {
		mockRepo.deleteFn = func(ctx context.Context, id string) error {
			return errors.New("db delete error")
		}

		err := svc.DeleteCategory(ctx, "cat-123")
		assert.Error(t, err)
	})

	t.Run("başarıyla kategori silmeli", func(t *testing.T) {
		called := false
		mockRepo.deleteFn = func(ctx context.Context, id string) error {
			called = true
			assert.Equal(t, "cat-123", id)
			return nil
		}

		err := svc.DeleteCategory(ctx, "cat-123")
		assert.NoError(t, err)
		assert.True(t, called)
	})
}

func TestCategoryService_UpdateCategory(t *testing.T) {
	mockRepo := &mockCategoryRepoForCategoryService{}
	svc := NewCategoryService(mockRepo)
	ctx := context.Background()

	t.Run("veritabanı Update hata verirse hatayı dönmeli", func(t *testing.T) {
		mockRepo.updateFn = func(ctx context.Context, id string, name string) error {
			return errors.New("db update error")
		}

		err := svc.UpdateCategory(ctx, "cat-123", "Giyim")
		assert.Error(t, err)
	})

	t.Run("başarıyla kategori güncellemeli", func(t *testing.T) {
		called := false
		mockRepo.updateFn = func(ctx context.Context, id string, name string) error {
			called = true
			assert.Equal(t, "cat-123", id)
			assert.Equal(t, "Giyim", name)
			return nil
		}

		err := svc.UpdateCategory(ctx, "cat-123", "Giyim")
		assert.NoError(t, err)
		assert.True(t, called)
	})
}

func TestCategoryService_GetAll(t *testing.T) {
	mockRepo := &mockCategoryRepoForCategoryService{}
	svc := NewCategoryService(mockRepo)
	ctx := context.Background()

	t.Run("veritabanı GetAll hata verirse hatayı dönmeli", func(t *testing.T) {
		mockRepo.getAllFn = func(ctx context.Context) ([]structs.GetCategory, error) {
			return nil, errors.New("db getall error")
		}

		res, err := svc.GetAll(ctx)
		assert.Error(t, err)
		assert.Nil(t, res)
	})

	t.Run("başarıyla tüm kategorileri listelemeli", func(t *testing.T) {
		expected := []structs.GetCategory{
			{Category_id: "cat-1", Category_name: "Kitap"},
			{Category_id: "cat-2", Category_name: "Müzik"},
		}
		mockRepo.getAllFn = func(ctx context.Context) ([]structs.GetCategory, error) {
			return expected, nil
		}

		res, err := svc.GetAll(ctx)
		assert.NoError(t, err)
		require.NotNil(t, res)
		assert.Equal(t, expected, res)
	})
}
