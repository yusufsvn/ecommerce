package repository

import (
	"context"
	"testing"

	models "ecommerceapi/Database/Models"
	st "ecommerceapi/Structs"
	"ecommerceapi/Utils"

	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestProductOrderAdd(t *testing.T) {
	productOrderRepo := NewProductOrderStore(testDbInstance)

	t.Run("siparişe ürün ekler", func(t *testing.T) {
		ctx := context.Background()
		productID := createTestProduct(t, ctx, testDbInstance)
		userID := uuid.NewString()
		createTestUser(t, ctx, testDbInstance, userID)

		orderID := uuid.NewString()
		orderCode, err := Utils.CreateNanoId(8)
		require.NoError(t, err)
		require.NoError(t, NewOrderStore(testDbInstance).Create(ctx, &st.OrderDTO{
			Order_id:        orderID,
			User_id:         userID,
			Status:          "pending",
			TotalAmount:     "1000",
			ShippingAddress: "shipping",
			BillingAddress:  "billing",
			OrderCode:       orderCode,
		}))

		item := models.ProductOrder{
			ProductID:       productID,
			OrderID:         orderID,
			Quantity:        2,
			PriceAtPurchase: 1500,
		}

		err = productOrderRepo.Add(ctx, item)
		require.NoError(t, err)

		var storedProductID, storedOrderID string
		var storedQuantity int
		var storedPrice int64
		err = testDbInstance.QueryRow(
			ctx,
			"SELECT product_id, order_id, quantity, price_at_purchase FROM productorders WHERE product_id=$1 AND order_id=$2",
			productID,
			orderID,
		).Scan(&storedProductID, &storedOrderID, &storedQuantity, &storedPrice)
		require.NoError(t, err)

		assert.Equal(t, item.ProductID, storedProductID)
		assert.Equal(t, item.OrderID, storedOrderID)
		assert.Equal(t, item.Quantity, storedQuantity)
		assert.Equal(t, item.PriceAtPurchase, storedPrice)
	})
}

func TestProductOrderAddMany(t *testing.T) {
	productOrderRepo := NewProductOrderStore(testDbInstance)
	orderRepo := NewOrderStore(testDbInstance)

	t.Run("birden fazla ürünü toplu ekler", func(t *testing.T) {
		ctx := context.Background()
		productIDOne := createTestProduct(t, ctx, testDbInstance)
		productIDTwo := createTestProduct(t, ctx, testDbInstance)
		userID := uuid.NewString()
		createTestUser(t, ctx, testDbInstance, userID)

		orderID := uuid.NewString()
		orderCode, err := Utils.CreateNanoId(8)
		require.NoError(t, err)
		// alt test aşaması
		require.NoError(t, orderRepo.Create(ctx, &st.OrderDTO{
			Order_id:        orderID,
			User_id:         userID,
			Status:          "pending",
			TotalAmount:     "2500",
			ShippingAddress: "shipping",
			BillingAddress:  "billing",
			OrderCode:       orderCode,
		}))

		items := []st.ProductsFromCartDTO{
			{ProductID: productIDOne, Quantity: 1, Price: 1000},
			{ProductID: productIDTwo, Quantity: 2, Price: 1500},
		}

		err = productOrderRepo.AddMany(ctx, items, orderID)
		require.NoError(t, err)

		var count int
		err = testDbInstance.QueryRow(ctx, "SELECT COUNT(*) FROM productorders WHERE order_id=$1", orderID).Scan(&count)
		require.NoError(t, err)
		assert.Equal(t, 2, count)
	})
}

func TestProductOrderGetByOrderID(t *testing.T) {
	productOrderRepo := NewProductOrderStore(testDbInstance)

	t.Run("siparişe ait ürün detaylarını döndürür", func(t *testing.T) {
		ctx := context.Background()
		productID := createTestProduct(t, ctx, testDbInstance)
		userID := uuid.NewString()
		createTestUser(t, ctx, testDbInstance, userID)

		orderID := uuid.NewString()
		orderCode, err := Utils.CreateNanoId(8)
		require.NoError(t, err)
		require.NoError(t, NewOrderStore(testDbInstance).Create(ctx, &st.OrderDTO{
			Order_id:        orderID,
			User_id:         userID,
			Status:          "pending",
			TotalAmount:     "1000",
			ShippingAddress: "shipping",
			BillingAddress:  "billing",
			OrderCode:       orderCode,
		}))

		item := models.ProductOrder{
			ProductID:       productID,
			OrderID:         orderID,
			Quantity:        3,
			PriceAtPurchase: 900,
		}
		require.NoError(t, productOrderRepo.Add(ctx, item))

		items, err := productOrderRepo.GetByOrderID(ctx, orderID)
		require.NoError(t, err)
		require.Len(t, items, 1)

		assert.Equal(t, productID, items[0].ProductID)
		assert.Equal(t, item.Quantity, items[0].Quantity)
		assert.Equal(t, item.PriceAtPurchase, items[0].PriceAtPurchase)
	})
}
