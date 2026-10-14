Alter table products drop column item_type;
alter table orders alter column total_amount type numeric(10,2);
alter table productorders alter column price_at_purchase type numeric(10,2);
alter table products rename column thumbnail_image to photos_url ;

alter table productorders drop column created_at Timestamptz Not Null Default now();
alter table productorders drop column updated_at Timestamptz;
alter table productProductcategories drop column created_at Timestamptz Not Null Default now();
alter table productProductcategories drop column updated_at Timestamptz;
alter table cartProducts drop column created_at Timestamptz Not Null Default now();
alter table cartProducts drop column updated_at Timestamptz;