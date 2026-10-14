package main

import (
	"context"
	models "ecommerceapi/Database/Models"
	repo "ecommerceapi/Database/Repository"
	services "ecommerceapi/Services"
	structs "ecommerceapi/Structs"
	"ecommerceapi/Utils"
	"fmt"
	"io"
	"net/http"
	"testing"
	"time"

	"github.com/google/uuid"
	"golang.org/x/crypto/bcrypt"
)

// -----------------------------------------------------------------------------
// 1. CPU & KRİPTO / YETKİLENDİRME BENCHMARK'LARI (In-Memory & Algoritmik)
// -----------------------------------------------------------------------------

// Bcrypt şifre karşılaştırma maliyetini ölçer (Login işlemi sırasında CPU yükünü gösterir)
func Benchmark_Auth_BcryptCompare(b *testing.B) {
	password := "SecretPass123!"
	hash, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		b.Fatalf("bcrypt generate failed: %v", err)
	}

	b.ReportAllocs()
	b.ResetTimer()

	for i := 0; i < b.N; i++ {
		_ = bcrypt.CompareHashAndPassword(hash, []byte(password))
	}
}

// Bcrypt şifre hash'leme maliyetini ölçer (Kayıt olma anında CPU yükünü gösterir)
func Benchmark_Auth_BcryptHash(b *testing.B) {
	password := []byte("SecretPass123!")

	b.ReportAllocs()
	b.ResetTimer()

	for i := 0; i < b.N; i++ {
		_, _ = bcrypt.GenerateFromPassword(password, bcrypt.DefaultCost)
	}
}

// JWT Token üretme hızını ve bellek tahsisini (allocs/op) ölçer
func Benchmark_Auth_JWT_GenToken(b *testing.B) {
	tokenStore := repo.NewTokenStore(integrationDBPool)
	tokenSvc := services.NewTokenService(tokenStore, secret_key)
	exp := time.Now().Add(15 * time.Minute)
	famID := uuid.NewString()
	uID := uuid.NewString()

	b.ReportAllocs()
	b.ResetTimer()

	for i := 0; i < b.N; i++ {
		_, err := tokenSvc.GenToken("user", exp, famID, uID)
		if err != nil {
			b.Fatalf("GenToken failed: %v", err)
		}
	}
}

// JWT Token imza doğrulama ve claims çıkarma hızını ölçer (Her korumalı istekte çalışır)
func Benchmark_Auth_JWT_ValidateToken(b *testing.B) {
	tokenStore := repo.NewTokenStore(integrationDBPool)
	tokenSvc := services.NewTokenService(tokenStore, secret_key)
	exp := time.Now().Add(15 * time.Minute)
	famID := uuid.NewString()
	uID := uuid.NewString()

	tokenData, err := tokenSvc.GenToken("user", exp, famID, uID)
	if err != nil {
		b.Fatalf("GenToken failed: %v", err)
	}

	b.ReportAllocs()
	b.ResetTimer()

	for i := 0; i < b.N; i++ {
		_, err := ValidateToken(tokenData.Token)
		if err != nil {
			b.Fatalf("ValidateToken failed: %v", err)
		}
	}
}

// Yüksek eşzamanlılıkta (Concurrent/Parallel) JWT doğrulama yük testi
func Benchmark_Parallel_Auth_JWT_ValidateToken(b *testing.B) {
	tokenStore := repo.NewTokenStore(integrationDBPool)
	tokenSvc := services.NewTokenService(tokenStore, secret_key)
	exp := time.Now().Add(15 * time.Minute)
	famID := uuid.NewString()
	uID := uuid.NewString()

	tokenData, err := tokenSvc.GenToken("user", exp, famID, uID)
	if err != nil {
		b.Fatalf("GenToken failed: %v", err)
	}

	b.ReportAllocs()
	b.ResetTimer()

	b.RunParallel(func(pb *testing.PB) {
		for pb.Next() {
			_, err := ValidateToken(tokenData.Token)
			if err != nil {
				b.Fatalf("ValidateToken failed: %v", err)
			}
		}
	})
}

// -----------------------------------------------------------------------------
// 2. YARDIMCI VE HESAPLAMA FONKSİYONLARI BENCHMARK'LARI
// -----------------------------------------------------------------------------

// Sepet toplam tutar hesaplama algoritmasının CPU ve bellek performansını ölçer
func Benchmark_Utils_CalcTotalPrice_10Items(b *testing.B) {
	items := make([]structs.ProductsFromCartDTO, 10)
	for i := 0; i < 10; i++ {
		items[i] = structs.ProductsFromCartDTO{
			ProductID:   uuid.NewString(),
			ProductName: fmt.Sprintf("Ürün %d", i),
			Price:       int64((i + 1) * 150),
			Quantity:    i + 1,
			Category:    "Elektronik",
			ItemType:    "PHYSICAL",
		}
	}

	b.ReportAllocs()
	b.ResetTimer()

	for i := 0; i < b.N; i++ {
		_ = Utils.CalcTotalPrice(items)
	}
}

func Benchmark_Utils_CalcTotalPrice_50Items(b *testing.B) {
	items := make([]structs.ProductsFromCartDTO, 50)
	for i := 0; i < 50; i++ {
		items[i] = structs.ProductsFromCartDTO{
			ProductID:   uuid.NewString(),
			ProductName: fmt.Sprintf("Ürün %d", i),
			Price:       int64((i + 1) * 150),
			Quantity:    i + 1,
			Category:    "Elektronik",
			ItemType:    "PHYSICAL",
		}
	}

	b.ReportAllocs()
	b.ResetTimer()

	for i := 0; i < b.N; i++ {
		_ = Utils.CalcTotalPrice(items)
	}
}

// -----------------------------------------------------------------------------
// 3. HTTP VE VERİTABANI UÇTAN UCA ENTEGRASYON BENCHMARK'LARI
// -----------------------------------------------------------------------------

// Public kategorileri listeleme endpoint'i (GET /api/public/getcategories)
func Benchmark_HTTP_Public_GetCategories(b *testing.B) {
	app := setupIntegrationApp(b)
	ctx := context.Background()

	// 5 adet kategori oluştur
	for i := 0; i < 5; i++ {
		_, _ = app.CategoryStore.Create(ctx, fmt.Sprintf("BenchCat-%d-%s", i, uuid.NewString()[:6]))
	}

	b.ReportAllocs()
	b.ResetTimer()

	for i := 0; i < b.N; i++ {
		resp := app.doRequest(b, http.MethodGet, "/api/public/getcategories", nil)
		_, _ = io.Copy(io.Discard, resp.Body)
		_ = resp.Body.Close()
	}
}

// Public ürünleri listeleme endpoint'i (GET /api/public/getallproducts)
func Benchmark_HTTP_Public_GetProducts(b *testing.B) {
	app := setupIntegrationApp(b)
	ctx := context.Background()

	// Test için ürün tohumlama (seeding)
	now := time.Now()
	for i := 0; i < 5; i++ {
		prod := models.Product{
			ProductID:           uuid.NewString(),
			ProductName:         fmt.Sprintf("Bench Prod %d", i),
			Stock:               100,
			Price:               500,
			ThumbnailImage:      uuid.NewString(),
			Description:         "Benchmark ürünü",
			CreatedAt:           &now,
			UpdatedAt:           &now,
			ItemType:            "PHYSICAL",
			ProductCode:         uuid.NewString()[:8],
			ProductPublishState: "published",
			TechnicalDetails:    []byte(`{}`),
		}
		_ = app.ProductStore.Create(ctx, &prod)
	}

	b.ReportAllocs()
	b.ResetTimer()

	for i := 0; i < b.N; i++ {
		resp := app.doRequest(b, http.MethodGet, "/api/public/getallproducts?limit=10&offset=0", nil)
		_, _ = io.Copy(io.Discard, resp.Body)
		_ = resp.Body.Close()
	}
}

// Eşzamanlı (Parallel) sanal kullanıcıların ürün listeleme yük testi (Throughput / RPS)
func Benchmark_Parallel_HTTP_Public_GetProducts(b *testing.B) {
	app := setupIntegrationApp(b)
	ctx := context.Background()

	now := time.Now()
	for i := 0; i < 5; i++ {
		prod := models.Product{
			ProductID:           uuid.NewString(),
			ProductName:         fmt.Sprintf("Parallel Bench Prod %d", i),
			Stock:               100,
			Price:               500,
			ThumbnailImage:      uuid.NewString(),
			Description:         "Benchmark paralel ürünü",
			CreatedAt:           &now,
			UpdatedAt:           &now,
			ItemType:            "PHYSICAL",
			ProductCode:         uuid.NewString()[:8],
			ProductPublishState: "published",
			TechnicalDetails:    []byte(`{}`),
		}
		_ = app.ProductStore.Create(ctx, &prod)
	}

	b.ReportAllocs()
	b.ResetTimer()

	b.RunParallel(func(pb *testing.PB) {
		client := app.Server.Client()
		for pb.Next() {
			req, err := http.NewRequest(http.MethodGet, app.Server.URL+"/api/public/getallproducts?limit=10&offset=0", nil)
			if err != nil {
				b.Fatalf("req err: %v", err)
			}
			resp, err := client.Do(req)
			if err != nil {
				b.Fatalf("do err: %v", err)
			}
			_, _ = io.Copy(io.Discard, resp.Body)
			_ = resp.Body.Close()
		}
	})
}

// PostgreSQL Full-Text Search (GIN indeksi) arama endpoint'i performansı
func Benchmark_HTTP_User_SearchProduct(b *testing.B) {
	app := setupIntegrationApp(b)
	ctx := context.Background()

	now := time.Now()
	searchTerm := fmt.Sprintf("ArananUrun%s", uuid.NewString()[:4])
	prod := models.Product{
		ProductID:           uuid.NewString(),
		ProductName:         searchTerm,
		Stock:               100,
		Price:               1500,
		ThumbnailImage:      uuid.NewString(),
		Description:         "Arama performans testi açıklaması",
		CreatedAt:           &now,
		UpdatedAt:           &now,
		ItemType:            "PHYSICAL",
		ProductCode:         uuid.NewString()[:8],
		ProductPublishState: "published",
		TechnicalDetails:    []byte(`{}`),
	}
	_ = app.ProductStore.Create(ctx, &prod)

	b.ReportAllocs()
	b.ResetTimer()

	for i := 0; i < b.N; i++ {
		resp := app.doRequest(b, http.MethodGet, fmt.Sprintf("/api/user/searchproduct?search=%s&limit=10", searchTerm), nil)
		_, _ = io.Copy(io.Discard, resp.Body)
		_ = resp.Body.Close()
	}
}

// -----------------------------------------------------------------------------
// 4. VERİTABANI REPOSITORY KATMANI DOĞRUDAN BENCHMARK'LARI
// -----------------------------------------------------------------------------

// ProductStore.GetAll veritabanı sorgusu ve pgx satır dönüştürme hızı
func Benchmark_DB_ProductStore_GetAll(b *testing.B) {
	dbpool := integrationDBPool
	prodStore := repo.NewProductStore(dbpool)
	ctx := context.Background()

	b.ReportAllocs()
	b.ResetTimer()

	for i := 0; i < b.N; i++ {
		_, err := prodStore.GetAll(10, 0, ctx)
		if err != nil {
			b.Fatalf("GetAll failed: %v", err)
		}
	}
}

// ProductStore.SearchProduct FTS veritabanı sorgusu hızı
func Benchmark_DB_ProductStore_Search(b *testing.B) {
	dbpool := integrationDBPool
	prodStore := repo.NewProductStore(dbpool)
	ctx := context.Background()

	b.ReportAllocs()
	b.ResetTimer()

	for i := 0; i < b.N; i++ {
		_, err := prodStore.SearchProduct(ctx, "test", "10")
		if err != nil {
			b.Fatalf("SearchProduct failed: %v", err)
		}
	}
}
