Create Table Users (
    user_id uuid Unique Default gen_random_uuid(),
    password varchar Not Null,
    name varchar Not Null,
    surname varchar Not Null,
    phone varchar Not Null,
    email varchar Not Null Unique,
    created_at Timestamptz Not Null Default now(),
    updated_at Timestamptz,
    Constraint Pk_User Primary Key(user_id)
);

Create Table Addresses (
    address_id uuid Unique Default gen_random_uuid(),
    user_id uuid,
    name varchar,
    post_code varchar, 
    country varchar,
    city varchar,
    district varchar,
    full_address varchar, 
    created_at Timestamptz Not Null Default now(),
    updated_at Timestamptz,

    Constraint Pk_Address Primary Key(address_id),
    Constraint Fk_User_Address Foreign Key (user_id) 
    REFERENCES Users(user_id) ON delete cascade
);

Create Table Orders(
    order_id uuid Unique Default gen_random_uuid(),
    user_id uuid,
    payment_method varchar,
    status varchar,
    total_amount Decimal(10, 2), 
    shipping_address varchar,
    billing_address varchar,

    created_at Timestamptz Not Null Default now(),
    updated_at Timestamptz,

    Constraint Pk_Order Primary Key(order_id),
    Constraint Fk_User_Id Foreign Key (user_id) 
    REFERENCES Users(user_id) ON delete cascade
);

Create Table Products (
    product_id uuid Unique Default gen_random_uuid(),
    product_name varchar,
    stock int,
    photos_url varchar,
    description varchar,
    price bigint,
    product_features_url varchar,
    created_at Timestamptz Default now(),
    updated_at Timestamptz,

    Constraint Pk_Product Primary Key(product_id)
);

Create Table ProductCategories(
    category_id uuid Unique Default gen_random_uuid(), 
    category_name varchar,
    category_enum_no int,

    created_at Timestamptz Not Null Default now(),
    updated_at Timestamptz,
    Constraint Pk_ProductCategories Primary Key(category_id)
);

Create Table Carts(
    cart_id uuid Unique Default gen_random_uuid(), 
    user_id uuid,
    status varchar,
    created_at Timestamptz Not Null Default now(),
    updated_at Timestamptz,
    
    Constraint Pk_Cart Primary Key(cart_id),
    Constraint Fk_User_Id Foreign Key (user_id) 
    REFERENCES Users(user_id) ON delete cascade
);

Create Table Payments(
    payment_id uuid Unique Default gen_random_uuid(), 
    order_id uuid, 
    payment_method varchar,
    amount Decimal(10, 2),
    status varchar,

    Constraint Pk_Payment Primary Key(payment_id),
    Constraint Fk_Order_Id Foreign Key (order_id) 
    REFERENCES Orders(order_id) ON delete cascade
);

-- İlişki Tabloları

Create Table cartProducts(
    cart_id uuid,
    product_id uuid,
    quantity int,
    Constraint Pk_Cart_Products Primary key (cart_id,product_id),
    Constraint Fk_Product_Id Foreign key (product_id) REFERENCES Products(product_id),
    Constraint Fk_Cart_Id Foreign key (cart_id) REFERENCES Carts(cart_id)
);

Create Table productProductcategories(
    product_id uuid,
    category_id uuid,
    
    Constraint Pk_Product_Product_Categories Primary key (product_id,category_id),
    Constraint Fk_Product_Id_Cat Foreign key (product_id) REFERENCES Products(product_id),
    Constraint Fk_Product_Category_Id Foreign key (category_id) REFERENCES ProductCategories(category_id)
);

Create Table productOrders(
    product_id uuid,
    order_id uuid,
    quantity int,
    price_at_purchase Decimal(10, 2), 
    
    Constraint Pk_Product_Orders Primary key (product_id,order_id),
    Constraint Fk_Product_Id_Ord Foreign key (product_id) REFERENCES Products(product_id),
    Constraint Fk_Order_Id_Ord Foreign key (order_id) REFERENCES Orders(order_id)
);