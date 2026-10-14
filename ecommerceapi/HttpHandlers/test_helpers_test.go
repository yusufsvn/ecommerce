package httphandlers

import (
	"context"
	"io"
	models "ecommerceapi/Database/Models"
	Entities "ecommerceapi/Database/Models"
	st "ecommerceapi/Structs"
	structs "ecommerceapi/Structs"
	"ecommerceapi/Types"
	"net/http"

	"github.com/golang-jwt/jwt/v5"
	"github.com/google/uuid"
)

// Helper: Context'e JWT claim'leri ekler
func withUserClaims(r *http.Request, userID, role string) *http.Request {
	claims := jwt.MapClaims{
		"user_id": userID,
		"role":    role,
	}
	ctx := context.WithValue(r.Context(), Types.UserContextKey, claims)
	return r.WithContext(ctx)
}

// ----------------------------------------------------
// Mock Address Repository
// ----------------------------------------------------
type mockAddressRepo struct {
	createFn      func(ctx context.Context, address *st.CreateAddressDTO) error
	deleteFn      func(ctx context.Context, addressId string) error
	updateFn      func(ctx context.Context, address st.Address, user_id string) error
	getByUserIdFn func(ctx context.Context, user_id string) (*[]st.GetAddressDTO, error)
}

func (m *mockAddressRepo) Create(ctx context.Context, address *st.CreateAddressDTO) error {
	if m.createFn != nil {
		return m.createFn(ctx, address)
	}
	return nil
}
func (m *mockAddressRepo) Delete(ctx context.Context, addressId string) error {
	if m.deleteFn != nil {
		return m.deleteFn(ctx, addressId)
	}
	return nil
}
func (m *mockAddressRepo) Update(ctx context.Context, address st.Address, user_id string) error {
	if m.updateFn != nil {
		return m.updateFn(ctx, address, user_id)
	}
	return nil
}
func (m *mockAddressRepo) GetByUserId(ctx context.Context, user_id string) (*[]st.GetAddressDTO, error) {
	if m.getByUserIdFn != nil {
		return m.getByUserIdFn(ctx, user_id)
	}
	return &[]st.GetAddressDTO{}, nil
}

// ----------------------------------------------------
// Mock Category Repository (ProductCategoryFunctions)
// ----------------------------------------------------
type mockCategoryRepo struct {
	getAllFn                  func(ctx context.Context) ([]structs.GetCategory, error)
	createFn                  func(ctx context.Context, category_name string) (structs.GetCategory, error)
	deleteFn                  func(ctx context.Context, categoryId string) error
	getByCategoryNameFn       func(ctx context.Context, categoryName string) (uuid.UUID, error)
	getCategoriesByProductIdFn func(ctx context.Context, product_id string) ([]structs.GetCategory, error)
	getByCategoryIdsFn        func(ctx context.Context, categoryNames []string) ([]string, error)
	updateFn                  func(ctx context.Context, category_id string, category_name string) error
}

func (m *mockCategoryRepo) GetAll(ctx context.Context) ([]structs.GetCategory, error) {
	if m.getAllFn != nil {
		return m.getAllFn(ctx)
	}
	return nil, nil
}
func (m *mockCategoryRepo) Create(ctx context.Context, category_name string) (structs.GetCategory, error) {
	if m.createFn != nil {
		return m.createFn(ctx, category_name)
	}
	return structs.GetCategory{}, nil
}
func (m *mockCategoryRepo) Delete(ctx context.Context, categoryId string) error {
	if m.deleteFn != nil {
		return m.deleteFn(ctx, categoryId)
	}
	return nil
}
func (m *mockCategoryRepo) GetByCategoryName(ctx context.Context, categoryName string) (uuid.UUID, error) {
	if m.getByCategoryNameFn != nil {
		return m.getByCategoryNameFn(ctx, categoryName)
	}
	return uuid.Nil, nil
}
func (m *mockCategoryRepo) GetCategoriesByProductId(ctx context.Context, product_id string) ([]structs.GetCategory, error) {
	if m.getCategoriesByProductIdFn != nil {
		return m.getCategoriesByProductIdFn(ctx, product_id)
	}
	return nil, nil
}
func (m *mockCategoryRepo) GetByCategoryIds(ctx context.Context, categoryNames []string) ([]string, error) {
	if m.getByCategoryIdsFn != nil {
		return m.getByCategoryIdsFn(ctx, categoryNames)
	}
	return nil, nil
}
func (m *mockCategoryRepo) Update(ctx context.Context, category_id string, category_name string) error {
	if m.updateFn != nil {
		return m.updateFn(ctx, category_id, category_name)
	}
	return nil
}

// ----------------------------------------------------
// Mock Cart Repository
// ----------------------------------------------------
type mockCartRepo struct {
	getCartByUserIDFn         func(ctx context.Context, user_id string) (string, error)
	createFn                  func(ctx context.Context, user_id string) (string, error)
	deleteFn                  func(ctx context.Context, cartId string) error
	getCartProductsByUserIDFn func(ctx context.Context, user_id string) ([]structs.CartProductDTO, error)
}

func (m *mockCartRepo) GetCartByUserID(ctx context.Context, user_id string) (string, error) {
	if m.getCartByUserIDFn != nil {
		return m.getCartByUserIDFn(ctx, user_id)
	}
	return "", nil
}
func (m *mockCartRepo) Create(ctx context.Context, user_id string) (string, error) {
	if m.createFn != nil {
		return m.createFn(ctx, user_id)
	}
	return "", nil
}
func (m *mockCartRepo) Delete(ctx context.Context, cartId string) error {
	if m.deleteFn != nil {
		return m.deleteFn(ctx, cartId)
	}
	return nil
}
func (m *mockCartRepo) GetCartProductsByUserID(ctx context.Context, user_id string) ([]structs.CartProductDTO, error) {
	if m.getCartProductsByUserIDFn != nil {
		return m.getCartProductsByUserIDFn(ctx, user_id)
	}
	return nil, nil
}

// ----------------------------------------------------
// Mock CartProduct Repository
// ----------------------------------------------------
type mockCartProductRepo struct {
	getByCartIDFn            func(ctx context.Context, cartID uuid.UUID) ([]Entities.CartProduct, error)
	addFn                    func(ctx context.Context, cartProduct *st.CartProductDTO) error
	removeFn                 func(ctx context.Context, cart_id string, product_id string) error
	updateQuantityFn         func(ctx context.Context, data *st.CartProductDTO) error
	getAllProductsFromCartFn func(ctx context.Context, cart_id string) ([]st.ProductsFromCartDb, error)
	removeAllFn              func(ctx context.Context, cart_id string) error
}

func (m *mockCartProductRepo) GetByCartID(ctx context.Context, cartID uuid.UUID) ([]Entities.CartProduct, error) {
	if m.getByCartIDFn != nil {
		return m.getByCartIDFn(ctx, cartID)
	}
	return nil, nil
}
func (m *mockCartProductRepo) Add(ctx context.Context, cartProduct *st.CartProductDTO) error {
	if m.addFn != nil {
		return m.addFn(ctx, cartProduct)
	}
	return nil
}
func (m *mockCartProductRepo) Remove(ctx context.Context, cart_id string, product_id string) error {
	if m.removeFn != nil {
		return m.removeFn(ctx, cart_id, product_id)
	}
	return nil
}
func (m *mockCartProductRepo) UpdateQuantity(ctx context.Context, data *st.CartProductDTO) error {
	if m.updateQuantityFn != nil {
		return m.updateQuantityFn(ctx, data)
	}
	return nil
}
func (m *mockCartProductRepo) GetAllProductsFromCart(ctx context.Context, cart_id string) ([]st.ProductsFromCartDb, error) {
	if m.getAllProductsFromCartFn != nil {
		return m.getAllProductsFromCartFn(ctx, cart_id)
	}
	return nil, nil
}
func (m *mockCartProductRepo) RemoveAll(ctx context.Context, cart_id string) error {
	if m.removeAllFn != nil {
		return m.removeAllFn(ctx, cart_id)
	}
	return nil
}

// ----------------------------------------------------
// Mock Token Repository
// ----------------------------------------------------
type mockTokenRepo struct {
	addTokenFn       func(ctx context.Context, token st.RefreshTokenSaveDTO) error
	removeTokenFn    func(ctx context.Context, token_id string) error
	getTokenByIdFn   func(ctx context.Context, token_id string) (*st.TokenQueryRes, error)
	updateTokenFn    func(ctx context.Context, token_id string, is_used bool) error
	panicUpdateFn    func(ctx context.Context, family_id string, is_revoked bool) error
}

func (m *mockTokenRepo) AddToken(ctx context.Context, token st.RefreshTokenSaveDTO) error {
	if m.addTokenFn != nil {
		return m.addTokenFn(ctx, token)
	}
	return nil
}
func (m *mockTokenRepo) RemoveToken(ctx context.Context, token_id string) error {
	if m.removeTokenFn != nil {
		return m.removeTokenFn(ctx, token_id)
	}
	return nil
}
func (m *mockTokenRepo) GetTokenById(ctx context.Context, token_id string) (*st.TokenQueryRes, error) {
	if m.getTokenByIdFn != nil {
		return m.getTokenByIdFn(ctx, token_id)
	}
	return nil, nil
}
func (m *mockTokenRepo) UpdateToken(ctx context.Context, token_id string, is_used bool) error {
	if m.updateTokenFn != nil {
		return m.updateTokenFn(ctx, token_id, is_used)
	}
	return nil
}
func (m *mockTokenRepo) PanicUpdate(ctx context.Context, family_id string, is_revoked bool) error {
	if m.panicUpdateFn != nil {
		return m.panicUpdateFn(ctx, family_id, is_revoked)
	}
	return nil
}

// ----------------------------------------------------
// Mock User Repository
// ----------------------------------------------------
type mockUserRepo struct {
	getAllFn               func(limit int, offset int, ctx context.Context) ([]Entities.User, error)
	createFn               func(ctx context.Context, user *Entities.User) (string, error)
	deleteFn               func(ctx context.Context, userId *uuid.UUID) error
	updateFn               func(ctx context.Context, user *Entities.User) error
	getUserByEmailFn       func(ctx context.Context, email string) (*st.UserLoginRes, error)
	getUserByIdFn          func(ctx context.Context, userId string) (string, error)
	getUserFn              func(ctx context.Context, user_id string) (*st.UserWithoutAddress, error)
	getUserCountFn         func(ctx context.Context) (int, error)
	addNewUserWithGoogleFn func(ctx context.Context, user *st.GoogleUser) (string, error)
	getUserByOauthIdFn     func(ctx context.Context, oauthId string) (*st.UserLoginRes, error)
}

func (m *mockUserRepo) GetAll(limit int, offset int, ctx context.Context) ([]Entities.User, error) {
	if m.getAllFn != nil {
		return m.getAllFn(limit, offset, ctx)
	}
	return nil, nil
}
func (m *mockUserRepo) Create(ctx context.Context, user *Entities.User) (string, error) {
	if m.createFn != nil {
		return m.createFn(ctx, user)
	}
	return "", nil
}
func (m *mockUserRepo) Delete(ctx context.Context, userId *uuid.UUID) error {
	if m.deleteFn != nil {
		return m.deleteFn(ctx, userId)
	}
	return nil
}
func (m *mockUserRepo) Update(ctx context.Context, user *Entities.User) error {
	if m.updateFn != nil {
		return m.updateFn(ctx, user)
	}
	return nil
}
func (m *mockUserRepo) GetUserByEmail(ctx context.Context, email string) (*st.UserLoginRes, error) {
	if m.getUserByEmailFn != nil {
		return m.getUserByEmailFn(ctx, email)
	}
	return nil, nil
}
func (m *mockUserRepo) GetUserById(ctx context.Context, userId string) (string, error) {
	if m.getUserByIdFn != nil {
		return m.getUserByIdFn(ctx, userId)
	}
	return "", nil
}
func (m *mockUserRepo) GetUser(ctx context.Context, user_id string) (*st.UserWithoutAddress, error) {
	if m.getUserFn != nil {
		return m.getUserFn(ctx, user_id)
	}
	return nil, nil
}
func (m *mockUserRepo) GetUserCount(ctx context.Context) (int, error) {
	if m.getUserCountFn != nil {
		return m.getUserCountFn(ctx)
	}
	return 0, nil
}
func (m *mockUserRepo) AddNewUserWithGoogle(ctx context.Context, user *st.GoogleUser) (string, error) {
	if m.addNewUserWithGoogleFn != nil {
		return m.addNewUserWithGoogleFn(ctx, user)
	}
	return "", nil
}
func (m *mockUserRepo) GetUserByOauthId(ctx context.Context, oauthId string) (*st.UserLoginRes, error) {
	if m.getUserByOauthIdFn != nil {
		return m.getUserByOauthIdFn(ctx, oauthId)
	}
	return nil, nil
}

// ----------------------------------------------------
// Mock Order Repository
// ----------------------------------------------------
type mockOrderRepo struct {
	getAllFn                 func(ctx context.Context, limit string, offset string) ([]st.OrderForAdmin, error)
	createFn                 func(ctx context.Context, order *st.OrderDTO) error
	deleteFn                 func(ctx context.Context, orderId uuid.UUID) error
	getOrdersByUserIdFn      func(ctx context.Context, user_id string) ([]st.GetOrder, error)
	updateOrderStateFn       func(ctx context.Context, newState string, order_id string) error
	getOrderCountFn          func(ctx context.Context) (int, error)
	updateMultiStateFn       func(ctx context.Context, idlist []string, new_state string) error
	getOrderCountByStatusFn  func(ctx context.Context, status string) (int, error)
	getOrdersByStatusFn      func(ctx context.Context, limit string, offset string, status string) ([]st.OrderForAdmin, error)
	getOrderByOrderCodeFn    func(ctx context.Context, order_code string) (*st.GetOrderAddress, error)
	getOrderCountForWaitingFn func(ctx context.Context) (int, error)
	searchOrderFn            func(ctx context.Context, limit string, search_query string) ([]st.OrderForAdmin, error)
	searchOrderCountFn       func(ctx context.Context, search_query string) (int, error)
}

func (m *mockOrderRepo) GetAll(ctx context.Context, limit string, offset string) ([]st.OrderForAdmin, error) {
	if m.getAllFn != nil {
		return m.getAllFn(ctx, limit, offset)
	}
	return nil, nil
}
func (m *mockOrderRepo) Create(ctx context.Context, order *st.OrderDTO) error {
	if m.createFn != nil {
		return m.createFn(ctx, order)
	}
	return nil
}
func (m *mockOrderRepo) Delete(ctx context.Context, orderId uuid.UUID) error {
	if m.deleteFn != nil {
		return m.deleteFn(ctx, orderId)
	}
	return nil
}
func (m *mockOrderRepo) GetOrdersByUserId(ctx context.Context, user_id string) ([]st.GetOrder, error) {
	if m.getOrdersByUserIdFn != nil {
		return m.getOrdersByUserIdFn(ctx, user_id)
	}
	return nil, nil
}
func (m *mockOrderRepo) UpdateOrderState(ctx context.Context, newState string, order_id string) error {
	if m.updateOrderStateFn != nil {
		return m.updateOrderStateFn(ctx, newState, order_id)
	}
	return nil
}
func (m *mockOrderRepo) GetOrderCount(ctx context.Context) (int, error) {
	if m.getOrderCountFn != nil {
		return m.getOrderCountFn(ctx)
	}
	return 0, nil
}
func (m *mockOrderRepo) UpdateMultiState(ctx context.Context, idlist []string, new_state string) error {
	if m.updateMultiStateFn != nil {
		return m.updateMultiStateFn(ctx, idlist, new_state)
	}
	return nil
}
func (m *mockOrderRepo) GetOrderCountByStatus(ctx context.Context, status string) (int, error) {
	if m.getOrderCountByStatusFn != nil {
		return m.getOrderCountByStatusFn(ctx, status)
	}
	return 0, nil
}
func (m *mockOrderRepo) GetOrdersByStatus(ctx context.Context, limit string, offset string, status string) ([]st.OrderForAdmin, error) {
	if m.getOrdersByStatusFn != nil {
		return m.getOrdersByStatusFn(ctx, limit, offset, status)
	}
	return nil, nil
}
func (m *mockOrderRepo) GetOrderByOrderCode(ctx context.Context, order_code string) (*st.GetOrderAddress, error) {
	if m.getOrderByOrderCodeFn != nil {
		return m.getOrderByOrderCodeFn(ctx, order_code)
	}
	return nil, nil
}
func (m *mockOrderRepo) GetOrderCountForWaiting(ctx context.Context) (int, error) {
	if m.getOrderCountForWaitingFn != nil {
		return m.getOrderCountForWaitingFn(ctx)
	}
	return 0, nil
}
func (m *mockOrderRepo) SearchOrder(ctx context.Context, limit string, search_query string) ([]st.OrderForAdmin, error) {
	if m.searchOrderFn != nil {
		return m.searchOrderFn(ctx, limit, search_query)
	}
	return nil, nil
}
func (m *mockOrderRepo) SearchOrderCount(ctx context.Context, search_query string) (int, error) {
	if m.searchOrderCountFn != nil {
		return m.searchOrderCountFn(ctx, search_query)
	}
	return 0, nil
}

// ----------------------------------------------------
// Mock ProductOrder Repository
// ----------------------------------------------------
type mockProductOrderRepo struct {
	getByOrderIDFn func(ctx context.Context, orderID string) ([]st.OrderFullDetail, error)
	addFn          func(ctx context.Context, item Entities.ProductOrder) error
	addManyFn      func(ctx context.Context, items []st.ProductsFromCartDTO, order_id string) error
}

func (m *mockProductOrderRepo) GetByOrderID(ctx context.Context, orderID string) ([]st.OrderFullDetail, error) {
	if m.getByOrderIDFn != nil {
		return m.getByOrderIDFn(ctx, orderID)
	}
	return nil, nil
}
func (m *mockProductOrderRepo) Add(ctx context.Context, item Entities.ProductOrder) error {
	if m.addFn != nil {
		return m.addFn(ctx, item)
	}
	return nil
}
func (m *mockProductOrderRepo) AddMany(ctx context.Context, items []st.ProductsFromCartDTO, order_id string) error {
	if m.addManyFn != nil {
		return m.addManyFn(ctx, items, order_id)
	}
	return nil
}

// ----------------------------------------------------
// Mock Product Repository
// ----------------------------------------------------
type mockProductRepo struct {
	createFn                     func(ctx context.Context, product *models.Product) error
	deleteFn                     func(ctx context.Context, productId *uuid.UUID) error
	updateFn                     func(ctx context.Context, product *models.Product) error
	patchFn                      func(ctx context.Context, product *models.Product) error
	changeThumbnailImageFn       func(ctx context.Context, new_id string, product_id string) error
	getAllFn                     func(limit int, offset int, ctx context.Context) ([]structs.Getproducts, error)
	getProductByIdFn             func(ctx context.Context, ProductID uuid.UUID) (*structs.ProductWithImage, error)
	getManyFn                    func(ctx context.Context, ids []string) ([]*structs.ProductForCart, error)
	getProductCountFn            func(ctx context.Context) (int, error)
	getProductsByCategoryIdFn    func(ctx context.Context, category_ids []string, limit string, offset string) ([]structs.Getproducts, error)
	getCountProductsByCategoryFn func(ctx context.Context, category_ids []string) (int, error)
	searchProductCountFn         func(ctx context.Context, searchQuery string) int
	searchProductFn              func(ctx context.Context, searchQuery string, limit string) ([]st.Getproducts, error)
	alternateSearchFn            func(ctx context.Context, searchQuery string, limit string) ([]st.Getproducts, error)
	filteredAndSortedProductsFn  func(ctx context.Context, category_ids []string, limit, offset, sortBy, sortOrder string) ([]structs.Getproducts, error)
	searchProductCountAdminFn    func(ctx context.Context, searchQuery string) int
	searchProductAdminFn         func(ctx context.Context, searchQuery string, limit string) ([]st.GetproductsAdmin, error)
	alternateSearchAdminFn       func(ctx context.Context, searchQuery string, limit string) ([]st.GetproductsAdmin, error)
	getProductCountAdminFn       func(ctx context.Context) (int, error)
	getAllAdminFn                func(limit int, offset int, ctx context.Context) ([]structs.GetproductsAdmin, error)
	changePublishStateFn         func(ctx context.Context, product_id string, newState string) error
}

func (m *mockProductRepo) Create(ctx context.Context, product *models.Product) error {
	if m.createFn != nil {
		return m.createFn(ctx, product)
	}
	return nil
}
func (m *mockProductRepo) Delete(ctx context.Context, productId *uuid.UUID) error {
	if m.deleteFn != nil {
		return m.deleteFn(ctx, productId)
	}
	return nil
}
func (m *mockProductRepo) Update(ctx context.Context, product *models.Product) error {
	if m.updateFn != nil {
		return m.updateFn(ctx, product)
	}
	return nil
}
func (m *mockProductRepo) Patch(ctx context.Context, product *models.Product) error {
	if m.patchFn != nil {
		return m.patchFn(ctx, product)
	}
	return nil
}
func (m *mockProductRepo) ChangeThumbnailImage(ctx context.Context, new_id string, product_id string) error {
	if m.changeThumbnailImageFn != nil {
		return m.changeThumbnailImageFn(ctx, new_id, product_id)
	}
	return nil
}
func (m *mockProductRepo) GetAll(limit int, offset int, ctx context.Context) ([]structs.Getproducts, error) {
	if m.getAllFn != nil {
		return m.getAllFn(limit, offset, ctx)
	}
	return nil, nil
}
func (m *mockProductRepo) GetProductById(ctx context.Context, ProductID uuid.UUID) (*structs.ProductWithImage, error) {
	if m.getProductByIdFn != nil {
		return m.getProductByIdFn(ctx, ProductID)
	}
	return nil, nil
}
func (m *mockProductRepo) GetMany(ctx context.Context, ids []string) ([]*structs.ProductForCart, error) {
	if m.getManyFn != nil {
		return m.getManyFn(ctx, ids)
	}
	return nil, nil
}
func (m *mockProductRepo) GetProductCount(ctx context.Context) (int, error) {
	if m.getProductCountFn != nil {
		return m.getProductCountFn(ctx)
	}
	return 0, nil
}
func (m *mockProductRepo) GetProductsByCategoryId(ctx context.Context, category_ids []string, limit string, offset string) ([]structs.Getproducts, error) {
	if m.getProductsByCategoryIdFn != nil {
		return m.getProductsByCategoryIdFn(ctx, category_ids, limit, offset)
	}
	return nil, nil
}
func (m *mockProductRepo) GetCountProductsByCategory(ctx context.Context, category_ids []string) (int, error) {
	if m.getCountProductsByCategoryFn != nil {
		return m.getCountProductsByCategoryFn(ctx, category_ids)
	}
	return 0, nil
}
func (m *mockProductRepo) SearchProductCount(ctx context.Context, searchQuery string) int {
	if m.searchProductCountFn != nil {
		return m.searchProductCountFn(ctx, searchQuery)
	}
	return 0
}
func (m *mockProductRepo) SearchProduct(ctx context.Context, searchQuery string, limit string) ([]st.Getproducts, error) {
	if m.searchProductFn != nil {
		return m.searchProductFn(ctx, searchQuery, limit)
	}
	return nil, nil
}
func (m *mockProductRepo) AlternateSearch(ctx context.Context, searchQuery string, limit string) ([]st.Getproducts, error) {
	if m.alternateSearchFn != nil {
		return m.alternateSearchFn(ctx, searchQuery, limit)
	}
	return nil, nil
}
func (m *mockProductRepo) FilteredAndSortedProducts(ctx context.Context, category_ids []string, limit, offset, sortBy, sortOrder string) ([]structs.Getproducts, error) {
	if m.filteredAndSortedProductsFn != nil {
		return m.filteredAndSortedProductsFn(ctx, category_ids, limit, offset, sortBy, sortOrder)
	}
	return nil, nil
}
func (m *mockProductRepo) SearchProductCountAdmin(ctx context.Context, searchQuery string) int {
	if m.searchProductCountAdminFn != nil {
		return m.searchProductCountAdminFn(ctx, searchQuery)
	}
	return 0
}
func (m *mockProductRepo) SearchProductAdmin(ctx context.Context, searchQuery string, limit string) ([]st.GetproductsAdmin, error) {
	if m.searchProductAdminFn != nil {
		return m.searchProductAdminFn(ctx, searchQuery, limit)
	}
	return nil, nil
}
func (m *mockProductRepo) AlternateSearchAdmin(ctx context.Context, searchQuery string, limit string) ([]st.GetproductsAdmin, error) {
	if m.alternateSearchAdminFn != nil {
		return m.alternateSearchAdminFn(ctx, searchQuery, limit)
	}
	return nil, nil
}
func (m *mockProductRepo) GetProductCountAdmin(ctx context.Context) (int, error) {
	if m.getProductCountAdminFn != nil {
		return m.getProductCountAdminFn(ctx)
	}
	return 0, nil
}
func (m *mockProductRepo) GetAllAdmin(limit int, offset int, ctx context.Context) ([]structs.GetproductsAdmin, error) {
	if m.getAllAdminFn != nil {
		return m.getAllAdminFn(limit, offset, ctx)
	}
	return nil, nil
}
func (m *mockProductRepo) ChangePublishState(ctx context.Context, product_id string, newState string) error {
	if m.changePublishStateFn != nil {
		return m.changePublishStateFn(ctx, product_id, newState)
	}
	return nil
}

// ----------------------------------------------------
// Mock ProductProductCategory Repository
// ----------------------------------------------------
type mockProductProductCategoryRepo struct {
	getByProductIDFn          func(ctx context.Context, productID uuid.UUID) ([]Entities.ProductProductCategory, error)
	getByCategoryIDFn          func(ctx context.Context, categoryID uuid.UUID) ([]Entities.ProductProductCategory, error)
	getAllProductCategoriesFn func(ctx context.Context, ids []string) ([]*structs.GetAllCategory, error)
	addFn                     func(ctx context.Context, item *Entities.ProductProductCategory) error
	removeFn                  func(ctx context.Context, productID uuid.UUID, categoryID uuid.UUID) error
	getByProductIDOneRowFn    func(ctx context.Context, productID uuid.UUID) (uuid.UUID, error)
	addManyFn                 func(ctx context.Context, items []string, product_id string) error
	removeManyFn              func(ctx context.Context, items []string, product_id string) error
}

func (m *mockProductProductCategoryRepo) GetByProductID(ctx context.Context, productID uuid.UUID) ([]Entities.ProductProductCategory, error) {
	if m.getByProductIDFn != nil {
		return m.getByProductIDFn(ctx, productID)
	}
	return nil, nil
}
func (m *mockProductProductCategoryRepo) GetByCategoryID(ctx context.Context, categoryID uuid.UUID) ([]Entities.ProductProductCategory, error) {
	if m.getByCategoryIDFn != nil {
		return m.getByCategoryIDFn(ctx, categoryID)
	}
	return nil, nil
}
func (m *mockProductProductCategoryRepo) GetAllProductCategories(ctx context.Context, ids []string) ([]*structs.GetAllCategory, error) {
	if m.getAllProductCategoriesFn != nil {
		return m.getAllProductCategoriesFn(ctx, ids)
	}
	return nil, nil
}
func (m *mockProductProductCategoryRepo) Add(ctx context.Context, item *Entities.ProductProductCategory) error {
	if m.addFn != nil {
		return m.addFn(ctx, item)
	}
	return nil
}
func (m *mockProductProductCategoryRepo) Remove(ctx context.Context, productID uuid.UUID, categoryID uuid.UUID) error {
	if m.removeFn != nil {
		return m.removeFn(ctx, productID, categoryID)
	}
	return nil
}
func (m *mockProductProductCategoryRepo) GetByProductIDOneRow(ctx context.Context, productID uuid.UUID) (uuid.UUID, error) {
	if m.getByProductIDOneRowFn != nil {
		return m.getByProductIDOneRowFn(ctx, productID)
	}
	return uuid.Nil, nil
}
func (m *mockProductProductCategoryRepo) AddMany(ctx context.Context, items []string, product_id string) error {
	if m.addManyFn != nil {
		return m.addManyFn(ctx, items, product_id)
	}
	return nil
}
func (m *mockProductProductCategoryRepo) RemoveMany(ctx context.Context, items []string, product_id string) error {
	if m.removeManyFn != nil {
		return m.removeManyFn(ctx, items, product_id)
	}
	return nil
}

// ----------------------------------------------------
// Mock Image Repository & Mock File Manager
// ----------------------------------------------------
type mockImageRepo struct {
	addOneImageFn      func(ctx context.Context, image models.Image) error
	deleteImageFn      func(ctx context.Context, images []string) error
	getProductImagesFn func(ctx context.Context, product_id string) ([]string, error)
}

func (m *mockImageRepo) AddOneImage(ctx context.Context, image models.Image) error {
	if m.addOneImageFn != nil {
		return m.addOneImageFn(ctx, image)
	}
	return nil
}
func (m *mockImageRepo) DeleteImage(ctx context.Context, images []string) error {
	if m.deleteImageFn != nil {
		return m.deleteImageFn(ctx, images)
	}
	return nil
}
func (m *mockImageRepo) GetProductImages(ctx context.Context, product_id string) ([]string, error) {
	if m.getProductImagesFn != nil {
		return m.getProductImagesFn(ctx, product_id)
	}
	return nil, nil
}

type mockFileManager struct {
	uploadFileFn             func(fileName string, fileSize int64, prefix string, ctx context.Context, reader io.Reader) (string, error)
	deleteFileFn             func(ctx context.Context, image_id string) error
	deleteMultipleFileFn     func(ctx context.Context, images []string) error
	getPresignedFileObjectFn func() (st.UrlResult, error)
}

func (m *mockFileManager) UploadFile(fileName string, fileSize int64, prefix string, ctx context.Context, reader io.Reader) (string, error) {
	if m.uploadFileFn != nil {
		return m.uploadFileFn(fileName, fileSize, prefix, ctx, reader)
	}
	return "uploaded-" + fileName, nil
}
func (m *mockFileManager) DeleteFile(ctx context.Context, image_id string) error {
	if m.deleteFileFn != nil {
		return m.deleteFileFn(ctx, image_id)
	}
	return nil
}
func (m *mockFileManager) DeleteMultipleFile(ctx context.Context, images []string) error {
	if m.deleteMultipleFileFn != nil {
		return m.deleteMultipleFileFn(ctx, images)
	}
	return nil
}
func (m *mockFileManager) GetPresignedFileObject() (st.UrlResult, error) {
	if m.getPresignedFileObjectFn != nil {
		return m.getPresignedFileObjectFn()
	}
	return st.UrlResult{}, nil
}

// ----------------------------------------------------
// Mock Payment Service
// ----------------------------------------------------
type mockPaymentService struct {
	createPaymentNonTDSFn  func(req structs.CreateOrderDTO) bool
	createPaymentWith3DSFn func(data structs.CreateOrderDTO) (string, error)
}

func (m *mockPaymentService) CreatePaymentNonTDS(req structs.CreateOrderDTO) bool {
	if m.createPaymentNonTDSFn != nil {
		return m.createPaymentNonTDSFn(req)
	}
	return true
}

func (m *mockPaymentService) CreatePaymentWith3DS(data structs.CreateOrderDTO) (string, error) {
	if m.createPaymentWith3DSFn != nil {
		return m.createPaymentWith3DSFn(data)
	}
	return "<html>3ds</html>", nil
}

