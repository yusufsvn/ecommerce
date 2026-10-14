package httphandlers

import (
	services "ecommerceapi/Services"
	st "ecommerceapi/Structs"
	"ecommerceapi/Utils"
	"encoding/json"
	"log"
	"net/http"
)

type AddressHandler struct {
	AddressService services.AddressService
}

func NewAddressHandler(_addressService services.AddressService) *AddressHandler {
	return &AddressHandler{
		AddressService: _addressService,
	}
}

func (ah *AddressHandler) GetAddress(w http.ResponseWriter, r *http.Request) {

	data, err := ah.AddressService.GetAddress(r.Context())

	if err != nil {
		log.Printf("adres getirilemedi--GetAdress--UserHandler")
		http.Error(w, "adres verisi bulunamadı", http.StatusInternalServerError)
		return
	}

	json.NewEncoder(w).Encode(data)
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
}

func (ah *AddressHandler) UpdateAddress(w http.ResponseWriter, r *http.Request) {

	var addressReq st.Address
	Utils.ReadJSON(w, r, &addressReq)

	err := ah.AddressService.UpdateAddress(r.Context(), addressReq)

	if err != nil {
		log.Printf("address güncellenemedi--updateAdress--userhandler--%v", err)
		http.Error(w, "adres güncellenemedi", http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusOK)

}

func (ah *AddressHandler) CreateAddress(w http.ResponseWriter, r *http.Request) {

	// veri modeli
	var formData st.Address
	Utils.ReadJSON(w, r, &formData)

	err := ah.AddressService.CreateAddress(r.Context(), formData)

	if err != nil {
		http.Error(w, "Adres Eklenirken Hatayla Karşılaşıldı", http.StatusInternalServerError)
		return
	}

}

func (ah *AddressHandler) DeleteAddress(w http.ResponseWriter, r *http.Request) {

	id := r.URL.Query().Get("addressId")

	if id != "" {
		err := ah.AddressService.DeleteAddress(r.Context(), id)

		if err != nil {
			http.Error(w, "Beklenmedik bir hatayla karşılaşıldı", http.StatusInternalServerError)
			return
		}
		return
	}
	http.Error(w, "Adres silinemedi", http.StatusBadRequest)
}
