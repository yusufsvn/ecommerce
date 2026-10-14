package httphandlers

import (
	e "ecommerceapi/CustomErrors"
	services "ecommerceapi/Services"
	st "ecommerceapi/Structs"
	"ecommerceapi/Utils"
	"encoding/json"
	"log"
	"net"
	"net/http"
)

type OrderHandler struct {
	orderService services.OrderService
}

func NewOrderHandler(_orderService services.OrderService) *OrderHandler {
	return &OrderHandler{
		orderService: _orderService,
	}
}

func (oh *OrderHandler) CreateOrder(w http.ResponseWriter, r *http.Request) {
	var req st.CreatePaymentDto
	Utils.ReadJSON(w, r, &req)

	ip, port, err_ip := net.SplitHostPort(r.RemoteAddr)
	log.Printf("ip: %s,port: %s", ip, port)
	if err_ip != nil {
		log.Printf("Error splitting host port: %v", err_ip)
		http.Error(w, "Error splitting host port", http.StatusInternalServerError)
		return
	}
	req.User_ip = ip

	err := oh.orderService.CreateOrder(r.Context(), req)

	if err != nil {
		log.Printf("Sipariş oluşturulurken hata oluştu %v", err)
		http.Error(w, "Sipariş oluşturulurken hata oluştu", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)

}

func (oh *OrderHandler) GetOrdersByUserId(w http.ResponseWriter, r *http.Request) {

	data, err := oh.orderService.GetOrdersByUserId(r.Context())

	if err != nil {
		log.Printf("order getirilemedi--OrderHandler--GetOrdersByUserId--%v", err)
		http.Error(w, "siparişler getirilemedi", http.StatusInternalServerError)
		return
	}

	json.NewEncoder(w).Encode(data)
	w.Header().Set("Content-Type", "application/json")

}

func (oh *OrderHandler) GetOrderDetails(w http.ResponseWriter, r *http.Request) {

	order_id := r.URL.Query().Get("orderId")

	if order_id == "" {
		http.Error(w, "order_id eksik", http.StatusBadRequest)
		return
	}

	data, err := oh.orderService.GetOrderDetails(r.Context(), order_id)

	if err != nil {
		log.Printf("order details error")
		http.Error(w, "beklenmeyen hatayla karşılaşıldı", http.StatusInternalServerError)
		return
	}

	json.NewEncoder(w).Encode(data)
	w.Header().Set("Content-Type", "application/json")

}

// only admin
func (oh *OrderHandler) GetAllOrders(w http.ResponseWriter, r *http.Request) {

	limit := r.URL.Query().Get("limit")
	offset := r.URL.Query().Get("offset")

	if limit == "" || offset == "" {
		http.Error(w, "query stringsler boş olamaz", http.StatusBadRequest)
		return
	}

	orders, err := oh.orderService.GetAllOrders(r.Context(), limit, offset)

	if err != nil {
		log.Printf("GetAllOrders da hata %v", err)
		http.Error(w, "Beklenmeyen Hatayla karşılaşıldı", http.StatusInternalServerError)
		return
	}

	json.NewEncoder(w).Encode(orders)
	w.Header().Set("Content-Type", "application/json")

}

func (oh *OrderHandler) UpdateOrderState(w http.ResponseWriter, r *http.Request) {

	var bodyDatas st.UpdateOrderStateReq
	Utils.ReadJSON(w, r, &bodyDatas)

	err := oh.orderService.UpdateOrderState(r.Context(), bodyDatas.NewState, bodyDatas.Order_id)

	if err != nil {
		http.Error(w, "Beklenmedik Hatayla karşılaşıldı", http.StatusInternalServerError)
		return
	}

}

func (oh *OrderHandler) UpdateMultiState(w http.ResponseWriter, r *http.Request) {

	var payload st.UpdateMultiOrderReq

	Utils.ReadJSON(w, r, &payload)

	err := oh.orderService.UpdateMultiState(r.Context(), &payload)

	if err != nil {
		http.Error(w, "Ürün State i değiştirilemedi", http.StatusInternalServerError)
		return
	}

}

func (oh *OrderHandler) GetOrdersByStatus(w http.ResponseWriter, r *http.Request) {

	limit := r.URL.Query().Get("limit")
	offset := r.URL.Query().Get("offset")
	status := r.URL.Query().Get("status")

	if limit == "" || offset == "" || status == "" {
		http.Error(w, "query stringsler boş olamaz", http.StatusBadRequest)
		return
	}

	orders, err := oh.orderService.GetOrdersByStatus(r.Context(), limit, offset, status)

	if err != nil {
		log.Printf("GetAllOrders da hata %v", err)
		http.Error(w, "Beklenmeyen Hatayla karşılaşıldı", http.StatusInternalServerError)
		return
	}

	json.NewEncoder(w).Encode(orders)
	w.Header().Set("Content-Type", "application/json")

}

func (oh *OrderHandler) GetOrderCountForWaiting(w http.ResponseWriter, r *http.Request) {
	count, err := oh.orderService.GetOrderCountForWaiting(r.Context())

	if err != nil {
		http.Error(w, "Sipariş sayısı Getirilirken hata oluştu", http.StatusInternalServerError)
		return
	}

	json.NewEncoder(w).Encode(count)
}

func (oh *OrderHandler) SearchOrder(w http.ResponseWriter, r *http.Request) {

	limit := r.URL.Query().Get("limit")
	search_query := r.URL.Query().Get("search_query")

	if limit == "" || search_query == "" {
		http.Error(w, "query stringsler boş olamaz", http.StatusBadRequest)
		return
	}

	orders, err := oh.orderService.SearchOrder(r.Context(), limit, search_query)

	if err == e.ZeroListResult {
		json.NewEncoder(w).Encode(orders)
		return
	}
	if err != nil {
		http.Error(w, "Beklenmedik Hatayla karşılaşıldı", http.StatusInternalServerError)
		return
	}

	json.NewEncoder(w).Encode(orders)
	w.Header().Set("Content-Type", "application/json")
}
