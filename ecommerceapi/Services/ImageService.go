package services

import (
	"context"
	models "ecommerceapi/Database/Models"
	repo "ecommerceapi/Database/Repository"
	st "ecommerceapi/Structs"
	"fmt"
	"io"
	"log"
	"mime/multipart"

	"github.com/google/uuid"
)

type FileManagerFunctions interface {
	UploadFile(fileName string, fileSize int64, prefix string, ctx context.Context, reader io.Reader) (string, error)
	DeleteFile(ctx context.Context, image_id string) error
	DeleteMultipleFile(ctx context.Context, images []string) error
	GetPresignedFileObject() (st.UrlResult, error)
}

type ImageService struct {
	imageStore  repo.ImageFunctions
	fileManager FileManagerFunctions
}

func NewImageService(imageStore repo.ImageFunctions, _fileManager FileManagerFunctions) *ImageService {

	return &ImageService{
		imageStore:  imageStore,
		fileManager: _fileManager,
	}
}

func (img *ImageService) DeleteMultipleImage(ctx context.Context, images []string, product_id string) error {

	err := img.imageStore.DeleteImage(ctx, images)

	if err != nil {
		return fmt.Errorf("Image Silinemedi %w", err)
	}

	//storagedan silme
	err = img.fileManager.DeleteMultipleFile(ctx, images)

	if err != nil {
		return fmt.Errorf("delete multiple da hata %w", err)
	}

	return nil
}

func (img *ImageService) DeleteImage(ctx context.Context, image_id string) error {

	temp := []string{image_id}
	err := img.imageStore.DeleteImage(ctx, temp)

	if err != nil {
		return fmt.Errorf("Image Silinemedi %w", err)
	}

	image_id = fmt.Sprintf("%s-%s", "", image_id)
	//storagedan silme
	img.fileManager.DeleteFile(ctx, image_id)

	return nil
}

func (img *ImageService) CreateImages(ctx context.Context, images []*multipart.FileHeader, product_id string, isChanged bool, thumbnail_id string) error {

	for i, image := range images {

		file, err := image.Open()
		file_size := image.Size
		var image_id string

		if err != nil {
			return fmt.Errorf("Dosya Açılmama Hatası %w", err)
		}

		if i == 0 && isChanged {
			image_id = thumbnail_id
		} else {
			image_id = uuid.NewString()
		}

		path, err2 := img.fileManager.UploadFile(image_id, file_size, product_id, context.Background(), file)

		file.Close()
		if err2 != nil || len(path) == 0 {
			i--
			log.Printf("image %d yüklenirken hatayla karşılaşıldı  %v", i, err2)
			continue
		}

		// resim db kaydı oluştur
		image := models.Image{
			ImageID:   image_id,
			ImagePath: path,
			ProductID: product_id,
		}

		if i == 0 && isChanged {
			continue
		}
		err = img.imageStore.AddOneImage(ctx, image)

		if err != nil {
			// storage rollback mekanizması
			img.fileManager.DeleteFile(ctx, image_id)
			return fmt.Errorf("resimler db ye kayıt edilemedi %w", err)
		}

	}

	return nil
}

func (img *ImageService) DeleteAllImage(ctx context.Context, product_id string) error {

	image_ids, err := img.imageStore.GetProductImages(ctx, product_id)

	if err != nil {
		return fmt.Errorf("image store da hata %w", err)
	}

	// delete from storage
	err = img.fileManager.DeleteMultipleFile(ctx, image_ids)

	if err != nil {
		return fmt.Errorf("storage da hata %w", err)
	}

	err = img.imageStore.DeleteImage(ctx, image_ids)

	if err != nil {
		return fmt.Errorf("resimler silinirken hatayla karşılaşıldı %w", err)
	}

	return nil
}

func (img *ImageService) AddThumbbailImage(ctx context.Context, image multipart.FileHeader, product_id string) error {

	file, err := image.Open()
	file_size := image.Size

	if err != nil {
		return fmt.Errorf("Dosya Açılmama Hatası %w", err)
	}

	image_id := uuid.NewString()

	path, err2 := img.fileManager.UploadFile(image_id, file_size, product_id, context.Background(), file)

	file.Close()
	if err2 != nil || len(path) == 0 {

		log.Printf("image  yüklenirken hatayla karşılaşıldı  %v", err2)

	}

	// resim db kaydı oluştur
	image_new := models.Image{
		ImageID:   image_id,
		ImagePath: path,
		ProductID: product_id,
	}

	err = img.imageStore.AddOneImage(ctx, image_new)

	if err != nil {
		// storage rollback mekanizması
		img.fileManager.DeleteFile(ctx, image_id)
		return fmt.Errorf("resimler db ye kayıt edilemedi %w", err)
	}

	return nil
}
