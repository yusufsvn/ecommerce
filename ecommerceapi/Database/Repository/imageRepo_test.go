package repository

import (
	"context"
	"fmt"
	"testing"

	models "ecommerceapi/Database/Models"
	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestAddOneImage(t *testing.T) {
	imageRepo := NewImageStore(testDbInstance)

	t.Run("ürün için görsel ekler", func(t *testing.T) {
		ctx := context.Background()
		productID := createTestProduct(t, ctx, testDbInstance)

		image := models.Image{
			ImageID:   uuid.NewString(),
			ImagePath: fmt.Sprintf("/images/%s.jpg", uuid.NewString()),
			ProductID: productID,
		}

		err := imageRepo.AddOneImage(ctx, image)
		require.NoError(t, err)

		var storedID, storedPath, storedProductID string
		err = testDbInstance.QueryRow(
			ctx,
			"SELECT image_id, image_path, product_id FROM images WHERE image_id=$1",
			image.ImageID,
		).Scan(&storedID, &storedPath, &storedProductID)
		require.NoError(t, err)

		assert.Equal(t, image.ImageID, storedID)
		assert.Equal(t, image.ImagePath, storedPath)
		assert.Equal(t, image.ProductID, storedProductID)
	})
}

func TestDeleteImage(t *testing.T) {
	imageRepo := NewImageStore(testDbInstance)

	t.Run("seçilen görselleri siler", func(t *testing.T) {
		ctx := context.Background()
		productID := createTestProduct(t, ctx, testDbInstance)

		firstImage := models.Image{
			ImageID:   uuid.NewString(),
			ImagePath: fmt.Sprintf("/images/%s.jpg", uuid.NewString()),
			ProductID: productID,
		}
		secondImage := models.Image{
			ImageID:   uuid.NewString(),
			ImagePath: fmt.Sprintf("/images/%s.jpg", uuid.NewString()),
			ProductID: productID,
		}

		require.NoError(t, imageRepo.AddOneImage(ctx, firstImage))
		require.NoError(t, imageRepo.AddOneImage(ctx, secondImage))

		err := imageRepo.DeleteImage(ctx, []string{firstImage.ImageID, secondImage.ImageID})
		require.NoError(t, err)

		var count int
		err = testDbInstance.QueryRow(
			ctx,
			"SELECT COUNT(*) FROM images WHERE image_id = ANY($1)",
			[]string{firstImage.ImageID, secondImage.ImageID},
		).Scan(&count)
		require.NoError(t, err)
		assert.Equal(t, 0, count)
	})
}

func TestGetProductImages(t *testing.T) {
	imageRepo := NewImageStore(testDbInstance)

	t.Run("ürünün görsellerini listeler", func(t *testing.T) {
		ctx := context.Background()
		productID := createTestProduct(t, ctx, testDbInstance)

		firstImage := models.Image{
			ImageID:   uuid.NewString(),
			ImagePath: fmt.Sprintf("/images/%s.jpg", uuid.NewString()),
			ProductID: productID,
		}
		secondImage := models.Image{
			ImageID:   uuid.NewString(),
			ImagePath: fmt.Sprintf("/images/%s.jpg", uuid.NewString()),
			ProductID: productID,
		}

		require.NoError(t, imageRepo.AddOneImage(ctx, firstImage))
		require.NoError(t, imageRepo.AddOneImage(ctx, secondImage))

		imageIDs, err := imageRepo.GetProductImages(ctx, productID)
		require.NoError(t, err)
		assert.ElementsMatch(t, []string{firstImage.ImageID, secondImage.ImageID}, imageIDs)
	})
}
