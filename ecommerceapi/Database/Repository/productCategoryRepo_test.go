package repository

import (
	"context"
	structs "ecommerceapi/Structs"
	"fmt"
	"testing"

	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

// createTestCategoryHelper testler için izole bir kategori oluşturur
func createTestCategoryHelper(t *testing.T, ctx context.Context, repo *ProductCategoryTable, namePrefix string) structs.GetCategory {
	t.Helper()

	categoryName := fmt.Sprintf("%s_%s", namePrefix, uuid.NewString()[:8])
	category, err := repo.Create(ctx, categoryName)
	require.NoError(t, err)
	require.NotEmpty(t, category.Category_id)
	require.Equal(t, categoryName, category.Category_name)

	return category
}

func TestGetAllCategories(t *testing.T) {
	categoryRepo := NewProductCategoryStore(testDbInstance)

	t.Run("tüm kategorileri başarıyla listelemeli", func(t *testing.T) {
		// arrange
		ctx := context.Background()
		cat1 := createTestCategoryHelper(t, ctx, categoryRepo, "GetAllCat1")
		cat2 := createTestCategoryHelper(t, ctx, categoryRepo, "GetAllCat2")

		// act
		categories, err := categoryRepo.GetAll(ctx)

		// assert
		assert.NoError(t, err)
		assert.NotEmpty(t, categories)

		found1, found2 := false, false
		for _, c := range categories {
			if c.Category_id == cat1.Category_id && c.Category_name == cat1.Category_name {
				found1 = true
			}
			if c.Category_id == cat2.Category_id && c.Category_name == cat2.Category_name {
				found2 = true
			}
		}
		assert.True(t, found1, "cat1 kategorisi listede bulunmalı")
		assert.True(t, found2, "cat2 kategorisi listede bulunmalı")
	})

	t.Run("bağlam iptal edildiğinde hata dönmeli", func(t *testing.T) {
		// arrange
		ctx, cancel := context.WithCancel(context.Background())
		cancel()

		// act
		categories, err := categoryRepo.GetAll(ctx)

		// assert
		assert.Error(t, err)
		assert.Nil(t, categories)
		assert.Contains(t, err.Error(), "ürün kategori sorgu hatası")
	})
}

func TestCreateCategory(t *testing.T) {
	categoryRepo := NewProductCategoryStore(testDbInstance)

	t.Run("başarıyla yeni kategori oluşturmalı", func(t *testing.T) {
		// arrange
		ctx := context.Background()
		categoryName := fmt.Sprintf("YeniKategori_%s", uuid.NewString()[:8])

		// act
		created, err := categoryRepo.Create(ctx, categoryName)

		// assert
		assert.NoError(t, err)
		assert.NotEmpty(t, created.Category_id)
		assert.Equal(t, categoryName, created.Category_name)

		parsedUUID, err := uuid.Parse(created.Category_id)
		assert.NoError(t, err, "Category_id geçerli bir UUID olmalı")

		// DB doğrulama
		var dbName string
		query := `SELECT category_name FROM productcategories WHERE category_id = $1`
		err = testDbInstance.QueryRow(ctx, query, parsedUUID).Scan(&dbName)
		assert.NoError(t, err)
		assert.Equal(t, categoryName, dbName)
	})

	t.Run("bağlam iptal edildiğinde hata dönmeli", func(t *testing.T) {
		// arrange
		ctx, cancel := context.WithCancel(context.Background())
		cancel()

		// act
		created, err := categoryRepo.Create(ctx, "CancelKategori")

		// assert
		assert.Error(t, err)
		assert.Empty(t, created.Category_id)
		assert.Contains(t, err.Error(), "ürün kategori oluşturma hatası")
	})
}

func TestDeleteCategory(t *testing.T) {
	categoryRepo := NewProductCategoryStore(testDbInstance)

	t.Run("kategoriyi başarıyla silmeli", func(t *testing.T) {
		// arrange
		ctx := context.Background()
		cat := createTestCategoryHelper(t, ctx, categoryRepo, "SilinecekCat")

		// act
		err := categoryRepo.Delete(ctx, cat.Category_id)

		// assert
		assert.NoError(t, err)

		// DB doğrulama: satır silinmiş olmalı
		var count int
		query := `SELECT count(*) FROM productcategories WHERE category_id = $1`
		err = testDbInstance.QueryRow(ctx, query, cat.Category_id).Scan(&count)
		assert.NoError(t, err)
		assert.Equal(t, 0, count, "Kategori veritabanından silinmiş olmalı")
	})

	t.Run("olmayan kategori silinmek istendiğinde hata dönmemeli", func(t *testing.T) {
		// arrange
		ctx := context.Background()
		randomID := uuid.NewString()

		// act
		err := categoryRepo.Delete(ctx, randomID)

		// assert
		assert.NoError(t, err)
	})

	t.Run("bağlam iptal edildiğinde hata dönmeli", func(t *testing.T) {
		// arrange
		ctx, cancel := context.WithCancel(context.Background())
		cancel()

		// act
		err := categoryRepo.Delete(ctx, uuid.NewString())

		// assert
		assert.Error(t, err)
		assert.Contains(t, err.Error(), "ürün kategori silme hatası")
	})
}

func TestGetByCategoryName(t *testing.T) {
	categoryRepo := NewProductCategoryStore(testDbInstance)

	t.Run("kategori ismi ile kategori ID'sini başarıyla getirmeli", func(t *testing.T) {
		// arrange
		ctx := context.Background()
		cat := createTestCategoryHelper(t, ctx, categoryRepo, "BulunacakCat")

		// act
		foundUUID, err := categoryRepo.GetByCategoryName(ctx, cat.Category_name)

		// assert
		assert.NoError(t, err)
		assert.Equal(t, cat.Category_id, foundUUID.String())
	})

	t.Run("olmayan kategori ismi sorgulandığında hata dönmeli", func(t *testing.T) {
		// arrange
		ctx := context.Background()
		nonExistentName := fmt.Sprintf("non_existent_%s", uuid.NewString())

		// act
		foundUUID, err := categoryRepo.GetByCategoryName(ctx, nonExistentName)

		// assert
		assert.Error(t, err)
		assert.Equal(t, uuid.Nil, foundUUID)
		assert.Contains(t, err.Error(), "select sorgusunda hatayla karsilasildi")
	})

	t.Run("bağlam iptal edildiğinde hata dönmeli", func(t *testing.T) {
		// arrange
		ctx, cancel := context.WithCancel(context.Background())
		cancel()

		// act
		foundUUID, err := categoryRepo.GetByCategoryName(ctx, "herhangi_cat")

		// assert
		assert.Error(t, err)
		assert.Equal(t, uuid.Nil, foundUUID)
	})
}

func TestGetCategoriesByProductId(t *testing.T) {
	categoryRepo := NewProductCategoryStore(testDbInstance)

	t.Run("ürüne bağlı kategorileri başarıyla listelemeli", func(t *testing.T) {
		// arrange
		ctx := context.Background()
		productID := createTestProduct(t, ctx, testDbInstance)
		cat1 := createTestCategoryHelper(t, ctx, categoryRepo, "ProdCat1")
		cat2 := createTestCategoryHelper(t, ctx, categoryRepo, "ProdCat2")

		// Ürün ile kategorileri eşleştir
		insertQuery := `INSERT INTO productproductcategories (product_id, category_id) VALUES ($1, $2), ($1, $3)`
		_, err := testDbInstance.Exec(ctx, insertQuery, productID, cat1.Category_id, cat2.Category_id)
		require.NoError(t, err)

		// act
		categories, err := categoryRepo.GetCategoriesByProductId(ctx, productID)

		// assert
		assert.NoError(t, err)
		assert.Len(t, categories, 2)

		catIDs := make([]string, len(categories))
		for i, c := range categories {
			catIDs[i] = c.Category_id
		}
		assert.Contains(t, catIDs, cat1.Category_id)
		assert.Contains(t, catIDs, cat2.Category_id)
	})

	t.Run("hiç kategorisi olmayan ürün için boş liste dönmeli", func(t *testing.T) {
		// arrange
		ctx := context.Background()
		productID := createTestProduct(t, ctx, testDbInstance)

		// act
		categories, err := categoryRepo.GetCategoriesByProductId(ctx, productID)

		// assert
		assert.NoError(t, err)
		assert.Empty(t, categories)
	})

	t.Run("bağlam iptal edildiğinde hata dönmeli", func(t *testing.T) {
		// arrange
		ctx, cancel := context.WithCancel(context.Background())
		cancel()

		// act
		categories, err := categoryRepo.GetCategoriesByProductId(ctx, uuid.NewString())

		// assert
		assert.Error(t, err)
		assert.Nil(t, categories)
	})
}

func TestGetByCategoryIds(t *testing.T) {
	categoryRepo := NewProductCategoryStore(testDbInstance)

	t.Run("kategori isim listesine göre kategori ID'lerini getirmeli", func(t *testing.T) {
		// arrange
		ctx := context.Background()
		cat1 := createTestCategoryHelper(t, ctx, categoryRepo, "ByNames1")
		cat2 := createTestCategoryHelper(t, ctx, categoryRepo, "ByNames2")
		cat3 := createTestCategoryHelper(t, ctx, categoryRepo, "ByNames3")

		namesToQuery := []string{cat1.Category_name, cat2.Category_name, cat3.Category_name}

		// act
		ids, err := categoryRepo.GetByCategoryIds(ctx, namesToQuery)

		// assert
		assert.NoError(t, err)
		assert.Len(t, ids, 3)
		assert.Contains(t, ids, cat1.Category_id)
		assert.Contains(t, ids, cat2.Category_id)
		assert.Contains(t, ids, cat3.Category_id)
	})

	t.Run("eşleşmeyen isimler verildiğinde boş liste dönmeli", func(t *testing.T) {
		// arrange
		ctx := context.Background()
		nonExistentNames := []string{
			fmt.Sprintf("not_exist_%s", uuid.NewString()),
			fmt.Sprintf("not_exist_%s", uuid.NewString()),
		}

		// act
		ids, err := categoryRepo.GetByCategoryIds(ctx, nonExistentNames)

		// assert
		assert.NoError(t, err)
		assert.Empty(t, ids)
	})

	t.Run("boş isim listesi verildiğinde boş liste dönmeli", func(t *testing.T) {
		// arrange
		ctx := context.Background()

		// act
		ids, err := categoryRepo.GetByCategoryIds(ctx, []string{})

		// assert
		assert.NoError(t, err)
		assert.Empty(t, ids)
	})

	t.Run("bağlam iptal edildiğinde hata dönmeli", func(t *testing.T) {
		// arrange
		ctx, cancel := context.WithCancel(context.Background())
		cancel()

		// act
		ids, err := categoryRepo.GetByCategoryIds(ctx, []string{"any_category"})

		// assert
		assert.Error(t, err)
		assert.Nil(t, ids)
	})
}

func TestUpdateCategory(t *testing.T) {
	categoryRepo := NewProductCategoryStore(testDbInstance)

	t.Run("kategori ismini başarıyla güncellemeli", func(t *testing.T) {
		// arrange
		ctx := context.Background()
		cat := createTestCategoryHelper(t, ctx, categoryRepo, "EskiIsim")
		newCategoryName := fmt.Sprintf("GuncelIsim_%s", uuid.NewString()[:8])

		// act
		err := categoryRepo.Update(ctx, cat.Category_id, newCategoryName)

		// assert
		assert.NoError(t, err)

		// DB doğrulama
		var dbName string
		query := `SELECT category_name FROM productcategories WHERE category_id = $1`
		err = testDbInstance.QueryRow(ctx, query, cat.Category_id).Scan(&dbName)
		assert.NoError(t, err)
		assert.Equal(t, newCategoryName, dbName)
	})

	t.Run("olmayan kategori güncellenmek istendiğinde hata dönmemeli", func(t *testing.T) {
		// arrange
		ctx := context.Background()
		randomID := uuid.NewString()

		// act
		err := categoryRepo.Update(ctx, randomID, "Yeni Isim")

		// assert
		assert.NoError(t, err)
	})

	t.Run("bağlam iptal edildiğinde hata dönmeli", func(t *testing.T) {
		// arrange
		ctx, cancel := context.WithCancel(context.Background())
		cancel()

		// act
		err := categoryRepo.Update(ctx, uuid.NewString(), "Iptal Isim")

		// assert
		assert.Error(t, err)
	})
}
