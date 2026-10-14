DROP INDEX IF EXISTS idx_products_name_trgm;

DROP EXTENSION IF EXISTS pg_trgm;

DROP INDEX IF EXISTS idx_search_vector;

ALTER TABLE products DROP COLUMN IF EXISTS search_vector;