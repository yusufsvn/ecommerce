package services

import (
	"context"
	repo "ecommerceapi/Database/Repository"
	structs "ecommerceapi/Structs"
	"log"
)

type CategoryService struct {
	PCategoryStore repo.ProductCategoryFunctions
}

func NewCategoryService(_PCategoryStore repo.ProductCategoryFunctions) *CategoryService {
	return &CategoryService{
		PCategoryStore: _PCategoryStore,
	}
}

func (cts *CategoryService) AddNewCategory(ctx context.Context, category_name string) (structs.GetCategory, error) {
	createddata, err := cts.PCategoryStore.Create(ctx, category_name)

	if err != nil {
		log.Printf("AddNewCategory: cannot create new category: %v", err)
		return structs.GetCategory{}, err
	}

	return createddata, nil

}

func (cts *CategoryService) DeleteCategory(ctx context.Context, category_id string) error {

	err := cts.PCategoryStore.Delete(ctx, category_id)

	if err != nil {
		return err
	}

	return nil

}

func (cts *CategoryService) UpdateCategory(ctx context.Context, category_id string, category_name string) error {

	err := cts.PCategoryStore.Update(ctx, category_id, category_name)

	if err != nil {
		return err
	}

	return nil
}

func (cts *CategoryService) GetAll(ctx context.Context) ([]structs.GetCategory, error) {
	data, err := cts.PCategoryStore.GetAll(ctx)

	if err != nil {
		return nil, err
	}
	return data, nil
}
