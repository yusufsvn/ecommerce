import type { UpdateProduct, Createproduct, FullProductresponse, patchProductModel, CartProductResponse, GetProductsResult, GetProductsWithCount, GetCategories, GetCategoryWithProduct } from "../models/ProductResponse";
import { api } from "../Interceptors/httpinterceptor";


const AdminUrl = "https://ecommerceapi.yusufsevinc.work/api/admin";
const publicUrl = "https://ecommerceapi.yusufsevinc.work/api/public";
const userUrl = "https://ecommerceapi.yusufsevinc.work/api/user";

export const productService = {


    createProduct: async (data: Createproduct) => {

        const formData = new FormData();

        // Append simple fields
        formData.append("ProductName", data.product_name);
        formData.append("Price", data.price.toString());
        formData.append("Description", data.description);
        if (data.product_category && data.product_category.length > 0) {
            data.product_category.forEach((cat) => {
                formData.append("Categories", cat);
            });
        }
        formData.append("ItemType", data.item_type)

        formData.append("Stock", data.stock.toString());
        if (data.technical_details && data.technical_details.length > 0) {
            const detailsMap: Record<string, string> = {};
            data.technical_details.forEach(feature => {
                detailsMap[feature.key] = feature.value;
            });
            formData.append("TechnicalDetails", JSON.stringify(detailsMap));
        }
        formData.append("IsChangedThumb", "true");
        formData.append("PublishState", data.publish_state.toString());
        // features form dataya eklenecek mi string olarak mı yoksa




        // Append each file. Assuming backend expects 'Images' or similar key for multiple files.
        // If backend expects 'Image', you might need to loop and append with same key.
        if (data.images && data.images.length > 0) {
            data.images.forEach((file) => {
                formData.append("Images", file);
            });
        }

        console.log(formData);

        const response = await api.post(AdminUrl + "/createproduct", formData, {
            headers: {
                'Content-Type': 'multipart/form-data'
            },
            withCredentials: true
        })
        return response.data;
    },

    // get products with category from admin panel
    getproducts: async (offset: number, limit: number): Promise<GetProductsResult> => {
        const response = await api.get<GetProductsResult>(AdminUrl + "/getallproducts", {
            params: {
                limit: limit,
                offset: offset
            }
        })
        return response.data;
    },

    // get products for public/guests
    getPublicProducts: async (offset: number, limit: number): Promise<GetProductsWithCount> => {
        const response = await api.get<GetProductsWithCount>(publicUrl + "/getallproducts", {
            params: {
                limit: limit,
                offset: offset
            }
        });
        return response.data;
    },

    // get filtered and sorted products
    getFilteredProducts: async (limit: number, offset: number, categoryIds: string[], sortBy: string): Promise<GetProductsWithCount> => {
        const response = await api.get<GetProductsWithCount>(publicUrl + "/getfilteredproducts", {
            params: {
                limit: limit,
                offset: offset,
                category: categoryIds && categoryIds.length > 0 ? categoryIds.join(',') : undefined,
                sortBy: sortBy
            }
        })
        return response.data;
    },
    // public
    getProductById: async (id: string) => {
        const response = await api.get<FullProductresponse>(publicUrl + "/getproductbyid", {
            params: {
                ProductID: id,
                product_id: id
            }
        });
        const data = response.data;
        if (data) {
            data.ProductID = data.product_id || data.ProductID;
            data.ProductName = data.product_name || data.ProductName;
            data.Stock = data.stock ?? data.Stock;
            data.Price = data.price ?? data.Price;
            data.Description = data.description || data.Description;
            data.Images = data.images || data.Images;
            data.Categories = data.categories || data.Categories;
            data.TechnicalDetails = data.technical_details || data.TechnicalDetails;
            data.ProductCode = data.product_code || data.ProductCode;
            data.ProductPublishState = data.product_publish_state || data.ProductPublishState;
            data.ItemType = data.item_type || data.ItemType;
        }
        return data;
    },

    getProductByIdForAdmin: async (id: string) => {
        const response = await api.get<FullProductresponse>(AdminUrl + "/getproductbyid", {
            params: {
                ProductID: id,
                product_id: id
            }
        });
        const data = response.data;
        if (data) {
            data.ProductID = data.product_id || data.ProductID;
            data.ProductName = data.product_name || data.ProductName;
            data.Stock = data.stock ?? data.Stock;
            data.Price = data.price ?? data.Price;
            data.Description = data.description || data.Description;
            data.Images = data.images || data.Images;
            data.Categories = data.categories || data.Categories;
            data.TechnicalDetails = data.technical_details || data.TechnicalDetails;
            data.ProductCode = data.product_code || data.ProductCode;
            data.ProductPublishState = data.product_publish_state || data.ProductPublishState;
            data.ItemType = data.item_type || data.ItemType;
        }
        return data;
    },

    deleteProduct: async (id: string) => {
        const response = await api.delete(AdminUrl + "/deleteproduct", {
            params: {
                ProductID: id
            },
            withCredentials: true
        });
        if (response.status === 200) {
            return true;
        }
        return false;
    },

    // Edit Product
    updateProduct: async (product: UpdateProduct) => {

        const response = await api.put(AdminUrl + "/updateproduct", product, {
            headers: {
                'Content-Type': 'application/json'
            },
            withCredentials: true
        })
        if (response.status === 200) {
            return true;
        }
        return false;
    },

    patchProduct: async (product: patchProductModel) => {
        const response = await api.patch(AdminUrl + "/patchproduct", product, { withCredentials: true })
        if (response.status === 200) {
            return true;
        }
        return false;
    },

    addProductImage: async (images: File[], productId: string, isChangedThumb: boolean) => {
        console.log("add image request");
        console.log(images);
        const formData = new FormData();

        console.log(isChangedThumb);
        formData.append("product_id", productId);
        formData.append("isChangedThumb", isChangedThumb.toString())
        // componentten bu metodu çağırdığı yerden parametre verilecek hem createproduct hem de edit.tsx de

        if (images && images.length > 0) {
            images.forEach((file) => {
                formData.append("Images", file);
            });
        }

        const response = await api.post(AdminUrl + "/addimage", formData, {
            headers: {
                'Content-Type': 'multipart/form-data'
            },
            withCredentials: true
        })

        if (response.status === 200) {
            return true;
        }
        return false;
    },

    getProductByIds: async (ids: string[]) => {
        const response = await api.get<CartProductResponse[]>(publicUrl + "/getproductsbyids", {
            params: {
                ids: ids.join(',')
            }
        })

        return response.data;
    },



    getCategories: async () => {
        const response = await api.get<GetCategories[]>(publicUrl + "/getcategories")
        return response.data;
    },


    searchProduct: async (searchQuery: string, limit: number) => {
        const response = await api.get<GetProductsWithCount>(userUrl + "/searchproduct", {
            params: {
                searchQuery: searchQuery,
                limit: limit
            }
        })
        return response.data;
    },

    searchAdminProduct: async (searchQuery: string, limit: number) => {
        const response = await api.get<GetCategoryWithProduct[]>(AdminUrl + "/products/search", {
            params: {
                searchQuery: searchQuery,
                limit: limit
            }
        })
        return response.data;
    },

}
