package services

import (
	"bytes"
	"context"
	models "ecommerceapi/Database/Models"
	st "ecommerceapi/Structs"
	"errors"
	"io"
	"mime/multipart"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

type mockImageRepoForImageService struct {
	addOneImageFn      func(ctx context.Context, image models.Image) error
	deleteImageFn      func(ctx context.Context, images []string) error
	getProductImagesFn func(ctx context.Context, product_id string) ([]string, error)
}

func (m *mockImageRepoForImageService) AddOneImage(ctx context.Context, image models.Image) error {
	if m.addOneImageFn != nil {
		return m.addOneImageFn(ctx, image)
	}
	return nil
}

func (m *mockImageRepoForImageService) DeleteImage(ctx context.Context, images []string) error {
	if m.deleteImageFn != nil {
		return m.deleteImageFn(ctx, images)
	}
	return nil
}

func (m *mockImageRepoForImageService) GetProductImages(ctx context.Context, product_id string) ([]string, error) {
	if m.getProductImagesFn != nil {
		return m.getProductImagesFn(ctx, product_id)
	}
	return nil, nil
}

type mockFileManagerForImageService struct {
	uploadFileFn             func(fileName string, fileSize int64, prefix string, ctx context.Context, reader io.Reader) (string, error)
	deleteFileFn             func(ctx context.Context, image_id string) error
	deleteMultipleFileFn     func(ctx context.Context, images []string) error
	getPresignedFileObjectFn func() (st.UrlResult, error)
}

func (m *mockFileManagerForImageService) UploadFile(fileName string, fileSize int64, prefix string, ctx context.Context, reader io.Reader) (string, error) {
	if m.uploadFileFn != nil {
		return m.uploadFileFn(fileName, fileSize, prefix, ctx, reader)
	}
	return fileName, nil
}

func (m *mockFileManagerForImageService) DeleteFile(ctx context.Context, image_id string) error {
	if m.deleteFileFn != nil {
		return m.deleteFileFn(ctx, image_id)
	}
	return nil
}

func (m *mockFileManagerForImageService) DeleteMultipleFile(ctx context.Context, images []string) error {
	if m.deleteMultipleFileFn != nil {
		return m.deleteMultipleFileFn(ctx, images)
	}
	return nil
}

func (m *mockFileManagerForImageService) GetPresignedFileObject() (st.UrlResult, error) {
	if m.getPresignedFileObjectFn != nil {
		return m.getPresignedFileObjectFn()
	}
	return st.UrlResult{}, nil
}

func createTestFileHeader(t *testing.T, filename string, content []byte) *multipart.FileHeader {
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

func TestImageService_DeleteMultipleImage(t *testing.T) {
	mockRepo := &mockImageRepoForImageService{}
	mockFM := &mockFileManagerForImageService{}
	svc := NewImageService(mockRepo, mockFM)
	ctx := context.Background()

	t.Run("veritabanı silme işlemi hata verirse hata dönmeli", func(t *testing.T) {
		mockRepo.deleteImageFn = func(ctx context.Context, images []string) error {
			return errors.New("db delete error")
		}

		err := svc.DeleteMultipleImage(ctx, []string{"img1", "img2"}, "prod-1")
		assert.Error(t, err)
		assert.Contains(t, err.Error(), "Image Silinemedi")
	})

	t.Run("storage silme işlemi hata verirse hata dönmeli", func(t *testing.T) {
		mockRepo.deleteImageFn = func(ctx context.Context, images []string) error {
			return nil
		}
		mockFM.deleteMultipleFileFn = func(ctx context.Context, images []string) error {
			return errors.New("storage delete error")
		}

		err := svc.DeleteMultipleImage(ctx, []string{"img1", "img2"}, "prod-1")
		assert.Error(t, err)
		assert.Contains(t, err.Error(), "delete multiple da hata")
	})

	t.Run("başarıyla hem db hem storage'dan silmeli", func(t *testing.T) {
		dbDeleted := false
		mockRepo.deleteImageFn = func(ctx context.Context, images []string) error {
			dbDeleted = true
			assert.Equal(t, []string{"img1", "img2"}, images)
			return nil
		}

		storageDeleted := false
		mockFM.deleteMultipleFileFn = func(ctx context.Context, images []string) error {
			storageDeleted = true
			assert.Equal(t, []string{"img1", "img2"}, images)
			return nil
		}

		err := svc.DeleteMultipleImage(ctx, []string{"img1", "img2"}, "prod-1")
		assert.NoError(t, err)
		assert.True(t, dbDeleted)
		assert.True(t, storageDeleted)
	})
}

func TestImageService_DeleteImage(t *testing.T) {
	mockRepo := &mockImageRepoForImageService{}
	mockFM := &mockFileManagerForImageService{}
	svc := NewImageService(mockRepo, mockFM)
	ctx := context.Background()

	t.Run("db silme hata verirse hata dönmeli", func(t *testing.T) {
		mockRepo.deleteImageFn = func(ctx context.Context, images []string) error {
			return errors.New("db delete error")
		}

		err := svc.DeleteImage(ctx, "img-single")
		assert.Error(t, err)
		assert.Contains(t, err.Error(), "Image Silinemedi")
	})

	t.Run("başarıyla tek görseli silmeli", func(t *testing.T) {
		dbDeleted := false
		mockRepo.deleteImageFn = func(ctx context.Context, images []string) error {
			dbDeleted = true
			assert.Equal(t, []string{"img-single"}, images)
			return nil
		}

		storageDeleted := false
		mockFM.deleteFileFn = func(ctx context.Context, imageID string) error {
			storageDeleted = true
			assert.Equal(t, "-img-single", imageID)
			return nil
		}

		err := svc.DeleteImage(ctx, "img-single")
		assert.NoError(t, err)
		assert.True(t, dbDeleted)
		assert.True(t, storageDeleted)
	})
}

func TestImageService_DeleteAllImage(t *testing.T) {
	mockRepo := &mockImageRepoForImageService{}
	mockFM := &mockFileManagerForImageService{}
	svc := NewImageService(mockRepo, mockFM)
	ctx := context.Background()

	t.Run("ürün resim ID'leri alınamazsa hata dönmeli", func(t *testing.T) {
		mockRepo.getProductImagesFn = func(ctx context.Context, product_id string) ([]string, error) {
			return nil, errors.New("get images error")
		}

		err := svc.DeleteAllImage(ctx, "prod-1")
		assert.Error(t, err)
		assert.Contains(t, err.Error(), "image store da hata")
	})

	t.Run("storage çoklu silme hata verirse hata dönmeli", func(t *testing.T) {
		mockRepo.getProductImagesFn = func(ctx context.Context, product_id string) ([]string, error) {
			return []string{"img1", "img2"}, nil
		}
		mockFM.deleteMultipleFileFn = func(ctx context.Context, images []string) error {
			return errors.New("storage error")
		}

		err := svc.DeleteAllImage(ctx, "prod-1")
		assert.Error(t, err)
		assert.Contains(t, err.Error(), "storage da hata")
	})

	t.Run("db'den resimler silinirken hata çıkarsa hata dönmeli", func(t *testing.T) {
		mockRepo.getProductImagesFn = func(ctx context.Context, product_id string) ([]string, error) {
			return []string{"img1"}, nil
		}
		mockFM.deleteMultipleFileFn = func(ctx context.Context, images []string) error {
			return nil
		}
		mockRepo.deleteImageFn = func(ctx context.Context, images []string) error {
			return errors.New("delete image db error")
		}

		err := svc.DeleteAllImage(ctx, "prod-1")
		assert.Error(t, err)
		assert.Contains(t, err.Error(), "resimler silinirken hatayla karşılaşıldı")
	})

	t.Run("başarıyla tüm görselleri silmeli", func(t *testing.T) {
		mockRepo.getProductImagesFn = func(ctx context.Context, product_id string) ([]string, error) {
			return []string{"img1", "img2"}, nil
		}
		mockFM.deleteMultipleFileFn = func(ctx context.Context, images []string) error {
			assert.Equal(t, []string{"img1", "img2"}, images)
			return nil
		}
		mockRepo.deleteImageFn = func(ctx context.Context, images []string) error {
			assert.Equal(t, []string{"img1", "img2"}, images)
			return nil
		}

		err := svc.DeleteAllImage(ctx, "prod-1")
		assert.NoError(t, err)
	})
}

func TestImageService_CreateImages(t *testing.T) {
	mockRepo := &mockImageRepoForImageService{}
	mockFM := &mockFileManagerForImageService{}
	svc := NewImageService(mockRepo, mockFM)
	ctx := context.Background()

	fh1 := createTestFileHeader(t, "pic1.jpg", []byte("file-content-1"))
	fh2 := createTestFileHeader(t, "pic2.jpg", []byte("file-content-2"))

	t.Run("db kayıt hatası oluşursa storage rollback yapılmalı ve hata dönmeli", func(t *testing.T) {
		mockFM.uploadFileFn = func(fileName string, fileSize int64, prefix string, ctx context.Context, reader io.Reader) (string, error) {
			return fileName, nil
		}

		mockRepo.addOneImageFn = func(ctx context.Context, image models.Image) error {
			return errors.New("db insert error")
		}

		rollbackCalled := false
		mockFM.deleteFileFn = func(ctx context.Context, imageID string) error {
			rollbackCalled = true
			return nil
		}

		err := svc.CreateImages(ctx, []*multipart.FileHeader{fh1}, "prod-1", false, "")
		assert.Error(t, err)
		assert.Contains(t, err.Error(), "resimler db ye kayıt edilemedi")
		assert.True(t, rollbackCalled, "Hata durumunda storage rollback DeleteFile çağrılmalı")
	})

	t.Run("başarıyla görselleri yükleyip db'ye kaydetmeli", func(t *testing.T) {
		uploadedFiles := 0
		mockFM.uploadFileFn = func(fileName string, fileSize int64, prefix string, ctx context.Context, reader io.Reader) (string, error) {
			uploadedFiles++
			assert.Equal(t, "prod-1", prefix)
			return fileName, nil
		}

		dbSaved := 0
		mockRepo.addOneImageFn = func(ctx context.Context, image models.Image) error {
			dbSaved++
			assert.Equal(t, "prod-1", image.ProductID)
			return nil
		}

		err := svc.CreateImages(ctx, []*multipart.FileHeader{fh1, fh2}, "prod-1", false, "")
		assert.NoError(t, err)
		assert.Equal(t, 2, uploadedFiles)
		assert.Equal(t, 2, dbSaved)
	})

	t.Run("isChanged true olduğunda ilk görsel thumbnail olarak atanmalı ve db'ye eklenmeden geçilmeli", func(t *testing.T) {
		uploadedFiles := 0
		mockFM.uploadFileFn = func(fileName string, fileSize int64, prefix string, ctx context.Context, reader io.Reader) (string, error) {
			uploadedFiles++
			if uploadedFiles == 1 {
				assert.Equal(t, "thumb-123", fileName)
			}
			return fileName, nil
		}

		dbSaved := 0
		mockRepo.addOneImageFn = func(ctx context.Context, image models.Image) error {
			dbSaved++
			return nil
		}

		err := svc.CreateImages(ctx, []*multipart.FileHeader{fh1, fh2}, "prod-1", true, "thumb-123")
		assert.NoError(t, err)
		assert.Equal(t, 2, uploadedFiles)
		assert.Equal(t, 1, dbSaved, "isChanged=true iken ilk eleman thumbnail olduğu için db ye kaydedilmez")
	})
}

func TestImageService_AddThumbbailImage(t *testing.T) {
	mockRepo := &mockImageRepoForImageService{}
	mockFM := &mockFileManagerForImageService{}
	svc := NewImageService(mockRepo, mockFM)
	ctx := context.Background()

	fh := createTestFileHeader(t, "thumb.jpg", []byte("thumbnail-bytes"))

	t.Run("db kayıt hatasında storage rollback yapılmalı ve hata dönmeli", func(t *testing.T) {
		mockFM.uploadFileFn = func(fileName string, fileSize int64, prefix string, ctx context.Context, reader io.Reader) (string, error) {
			return fileName, nil
		}

		mockRepo.addOneImageFn = func(ctx context.Context, image models.Image) error {
			return errors.New("db insert fail")
		}

		rollbackCalled := false
		mockFM.deleteFileFn = func(ctx context.Context, imageID string) error {
			rollbackCalled = true
			return nil
		}

		err := svc.AddThumbbailImage(ctx, *fh, "prod-1")
		assert.Error(t, err)
		assert.Contains(t, err.Error(), "resimler db ye kayıt edilemedi")
		assert.True(t, rollbackCalled)
	})

	t.Run("başarıyla thumbnail görselini eklemeli", func(t *testing.T) {
		uploadCalled := false
		mockFM.uploadFileFn = func(fileName string, fileSize int64, prefix string, ctx context.Context, reader io.Reader) (string, error) {
			uploadCalled = true
			return fileName, nil
		}

		dbCalled := false
		mockRepo.addOneImageFn = func(ctx context.Context, image models.Image) error {
			dbCalled = true
			assert.Equal(t, "prod-1", image.ProductID)
			return nil
		}

		err := svc.AddThumbbailImage(ctx, *fh, "prod-1")
		assert.NoError(t, err)
		assert.True(t, uploadCalled)
		assert.True(t, dbCalled)
	})
}
