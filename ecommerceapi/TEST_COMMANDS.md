# 🧪 E-Commerce API - Test Çalıştırma Rehberi

Bu dosya; projedeki birim (unit), entegrasyon (integration), performans/benchmark testlerini ve test kapsamı (coverage) raporlarını çalıştırmak için gereken tüm komutları içerir.

---

## 📌 Hızlı Başlangıç

| Hedef | Komut |
| :--- | :--- |
| **Tüm Proje Testlerini Çalıştır** | `go test ./...` |
| **Önbelleksiz (Temiz) Tüm Testler** | `go test -count=1 ./...` |
| **Sadece Entegrasyon Testleri** | `go test -v -run TestIntegration .` |
| **Sadece Birim (Unit) Testleri** | `go test -v ./Database/Repository/... ./Services/... ./HttpHandlers/... ./Utils/...` |
| **Yük & Performans (Benchmark) Testleri** | `go test -run=^$ -bench=. -benchmem .` |
| **Test Kapsamını (Coverage) Gör** | `go test -cover ./...` |

---

## 1. 🧩 Birim (Unit) Test Komutları

Birim testler veritabanı gerektirmeden in-memory mock'lar ile anında çalışır.

### Tüm Birim Testleri Birlikte Çalıştırma:
```powershell
go test -v ./Database/Repository/... ./Services/... ./HttpHandlers/... ./Utils/...
```

### Katman Bazında Ayrı Ayrı Çalıştırma:

- **Repository (Veritabanı Katmanı) Testleri:**
  ```powershell
  go test -v ./Database/Repository/...
  ```

- **Services (İş Mantığı Katmanı) Testleri:**
  ```powershell
  go test -v ./Services/...
  ```

- **HttpHandlers (HTTP Denetleyici Katmanı) Testleri:**
  ```powershell
  go test -v ./HttpHandlers/...
  ```

- **Utils (Yardımcı Fonksiyonlar) Testleri:**
  ```powershell
  go test -v ./Utils/...
  ```

### Belirli Bir Test Fonksiyonunu Çalıştırma:
```powershell
# Örnek: Sadece UserService_Login testini çalıştır
go test -v -run TestUserService_Login ./Services/...

# Örnek: Sadece ProductHandler testlerini çalıştır
go test -v -run TestProductHandler ./HttpHandlers/...
```

---

## 2. 🐳 Entegrasyon (Integration) Test Komutları

> **Not:** Entegrasyon testleri Docker üzerinde geçici bir `postgres:17.5-alpine` konteyneri ayağa kaldırır ve 14 migrasyonu otomatik uygular. Docker Desktop açık olmalıdır.

### Tüm Entegrasyon Testlerini Çalıştırma:
```powershell
go test -v -run TestIntegration .
```

### Modül Bazında Entegrasyon Testleri:

- **Auth & Yetkilendirme Testleri (Register, Login, Token, Admin RBAC):**
  ```powershell
  go test -v -run TestIntegration_Auth .
  ```

- **Ürün & Kategori Testleri (CRUD, Publish State, FTS Search):**
  ```powershell
  go test -v -run "TestIntegration_Product|TestIntegration_Category" .
  ```

- **Sepet (Cart) Testleri:**
  ```powershell
  go test -v -run TestIntegration_Cart .
  ```

- **Adres (Address) Testleri:**
  ```powershell
  go test -v -run TestIntegration_Address .
  ```

- **Sipariş (Order) Yaşam Döngüsü & Stok Testleri:**
  ```powershell
  go test -v -run TestIntegration_Order .
  ```

---

## 3. ⚡ Yük ve Performans (Benchmark) Komutları

Go'nun yerleşik `testing.B` altyapısı ile CPU maliyeti, nanosaniye gecikmesi, bellek tahsisi ve paralel eşzamanlı sanal kullanıcı simülasyonları ölçülür.

### Tüm Benchmark Testlerini Çalıştırma:
```powershell
go test -run=^$ -bench=. -benchmem .
```

### Kategori Bazlı Benchmarklar:

- **CPU & Kripto Maliyeti (Bcrypt Hash & Compare):**
  ```powershell
  go test -run=^$ -bench=Benchmark_Auth_Bcrypt . -benchmem
  ```

- **JWT Token Üretim ve Doğrulama Hızı:**
  ```powershell
  go test -run=^$ -bench=Benchmark_Auth_JWT . -benchmem
  ```

- **Yüksek Eşzamanlılık (Paralel / Concurrent Kullanıcı Yükü):**
  ```powershell
  # Saniyede ~200.000 paralel token doğrulama simülasyonu
  go test -run=^$ -bench=Benchmark_Parallel_Auth_JWT . -benchmem

  # Eşzamanlı HTTP istekleri ve DB sorgusu yük simülasyonu
  go test -run=^$ -bench=Benchmark_Parallel_HTTP -benchmem -benchtime=500x .
  ```

- **Veritabanı FTS ve Sorgu Performansı:**
  ```powershell
  go test -run=^$ -bench=Benchmark_DB -benchmem -benchtime=300x .
  ```

### 🔬 CPU ve Bellek Profil Çıkarma (pprof):
```powershell
# 1. Profil dosyasını üret
go test -run=^$ -bench=Benchmark_Auth_BcryptCompare -cpuprofile=cpu.prof -memprofile=mem.prof .

# 2. Tarayıcıda grafiksel olarak incele (Alev grafiği / Flamegraph)
go tool pprof -http=:8081 cpu.prof
```

---

## 4. 📈 Test Kapsamı (Coverage) Komutları

### Konsolda Paket Bazlı Kapsam Özeti:
```powershell
go test -cover ./...
```

### Detaylı HTML Kapsam Raporu Üretme ve Açma:
```powershell
# Kapsam profilini oluştur
go test -coverprofile=coverage.out ./...

# Tarayıcıda hangi satırların test edildiğini renkli olarak aç
go tool cover -html=coverage.out
```

---

## 5. 💡 Sık Kullanılan Bayraklar (Flags) Anlamları

- `-v` : Verbose mode. Çalışan her alt testin adını ve loglarını ayrıntılı yazdırır.
- `-count=1` : Go'nun test önbelleğini (test cache) devre dışı bırakır, testi sıfırdan koşturur.
- `-run <Pattern>` : Sadece adı belirtilen regex deseni ile eşleşen testleri çalıştırır.
- `-bench <Pattern>` : Sadece belirtilen isimdeki benchmark fonksiyonlarını çalıştırır.
- `-run=^$` : Normal birim testlerini atlar, sadece benchmark'ların çalışmasını sağlar.
- `-benchmem` : Benchmark çıktısında bellek tahsisatı (`B/op` ve `allocs/op`) sütunlarını gösterir.
- `-benchtime=Nx` veya `-benchtime=Ts` : Benchmark'ın N defa veya T saniye boyunca tekrarlanmasını sağlar (örn: `-benchtime=500x` veya `-benchtime=2s`).
