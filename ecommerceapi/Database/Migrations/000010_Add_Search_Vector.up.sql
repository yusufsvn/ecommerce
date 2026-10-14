alter table products add column search_vector tsvector 
generated always as (
    setweight(to_tsvector('turkish',coalesce(product_name,'')), 'A') ||
    setweight(to_tsvector('turkish',coalesce(description,'')), 'B') ||
    setweight(to_tsvector('turkish',coalesce(technical_details,'{}'::jsonb)),'C') ) stored ;

create index idx_search_vector on products using gin (search_vector);

CREATE EXTENSION IF NOT EXISTS pg_trgm;
CREATE INDEX idx_products_name_trgm ON products USING GIN (product_name gin_trgm_ops);