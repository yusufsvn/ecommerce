import type { ProductFeature } from '../components/admin/CreateProductFeature';
import type { cartItem } from './cartmodel';


export interface ProductResponse {
    ProductID: string,
    ProductName: string,
    Stock: number,
    Description: string,
    Price: number,
    Images: Image[],
    Categories: CategoryResponse[],
    product_publish_state?: string,
    ProductPublishState?: string,
    item_type?: string,
    ItemType?: string
}

export interface CategoryResponse {
    category_name: string
    category_id: string
}

export interface GetProducts {
    product_id: string,
    product_name: string,
    stock: number,
    price: number,
    thumbnail_image: string
}

export interface GetProductsResult {
    count: number;
    data: GetCategoryWithProduct[];
}

export interface GetProductsWithCount {
    Count: number;
    Data: GetProducts[];
}

export interface GetCategories {
    category_id: string;
    category_name: string;
}

export interface GetCategoryWithProduct {
    Product: GetProducts,
    Categories: GetCategories[]
}

export interface ProductResult {
    count: number;
    data: ProductResponse[];
}


export interface Image {
    image_id: string,
    path: string
}

export interface UpdateProduct {
    ProductID: string,
    ProductName: string,
    Stock: number,
    Description: string,
    Price: number,
    DeletedImages: string[],
    Categories: {
        NewCategories: string[];
        DeletedCategories: string[];
    },
    TechnicalDetails?: Record<string, string>;
    ProductPublishState?: string;
    ItemType?: string;
}

export interface GetCategory {
    category_id: string;
    category_name: string;
}

export interface ProductImage {
    image_id: string;
    path: string;
}

export interface FullProductResponse {
    product_id: string;
    product_name: string;
    stock: number;
    description: string;
    price: number;
    technical_details?: Record<string, string>;
    categories: GetCategory[];
    images: ProductImage[];
    product_code?: string;
    product_publish_state?: string;
    item_type?: string;

    // Optional backward compatibility properties
    ProductID?: string;
    ProductName?: string;
    Stock?: number;
    Description?: string;
    Price?: number;
    TechnicalDetails?: Record<string, string>;
    Categories?: GetCategory[];
    Images?: ProductImage[];
    ProductCode?: string;
    ProductPublishState?: string;
    ItemType?: string;
    ProductFeaturesUrl?: string;
}

export type FullProductresponse = FullProductResponse;

export interface patchProductModel {
    ProductID: string;
    ProductName: string;
    Stock: number;
    Price: number;
}

export interface EditProductModel {
    ProductID: string;
    ProductName: string;
    Stock: number;
    Price: number;
    Description: string;
    PhotosURL: string;
    DeletedImages: string[];
}

export interface CartProductResponse {
    product_id: string;
    product_name: string;
    price: number;
    stock: number;
    thumbnail_image: string
}

export interface product {
    id: number;
    name: string;
    category: string;
    price: number;
    description: string;
    image: string;
}

export interface productprops {
    product: GetProducts
    addtocart: (product: cartItem) => void
    isInCart?: boolean
    isDarkContainer?: boolean
}

export interface ProductDetailProps {
    addtocart: (product: cartItem) => void;
    cartItems: cartItem[];
}

export interface Createproduct {
    product_name: string;
    product_category: string[];
    price: string;
    description: string;
    stock: string;
    images: File[];
    item_type: string;
    technical_details?: ProductFeature[];
    publish_state:string;
}

export interface UpdateCartItem {
    product_id: string
    quantity: number
}