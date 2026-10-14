package repository

import (
	"context"
	constants "ecommerceapi/Constants"
	models "ecommerceapi/Database/Models"
	"ecommerceapi/Utils"
	"fmt"
	"strings"
	"testing"
	"time"

	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
)

func TestCreateProduct(t *testing.T) {

	productRepo := NewProductStore(testDbInstance)

	t.Run("başarılı ürün oluşturma", func(t *testing.T) {

		// arrange

		productID := uuid.NewString()
		thumbnailImage := uuid.NewString()
		now := time.Now()

		product := models.Product{
			ProductID:           productID,
			ProductName:         "Basket Topu",
			Stock:               10,
			Price:               165000,
			ThumbnailImage:      thumbnailImage,
			Description:         "Gerçek basketbol topu açıklaması",
			CreatedAt:           &now,
			UpdatedAt:           &now,
			ItemType:            "PHYSICAL",
			ProductCode:         "ABCDEF12",
			ProductPublishState: "published",
		}
		ctx := context.Background()
		// act
		err := productRepo.Create(ctx, &product)

		// assert
		assert.NoErrorf(t, err, "Ürün Oluşturulamadı: %v", err)
	})

	t.Run("boş publish state gönderildiğinde varsayılan değer atanmalı ve zaman damgaları oluşmalı", func(t *testing.T) {
		// arrange
		product := models.Product{
			ProductID:   uuid.NewString(),
			ProductName: "State Test Ürünü",
			Price:       100,
			// ProductPublishState bilerek boş bırakıldı
			// CreatedAt ve UpdatedAt bilerek nil bırakıldı
		}
		ctx := context.Background()

		// act
		err := productRepo.Create(ctx, &product)

		// assert
		assert.NoError(t, err)
		assert.Equal(t, constants.PublishStateUnpublished, product.ProductPublishState)
		assert.NotNil(t, product.CreatedAt)
		assert.NotNil(t, product.UpdatedAt)
	})

	t.Run("veritabanı hatası oluştuğunda beklenen formatta hata dönmeli (Duplicate ID)", func(t *testing.T) {
		// arrange
		productID := uuid.NewString()
		product := models.Product{
			ProductID:           productID,
			ProductName:         "Hata Test Ürünü",
			ProductPublishState: "published",
			ProductCode:         "ABCDEF23",
		}
		ctx := context.Background()

		// İlk kaydı başarılı şekilde at
		err1 := productRepo.Create(ctx, &product)
		assert.NoError(t, err1)

		// act
		// Aynı ProductID ile ikinci kez kayıt atmaya çalış (Primary Key kısıtlamasına takılacak)
		err2 := productRepo.Create(ctx, &product)

		// assert
		assert.Error(t, err2)
		assert.Contains(t, err2.Error(), "ürün oluşturma hatası")
	})

}

func TestDeleteProduct(t *testing.T) {

	productRepo := NewProductStore(testDbInstance)

	t.Run("başarılı silme işlemi (var olan kayıt)", func(t *testing.T) {
		// arrange: Silinecek veriyi önce oluştur
		id := uuid.New()
		product := models.Product{
			ProductID:           id.String(),
			ProductName:         "Silinecek Ürün",
			ProductCode:         "ABCDEF34",
			ProductPublishState: "published",
		}
		ctx := context.Background()

		err := productRepo.Create(ctx, &product)
		assert.NoError(t, err)

		// act
		err = productRepo.Delete(ctx, &id)

		// assert
		assert.NoError(t, err)
	})

	t.Run("olmayan bir ürünü silme denemesi veritabanı hatası fırlatmaz", func(t *testing.T) {
		// arrange: Veritabanında olmayan rastgele bir UUID
		id := uuid.New()
		ctx := context.Background()

		// act
		err := productRepo.Delete(ctx, &id)

		// assert
		// Not: SQL'de DELETE işlemi eşleşen kayıt bulamazsa hata fırlatmaz, sadece etkilenen satır sayısını 0 döner.
		// Kodda RowsAffected kontrolü yapılmadığı için bu işlemin nil dönmesi beklenir.
		assert.NoError(t, err)
	})

	t.Run("veritabanı hatası durumunda beklenen formatta hata dönmeli (Context Canceled)", func(t *testing.T) {
		// arrange
		id := uuid.New()
		ctx, cancel := context.WithCancel(context.Background())
		cancel() // Veritabanı işlemini iptal ederek zorla hata ürettirilir

		// act
		err := productRepo.Delete(ctx, &id)

		// assert
		assert.Error(t, err)
		assert.Contains(t, err.Error(), "ürün silme hatası")
	})
}

func TestUpdateProduct(t *testing.T) {
	productRepo := NewProductStore(testDbInstance)

	t.Run("başarılı ürün güncelleme (tüm alanlar)", func(t *testing.T) {
		// arrange
		ctx := context.Background()
		productID := uuid.NewString()
		productCode, err := Utils.CreateNanoId(8)
		assert.NoError(t, err)

		initialProduct := models.Product{
			ProductID:           productID,
			ProductName:         "Eski Ürün Adı",
			Stock:               5,
			Price:               50000,
			ThumbnailImage:      "old_thumbnail.jpg",
			Description:         "Eski açıklama",
			ItemType:            "PHYSICAL",
			ProductCode:         productCode,
			ProductPublishState: constants.PublishStateUnpublished,
		}

		err = productRepo.Create(ctx, &initialProduct)
		assert.NoError(t, err)

		updatedProduct := models.Product{
			ProductID:           productID,
			ProductName:         "Güncellenmiş Ürün Adı",
			Stock:               25,
			Price:               75000,
			Description:         "Güncellenmiş detaylı açıklama",
			ItemType:            "DIGITAL",
			TechnicalDetails:    []byte(`{"color":"black","weight":200}`),
			ProductPublishState: constants.PublishStatePublished,
		}

		// act
		err = productRepo.Update(ctx, &updatedProduct)

		// assert
		assert.NoError(t, err)
		assert.NotNil(t, updatedProduct.UpdatedAt)

		// Veritabanındaki değerleri doğrudan doğrula
		var (
			name         string
			stock        int
			desc         string
			price        int64
			itemType     string
			techDetails  []byte
			publishState string
			updatedAt    *time.Time
		)
		query := `SELECT product_name, stock, description, price, item_type, technical_details, product_publish_state, updated_at FROM products WHERE product_id = $1`
		err = testDbInstance.QueryRow(ctx, query, productID).Scan(
			&name, &stock, &desc, &price, &itemType, &techDetails, &publishState, &updatedAt,
		)
		assert.NoError(t, err)
		assert.Equal(t, updatedProduct.ProductName, name)
		assert.Equal(t, updatedProduct.Stock, stock)
		assert.Equal(t, updatedProduct.Description, desc)
		assert.Equal(t, updatedProduct.Price, price)
		assert.Equal(t, updatedProduct.ItemType, itemType)
		assert.JSONEq(t, string(updatedProduct.TechnicalDetails), string(techDetails))
		assert.Equal(t, updatedProduct.ProductPublishState, publishState)
		assert.NotNil(t, updatedAt)
	})

	t.Run("olmayan bir ürünü güncelleme denemesi veritabanı hatası fırlatmaz", func(t *testing.T) {
		// arrange
		ctx := context.Background()
		nonExistentProduct := models.Product{
			ProductID:           uuid.NewString(),
			ProductName:         "Olmayan Ürün",
			Stock:               1,
			Price:               100,
			Description:         "Açıklama",
			ItemType:            "PHYSICAL",
			ProductPublishState: constants.PublishStatePublished,
		}

		// act
		err := productRepo.Update(ctx, &nonExistentProduct)

		// assert
		assert.NoError(t, err)
	})

	t.Run("geçersiz UUID formatında veritabanı hatası dönmeli", func(t *testing.T) {
		// arrange
		ctx := context.Background()
		product := models.Product{
			ProductID:   "invalid-uuid-format",
			ProductName: "Geçersiz ID'li Ürün",
		}

		// act
		err := productRepo.Update(ctx, &product)

		// assert
		assert.Error(t, err)
		assert.Contains(t, err.Error(), "güncelleme sorgusu hatası")
	})

	t.Run("veritabanı hatası durumunda beklenen formatta hata dönmeli (Context Canceled)", func(t *testing.T) {
		// arrange
		ctx, cancel := context.WithCancel(context.Background())
		cancel()

		product := models.Product{
			ProductID:   uuid.NewString(),
			ProductName: "Hata Ürünü",
		}

		// act
		err := productRepo.Update(ctx, &product)

		// assert
		assert.Error(t, err)
		assert.Contains(t, err.Error(), "güncelleme sorgusu hatası")
	})
}

func TestPatchProduct(t *testing.T) {
	productRepo := NewProductStore(testDbInstance)

	t.Run("başarılı kısmi güncelleme (Patch) ve diğer alanların korunması", func(t *testing.T) {
		// arrange
		ctx := context.Background()
		productID := uuid.NewString()
		productCode, err := Utils.CreateNanoId(8)
		assert.NoError(t, err)

		initialProduct := models.Product{
			ProductID:           productID,
			ProductName:         "Orijinal Ürün",
			Stock:               10,
			Price:               1000,
			ThumbnailImage:      "original_thumb.png",
			Description:         "Bu açıklama patch sonrası aynı kalmalıdır",
			ItemType:            "PHYSICAL",
			TechnicalDetails:    []byte(`{"material":"cotton"}`),
			ProductCode:         productCode,
			ProductPublishState: constants.PublishStatePublished,
		}

		err = productRepo.Create(ctx, &initialProduct)
		assert.NoError(t, err)

		patchData := models.Product{
			ProductID:   productID,
			ProductName: "Yamalanmış Ürün",
			Stock:       35,
			Price:       2500,
		}

		// act
		err = productRepo.Patch(ctx, &patchData)

		// assert
		assert.NoError(t, err)
		assert.NotNil(t, patchData.UpdatedAt)

		// Veritabanındaki güncellenen ve korunması gereken alanları doğrula
		var (
			name         string
			stock        int
			price        int64
			desc         string
			thumb        string
			itemType     string
			techDetails  []byte
			publishState string
			updatedAt    *time.Time
		)
		query := `SELECT product_name, stock, price, description, thumbnail_image, item_type, technical_details, product_publish_state, updated_at FROM products WHERE product_id = $1`
		err = testDbInstance.QueryRow(ctx, query, productID).Scan(
			&name, &stock, &price, &desc, &thumb, &itemType, &techDetails, &publishState, &updatedAt,
		)
		assert.NoError(t, err)

		// Güncellenen alanlar:
		assert.Equal(t, "Yamalanmış Ürün", name)
		assert.Equal(t, 35, stock)
		assert.Equal(t, int64(2500), price)
		assert.NotNil(t, updatedAt)

		// Korunması gereken alanlar (Patch tarafından değiştirilmemeli):
		assert.Equal(t, initialProduct.Description, desc)
		assert.Equal(t, initialProduct.ThumbnailImage, thumb)
		assert.Equal(t, initialProduct.ItemType, itemType)
		assert.JSONEq(t, string(initialProduct.TechnicalDetails), string(techDetails))
		assert.Equal(t, initialProduct.ProductPublishState, publishState)
	})

	t.Run("olmayan bir ürünü patch etme denemesi veritabanı hatası fırlatmaz", func(t *testing.T) {
		// arrange
		ctx := context.Background()
		nonExistentProduct := models.Product{
			ProductID:   uuid.NewString(),
			ProductName: "Olmayan Ürün",
			Stock:       5,
			Price:       100,
		}

		// act
		err := productRepo.Patch(ctx, &nonExistentProduct)

		// assert
		assert.NoError(t, err)
	})

	t.Run("geçersiz UUID formatında veritabanı hatası dönmeli", func(t *testing.T) {
		// arrange
		ctx := context.Background()
		product := models.Product{
			ProductID:   "invalid-uuid-format",
			ProductName: "Geçersiz ID'li Ürün",
			Stock:       5,
			Price:       100,
		}

		// act
		err := productRepo.Patch(ctx, &product)

		// assert
		assert.Error(t, err)
		assert.Contains(t, err.Error(), "güncelleme sorgusu hatası")
	})

	t.Run("veritabanı hatası durumunda beklenen formatta hata dönmeli (Context Canceled)", func(t *testing.T) {
		// arrange
		ctx, cancel := context.WithCancel(context.Background())
		cancel()

		product := models.Product{
			ProductID:   uuid.NewString(),
			ProductName: "Hata Ürünü",
			Stock:       5,
			Price:       100,
		}

		// act
		err := productRepo.Patch(ctx, &product)

		// assert
		assert.Error(t, err)
		assert.Contains(t, err.Error(), "güncelleme sorgusu hatası")
	})
}

func TestChangeThumbnailImage(t *testing.T) {
	productRepo := NewProductStore(testDbInstance)

	t.Run("başarılı thumbnail_image güncelleme", func(t *testing.T) {
		// arrange
		ctx := context.Background()
		productID := uuid.NewString()
		productCode, err := Utils.CreateNanoId(8)
		assert.NoError(t, err)

		product := models.Product{
			ProductID:           productID,
			ProductName:         "Thumbnail Test Ürünü",
			ThumbnailImage:      "initial_thumbnail.jpg",
			ProductCode:         productCode,
			ProductPublishState: constants.PublishStatePublished,
		}

		err = productRepo.Create(ctx, &product)
		assert.NoError(t, err)

		newThumbnailID := uuid.NewString()

		// act
		err = productRepo.ChangeThumbnailImage(ctx, newThumbnailID, productID)

		// assert
		assert.NoError(t, err)

		// Veritabanını doğrula
		var actualThumbnail string
		err = testDbInstance.QueryRow(ctx, "SELECT thumbnail_image FROM products WHERE product_id = $1", productID).Scan(&actualThumbnail)
		assert.NoError(t, err)
		assert.Equal(t, newThumbnailID, actualThumbnail)
	})

	t.Run("olmayan bir ürünün thumbnailini güncelleme denemesi veritabanı hatası fırlatmaz", func(t *testing.T) {
		// arrange
		ctx := context.Background()
		nonExistentID := uuid.NewString()

		// act
		err := productRepo.ChangeThumbnailImage(ctx, "new_thumb.png", nonExistentID)

		// assert
		assert.NoError(t, err)
	})

	t.Run("geçersiz UUID formatında veritabanı hatası dönmeli", func(t *testing.T) {
		// arrange
		ctx := context.Background()

		// act
		err := productRepo.ChangeThumbnailImage(ctx, "thumb.png", "invalid-uuid-format")

		// assert
		assert.Error(t, err)
		assert.Contains(t, err.Error(), "Beklenmedik Bir Hatayla Karşılaşıldı")
	})

	t.Run("veritabanı hatası durumunda beklenen formatta hata dönmeli (Context Canceled)", func(t *testing.T) {
		// arrange
		ctx, cancel := context.WithCancel(context.Background())
		cancel()

		// act
		err := productRepo.ChangeThumbnailImage(ctx, "thumb.png", uuid.NewString())

		// assert
		assert.Error(t, err)
		assert.Contains(t, err.Error(), "Beklenmedik Bir Hatayla Karşılaşıldı")
	})
}

func TestGetAllProducts(t *testing.T) {
	productRepo := NewProductStore(testDbInstance)

	t.Run("yalnızca yayınlanmış (published) ürünleri getirmeli ve yayınlanmamışları dahil etmemeli", func(t *testing.T) {
		// arrange
		ctx := context.Background()

		code1, err := Utils.CreateNanoId(8)
		assert.NoError(t, err)
		publishedProduct := models.Product{
			ProductID:           uuid.NewString(),
			ProductName:         "Yayınlanan Test Ürünü",
			Stock:               12,
			Price:               5500,
			ThumbnailImage:      "pub.png",
			ProductCode:         code1,
			ProductPublishState: constants.PublishStatePublished,
		}
		err = productRepo.Create(ctx, &publishedProduct)
		assert.NoError(t, err)

		code2, err := Utils.CreateNanoId(8)
		assert.NoError(t, err)
		unpublishedProduct := models.Product{
			ProductID:           uuid.NewString(),
			ProductName:         "Yayınlanmayan Test Ürünü",
			Stock:               6,
			Price:               3200,
			ThumbnailImage:      "unpub.png",
			ProductCode:         code2,
			ProductPublishState: constants.PublishStateUnpublished,
		}
		err = productRepo.Create(ctx, &unpublishedProduct)
		assert.NoError(t, err)

		// act
		products, err := productRepo.GetAll(100, 0, ctx)

		// assert
		assert.NoError(t, err)
		assert.NotEmpty(t, products)

		var foundPublished, foundUnpublished bool
		for _, p := range products {
			if p.ProductID == publishedProduct.ProductID {
				foundPublished = true
				assert.Equal(t, publishedProduct.ProductName, p.ProductName)
				assert.Equal(t, publishedProduct.Stock, p.Stock)
				assert.Equal(t, publishedProduct.Price, p.Price)
				assert.Equal(t, publishedProduct.ThumbnailImage, p.ThumbnailImage)
			}
			if p.ProductID == unpublishedProduct.ProductID {
				foundUnpublished = true
			}
		}
		assert.True(t, foundPublished, "Yayınlanmış ürün sonuçlarda bulunmalıdır")
		assert.False(t, foundUnpublished, "Yayınlanmamış ürün sonuçlarda bulunmamalıdır")
	})

	t.Run("limit ve offset parametreleri ile sayfalama doğru çalışmalı", func(t *testing.T) {
		// arrange
		ctx := context.Background()

		// act
		// limit = 1 ile en fazla 1 ürün gelmeli
		p1, err := productRepo.GetAll(1, 0, ctx)
		assert.NoError(t, err)
		assert.LessOrEqual(t, len(p1), 1)

		// offset = 1 ile sonraki kayıt gelmeli
		p2, err := productRepo.GetAll(1, 1, ctx)
		assert.NoError(t, err)

		if len(p1) > 0 && len(p2) > 0 {
			assert.NotEqual(t, p1[0].ProductID, p2[0].ProductID, "Farklı offset değerleri farklı kayıtlar getirmelidir")
		}
	})

	t.Run("yüksek offset verildiğinde boş liste dönmeli", func(t *testing.T) {
		// arrange
		ctx := context.Background()

		// act
		products, err := productRepo.GetAll(10, 999999, ctx)

		// assert
		assert.NoError(t, err)
		assert.Empty(t, products)
	})

	t.Run("veritabanı hatası durumunda beklenen formatta hata dönmeli (Context Canceled)", func(t *testing.T) {
		// arrange
		ctx, cancel := context.WithCancel(context.Background())
		cancel()

		// act
		products, err := productRepo.GetAll(10, 0, ctx)

		// assert
		assert.Error(t, err)
		assert.Nil(t, products)
		assert.Contains(t, err.Error(), "ürün sorgu hatası")
	})
}

func TestGetProductById(t *testing.T) {
	productRepo := NewProductStore(testDbInstance)

	t.Run("başarılı ürün getirme (görseller ve teknik detaylar dahil)", func(t *testing.T) {
		// arrange
		ctx := context.Background()
		productID := uuid.NewString()
		productCode, err := Utils.CreateNanoId(8)
		assert.NoError(t, err)

		techDetails := []byte(`{"screen":"AMOLED","ram":"8GB"}`)
		product := models.Product{
			ProductID:           productID,
			ProductName:         "Detay Test Telefonu",
			Stock:               15,
			Price:               3500000,
			ThumbnailImage:      "thumb_phone.jpg",
			Description:         "Kapsamlı ürün detay açıklaması",
			ItemType:            "PHYSICAL",
			TechnicalDetails:    techDetails,
			ProductCode:         productCode,
			ProductPublishState: constants.PublishStatePublished,
		}

		err = productRepo.Create(ctx, &product)
		assert.NoError(t, err)

		// Ürüne ait 2 görsel ekle
		img1ID := uuid.New()
		img2ID := uuid.New()
		_, err = testDbInstance.Exec(ctx, `INSERT INTO images (image_id, image_path, product_id) VALUES ($1, $2, $3), ($4, $5, $6)`,
			img1ID, "/img/phone1.jpg", productID,
			img2ID, "/img/phone2.jpg", productID,
		)
		assert.NoError(t, err)

		productUUID, err := uuid.Parse(productID)
		assert.NoError(t, err)

		// act
		result, err := productRepo.GetProductById(ctx, productUUID)

		// assert
		assert.NoError(t, err)
		assert.NotNil(t, result)
		assert.Equal(t, productID, result.ProductID)
		assert.Equal(t, product.ProductName, result.ProductName)
		assert.Equal(t, product.Price, result.Price)
		assert.Equal(t, product.Stock, result.Stock)
		assert.Equal(t, product.ThumbnailImage, result.ThumbnailImage)
		assert.Equal(t, product.Description, result.Description)
		assert.JSONEq(t, string(techDetails), string(result.TechnicalDetails))
		assert.Len(t, result.Images, 2)
	})

	t.Run("yayınlanmamış (unpublished) ürün getirilemez", func(t *testing.T) {
		// arrange
		ctx := context.Background()
		productID := uuid.NewString()
		productCode, err := Utils.CreateNanoId(8)
		assert.NoError(t, err)

		product := models.Product{
			ProductID:           productID,
			ProductName:         "Gizli Ürün",
			Price:               1000,
			ProductCode:         productCode,
			ProductPublishState: constants.PublishStateUnpublished,
		}

		err = productRepo.Create(ctx, &product)
		assert.NoError(t, err)

		productUUID, err := uuid.Parse(productID)
		assert.NoError(t, err)

		// act
		result, err := productRepo.GetProductById(ctx, productUUID)

		// assert
		assert.Error(t, err)
		assert.Nil(t, result)
		assert.Contains(t, err.Error(), "CollectOneRow işlemi sırasında hatayla karşılaşıldı")
	})

	t.Run("olmayan ürün ID'si ile sorgulandığında hata dönmeli", func(t *testing.T) {
		// arrange
		ctx := context.Background()
		nonExistentUUID := uuid.New()

		// act
		result, err := productRepo.GetProductById(ctx, nonExistentUUID)

		// assert
		assert.Error(t, err)
		assert.Nil(t, result)
		assert.Contains(t, err.Error(), "CollectOneRow işlemi sırasında hatayla karşılaşıldı")
	})

	t.Run("veritabanı hatası durumunda beklenen formatta hata dönmeli (Context Canceled)", func(t *testing.T) {
		// arrange
		ctx, cancel := context.WithCancel(context.Background())
		cancel()

		// act
		result, err := productRepo.GetProductById(ctx, uuid.New())

		// assert
		assert.Error(t, err)
		assert.Nil(t, result)
		assert.Contains(t, err.Error(), "Veritabanı sorgusunda hatayla karşılaşıldı")
	})
}

func TestGetMany(t *testing.T) {
	productRepo := NewProductStore(testDbInstance)

	t.Run("başarılı çoklu ürün getirme (yalnızca yayınlanmış olanlar)", func(t *testing.T) {
		// arrange
		ctx := context.Background()

		code1, err := Utils.CreateNanoId(8)
		assert.NoError(t, err)
		prod1 := models.Product{
			ProductID:           uuid.NewString(),
			ProductName:         "Sepet Ürünü 1",
			Stock:               20,
			Price:               15000,
			ThumbnailImage:      "cart_thumb_1.jpg",
			ProductCode:         code1,
			ProductPublishState: constants.PublishStatePublished,
		}
		err = productRepo.Create(ctx, &prod1)
		assert.NoError(t, err)

		code2, err := Utils.CreateNanoId(8)
		assert.NoError(t, err)
		prod2 := models.Product{
			ProductID:           uuid.NewString(),
			ProductName:         "Sepet Ürünü 2",
			Stock:               10,
			Price:               25000,
			ThumbnailImage:      "cart_thumb_2.jpg",
			ProductCode:         code2,
			ProductPublishState: constants.PublishStatePublished,
		}
		err = productRepo.Create(ctx, &prod2)
		assert.NoError(t, err)

		code3, err := Utils.CreateNanoId(8)
		assert.NoError(t, err)
		prodUnpublished := models.Product{
			ProductID:           uuid.NewString(),
			ProductName:         "Sepet Ürünü Yayınlanmamış",
			Stock:               5,
			Price:               5000,
			ThumbnailImage:      "cart_thumb_unpub.jpg",
			ProductCode:         code3,
			ProductPublishState: constants.PublishStateUnpublished,
		}
		err = productRepo.Create(ctx, &prodUnpublished)
		assert.NoError(t, err)

		// act
		queryIDs := []string{prod1.ProductID, prod2.ProductID, prodUnpublished.ProductID}
		results, err := productRepo.GetMany(ctx, queryIDs)

		// assert
		assert.NoError(t, err)
		assert.Len(t, results, 2)

		var found1, found2, foundUnpublished bool
		for _, item := range results {
			if item.ProductID == prod1.ProductID {
				found1 = true
				assert.Equal(t, prod1.ProductName, item.ProductName)
				assert.Equal(t, float64(prod1.Price), item.Price)
				assert.Equal(t, prod1.Stock, item.Stock)
				assert.Equal(t, prod1.ThumbnailImage, item.ThumbnailImage)
			}
			if item.ProductID == prod2.ProductID {
				found2 = true
				assert.Equal(t, prod2.ProductName, item.ProductName)
				assert.Equal(t, float64(prod2.Price), item.Price)
				assert.Equal(t, prod2.Stock, item.Stock)
				assert.Equal(t, prod2.ThumbnailImage, item.ThumbnailImage)
			}
			if item.ProductID == prodUnpublished.ProductID {
				foundUnpublished = true
			}
		}
		assert.True(t, found1, "Yayınlanmış prod1 sonuçlarda olmalıdır")
		assert.True(t, found2, "Yayınlanmış prod2 sonuçlarda olmalıdır")
		assert.False(t, foundUnpublished, "Yayınlanmamış ürün sonuçlarda olmamalıdır")
	})

	t.Run("boş ID listesi gönderildiğinde boş liste dönmeli", func(t *testing.T) {
		// arrange
		ctx := context.Background()

		// act
		results, err := productRepo.GetMany(ctx, []string{})

		// assert
		assert.NoError(t, err)
		assert.Empty(t, results)
	})

	t.Run("geçersiz UUID formatı içeren listede veritabanı hatası dönmeli", func(t *testing.T) {
		// arrange
		ctx := context.Background()

		// act
		results, err := productRepo.GetMany(ctx, []string{"geçersiz-uuid"})

		// assert
		assert.Error(t, err)
		assert.Nil(t, results)
		assert.Contains(t, err.Error(), "CollectRows işlemi sırasında hatayla karşılaşıldı")
	})

	t.Run("veritabanı hatası durumunda beklenen formatta hata dönmeli (Context Canceled)", func(t *testing.T) {
		// arrange
		ctx, cancel := context.WithCancel(context.Background())
		cancel()

		// act
		results, err := productRepo.GetMany(ctx, []string{uuid.NewString()})

		// assert
		assert.Error(t, err)
		assert.Nil(t, results)
		assert.Contains(t, err.Error(), "Veritabanı sorgusunda hatayla karşılaşıldı")
	})
}

func TestGetProductCount(t *testing.T) {
	productRepo := NewProductStore(testDbInstance)

	t.Run("yayınlanmış ürün sayısını doğru dönmeli ve yayınlanmamışları saymamalı", func(t *testing.T) {
		// arrange
		ctx := context.Background()

		initialCount, err := productRepo.GetProductCount(ctx)
		assert.NoError(t, err)
		assert.GreaterOrEqual(t, initialCount, 0)

		// 1 adet yayınlanmış ürün ekle
		code1, err := Utils.CreateNanoId(8)
		assert.NoError(t, err)
		pubProduct := models.Product{
			ProductID:           uuid.NewString(),
			ProductName:         "Sayım Test Yayınlanan",
			ProductCode:         code1,
			ProductPublishState: constants.PublishStatePublished,
		}
		err = productRepo.Create(ctx, &pubProduct)
		assert.NoError(t, err)

		// act 1
		countAfterPub, err := productRepo.GetProductCount(ctx)

		// assert 1
		assert.NoError(t, err)
		assert.Equal(t, initialCount+1, countAfterPub)

		// 1 adet yayınlanmamış ürün ekle
		code2, err := Utils.CreateNanoId(8)
		assert.NoError(t, err)
		unpubProduct := models.Product{
			ProductID:           uuid.NewString(),
			ProductName:         "Sayım Test Yayınlanmayan",
			ProductCode:         code2,
			ProductPublishState: constants.PublishStateUnpublished,
		}
		err = productRepo.Create(ctx, &unpubProduct)
		assert.NoError(t, err)

		// act 2
		countAfterUnpub, err := productRepo.GetProductCount(ctx)

		// assert 2: Yayınlanmamış ürün sayıyı artırmamalı
		assert.NoError(t, err)
		assert.Equal(t, countAfterPub, countAfterUnpub)
	})

	t.Run("veritabanı hatası durumunda -1 ve hata dönmeli (Context Canceled)", func(t *testing.T) {
		// arrange
		ctx, cancel := context.WithCancel(context.Background())
		cancel()

		// act
		count, err := productRepo.GetProductCount(ctx)

		// assert
		assert.Error(t, err)
		assert.Equal(t, -1, count)
	})
}

func TestGetProductsByCategoryId(t *testing.T) {
	productRepo := NewProductStore(testDbInstance)

	t.Run("kategoriye ait yayınlanmış ürünleri doğru listelemeli", func(t *testing.T) {
		// arrange
		ctx := context.Background()

		var categoryID string
		err := testDbInstance.QueryRow(ctx, "INSERT INTO productcategories (category_name) VALUES ($1) RETURNING category_id", "Kategori Test").Scan(&categoryID)
		assert.NoError(t, err)

		// Kategoriye ait yayınlanmış ürün
		code1, err := Utils.CreateNanoId(8)
		assert.NoError(t, err)
		prodPub := models.Product{
			ProductID:           uuid.NewString(),
			ProductName:         "Kategori Yayınlanan Ürün",
			Stock:               8,
			Price:               4500,
			ThumbnailImage:      "cat_pub.png",
			ProductCode:         code1,
			ProductPublishState: constants.PublishStatePublished,
		}
		err = productRepo.Create(ctx, &prodPub)
		assert.NoError(t, err)

		// Kategoriye ait yayınlanmamış ürün
		code2, err := Utils.CreateNanoId(8)
		assert.NoError(t, err)
		prodUnpub := models.Product{
			ProductID:           uuid.NewString(),
			ProductName:         "Kategori Yayınlanmayan Ürün",
			Stock:               3,
			Price:               2500,
			ThumbnailImage:      "cat_unpub.png",
			ProductCode:         code2,
			ProductPublishState: constants.PublishStateUnpublished,
		}
		err = productRepo.Create(ctx, &prodUnpub)
		assert.NoError(t, err)

		// Başka kategoriye ait ürün
		code3, err := Utils.CreateNanoId(8)
		assert.NoError(t, err)
		prodOther := models.Product{
			ProductID:           uuid.NewString(),
			ProductName:         "Başka Kategori Ürün",
			ProductCode:         code3,
			ProductPublishState: constants.PublishStatePublished,
		}
		err = productRepo.Create(ctx, &prodOther)
		assert.NoError(t, err)

		// İlk iki ürünü kategoriye bağla
		_, err = testDbInstance.Exec(ctx, "INSERT INTO productproductcategories (product_id, category_id) VALUES ($1, $2), ($3, $4)",
			prodPub.ProductID, categoryID,
			prodUnpub.ProductID, categoryID,
		)
		assert.NoError(t, err)

		// act
		results, err := productRepo.GetProductsByCategoryId(ctx, []string{categoryID}, "10", "0")

		// assert
		assert.NoError(t, err)
		assert.NotEmpty(t, results)

		var foundPub, foundUnpub, foundOther bool
		for _, p := range results {
			if p.ProductID == prodPub.ProductID {
				foundPub = true
				assert.Equal(t, prodPub.ProductName, p.ProductName)
				assert.Equal(t, prodPub.Stock, p.Stock)
				assert.Equal(t, prodPub.Price, p.Price)
				assert.Equal(t, prodPub.ThumbnailImage, p.ThumbnailImage)
			}
			if p.ProductID == prodUnpub.ProductID {
				foundUnpub = true
			}
			if p.ProductID == prodOther.ProductID {
				foundOther = true
			}
		}
		assert.True(t, foundPub, "Kategoriye bağlı yayınlanmış ürün sonuçlarda olmalı")
		assert.False(t, foundUnpub, "Kategoriye bağlı yayınlanmamış ürün sonuçlarda olmamalı")
		assert.False(t, foundOther, "Farklı kategori ürünü sonuçlarda olmamalı")
	})

	t.Run("sayfalama (limit ve offset) parametreleri doğru çalışmalı", func(t *testing.T) {
		// arrange
		ctx := context.Background()

		var categoryID string
		err := testDbInstance.QueryRow(ctx, "INSERT INTO productcategories (category_name) VALUES ($1) RETURNING category_id", "Sayfalama Kategori").Scan(&categoryID)
		assert.NoError(t, err)

		// 2 yayınlanmış ürün ekle ve kategoriye bağla
		code1, _ := Utils.CreateNanoId(8)
		p1 := models.Product{ProductID: uuid.NewString(), ProductName: "P1", ProductCode: code1, ProductPublishState: constants.PublishStatePublished}
		assert.NoError(t, productRepo.Create(ctx, &p1))

		code2, _ := Utils.CreateNanoId(8)
		p2 := models.Product{ProductID: uuid.NewString(), ProductName: "P2", ProductCode: code2, ProductPublishState: constants.PublishStatePublished}
		assert.NoError(t, productRepo.Create(ctx, &p2))

		_, err = testDbInstance.Exec(ctx, "INSERT INTO productproductcategories (product_id, category_id) VALUES ($1, $2), ($3, $4)",
			p1.ProductID, categoryID,
			p2.ProductID, categoryID,
		)
		assert.NoError(t, err)

		// act: limit 1
		resLimit1, err := productRepo.GetProductsByCategoryId(ctx, []string{categoryID}, "1", "0")
		assert.NoError(t, err)
		assert.Len(t, resLimit1, 1)

		// act: offset 1
		resOffset1, err := productRepo.GetProductsByCategoryId(ctx, []string{categoryID}, "1", "1")
		assert.NoError(t, err)
		assert.Len(t, resOffset1, 1)
		assert.NotEqual(t, resLimit1[0].ProductID, resOffset1[0].ProductID)
	})

	t.Run("kategoriye ait ürün bulunamadığında boş liste dönmeli", func(t *testing.T) {
		// arrange
		ctx := context.Background()
		emptyCatID := uuid.NewString()

		// act
		results, err := productRepo.GetProductsByCategoryId(ctx, []string{emptyCatID}, "10", "0")

		// assert
		assert.NoError(t, err)
		assert.Empty(t, results)
	})

	t.Run("veritabanı hatası durumunda hata dönmeli (Context Canceled)", func(t *testing.T) {
		// arrange
		ctx, cancel := context.WithCancel(context.Background())
		cancel()

		// act
		results, err := productRepo.GetProductsByCategoryId(ctx, []string{uuid.NewString()}, "10", "0")

		// assert
		assert.Error(t, err)
		assert.Nil(t, results)
	})
}

func TestGetCountProductsByCategory(t *testing.T) {
	productRepo := NewProductStore(testDbInstance)

	t.Run("kategoriye ait yayınlanmış ürün sayısını doğru dönmeli", func(t *testing.T) {
		// arrange
		ctx := context.Background()

		var categoryID string
		err := testDbInstance.QueryRow(ctx, "INSERT INTO productcategories (category_name) VALUES ($1) RETURNING category_id", "Sayım Kategori").Scan(&categoryID)
		assert.NoError(t, err)

		// Başlangıçta 0 olmalı
		initCount, err := productRepo.GetCountProductsByCategory(ctx, []string{categoryID})
		assert.NoError(t, err)
		assert.Equal(t, 0, initCount)

		// 1 yayınlanmış ürün ekle ve kategoriye bağla
		code1, _ := Utils.CreateNanoId(8)
		pubProd := models.Product{ProductID: uuid.NewString(), ProductName: "Kategori Sayım 1", ProductCode: code1, ProductPublishState: constants.PublishStatePublished}
		assert.NoError(t, productRepo.Create(ctx, &pubProd))
		_, err = testDbInstance.Exec(ctx, "INSERT INTO productproductcategories (product_id, category_id) VALUES ($1, $2)", pubProd.ProductID, categoryID)
		assert.NoError(t, err)

		// act 1
		count1, err := productRepo.GetCountProductsByCategory(ctx, []string{categoryID})

		// assert 1
		assert.NoError(t, err)
		assert.Equal(t, 1, count1)

		// 1 yayınlanmamış ürün ekle ve kategoriye bağla
		code2, _ := Utils.CreateNanoId(8)
		unpubProd := models.Product{ProductID: uuid.NewString(), ProductName: "Kategori Sayım 2 Unpub", ProductCode: code2, ProductPublishState: constants.PublishStateUnpublished}
		assert.NoError(t, productRepo.Create(ctx, &unpubProd))
		_, err = testDbInstance.Exec(ctx, "INSERT INTO productproductcategories (product_id, category_id) VALUES ($1, $2)", unpubProd.ProductID, categoryID)
		assert.NoError(t, err)

		// act 2
		count2, err := productRepo.GetCountProductsByCategory(ctx, []string{categoryID})

		// assert 2: Yayınlanmamış ürün sayıyı değiştirmemeli
		assert.NoError(t, err)
		assert.Equal(t, 1, count2)
	})

	t.Run("eşleşen ürün olmadığında 0 dönmeli", func(t *testing.T) {
		// arrange
		ctx := context.Background()

		// act
		count, err := productRepo.GetCountProductsByCategory(ctx, []string{uuid.NewString()})

		// assert
		assert.NoError(t, err)
		assert.Equal(t, 0, count)
	})

	t.Run("veritabanı hatası durumunda -1 ve hata dönmeli (Context Canceled)", func(t *testing.T) {
		// arrange
		ctx, cancel := context.WithCancel(context.Background())
		cancel()

		// act
		count, err := productRepo.GetCountProductsByCategory(ctx, []string{uuid.NewString()})

		// assert
		assert.Error(t, err)
		assert.Equal(t, -1, count)
	})
}

func TestSearchProductCount(t *testing.T) {
	productRepo := NewProductStore(testDbInstance)

	t.Run("yayınlanmış ürünler arasında arama sonucunu doğru saymalı ve yayınlanmamışları dahil etmemeli", func(t *testing.T) {
		// arrange
		ctx := context.Background()
		keyword := fmt.Sprintf("klavyetest%s", strings.ToLower(uuid.NewString()[:6]))

		code1, _ := Utils.CreateNanoId(8)
		pubProd := models.Product{
			ProductID:           uuid.NewString(),
			ProductName:         fmt.Sprintf("Mekanik %s RGB", keyword),
			ProductCode:         code1,
			ProductPublishState: constants.PublishStatePublished,
		}
		assert.NoError(t, productRepo.Create(ctx, &pubProd))

		code2, _ := Utils.CreateNanoId(8)
		unpubProd := models.Product{
			ProductID:           uuid.NewString(),
			ProductName:         fmt.Sprintf("Kablosuz %s Standart", keyword),
			ProductCode:         code2,
			ProductPublishState: constants.PublishStateUnpublished,
		}
		assert.NoError(t, productRepo.Create(ctx, &unpubProd))

		// act
		count := productRepo.SearchProductCount(ctx, keyword)

		// assert
		assert.Equal(t, 1, count, "Yalnızca yayınlanmış ürün sayılmalı")
	})

	t.Run("eşleşen ürün bulunamadığında 0 dönmeli", func(t *testing.T) {
		// arrange
		ctx := context.Background()

		// act
		count := productRepo.SearchProductCount(ctx, "eslesmeyencokozelkelime12345")

		// assert
		assert.Equal(t, 0, count)
	})

	t.Run("veritabanı hatası durumunda 0 dönmeli (Context Canceled)", func(t *testing.T) {
		// arrange
		ctx, cancel := context.WithCancel(context.Background())
		cancel()

		// act
		count := productRepo.SearchProductCount(ctx, "arama")

		// assert
		assert.Equal(t, 0, count)
	})
}

func TestSearchProduct(t *testing.T) {
	productRepo := NewProductStore(testDbInstance)

	t.Run("yayınlanmış ürünlerde arama yapmalı ve sonuçları döndürmeli", func(t *testing.T) {
		// arrange
		ctx := context.Background()
		keyword := fmt.Sprintf("kulakliktest%s", strings.ToLower(uuid.NewString()[:6]))

		code1, _ := Utils.CreateNanoId(8)
		pubProd := models.Product{
			ProductID:           uuid.NewString(),
			ProductName:         fmt.Sprintf("Bluetooth %s Pro", keyword),
			Stock:               20,
			Price:               75000,
			ThumbnailImage:      "kulaklik.png",
			ProductCode:         code1,
			ProductPublishState: constants.PublishStatePublished,
		}
		assert.NoError(t, productRepo.Create(ctx, &pubProd))

		code2, _ := Utils.CreateNanoId(8)
		unpubProd := models.Product{
			ProductID:           uuid.NewString(),
			ProductName:         fmt.Sprintf("Kablolu %s Basic", keyword),
			ProductCode:         code2,
			ProductPublishState: constants.PublishStateUnpublished,
		}
		assert.NoError(t, productRepo.Create(ctx, &unpubProd))

		// act
		results, err := productRepo.SearchProduct(ctx, keyword, "10")

		// assert
		assert.NoError(t, err)
		assert.Len(t, results, 1)
		assert.Equal(t, pubProd.ProductID, results[0].ProductID)
		assert.Equal(t, pubProd.ProductName, results[0].ProductName)
		assert.Equal(t, pubProd.Price, results[0].Price)
		assert.Equal(t, pubProd.Stock, results[0].Stock)
		assert.Equal(t, pubProd.ThumbnailImage, results[0].ThumbnailImage)
	})

	t.Run("eşleşen ürün bulunamadığında boş liste dönmeli", func(t *testing.T) {
		// arrange
		ctx := context.Background()

		// act
		results, err := productRepo.SearchProduct(ctx, "bulunamayacakkelime9999", "10")

		// assert
		assert.NoError(t, err)
		assert.Empty(t, results)
	})

	t.Run("veritabanı hatası durumunda hata dönmeli (Context Canceled)", func(t *testing.T) {
		// arrange
		ctx, cancel := context.WithCancel(context.Background())
		cancel()

		// act
		results, err := productRepo.SearchProduct(ctx, "test", "10")

		// assert
		assert.Error(t, err)
		assert.Nil(t, results)
	})
}

func TestAlternateSearch(t *testing.T) {
	productRepo := NewProductStore(testDbInstance)

	t.Run("trigram benzerliği ile yayınlanmış ürünleri aramalı", func(t *testing.T) {
		// arrange
		ctx := context.Background()
		uniquePart := strings.ToLower(uuid.NewString()[:6])
		uniqueName := fmt.Sprintf("MonitorOzel%s", uniquePart)

		code1, _ := Utils.CreateNanoId(8)
		pubProd := models.Product{
			ProductID:           uuid.NewString(),
			ProductName:         uniqueName,
			Stock:               5,
			Price:               450000,
			ThumbnailImage:      "monitor.png",
			ProductCode:         code1,
			ProductPublishState: constants.PublishStatePublished,
		}
		assert.NoError(t, productRepo.Create(ctx, &pubProd))

		code2, _ := Utils.CreateNanoId(8)
		unpubProd := models.Product{
			ProductID:           uuid.NewString(),
			ProductName:         fmt.Sprintf("MonitorGizli%s", uniquePart),
			ProductCode:         code2,
			ProductPublishState: constants.PublishStateUnpublished,
		}
		assert.NoError(t, productRepo.Create(ctx, &unpubProd))

		// act: benzeyen kelime ile ara
		results, err := productRepo.AlternateSearch(ctx, uniqueName, "5")

		// assert
		assert.NoError(t, err)
		assert.NotEmpty(t, results)
		assert.Equal(t, pubProd.ProductID, results[0].ProductID, "En benzer ve yayınlanmış ürün ilk sırada gelmeli")

		// Yayınlanmamış ürün sonuçlarda olmamalı
		for _, item := range results {
			assert.NotEqual(t, unpubProd.ProductID, item.ProductID)
		}
	})

	t.Run("veritabanı hatası durumunda hata dönmeli (Context Canceled)", func(t *testing.T) {
		// arrange
		ctx, cancel := context.WithCancel(context.Background())
		cancel()

		// act
		results, err := productRepo.AlternateSearch(ctx, "test", "5")

		// assert
		assert.Error(t, err)
		assert.Nil(t, results)
	})
}

func TestFilteredAndSortedProducts(t *testing.T) {
	productRepo := NewProductStore(testDbInstance)

	t.Run("kategori filtresi ve fiyata göre artan (ASC) sıralama", func(t *testing.T) {
		// arrange
		ctx := context.Background()

		var categoryID string
		err := testDbInstance.QueryRow(ctx, "INSERT INTO productcategories (category_name) VALUES ($1) RETURNING category_id", "Filtre Kategori").Scan(&categoryID)
		assert.NoError(t, err)

		code1, _ := Utils.CreateNanoId(8)
		pCheap := models.Product{
			ProductID:           uuid.NewString(),
			ProductName:         "Ucuz Ürün",
			Stock:               10,
			Price:               1000,
			ProductCode:         code1,
			ProductPublishState: constants.PublishStatePublished,
		}
		assert.NoError(t, productRepo.Create(ctx, &pCheap))

		code2, _ := Utils.CreateNanoId(8)
		pExpensive := models.Product{
			ProductID:           uuid.NewString(),
			ProductName:         "Pahalı Ürün",
			Stock:               5,
			Price:               5000,
			ProductCode:         code2,
			ProductPublishState: constants.PublishStatePublished,
		}
		assert.NoError(t, productRepo.Create(ctx, &pExpensive))

		// Kategoriye bağla
		_, err = testDbInstance.Exec(ctx, "INSERT INTO productproductcategories (product_id, category_id) VALUES ($1, $2), ($3, $4)",
			pCheap.ProductID, categoryID,
			pExpensive.ProductID, categoryID,
		)
		assert.NoError(t, err)

		// act: price asc
		results, err := productRepo.FilteredAndSortedProducts(ctx, []string{categoryID}, "10", "0", constants.Price, constants.Asc)

		// assert
		assert.NoError(t, err)
		assert.GreaterOrEqual(t, len(results), 2)
		assert.Equal(t, pCheap.ProductID, results[0].ProductID)
		assert.Equal(t, pExpensive.ProductID, results[1].ProductID)
		assert.LessOrEqual(t, results[0].Price, results[1].Price)
	})

	t.Run("kategori filtresi ve fiyata göre azalan (DESC) sıralama", func(t *testing.T) {
		// arrange
		ctx := context.Background()

		var categoryID string
		err := testDbInstance.QueryRow(ctx, "INSERT INTO productcategories (category_name) VALUES ($1) RETURNING category_id", "Filtre Kategori DESC").Scan(&categoryID)
		assert.NoError(t, err)

		code1, _ := Utils.CreateNanoId(8)
		p1 := models.Product{ProductID: uuid.NewString(), ProductName: "A", Price: 2000, ProductCode: code1, ProductPublishState: constants.PublishStatePublished}
		assert.NoError(t, productRepo.Create(ctx, &p1))

		code2, _ := Utils.CreateNanoId(8)
		p2 := models.Product{ProductID: uuid.NewString(), ProductName: "B", Price: 8000, ProductCode: code2, ProductPublishState: constants.PublishStatePublished}
		assert.NoError(t, productRepo.Create(ctx, &p2))

		_, err = testDbInstance.Exec(ctx, "INSERT INTO productproductcategories (product_id, category_id) VALUES ($1, $2), ($3, $4)",
			p1.ProductID, categoryID,
			p2.ProductID, categoryID,
		)
		assert.NoError(t, err)

		// act: price desc
		results, err := productRepo.FilteredAndSortedProducts(ctx, []string{categoryID}, "10", "0", constants.Price, constants.Desc)

		// assert
		assert.NoError(t, err)
		assert.GreaterOrEqual(t, len(results), 2)
		assert.Equal(t, p2.ProductID, results[0].ProductID)
		assert.Equal(t, p1.ProductID, results[1].ProductID)
		assert.GreaterOrEqual(t, results[0].Price, results[1].Price)
	})

	t.Run("kategori filtresi olmadan ürünleri listeleme", func(t *testing.T) {
		// arrange
		ctx := context.Background()

		// act
		results, err := productRepo.FilteredAndSortedProducts(ctx, nil, "5", "0", constants.Price, constants.Asc)

		// assert
		assert.NoError(t, err)
		assert.NotEmpty(t, results)
		assert.LessOrEqual(t, len(results), 5)
	})

	t.Run("veritabanı hatası durumunda hata dönmeli (Context Canceled)", func(t *testing.T) {
		// arrange
		ctx, cancel := context.WithCancel(context.Background())
		cancel()

		// act
		results, err := productRepo.FilteredAndSortedProducts(ctx, nil, "10", "0", constants.Price, constants.Asc)

		// assert
		assert.Error(t, err)
		assert.Nil(t, results)
	})
}

func TestSearchProductCountAdmin(t *testing.T) {
	productRepo := NewProductStore(testDbInstance)

	t.Run("yayınlanmış ve yayınlanmamış tüm ürünleri saymalı", func(t *testing.T) {
		// arrange
		ctx := context.Background()
		keyword := fmt.Sprintf("admincount%s", strings.ToLower(uuid.NewString()[:6]))

		code1, _ := Utils.CreateNanoId(8)
		pubProd := models.Product{
			ProductID:           uuid.NewString(),
			ProductName:         fmt.Sprintf("Ürün 1 %s", keyword),
			ProductCode:         code1,
			ProductPublishState: constants.PublishStatePublished,
		}
		assert.NoError(t, productRepo.Create(ctx, &pubProd))

		code2, _ := Utils.CreateNanoId(8)
		unpubProd := models.Product{
			ProductID:           uuid.NewString(),
			ProductName:         fmt.Sprintf("Ürün 2 %s", keyword),
			ProductCode:         code2,
			ProductPublishState: constants.PublishStateUnpublished,
		}
		assert.NoError(t, productRepo.Create(ctx, &unpubProd))

		// act
		count := productRepo.SearchProductCountAdmin(ctx, keyword)

		// assert: Her ikisi de (published + unpublished) sayılmalı
		assert.Equal(t, 2, count)
	})

	t.Run("eşleşen ürün bulunamadığında 0 dönmeli", func(t *testing.T) {
		// arrange
		ctx := context.Background()

		// act
		count := productRepo.SearchProductCountAdmin(ctx, "eslesmeyenadminkelime999")

		// assert
		assert.Equal(t, 0, count)
	})
}

func TestSearchProductAdmin(t *testing.T) {
	productRepo := NewProductStore(testDbInstance)

	t.Run("admin için hem yayınlanan hem yayınlanmayan ürünleri getirmeli", func(t *testing.T) {
		// arrange
		ctx := context.Background()
		keyword := fmt.Sprintf("adminsearch%s", strings.ToLower(uuid.NewString()[:6]))

		code1, _ := Utils.CreateNanoId(8)
		pubProd := models.Product{
			ProductID:           uuid.NewString(),
			ProductName:         fmt.Sprintf("Admin Ürün %s Pub", keyword),
			Stock:               10,
			Price:               5000,
			ThumbnailImage:      "pub.png",
			ProductCode:         code1,
			ProductPublishState: constants.PublishStatePublished,
		}
		assert.NoError(t, productRepo.Create(ctx, &pubProd))

		code2, _ := Utils.CreateNanoId(8)
		unpubProd := models.Product{
			ProductID:           uuid.NewString(),
			ProductName:         fmt.Sprintf("Admin Ürün %s Unpub", keyword),
			Stock:               5,
			Price:               2500,
			ThumbnailImage:      "unpub.png",
			ProductCode:         code2,
			ProductPublishState: constants.PublishStateUnpublished,
		}
		assert.NoError(t, productRepo.Create(ctx, &unpubProd))

		// act
		results, err := productRepo.SearchProductAdmin(ctx, keyword, "10")

		// assert
		assert.NoError(t, err)
		assert.Len(t, results, 2)

		var foundPub, foundUnpub bool
		for _, item := range results {
			if item.ProductID == pubProd.ProductID {
				foundPub = true
				assert.Equal(t, pubProd.ProductCode, item.ProductCode)
				assert.Equal(t, constants.PublishStatePublished, item.ProductPublishState)
			}
			if item.ProductID == unpubProd.ProductID {
				foundUnpub = true
				assert.Equal(t, unpubProd.ProductCode, item.ProductCode)
				assert.Equal(t, constants.PublishStateUnpublished, item.ProductPublishState)
			}
		}
		assert.True(t, foundPub)
		assert.True(t, foundUnpub)
	})

	t.Run("veritabanı hatası durumunda hata dönmeli (Context Canceled)", func(t *testing.T) {
		// arrange
		ctx, cancel := context.WithCancel(context.Background())
		cancel()

		// act
		results, err := productRepo.SearchProductAdmin(ctx, "test", "10")

		// assert
		assert.Error(t, err)
		assert.Nil(t, results)
	})
}

func TestAlternateSearchAdmin(t *testing.T) {
	productRepo := NewProductStore(testDbInstance)

	t.Run("trigram araması ile admin için tüm durumdaki ürünleri getirmeli", func(t *testing.T) {
		// arrange
		ctx := context.Background()
		uniquePart := strings.ToLower(uuid.NewString()[:6])
		uniqueName := fmt.Sprintf("AdminTrigram%s", uniquePart)

		code1, _ := Utils.CreateNanoId(8)
		pubProd := models.Product{
			ProductID:           uuid.NewString(),
			ProductName:         uniqueName,
			ProductCode:         code1,
			ProductPublishState: constants.PublishStatePublished,
		}
		assert.NoError(t, productRepo.Create(ctx, &pubProd))

		code2, _ := Utils.CreateNanoId(8)
		unpubProd := models.Product{
			ProductID:           uuid.NewString(),
			ProductName:         fmt.Sprintf("AdminTrigramUnpub%s", uniquePart),
			ProductCode:         code2,
			ProductPublishState: constants.PublishStateUnpublished,
		}
		assert.NoError(t, productRepo.Create(ctx, &unpubProd))

		// act
		results, err := productRepo.AlternateSearchAdmin(ctx, uniqueName, "10")

		// assert
		assert.NoError(t, err)
		assert.NotEmpty(t, results)

		var foundPub, foundUnpub bool
		for _, item := range results {
			if item.ProductID == pubProd.ProductID {
				foundPub = true
			}
			if item.ProductID == unpubProd.ProductID {
				foundUnpub = true
			}
		}
		assert.True(t, foundPub)
		assert.True(t, foundUnpub)
	})

	t.Run("veritabanı hatası durumunda hata dönmeli (Context Canceled)", func(t *testing.T) {
		// arrange
		ctx, cancel := context.WithCancel(context.Background())
		cancel()

		// act
		results, err := productRepo.AlternateSearchAdmin(ctx, "test", "5")

		// assert
		assert.Error(t, err)
		assert.Nil(t, results)
	})
}

func TestGetProductCountAdmin(t *testing.T) {
	productRepo := NewProductStore(testDbInstance)

	t.Run("tüm durumdaki ürünlerin toplam sayısını dönmeli", func(t *testing.T) {
		// arrange
		ctx := context.Background()

		initialCount, err := productRepo.GetProductCountAdmin(ctx)
		assert.NoError(t, err)

		// 1 yayınlanmış ürün ekle
		code1, _ := Utils.CreateNanoId(8)
		pubProd := models.Product{ProductID: uuid.NewString(), ProductName: "Count Admin 1", ProductCode: code1, ProductPublishState: constants.PublishStatePublished}
		assert.NoError(t, productRepo.Create(ctx, &pubProd))

		c1, err := productRepo.GetProductCountAdmin(ctx)
		assert.NoError(t, err)
		assert.Equal(t, initialCount+1, c1)

		// 1 yayınlanmamış ürün ekle
		code2, _ := Utils.CreateNanoId(8)
		unpubProd := models.Product{ProductID: uuid.NewString(), ProductName: "Count Admin 2", ProductCode: code2, ProductPublishState: constants.PublishStateUnpublished}
		assert.NoError(t, productRepo.Create(ctx, &unpubProd))

		c2, err := productRepo.GetProductCountAdmin(ctx)
		assert.NoError(t, err)
		assert.Equal(t, initialCount+2, c2)
	})

	t.Run("veritabanı hatası durumunda -1 ve hata dönmeli (Context Canceled)", func(t *testing.T) {
		// arrange
		ctx, cancel := context.WithCancel(context.Background())
		cancel()

		// act
		count, err := productRepo.GetProductCountAdmin(ctx)

		// assert
		assert.Error(t, err)
		assert.Equal(t, -1, count)
	})
}

func TestGetAllAdmin(t *testing.T) {
	productRepo := NewProductStore(testDbInstance)

	t.Run("yayınlanmış ve yayınlanmamış tüm ürünleri admin için listelemeli", func(t *testing.T) {
		// arrange
		ctx := context.Background()

		code1, _ := Utils.CreateNanoId(8)
		pubProd := models.Product{
			ProductID:           uuid.NewString(),
			ProductName:         "GetAllAdmin Yayında",
			Stock:               10,
			Price:               1000,
			ThumbnailImage:      "pub.jpg",
			ProductCode:         code1,
			ProductPublishState: constants.PublishStatePublished,
		}
		assert.NoError(t, productRepo.Create(ctx, &pubProd))

		code2, _ := Utils.CreateNanoId(8)
		unpubProd := models.Product{
			ProductID:           uuid.NewString(),
			ProductName:         "GetAllAdmin Taslak",
			Stock:               5,
			Price:               500,
			ThumbnailImage:      "draft.jpg",
			ProductCode:         code2,
			ProductPublishState: constants.PublishStateUnpublished,
		}
		assert.NoError(t, productRepo.Create(ctx, &unpubProd))

		// act
		results, err := productRepo.GetAllAdmin(100, 0, ctx)

		// assert
		assert.NoError(t, err)
		assert.NotEmpty(t, results)

		var foundPub, foundUnpub bool
		for _, item := range results {
			if item.ProductID == pubProd.ProductID {
				foundPub = true
				assert.Equal(t, pubProd.ProductName, item.ProductName)
				assert.Equal(t, pubProd.ProductCode, item.ProductCode)
				assert.Equal(t, constants.PublishStatePublished, item.ProductPublishState)
			}
			if item.ProductID == unpubProd.ProductID {
				foundUnpub = true
				assert.Equal(t, unpubProd.ProductName, item.ProductName)
				assert.Equal(t, unpubProd.ProductCode, item.ProductCode)
				assert.Equal(t, constants.PublishStateUnpublished, item.ProductPublishState)
			}
		}
		assert.True(t, foundPub)
		assert.True(t, foundUnpub)
	})

	t.Run("sayfalama (limit ve offset) parametreleri doğru çalışmalı", func(t *testing.T) {
		// arrange
		ctx := context.Background()

		// act: limit 1
		res1, err := productRepo.GetAllAdmin(1, 0, ctx)
		assert.NoError(t, err)
		assert.LessOrEqual(t, len(res1), 1)

		// act: offset 1
		res2, err := productRepo.GetAllAdmin(1, 1, ctx)
		assert.NoError(t, err)
		if len(res1) > 0 && len(res2) > 0 {
			assert.NotEqual(t, res1[0].ProductID, res2[0].ProductID)
		}
	})

	t.Run("veritabanı hatası durumunda beklenen formatta hata dönmeli (Context Canceled)", func(t *testing.T) {
		// arrange
		ctx, cancel := context.WithCancel(context.Background())
		cancel()

		// act
		results, err := productRepo.GetAllAdmin(10, 0, ctx)

		// assert
		assert.Error(t, err)
		assert.Nil(t, results)
		assert.Contains(t, err.Error(), "ürün sorgu hatası")
	})
}

func TestChangePublishState(t *testing.T) {
	productRepo := NewProductStore(testDbInstance)

	t.Run("başarılı yayın durumu güncelleme", func(t *testing.T) {
		// arrange
		ctx := context.Background()
		code, _ := Utils.CreateNanoId(8)
		prod := models.Product{
			ProductID:           uuid.NewString(),
			ProductName:         "Durum Test Ürünü",
			ProductCode:         code,
			ProductPublishState: constants.PublishStateUnpublished,
		}
		assert.NoError(t, productRepo.Create(ctx, &prod))

		// act: unpublished -> published
		err := productRepo.ChangePublishState(ctx, prod.ProductID, constants.PublishStatePublished)

		// assert
		assert.NoError(t, err)

		var currentState string
		var updatedAt *time.Time
		err = testDbInstance.QueryRow(ctx, "SELECT product_publish_state, updated_at FROM products WHERE product_id = $1", prod.ProductID).Scan(&currentState, &updatedAt)
		assert.NoError(t, err)
		assert.Equal(t, constants.PublishStatePublished, currentState)
		assert.NotNil(t, updatedAt)
	})

	t.Run("olmayan bir ürünün durumunu güncelleme hata fırlatmaz", func(t *testing.T) {
		// arrange
		ctx := context.Background()

		// act
		err := productRepo.ChangePublishState(ctx, uuid.NewString(), constants.PublishStatePublished)

		// assert
		assert.NoError(t, err)
	})

	t.Run("geçersiz UUID formatında veritabanı hatası dönmeli", func(t *testing.T) {
		// arrange
		ctx := context.Background()

		// act
		err := productRepo.ChangePublishState(ctx, "gecersiz-uuid", constants.PublishStatePublished)

		// assert
		assert.Error(t, err)
		assert.Contains(t, err.Error(), "ürün yayın durumu güncellenemedi")
	})

	t.Run("veritabanı hatası durumunda beklenen formatta hata dönmeli (Context Canceled)", func(t *testing.T) {
		// arrange
		ctx, cancel := context.WithCancel(context.Background())
		cancel()

		// act
		err := productRepo.ChangePublishState(ctx, uuid.NewString(), constants.PublishStatePublished)

		// assert
		assert.Error(t, err)
		assert.Contains(t, err.Error(), "ürün yayın durumu güncellenemedi")
	})
}



