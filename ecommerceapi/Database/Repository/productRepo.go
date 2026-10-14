package repository

import (
	"context"
	constants "ecommerceapi/Constants"
	Models "ecommerceapi/Database/Models"
	st "ecommerceapi/Structs"
	structs "ecommerceapi/Structs"
	"fmt"
	"log"
	"time"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
)

// Product entity

// fonksiyonlar için interface

type ProductFunctions interface {
	Create(ctx context.Context, product *Models.Product) error
	Delete(ctx context.Context, productId *uuid.UUID) error
	Update(ctx context.Context, product *Models.Product) error
	Patch(ctx context.Context, product *Models.Product) error
	ChangeThumbnailImage(ctx context.Context, new_id string, product_id string) error
	// Get Methods
	GetAll(limit int, offset int, ctx context.Context) ([]structs.Getproducts, error)
	GetProductById(ctx context.Context, ProductID uuid.UUID) (*structs.ProductWithImage, error)
	GetMany(ctx context.Context, ids []string) ([]*structs.ProductForCart, error)
	GetProductCount(ctx context.Context) (int, error)
	GetProductsByCategoryId(ctx context.Context, category_ids []string, limit string, offset string) ([]structs.Getproducts, error)
	GetCountProductsByCategory(ctx context.Context, category_ids []string) (int, error)
	SearchProductCount(ctx context.Context, searchQuery string) int
	SearchProduct(ctx context.Context, searchQuery string, limit string) ([]st.Getproducts, error)
	AlternateSearch(ctx context.Context, searchQuery string, limit string) ([]st.Getproducts, error)
	FilteredAndSortedProducts(ctx context.Context, category_ids []string, limit, offset, sortBy, sortOrder string) ([]structs.Getproducts, error)
	// Admine Özel Get Metodları
	GetProductByIdForAdmin(ctx context.Context, ProductID uuid.UUID) (*structs.ProductWithImage, error)
	SearchProductCountAdmin(ctx context.Context, searchQuery string) int
	SearchProductAdmin(ctx context.Context, searchQuery string, limit string) ([]st.GetproductsAdmin, error)
	AlternateSearchAdmin(ctx context.Context, searchQuery string, limit string) ([]st.GetproductsAdmin, error)
	GetProductCountAdmin(ctx context.Context) (int, error)
	GetAllAdmin(limit int, offset int, ctx context.Context) ([]structs.GetproductsAdmin, error)
	ChangePublishState(ctx context.Context, product_id string, newState string) error
}

// db struct ı

type ProductTable struct {
	db *pgxpool.Pool
}

// Dışarıdan Erişim
// Constructor gibi çalışan foknsiyon
func NewProductStore(dbpool *pgxpool.Pool) *ProductTable {

	return &ProductTable{
		db: dbpool,
	}
}

// Concrete Functions

// Get Product Properties : id,name,price,thumbnail,stock,category
func (s *ProductTable) GetAll(limit int, offset int, ctx context.Context) ([]structs.Getproducts, error) {
	fmt.Println("Product Get All çalışıyor")

	var query string = `select p.product_id,p.product_name, p.price,p.thumbnail_image,p.stock ,p.created_at
							from products p where p.product_publish_state =$1 order by p.created_at desc limit $2 offset $3;`

	rows, err := s.db.Query(ctx, query, constants.PublishStatePublished, limit, offset)

	if err != nil {
		log.Printf("productRepo--GetAll--error:%v", err)
		return nil, fmt.Errorf("ürün sorgu hatası: %w", err)
	}

	// defer rows.Close() işlemini otomatik yapıyor
	products, err := pgx.CollectRows(rows, pgx.RowToStructByName[structs.Getproducts])

	if err != nil {
		log.Printf("productRepo--GetAll--error:%v", err)
		return nil, fmt.Errorf("satır dönüşüm hatası: %w", err)
	}
	return products, nil
}

func (s *ProductTable) Create(ctx context.Context, product *Models.Product) error {

	var Query string = `Insert into Products( 
	product_id,
	product_name, 
	stock, 
	thumbnail_image, 
	description, 
	price, 
	created_at, 
	updated_at,
	item_type,
	technical_details,
	product_code,
	product_publish_state) Values ($1,$2,$3,$4,$5,$6,$7,$8,$9,$10,$11,$12)`

	now := time.Now()
	product.CreatedAt = &now
	product.UpdatedAt = &now

	log.Printf("publish state:%s", product.ProductPublishState)

	if product.ProductPublishState == "" {
		product.ProductPublishState = constants.PublishStateUnpublished
	}

	_, err := s.db.Exec(ctx, Query,
		product.ProductID,
		product.ProductName,
		product.Stock,
		product.ThumbnailImage,
		product.Description,
		product.Price,
		product.CreatedAt,
		product.UpdatedAt,
		product.ItemType,
		product.TechnicalDetails,
		product.ProductCode,
		product.ProductPublishState)

	if err != nil {
		log.Printf("productRepo--Create--error:%v", err)
		return fmt.Errorf("ürün oluşturma hatası: %v", err)
	}
	fmt.Println("Yeni ürün ID'si:", product.ProductID)
	return nil
}

func (s *ProductTable) Delete(ctx context.Context, productId *uuid.UUID) error {

	query := `Delete From Products where product_id=$1`

	// AI Açıklaması:
	// Veri okumadığınız, sadece işlem yaptığınız durumlarda
	// Exec kullanılır. Exec işini yapar ve bağlantıyı hemen bırakır.
	a, err := s.db.Exec(ctx, query, productId)
	fmt.Println(a)

	if err != nil {
		log.Printf("productRepo--Delete--error:%v", err)
		return fmt.Errorf("ürün silme hatası: %w", err)
	}
	return nil
}

func (s *ProductTable) Update(ctx context.Context, product *Models.Product) error {
	var Query string = `Update products set product_name = $1,stock = $2,description = $3,
	price = $4,updated_at = $5,item_type=$6,technical_details=$7,product_publish_state=$8 WHERE product_id = $9`

	now := time.Now()
	product.UpdatedAt = &now

	_, queryerr := s.db.Exec(ctx, Query,
		product.ProductName,
		product.Stock,
		product.Description,
		product.Price,
		product.UpdatedAt,
		product.ItemType,
		product.TechnicalDetails,
		product.ProductPublishState,
		product.ProductID)

	if queryerr != nil {
		log.Printf("productRepo--Update--error:%v", queryerr)
		return fmt.Errorf("güncelleme sorgusu hatası: %w", queryerr)
	}
	return nil
}

// eğer tek bir endpoint yapsaydık admin yayınlanmamış ürünleri göremiyordu
func (s *ProductTable) GetProductByIdForAdmin(ctx context.Context, ProductID uuid.UUID) (*structs.ProductWithImage, error) {

	var query string = `SELECT 
    p.product_id, 
    p.product_name, 
    p.price,
	p.description,
	p.created_at,
	p.thumbnail_image,
	p.updated_at,
	p.stock,
	p.technical_details,
	p.product_code,
	p.product_publish_state,
	p.item_type,
    COALESCE(
        json_agg(
            json_build_object(
				'image_id', i.image_id,
                'path', i.image_path
            )
        ) FILTER (WHERE i.image_id IS NOT NULL), 
        '[]'
    ) AS images
FROM products p
LEFT JOIN images i ON p.product_id = i.product_id
where p.product_id = $1
GROUP BY p.product_id`

	rows, err := s.db.Query(ctx, query, ProductID)

	if err != nil {
		log.Printf("productRepo--GetProductById--Query--error:%v", err)
		return nil, fmt.Errorf("Veritabanı sorgusunda hatayla karşılaşıldı %w", err)
	}

	product, err := pgx.CollectOneRow(rows, pgx.RowToAddrOfStructByName[structs.ProductWithImage])

	if err != nil {
		log.Printf("productRepo--GetProductById--CollectOneRow--error:%v", err)
		return nil, fmt.Errorf("CollectOneRow işlemi sırasında hatayla karşılaşıldı %w", err)
	}
	return product, nil
}

// public endpoint : puslished ise ürünü gösterir
func (s *ProductTable) GetProductById(ctx context.Context, ProductID uuid.UUID) (*structs.ProductWithImage, error) {

	var query string = fmt.Sprintf(`SELECT 
    p.product_id, 
    p.product_name, 
    p.price,
	p.description,
	p.created_at,
	p.thumbnail_image,
	p.updated_at,
	p.stock,
	p.technical_details,
	p.product_code,
	p.product_publish_state,
	p.item_type,
    COALESCE(
        json_agg(
            json_build_object(
				'image_id', i.image_id,
                'path', i.image_path
            )
        ) FILTER (WHERE i.image_id IS NOT NULL), 
        '[]'
    ) AS images
FROM products p
LEFT JOIN images i ON p.product_id = i.product_id
where p.product_id = $1 and p.product_publish_state = '%s'
GROUP BY p.product_id;
`, constants.PublishStatePublished)
	rows, err := s.db.Query(ctx, query, ProductID)

	if err != nil {
		log.Printf("productRepo--GetProductById--Query--error:%v", err)
		return nil, fmt.Errorf("Veritabanı sorgusunda hatayla karşılaşıldı %w", err)
	}

	product, err := pgx.CollectOneRow(rows, pgx.RowToAddrOfStructByName[structs.ProductWithImage])

	if err != nil {
		log.Printf("productRepo--GetProductById--CollectOneRow--error:%v", err)
		return nil, fmt.Errorf("CollectOneRow işlemi sırasında hatayla karşılaşıldı %w", err)
	}
	return product, nil
}

func (s *ProductTable) Patch(ctx context.Context, product *Models.Product) error {
	var Query string = `Update products set product_name = $1,stock = $2,price = $3,updated_at = $4 WHERE product_id = $5`

	now := time.Now()
	product.UpdatedAt = &now

	_, queryerr := s.db.Exec(ctx, Query,
		product.ProductName,
		product.Stock,
		product.Price,
		product.UpdatedAt,
		product.ProductID)

	if queryerr != nil {
		log.Printf("productRepo--Patch--error:%v", queryerr)
		return fmt.Errorf("güncelleme sorgusu hatası: %w", queryerr)
	}
	return nil
}

// get many products with its thumbnails
// get for cart
func (s *ProductTable) GetMany(ctx context.Context, ids []string) ([]*structs.ProductForCart, error) {
	var query string = fmt.Sprintf(`SELECT 
    p.product_id, 
    p.product_name, 
    p.price,
    p.stock,
	p.thumbnail_image
    from products p where p.product_id = ANY($1::uuid[]) and p.product_publish_state = '%s'
	limit $2`, constants.PublishStatePublished)

	rows, err := s.db.Query(ctx, query, ids, len(ids))

	if err != nil {
		log.Printf("productRepo--GetMany--query--error:%v", err)
		return nil, fmt.Errorf("Veritabanı sorgusunda hatayla karşılaşıldı %w", err)
	}

	products, err := pgx.CollectRows(rows, pgx.RowToAddrOfStructByName[structs.ProductForCart])

	if err != nil {
		log.Printf("productRepo--GetMany--collectRows--error:%v", err)
		return nil, fmt.Errorf("CollectRows işlemi sırasında hatayla karşılaşıldı %w", err)
	}
	return products, nil
}

func (s *ProductTable) GetProductCount(ctx context.Context) (int, error) {

	var count int
	query := `select count(*) from products where product_publish_state =$1`

	err := s.db.QueryRow(ctx, query, constants.PublishStatePublished).Scan(&count)

	if err != nil {
		return -1, err
	}

	return count, nil

}

func (s *ProductTable) ChangeThumbnailImage(ctx context.Context, new_id string, product_id string) error {

	query := `Update products set thumbnail_image=$1 where product_id=$2`

	_, err := s.db.Exec(ctx, query, new_id, product_id)

	if err != nil {
		log.Printf("thumbnail güncellenemedi %v", err)
		return fmt.Errorf("Beklenmedik Bir Hatayla Karşılaşıldı")
	}

	return nil
}

// filter and sort products by category and price
// Params
// - sortBy: price, name, stock
// - sortOrder: asc, desc
func (s *ProductTable) FilteredAndSortedProducts(ctx context.Context, category_ids []string, limit, offset, sortBy, sortOrder string) ([]structs.Getproducts, error) {

	var mainQuery string
	var HasFilter bool
	var rows pgx.Rows
	var err error

	withCategoyFilter := fmt.Sprintf(`select p.product_id,p.product_name, p.price,p.thumbnail_image,p.stock,p.created_at from products p 
	where exists(
	select 1 from productproductcategories ppc where ppc.product_id =p.product_id 
	and ppc.category_id=any($1)) and p.product_publish_state = '%s'`, constants.PublishStatePublished)

	withoutFilter := fmt.Sprintf(`select p.product_id,p.product_name, p.price,p.thumbnail_image,p.stock,p.created_at from products p where p.product_publish_state = '%s'`, constants.PublishStatePublished)

	if len(category_ids) > 0 {
		mainQuery = withCategoyFilter
		HasFilter = true
	} else {
		mainQuery = withoutFilter
		HasFilter = false
	}

	orderClause := ""
	if dbCol, exists := constants.ValidSortColumns[sortBy]; exists {
		orderClause = fmt.Sprintf(" ORDER BY %s %s", dbCol, sortOrder)
	}

	if HasFilter {
		finalQuery := mainQuery + orderClause + " LIMIT $2 OFFSET $3;"
		rows, err = s.db.Query(ctx, finalQuery, category_ids, limit, offset)

		if err != nil {
			log.Printf("db error %v", err)
			return nil, err
		}
	} else {
		finalQuery := mainQuery + orderClause + " LIMIT $1 OFFSET $2;"
		rows, err = s.db.Query(ctx, finalQuery, limit, offset)

		if err != nil {
			log.Printf("db error %v", err)
			return nil, err
		}
	}

	data, err1 := pgx.CollectRows(rows, pgx.RowToStructByName[structs.Getproducts])

	if err1 != nil {
		log.Printf("collect rows error %v", err1)
		return nil, err1
	}

	return data, nil

}

func (s *ProductTable) GetProductsByCategoryId(ctx context.Context, category_ids []string, limit string, offset string) ([]structs.Getproducts, error) {

	query := fmt.Sprintf(`select p.product_id,p.product_name, p.price,p.thumbnail_image,p.stock,p.created_at from products p 
	where exists(
	select 1 from productproductcategories ppc where ppc.product_id =p.product_id 
	and ppc.category_id=any($1))and p.product_publish_state = '%s' limit $2 offset $3;`, constants.PublishStatePublished)

	rows, err := s.db.Query(ctx, query, category_ids, limit, offset)

	if err != nil {
		log.Printf("db error %v", err)
		return nil, err
	}

	data, err1 := pgx.CollectRows(rows, pgx.RowToStructByName[structs.Getproducts])

	if err1 != nil {
		log.Printf("collect rows error %v", err1)
		return nil, err1
	}

	return data, nil
}

func (s *ProductTable) GetCountProductsByCategory(ctx context.Context, category_ids []string) (int, error) {

	query := fmt.Sprintf(`select count(*) from products p 
	where exists(
	select 1 from productproductcategories ppc where ppc.product_id =p.product_id 
	and ppc.category_id=any($1)) and p.product_publish_state = '%s'`, constants.PublishStatePublished)

	var count int

	err := s.db.QueryRow(ctx, query, category_ids).Scan(&count)

	if err != nil {
		return -1, err
	}
	return count, nil

}

// Search Product Functions

func (s *ProductTable) SearchProductCount(ctx context.Context, searchQuery string) int {

	query := fmt.Sprintf(`Select count(*) from products where search_vector @@ websearch_to_tsquery('turkish',$1) and product_publish_state = '%s'`, constants.PublishStatePublished)

	var count int

	s.db.QueryRow(ctx, query, searchQuery).Scan(&count)

	return count

}

func (s *ProductTable) SearchProduct(ctx context.Context, searchQuery string, limit string) ([]st.Getproducts, error) {

	query := fmt.Sprintf(`select product_id,product_name,stock,price,thumbnail_image,created_at from products where search_vector @@ websearch_to_tsquery('turkish',$1) 
	and product_publish_state = '%s'
	order by ts_rank(search_vector,websearch_to_tsquery('turkish',$1)) desc limit $2`, constants.PublishStatePublished)

	rows, err := s.db.Query(ctx, query, searchQuery, limit)

	if err != nil {
		return nil, err
	}

	data, err1 := pgx.CollectRows(rows, pgx.RowToStructByName[st.Getproducts])

	if err1 != nil {
		return nil, err1
	}

	return data, nil

}

// using trigram similarity
func (s *ProductTable) AlternateSearch(ctx context.Context, searchQuery string, limit string) ([]st.Getproducts, error) {

	query := fmt.Sprintf(`SELECT product_id,product_name,stock,price,thumbnail_image,created_at FROM products where product_publish_state = '%s'
	ORDER BY product_name <-> $1 LIMIT $2`, constants.PublishStatePublished)

	rows, err := s.db.Query(ctx, query, searchQuery, limit)

	if err != nil {
		return nil, err
	}

	data, err1 := pgx.CollectRows(rows, pgx.RowToStructByName[st.Getproducts])

	if err1 != nil {
		return nil, err1
	}

	return data, nil

}

func (s *ProductTable) SearchProductCountAdmin(ctx context.Context, searchQuery string) int {

	query := `Select count(*) from products where search_vector @@ websearch_to_tsquery('turkish',$1)`

	var count int

	s.db.QueryRow(ctx, query, searchQuery).Scan(&count)

	return count

}

func (s *ProductTable) SearchProductAdmin(ctx context.Context, searchQuery string, limit string) ([]st.GetproductsAdmin, error) {

	query := `select product_id,product_name,stock,price,thumbnail_image,created_at,product_code,product_publish_state from products where search_vector @@ websearch_to_tsquery('turkish',$1) 
	order by ts_rank(search_vector,websearch_to_tsquery('turkish',$1)) desc limit $2`

	rows, err := s.db.Query(ctx, query, searchQuery, limit)

	if err != nil {
		return nil, err
	}

	data, err1 := pgx.CollectRows(rows, pgx.RowToStructByName[st.GetproductsAdmin])

	if err1 != nil {
		return nil, err1
	}

	return data, nil

}

func (s *ProductTable) AlternateSearchAdmin(ctx context.Context, searchQuery string, limit string) ([]st.GetproductsAdmin, error) {

	query := `SELECT product_id,product_name,stock,price,thumbnail_image,created_at,product_code,product_publish_state FROM products 
	ORDER BY product_name <-> $1 LIMIT $2`

	rows, err := s.db.Query(ctx, query, searchQuery, limit)

	if err != nil {
		return nil, err
	}

	data, err1 := pgx.CollectRows(rows, pgx.RowToStructByName[st.GetproductsAdmin])

	if err1 != nil {
		return nil, err1
	}

	return data, nil

}

func (s *ProductTable) GetProductCountAdmin(ctx context.Context) (int, error) {

	var count int
	query := `select count(*) from products`

	err := s.db.QueryRow(ctx, query).Scan(&count)

	if err != nil {
		return -1, err
	}

	return count, nil

}

func (s *ProductTable) GetAllAdmin(limit int, offset int, ctx context.Context) ([]structs.GetproductsAdmin, error) {
	fmt.Println("Product GetAllAdmin çalışıyor")

	var query string = `select p.product_id,p.product_name, p.price,p.thumbnail_image,p.stock ,p.created_at, p.product_code, p.product_publish_state
							from products p order by p.created_at desc limit $1 offset $2;`

	rows, err := s.db.Query(ctx, query, limit, offset)

	if err != nil {
		log.Printf("productRepo--GetAllAdmin--error:%v", err)
		return nil, fmt.Errorf("ürün sorgu hatası: %w", err)
	}

	products, err := pgx.CollectRows(rows, pgx.RowToStructByName[structs.GetproductsAdmin])

	if err != nil {
		log.Printf("productRepo--GetAllAdmin--error:%v", err)
		return nil, fmt.Errorf("satır dönüşüm hatası: %w", err)
	}
	return products, nil
}

func (s *ProductTable) ChangePublishState(ctx context.Context, product_id string, newState string) error {
	query := `Update products set product_publish_state=$1, updated_at=$2 where product_id=$3`

	_, err := s.db.Exec(ctx, query, newState, time.Now(), product_id)

	if err != nil {
		log.Printf("productRepo--ChangePublishState--error:%v", err)
		return fmt.Errorf("ürün yayın durumu güncellenemedi: %w", err)
	}

	return nil
}
