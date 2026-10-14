package httphandlers

import (
	services "ecommerceapi/Services"
	"encoding/json"
	"net/http"
)

type CategoryHandler struct {
	categoryService services.CategoryService
}

func NewCategoryHandler(categoryService services.CategoryService) *CategoryHandler {
	return &CategoryHandler{
		categoryService: categoryService,
	}
}

func (ch *CategoryHandler) AddNewCategory(w http.ResponseWriter, r *http.Request) {

	categoryName := r.URL.Query().Get("category_name")

	if categoryName == "" {
		http.Error(w, "category_name is required", http.StatusBadRequest)
		return
	}

	createddata, err := ch.categoryService.AddNewCategory(r.Context(), categoryName)

	if err != nil {
		http.Error(w, "Failed to add new category", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(createddata)

}

func (ch *CategoryHandler) DeleteCategory(w http.ResponseWriter, r *http.Request) {

	categoryID := r.URL.Query().Get("category_id")

	if categoryID == "" {
		http.Error(w, "category_id is required", http.StatusBadRequest)
		return
	}

	err := ch.categoryService.DeleteCategory(r.Context(), categoryID)

	if err != nil {
		http.Error(w, "Failed to delete category", http.StatusInternalServerError)
		return
	}
}

func (ch *CategoryHandler) UpdateCategory(w http.ResponseWriter, r *http.Request) {

	categoryID := r.URL.Query().Get("category_id")
	categoryName := r.URL.Query().Get("category_name")

	if categoryID == "" || categoryName == "" {
		http.Error(w, "category_id and category_name are required", http.StatusBadRequest)
		return
	}

	err := ch.categoryService.UpdateCategory(r.Context(), categoryID, categoryName)

	if err != nil {
		http.Error(w, "Failed to update category", http.StatusInternalServerError)
		return
	}
}

func (ch *CategoryHandler) GetAll(w http.ResponseWriter, r *http.Request) {

	data, err := ch.categoryService.GetAll(r.Context())

	if err != nil {
		http.Error(w, "beklenmedik hataayla karşılaşıldı", http.StatusInternalServerError)
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(data)
}
