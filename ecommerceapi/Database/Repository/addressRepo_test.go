package repository

import (
	"context"
	st "ecommerceapi/Structs"
	"os"
	"testing"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/stretchr/testify/assert"
)

var testDbInstance *pgxpool.Pool

// Otomatik algınanır ve ilk olarak çalıştırılır.
func TestMain(m *testing.M) {

	testDB := SetupTestDatabase()
	testDbInstance = testDB.DbInstance

	defer testDB.TearDown()
	os.Exit(m.Run())
}

// data model: CreateAddressDTO
func TestCreateAddress(t *testing.T) {

	addressRepo := NewAddressStore(testDbInstance)

	t.Run("Başarılı Adres Oluşturma", func(t *testing.T) {

		// arrange
		id := uuid.NewString()
		t.Logf("Test user id: %s", id)

		address := &st.CreateAddressDTO{
			UserId:      id,
			Name:        "Ev Adresi",
			PostCode:    "12345",
			Country:     "Turkey",
			District:    "Kadikoy",
			City:        "Istanbul",
			FullAddress: "açık adres tarifi",
		}

		ctx := context.Background()
		createTestUser(t, ctx, addressRepo.db, id)
		// act
		err := addressRepo.Create(ctx, address)
		if err != nil {
			t.Fatalf("Failed to create address: %v", err)
		}

		// assert
		var validationAddress struct {
			UserId      string `db:"user_id" json:"user_id"`
			AddressId   string `db:"address_id" json:"address_id"`
			Name        string `db:"name" json:"name"`
			PostCode    string `db:"post_code" json:"post_code"`
			Country     string `db:"country" json:"country"`
			District    string `db:"district" json:"district"`
			City        string `db:"city" json:"city"`
			FullAddress string `db:"full_address" json:"full_address"`
		}

		query := "SELECT user_id, address_id, name, post_code, country, district, city, full_address FROM Addresses WHERE user_id=$1"

		testDbInstance.QueryRow(ctx, query, id).Scan(
			&validationAddress.UserId,
			&validationAddress.AddressId,
			&validationAddress.Name,
			&validationAddress.PostCode,
			&validationAddress.Country,
			&validationAddress.District,
			&validationAddress.City,
			&validationAddress.FullAddress,
		)

		assert.Equal(t, address.UserId, validationAddress.UserId)
		assert.Equal(t, address.Name, validationAddress.Name)
		assert.Equal(t, address.PostCode, validationAddress.PostCode)
		assert.Equal(t, address.Country, validationAddress.Country)
		assert.Equal(t, address.District, validationAddress.District)
		assert.Equal(t, address.City, validationAddress.City)
		assert.Equal(t, address.FullAddress, validationAddress.FullAddress)

	})

}

func TestDeleteAddress(t *testing.T) {
	addressRepo := NewAddressStore(testDbInstance)

	t.Run("Adres Silme Testi", func(t *testing.T) {

		// arrange
		ctx := context.Background()
		id := uuid.NewString()

		// 1. user ekle
		createTestUser(t, ctx, addressRepo.db, id)
		// 2. address ekle
		address := &st.CreateAddressDTO{
			UserId:      id,
			Name:        "Ev Adresi",
			PostCode:    "12345",
			Country:     "Turkey",
			District:    "Kadikoy",
			City:        "Istanbul",
			FullAddress: "açık adres tarifi",
		}
		addressId := createTestAddress(t, ctx, addressRepo.db, address)
		// act
		err := addressRepo.Delete(ctx, addressId)
		if err != nil {
			t.Fatalf("Test Başarısız %v", err)
		}
		// assert

		var exists bool
		checkQuery := `SELECT EXISTS(SELECT 1 FROM addresses WHERE address_id = $1)`
		err = addressRepo.db.QueryRow(ctx, checkQuery, addressId).Scan(&exists)
		if err != nil {
			t.Fatalf("Veritabanı kontrol sorgusu başarısız oldu: %v", err)
		}

		if exists {
			t.Errorf("Adres silinmedi: ID %s olan kayıt veritabanında hala mevcut", addressId)
		}

	})
}
func TestUpdateAddress(t *testing.T) {
	addressRepo := NewAddressStore(testDbInstance)

	t.Run("Adres Güncelleme Testi", func(t *testing.T) {
		ctx := context.Background()
		userID := uuid.NewString()
		createTestUser(t, ctx, addressRepo.db, userID)

		address := &st.CreateAddressDTO{
			UserId:      userID,
			Name:        "Ev Adresi",
			PostCode:    "12345",
			Country:     "Turkey",
			District:    "Kadikoy",
			City:        "Istanbul",
			FullAddress: "Eski adres",
		}
		addressId := createTestAddress(t, ctx, addressRepo.db, address)

		updatedAddress := st.Address{
			AddressStr: "Yeni adres bilgisi",
			ZipCode:    "54321",
			City:       "Ankara",
			District:   "Cankaya",
			Country:    "Turkey",
		}

		err := addressRepo.Update(ctx, updatedAddress, userID)
		if err != nil {
			t.Fatalf("Adres güncelleme başarısız oldu: %v", err)
		}

		var validationAddress struct {
			AddressID   string `db:"address_id"`
			UserID      string `db:"user_id"`
			PostCode    string `db:"post_code"`
			Country     string `db:"country"`
			District    string `db:"district"`
			City        string `db:"city"`
			FullAddress string `db:"full_address"`
		}

		query := `SELECT address_id, user_id, post_code, country, district, city, full_address FROM addresses WHERE address_id=$1`
		err = addressRepo.db.QueryRow(ctx, query, addressId).Scan(
			&validationAddress.AddressID,
			&validationAddress.UserID,
			&validationAddress.PostCode,
			&validationAddress.Country,
			&validationAddress.District,
			&validationAddress.City,
			&validationAddress.FullAddress,
		)
		if err != nil {
			t.Fatalf("Güncellenen adres verisi alınamadı: %v", err)
		}

		assert.Equal(t, userID, validationAddress.UserID)
		assert.Equal(t, updatedAddress.ZipCode, validationAddress.PostCode)
		assert.Equal(t, updatedAddress.Country, validationAddress.Country)
		assert.Equal(t, updatedAddress.District, validationAddress.District)
		assert.Equal(t, updatedAddress.City, validationAddress.City)
		assert.Equal(t, updatedAddress.AddressStr, validationAddress.FullAddress)
		assert.Equal(t, addressId, validationAddress.AddressID)
	})
}

func TestGetByUserId(t *testing.T) {
	addressRepo := NewAddressStore(testDbInstance)

	t.Run("Kullanıcıya Göre Adresleri Getirir", func(t *testing.T) {
		ctx := context.Background()
		userID := uuid.NewString()
		createTestUser(t, ctx, addressRepo.db, userID)

		addressOne := &st.CreateAddressDTO{
			UserId:      userID,
			Name:        "Ev Adresi",
			PostCode:    "12345",
			Country:     "Turkey",
			District:    "Kadikoy",
			City:        "Istanbul",
			FullAddress: "Birinci adres",
		}
		addressTwo := &st.CreateAddressDTO{
			UserId:      userID,
			Name:        "İş Adresi",
			PostCode:    "67890",
			Country:     "Turkey",
			District:    "Besiktas",
			City:        "Istanbul",
			FullAddress: "İkinci adres",
		}

		createTestAddress(t, ctx, addressRepo.db, addressOne)
		createTestAddress(t, ctx, addressRepo.db, addressTwo)
		// act
		addresses, err := addressRepo.GetByUserId(ctx, userID)
		if err != nil {
			t.Fatalf("Adresleri getirirken hata oluştu: %v", err)
		}

		// assert
		if addresses == nil {
			t.Fatal("Adres listesi nil döndü")
		}

		assert.Len(t, *addresses, 2)

		assert.Contains(t, []string{(*addresses)[0].Name, (*addresses)[1].Name}, addressOne.Name)
		assert.Contains(t, []string{(*addresses)[0].Name, (*addresses)[1].Name}, addressTwo.Name)
		assert.Contains(t, []string{(*addresses)[0].PostCode, (*addresses)[1].PostCode}, addressOne.PostCode)
		assert.Contains(t, []string{(*addresses)[0].PostCode, (*addresses)[1].PostCode}, addressTwo.PostCode)
	})
}

// helper methods

func createTestAddress(t *testing.T, ctx context.Context, db *pgxpool.Pool, dto *st.CreateAddressDTO) string {
	t.Helper() // Test loglarında hatanın bu fonksiyonda değil, çağrıldığı satırda görünmesini sağlar.

	var addressId string

	query := `
		INSERT INTO addresses (user_id, name, post_code, country, district, city, full_address)
		VALUES ($1, $2, $3, $4, $5, $6, $7) returning address_id`

	err := db.QueryRow(ctx,
		query,
		dto.UserId,
		dto.Name,
		dto.PostCode,
		dto.Country,
		dto.District,
		dto.City,
		dto.FullAddress,
	).Scan(&addressId)

	if err != nil {
		t.Fatalf("Test için adres verisi oluşturulamadı: %v", err)
	}

	return addressId

}
