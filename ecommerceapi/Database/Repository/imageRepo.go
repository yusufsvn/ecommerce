package repository

import (
	"context"
	models "ecommerceapi/Database/Models"
	"fmt"
	"log"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
)

type ImageFunctions interface {
	AddOneImage(ctx context.Context, image models.Image) error
	DeleteImage(ctx context.Context, images []string) error
	GetProductImages(ctx context.Context, product_id string) ([]string, error)
}

type ImageTable struct {
	db *pgxpool.Pool
}

func NewImageStore(dbpool *pgxpool.Pool) *ImageTable {
	return &ImageTable{
		db: dbpool,
	}
}

func (s *ImageTable) AddOneImage(ctx context.Context, image models.Image) error {

	_, err := s.db.Exec(
		ctx,
		"insert into images (image_id,image_path,product_id) values ($1,$2,$3)",
		image.ImageID,
		image.ImagePath,
		image.ProductID)

	if err != nil {
		log.Printf("imageRepo--AddOneImage--error:%v", err)
		return fmt.Errorf("Database Error -Creating Image %w", err)
	}

	return nil

}

func (s *ImageTable) DeleteImage(ctx context.Context, images []string) error {

	_, err := s.db.Exec(ctx, "delete from images where image_id = any($1)", images)

	if err != nil {
		log.Printf("imageRepo--DeleteImage--error:%v", err)
		return fmt.Errorf("Image Table da db hatası %w", err)
	}

	return nil

}

func (s *ImageTable) GetProductImages(ctx context.Context, product_id string) ([]string, error) {

	var query string = "Select image_id from images where product_id =$1"

	rows, err := s.db.Query(ctx, query, product_id)

	if err != nil {
		log.Printf("imageRepo--GetProductImages--error:%v", err)
		return nil, fmt.Errorf("select image sorgusunda hata from imageRepo.GetProductImages %w", err)
	}

	image_ids, err := pgx.CollectRows(rows, pgx.RowTo[string])

	if err != nil {
		log.Printf("imageRepo--GetProductImages--error:%v", err)
		return nil, fmt.Errorf("collect rows da hata oluştu from imageRepo.GetProductImages %w", err)
	}

	return image_ids, nil

}
