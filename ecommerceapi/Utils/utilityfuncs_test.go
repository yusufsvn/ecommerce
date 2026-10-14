package Utils

import (
	"strings"
	"testing"
)

// IsValidState() Testleri (1)
func TestIsValidState(t *testing.T) {

	test := []struct {
		name  string
		input string
		want  bool
	}{
		{"delivered String", "delivered", true},
		{"cancelled String", "cancelled", true},
		{"shipped String", "shipped", true},
		{"pending String", "pending", true},
		{"preparing String", "preparing", true},
		{"invalid string", "unknown_state", false},
		{"empty string", "", false},
	}

	for _, tt := range test {
		t.Run(tt.name, func(t *testing.T) {
			got := IsValidState(tt.input)
			if got != tt.want {
				t.Errorf("IsValidState(%q) = %v, want %v", tt.input, got, tt.want)
				return
			}
		})
	}

}

// CalcTotalPrice - GetUserIdFromClaims - MatchWithURL - MatchWithURLSingle

// CreateNanoId() Testleri (2)
// Uzunluk - Uygun Karakter - Doğru Hata Kontrol Testi
func TestCreateNanoId1(t *testing.T) {

	test := []struct {
		name             string
		Length           int
		expectedCharSet  string
		wanterror        bool
		expectederrorstr string
	}{
		{"used size", 8, "ABCDEFGHIJKLMNPQRSTUVWXYZ123456789", false, ""},
		{"1 Char Code", 1, "ABCDEFGHIJKLMNPQRSTUVWXYZ123456789", false, ""},
		{"Zero Char Code", 0, "ABCDEFGHIJKLMNPQRSTUVWXYZ123456789", true, "size must be positive"},
		{"Invalid Size", -1, "ABCDEFGHIJKLMNPQRSTUVWXYZ123456789", true, "size must be positive"},
		{"size 10", 10, "ABCDEFGHIJKLMNPQRSTUVWXYZ123456789", false, ""},
		{"size 11", 11, "ABCDEFGHIJKLMNPQRSTUVWXYZ123456789", false, ""},
		{"long char code", 50, "ABCDEFGHIJKLMNPQRSTUVWXYZ123456789", false, ""},
	}

	for _, tt := range test {
		t.Run(tt.name, func(t *testing.T) {

			code, err := CreateNanoId(tt.Length)
			codelength := len(code)

			if tt.wanterror {
				if err == nil {
					t.Fatalf("hata bekleniyordu ama gelmedi")
				}
				if !strings.Contains(err.Error(), tt.expectederrorstr) {
					t.Fatalf("hata mesajı uyuşmadı, gelen: %v, istenen: %s", err, tt.expectederrorstr)
				}
				return // beklenen hata geldi, test başarılı
			}

			if err != nil {
				t.Fatalf("beklenmeyen hata: %v", err)
			}

			if codelength != tt.Length {
				t.Errorf("Verilen Uzunlukla Uyuşmayan kod üretildi istenen:%d üretilen:%d", tt.Length, codelength)
				return
			}

			for _, char := range code {
				if !strings.ContainsRune(tt.expectedCharSet, char) {
					t.Errorf("Verilen Alfabeyle uyuşmayan karakter içeriyor")
					return
				}

			}

		})

	}
}

// Unique testi
func TestCreateNanoId_Uniqueness(t *testing.T) {
	tests := []struct {
		name       string
		length     int
		iterations int
	}{
		{"short code", 8, 10_000},
		{"medium code", 8, 50_000},
		{"long code", 8, 100_000},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			seen := make(map[string]struct{}, tt.iterations)

			for i := 0; i < tt.iterations; i++ {
				code, err := CreateNanoId(tt.length)
				if err != nil {
					t.Fatalf("iterasyon %d: beklenmeyen hata: %v", i, err)
				}

				if _, exists := seen[code]; exists {
					t.Fatalf("iterasyon %d: benzersizlik ihlali, tekrar eden kod: %q", i, code)
				}

				seen[code] = struct{}{}
			}
		})
	}
}
