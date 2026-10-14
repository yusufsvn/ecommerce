import type { cartItem } from "./cartmodel";


export interface ProductListProps {
    addtocart: (product: cartItem) => void
    cartItems?: cartItem[];
}

export interface NavbarProps {
    getTotalItems: () => number;
}