package main

import (
	structs "ecommerceapi/Structs"
	"encoding/json"
	"fmt"
	"net/http"
	"testing"

	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
)

func TestIntegration_Address_Lifecycle(t *testing.T) {
	app := setupIntegrationApp(t)

	userEmail := fmt.Sprintf("addr-user-%s@example.com", uuid.NewString()[:8])
	_, accessCookie, _ := app.createTestUser(t, userEmail, "Pass123!", "user")

	addrName := fmt.Sprintf("Ev-%s", uuid.NewString()[:4])
	var createdAddrID string

	t.Run("kullanıcı yeni adres oluşturabilmeli", func(t *testing.T) {
		createPayload := structs.Address{
			AddressName: addrName,
			AddressStr:  "Atatürk Cad. No: 42 D: 5",
			City:        "Ankara",
			District:    "Çankaya",
			Country:     "Türkiye",
			ZipCode:     "06690",
		}

		resp := app.doRequest(t, http.MethodPost, "/api/user/createaddress", createPayload, accessCookie)
		assert.Equal(t, http.StatusOK, resp.StatusCode)
	})

	t.Run("kullanıcı adreslerini listeleyebilmeli", func(t *testing.T) {
		resp := app.doRequest(t, http.MethodGet, "/api/user/getaddresses", nil, accessCookie)
		assert.Equal(t, http.StatusOK, resp.StatusCode)

		var addrList structs.GetAddressList
		err := json.NewDecoder(resp.Body).Decode(&addrList)
		assert.NoError(t, err)
		assert.NotNil(t, addrList.Addresses)
		assert.NotEmpty(t, *addrList.Addresses)

		var found bool
		for _, a := range *addrList.Addresses {
			if a.Name == addrName {
				found = true
				createdAddrID = a.AddressId
				assert.Equal(t, "Ankara", a.City)
				assert.Equal(t, "Çankaya", a.District)
				break
			}
		}
		assert.True(t, found, "Oluşturulan adres listede yer almalı")
	})

	t.Run("kullanıcı adresini güncelleyebilmeli", func(t *testing.T) {
		updatePayload := structs.Address{
			AddressName: addrName,
			AddressStr:  "İnönü Bulvarı No: 100",
			City:        "Ankara",
			District:    "Çankaya",
			Country:     "Türkiye",
			ZipCode:     "06690",
		}

		resp := app.doRequest(t, http.MethodPut, "/api/user/updateaddress", updatePayload, accessCookie)
		assert.Equal(t, http.StatusOK, resp.StatusCode)
	})

	t.Run("kullanıcı adresini silebilmeli", func(t *testing.T) {
		assert.NotEmpty(t, createdAddrID)
		path := fmt.Sprintf("/api/user/deleteaddress?addressId=%s", createdAddrID)
		resp := app.doRequest(t, http.MethodDelete, path, nil, accessCookie)
		assert.Equal(t, http.StatusOK, resp.StatusCode)

		// Listede artık olmamalı
		getResp := app.doRequest(t, http.MethodGet, "/api/user/getaddresses", nil, accessCookie)
		var addrList structs.GetAddressList
		_ = json.NewDecoder(getResp.Body).Decode(&addrList)
		if addrList.Addresses != nil {
			for _, a := range *addrList.Addresses {
				assert.NotEqual(t, createdAddrID, a.AddressId, "Silinen adres artık listede olmamalı")
			}
		}
	})
}
