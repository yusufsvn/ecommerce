package main

import (
	"bytes"
	"context"
	models "ecommerceapi/Database/Models"
	repo "ecommerceapi/Database/Repository"
	httphandlers "ecommerceapi/HttpHandlers"
	services "ecommerceapi/Services"
	structs "ecommerceapi/Structs"
	"ecommerceapi/Utils"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"log"
	"net/http"
	"net/http/httptest"
	"os"
	"path/filepath"
	"runtime"
	"testing"
	"time"

	"github.com/golang-migrate/migrate/v4"
	_ "github.com/golang-migrate/migrate/v4/database/postgres"
	_ "github.com/golang-migrate/migrate/v4/source/file"
	"github.com/google/uuid"
	"github.com/jackc/pgx/v5/pgxpool"
	_ "github.com/jackc/pgx/v5/stdlib"
	"github.com/testcontainers/testcontainers-go"
	"github.com/testcontainers/testcontainers-go/wait"
)

var (
	integrationDBPool    *pgxpool.Pool
	integrationContainer testcontainers.Container
)

const (
	intDbName = "integration_db"
	intDbUser = "test_user"
	intDbPass = "test_password"
)

func TestMain(m *testing.M) {
	secret_key = "integration-test-secret-key-32-chars!!"

	ctx, cancel := context.WithTimeout(context.Background(), 90*time.Second)
	defer cancel()

	container, pool, dbAddr, err := startIntegrationPostgres(ctx)
	if err != nil {
		log.Fatalf("Integration DB container başlatılamadı: %v", err)
	}
	integrationContainer = container
	integrationDBPool = pool

	if err := runMigrations(dbAddr); err != nil {
		log.Fatalf("Integration DB migrasyonları çalıştırılamadı: %v", err)
	}

	code := m.Run()

	integrationDBPool.Close()
	_ = integrationContainer.Terminate(context.Background())

	os.Exit(code)
}

func startIntegrationPostgres(ctx context.Context) (testcontainers.Container, *pgxpool.Pool, string, error) {
	env := map[string]string{
		"POSTGRES_PASSWORD": intDbPass,
		"POSTGRES_USER":     intDbUser,
		"POSTGRES_DB":       intDbName,
	}

	req := testcontainers.GenericContainerRequest{
		ContainerRequest: testcontainers.ContainerRequest{
			Image:        "postgres:17.5-alpine",
			ExposedPorts: []string{"5432/tcp"},
			Env:          env,
			WaitingFor:   wait.ForLog("database system is ready to accept connections"),
		},
		Started: true,
	}

	container, err := testcontainers.GenericContainer(ctx, req)
	if err != nil {
		return nil, nil, "", fmt.Errorf("container oluşturulamadı: %w", err)
	}

	mappedPort, err := container.MappedPort(ctx, "5432")
	if err != nil {
		return container, nil, "", fmt.Errorf("port okunamadı: %w", err)
	}

	time.Sleep(time.Second)

	dbAddr := fmt.Sprintf("localhost:%s", mappedPort.Port())
	connStr := fmt.Sprintf("postgres://%s:%s@%s/%s?sslmode=disable", intDbUser, intDbPass, dbAddr, intDbName)

	pool, err := pgxpool.New(ctx, connStr)
	if err != nil {
		return container, nil, dbAddr, fmt.Errorf("pgxpool bağlantısı kurulamadı: %w", err)
	}

	// PostgreSQL hazır olana kadar ping döngüsü
	for i := 0; i < 10; i++ {
		if err := pool.Ping(ctx); err == nil {
			break
		}
		time.Sleep(500 * time.Millisecond)
	}

	return container, pool, dbAddr, nil
}

func runMigrations(dbAddr string) error {
	_, filename, _, ok := runtime.Caller(0)
	if !ok {
		return errors.New("runtime.Caller başarısız")
	}

	dir := filepath.Dir(filename)
	migrationsPath := filepath.ToSlash(filepath.Join(dir, "Database", "Migrations"))

	dbURL := fmt.Sprintf("postgres://%s:%s@%s/%s?sslmode=disable", intDbUser, intDbPass, dbAddr, intDbName)
	m, err := migrate.New(fmt.Sprintf("file:%s", migrationsPath), dbURL)
	if err != nil {
		return fmt.Errorf("migrator oluşturulamadı: %w", err)
	}
	defer m.Close()

	if err := m.Up(); err != nil && !errors.Is(err, migrate.ErrNoChange) {
		return fmt.Errorf("migration up hatası: %w", err)
	}

	log.Println("Entegrasyon veritabanı migrasyonları tamamlandı")
	return nil
}

// Mock external services for integration
type mockIntegrationPaymentService struct {
	CreatePaymentNonTDSFn  func(req structs.CreateOrderDTO) bool
	CreatePaymentWith3DSFn func(data structs.CreateOrderDTO) (string, error)
}

func (m *mockIntegrationPaymentService) CreatePaymentNonTDS(req structs.CreateOrderDTO) bool {
	if m.CreatePaymentNonTDSFn != nil {
		return m.CreatePaymentNonTDSFn(req)
	}
	return true
}

func (m *mockIntegrationPaymentService) CreatePaymentWith3DS(data structs.CreateOrderDTO) (string, error) {
	if m.CreatePaymentWith3DSFn != nil {
		return m.CreatePaymentWith3DSFn(data)
	}
	return "https://mock-payment.com/3ds", nil
}

type mockIntegrationFileManager struct {
	UploadFileFn             func(fileName string, fileSize int64, prefix string, ctx context.Context, reader io.Reader) (string, error)
	DeleteFileFn             func(ctx context.Context, image_id string) error
	DeleteMultipleFileFn     func(ctx context.Context, images []string) error
	GetPresignedFileObjectFn func() (structs.UrlResult, error)
}

func (m *mockIntegrationFileManager) UploadFile(fileName string, fileSize int64, prefix string, ctx context.Context, reader io.Reader) (string, error) {
	if m.UploadFileFn != nil {
		return m.UploadFileFn(fileName, fileSize, prefix, ctx, reader)
	}
	return fileName, nil
}

func (m *mockIntegrationFileManager) DeleteFile(ctx context.Context, image_id string) error {
	if m.DeleteFileFn != nil {
		return m.DeleteFileFn(ctx, image_id)
	}
	return nil
}

func (m *mockIntegrationFileManager) DeleteMultipleFile(ctx context.Context, images []string) error {
	if m.DeleteMultipleFileFn != nil {
		return m.DeleteMultipleFileFn(ctx, images)
	}
	return nil
}

func (m *mockIntegrationFileManager) GetPresignedFileObject() (structs.UrlResult, error) {
	if m.GetPresignedFileObjectFn != nil {
		return m.GetPresignedFileObjectFn()
	}
	return structs.UrlResult{}, nil
}

// TestApp encapsulates router, server, and DB pool
type testApp struct {
	Router         *http.ServeMux
	Server         *httptest.Server
	DB             *pgxpool.Pool
	ProductStore   repo.ProductFunctions
	CategoryStore  repo.ProductCategoryFunctions
	PPCStore       repo.ProductProductCategoryFunctions
	UserStore      repo.UserFunctions
	TokenStore     repo.TokenFunctions
	CartStore      repo.CartFunctions
	CartProdStore  repo.CartProductFunctions
	OrderStore     repo.OrderFunctions
	ProdOrderStore repo.ProductOrderFunctions
	AddressStore   repo.AddressFunctions
	ImageStore     repo.ImageFunctions
	PaymentSvc     *mockIntegrationPaymentService
	FileMgrSvc     *mockIntegrationFileManager
}

func setupIntegrationApp(t testing.TB) *testApp {
	dbpool := integrationDBPool

	prodStore := repo.NewProductStore(dbpool)
	imgStore := repo.NewImageStore(dbpool)
	tokenStore := repo.NewTokenStore(dbpool)
	catStore := repo.NewProductCategoryStore(dbpool)
	ppcStore := repo.NewProductProductCategoryStore(dbpool)
	userStore := repo.NewUserStore(dbpool)
	addrStore := repo.NewAddressStore(dbpool)
	cartStore := repo.NewCartStore(dbpool)
	orderStore := repo.NewOrderStore(dbpool)
	prodOrderStore := repo.NewProductOrderStore(dbpool)
	cartProdStore := repo.NewCartProductStore(dbpool)

	fileMgrMock := &mockIntegrationFileManager{}
	imageSvc := services.NewImageService(imgStore, fileMgrMock)
	tokenSvc := services.NewTokenService(tokenStore, secret_key)
	userSvc := services.NewUserService(userStore, *tokenSvc, cartStore, addrStore, "mock-google-client-id")
	cartSvc := services.NewCartService(cartStore, cartProdStore)
	paymentMock := &mockIntegrationPaymentService{}
	addrSvc := services.NewAddressService(addrStore)
	orderSvc := services.NewOrderService(
		cartSvc,
		cartProdStore,
		userStore,
		paymentMock,
		orderStore,
		prodOrderStore,
		prodStore,
		ppcStore,
		"http://localhost:9000",
		"http://localhost:9000",
		"ecommerce-bucket",
		*addrSvc,
	)
	productSvc := services.NewProductService(
		prodStore,
		catStore,
		ppcStore,
		"http://localhost:9000",
		"http://localhost:9000",
		"ecommerce-bucket",
		*imageSvc,
	)
	catSvc := services.NewCategoryService(catStore)

	prodHandler := httphandlers.NewProductHandler(*productSvc)
	userHandler := httphandlers.NewUserHandler(*userSvc)
	tokenHandler := httphandlers.NewTokenHandler(*tokenSvc)
	adminHandler := httphandlers.NewAdminHandler(*userSvc)
	cartHandler := httphandlers.NewCartHandler(*cartSvc)
	orderHandler := httphandlers.NewOrderHandler(*orderSvc)
	addrHandler := httphandlers.NewAddressHandler(*addrSvc)
	catHandler := httphandlers.NewCategoryHandler(*catSvc)

	mux := http.NewServeMux()

	// Product endpoints
	mux.HandleFunc("GET /api/admin/getallproducts", UseAuthentication(UseAuthorization("admin")(prodHandler.Getallproducts)))
	mux.HandleFunc("GET /api/public/getallproducts", prodHandler.PubGetAllProducts)
	mux.HandleFunc("GET /api/public/getproductbyid", prodHandler.GetProductById)
	mux.HandleFunc("GET /api/public/getproductsbycategory", prodHandler.GetProductsByCategoryId)
	mux.HandleFunc("GET /api/admin/getproductcount", UseAuthentication(UseAuthorization("admin")(prodHandler.GetProductCountAdmin)))
	mux.HandleFunc("GET /api/public/getfilteredproducts", prodHandler.GetFilteredProducts)

	// Category endpoints
	mux.HandleFunc("GET /api/public/getcategories", catHandler.GetAll)
	mux.HandleFunc("POST /api/admin/createcategory", UseAuthentication(UseAuthorization("admin")(catHandler.AddNewCategory)))
	mux.HandleFunc("DELETE /api/admin/deletecategory", UseAuthentication(UseAuthorization("admin")(catHandler.DeleteCategory)))
	mux.HandleFunc("PUT /api/admin/updatecategory", UseAuthentication(UseAuthorization("admin")(catHandler.UpdateCategory)))

	// Product CRUD endpoints
	mux.HandleFunc("POST /api/admin/createproduct", UseAuthentication(UseAuthorization("admin")(prodHandler.CreateProduct)))
	mux.HandleFunc("DELETE /api/admin/deleteproduct", UseAuthentication(UseAuthorization("admin")(prodHandler.DeleteProduct)))
	mux.HandleFunc("PUT /api/admin/updateproduct", UseAuthentication(UseAuthorization("admin")(prodHandler.UpdateProduct)))
	mux.HandleFunc("PATCH /api/admin/patchproduct", UseAuthentication(UseAuthorization("admin")(prodHandler.PatchProduct)))
	mux.HandleFunc("POST /api/admin/addimage", UseAuthentication(UseAuthorization("admin")(prodHandler.AddProductImage)))
	mux.HandleFunc("PUT /api/admin/product/publishstate", UseAuthentication(UseAuthorization("admin")(prodHandler.UpdatePublishState)))

	// Rate Limiter tanımları
	loginLimiter := Utils.NewIPRateLimiter(10, 5.0, 1*time.Minute)
	registerLimiter := Utils.NewIPRateLimiter(10, 5.0, 1*time.Minute)

	// Auth & User endpoints
	mux.HandleFunc("GET /api/user/me", UseAuthentication(userHandler.GetUserRole))
	mux.HandleFunc("POST /api/user/registeruser", Utils.RateLimitMiddleware(registerLimiter)(userHandler.Register))
	mux.HandleFunc("POST /api/admin/login", Utils.RateLimitMiddleware(loginLimiter)(adminHandler.Login))
	mux.HandleFunc("POST /api/user/login", Utils.RateLimitMiddleware(loginLimiter)(userHandler.Login))
	mux.HandleFunc("POST /api/user/loginwithgoogle", userHandler.LoginWithGoogle)
	mux.HandleFunc("GET /api/admin/getusercount", UseAuthentication(UseAuthorization("admin")(userHandler.GetUserCount)))

	// Order endpoints
	mux.HandleFunc("GET /api/admin/getordercountforwait", UseAuthentication(UseAuthorization("admin")(orderHandler.GetOrderCountForWaiting)))
	mux.HandleFunc("POST /api/user/createorder", UseAuthentication(UseAuthorization("user")(orderHandler.CreateOrder)))
	mux.HandleFunc("GET /api/user/getorderbyuserid", UseAuthentication(UseAuthorization("user")(orderHandler.GetOrdersByUserId)))
	mux.HandleFunc("GET /api/user/getorderdetails", UseAuthentication(UseAuthorization("user", "admin")(orderHandler.GetOrderDetails)))
	mux.HandleFunc("GET /api/admin/getallorders", UseAuthentication(UseAuthorization("admin")(orderHandler.GetAllOrders)))
	mux.HandleFunc("PUT /api/admin/updateorderstate", UseAuthentication(UseAuthorization("admin")(orderHandler.UpdateOrderState)))
	mux.HandleFunc("PUT /api/admin/updatestates", UseAuthentication(UseAuthorization("admin")(orderHandler.UpdateMultiState)))
	mux.HandleFunc("GET /api/admin/getordersbystatus", UseAuthentication(UseAuthorization("admin")(orderHandler.GetOrdersByStatus)))
	mux.HandleFunc("GET /api/admin/searchorder", UseAuthentication(UseAuthorization("admin")(orderHandler.SearchOrder)))

	// Cart endpoints
	mux.HandleFunc("GET /api/public/getproductsbyids", UseAuthentication(UseAuthorization("user")(prodHandler.GetManyProducts)))
	mux.HandleFunc("GET /api/user/getcart", UseAuthentication(UseAuthorization("user")(cartHandler.GetCart)))
	mux.HandleFunc("POST /api/user/addtocart", UseAuthentication(UseAuthorization("user")(cartHandler.AddToCart)))
	mux.HandleFunc("DELETE /api/user/removefromcart", UseAuthentication(UseAuthorization("user")(cartHandler.RemoveFromCart)))
	mux.HandleFunc("PUT /api/user/updatecartitem", UseAuthentication(UseAuthorization("user")(cartHandler.UpdateCart)))

	// Refresh Token endpoint
	mux.HandleFunc("POST /api/refreshtoken", ValidateRefToken(tokenHandler.GetNewAccessToken))

	// Address endpoints
	mux.HandleFunc("GET /api/user/getaddresses", UseAuthentication(UseAuthorization("user")(addrHandler.GetAddress)))
	mux.HandleFunc("PUT /api/user/updateaddress", UseAuthentication(UseAuthorization("user")(addrHandler.UpdateAddress)))
	mux.HandleFunc("POST /api/user/createaddress", UseAuthentication(UseAuthorization("user")(addrHandler.CreateAddress)))
	mux.HandleFunc("DELETE /api/user/deleteaddress", UseAuthentication(UseAuthorization("user")(addrHandler.DeleteAddress)))

	// Search endpoints
	mux.HandleFunc("GET /api/user/searchproduct", prodHandler.SearchProduct)
	mux.HandleFunc("GET /api/admin/products/search", prodHandler.SearchProductAdmin)

	server := httptest.NewServer(mux)
	t.Cleanup(func() {
		server.Close()
	})

	return &testApp{
		Router:         mux,
		Server:         server,
		DB:             dbpool,
		ProductStore:   prodStore,
		CategoryStore:  catStore,
		PPCStore:       ppcStore,
		UserStore:      userStore,
		TokenStore:     tokenStore,
		CartStore:      cartStore,
		CartProdStore:  cartProdStore,
		OrderStore:     orderStore,
		ProdOrderStore: prodOrderStore,
		AddressStore:   addrStore,
		ImageStore:     imgStore,
		PaymentSvc:     paymentMock,
		FileMgrSvc:     fileMgrMock,
	}
}

// Helper: creates a user in DB and returns access token, refresh token and user ID
func (app *testApp) createTestUser(t testing.TB, email, password, role string) (string, *http.Cookie, *http.Cookie) {
	ctx := context.Background()

	user := models.User{
		Email:    email,
		Password: password,
		Name:     "Test",
		Surname:  "User",
		Phone:    "05551234567",
	}

	insertedUserID, err := app.UserStore.Create(ctx, &user)
	if err != nil {
		t.Fatalf("kullanıcı oluşturulamadı: %v", err)
	}

	if role != "user" && role != "" {
		if _, err := app.DB.Exec(ctx, "UPDATE users SET user_role = $1 WHERE user_id = $2", role, insertedUserID); err != nil {
			t.Fatalf("kullanıcı rolü güncellenemedi: %v", err)
		}
	}

	if _, err := app.CartStore.Create(ctx, insertedUserID); err != nil {
		t.Fatalf("sepet oluşturulamadı: %v", err)
	}

	tokenSvc := services.NewTokenService(app.TokenStore, secret_key)
	familyID := uuid.NewString()
	acExp := time.Now().Add(15 * time.Minute)
	rfExp := time.Now().Add(7 * 24 * time.Hour)

	accessToken, err := tokenSvc.GenToken(role, acExp, familyID, insertedUserID)
	if err != nil {
		t.Fatalf("access token üretilemedi: %v", err)
	}

	refreshToken, err := tokenSvc.GenToken(role, rfExp, familyID, insertedUserID)
	if err != nil {
		t.Fatalf("refresh token üretilemedi: %v", err)
	}

	if err := tokenSvc.SaveRefreshToken(ctx, refreshToken.Token, insertedUserID, rfExp, refreshToken.Jti, familyID); err != nil {
		t.Fatalf("refresh token kaydedilemedi: %v", err)
	}

	accessCookie := &http.Cookie{
		Name:     "access_token",
		Value:    accessToken.Token,
		Path:     "/",
		HttpOnly: true,
	}
	refreshCookie := &http.Cookie{
		Name:     "refresh_token",
		Value:    refreshToken.Token,
		Path:     "/",
		HttpOnly: true,
	}

	return insertedUserID, accessCookie, refreshCookie
}

// Helper: executes HTTP request with optional cookies and json body
func (app *testApp) doRequest(
	t testing.TB,
	method, path string,
	body any,
	cookies ...*http.Cookie,
) *http.Response {
	var bodyReader io.Reader
	if body != nil {
		switch b := body.(type) {
		case io.Reader:
			bodyReader = b
		case string:
			bodyReader = bytes.NewReader([]byte(b))
		case []byte:
			bodyReader = bytes.NewReader(b)
		default:
			jsonData, err := json.Marshal(b)
			if err != nil {
				t.Fatalf("gövde json'a serileştirilemedi: %v", err)
			}
			bodyReader = bytes.NewReader(jsonData)
		}
	}

	req, err := http.NewRequest(method, app.Server.URL+path, bodyReader)
	if err != nil {
		t.Fatalf("istek oluşturulamadı: %v", err)
	}

	if body != nil {
		if _, ok := body.(io.Reader); !ok {
			req.Header.Set("Content-Type", "application/json")
		}
	}

	for _, c := range cookies {
		if c != nil {
			req.AddCookie(c)
		}
	}

	client := app.Server.Client()
	resp, err := client.Do(req)
	if err != nil {
		t.Fatalf("istek gerçekleştirilemedi: %v", err)
	}
	return resp
}
