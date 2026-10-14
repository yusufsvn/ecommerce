package Utils

import (
	"math"
	"net"
	"net/http"
	"strings"
	"sync"
	"time"
)

// TokenBucket her bir istemci (IP) için jeton havuzunu temsil eder.
// Token Bucket Algoritması:
// 1. Kovaya saniyede 'refillRate' hızında yeni jetonlar damlar.
// 2. Kovanın alabileceği maksimum jeton miktarı 'capacity' kadardır (Burst limiti).
// 3. Her istek geldiğinde kovadan 1 jeton alınır. Kova boşsa istek 429 ile reddedilir.
type TokenBucket struct {
	capacity   float64   // Kovanın maksimum jeton kapasitesi (Burst)
	tokens     float64   // Kovadaki mevcut jeton sayısı
	refillRate float64   // Saniyede kovaya eklenen jeton miktarı (Jeton/saniye)
	lastRefill time.Time // Son jeton hesaplama zamanı
	mu         sync.Mutex
}

// NewTokenBucket belirtilen kapasite ve saniyelik dolum hızına göre yeni bir kova oluşturur.
func NewTokenBucket(capacity float64, refillRate float64) *TokenBucket {
	return &TokenBucket{
		capacity:   capacity,
		tokens:     capacity, // Başlangıçta kova tam doludur
		refillRate: refillRate,
		lastRefill: time.Now(),
	}
}

// Allow bir istek geldiğinde kovadan 1 jeton harcamaya çalışır.
// Jeton varsa true, kova boşsa false döner.
func (tb *TokenBucket) Allow() bool {
	tb.mu.Lock()
	defer tb.mu.Unlock()

	now := time.Now()
	elapsed := now.Sub(tb.lastRefill).Seconds()

	// Geçen süre boyunca kovaya damlayan yeni jetonları hesapla ve ekle
	tb.tokens = math.Min(tb.capacity, tb.tokens+(elapsed*tb.refillRate))
	tb.lastRefill = now

	// Kovada en az 1 jeton varsa harca ve izin ver
	if tb.tokens >= 1.0 {
		tb.tokens -= 1.0
		return true
	}

	return false
}

// clientEntry haritada saklanan istemci verisi ve son görülme zamanı
type clientEntry struct {
	bucket   *TokenBucket
	lastSeen time.Time
}

// IPRateLimiter IP adreslerine göre Token Bucket sınırlaması uygular.
type IPRateLimiter struct {
	mu         sync.RWMutex
	clients    map[string]*clientEntry
	capacity   float64
	refillRate float64
	clientTTL  time.Duration
	stopCh     chan struct{}
}

// NewIPRateLimiter yeni bir rate limiter başlatır ve arka planda temizleme döngüsünü çalıştırır.
// capacity: Anlık maksimum istek sayısı (Burst, örn: 5)
// refillRate: Saniyede kovaya eklenen jeton (örn: 5/60.0 -> dakikada 5 istek)
// ttl: İnaktif IP'lerin bellekten temizlenme süresi (Memory leak önlemi)
func NewIPRateLimiter(capacity float64, refillRate float64, ttl time.Duration) *IPRateLimiter {
	limiter := &IPRateLimiter{
		clients:    make(map[string]*clientEntry),
		capacity:   capacity,
		refillRate: refillRate,
		clientTTL:  ttl,
		stopCh:     make(chan struct{}),
	}

	go limiter.cleanupLoop()

	return limiter
}

// Stop arka plandaki temizleme goroutine'ini durdurur.
func (rl *IPRateLimiter) Stop() {
	select {
	case <-rl.stopCh:
		// Zaten durdurulmuş
	default:
		close(rl.stopCh)
	}
}

// GetLimiter belirtilen IP için mevcut TokenBucket'ı döner veya yenisini oluşturur.
func (rl *IPRateLimiter) GetLimiter(ip string) *TokenBucket {
	rl.mu.Lock()
	defer rl.mu.Unlock()

	client, exists := rl.clients[ip]
	if !exists {
		bucket := NewTokenBucket(rl.capacity, rl.refillRate)
		rl.clients[ip] = &clientEntry{
			bucket:   bucket,
			lastSeen: time.Now(),
		}
		return bucket
	}

	client.lastSeen = time.Now()
	return client.bucket
}

// cleanupLoop periyodik olarak inaktif IP'leri bellekten temizler.
func (rl *IPRateLimiter) cleanupLoop() {
	ticker := time.NewTicker(rl.clientTTL / 2)
	defer ticker.Stop()

	for {
		select {
		case <-ticker.C:
			rl.cleanup()
		case <-rl.stopCh:
			return
		}
	}
}

// cleanup süresi dolmuş IP kayıtlarını haritadan siler.
func (rl *IPRateLimiter) cleanup() {
	rl.mu.Lock()
	defer rl.mu.Unlock()

	now := time.Now()
	for ip, client := range rl.clients {
		if now.Sub(client.lastSeen) > rl.clientTTL {
			delete(rl.clients, ip)
		}
	}
}

// GetClientIP HTTP isteğinden gerçek istemci IP'sini ayıklar.
func GetClientIP(r *http.Request) string {
	xff := r.Header.Get("X-Forwarded-For")
	if xff != "" {
		ips := strings.Split(xff, ",")
		clientIP := strings.TrimSpace(ips[0])
		if clientIP != "" {
			return clientIP
		}
	}

	xRealIP := r.Header.Get("X-Real-IP")
	if xRealIP != "" {
		return strings.TrimSpace(xRealIP)
	}

	ip, _, err := net.SplitHostPort(r.RemoteAddr)
	if err != nil {
		return r.RemoteAddr
	}

	return ip
}

// RateLimitMiddleware belirtilen IPRateLimiter ile HTTP Handler'ı sarmalar.
func RateLimitMiddleware(limiter *IPRateLimiter) func(http.HandlerFunc) http.HandlerFunc {
	return func(next http.HandlerFunc) http.HandlerFunc {
		return func(w http.ResponseWriter, r *http.Request) {
			ip := GetClientIP(r)
			bucket := limiter.GetLimiter(ip)

			if !bucket.Allow() {
				w.Header().Set("Retry-After", "60")
				w.Header().Set("Content-Type", "application/json; charset=utf-8")
				w.WriteHeader(http.StatusTooManyRequests)
				_, _ = w.Write([]byte(`{"error":"Çok fazla istek yapıldı. Lütfen 1 dakika sonra tekrar deneyiniz."}`))
				return
			}

			next.ServeHTTP(w, r)
		}
	}
}
