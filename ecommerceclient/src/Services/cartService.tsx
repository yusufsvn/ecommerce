import { api } from "../Interceptors/httpinterceptor";
import type { cartItem, cartResponse } from "../models/cartmodel";
import { productService } from "./productService";
import type { CartProductResponse, UpdateCartItem } from "../models/ProductResponse";

const baseUrl = "https://ecommerceapi.yusufsevinc.work/api/user";


const isAuth = () => localStorage.getItem("isAuthenticated") === "true";
export const cartService = {

    updateCartItemAmount: async (updatePayload: UpdateCartItem): Promise<number> => {
        if (isAuth()) {
            const response = await api.put<UpdateCartItem>(
                `${baseUrl}/updatecartitem`,
                updatePayload,
                { withCredentials: true }
            );
            return response.status;
        }
        return -1;
    },
    addCartItem: async (productId: string, amount: number): Promise<number> => {
        if (isAuth()) {
            const response = await api.post(
                `${baseUrl}/addtocart`,
                { product_id: productId, quantity: amount },
                { withCredentials: true }
            );
            return response.status;
        }
        return -1;
    },
    deleteCartItem: async (productId: string): Promise<number> => {
        if (isAuth()) {
            const response = await api.delete(
                `${baseUrl}/removefromcart`,
                {
                    data: { product_id: productId },
                    withCredentials: true
                }
            );
            return response.status;
        }
        return -1;
    },
    getCartItems: async (): Promise<cartResponse> => {
        if (isAuth()) {
            const response = await api.get<cartResponse>(
                `${baseUrl}/getcart`,
                { withCredentials: true }
            );
            return response.data;
        }
        return { cart_id: "", products: [] };
    },

    getandProcessCartItems: async (): Promise<cartItem[]> => {
        const cartitems = await cartService.getCartItems().then(async (data: cartResponse) => {
            if (data) {
                const cartitems = await productService.getProductByIds(data.products.map((product) => product.product_id))
                    .then((products: CartProductResponse[]) => {
                        const cartItems: cartItem[] = [];

                        for (let i = 0; i < data.products.length; i++) {
                            cartItems.push({
                                product_id: data.products[i].product_id,
                                product_name: products[i].product_name,
                                price: products[i].price,
                                stock: products[i].stock,
                                thumbnail_image: products[i].thumbnail_image,
                                quantity: data.products[i].quantity
                            });
                        }
                        return cartItems;
                    });
                return cartitems;
            }
            return [];
        });
        return cartitems;
    }
};


