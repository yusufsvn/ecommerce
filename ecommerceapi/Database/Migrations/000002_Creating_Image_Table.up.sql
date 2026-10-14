Create Table Images(
    image_id uuid Unique Default gen_random_uuid(),
    image_path varchar Unique,
    product_id uuid,
    created_at Timestamptz Not Null Default now(),
    updated_at Timestamptz,  

    Constraint PK_Image Primary Key(image_id),
    Constraint Fk_Product_Image Foreign Key(product_id) REFERENCES Products(product_id) ON Delete cascade
);