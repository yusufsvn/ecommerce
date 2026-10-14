package services

import (
	"context"
	constants "ecommerceapi/Constants"
	e "ecommerceapi/CustomErrors"
	db "ecommerceapi/Database/Repository"
	st "ecommerceapi/Structs"
	structs "ecommerceapi/Structs"
	"ecommerceapi/Utils"
	"fmt"
	"log"

	"github.com/google/uuid"
	"golang.org/x/sync/errgroup"
)

type PaymentFunctions interface {
	CreatePaymentNonTDS(req structs.CreateOrderDTO) bool
	CreatePaymentWith3DS(data structs.CreateOrderDTO) (string, error)
}

type OrderService struct {
	cartService       CartService
	cartProducts      db.CartProductFunctions
	users             db.UserFunctions
	paymentService    PaymentFunctions
	orderFunc         db.OrderFunctions
	productorders     db.ProductOrderFunctions
	products          db.ProductFunctions
	productcategories db.ProductProductCategoryFunctions
	StorageAPIUrl     string
	storagePublicUrl  string
	BucketName        string
	addressService    AddressService
}

func NewOrderService(
	_cartService *CartService,
	_cartProductsFunc db.CartProductFunctions,
	_userFunc db.UserFunctions,
	_paymentService PaymentFunctions,
	_orderFunc db.OrderFunctions,
	_productorders db.ProductOrderFunctions,
	_productsFunc db.ProductFunctions,
	_productCategories db.ProductProductCategoryFunctions,
	_storageAPIUrl string,
	_storagePublicUrl string,
	_bucketName string,
	_addressService AddressService) *OrderService {
	return &OrderService{
		cartService:       *_cartService,
		cartProducts:      _cartProductsFunc,
		users:             _userFunc,
		paymentService:    _paymentService,
		orderFunc:         _orderFunc,
		productorders:     _productorders,
		products:          _productsFunc,
		productcategories: _productCategories,
		StorageAPIUrl:     _storageAPIUrl,
		storagePublicUrl:  _storagePublicUrl,
		BucketName:        _bucketName,
		addressService:    _addressService,
	}
}

func (ors *OrderService) CreateOrder(ctx context.Context, req st.CreatePaymentDto) error {

	// byte değerini kullanarak karşılaştırır
	if len(req.PhoneNumber) == 0 || req.PhoneNumber[0] != byte(53) {
		log.Println("hatalı telefon numarası")
		return e.ErrInvalidPhoneNumber
	}

	// sepet bilgisi
	ids, err_id := ors.cartService.GetUserCartId(ctx)

	if err_id != nil {
		log.Printf("id hatası-OrderService--CreateOrder--%v", err_id)
		return fmt.Errorf("id error")
	}

	// gen new order id
	orderId, err := uuid.NewUUID()

	if err != nil {
		log.Println("id error--OrderService--CreateOrder")
	}

	// errgroup ve iptal edilebilir context
	g, gCtx := errgroup.WithContext(ctx)

	var basketItemswithcategory []st.ProductsFromCartDTO
	var totalPrice int64
	var userwadd *structs.UserWithoutAddress

	// Goroutine 1 :
	g.Go(func() error {

		// gCtx kullanılarak olası bir iptal durumunda veritabanı sorgusunun kesilmesi sağlanır
		// sepetteki tüm ürünleri getir. (name,price,category,itemType)
		basketItems, err1 := ors.cartProducts.GetAllProductsFromCart(gCtx, ids.Cart_id)

		if err1 != nil {
			return fmt.Errorf("sepet verileri getirelemedi")
		}

		var productIds []string

		log.Println("Debug Basket Item type:")
		for i := range basketItems {
			log.Printf("i : %d basket item type %s:", i, basketItems[i].ItemType)
			productIds = append(productIds, basketItems[i].ProductID)
		}

		// item type da hata vermişti ama şu an kategori düzenlemelerini yapıyoruz ödeme için
		categories, errct := ors.productcategories.GetAllProductCategories(gCtx, productIds)

		if errct != nil {
			return errct
		}

		// dto eşlemesi
		// sliced just category name
		for i, el := range basketItems {
			temp := st.ProductsFromCartDTO{
				ProductID:   el.ProductID,
				ProductName: el.ProductName,
				Price:       el.Price,
				ItemType:    el.ItemType,
				Quantity:    el.Quantity,
				Category:    categories[i].Categories[0].CategoryName,
			}

			basketItemswithcategory = append(basketItemswithcategory, temp)
		}

		// total price hesapla
		totalPrice = Utils.CalcTotalPrice(basketItemswithcategory)

		if totalPrice <= 0 {
			log.Println("total price 0 veya eksi değer olamaz")
			return fmt.Errorf("total price 0")
		}
		return nil
	})

	// Goroutine 2:

	g.Go(func() error {
		var err1 error
		// user bilgilerini getir (name,surname,mail,phone,Address: object)
		userwadd, err1 = ors.users.GetUser(gCtx, ids.User_id)

		if err1 != nil {
			log.Printf("user getirilemedi")
			return fmt.Errorf("user cannot fetched")
		}

		return nil
	})

	// Wait(), başlatılan tüm goroutine'lerin bitmesini bekler.
	// Eğer herhangi bir g.Go() bloğu nil dışında bir error döndürürse,
	// Wait() ilk karşılaştığı hatayı döndürür ve gCtx iptal (cancel) edilir.
	if err := g.Wait(); err != nil {
		return err
	}

	OrderDTO := st.CreateOrderDTO{
		Order_id:   orderId.String(),
		Cart_id:    ids.Cart_id,
		TotalPrice: fmt.Sprint(totalPrice),
		UserData: st.UserOrderDTO{
			User_id: ids.User_id,
			Name:    userwadd.Name,
			Surname: userwadd.Surname,
			Mail:    userwadd.Email,
			Phone:   userwadd.Phone,
			Address: structs.Address{
				AddressStr: req.BillingAddress.AddressStr,
				City:       req.BillingAddress.City,
				Country:    req.BillingAddress.Country,
				ZipCode:    req.BillingAddress.ZipCode,
			},
		},
		CartItems: basketItemswithcategory,
		ReqData:   req,
	}

	// ödeme servisi
	if !req.SecurePay {
		isPaid := ors.paymentService.CreatePaymentNonTDS(OrderDTO)

		if !isPaid {
			return fmt.Errorf("ödeme başarısız")
		}
	} else {
		// decodedHtmlContent, err := ors.paymentService.CreatePaymentWith3DS(OrderDTO)
		// if err != nil {
		// 	return fmt.Errorf("ödeme başarısız")
		// }
		// return decodedHtmlContent
		return fmt.Errorf("not supported 3ds payment")
	}

	// db save - orderRepo/cartproductsRepo
	ordercode, errC := Utils.CreateNanoId(8)

	if errC != nil || ordercode == "" {
		return errC
	}
	orderDto2 := st.OrderDTO{
		Order_id:        orderId.String(),
		User_id:         ids.User_id,
		Status:          constants.OrderStatePending,
		TotalAmount:     fmt.Sprint(totalPrice),
		ShippingAddress: req.ShippingAddress.AddressStr,
		BillingAddress:  req.BillingAddress.AddressStr,
		OrderCode:       ordercode,
	}

	err2 := ors.orderFunc.Create(ctx, &orderDto2)

	if err2 != nil {
		log.Printf("order save error--orderservice--CreateOrder-- %v", err2)
		return fmt.Errorf("order save error")
	}

	err5 := ors.productorders.AddMany(ctx, basketItemswithcategory, orderDto2.Order_id)

	if err5 != nil {
		log.Printf("db error createorder bulk insert: %v", err5)
		return fmt.Errorf("beklenmeyen hatayla karşılaşıldı")
	}

	// second process : sepet temizleme ve varsa adres ekleme işlemi
	g2, gCtx2 := errgroup.WithContext(ctx)

	g2.Go(func() error {
		// delete cart products
		err3 := ors.cartProducts.RemoveAll(gCtx2, ids.Cart_id)

		if err3 != nil {
			log.Printf("sepet ürünleri silinemedi")
			return fmt.Errorf("sepet ürünleri silinemedi")
		}
		log.Println(orderDto2)
		return nil
	})

	if req.SaveAddress {
		g2.Go(func() error {
			errAdd := ors.addressService.CreateAddress(gCtx2, req.ShippingAddress)
			if errAdd != nil {
				return errAdd
			}
			return nil
		})
	}

	if err := g2.Wait(); err != nil {
		return err
	}

	return nil
}

// get from ordertable - return obj: orderId,totalPrice,Status,time
func (ors *OrderService) GetOrdersByUserId(ctx context.Context) ([]st.GetOrder, error) {

	user_id, err_id := Utils.GetUserIdFromClaims(ctx)

	if err_id != nil || user_id == "" {
		log.Printf("user_id zero length")
		return nil, fmt.Errorf("beklenmeyen hatayla karşılaşıldı")
	}

	data, err := ors.orderFunc.GetOrdersByUserId(ctx, user_id)

	if err != nil {
		log.Printf("db error:OrderService--GetOrdersByUserId--%v", err)
		return nil, fmt.Errorf("db error")
	}

	return data, nil

}

func (ors *OrderService) GetOrderDetails(ctx context.Context, order_code string) (*st.OrderResponse, error) {

	// orders tablosundan id ve adresleri getir : order code a göre
	data1, err := ors.orderFunc.GetOrderByOrderCode(ctx, order_code)

	if err != nil {
		log.Printf("GetORderDetails: %v ", err)
		return nil, err
	}

	log.Printf("id : %s", data1.OrderId)

	// prdocut detaylarını ve sipariş edilme bilgilerini getir
	prodOrderResult, err1 := ors.productorders.GetByOrderID(ctx, data1.OrderId)

	if err1 != nil {
		log.Printf("productOrder da hata-- GetOrderDetails-- %v", err1)
		return nil, fmt.Errorf("beklenmeyen hatayla karşılaşıldı")
	}

	for i := range prodOrderResult {
		prodOrderResult[i].ThumbnailImage =
			Utils.MatchWithURLSingle(ors.storagePublicUrl, ors.BucketName, prodOrderResult[i].ThumbnailImage)
	}
	// prodOrderResult ve data birleştirilip client a gönderilecek

	response := st.OrderResponse{
		Products:        prodOrderResult,
		ShippingAddress: data1.ShippingAddress,
		BillingAddress:  data1.BillingAddress,
	}

	return &response, nil

}

// get all for admin
func (ors *OrderService) GetAllOrders(ctx context.Context, limit string, offset string) (*st.OrderResult, error) {

	totaldata, err1 := ors.orderFunc.GetOrderCount(ctx)

	if err1 != nil {
		return nil, fmt.Errorf("Beklenmeyen Hatayla Karşılaşıldı")

	}

	orders, err := ors.orderFunc.GetAll(ctx, limit, offset)

	if err != nil {
		log.Printf("db error: getallorders--%v", err)
		return nil, err
	}

	result := st.OrderResult{
		Count: totaldata,
		Data:  orders,
	}

	return &result, nil

}

func (ors *OrderService) UpdateOrderState(ctx context.Context, newState string, order_id string) error {

	result := Utils.IsValidState(newState)
	if !result {
		log.Printf("state eşleşmedi")
		return fmt.Errorf("Hatalı State Girdiniz")
	}

	err1 := ors.orderFunc.UpdateOrderState(ctx, newState, order_id)

	if err1 != nil {
		return fmt.Errorf("Beklenmeyen Hatayla Karşılaşıldı")
	}

	return nil
}

func (ors *OrderService) UpdateMultiState(ctx context.Context, request *st.UpdateMultiOrderReq) error {

	err := ors.orderFunc.UpdateMultiState(ctx, request.OrderIds, request.NewState)

	if err != nil {
		return err
	}

	return nil
}

func (ors *OrderService) GetOrdersByStatus(ctx context.Context, limit string, offset string, status string) (*st.OrderResult, error) {

	totaldata, err1 := ors.orderFunc.GetOrderCountByStatus(ctx, status)

	if err1 != nil || totaldata == -1 {
		return nil, err1

	}

	orders, err := ors.orderFunc.GetOrdersByStatus(ctx, limit, offset, status)

	if err != nil {
		log.Printf("db error: getallorders--%v", err)
		return nil, err
	}

	result := st.OrderResult{
		Count: totaldata,
		Data:  orders,
	}

	return &result, nil

}

// func (ors *OrderService) CreateOrderORJ(ctx context.Context, req st.CreatePaymentDto) error {

// 	// sepet bilgisi
// 	ids, err_id := ors.cartService.GetUserCartId(ctx)

// 	if err_id != nil {
// 		log.Printf("id hatası-OrderService--CreateOrder--%v", err_id)
// 		return fmt.Errorf("id error")
// 	}

// 	// gen new order id
// 	orderId, err := uuid.NewUUID()

// 	if err != nil {
// 		log.Println("id error--OrderService--CreateOrder")
// 	}

// 	// sepetteki tüm ürünleri getir. (name,price,category,itemType)
// 	basketItems, err1 := ors.cartProducts.GetAllProductsFromCart(ctx, ids.Cart_id)

// 	if err1 != nil {
// 		return fmt.Errorf("sepet verileri getirelemedi")
// 	}

// 	var productIds []string

// 	log.Println("Debug Basket Item type:")
// 	for i := range basketItems {
// 		log.Printf("i : %d basket item type %s:", i, basketItems[i].ItemType)
// 		productIds = append(productIds, basketItems[i].ProductID)
// 	}

// 	// item type da hata vermişti ama şu an kategori düzenlemelerini yapıyoruz ödeme için
// 	categories, errct := ors.productcategories.GetAllProductCategories(ctx, productIds)

// 	if errct != nil {
// 		return errct
// 	}

// 	// dto eşlemesi
// 	var basketItemswithcategory []st.ProductsFromCartDTO
// 	// sliced just category name
// 	for i, el := range basketItems {
// 		temp := st.ProductsFromCartDTO{
// 			ProductID:   el.ProductID,
// 			ProductName: el.ProductName,
// 			Price:       el.Price,
// 			ItemType:    el.ItemType,
// 			Quantity:    el.Quantity,
// 			Category:    categories[i].Categories[0].CategoryName,
// 		}

// 		basketItemswithcategory = append(basketItemswithcategory, temp)
// 	}

// 	// total price hesapla
// 	totalPrice := Utils.CalcTotalPrice(basketItemswithcategory)

// 	if totalPrice <= 0 {
// 		log.Println("total price 0 veya eksi değer olamaz")
// 		return fmt.Errorf("total price 0")
// 	}
// 	// user bilgilerini getir (name,surname,mail,phone,Address: object)
// 	userwadd, err1 := ors.users.GetUser(ctx, ids.User_id)

// 	if err1 != nil {
// 		log.Printf("user getirilemedi")
// 		return fmt.Errorf("user cannot fetched")
// 	}
// 	OrderDTO := st.CreateOrderDTO{
// 		Order_id:   orderId.String(),
// 		Cart_id:    ids.Cart_id,
// 		TotalPrice: fmt.Sprint(totalPrice),
// 		UserData: st.UserOrderDTO{
// 			User_id: ids.User_id,
// 			Name:    userwadd.Name,
// 			Surname: userwadd.Surname,
// 			Mail:    userwadd.Email,
// 			Phone:   req.PhoneNumber,
// 			Address: structs.Address{
// 				AddressStr: req.ShippingAddress.AddressStr,
// 				City:       userwadd.City,
// 				Country:    userwadd.Country,
// 				ZipCode:    userwadd.ZipCode,
// 			},
// 		},
// 		CartItems: basketItemswithcategory,
// 		ReqData:   req,
// 	}

// 	// ödeme servisi
// 	if !req.SecurePay {
// 		isPaid := ors.paymentService.CreatePaymentNonTDS(OrderDTO)

// 		if !isPaid {
// 			return fmt.Errorf("ödeme başarısız")
// 		}
// 	} else {
// 		// decodedHtmlContent, err := ors.paymentService.CreatePaymentWith3DS(OrderDTO)
// 		// if err != nil {
// 		// 	return fmt.Errorf("ödeme başarısız")
// 		// }
// 		// return decodedHtmlContent
// 		return fmt.Errorf("not supported 3ds payment")
// 	}

// 	// db save - orderRepo/cartproductsRepo
// 	ordercode, errC := Utils.CreateNanoId(8)

// 	if errC != nil || ordercode == "" {
// 		return errC
// 	}
// 	orderDto2 := st.OrderDTO{
// 		Order_id:        orderId.String(),
// 		User_id:         ids.User_id,
// 		Status:          "paid",
// 		TotalAmount:     fmt.Sprint(totalPrice),
// 		ShippingAddress: req.ShippingAddress.AddressStr,
// 		BillingAddress:  req.BillingAddress.AddressStr,
// 		OrderCode:       ordercode,
// 	}

// 	err2 := ors.orderFunc.Create(ctx, &orderDto2)

// 	if err2 != nil {
// 		log.Printf("order save error--orderservice--CreateOrder-- %v", err2)
// 		return fmt.Errorf("order save error")
// 	}

// 	err5 := ors.productorders.AddMany(ctx, basketItemswithcategory, orderDto2.Order_id)

// 	if err5 != nil {
// 		log.Printf("db error createorder bulk insert: %v", err5)
// 		return fmt.Errorf("beklenmeyen hatayla karşılaşıldı")
// 	}

// 	// delete cart products
// 	err3 := ors.cartProducts.RemoveAll(ctx, ids.Cart_id)

// 	if err3 != nil {
// 		log.Printf("sepet ürünleri silinemedi")
// 		return fmt.Errorf("sepet ürünleri silinemedi")
// 	}
// 	log.Println(orderDto2)

// 	if req.SaveAddress {
// 		errAdd := ors.addressService.CreateAddress(ctx, req.ShippingAddress)
// 		if errAdd != nil {
// 			return errAdd
// 		}
// 	}

// 	return nil
// }

func (ors *OrderService) GetOrderCountForWaiting(ctx context.Context) (int, error) {

	count, err := ors.orderFunc.GetOrderCountForWaiting(ctx)

	if err != nil {
		return -1, err
	}

	return count, nil
}

func (ors *OrderService) SearchOrder(ctx context.Context, limit string, search_query string) (*st.OrderResult, error) {

	count, err := ors.orderFunc.SearchOrderCount(ctx, search_query)

	if err != nil {
		log.Printf("orderRepo--SearchOrderCount--error:%v", err)
		return nil, fmt.Errorf("sipariş arama sayısı hatası: %w", err)
	}

	if count == 0 {
		return nil, e.ZeroListResult
	}

	orders, err := ors.orderFunc.SearchOrder(ctx, limit, search_query)

	if err != nil {
		log.Printf("orderRepo--SearchOrder--error:%v", err)
		return nil, fmt.Errorf("sipariş arama hatası: %w", err)
	}

	result := st.OrderResult{
		Count: count,
		Data:  orders,
	}

	return &result, nil
}

func (ors *OrderService) RemoveByProductId(ctx context.Context, product_id *uuid.UUID) error {

	log.Printf("orderService--RemoveByProductId--product_id:%s", product_id)

	err := ors.productorders.Delete(ctx, *product_id)

	if err != nil {
		log.Printf("orderService--RemoveByProductId--error:%v", err)
		return fmt.Errorf("ürün silme hatası: %w", err)
	}

	return nil
}
