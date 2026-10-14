package main

import (
	"context"
	Customerrors "ecommerceapi/CustomErrors"
	repo "ecommerceapi/Database/Repository"
	httphandlers "ecommerceapi/HttpHandlers"
	services "ecommerceapi/Services"
	"ecommerceapi/Types"
	"ecommerceapi/Utils"
	"fmt"
	"log"
	"net"
	"net/http"
	"os"
	"slices"
	"strconv"
	"strings"
	"time"

	"github.com/JspBack/iyzipay-go"
	"github.com/golang-jwt/jwt/v5"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/minio/minio-go/v7"
	"github.com/minio/minio-go/v7/pkg/credentials"
)

// JWT Secret Key Variable
var secret_key string

func main() {

	secret_key = os.Getenv("SECRET_KEY")
	var (
		// DB Enviroment Variables
		db_name     string = os.Getenv("DB_NAME")
		db_host     string = os.Getenv("DB_HOST")
		db_password string = os.Getenv("DB_PASSWORD")
		db_user     string = os.Getenv("DB_USER")
		db_port     string = os.Getenv("DB_PORT")
		db_sslmode  string = os.Getenv("DB_SSLMODE")
		// Storage Env Variables
		accessKeyID      string = os.Getenv("MINIO_ROOT_USER")
		secretKey        string = os.Getenv("MINIO_ROOT_PASSWORD")
		storageAPIUrl    string = os.Getenv("MINIO_ENDPOINT")
		storagePublicUrl string = os.Getenv("MINIO_PUBLIC_URL")
		bucketName       string = os.Getenv("MINIO_BUCKET_NAME")
		location         string = os.Getenv("MINIO_LOCATION")
		port             string = os.Getenv("API_PORT")
		sslmodevar       string = os.Getenv("STORAGE_SSLMODE_IS_ACTIVE")
		// Iyzipay Env Variables
		iyziApiKey    string = os.Getenv("IYZIPAY_API_KEY")
		iyziApiSecret string = os.Getenv("IYZIPAY_API_SECRET")
		// Google Oauth Env Variables
		googleClientId string = os.Getenv("CLIENT_ID")
	)
	sslmode, err := strconv.ParseBool(sslmodevar)
	// Api settings
	port = ":" + port

	var connectionString string = fmt.Sprintf("user=%s dbname=%s password=%s host=%s port=%s sslmode=%s",
		db_user,
		db_name,
		db_password,
		db_host,
		db_port,
		db_sslmode)

	poolConfig, err := pgxpool.ParseConfig(connectionString)
	if err != nil {
		log.Fatalf("Veritabanı yapılandırma hatası: %v", err)
	}
	poolConfig.MaxConns = 25
	poolConfig.MinConns = 5
	poolConfig.MaxConnLifetime = 1 * time.Hour
	poolConfig.MaxConnIdleTime = 30 * time.Minute
	poolConfig.HealthCheckPeriod = 1 * time.Minute

	dbpool, err := pgxpool.NewWithConfig(context.Background(), poolConfig)
	if err != nil {
		log.Fatal(err)
	}

	// DB Tables Init
	ProductTable := repo.NewProductStore(dbpool)
	ImageTable := repo.NewImageStore(dbpool)
	TokenTable := repo.NewTokenStore(dbpool)

	ProductCategoryStore := repo.NewProductCategoryStore(dbpool)
	ProductProductCategoryStore := repo.NewProductProductCategoryStore(dbpool)
	UserTable := repo.NewUserStore(dbpool)
	AddressTable := repo.NewAddressStore(dbpool)
	CartTable := repo.NewCartStore(dbpool)
	OrderTable := repo.NewOrderStore(dbpool)
	ProductOrderTable := repo.NewProductOrderStore(dbpool)
	// repo.NewPaymentStore(dbpool)
	CartProductsTable := repo.NewCartProductStore(dbpool)

	defer dbpool.Close()

	// Minio Init
	minioClient := MinIOMain(accessKeyID, secretKey, sslmode, storageAPIUrl)
	CreateBucketIfNotExist(minioClient, context.Background(), bucketName, location)
	CheckBucketForPublic(minioClient, bucketName, context.Background())
	// Iyzipay Init

	client, err := iyzipay.New(iyziApiKey, iyziApiSecret)
	if err != nil {
		fmt.Println(err)
	}

	log.Println(client)
	// Google Oauth init

	// Services Init
	fileManagerService := services.NewFileManagerService(minioClient, bucketName)
	imageService := services.NewImageService(ImageTable, *fileManagerService)
	tokenService := services.NewTokenService(TokenTable, secret_key)
	userService := services.NewUserService(UserTable, *tokenService, CartTable, AddressTable, googleClientId)
	cartService := services.NewCartService(CartTable, CartProductsTable)
	paymentService := services.NewPaymentService(client)
	addressService := services.NewAddressService(AddressTable)
	orderService := services.NewOrderService(cartService, CartProductsTable, UserTable,
		paymentService, OrderTable, ProductOrderTable, ProductTable, ProductProductCategoryStore,
		storageAPIUrl, storagePublicUrl, bucketName, *addressService)
	productService := services.NewProductService(ProductTable, ProductCategoryStore,
		ProductProductCategoryStore, storageAPIUrl, storagePublicUrl, bucketName, *imageService, *orderService)
	categoryService := services.NewCategoryService(ProductCategoryStore)

	// Handler init
	producthandler := httphandlers.NewProductHandler(*productService)
	userhandler := httphandlers.NewUserHandler(*userService)
	tokenhandler := httphandlers.NewTokenHandler(*tokenService)
	adminhandler := httphandlers.NewAdminHandler(*userService)
	carthandler := httphandlers.NewCartHandler(*cartService)
	orderhandler := httphandlers.NewOrderHandler(*orderService)
	addresshandler := httphandlers.NewAddressHandler(*addressService)
	categoryhandler := httphandlers.NewCategoryHandler(*categoryService)

	err = dbpool.Ping(context.Background())
	if err != nil {
		log.Fatalf("Veritabanına bağlanılamadı: %v", err)
	}

	fmt.Println("Veritabanına başarıyla bağlanıldı!")
	ips := getLocalIPs()
	log.Printf(" aa Makinenin aktif yerel IPv4 adresleri: %v", ips)

	// Router
	mux := http.NewServeMux()

	// get product
	mux.HandleFunc("GET /api/admin/getallproducts", UseAuthentication(UseAuthorization("admin")(producthandler.Getallproducts)))
	mux.HandleFunc("GET /api/public/getallproducts", producthandler.PubGetAllProducts)
	mux.HandleFunc("GET /api/public/getproductbyid", producthandler.GetProductById)
	mux.HandleFunc("GET /api/admin/getproductbyid", UseAuthentication(UseAuthorization("admin")(producthandler.GetProductByIdForAdmin)))
	mux.HandleFunc("GET /api/public/getproductsbycategory", producthandler.GetProductsByCategoryId)
	mux.HandleFunc("GET /api/admin/getproductcount", UseAuthentication(UseAuthorization("admin")(producthandler.GetProductCountAdmin)))

	// sıralamalı+ filtre get products
	mux.HandleFunc("GET /api/public/getfilteredproducts", producthandler.GetFilteredProducts)

	// category
	mux.HandleFunc("GET /api/public/getcategories", categoryhandler.GetAll)
	mux.HandleFunc("POST /api/admin/createcategory", UseAuthentication(UseAuthorization("admin")(categoryhandler.AddNewCategory)))
	mux.HandleFunc("DELETE /api/admin/deletecategory", UseAuthentication(UseAuthorization("admin")(categoryhandler.DeleteCategory)))
	mux.HandleFunc("PUT /api/admin/updatecategory", UseAuthentication(UseAuthorization("admin")(categoryhandler.UpdateCategory)))

	// product CRUD endpoints
	mux.HandleFunc("POST /api/admin/createproduct", UseAuthentication(UseAuthorization("admin")(producthandler.CreateProduct)))
	mux.HandleFunc("DELETE /api/admin/deleteproduct", UseAuthentication(UseAuthorization("admin")(producthandler.DeleteProduct)))
	mux.HandleFunc("PUT /api/admin/updateproduct", UseAuthentication(UseAuthorization("admin")(producthandler.UpdateProduct)))
	mux.HandleFunc("PATCH /api/admin/patchproduct", UseAuthentication(UseAuthorization("admin")(producthandler.PatchProduct)))
	mux.HandleFunc("POST /api/admin/addimage", UseAuthentication(UseAuthorization("admin")(producthandler.AddProductImage)))
	mux.HandleFunc("PUT /api/admin/product/publishstate", UseAuthentication(UseAuthorization("admin")(producthandler.UpdatePublishState)))

	// Rate Limiter tanımları (Bcrypt DoS ve Brute-Force Koruması)
	// Login: Anlık 5 istek kapasitesi (Burst), dakikada 10 istek dolum hızı (10/60 sn)
	loginLimiter := Utils.NewIPRateLimiter(5, 10.0/60.0, 3*time.Minute)
	defer loginLimiter.Stop()

	// Register: Anlık 3 istek kapasitesi (Burst), dakikada 5 istek dolum hızı (5/60 sn)
	registerLimiter := Utils.NewIPRateLimiter(3, 5.0/60.0, 3*time.Minute)
	defer registerLimiter.Stop()

	// login and register endpoints
	mux.HandleFunc("GET /api/user/me", UseAuthentication(userhandler.GetUserRole))
	mux.HandleFunc("POST /api/user/registeruser", Utils.RateLimitMiddleware(registerLimiter)(userhandler.Register))
	mux.HandleFunc("POST /api/admin/login", Utils.RateLimitMiddleware(loginLimiter)(adminhandler.Login))
	mux.HandleFunc("POST /api/user/login", Utils.RateLimitMiddleware(loginLimiter)(userhandler.Login))
	mux.HandleFunc("POST /api/user/loginwithgoogle", userhandler.LoginWithGoogle)

	// user count
	mux.HandleFunc("GET /api/admin/getusercount", UseAuthentication(UseAuthorization("admin")(userhandler.GetUserCount)))

	// order endpoints
	mux.HandleFunc("GET /api/admin/getordercountforwait", UseAuthentication(UseAuthorization("admin")(orderhandler.GetOrderCountForWaiting)))
	mux.HandleFunc("POST /api/user/createorder", UseAuthentication(UseAuthorization("user")(orderhandler.CreateOrder)))
	mux.HandleFunc("GET /api/user/getorderbyuserid", UseAuthentication(UseAuthorization("user")(orderhandler.GetOrdersByUserId)))
	mux.HandleFunc("GET /api/user/getorderdetails", UseAuthentication(UseAuthorization("user", "admin")(orderhandler.GetOrderDetails)))
	mux.HandleFunc("GET /api/admin/getallorders", UseAuthentication(UseAuthorization("admin")(orderhandler.GetAllOrders)))
	mux.HandleFunc("PUT /api/admin/updateorderstate", UseAuthentication(UseAuthorization("admin")(orderhandler.UpdateOrderState)))
	mux.HandleFunc("PUT /api/admin/updatestates", UseAuthentication(UseAuthorization("admin")(orderhandler.UpdateMultiState)))
	mux.HandleFunc("GET /api/admin/getordersbystatus", UseAuthentication(UseAuthorization("admin")(orderhandler.GetOrdersByStatus)))
	mux.HandleFunc("GET /api/admin/searchorder", UseAuthentication(UseAuthorization("admin")(orderhandler.SearchOrder)))
	// user cart endpoints
	mux.HandleFunc("GET /api/public/getproductsbyids", UseAuthentication(UseAuthorization("user")(producthandler.GetManyProducts)))
	mux.HandleFunc("GET /api/user/getcart", UseAuthentication(UseAuthorization("user")(carthandler.GetCart)))
	mux.HandleFunc("POST /api/user/addtocart", UseAuthentication(UseAuthorization("user")(carthandler.AddToCart)))
	mux.HandleFunc("DELETE /api/user/removefromcart", UseAuthentication(UseAuthorization("user")(carthandler.RemoveFromCart)))
	mux.HandleFunc("PUT /api/user/updatecartitem", UseAuthentication(UseAuthorization("user")(carthandler.UpdateCart)))

	// refresh token endpoint
	mux.HandleFunc("POST /api/refreshtoken", ValidateRefToken(tokenhandler.GetNewAccessToken))

	// address endpoints
	mux.HandleFunc("GET /api/user/getaddresses", UseAuthentication(UseAuthorization("user")(addresshandler.GetAddress)))
	mux.HandleFunc("PUT /api/user/updateaddress", UseAuthentication(UseAuthorization("user")(addresshandler.UpdateAddress)))
	mux.HandleFunc("POST /api/user/createaddress", UseAuthentication(UseAuthorization("user")(addresshandler.CreateAddress)))
	mux.HandleFunc("DELETE /api/user/deleteaddress", UseAuthentication(UseAuthorization("user")(addresshandler.DeleteAddress)))

	// search
	mux.HandleFunc("GET /api/user/searchproduct", producthandler.SearchProduct)
	mux.HandleFunc("GET /api/admin/products/search", producthandler.SearchProductAdmin)

	fmt.Println("Server 7116 portu üzerinde çalışıyor")

	log.Fatal(http.ListenAndServe(port, enableCORS(mux)))
}

func enableCORS(next http.Handler) http.Handler {

	allowedOrigins := map[string]bool{
		"https://ecommerce.yusufsevinc.work": true,
		"http://localhost:5173":              true,
	}
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {

		origin := r.Header.Get("Origin")

		if allowedOrigins[origin] {
			w.Header().Set("Access-Control-Allow-Origin", origin)
		}
		// React uygulamasının adresi
		w.Header().Set("Access-Control-Allow-Methods", "POST, GET, OPTIONS, PUT, DELETE,PATCH")
		w.Header().Set("Access-Control-Allow-Headers", "Accept, Content-Type, Content-Length, Accept-Encoding, X-CSRF-Token, Authorization")
		w.Header().Set("Access-Control-Allow-Credentials", "true")
		// Preflight (OPTIONS) isteği gelirse hemen 200 dön
		if r.Method == "OPTIONS" {
			return
		}

		next.ServeHTTP(w, r)
	})
}

func MinIOMain(accessKey string, secretKey string, sslMode bool, endpoint string) *minio.Client {

	// Initialize minio client object.
	minioClient, err := minio.New(endpoint, &minio.Options{
		Creds:  credentials.NewStaticV4(accessKey, secretKey, ""),
		Secure: sslMode,
	})
	if err != nil {
		log.Fatalln("bağlantı kurulamadı %w", err)
	}

	log.Printf("%#v\n", minioClient) // minioClient is now set up
	return minioClient
}

func CreateBucketIfNotExist(minioClient *minio.Client, ctx context.Context, bucketName string, location string) {

	err := minioClient.MakeBucket(ctx, bucketName, minio.MakeBucketOptions{Region: location})
	if err != nil {
		// Check to see if we already own this bucket (which happens if you run this twice)
		exists, errBucketExists := minioClient.BucketExists(ctx, bucketName)
		if errBucketExists == nil && exists {
			log.Printf("We already own %s\n", bucketName)
			return
		} else {
			log.Fatalln(err)
		}
	} else {
		log.Printf("Successfully created %s\n", bucketName)
		return
	}
}

func CheckBucketForPublic(minioClient *minio.Client, bucketName string, ctx context.Context) error {
	policy := fmt.Sprintf(`{
		"Version": "2012-10-17",
		"Statement": [
			{
				"Action": ["s3:GetObject"],
				"Effect": "Allow",
				"Principal": {"AWS": ["*"]},
				"Resource": ["arn:aws:s3:::%s/*"]
			}
		]
	}`, bucketName)

	// Politikayı bucket'a uygulama
	err := minioClient.SetBucketPolicy(ctx, bucketName, policy)
	if err != nil {
		log.Fatalln("Bucket politikası ayarlanamadı:", err)
	}

	fmt.Println("Bucket başarıyla public (herkese açık) hale getirildi.")
	return nil
}

// userClaims den isim değişti, Types klasörüne taşındı

func UseAuthentication(next http.HandlerFunc) http.HandlerFunc {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {

		authHeader := r.Header.Get("Authorization")

		if authHeader == "" {

			cookie, err := r.Cookie("access_token")
			if err != nil {
				log.Println("cookie okunamadı")
				http.Error(w, "Hatalı Token", http.StatusUnauthorized)
				return
			}

			tokenString := cookie.Value
			claims, err1 := ValidateToken(tokenString)

			if err1 == nil {
				ctx := context.WithValue(r.Context(), Types.UserContextKey, claims)
				// Yeni context ile isteği (request) güncelle
				r = r.WithContext(ctx)
				next.ServeHTTP(w, r)
				return
			}

			log.Println("Authorization Başlığı eksik %w", err1)
			http.Error(w, "Hatalı Token", http.StatusUnauthorized)
			return
		}

		parts := strings.Split(authHeader, " ")
		if len(parts) != 2 || parts[0] != "Bearer" {
			http.Error(w, "Geçersiz token formatı", http.StatusUnauthorized)
			return
		}

		tokenString := parts[1]

		claims, err := ValidateToken(tokenString)
		if err != nil {
			log.Println("Authorization Başlığı eksik")
			http.Error(w, "Hatalı Token", http.StatusUnauthorized)
			return
		}
		ctx := context.WithValue(r.Context(), Types.UserContextKey, claims)
		// Yeni context ile isteği (request) güncelle
		r = r.WithContext(ctx)
		next.ServeHTTP(w, r)
	})
}

func UseAuthorization(requiredRoles ...string) func(http.HandlerFunc) http.HandlerFunc {
	return func(next http.HandlerFunc) http.HandlerFunc {
		return func(w http.ResponseWriter, r *http.Request) {
			// Context'ten kimlik doğrulama adımında eklenen claims verisini çek
			claims, ok := r.Context().Value(Types.UserContextKey).(jwt.MapClaims)
			if !ok {
				// Bu hatanın alınması JWTMiddleware'in çalışmadığı veya sırasının yanlış olduğu anlamına gelir
				http.Error(w, "Yetki verisi bulunamadı", http.StatusInternalServerError)
				return
			}

			// JWT içindeki "role" alanını oku
			userRole, ok := claims["role"].(string)

			// Rol bilgisi yoksa veya istenen rolle eşleşmiyorsa 403 Forbidden döndür
			if !ok || !slices.Contains(requiredRoles, userRole) {
				http.Error(w, "Bu işlemi gerçekleştirmek için yetkiniz yok", http.StatusForbidden)
				return
			}

			// Yetkilendirme başarılı, isteği asıl handler'a ilet
			next.ServeHTTP(w, r)
		}
	}
}

// parse token string and validate it
func ValidateToken(tokenString string) (jwt.MapClaims, error) {

	token, err := jwt.Parse(tokenString, func(token *jwt.Token) (any, error) {
		// Beklenen imzalama algoritmasının kullanıldığını kontrol et
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("beklenmeyen imzalama yöntemi: %v", token.Header["alg"])
		}
		return []byte(secret_key), nil
	})

	if err != nil || !token.Valid {
		return nil, Customerrors.ErrInvalidToken
	}

	// Token geçerliyse payload (claims) bilgilerini Context'e ekle
	if claims, ok := token.Claims.(jwt.MapClaims); ok {
		return claims, nil
	} else {
		log.Println("Token payload'u okunamadı")
		return nil, Customerrors.ErrInvalidToken
	}
}

func ValidateRefToken(next http.HandlerFunc) http.HandlerFunc {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		var refreshTokenString string
		// 1. Önce Header Kontrol Et (Mobil Cihazlar veya özel istemciler için)
		authHeader := r.Header.Get("Authorization")
		if authHeader != "" {
			parts := strings.Split(authHeader, " ")
			if len(parts) == 2 && parts[0] == "Bearer" {
				refreshTokenString = parts[1]
			}
		}
		// 2. Header yoksa Cookie Kontrol Et (Web Tarayıcıları için)
		if refreshTokenString == "" {
			cookie, err := r.Cookie("refresh_token")
			if err != nil {
				log.Println("Refresh token bulunamadı (Header veya Cookie)")
				http.Error(w, "Yetkisiz Erişim", http.StatusUnauthorized)
				return
			}
			refreshTokenString = cookie.Value
		}
		// 3. Token Doğrulama
		claims, err := ValidateToken(refreshTokenString)
		if err != nil {
			log.Println("Geçersiz Refresh Token:", err)
			http.Error(w, "Geçersiz Refresh Token", http.StatusUnauthorized)
			return
		}
		// 4. Context'i Güncelle ve Devam Et
		ctx := context.WithValue(r.Context(), Types.UserContextKey, claims)
		r = r.WithContext(ctx)
		next.ServeHTTP(w, r)
	})
}

// getLocalIPs, makineye atanmış loopback (127.0.0.1) harici IPv4 adreslerini döndürür.
func getLocalIPs() []string {
	var ips []string
	addrs, err := net.InterfaceAddrs()
	if err != nil {
		return ips
	}

	for _, address := range addrs {
		if ipnet, ok := address.(*net.IPNet); ok && !ipnet.IP.IsLoopback() {
			if ipnet.IP.To4() != nil {
				ips = append(ips, ipnet.IP.String())
			}
		}
	}
	return ips
}
