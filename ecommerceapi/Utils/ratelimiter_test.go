package Utils

import (
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

func TestGetClientIP(t *testing.T) {
	t.Run("X-Forwarded-For başlığından ilk IP'yi almalı", func(t *testing.T) {
		req := httptest.NewRequest(http.MethodGet, "/", nil)
		req.Header.Set("X-Forwarded-For", "203.0.113.195, 70.41.3.18, 150.172.238.178")
		req.RemoteAddr = "10.0.0.1:1234"

		ip := GetClientIP(req)
		assert.Equal(t, "203.0.113.195", ip)
	})

	t.Run("X-Real-IP başlığından IP almalı", func(t *testing.T) {
		req := httptest.NewRequest(http.MethodGet, "/", nil)
		req.Header.Set("X-Real-IP", "198.51.100.1")
		req.RemoteAddr = "10.0.0.1:1234"

		ip := GetClientIP(req)
		assert.Equal(t, "198.51.100.1", ip)
	})

	t.Run("Başlık yoksa RemoteAddr adresindeki portu ayıklayarak IP'yi almalı", func(t *testing.T) {
		req := httptest.NewRequest(http.MethodGet, "/", nil)
		req.RemoteAddr = "192.168.1.50:54321"

		ip := GetClientIP(req)
		assert.Equal(t, "192.168.1.50", ip)
	})
}

func TestTokenBucket_BasicLogic(t *testing.T) {
	t.Run("Kapasite kadar istek izin verilmeli, kova boşalınca reddedilmeli", func(t *testing.T) {
		// Kapasite = 2, Saniyede 0.5 jeton
		bucket := NewTokenBucket(2, 0.5)

		assert.True(t, bucket.Allow(), "1. istek geçmeli")
		assert.True(t, bucket.Allow(), "2. istek geçmeli (kova doluydu)")
		assert.False(t, bucket.Allow(), "3. istek reddedilmeli (kova boşaldı)")
	})
}

func TestRateLimitMiddleware(t *testing.T) {
	t.Run("Burst limitine kadar istekler 200 dönmeli, limiti aşınca 429 dönmeli", func(t *testing.T) {
		// Saniyede 1 jeton eklenen, maksimum 3 jeton kapasiteli limiter
		limiter := NewIPRateLimiter(3, 1.0, 1*time.Minute)
		defer limiter.Stop()

		hitCount := 0
		dummyHandler := func(w http.ResponseWriter, r *http.Request) {
			hitCount++
			w.WriteHeader(http.StatusOK)
		}

		protectedHandler := RateLimitMiddleware(limiter)(dummyHandler)

		// İlk 3 istek başarıyla geçmeli (Burst = 3)
		for i := 1; i <= 3; i++ {
			req := httptest.NewRequest(http.MethodPost, "/login", nil)
			req.RemoteAddr = "1.2.3.4:1111"
			rr := httptest.NewRecorder()

			protectedHandler.ServeHTTP(rr, req)
			assert.Equal(t, http.StatusOK, rr.Code, "İstek %d 200 dönmeli", i)
		}
		assert.Equal(t, 3, hitCount)

		// 4. istek kova boş olduğu için 429 Too Many Requests almalı
		req := httptest.NewRequest(http.MethodPost, "/login", nil)
		req.RemoteAddr = "1.2.3.4:1111"
		rr := httptest.NewRecorder()

		protectedHandler.ServeHTTP(rr, req)
		assert.Equal(t, http.StatusTooManyRequests, rr.Code)
		assert.Equal(t, "60", rr.Header().Get("Retry-After"))
		assert.Equal(t, 3, hitCount, "Handler çalıştırılmamalı")
	})

	t.Run("Farklı IP adresleri birbirinin kotasını etkilememeli", func(t *testing.T) {
		limiter := NewIPRateLimiter(1, 1.0, 1*time.Minute)
		defer limiter.Stop()

		protectedHandler := RateLimitMiddleware(limiter)(func(w http.ResponseWriter, r *http.Request) {
			w.WriteHeader(http.StatusOK)
		})

		// IP A - 1. istek geçer
		reqA1 := httptest.NewRequest(http.MethodPost, "/login", nil)
		reqA1.RemoteAddr = "10.0.0.1:1111"
		rrA1 := httptest.NewRecorder()
		protectedHandler.ServeHTTP(rrA1, reqA1)
		assert.Equal(t, http.StatusOK, rrA1.Code)

		// IP A - 2. istek bloklanır (429)
		reqA2 := httptest.NewRequest(http.MethodPost, "/login", nil)
		reqA2.RemoteAddr = "10.0.0.1:1111"
		rrA2 := httptest.NewRecorder()
		protectedHandler.ServeHTTP(rrA2, reqA2)
		assert.Equal(t, http.StatusTooManyRequests, rrA2.Code)

		// IP B - 1. istek geçer (IP A'nın kotası IP B'yi etkilemez)
		reqB1 := httptest.NewRequest(http.MethodPost, "/login", nil)
		reqB1.RemoteAddr = "10.0.0.2:2222"
		rrB1 := httptest.NewRecorder()
		protectedHandler.ServeHTTP(rrB1, reqB1)
		assert.Equal(t, http.StatusOK, rrB1.Code)
	})

	t.Run("Zaman geçince kovaya yeni jeton damlamalı ve tekrar istek yapılabilmeli", func(t *testing.T) {
		// Saniyede 10 jeton (her 100ms'de 1 jeton), burst = 1
		limiter := NewIPRateLimiter(1, 10.0, 1*time.Minute)
		defer limiter.Stop()

		protectedHandler := RateLimitMiddleware(limiter)(func(w http.ResponseWriter, r *http.Request) {
			w.WriteHeader(http.StatusOK)
		})

		req := httptest.NewRequest(http.MethodPost, "/login", nil)
		req.RemoteAddr = "5.6.7.8:1234"

		// 1. istek geçer
		rr1 := httptest.NewRecorder()
		protectedHandler.ServeHTTP(rr1, req)
		assert.Equal(t, http.StatusOK, rr1.Code)

		// Hemen ardından 2. istek bloklanır
		rr2 := httptest.NewRecorder()
		protectedHandler.ServeHTTP(rr2, req)
		assert.Equal(t, http.StatusTooManyRequests, rr2.Code)

		// 120ms bekle (1 yeni jeton damlamalı)
		time.Sleep(120 * time.Millisecond)

		// 3. istek başarılı olmalı
		rr3 := httptest.NewRecorder()
		protectedHandler.ServeHTTP(rr3, req)
		assert.Equal(t, http.StatusOK, rr3.Code)
	})

	t.Run("Cleanup süresi dolmuş IP'leri bellekten temizlemeli", func(t *testing.T) {
		ttl := 50 * time.Millisecond
		limiter := NewIPRateLimiter(5, 1.0, ttl)
		defer limiter.Stop()

		// 1 IP kaydı oluştur
		_ = limiter.GetLimiter("1.1.1.1")

		limiter.mu.RLock()
		assert.Equal(t, 1, len(limiter.clients))
		limiter.mu.RUnlock()

		// TTL süresinden fazla bekle
		time.Sleep(80 * time.Millisecond)
		limiter.cleanup()

		limiter.mu.RLock()
		assert.Equal(t, 0, len(limiter.clients), "Süresi dolan IP bellekten silinmeli")
		limiter.mu.RUnlock()
	})
}
