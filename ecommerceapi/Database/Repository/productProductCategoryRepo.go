package repository

import (
	"context"
	Entities "ecommerceapi/Database/Models"
	structs "ecommerceapi/Structs"
	"fmt"
	"log"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
)

type ProductProductCategoryFunctions interface {
	GetByProductID(ctx context.Context, productID uuid.UUID) ([]Entities.ProductProductCategory, error)
	GetByCategoryID(ctx context.Context, categoryID uuid.UUID) ([]Entities.ProductProductCategory, error)
	GetAllProductCategories(ctx context.Context, ids []string) ([]*structs.GetAllCategory, error)
	Add(ctx context.Context, item *Entities.ProductProductCategory) error
	Remove(ctx context.Context, productID uuid.UUID, categoryID uuid.UUID) error
	GetByProductIDOneRow(ctx context.Context, productID uuid.UUID) (uuid.UUID, error)
	AddMany(ctx context.Context, items []string, product_id string) error
	RemoveMany(ctx context.Context, items []string, product_id string) error
}

type ProductProductCategoryTable struct {
	db *pgxpool.Pool
}

func NewProductProductCategoryStore(dbpool *pgxpool.Pool) *ProductProductCategoryTable {
	return &ProductProductCategoryTable{
		db: dbpool,
	}
}

func (s *ProductProductCategoryTable) GetByProductID(ctx context.Context, productID uuid.UUID) ([]Entities.ProductProductCategory, error) {
	rows, err := s.db.Query(ctx, "Select * from ProductProductCategories where product_id=$1", productID)
	if err != nil {
		log.Printf("productProductCategoryRepo--GetByProductID--error:%v", err)
		return nil, fmt.Errorf("ürün kategori ilişkisi sorgu hatası: %w", err)
	}
	items, err := pgx.CollectRows(rows, pgx.RowToStructByName[Entities.ProductProductCategory])
	if err != nil {
		log.Printf("productProductCategoryRepo--GetByProductID--error:%v", err)
		return nil, fmt.Errorf("satır dönüşüm hatası: %w", err)
	}
	return items, nil
}

func (s *ProductProductCategoryTable) GetByCategoryID(ctx context.Context, categoryID uuid.UUID) ([]Entities.ProductProductCategory, error) {
	rows, err := s.db.Query(ctx, "Select * from ProductProductCategories where category_id=$1", categoryID)
	if err != nil {
		log.Printf("productProductCategoryRepo--GetByCategoryID--error:%v", err)
		return nil, fmt.Errorf("ürün kategori ilişkisi sorgu hatası: %w", err)
	}
	items, err := pgx.CollectRows(rows, pgx.RowToStructByName[Entities.ProductProductCategory])
	if err != nil {
		log.Printf("productProductCategoryRepo--GetByCategoryID--error:%v", err)
		return nil, fmt.Errorf("satır dönüşüm hatası: %w", err)
	}
	return items, nil
}

func (s *ProductProductCategoryTable) Add(ctx context.Context, item *Entities.ProductProductCategory) error {
	var Query string = `Insert into ProductProductCategories(product_id, category_id) Values ($1,$2)`
	_, err := s.db.Exec(ctx, Query,
		item.ProductID,
		item.CategoryID,
	)
	if err != nil {
		log.Printf("productProductCategoryRepo--Add--error:%v", err)
		return fmt.Errorf("ürün kategori ilişkisi oluşturma hatası: %w", err)
	}
	return nil
}

func (s *ProductProductCategoryTable) Remove(ctx context.Context, productID uuid.UUID, categoryID uuid.UUID) error {
	query := `Delete From ProductProductCategories where product_id=$1 AND category_id=$2`
	_, err := s.db.Exec(ctx, query, productID, categoryID)
	if err != nil {
		log.Printf("productProductCategoryRepo--Remove--error:%v", err)
		return fmt.Errorf("ürün kategori ilişkisi silme hatası: %w", err)
	}
	return nil
}

func (s *ProductProductCategoryTable) GetByProductIDOneRow(ctx context.Context, productID uuid.UUID) (uuid.UUID, error) {
	query := `select category_id from ProductProductCategories where product_id = $1`
	var category_id uuid.UUID
	err := s.db.QueryRow(ctx, query, productID).Scan(&category_id)

	if err != nil {
		log.Printf("productProductCategoryRepo--GetByProductIDOneRow--error:%v", err)
		return uuid.Nil, fmt.Errorf("category id çekilemedi %w", err)
	}

	return category_id, nil
}

func (s *ProductProductCategoryTable) AddMany(ctx context.Context, items []string, product_id string) error {

	// bulk insert
	_, err := s.db.CopyFrom(
		ctx,
		pgx.Identifier{"productproductcategories"},
		[]string{"product_id", "category_id"},

		// CopyFromSlice, veri dizisinin uzunluğunu ve her adımda
		// satır verisini döndüren bir callback (geri çağırma) fonksiyonu alır.
		pgx.CopyFromSlice(len(items), func(i int) ([]any, error) {
			return []any{
				product_id,
				items[i],
			}, nil
		}),
	)

	if err != nil {
		log.Printf("db error : productcategory/AddMany: %v", err)
		return fmt.Errorf("db error")
	}

	return nil
}

func (s *ProductProductCategoryTable) RemoveMany(ctx context.Context, items []string, product_id string) error {
	query := `delete from productproductcategories where product_id =$1 and category_id=Any($2)`

	_, err := s.db.Exec(ctx, query, product_id, items)

	if err != nil {
		log.Printf("productcategory silme hatası:%v", err)
		return err
	}

	return nil

}

// Verilen Product ID listesine göre kategorileri json formatında 2 sütun şeklinde getirir
// productid | json {category_id,category_name}
func (s *ProductProductCategoryTable) GetAllProductCategories(ctx context.Context, ids []string) ([]*structs.GetAllCategory, error) {

	query := `SELECT ppc.product_id,
	json_agg(
        json_build_object(
            'category_id', pc.category_id,
            'category_name', pc.category_name
        )) AS categories_json
	FROM productproductcategories ppc
	INNER JOIN productcategories pc ON ppc.category_id = pc.category_id
	WHERE ppc.product_id = ANY($1)
	GROUP BY ppc.product_id Order by product_id;`

	rows, err := s.db.Query(ctx, query, ids)

	if err != nil {
		log.Printf("get all category error %v", err)
		return nil, err
	}

	data, errCR := pgx.CollectRows(rows, pgx.RowToAddrOfStructByName[structs.GetAllCategory])
	if errCR != nil {
		log.Printf("get all category collect rows error:%v", errCR)
		return nil, errCR
	}

	for i := range data {
		if len(data[i].Categories) > 0 {
			log.Printf("debug session:ppc %s |%s", data[i].ProductId, data[i].Categories[0].CategoryName)
		}
	}

	return data, nil
}
