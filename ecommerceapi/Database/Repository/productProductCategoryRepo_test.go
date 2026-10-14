package repository

import (
	"context"
	Entities "ecommerceapi/Database/Models"
	"fmt"
	"testing"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

// createCategoryForTest testler için izole bir kategori oluşturur
func createCategoryForTest(t *testing.T, ctx context.Context, db *pgxpool.Pool) (string, string) {
	t.Helper()

	catID := uuid.NewString()
	name := fmt.Sprintf("Cat_%s", uuid.NewString()[:8])
	query := `INSERT INTO productcategories (category_id, category_name) VALUES ($1, $2)`
	_, err := db.Exec(ctx, query, catID, name)
	if err != nil {
		t.Fatalf("Test kategorisi oluşturulamadı: %v", err)
	}

	return catID, name
}

func TestAddProductProductCategory(t *testing.T) {
	ppcRepo := NewProductProductCategoryStore(testDbInstance)

	t.Run("başarıyla ürün-kategori ilişkisi eklemeli", func(t *testing.T) {
		// arrange
		ctx := context.Background()
		productID := createTestProduct(t, ctx, testDbInstance)
		categoryID, _ := createCategoryForTest(t, ctx, testDbInstance)

		item := &Entities.ProductProductCategory{
			ProductID:  productID,
			CategoryID: categoryID,
		}

		// act
		err := ppcRepo.Add(ctx, item)

		// assert
		assert.NoError(t, err)

		// DB doğrulama
		var count int
		query := `SELECT count(*) FROM productproductcategories WHERE product_id = $1 AND category_id = $2`
		err = testDbInstance.QueryRow(ctx, query, productID, categoryID).Scan(&count)
		assert.NoError(t, err)
		assert.Equal(t, 1, count)
	})

	t.Run("olmayan product_id ile eklenmeye çalışıldığında foreign key hatası dönmeli", func(t *testing.T) {
		// arrange
		ctx := context.Background()
		categoryID, _ := createCategoryForTest(t, ctx, testDbInstance)
		nonExistentProductID := uuid.NewString()

		item := &Entities.ProductProductCategory{
			ProductID:  nonExistentProductID,
			CategoryID: categoryID,
		}

		// act
		err := ppcRepo.Add(ctx, item)

		// assert
		assert.Error(t, err)
		assert.Contains(t, err.Error(), "ürün kategori ilişkisi oluşturma hatası")
	})

	t.Run("olmayan category_id ile eklenmeye çalışıldığında foreign key hatası dönmeli", func(t *testing.T) {
		// arrange
		ctx := context.Background()
		productID := createTestProduct(t, ctx, testDbInstance)
		nonExistentCategoryID := uuid.NewString()

		item := &Entities.ProductProductCategory{
			ProductID:  productID,
			CategoryID: nonExistentCategoryID,
		}

		// act
		err := ppcRepo.Add(ctx, item)

		// assert
		assert.Error(t, err)
		assert.Contains(t, err.Error(), "ürün kategori ilişkisi oluşturma hatası")
	})

	t.Run("bağlam iptal edildiğinde (context canceled) hata dönmeli", func(t *testing.T) {
		// arrange
		ctx, cancel := context.WithCancel(context.Background())
		cancel()

		item := &Entities.ProductProductCategory{
			ProductID:  uuid.NewString(),
			CategoryID: uuid.NewString(),
		}

		// act
		err := ppcRepo.Add(ctx, item)

		// assert
		assert.Error(t, err)
	})
}

func TestGetByProductID(t *testing.T) {
	ppcRepo := NewProductProductCategoryStore(testDbInstance)

	t.Run("ürün ID ile tüm kategori ilişkilerini başarıyla getirmeli", func(t *testing.T) {
		// arrange
		ctx := context.Background()
		productID := createTestProduct(t, ctx, testDbInstance)
		productUUID, err := uuid.Parse(productID)
		require.NoError(t, err)

		cat1ID, _ := createCategoryForTest(t, ctx, testDbInstance)
		cat2ID, _ := createCategoryForTest(t, ctx, testDbInstance)

		require.NoError(t, ppcRepo.Add(ctx, &Entities.ProductProductCategory{ProductID: productID, CategoryID: cat1ID}))
		require.NoError(t, ppcRepo.Add(ctx, &Entities.ProductProductCategory{ProductID: productID, CategoryID: cat2ID}))

		// act
		items, err := ppcRepo.GetByProductID(ctx, productUUID)

		// assert
		assert.NoError(t, err)
		assert.Len(t, items, 2)

		catIDs := []string{items[0].CategoryID, items[1].CategoryID}
		assert.Contains(t, catIDs, cat1ID)
		assert.Contains(t, catIDs, cat2ID)
	})

	t.Run("kategorisi olmayan ürün için boş liste dönmeli", func(t *testing.T) {
		// arrange
		ctx := context.Background()
		productID := createTestProduct(t, ctx, testDbInstance)
		productUUID, err := uuid.Parse(productID)
		require.NoError(t, err)

		// act
		items, err := ppcRepo.GetByProductID(ctx, productUUID)

		// assert
		assert.NoError(t, err)
		assert.Empty(t, items)
	})

	t.Run("bağlam iptal edildiğinde hata dönmeli", func(t *testing.T) {
		// arrange
		ctx, cancel := context.WithCancel(context.Background())
		cancel()

		// act
		items, err := ppcRepo.GetByProductID(ctx, uuid.New())

		// assert
		assert.Error(t, err)
		assert.Nil(t, items)
		assert.Contains(t, err.Error(), "ürün kategori ilişkisi sorgu hatası")
	})
}

func TestGetByCategoryID(t *testing.T) {
	ppcRepo := NewProductProductCategoryStore(testDbInstance)

	t.Run("kategori ID ile tüm ürün ilişkilerini başarıyla getirmeli", func(t *testing.T) {
		// arrange
		ctx := context.Background()
		catID, _ := createCategoryForTest(t, ctx, testDbInstance)
		catUUID, err := uuid.Parse(catID)
		require.NoError(t, err)

		p1ID := createTestProduct(t, ctx, testDbInstance)
		p2ID := createTestProduct(t, ctx, testDbInstance)

		require.NoError(t, ppcRepo.Add(ctx, &Entities.ProductProductCategory{ProductID: p1ID, CategoryID: catID}))
		require.NoError(t, ppcRepo.Add(ctx, &Entities.ProductProductCategory{ProductID: p2ID, CategoryID: catID}))

		// act
		items, err := ppcRepo.GetByCategoryID(ctx, catUUID)

		// assert
		assert.NoError(t, err)
		assert.Len(t, items, 2)

		pIDs := []string{items[0].ProductID, items[1].ProductID}
		assert.Contains(t, pIDs, p1ID)
		assert.Contains(t, pIDs, p2ID)
	})

	t.Run("ürünü olmayan kategori için boş liste dönmeli", func(t *testing.T) {
		// arrange
		ctx := context.Background()
		catID, _ := createCategoryForTest(t, ctx, testDbInstance)
		catUUID, err := uuid.Parse(catID)
		require.NoError(t, err)

		// act
		items, err := ppcRepo.GetByCategoryID(ctx, catUUID)

		// assert
		assert.NoError(t, err)
		assert.Empty(t, items)
	})

	t.Run("bağlam iptal edildiğinde hata dönmeli", func(t *testing.T) {
		// arrange
		ctx, cancel := context.WithCancel(context.Background())
		cancel()

		// act
		items, err := ppcRepo.GetByCategoryID(ctx, uuid.New())

		// assert
		assert.Error(t, err)
		assert.Nil(t, items)
		assert.Contains(t, err.Error(), "ürün kategori ilişkisi sorgu hatası")
	})
}

func TestRemoveProductProductCategory(t *testing.T) {
	ppcRepo := NewProductProductCategoryStore(testDbInstance)

	t.Run("ürün-kategori ilişkisini başarıyla silmeli", func(t *testing.T) {
		// arrange
		ctx := context.Background()
		productID := createTestProduct(t, ctx, testDbInstance)
		catID, _ := createCategoryForTest(t, ctx, testDbInstance)

		productUUID, err := uuid.Parse(productID)
		require.NoError(t, err)
		catUUID, err := uuid.Parse(catID)
		require.NoError(t, err)

		require.NoError(t, ppcRepo.Add(ctx, &Entities.ProductProductCategory{ProductID: productID, CategoryID: catID}))

		// act
		err = ppcRepo.Remove(ctx, productUUID, catUUID)

		// assert
		assert.NoError(t, err)

		// DB doğrulama: satır silinmiş olmalı
		var count int
		query := `SELECT count(*) FROM productproductcategories WHERE product_id = $1 AND category_id = $2`
		err = testDbInstance.QueryRow(ctx, query, productUUID, catUUID).Scan(&count)
		assert.NoError(t, err)
		assert.Equal(t, 0, count)
	})

	t.Run("olmayan ilişki silinmek istendiğinde hata dönmemeli", func(t *testing.T) {
		// arrange
		ctx := context.Background()

		// act
		err := ppcRepo.Remove(ctx, uuid.New(), uuid.New())

		// assert
		assert.NoError(t, err)
	})

	t.Run("bağlam iptal edildiğinde hata dönmeli", func(t *testing.T) {
		// arrange
		ctx, cancel := context.WithCancel(context.Background())
		cancel()

		// act
		err := ppcRepo.Remove(ctx, uuid.New(), uuid.New())

		// assert
		assert.Error(t, err)
		assert.Contains(t, err.Error(), "ürün kategori ilişkisi silme hatası")
	})
}

func TestGetByProductIDOneRow(t *testing.T) {
	ppcRepo := NewProductProductCategoryStore(testDbInstance)

	t.Run("ürüne ait ilk kategori ID'sini başarıyla getirmeli", func(t *testing.T) {
		// arrange
		ctx := context.Background()
		productID := createTestProduct(t, ctx, testDbInstance)
		catID, _ := createCategoryForTest(t, ctx, testDbInstance)

		productUUID, err := uuid.Parse(productID)
		require.NoError(t, err)
		catUUID, err := uuid.Parse(catID)
		require.NoError(t, err)

		require.NoError(t, ppcRepo.Add(ctx, &Entities.ProductProductCategory{ProductID: productID, CategoryID: catID}))

		// act
		foundCatUUID, err := ppcRepo.GetByProductIDOneRow(ctx, productUUID)

		// assert
		assert.NoError(t, err)
		assert.Equal(t, catUUID, foundCatUUID)
	})

	t.Run("ilişkisi olmayan ürün için hata dönmeli", func(t *testing.T) {
		// arrange
		ctx := context.Background()
		productID := createTestProduct(t, ctx, testDbInstance)
		productUUID, err := uuid.Parse(productID)
		require.NoError(t, err)

		// act
		foundCatUUID, err := ppcRepo.GetByProductIDOneRow(ctx, productUUID)

		// assert
		assert.Error(t, err)
		assert.Equal(t, uuid.Nil, foundCatUUID)
		assert.Contains(t, err.Error(), "category id çekilemedi")
	})

	t.Run("bağlam iptal edildiğinde hata dönmeli", func(t *testing.T) {
		// arrange
		ctx, cancel := context.WithCancel(context.Background())
		cancel()

		// act
		foundCatUUID, err := ppcRepo.GetByProductIDOneRow(ctx, uuid.New())

		// assert
		assert.Error(t, err)
		assert.Equal(t, uuid.Nil, foundCatUUID)
	})
}

func TestAddManyProductProductCategory(t *testing.T) {
	ppcRepo := NewProductProductCategoryStore(testDbInstance)

	t.Run("toplu olarak ürün-kategori ilişkilerini eklemeli", func(t *testing.T) {
		// arrange
		ctx := context.Background()
		productID := createTestProduct(t, ctx, testDbInstance)

		cat1ID, _ := createCategoryForTest(t, ctx, testDbInstance)
		cat2ID, _ := createCategoryForTest(t, ctx, testDbInstance)
		cat3ID, _ := createCategoryForTest(t, ctx, testDbInstance)

		items := []string{cat1ID, cat2ID, cat3ID}

		// act
		err := ppcRepo.AddMany(ctx, items, productID)

		// assert
		assert.NoError(t, err)

		// DB doğrulama: 3 ilişki de eklenmiş olmalı
		var count int
		query := `SELECT count(*) FROM productproductcategories WHERE product_id = $1`
		err = testDbInstance.QueryRow(ctx, query, productID).Scan(&count)
		assert.NoError(t, err)
		assert.Equal(t, 3, count)
	})

	t.Run("boş kategori listesi verildiğinde hata dönmemeli", func(t *testing.T) {
		// arrange
		ctx := context.Background()
		productID := createTestProduct(t, ctx, testDbInstance)

		// act
		err := ppcRepo.AddMany(ctx, []string{}, productID)

		// assert
		assert.NoError(t, err)
	})

	t.Run("geçersiz UUID verildiğinde hata dönmeli", func(t *testing.T) {
		// arrange
		ctx := context.Background()
		productID := createTestProduct(t, ctx, testDbInstance)

		// act
		err := ppcRepo.AddMany(ctx, []string{"invalid-uuid"}, productID)

		// assert
		assert.Error(t, err)
		assert.Contains(t, err.Error(), "db error")
	})

	t.Run("bağlam iptal edildiğinde hata dönmeli", func(t *testing.T) {
		// arrange
		ctx, cancel := context.WithCancel(context.Background())
		cancel()

		// act
		err := ppcRepo.AddMany(ctx, []string{uuid.NewString()}, uuid.NewString())

		// assert
		assert.Error(t, err)
	})
}

func TestRemoveManyProductProductCategory(t *testing.T) {
	ppcRepo := NewProductProductCategoryStore(testDbInstance)

	t.Run("belirtilen kategorileri toplu olarak silmeli", func(t *testing.T) {
		// arrange
		ctx := context.Background()
		productID := createTestProduct(t, ctx, testDbInstance)

		cat1ID, _ := createCategoryForTest(t, ctx, testDbInstance)
		cat2ID, _ := createCategoryForTest(t, ctx, testDbInstance)
		cat3ID, _ := createCategoryForTest(t, ctx, testDbInstance)

		// 3 ilişki ekle
		require.NoError(t, ppcRepo.AddMany(ctx, []string{cat1ID, cat2ID, cat3ID}, productID))

		// act: cat1 ve cat2'yi sil
		err := ppcRepo.RemoveMany(ctx, []string{cat1ID, cat2ID}, productID)

		// assert
		assert.NoError(t, err)

		// DB doğrulama: geriye sadece cat3 kalmalı
		var count int
		err = testDbInstance.QueryRow(ctx, "SELECT count(*) FROM productproductcategories WHERE product_id = $1", productID).Scan(&count)
		assert.NoError(t, err)
		assert.Equal(t, 1, count)

		var remainingCatID string
		err = testDbInstance.QueryRow(ctx, "SELECT category_id FROM productproductcategories WHERE product_id = $1", productID).Scan(&remainingCatID)
		assert.NoError(t, err)
		assert.Equal(t, cat3ID, remainingCatID)
	})

	t.Run("boş liste veya olmayan kategoriler için hata dönmemeli", func(t *testing.T) {
		// arrange
		ctx := context.Background()
		productID := createTestProduct(t, ctx, testDbInstance)

		// act
		err := ppcRepo.RemoveMany(ctx, []string{}, productID)

		// assert
		assert.NoError(t, err)
	})

	t.Run("bağlam iptal edildiğinde hata dönmeli", func(t *testing.T) {
		// arrange
		ctx, cancel := context.WithCancel(context.Background())
		cancel()

		// act
		err := ppcRepo.RemoveMany(ctx, []string{uuid.NewString()}, uuid.NewString())

		// assert
		assert.Error(t, err)
	})
}

func TestGetAllProductCategories(t *testing.T) {
	ppcRepo := NewProductProductCategoryStore(testDbInstance)

	t.Run("verilen ürün ID listesine göre kategorileri JSON formatında getirmeli", func(t *testing.T) {
		// arrange
		ctx := context.Background()
		p1 := createTestProduct(t, ctx, testDbInstance)
		p2 := createTestProduct(t, ctx, testDbInstance)

		catA_ID, catA_Name := createCategoryForTest(t, ctx, testDbInstance)
		catB_ID, catB_Name := createCategoryForTest(t, ctx, testDbInstance)
		catC_ID, catC_Name := createCategoryForTest(t, ctx, testDbInstance)

		// p1 -> CatA, CatB
		require.NoError(t, ppcRepo.AddMany(ctx, []string{catA_ID, catB_ID}, p1))
		// p2 -> CatC
		require.NoError(t, ppcRepo.AddMany(ctx, []string{catC_ID}, p2))

		// act
		results, err := ppcRepo.GetAllProductCategories(ctx, []string{p1, p2})

		// assert
		assert.NoError(t, err)
		assert.Len(t, results, 2)

		p1ResultFound := false
		p2ResultFound := false

		for _, r := range results {
			if r.ProductId == p1 {
				p1ResultFound = true
				assert.Len(t, r.Categories, 2)
				catNames := []string{r.Categories[0].CategoryName, r.Categories[1].CategoryName}
				assert.Contains(t, catNames, catA_Name)
				assert.Contains(t, catNames, catB_Name)
			}
			if r.ProductId == p2 {
				p2ResultFound = true
				assert.Len(t, r.Categories, 1)
				assert.Equal(t, catC_Name, r.Categories[0].CategoryName)
				assert.Equal(t, catC_ID, r.Categories[0].CategoryId)
			}
		}

		assert.True(t, p1ResultFound, "p1 sonuçlarda yer almalı")
		assert.True(t, p2ResultFound, "p2 sonuçlarda yer almalı")
	})

	t.Run("eşleşmeyen ürün ID'leri için boş liste dönmeli", func(t *testing.T) {
		// arrange
		ctx := context.Background()
		nonExistentID := uuid.NewString()

		// act
		results, err := ppcRepo.GetAllProductCategories(ctx, []string{nonExistentID})

		// assert
		assert.NoError(t, err)
		assert.Empty(t, results)
	})

	t.Run("boş ID listesi için boş liste dönmeli", func(t *testing.T) {
		// arrange
		ctx := context.Background()

		// act
		results, err := ppcRepo.GetAllProductCategories(ctx, []string{})

		// assert
		assert.NoError(t, err)
		assert.Empty(t, results)
	})

	t.Run("bağlam iptal edildiğinde hata dönmeli", func(t *testing.T) {
		// arrange
		ctx, cancel := context.WithCancel(context.Background())
		cancel()

		// act
		results, err := ppcRepo.GetAllProductCategories(ctx, []string{uuid.NewString()})

		// assert
		assert.Error(t, err)
		assert.Nil(t, results)
	})
}
