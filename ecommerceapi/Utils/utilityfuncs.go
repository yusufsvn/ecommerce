package Utils

import (
	"context"
	c "ecommerceapi/Constants"
	st "ecommerceapi/Structs"
	structs "ecommerceapi/Structs"
	"ecommerceapi/Types"
	"fmt"
	"log"

	"github.com/golang-jwt/jwt/v5"
	gonanoid "github.com/matoous/go-nanoid"
)

// Verilen Sepet ürünleri ile toplam tutar hesaplar
func CalcTotalPrice(items []st.ProductsFromCartDTO) int64 {

	var totalPrice int64 = 0

	for i := 0; i < len(items); i++ {
		tempprice := items[i].Price * int64(items[i].Quantity)
		totalPrice += tempprice
	}

	return totalPrice
}

// Context Üzerinden Claims a ulaşır ve tokendaki user_id gönderir
func GetUserIdFromClaims(ctx context.Context) (string, error) {

	claims, ok := ctx.Value(Types.UserContextKey).(jwt.MapClaims)
	if !ok {
		log.Printf("invalid token claims")
		return "", fmt.Errorf("invalid token claims")
	}

	user_id, ok := claims["user_id"].(string)
	if !ok || user_id == "" {
		log.Printf("invalid or missing user id claim")
		return "", fmt.Errorf("invalid token claims")
	}

	return user_id, nil
}

func IsValidState(state string) bool {
	switch state {
	case c.OrderStateCancelled, c.OrderStateDelivered, c.OrderStatePending, c.OrderStatePreparing, c.OrderStateShipped:
		return true
	}
	return false
}

// size must be positive
func CreateNanoId(size int) (string, error) {

	var alphabet string = "ABCDEFGHIJKLMNPQRSTUVWXYZ123456789"

	nanoid, err := gonanoid.Generate(alphabet, size)

	if err != nil {
		log.Printf("nano id üretilemedi:%v", err)
		return "", err
	}

	return nanoid, nil

}

// Storageden erişilebilmesi için gerekli url formatını düzenler
func MatchWithURL(url string, bucketname string, images []structs.ProductImage) []structs.ProductImage {
	for i := range images {
		path := images[i].Path
		images[i].Path = fmt.Sprintf("https://%s/%s/%s", url, bucketname, path)
	}
	return images
}

func MatchWithURLSingle(url string, bucketname string, image string) string {
	return fmt.Sprintf("https://%s/%s/%s", url, bucketname, image)
}
