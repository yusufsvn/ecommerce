package services

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"mime/multipart"
	"strings"

	"github.com/google/uuid"

	constants "ecommerceapi/Constants"
	Customerrors "ecommerceapi/CustomErrors"
	models "ecommerceapi/Database/Models"
	repo "ecommerceapi/Database/Repository"
	structs "ecommerceapi/Structs"
	"ecommerceapi/Utils"
)

type ProductService struct {
	ProductStore         repo.ProductFunctions
	ProductCategoryStore repo.ProductProductCategoryFunctions
	CategoryStore        repo.ProductCategoryFunctions
	StorageAPIUrl        string
	storagePublicUrl     string
	BucketName           string
	Image_service        ImageService
	OrderService         OrderService
}

func NewProductService(
	productStore repo.ProductFunctions,
	categoryStore repo.ProductCategoryFunctions,
	productproductcategoryStore repo.ProductProductCategoryFunctions,
	_storageAPIUrl string,
	_storagePublicUrl string,
	_bucketName string,
	_image_service ImageService,
	_orderService OrderService,
) *ProductService {

	return &ProductService{
		ProductStore:         productStore,
		ProductCategoryStore: productproductcategoryStore,
		CategoryStore:        categoryStore,
		StorageAPIUrl:        _storageAPIUrl,
		storagePublicUrl:     _storagePublicUrl,
		BucketName:           _bucketName,
		Image_service:        _image_service,
		OrderService:         _orderService,
	}
}

// create product by admin
func (p *ProductService) CreateProduct(ctx context.Context, request *structs.CreateProductRequest) error {

	product_id := uuid.NewString()
	thumbnail_id := uuid.NewString()

	// add control item type PHYSICAL or VIRTUAL
	switch request.ItemType {
	case constants.PhysicalProduct:
	case constants.VirtualProduct:
	default:
		return Customerrors.ErrInvalidItemType
	}

	TechDetailsBytes, errtd := json.Marshal(request.TechnicalDetails)

	if errtd != nil {
		return errtd
	}

	productCode, err := Utils.CreateNanoId(8)
	if err != nil {
		return fmt.Errorf("ürün kodu oluşturulamadı: %w", err)
	}

	publishState := request.PublishState
	log.Printf("publish state:%s", request.PublishState)
	if publishState == "" {
		publishState = constants.PublishStateUnpublished
	}

	// init product
	product := models.Product{
		ProductID:           product_id,
		ProductName:         request.ProductName,
		Stock:               request.Stock,
		Price:               request.Price,
		Description:         request.Description,
		ItemType:            request.ItemType,
		ThumbnailImage:      thumbnail_id,
		TechnicalDetails:    TechDetailsBytes,
		ProductCode:         productCode,
		ProductPublishState: publishState,
	}

	// product datasını interface ile veritabanı fonksiyonlarına erişip create oluşturuyor.
	errk := p.ProductStore.Create(ctx, &product)

	if errk != nil {
		return fmt.Errorf("Product oluşturulamadı %w", errk)
	}

	// image process
	imgerr := p.Image_service.CreateImages(ctx, request.Images, product_id, true, thumbnail_id)

	if imgerr != nil {
		return fmt.Errorf("image servisinde hata create image error %w", imgerr)
	}

	// kategori adı verilip idsi alınıyor
	category_ids, err := p.CategoryStore.GetByCategoryIds(ctx, request.ProductCategory)

	if err != nil && category_ids == nil {
		return fmt.Errorf("Category id cannot fetched %w", err)
	}

	categoryLength := len(category_ids)

	if categoryLength == 1 {
		aratablomodel := models.ProductProductCategory{
			ProductID:  product_id,
			CategoryID: category_ids[0],
		}
		err2 := p.ProductCategoryStore.Add(ctx, &aratablomodel)
		if err2 != nil {
			return fmt.Errorf("ara tabloya eklenirken hata oluştu %w", err2)
		}
	}
	if categoryLength > 1 {
		errCT1 := p.ProductCategoryStore.AddMany(ctx, category_ids, product_id)

		if errCT1 != nil {
			log.Printf("producut service/CreateProduct line 100:%v", errCT1)
			return errCT1
		}
	}

	return nil
}

// product list for admin panel
func (p *ProductService) GetAllProducts(ctx context.Context, limit int, offset int) (*structs.ProductResultAdmin, error) {

	totaldata, err1 := p.ProductStore.GetProductCountAdmin(ctx)

	if err1 != nil {
		return nil, err1
	}

	if totaldata == 0 {
		return &structs.ProductResultAdmin{
			Count: totaldata,
			Data:  []structs.GetAllProductResponseAdmin{},
		}, nil
	}

	data, err := p.ProductStore.GetAllAdmin(limit, offset, ctx)
	if err != nil {
		return nil, fmt.Errorf("GetAllProducts Service de db hatası %w", err)
	}

	var productIds []string

	for i := range data {
		productIds = append(productIds, data[i].ProductID)

	}

	categories, errCT := p.ProductCategoryStore.GetAllProductCategories(ctx, productIds)

	if errCT != nil {
		return nil, errCT
	}

	var productData []structs.GetAllProductResponseAdmin

	for i := range data {
		data[i].ThumbnailImage = Utils.MatchWithURLSingle(p.storagePublicUrl, p.BucketName, data[i].ThumbnailImage)
		for j := range categories {
			if data[i].ProductID == categories[j].ProductId {
				temp := structs.GetAllProductResponseAdmin{
					Product:    data[i],
					Categories: categories[j].Categories,
				}
				productData = append(productData, temp)
			}
		}
	}

	result := structs.ProductResultAdmin{
		Count: totaldata,
		Data:  productData,
	}

	return &result, nil
}

// product list by category for public list
func (p *ProductService) PubGetAllProducts(ctx context.Context, limit int, offset int) (*structs.CategoryFilteredProducts, error) {

	totaldata, err1 := p.ProductStore.GetProductCount(ctx)

	if err1 != nil {
		return nil, err1
	}

	if totaldata != 0 {
		data, err := p.ProductStore.GetAll(limit, offset, ctx)
		if err != nil {
			return nil, err
		}

		for i := range data {
			data[i].ThumbnailImage = Utils.MatchWithURLSingle(p.storagePublicUrl, p.BucketName, data[i].ThumbnailImage)
		}

		response := structs.CategoryFilteredProducts{
			Count: totaldata,
			Data:  data,
		}

		return &response, nil
	}

	return nil, nil

}

func (p *ProductService) GetProductById(ctx context.Context, ProductID string) (*structs.FullProductResponse, error) {

	ProductIDuuid, err := uuid.Parse(ProductID)

	if err != nil {
		return nil, fmt.Errorf("uuid dönüştürülürken hatayla karşılaşıldı %w", err)
	}
	// db call
	product, err := p.ProductStore.GetProductById(ctx, ProductIDuuid)

	if err != nil {
		return nil, fmt.Errorf("Product Verisi Alınamadı %w", err)
	}

	// get product category
	categories, errCT := p.CategoryStore.GetCategoriesByProductId(ctx, ProductID)

	if errCT != nil {
		log.Printf("Kategoriler getirilemedi :%v", errCT)
		return nil, errCT
	}

	// response model mapping

	thumbId, err1 := uuid.Parse(product.ThumbnailImage)

	if err1 != nil {
		return nil, fmt.Errorf("thumbnail uuid geçersiz formatta: %w", err)
	}

	thumbnailobj := structs.ProductImage{
		ImageID: thumbId,
		Path:    product.ThumbnailImage,
	}

	allImages := make([]structs.ProductImage, 0, len(product.Images)+1)

	allImages = append(allImages, thumbnailobj)
	allImages = append(allImages, product.Images...)

	techDetails := make(map[string]string)

	errdecode := json.Unmarshal(product.TechnicalDetails, &techDetails)
	if errdecode != nil {
		return nil, errdecode
	}

	productResponse := structs.FullProductResponse{
		ProductID:           product.ProductID,
		ProductName:         product.ProductName,
		Stock:               product.Stock,
		Description:         product.Description,
		Price:               product.Price,
		Images:              Utils.MatchWithURL(p.storagePublicUrl, p.BucketName, allImages),
		Categories:          categories,
		TechnicalDetails:    techDetails,
		ProductCode:         product.ProductCode,
		ProductPublishState: product.ProductPublishState,
		ItemType:            product.ItemType,
	}

	return &productResponse, nil
}

func (p *ProductService) GetProductByIdForAdmin(ctx context.Context, ProductID string) (*structs.FullProductResponse, error) {

	ProductIDuuid, err := uuid.Parse(ProductID)

	if err != nil {
		return nil, fmt.Errorf("uuid dönüştürülürken hatayla karşılaşıldı %w", err)
	}
	// db call
	product, err := p.ProductStore.GetProductByIdForAdmin(ctx, ProductIDuuid)

	if err != nil {
		return nil, fmt.Errorf("Product Verisi Alınamadı %w", err)
	}

	// get product category
	categories, errCT := p.CategoryStore.GetCategoriesByProductId(ctx, ProductID)

	if errCT != nil {
		log.Printf("Kategoriler getirilemedi :%v", errCT)
		return nil, errCT
	}

	// response model mapping

	thumbId, err1 := uuid.Parse(product.ThumbnailImage)

	if err1 != nil {
		return nil, fmt.Errorf("thumbnail uuid geçersiz formatta: %w", err)
	}

	thumbnailobj := structs.ProductImage{
		ImageID: thumbId,
		Path:    product.ThumbnailImage,
	}

	allImages := make([]structs.ProductImage, 0, len(product.Images)+1)

	allImages = append(allImages, thumbnailobj)
	allImages = append(allImages, product.Images...)

	techDetails := make(map[string]string)

	errdecode := json.Unmarshal(product.TechnicalDetails, &techDetails)
	if errdecode != nil {
		return nil, errdecode
	}

	productResponse := structs.FullProductResponse{
		ProductID:           product.ProductID,
		ProductName:         product.ProductName,
		Stock:               product.Stock,
		Description:         product.Description,
		Price:               product.Price,
		Images:              Utils.MatchWithURL(p.storagePublicUrl, p.BucketName, allImages),
		Categories:          categories,
		TechnicalDetails:    techDetails,
		ProductCode:         product.ProductCode,
		ProductPublishState: product.ProductPublishState,
		ItemType:            product.ItemType,
	}

	return &productResponse, nil
}

func (p *ProductService) DeleteProduct(ctx context.Context, ProductID string) error {

	id, err := uuid.Parse(ProductID)

	if err != nil {
		log.Println("Fatal: uuid convert error")
		return fmt.Errorf("uuid dönüşümünde hatayla karşılaşıldı %w", err)
	}

	// image id lerini getirtme
	err1 := p.Image_service.DeleteAllImage(ctx, ProductID)

	if err1 != nil {
		return fmt.Errorf("imageler silinirken hatayla karşılaşıldı %w", err1)
	}

	// db query
	category_id, err2 := p.ProductCategoryStore.GetByProductIDOneRow(ctx, id)

	if err2 != nil {
		//log.Fatalln("Fatal: DB ERROR")
		return fmt.Errorf("kategori id getirme hatası %w", err)
	}

	err3 := p.ProductCategoryStore.Remove(ctx, id, category_id)

	if err3 != nil {
		log.Println("Product, Category Kayıtlarından Silinemedi")
		return err3
	}

	err4 := p.OrderService.RemoveByProductId(ctx, &id)
	if err4 != nil {
		log.Printf("Siparişlerden Silinemedi %v", err4)
		return err4
	}

	err5 := p.ProductStore.Delete(ctx, &id)

	if err5 != nil {
		log.Println("Fatal: DB ERROR")
		return fmt.Errorf("Db silme sorgusunda hata %w", err5)
	}

	return nil
}

func (p *ProductService) UpdateProduct(ctx context.Context, data *structs.UpdateProductRequest) error {

	if data == nil {
		return fmt.Errorf("boş veri gönderildi ")
	}

	TechDetailsBytes, errtd := json.Marshal(data.TechnicalDetails)

	if errtd != nil {
		return errtd
	}

	product := models.Product{
		ProductID:           data.ProductID,
		ProductName:         data.ProductName,
		Stock:               data.Stock,
		Description:         data.Description,
		Price:               data.Price,
		TechnicalDetails:    TechDetailsBytes,
		ProductCode:         data.ProductCode,
		ProductPublishState: data.ProductPublishState,
	}

	log.Printf("Log 1. ItemType: %s, PublishState: %s, productCode:%s ", data.ItemType, data.ProductPublishState, data.ProductCode)
	if data.ItemType == "" {
		data.ItemType = constants.PhysicalProduct
	}

	log.Printf("Log 2: İtemtypeControl. ItemType: %s, PublishState: %s, productCode:%s ", data.ItemType, data.ProductPublishState, data.ProductCode)

	err := p.ProductStore.Update(ctx, &product)

	if len(data.DeletedImages) != 0 {
		p.Image_service.DeleteMultipleImage(ctx, data.DeletedImages, data.ProductID)
	}

	if err != nil {
		return fmt.Errorf("db sorgusunda hata %w", err)
	}

	// isimden kategori id getirme
	newCategoryIds, errCat := p.CategoryStore.GetByCategoryIds(ctx, data.Categories.NewCategories)

	if errCat != nil {
		return errCat
	}

	delCategoryIds, errCat2 := p.CategoryStore.GetByCategoryIds(ctx, data.Categories.DeletedCategories)

	if errCat2 != nil {
		return errCat2
	}

	// Yeni Kategori Ekleme
	if len(data.Categories.NewCategories) > 0 {
		errCT1 := p.ProductCategoryStore.AddMany(ctx, newCategoryIds, data.ProductID)

		if errCT1 != nil {
			return errCT1
		}
	}

	// Kategori Silme
	if len(data.Categories.DeletedCategories) > 0 {
		errCT2 := p.ProductCategoryStore.RemoveMany(ctx, delCategoryIds, data.ProductID)

		if errCT2 != nil {
			return errCT2
		}
	}
	return nil
}

func (p *ProductService) PatchProduct(ctx context.Context, data *structs.PatchProductRequest) error {

	if data == nil {
		return fmt.Errorf("boş veri gönderildi ")
	}

	product := models.Product{
		ProductID:   data.ProductID,
		ProductName: data.ProductName,
		Stock:       data.Stock,
		Price:       data.Price,
	}

	err := p.ProductStore.Patch(ctx, &product)

	if err != nil {
		return fmt.Errorf("db sorgusunda hata %w", err)
	}

	return nil
}

// Var olan Producta Resim Ekleme
// isChanged Parametresi yalnızca kapak fotosu değiştirilmek istendiğinde
func (p *ProductService) AddProductImage(ctx context.Context, images []*multipart.FileHeader, product_id string, isChanged bool) error {

	// db kaydı oluşturma
	// dosyayı okuma
	new_thumb_id := uuid.NewString()

	imgerr1 := p.Image_service.CreateImages(ctx, images, product_id, isChanged, new_thumb_id)

	if imgerr1 != nil {
		return fmt.Errorf("image kayıt sırasında hatayla karşılaşıldı %w", imgerr1)
	}

	// is thumbnail image changed
	if isChanged {
		log.Printf("thumbnail_id budur:%s", new_thumb_id)
		err := p.ProductStore.ChangeThumbnailImage(ctx, new_thumb_id, product_id)

		if err != nil {
			return err
		}

	}
	return nil
}

func (p *ProductService) GetManyProducts(ctx context.Context, ids []string) ([]*structs.ProductForCart, error) {

	products, err := p.ProductStore.GetMany(ctx, ids)

	if err != nil {
		log.Println("db sorgusunda hata %w", err)
		return nil, fmt.Errorf("db sorgusunda hata %w", err)
	}

	for i := range products {
		products[i].ThumbnailImage = Utils.MatchWithURLSingle(p.storagePublicUrl, p.BucketName, products[i].ThumbnailImage)
	}

	return products, nil
}

func (p *ProductService) GetProductsByCategoryId(ctx context.Context, category_ids []string, limit string, offset string) (*structs.CategoryFilteredProducts, error) {

	count, errC := p.ProductStore.GetCountProductsByCategory(ctx, category_ids)

	if errC != nil && count == -1 {
		return nil, errC
	}

	if count != 0 {
		data, err := p.ProductStore.GetProductsByCategoryId(ctx, category_ids, limit, offset)

		if err != nil {
			return nil, err
		}

		for i := range data {
			data[i].ThumbnailImage = Utils.MatchWithURLSingle(p.storagePublicUrl, p.BucketName, data[i].ThumbnailImage)
		}

		response := structs.CategoryFilteredProducts{
			Count: count,
			Data:  data,
		}

		return &response, nil
	}

	// 0 sonuç bulunduğu zaman özel bir hata yapılabilir
	return &structs.CategoryFilteredProducts{
		Count: 0,
		Data:  []structs.Getproducts{},
	}, nil

}

func (p *ProductService) GetFilteredAndSortedProducts(
	ctx context.Context, category_ids []string, limit, offset, sortBy string) (*structs.CategoryFilteredProducts, error) {

	var sortOrder string

	switch sortBy {
	case "priceASC":
		sortBy = constants.Price
		sortOrder = constants.Asc

	case "priceDESC":
		sortBy = constants.Price
		sortOrder = constants.Desc

	case "nameASC":
		sortBy = constants.Name
		sortOrder = constants.Asc

	case "nameDESC":
		sortBy = constants.Name
		sortOrder = constants.Desc

	case "dateDESC":
		sortBy = constants.Date
		sortOrder = constants.Desc
	default:
		sortBy = constants.Date
		sortOrder = constants.Desc
	}

	if len(category_ids) != 0 {
		category_ids = strings.Split(category_ids[0], ",")

	}

	data, err := p.ProductStore.FilteredAndSortedProducts(ctx, category_ids, limit, offset, sortBy, sortOrder)

	if err != nil {
		return nil, err
	}

	for i := range data {
		data[i].ThumbnailImage = Utils.MatchWithURLSingle(p.storagePublicUrl, p.BucketName, data[i].ThumbnailImage)
	}

	count, errC := p.ProductStore.GetCountProductsByCategory(ctx, category_ids)

	if errC != nil {
		return nil, errC
	}

	response := structs.CategoryFilteredProducts{
		Count: count,
		Data:  data,
	}

	return &response, nil

}

// bunun burada ne işi var
func (p *ProductService) GetCategories(ctx context.Context) ([]structs.GetCategory, error) {

	categories, err := p.CategoryStore.GetAll(ctx)

	if err != nil {
		return nil, err
	}

	return categories, nil

}

// arama fonksiyonları

func (p *ProductService) SearchProduct(ctx context.Context, limit string, searchQuery string) ([]structs.Getproducts, error) {

	// count
	count := p.ProductStore.SearchProductCount(ctx, searchQuery)
	var products []structs.Getproducts
	var err error
	if count == 0 {
		// call alternate search

		products, err = p.ProductStore.AlternateSearch(ctx, searchQuery, limit)

		if err != nil {
			log.Printf("cpunt 0 blogu err:%v", err)

			return nil, err
		}

	} else {
		products, err = p.ProductStore.SearchProduct(ctx, searchQuery, limit)
		if err != nil {
			log.Printf("else blogu err:%v", err)
			return nil, err
		}

	}

	for i := 0; i < len(products); i++ {
		products[i].ThumbnailImage = Utils.MatchWithURLSingle(p.storagePublicUrl, p.BucketName, products[i].ThumbnailImage)
	}

	return products, nil

}

func (p *ProductService) SearchProductAdmin(ctx context.Context, limit string, searchQuery string) ([]structs.GetAllProductResponseAdmin, error) {

	// count
	count := p.ProductStore.SearchProductCountAdmin(ctx, searchQuery)
	var products []structs.GetproductsAdmin
	var err error
	if count == 0 {
		// call alternate search

		products, err = p.ProductStore.AlternateSearchAdmin(ctx, searchQuery, limit)

		if err != nil {
			log.Printf("cpunt 0 blogu err:%v", err)

			return nil, err
		}

	} else {
		products, err = p.ProductStore.SearchProductAdmin(ctx, searchQuery, limit)
		if err != nil {
			log.Printf("else blogu err:%v", err)
			return nil, err
		}

	}

	var productIds []string

	for i := range products {
		productIds = append(productIds, products[i].ProductID)

	}

	categories, errCT := p.ProductCategoryStore.GetAllProductCategories(ctx, productIds)

	if errCT != nil {
		return nil, errCT
	}

	var productData []structs.GetAllProductResponseAdmin

	for i := range products {
		products[i].ThumbnailImage = Utils.MatchWithURLSingle(p.storagePublicUrl, p.BucketName, products[i].ThumbnailImage)
		for j := range categories {
			if products[i].ProductID == categories[j].ProductId {
				temp := structs.GetAllProductResponseAdmin{
					Product: structs.GetproductsAdmin{
						ProductID:           products[i].ProductID,
						ProductName:         products[i].ProductName,
						Stock:               products[i].Stock,
						ThumbnailImage:      products[i].ThumbnailImage,
						Price:               products[i].Price,
						CreatedAt:           products[i].CreatedAt,
						ProductCode:         products[i].ProductCode,
						ProductPublishState: products[i].ProductPublishState,
					},
					Categories: categories[j].Categories,
				}
				productData = append(productData, temp)
			}
		}
	}

	return productData, nil

}

func (p *ProductService) UpdateProductPublishState(ctx context.Context, productID string, newState string) error {
	if newState != constants.PublishStatePublished && newState != constants.PublishStateUnpublished {
		return fmt.Errorf("geçersiz yayın durumu: %s", newState)
	}

	err := p.ProductStore.ChangePublishState(ctx, productID, newState)
	if err != nil {
		return err
	}
	return nil
}
