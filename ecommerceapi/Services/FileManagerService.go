package services

import (
	"context"
	"fmt"
	"io"
	"log"
	"net/url"
	"time"

	st "ecommerceapi/Structs"

	"github.com/minio/minio-go/v7"
)

type FileManagerService struct {
	minioClient *minio.Client
	bucketName  string
}

func NewFileManagerService(Client *minio.Client, _bucketName string) *FileManagerService {
	return &FileManagerService{
		minioClient: Client,
		bucketName:  _bucketName,
	}
}

func (f FileManagerService) UploadFile(fileName string, fileSize int64, prefix string, ctx context.Context, reader io.Reader) (string, error) {

	//fileName = fmt.Sprintf("%s/%s", prefix, fileName)

	_, err := f.minioClient.PutObject(
		ctx,
		f.bucketName,
		fileName,
		reader,
		fileSize,
		minio.PutObjectOptions{ContentType: "image/webp"})

	if err != nil {
		return "", fmt.Errorf("storage'a yüklerken hatayla karşılaşıldı %w ", err)
	}

	return fileName, nil
}

func (f FileManagerService) DeleteFile(ctx context.Context, image_id string) error {

	err := f.minioClient.RemoveObject(ctx, f.bucketName, image_id, minio.RemoveObjectOptions{})

	if err != nil {
		return fmt.Errorf("storage remove işleminde hata %w", err)
	}

	return nil

}

// prefix parametresi imagelerin prefix'i olan klasörü simgeler
func (f FileManagerService) DeleteMultipleFile(ctx context.Context, images []string) error {

	// Silinecek nesneler için bir channel oluşturuyoruz
	objectsCh := make(chan minio.ObjectInfo)

	// Channel'a veri gönderme işlemini bir goroutine içinde yapmalıyız
	go func() {
		defer close(objectsCh)
		// Sadece sizin verdiğiniz listedeki dosyaları channel'a ekliyoruz
		for _, name := range images {
			objectsCh <- minio.ObjectInfo{
				Key: name,
			}
		}
	}()

	// RemoveObjects çoklu silme işlemi yapar
	opts := minio.RemoveObjectsOptions{
		GovernanceBypass: true, // Eğer object locking varsa bypass etmek için
	}

	// Hata dönüşlerini yakalamak için loop
	errorCh := f.minioClient.RemoveObjects(ctx, f.bucketName, objectsCh, opts)

	// Eğer silme sırasında hata oluşan dosyalar varsa burada listelenir
	for e := range errorCh {
		fmt.Printf("Hata oluştu - Dosya: %s, Hata: %s\n", e.ObjectName, e.Err)
	}

	return nil

}

// geçici url
func (f FileManagerService) GetPresignedFileObject() (st.UrlResult, error) {

	expireTime := time.Second * 24 * 60 * 60
	// Set request parameters for content-disposition.
	reqParams := make(url.Values)
	reqParams.Set("response-content-disposition", "inline; filename=\"a.jpeg\"")

	// Generates a presigned url which expires in a day.
	presignedURL, err := f.minioClient.PresignedGetObject(context.Background(), "go-ecommerce", "a.jpeg", expireTime, reqParams)
	if err != nil {
		log.Printf("MinIO PresignedGetObject hatası: %v", err)

		// HATA DÖNDÜRÜLÜRKEN: %w kullanarak hatayı sarmalayın (wrap), böylece üst katman asıl sorunu görebilir.
		return st.UrlResult{}, fmt.Errorf("dosya servisinde imzalama hatası: %w", err)
	}
	createdDate := time.Now()

	Result := st.UrlResult{
		CreatedDate:  createdDate,
		ExpireDate:   expireTime,
		PresignedURL: presignedURL,
	}

	fmt.Println("Successfully generated presigned URL", presignedURL)
	return Result, nil
}
