DROP INDEX IF EXISTS idx_search_vector;

ALTER TABLE products DROP COLUMN IF EXISTS search_vector;

ALTER TABLE products ADD COLUMN search_vector tsvector 
GENERATED ALWAYS AS (
    setweight(to_tsvector('turkish', coalesce(product_name, '')), 'A') ||
    setweight(to_tsvector('turkish', coalesce(description, '')), 'B') ||
    setweight(to_tsvector('turkish', coalesce(technical_details, '{}'::jsonb)), 'C')
) STORED;

CREATE INDEX idx_search_vector ON products USING GIN (search_vector);

-- order search vector

DROP INDEX IF EXISTS idx_order_search;
ALTER TABLE orders DROP COLUMN IF EXISTS search_vector;
