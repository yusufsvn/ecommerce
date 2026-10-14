export interface CreateOrderRequest {
    Tckn: string;
    PhoneNumber: string;
    CardData: creditCardModel;
    BillingAddress: addressModel;
    ShippingAddress: addressModel;
    SecurePay: boolean;
    SaveUserAddress: boolean;
}

export interface addressModel {
    AddressStr: string;
    City: string;
    Country: string;
    ZipCode: string;
}

export interface creditCardModel {
    CardName: string;
    CardNumber: string;
    ExpireYear: string
    ExpireMonth: string
    Cvc: string;
}

export interface OrderResponse {
    order_id: string;
    user_id: string;
    created_at: string;
    status: string;
    total_amount: string;
    total_quantity: number
    order_code: string;
}

export interface OrderItem {
    products: OrderProducts[]
    shipping_address: string;
    billing_address: string;
}

export interface OrderProducts {
    product_name: string;
    product_id: string
    price_at_purchase: string;
    quantity: number;
    thumbnail_image: string;
}

export interface OrderForAdmin {
    order_id: string;
    created_at: string;
    status: string;
    total_amount: string;
    name: string
    surname: string
}

export interface OrderResult {
    count: number;
    data: OrderForAdmin[];
}