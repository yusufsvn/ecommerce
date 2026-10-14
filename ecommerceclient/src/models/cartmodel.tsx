import type { ProductResponse } from "./ProductResponse";


export interface cartmodel {
    product: ProductResponse
    amount: number
}

export interface cartProps {
    products: cartItem[]
    deleteProduct: (id: string) => void
    handleDecrement: (id: string) => void
    handleIncrement: (id: string) => void
    calculateTotalPrice: () => number
}


export interface cartResponse {
    cart_id: string
    products: cartProducts[]
}

export interface cartProducts {
    product_id: string
    quantity: number
}

export interface cartItem {
    product_id: string
    product_name: string
    price: number
    stock: number
    thumbnail_image: string
    quantity: number
}

export interface cartImage {
    image_id: string;
    path: string;
}

export interface UpdateCart {
    productId: string;
    quantity: number;
}

export interface AddCart {
    productId: string;
    quantity: number;
}

export interface DeleteCart {
    productId: string;
    cartId: string;
}