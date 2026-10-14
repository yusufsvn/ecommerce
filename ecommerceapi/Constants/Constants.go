package constants

const PhysicalProduct string = "PHYSICAL"
const VirtualProduct string = "VIRTUAL"

// Product publish states
const PublishStatePublished string = "PUBLISHED"
const PublishStateUnpublished string = "UNPUBLISHED"

// Product order states
const OrderStateCancelled string = "cancelled"
const OrderStateDelivered string = "delivered"
const OrderStatePreparing string = "preparing"
const OrderStateShipped string = "shipped"
const OrderStatePending string = "pending"

// Constants for product sorting
const Price string = "price"
const Name string = "name"
const Date string = "date"
const Desc string = "desc"
const Asc string = "asc"

var ValidSortColumns = map[string]string{
	"price": "p.price",
	"name":  "p.product_name",
	"date":  "p.created_at",
}
