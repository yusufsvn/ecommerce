package services

import (
	"bytes"
	"context"
	"errors"
	"io"
	"mime/multipart"
	"testing"

	constants "ecommerceapi/Constants"
	Customerrors "ecommerceapi/CustomErrors"
	Entities "ecommerceapi/Database/Models"
	models "ecommerceapi/Database/Models"
	st "ecommerceapi/Structs"
	structs "ecommerceapi/Structs"

	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

type mockProductRepoForProductService struct {
	createFn                     func(ctx context.Context, product *models.Product) error
	deleteFn                     func(ctx context.Context, productId *uuid.UUID) error
	updateFn                     func(ctx context.Context, product *models.Product) error
	patchFn                      func(ctx context.Context, product *models.Product) error
	changeThumbnailImageFn       func(ctx context.Context, new_id string, product_id string) error
	getAllFn                     func(limit int, offset int, ctx context.Context) ([]structs.Getproducts, error)
	getProductByIdFn             func(ctx context.Context, ProductID uuid.UUID) (*structs.ProductWithImage, error)
	getManyFn                    func(ctx context.Context, ids []string) ([]*structs.ProductForCart, error)
	getProductCountFn            func(ctx context.Context) (int, error)
	getProductsByCategoryIdFn    func(ctx context.Context, category_ids []string, limit string, offset string) ([]structs.Getproducts, error)
	getCountProductsByCategoryFn func(ctx context.Context, category_ids []string) (int, error)
	searchProductCountFn         func(ctx context.Context, searchQuery string) int
	searchProductFn              func(ctx context.Context, searchQuery string, limit string) ([]st.Getproducts, error)
	alternateSearchFn            func(ctx context.Context, searchQuery string, limit string) ([]st.Getproducts, error)
	filteredAndSortedProductsFn  func(ctx context.Context, category_ids []string, limit, offset, sortBy, sortOrder string) ([]structs.Getproducts, error)
	searchProductCountAdminFn    func(ctx context.Context, searchQuery string) int
	searchProductAdminFn         func(ctx context.Context, searchQuery string, limit string) ([]st.GetproductsAdmin, error)
	alternateSearchAdminFn       func(ctx context.Context, searchQuery string, limit string) ([]st.GetproductsAdmin, error)
	getProductCountAdminFn       func(ctx context.Context) (int, error)
	getAllAdminFn                func(limit int, offset int, ctx context.Context) ([]structs.GetproductsAdmin, error)
	changePublishStateFn         func(ctx context.Context, product_id string, newState string) error
}

func (m *mockProductRepoForProductService) Create(ctx context.Context, product *models.Product) error {
	if m.createFn != nil {
		return m.createFn(ctx, product)
	}
	return nil
}
func (m *mockProductRepoForProductService) Delete(ctx context.Context, productId *uuid.UUID) error {
	if m.deleteFn != nil {
		return m.deleteFn(ctx, productId)
	}
	return nil
}
func (m *mockProductRepoForProductService) Update(ctx context.Context, product *models.Product) error {
	if m.updateFn != nil {
		return m.updateFn(ctx, product)
	}
	return nil
}
func (m *mockProductRepoForProductService) Patch(ctx context.Context, product *models.Product) error {
	if m.patchFn != nil {
		return m.patchFn(ctx, product)
	}
	return nil
}
func (m *mockProductRepoForProductService) ChangeThumbnailImage(ctx context.Context, new_id string, product_id string) error {
	if m.changeThumbnailImageFn != nil {
		return m.changeThumbnailImageFn(ctx, new_id, product_id)
	}
	return nil
}
func (m *mockProductRepoForProductService) GetAll(limit int, offset int, ctx context.Context) ([]structs.Getproducts, error) {
	if m.getAllFn != nil {
		return m.getAllFn(limit, offset, ctx)
	}
	return nil, nil
}
func (m *mockProductRepoForProductService) GetProductById(ctx context.Context, ProductID uuid.UUID) (*structs.ProductWithImage, error) {
	if m.getProductByIdFn != nil {
		return m.getProductByIdFn(ctx, ProductID)
	}
	return nil, nil
}
func (m *mockProductRepoForProductService) GetMany(ctx context.Context, ids []string) ([]*structs.ProductForCart, error) {
	if m.getManyFn != nil {
		return m.getManyFn(ctx, ids)
	}
	return nil, nil
}
func (m *mockProductRepoForProductService) GetProductCount(ctx context.Context) (int, error) {
	if m.getProductCountFn != nil {
		return m.getProductCountFn(ctx)
	}
	return 0, nil
}
func (m *mockProductRepoForProductService) GetProductsByCategoryId(ctx context.Context, category_ids []string, limit string, offset string) ([]structs.Getproducts, error) {
	if m.getProductsByCategoryIdFn != nil {
		return m.getProductsByCategoryIdFn(ctx, category_ids, limit, offset)
	}
	return nil, nil
}
func (m *mockProductRepoForProductService) GetCountProductsByCategory(ctx context.Context, category_ids []string) (int, error) {
	if m.getCountProductsByCategoryFn != nil {
		return m.getCountProductsByCategoryFn(ctx, category_ids)
	}
	return 0, nil
}
func (m *mockProductRepoForProductService) SearchProductCount(ctx context.Context, searchQuery string) int {
	if m.searchProductCountFn != nil {
		return m.searchProductCountFn(ctx, searchQuery)
	}
	return 0
}
func (m *mockProductRepoForProductService) SearchProduct(ctx context.Context, searchQuery string, limit string) ([]st.Getproducts, error) {
	if m.searchProductFn != nil {
		return m.searchProductFn(ctx, searchQuery, limit)
	}
	return nil, nil
}
func (m *mockProductRepoForProductService) AlternateSearch(ctx context.Context, searchQuery string, limit string) ([]st.Getproducts, error) {
	if m.alternateSearchFn != nil {
		return m.alternateSearchFn(ctx, searchQuery, limit)
	}
	return nil, nil
}
func (m *mockProductRepoForProductService) FilteredAndSortedProducts(ctx context.Context, category_ids []string, limit, offset, sortBy, sortOrder string) ([]structs.Getproducts, error) {
	if m.filteredAndSortedProductsFn != nil {
		return m.filteredAndSortedProductsFn(ctx, category_ids, limit, offset, sortBy, sortOrder)
	}
	return nil, nil
}
func (m *mockProductRepoForProductService) SearchProductCountAdmin(ctx context.Context, searchQuery string) int {
	if m.searchProductCountAdminFn != nil {
		return m.searchProductCountAdminFn(ctx, searchQuery)
	}
	return 0
}
func (m *mockProductRepoForProductService) SearchProductAdmin(ctx context.Context, searchQuery string, limit string) ([]st.GetproductsAdmin, error) {
	if m.searchProductAdminFn != nil {
		return m.searchProductAdminFn(ctx, searchQuery, limit)
	}
	return nil, nil
}
func (m *mockProductRepoForProductService) AlternateSearchAdmin(ctx context.Context, searchQuery string, limit string) ([]st.GetproductsAdmin, error) {
	if m.alternateSearchAdminFn != nil {
		return m.alternateSearchAdminFn(ctx, searchQuery, limit)
	}
	return nil, nil
}
func (m *mockProductRepoForProductService) GetProductCountAdmin(ctx context.Context) (int, error) {
	if m.getProductCountAdminFn != nil {
		return m.getProductCountAdminFn(ctx)
	}
	return 0, nil
}
func (m *mockProductRepoForProductService) GetAllAdmin(limit int, offset int, ctx context.Context) ([]structs.GetproductsAdmin, error) {
	if m.getAllAdminFn != nil {
		return m.getAllAdminFn(limit, offset, ctx)
	}
	return nil, nil
}
func (m *mockProductRepoForProductService) ChangePublishState(ctx context.Context, product_id string, newState string) error {
	if m.changePublishStateFn != nil {
		return m.changePublishStateFn(ctx, product_id, newState)
	}
	return nil
}

type mockProductProductCategoryRepoForProductService struct {
	getByProductIDFn          func(ctx context.Context, productID uuid.UUID) ([]Entities.ProductProductCategory, error)
	getByCategoryIDFn          func(ctx context.Context, categoryID uuid.UUID) ([]Entities.ProductProductCategory, error)
	getAllProductCategoriesFn func(ctx context.Context, ids []string) ([]*structs.GetAllCategory, error)
	addFn                     func(ctx context.Context, item *Entities.ProductProductCategory) error
	removeFn                  func(ctx context.Context, productID uuid.UUID, categoryID uuid.UUID) error
	getByProductIDOneRowFn    func(ctx context.Context, productID uuid.UUID) (uuid.UUID, error)
	addManyFn                 func(ctx context.Context, items []string, product_id string) error
	removeManyFn              func(ctx context.Context, items []string, product_id string) error
}

func (m *mockProductProductCategoryRepoForProductService) GetByProductID(ctx context.Context, productID uuid.UUID) ([]Entities.ProductProductCategory, error) {
	if m.getByProductIDFn != nil {
		return m.getByProductIDFn(ctx, productID)
	}
	return nil, nil
}
func (m *mockProductProductCategoryRepoForProductService) GetByCategoryID(ctx context.Context, categoryID uuid.UUID) ([]Entities.ProductProductCategory, error) {
	if m.getByCategoryIDFn != nil {
		return m.getByCategoryIDFn(ctx, categoryID)
	}
	return nil, nil
}
func (m *mockProductProductCategoryRepoForProductService) GetAllProductCategories(ctx context.Context, ids []string) ([]*structs.GetAllCategory, error) {
	if m.getAllProductCategoriesFn != nil {
		return m.getAllProductCategoriesFn(ctx, ids)
	}
	return nil, nil
}
func (m *mockProductProductCategoryRepoForProductService) Add(ctx context.Context, item *Entities.ProductProductCategory) error {
	if m.addFn != nil {
		return m.addFn(ctx, item)
	}
	return nil
}
func (m *mockProductProductCategoryRepoForProductService) Remove(ctx context.Context, productID uuid.UUID, categoryID uuid.UUID) error {
	if m.removeFn != nil {
		return m.removeFn(ctx, productID, categoryID)
	}
	return nil
}
func (m *mockProductProductCategoryRepoForProductService) GetByProductIDOneRow(ctx context.Context, productID uuid.UUID) (uuid.UUID, error) {
	if m.getByProductIDOneRowFn != nil {
		return m.getByProductIDOneRowFn(ctx, productID)
	}
	return uuid.Nil, nil
}
func (m *mockProductProductCategoryRepoForProductService) AddMany(ctx context.Context, items []string, product_id string) error {
	if m.addManyFn != nil {
		return m.addManyFn(ctx, items, product_id)
	}
	return nil
}
func (m *mockProductProductCategoryRepoForProductService) RemoveMany(ctx context.Context, items []string, product_id string) error {
	if m.removeManyFn != nil {
		return m.removeManyFn(ctx, items, product_id)
	}
	return nil
}

func setupProductServiceTest() (
	*ProductService,
	*mockProductRepoForProductService,
	*mockCategoryRepoForCategoryService,
	*mockProductProductCategoryRepoForProductService,
	*mockImageRepoForImageService,
	*mockFileManagerForImageService,
) {
	prodRepo := &mockProductRepoForProductService{}
	catRepo := &mockCategoryRepoForCategoryService{}
	ppcRepo := &mockProductProductCategoryRepoForProductService{}
	imgRepo := &mockImageRepoForImageService{}
	fm := &mockFileManagerForImageService{}

	imageService := NewImageService(imgRepo, fm)
	svc := NewProductService(
		prodRepo,
		catRepo,
		ppcRepo,
		"http://storage.api",
		"http://public.storage",
		"test-bucket",
		*imageService,
	)

	return svc, prodRepo, catRepo, ppcRepo, imgRepo, fm
}

func createProductTestFileHeader(t *testing.T, filename string, content []byte) *multipart.FileHeader {
	body := new(bytes.Buffer)
	writer := multipart.NewWriter(body)
	part, err := writer.CreateFormFile("file", filename)
	require.NoError(t, err)
	_, err = part.Write(content)
	require.NoError(t, err)
	err = writer.Close()
	require.NoError(t, err)

	reader := multipart.NewReader(body, writer.Boundary())
	form, err := reader.ReadForm(1024 * 1024)
	require.NoError(t, err)
	return form.File["file"][0]
}

func TestProductService_CreateProduct(t *testing.T) {
	svc, prodRepo, catRepo, ppcRepo, _, _ := setupProductServiceTest()
	ctx := context.Background()

	t.Run("geçersiz ItemType verilirse ErrInvalidItemType dönmeli", func(t *testing.T) {
		req := &structs.CreateProductRequest{
			ItemType: "INVALID_TYPE",
		}
		err := svc.CreateProduct(ctx, req)
		assert.ErrorIs(t, err, Customerrors.ErrInvalidItemType)
	})

	t.Run("ProductStore.Create hata verirse hata dönmeli", func(t *testing.T) {
		req := &structs.CreateProductRequest{
			ItemType:    constants.PhysicalProduct,
			ProductName: "Test Ürün",
			Price:       100,
			Stock:       10,
		}

		prodRepo.createFn = func(ctx context.Context, product *models.Product) error {
			return errors.New("db create error")
		}

		err := svc.CreateProduct(ctx, req)
		assert.Error(t, err)
		assert.Contains(t, err.Error(), "Product oluşturulamadı")
	})

	t.Run("tek kategori ile başarıyla ürün ve kategori ilişkisi oluşturmalı", func(t *testing.T) {
		req := &structs.CreateProductRequest{
			ItemType:        constants.PhysicalProduct,
			ProductName:     "Laptop",
			Price:           15000,
			Stock:           5,
			ProductCategory: []string{"Elektronik"},
		}

		prodCreated := false
		prodRepo.createFn = func(ctx context.Context, product *models.Product) error {
			prodCreated = true
			assert.Equal(t, "Laptop", product.ProductName)
			assert.Equal(t, constants.PublishStateUnpublished, product.ProductPublishState)
			return nil
		}

		catRepo.getByCategoryIdsFn = func(ctx context.Context, names []string) ([]string, error) {
			return []string{"cat-1"}, nil
		}

		ppcAdded := false
		ppcRepo.addFn = func(ctx context.Context, item *Entities.ProductProductCategory) error {
			ppcAdded = true
			assert.Equal(t, "cat-1", item.CategoryID)
			return nil
		}

		err := svc.CreateProduct(ctx, req)
		assert.NoError(t, err)
		assert.True(t, prodCreated)
		assert.True(t, ppcAdded)
	})

	t.Run("birden fazla kategori ile başarıyla AddMany çağrılmalı", func(t *testing.T) {
		req := &structs.CreateProductRequest{
			ItemType:        constants.PhysicalProduct,
			ProductName:     "Telefon",
			Price:           8000,
			Stock:           20,
			ProductCategory: []string{"Elektronik", "Mobil"},
		}

		prodRepo.createFn = func(ctx context.Context, product *models.Product) error {
			return nil
		}

		catRepo.getByCategoryIdsFn = func(ctx context.Context, names []string) ([]string, error) {
			return []string{"cat-1", "cat-2"}, nil
		}

		addManyCalled := false
		ppcRepo.addManyFn = func(ctx context.Context, items []string, product_id string) error {
			addManyCalled = true
			assert.Equal(t, []string{"cat-1", "cat-2"}, items)
			return nil
		}

		err := svc.CreateProduct(ctx, req)
		assert.NoError(t, err)
		assert.True(t, addManyCalled)
	})
}

func TestProductService_GetAllProducts(t *testing.T) {
	svc, prodRepo, _, ppcRepo, _, _ := setupProductServiceTest()
	ctx := context.Background()

	t.Run("GetProductCountAdmin hata verirse hata dönmeli", func(t *testing.T) {
		prodRepo.getProductCountAdminFn = func(ctx context.Context) (int, error) {
			return -1, errors.New("count error")
		}

		res, err := svc.GetAllProducts(ctx, 10, 0)
		assert.Error(t, err)
		assert.Nil(t, res)
	})

	t.Run("toplam sayı 0 ise boş liste dönmeli", func(t *testing.T) {
		prodRepo.getProductCountAdminFn = func(ctx context.Context) (int, error) {
			return 0, nil
		}

		res, err := svc.GetAllProducts(ctx, 10, 0)
		assert.NoError(t, err)
		require.NotNil(t, res)
		assert.Equal(t, 0, res.Count)
		assert.Empty(t, res.Data)
	})

	t.Run("başarıyla admin ürün listesini ve kategorilerini dönmeli", func(t *testing.T) {
		prodRepo.getProductCountAdminFn = func(ctx context.Context) (int, error) {
			return 1, nil
		}
		prodRepo.getAllAdminFn = func(limit int, offset int, ctx context.Context) ([]structs.GetproductsAdmin, error) {
			return []structs.GetproductsAdmin{
				{ProductID: "prod-1", ProductName: "Klavye", ThumbnailImage: "thumb.jpg"},
			}, nil
		}
		ppcRepo.getAllProductCategoriesFn = func(ctx context.Context, ids []string) ([]*structs.GetAllCategory, error) {
			return []*structs.GetAllCategory{
				{ProductId: "prod-1", Categories: []structs.CategoryDbJson{{CategoryName: "Aksesuar"}}},
			}, nil
		}

		res, err := svc.GetAllProducts(ctx, 10, 0)
		assert.NoError(t, err)
		require.NotNil(t, res)
		assert.Equal(t, 1, res.Count)
		assert.Len(t, res.Data, 1)
		assert.Contains(t, res.Data[0].Product.ThumbnailImage, "test-bucket")
		assert.Len(t, res.Data[0].Categories, 1)
	})
}

func TestProductService_PubGetAllProducts(t *testing.T) {
	svc, prodRepo, _, _, _, _ := setupProductServiceTest()
	ctx := context.Background()

	t.Run("GetProductCount hata verirse hata dönmeli", func(t *testing.T) {
		prodRepo.getProductCountFn = func(ctx context.Context) (int, error) {
			return 0, errors.New("count error")
		}

		res, err := svc.PubGetAllProducts(ctx, 10, 0)
		assert.Error(t, err)
		assert.Nil(t, res)
	})

	t.Run("başarıyla herkese açık ürün listesini dönmeli", func(t *testing.T) {
		prodRepo.getProductCountFn = func(ctx context.Context) (int, error) {
			return 1, nil
		}
		prodRepo.getAllFn = func(limit int, offset int, ctx context.Context) ([]structs.Getproducts, error) {
			return []structs.Getproducts{
				{ProductID: "prod-1", ProductName: "Monitör", ThumbnailImage: "monitor.jpg"},
			}, nil
		}

		res, err := svc.PubGetAllProducts(ctx, 10, 0)
		assert.NoError(t, err)
		require.NotNil(t, res)
		assert.Equal(t, 1, res.Count)
		assert.Len(t, res.Data, 1)
		assert.Contains(t, res.Data[0].ThumbnailImage, "test-bucket")
	})
}

func TestProductService_GetProductById(t *testing.T) {
	svc, prodRepo, catRepo, _, _, _ := setupProductServiceTest()
	ctx := context.Background()

	t.Run("geçersiz UUID verildiğinde hata dönmeli", func(t *testing.T) {
		res, err := svc.GetProductById(ctx, "invalid-uuid")
		assert.Error(t, err)
		assert.Nil(t, res)
	})

	t.Run("başarıyla ürün detayını dönmeli", func(t *testing.T) {
		prodID := uuid.New()
		thumbUUID := uuid.New()
		imgUUID := uuid.New()
		prodRepo.getProductByIdFn = func(ctx context.Context, pID uuid.UUID) (*structs.ProductWithImage, error) {
			assert.Equal(t, prodID, pID)
			return &structs.ProductWithImage{
				ProductID:        prodID.String(),
				ProductName:      "Saat",
				ThumbnailImage:   thumbUUID.String(),
				TechnicalDetails: []byte(`{"renk":"siyah"}`),
				Images: []structs.ProductImage{
					{ImageID: imgUUID, Path: "img1.jpg"},
				},
			}, nil
		}
		catRepo.getCategoriesByProductIdFn = func(ctx context.Context, product_id string) ([]structs.GetCategory, error) {
			return []structs.GetCategory{
				{Category_id: "cat-1", Category_name: "Moda"},
			}, nil
		}

		res, err := svc.GetProductById(ctx, prodID.String())
		assert.NoError(t, err)
		require.NotNil(t, res)
		assert.Equal(t, "Saat", res.ProductName)
		assert.Len(t, res.Categories, 1)
		assert.Len(t, res.Images, 2, "Thumbnail + resim = 2 görsel")
		assert.Equal(t, "siyah", res.TechnicalDetails["renk"])
	})
}

func TestProductService_DeleteProduct(t *testing.T) {
	svc, prodRepo, _, ppcRepo, imgRepo, fm := setupProductServiceTest()
	ctx := context.Background()

	t.Run("geçersiz UUID verildiğinde hata dönmeli", func(t *testing.T) {
		err := svc.DeleteProduct(ctx, "not-a-uuid")
		assert.Error(t, err)
	})

	t.Run("başarıyla ürünü ve ilişkilerini silmeli", func(t *testing.T) {
		prodID := uuid.New()
		catID := uuid.New()

		imgRepo.getProductImagesFn = func(ctx context.Context, product_id string) ([]string, error) {
			return []string{"img1"}, nil
		}
		fm.deleteMultipleFileFn = func(ctx context.Context, images []string) error { return nil }
		imgRepo.deleteImageFn = func(ctx context.Context, images []string) error { return nil }

		ppcRepo.getByProductIDOneRowFn = func(ctx context.Context, productID uuid.UUID) (uuid.UUID, error) {
			assert.Equal(t, prodID, productID)
			return catID, nil
		}

		ppcDeleted := false
		ppcRepo.removeFn = func(ctx context.Context, productID uuid.UUID, categoryID uuid.UUID) error {
			ppcDeleted = true
			assert.Equal(t, prodID, productID)
			assert.Equal(t, catID, categoryID)
			return nil
		}

		prodDeleted := false
		prodRepo.deleteFn = func(ctx context.Context, productId *uuid.UUID) error {
			prodDeleted = true
			assert.Equal(t, prodID, *productId)
			return nil
		}

		err := svc.DeleteProduct(ctx, prodID.String())
		assert.NoError(t, err)
		assert.True(t, ppcDeleted)
		assert.True(t, prodDeleted)
	})

	t.Run("görseller silinirken hata oluştuğunda hata dönmeli", func(t *testing.T) {
		prodID := uuid.New()
		imgRepo.getProductImagesFn = func(ctx context.Context, product_id string) ([]string, error) {
			return nil, errors.New("image fetch failed")
		}

		err := svc.DeleteProduct(ctx, prodID.String())
		assert.Error(t, err)
		assert.Contains(t, err.Error(), "imageler silinirken hatayla karşılaşıldı")
	})

	t.Run("kategori ilişkisi getirilirken hata oluşursa hata dönmeli", func(t *testing.T) {
		prodID := uuid.New()
		imgRepo.getProductImagesFn = func(ctx context.Context, product_id string) ([]string, error) {
			return []string{}, nil
		}
		ppcRepo.getByProductIDOneRowFn = func(ctx context.Context, productID uuid.UUID) (uuid.UUID, error) {
			return uuid.Nil, errors.New("cat not found")
		}

		err := svc.DeleteProduct(ctx, prodID.String())
		assert.Error(t, err)
		assert.Contains(t, err.Error(), "kategori id getirme hatası")
	})

	t.Run("ürün dbden silinirken hata oluşursa hata dönmeli", func(t *testing.T) {
		prodID := uuid.New()
		imgRepo.getProductImagesFn = func(ctx context.Context, product_id string) ([]string, error) {
			return []string{}, nil
		}
		ppcRepo.getByProductIDOneRowFn = func(ctx context.Context, productID uuid.UUID) (uuid.UUID, error) {
			return uuid.New(), nil
		}
		ppcRepo.removeFn = func(ctx context.Context, productID uuid.UUID, categoryID uuid.UUID) error {
			return nil
		}
		prodRepo.deleteFn = func(ctx context.Context, productId *uuid.UUID) error {
			return errors.New("db delete error")
		}

		err := svc.DeleteProduct(ctx, prodID.String())
		assert.Error(t, err)
		assert.Contains(t, err.Error(), "Db silme sorgusunda hata")
	})
}

func TestProductService_UpdateProduct(t *testing.T) {
	svc, prodRepo, catRepo, ppcRepo, _, _ := setupProductServiceTest()
	ctx := context.Background()

	t.Run("boş veri verildiğinde hata dönmeli", func(t *testing.T) {
		err := svc.UpdateProduct(ctx, nil)
		assert.Error(t, err)
	})

	t.Run("başarıyla ürünü ve kategorilerini güncellemeli", func(t *testing.T) {
		data := &structs.UpdateProductRequest{
			ProductID:   "prod-1",
			ProductName: "Güncel Ürün",
			Price:       200,
			Stock:       15,
			Categories: structs.UpdateCategories{
				NewCategories:     []string{"YeniKat"},
				DeletedCategories: []string{"EskiKat"},
			},
		}

		prodUpdated := false
		prodRepo.updateFn = func(ctx context.Context, product *models.Product) error {
			prodUpdated = true
			assert.Equal(t, "Güncel Ürün", product.ProductName)
			return nil
		}

		catRepo.getByCategoryIdsFn = func(ctx context.Context, names []string) ([]string, error) {
			if len(names) > 0 && names[0] == "YeniKat" {
				return []string{"cat-new"}, nil
			}
			return []string{"cat-old"}, nil
		}

		newAdded := false
		ppcRepo.addManyFn = func(ctx context.Context, items []string, product_id string) error {
			newAdded = true
			assert.Equal(t, []string{"cat-new"}, items)
			return nil
		}

		oldRemoved := false
		ppcRepo.removeManyFn = func(ctx context.Context, items []string, product_id string) error {
			oldRemoved = true
			assert.Equal(t, []string{"cat-old"}, items)
			return nil
		}

		err := svc.UpdateProduct(ctx, data)
		assert.NoError(t, err)
		assert.True(t, prodUpdated)
		assert.True(t, newAdded)
		assert.True(t, oldRemoved)
	})
}

func TestProductService_PatchProduct(t *testing.T) {
	svc, prodRepo, _, _, _, _ := setupProductServiceTest()
	ctx := context.Background()

	t.Run("boş veri verildiğinde hata dönmeli", func(t *testing.T) {
		err := svc.PatchProduct(ctx, nil)
		assert.Error(t, err)
	})

	t.Run("başarıyla patch uygulamalı", func(t *testing.T) {
		data := &structs.PatchProductRequest{
			ProductID:   "prod-1",
			ProductName: "Yeni İsim",
			Price:       500,
		}

		patched := false
		prodRepo.patchFn = func(ctx context.Context, product *models.Product) error {
			patched = true
			assert.Equal(t, "prod-1", product.ProductID)
			assert.Equal(t, "Yeni İsim", product.ProductName)
			assert.Equal(t, int64(500), product.Price)
			return nil
		}

		err := svc.PatchProduct(ctx, data)
		assert.NoError(t, err)
		assert.True(t, patched)
	})
}

func TestProductService_AddProductImage(t *testing.T) {
	svc, prodRepo, _, _, imgRepo, fm := setupProductServiceTest()
	ctx := context.Background()

	fh := createProductTestFileHeader(t, "new_pic.jpg", []byte("image-data"))

	t.Run("isChanged true iken ChangeThumbnailImage çağrılmalı", func(t *testing.T) {
		fm.uploadFileFn = func(fileName string, fileSize int64, prefix string, ctx context.Context, reader io.Reader) (string, error) {
			return fileName, nil
		}
		imgRepo.addOneImageFn = func(ctx context.Context, image models.Image) error {
			return nil
		}

		thumbChanged := false
		prodRepo.changeThumbnailImageFn = func(ctx context.Context, new_id string, product_id string) error {
			thumbChanged = true
			assert.Equal(t, "prod-1", product_id)
			assert.NotEmpty(t, new_id)
			return nil
		}

		err := svc.AddProductImage(ctx, []*multipart.FileHeader{fh}, "prod-1", true)
		assert.NoError(t, err)
		assert.True(t, thumbChanged)
	})
}

func TestProductService_GetProductsByCategoryId(t *testing.T) {
	svc, prodRepo, _, _, _, _ := setupProductServiceTest()
	ctx := context.Background()

	t.Run("sayı sorgusu hata verirse hata dönmeli", func(t *testing.T) {
		prodRepo.getCountProductsByCategoryFn = func(ctx context.Context, category_ids []string) (int, error) {
			return -1, errors.New("count error")
		}

		res, err := svc.GetProductsByCategoryId(ctx, []string{"cat-1"}, "10", "0")
		assert.Error(t, err)
		assert.Nil(t, res)
	})

	t.Run("başarıyla kategori ürünlerini dönmeli", func(t *testing.T) {
		prodRepo.getCountProductsByCategoryFn = func(ctx context.Context, category_ids []string) (int, error) {
			return 1, nil
		}
		prodRepo.getProductsByCategoryIdFn = func(ctx context.Context, category_ids []string, limit string, offset string) ([]structs.Getproducts, error) {
			return []structs.Getproducts{
				{ProductID: "p-1", ProductName: "Kulaklık", ThumbnailImage: "thumb.jpg"},
			}, nil
		}

		res, err := svc.GetProductsByCategoryId(ctx, []string{"cat-1"}, "10", "0")
		assert.NoError(t, err)
		require.NotNil(t, res)
		assert.Equal(t, 1, res.Count)
		assert.Len(t, res.Data, 1)
		assert.Contains(t, res.Data[0].ThumbnailImage, "test-bucket")
	})
}

func TestProductService_GetFilteredAndSortedProducts(t *testing.T) {
	svc, prodRepo, _, _, _, _ := setupProductServiceTest()
	ctx := context.Background()

	t.Run("başarıyla filtrelenmiş ve sıralanmış ürünleri dönmeli", func(t *testing.T) {
		prodRepo.filteredAndSortedProductsFn = func(ctx context.Context, category_ids []string, limit, offset, sortBy, sortOrder string) ([]structs.Getproducts, error) {
			assert.Equal(t, constants.Price, sortBy)
			assert.Equal(t, constants.Asc, sortOrder)
			return []structs.Getproducts{
				{ProductID: "p-1", ProductName: "Ucuz Ürün", ThumbnailImage: "thumb.jpg"},
			}, nil
		}
		prodRepo.getCountProductsByCategoryFn = func(ctx context.Context, category_ids []string) (int, error) {
			return 1, nil
		}

		res, err := svc.GetFilteredAndSortedProducts(ctx, []string{"cat-1"}, "10", "0", "priceASC")
		assert.NoError(t, err)
		require.NotNil(t, res)
		assert.Equal(t, 1, res.Count)
		assert.Len(t, res.Data, 1)
	})
}

func TestProductService_UpdateProductPublishState(t *testing.T) {
	svc, prodRepo, _, _, _, _ := setupProductServiceTest()
	ctx := context.Background()

	t.Run("geçersiz durum verildiğinde hata dönmeli", func(t *testing.T) {
		err := svc.UpdateProductPublishState(ctx, "prod-1", "invalid_state")
		assert.Error(t, err)
		assert.Contains(t, err.Error(), "geçersiz yayın durumu")
	})

	t.Run("başarıyla yayın durumunu güncellemeli", func(t *testing.T) {
		called := false
		prodRepo.changePublishStateFn = func(ctx context.Context, product_id string, newState string) error {
			called = true
			assert.Equal(t, "prod-1", product_id)
			assert.Equal(t, constants.PublishStatePublished, newState)
			return nil
		}

		err := svc.UpdateProductPublishState(ctx, "prod-1", constants.PublishStatePublished)
		assert.NoError(t, err)
		assert.True(t, called)
	})
}

func TestProductService_GetCategories(t *testing.T) {
	svc, _, catRepo, _, _, _ := setupProductServiceTest()
	ctx := context.Background()

	t.Run("başarıyla kategorileri dönmeli", func(t *testing.T) {
		expected := []structs.GetCategory{
			{Category_id: "c-1", Category_name: "Kitap"},
		}
		catRepo.getAllFn = func(ctx context.Context) ([]structs.GetCategory, error) {
			return expected, nil
		}

		res, err := svc.GetCategories(ctx)
		assert.NoError(t, err)
		assert.Equal(t, expected, res)
	})
}

func TestProductService_GetManyProducts(t *testing.T) {
	svc, prodRepo, _, _, _, _ := setupProductServiceTest()
	ctx := context.Background()

	t.Run("başarıyla çoklu ürünleri dönmeli", func(t *testing.T) {
		expected := []*structs.ProductForCart{
			{ProductID: "p-1", ProductName: "Ürün 1", Price: 100},
		}
		prodRepo.getManyFn = func(ctx context.Context, ids []string) ([]*structs.ProductForCart, error) {
			return expected, nil
		}

		res, err := svc.GetManyProducts(ctx, []string{"p-1"})
		assert.NoError(t, err)
		assert.Equal(t, expected, res)
	})
}

func TestProductService_SearchProduct(t *testing.T) {
	svc, prodRepo, _, _, _, _ := setupProductServiceTest()
	ctx := context.Background()

	t.Run("sayı 0 ise boş liste dönmeli", func(t *testing.T) {
		prodRepo.searchProductCountFn = func(ctx context.Context, searchQuery string) int {
			return 0
		}

		res, err := svc.SearchProduct(ctx, "10", "Bulunamayan")
		assert.NoError(t, err)
		assert.Empty(t, res)
	})

	t.Run("sonuç varsa listeyi dönmeli", func(t *testing.T) {
		prodRepo.searchProductCountFn = func(ctx context.Context, searchQuery string) int {
			return 1
		}
		prodRepo.searchProductFn = func(ctx context.Context, searchQuery string, limit string) ([]st.Getproducts, error) {
			return []st.Getproducts{
				{ProductID: "p-1", ProductName: "Laptop", ThumbnailImage: "thumb.jpg"},
			}, nil
		}

		res, err := svc.SearchProduct(ctx, "10", "Laptop")
		assert.NoError(t, err)
		assert.Len(t, res, 1)
		assert.Contains(t, res[0].ThumbnailImage, "test-bucket")
	})
}

func TestProductService_SearchProductAdmin(t *testing.T) {
	svc, prodRepo, _, ppcRepo, _, _ := setupProductServiceTest()
	ctx := context.Background()

	t.Run("admin arama sonucu 0 ise boş liste dönmeli", func(t *testing.T) {
		prodRepo.searchProductCountAdminFn = func(ctx context.Context, searchQuery string) int {
			return 0
		}

		res, err := svc.SearchProductAdmin(ctx, "10", "Yok")
		assert.NoError(t, err)
		assert.Empty(t, res)
	})

	t.Run("admin arama sonuçlarını kategorileriyle birlikte dönmeli", func(t *testing.T) {
		prodRepo.searchProductCountAdminFn = func(ctx context.Context, searchQuery string) int {
			return 1
		}
		prodRepo.searchProductAdminFn = func(ctx context.Context, searchQuery string, limit string) ([]st.GetproductsAdmin, error) {
			return []st.GetproductsAdmin{
				{ProductID: "p-1", ProductName: "Telefon", ThumbnailImage: "phone.jpg"},
			}, nil
		}
		ppcRepo.getAllProductCategoriesFn = func(ctx context.Context, ids []string) ([]*structs.GetAllCategory, error) {
			return []*structs.GetAllCategory{
				{ProductId: "p-1", Categories: []structs.CategoryDbJson{{CategoryName: "Mobil"}}},
			}, nil
		}

		res, err := svc.SearchProductAdmin(ctx, "10", "Telefon")
		assert.NoError(t, err)
		assert.Len(t, res, 1)
		assert.Equal(t, "Telefon", res[0].Product.ProductName)
		assert.Contains(t, res[0].Product.ThumbnailImage, "test-bucket")
	})
}
