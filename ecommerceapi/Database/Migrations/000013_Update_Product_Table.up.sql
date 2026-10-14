alter table products add column product_code varchar(8) unique not null;
alter table products add column product_publish_state varchar not null default 'unpublished';
