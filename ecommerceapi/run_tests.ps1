param (
    [Parameter(Position = 0)]
    [ValidateSet("all", "unit", "integration", "bench", "coverage", "clean")]
    [string]$Type = "all"
)

Write-Host "==========================================" -ForegroundColor Cyan
Write-Host "  E-Commerce API Test Runner ($Type)" -ForegroundColor Cyan
Write-Host "==========================================" -ForegroundColor Cyan

switch ($Type) {
    "all" {
        Write-Host "`n[1/2] Birim ve Entegrasyon Testleri Koşuluyor..." -ForegroundColor Yellow
        go test -v -count=1 ./...
    }
    "unit" {
        Write-Host "`nBirim (Unit) Testleri Koşuluyor..." -ForegroundColor Yellow
        go test -v ./Database/Repository/... ./Services/... ./HttpHandlers/... ./Utils/...
    }
    "integration" {
        Write-Host "`nEntegrasyon (Integration) Testleri Koşuluyor (Docker PostgreSQL)..." -ForegroundColor Yellow
        go test -v -run TestIntegration .
    }
    "bench" {
        Write-Host "`nYük ve Performans (Benchmark) Testleri Koşuluyor..." -ForegroundColor Yellow
        go test -run=^$ -bench=. -benchmem .
    }
    "coverage" {
        Write-Host "`nTest Kapsamı (Coverage) Hesaplanıyor..." -ForegroundColor Yellow
        go test -coverprofile=coverage.out ./...
        Write-Host "Rapor oluşturuldu: coverage.out" -ForegroundColor Green
        Write-Host "Tarayıcıda açmak için: go tool cover -html=coverage.out" -ForegroundColor Green
    }
    "clean" {
        Write-Host "`nTest önbelleği temizleniyor..." -ForegroundColor Yellow
        go clean -testcache
        Write-Host "Önbellek temizlendi." -ForegroundColor Green
    }
}

if ($LASTEXITCODE -eq 0) {
    Write-Host "`n[BAŞARILI] Tüm testler başarıyla tamamlandı!" -ForegroundColor Green
} else {
    Write-Host "`n[HATA] Bazı testler başarısız oldu (Hata Kodu: $LASTEXITCODE)." -ForegroundColor Red
}
