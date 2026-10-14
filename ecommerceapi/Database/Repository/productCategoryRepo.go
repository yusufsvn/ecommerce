package repository

import (
	"context"
	structs "ecommerceapi/Structs"
	"fmt"
	"log"
	"time"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
)

type ProductCategoryFunctions interface {
	GetAll(ctx context.Context) ([]structs.GetCategory, error)
	Create(ctx context.Context, category_name string) (structs.GetCategory, error)
	Delete(ctx context.Context, categoryId string) error
	GetByCategoryName(ctx context.Context, categoryName string) (uuid.UUID, error)
	GetCategoriesByProductId(ctx context.Context, product_id string) ([]structs.GetCategory, error)
	GetByCategoryIds(ctx context.Context, categoryNames []string) ([]string, error)
	Update(ctx context.Context, category_id string, category_name string) error
}

type ProductCategoryTable struct {
	db *pgxpool.Pool
}

func NewProductCategoryStore(dbpool *pgxpool.Pool) *ProductCategoryTable {
	return &ProductCategoryTable{
		db: dbpool,
	}
}

func (s *ProductCategoryTable) GetAll(ctx context.Context) ([]structs.GetCategory, error) {
	rows, err := s.db.Query(ctx, "Select category_id,category_name from ProductCategories")
	if err != nil {
		log.Printf("productCategoryRepo--GetAll--error:%v", err)
		return nil, fmt.Errorf("ürün kategori sorgu hatası: %w", err)
	}
	categories, err := pgx.CollectRows(rows, pgx.RowToStructByName[structs.GetCategory])
	if err != nil {
		log.Printf("productCategoryRepo--GetAll--error:%v", err)
		return nil, fmt.Errorf("satır dönüşüm hatası: %w", err)
	}
	return categories, nil
}

func (s *ProductCategoryTable) Create(ctx context.Context, category_name string) (structs.GetCategory, error) {
	now := time.Now()
	id := uuid.NewString()

	var Query string = `Insert into ProductCategories(category_id,category_name,updated_at) Values ($1,$2,$3)`
	_, err := s.db.Exec(ctx, Query, id, category_name, now)
	if err != nil {
		log.Printf("productCategoryRepo--Create--error:%v", err)
		return structs.GetCategory{}, fmt.Errorf("ürün kategori oluşturma hatası: %w", err)
	}

	result := structs.GetCategory{
		Category_id:   id,
		Category_name: category_name,
	}

	return result, nil
}

func (s *ProductCategoryTable) Delete(ctx context.Context, categoryId string) error {
	query := `Delete From ProductCategories where category_id=$1`
	_, err := s.db.Exec(ctx, query, categoryId)
	if err != nil {
		log.Printf("productCategoryRepo--Delete--error:%v", err)
		return fmt.Errorf("ürün kategori silme hatası: %w", err)
	}
	return nil
}

// Veritabanında EnumNo ile belirlenmiş bir sütun için function
func (s *ProductCategoryTable) GetByCategoryName(ctx context.Context, categoryName string) (uuid.UUID, error) {
	var categoryId uuid.UUID
	var query string = `Select category_id from ProductCategories where category_name = $1`

	err := s.db.QueryRow(ctx, query, categoryName).Scan(&categoryId)

	if err != nil {
		log.Printf("productCategoryRepo--GetByCategoryName--error:%v", err)
		return uuid.Nil, fmt.Errorf("select sorgusunda hatayla karsilasildi %w", err)
	}
	return categoryId, nil

}

func (s *ProductCategoryTable) GetCategoriesByProductId(ctx context.Context, product_id string) ([]structs.GetCategory, error) {

	query := `Select category_name,pc.category_id from productcategories pc 
				left join productproductcategories ppc on ppc.product_id = $1
				where pc.category_id = ppc.category_id;`

	rows, err := s.db.Query(ctx, query, product_id)

	if err != nil {
		log.Printf("pc repo da hata :%v", err)
		return nil, err
	}

	data, err1 := pgx.CollectRows(rows, pgx.RowToStructByName[structs.GetCategory])

	if err1 != nil {
		log.Printf("collect rows error :%v", err1)
		return nil, err1
	}

	return data, nil

}

func (s *ProductCategoryTable) GetByCategoryIds(ctx context.Context, categoryNames []string) ([]string, error) {

	var query string = `select category_id from productcategories where category_name=any($1);`

	rows, err := s.db.Query(ctx, query, categoryNames)

	if err != nil {
		log.Printf("productCategoryRepo--GetByCategoryName--error:%v", err)
		return nil, fmt.Errorf("select sorgusunda hatayla karsilasildi %w", err)
	}

	categories, err1 := pgx.CollectRows(rows, pgx.RowTo[string])

	if err1 != nil {
		log.Printf("productCategoryRepo--GetByCategoryNames--Collect rows--:%v", err1)
		return nil, err1
	}

	return categories, nil

}

func (s *ProductCategoryTable) Update(ctx context.Context, category_id string, category_name string) error {
	query := `update productcategories set category_name=$1 where category_id =$2`

	_, err := s.db.Exec(ctx, query, category_name, category_id)

	if err != nil {
		return err
	}
	return nil
}
