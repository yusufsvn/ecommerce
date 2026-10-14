-- Kısıtlamaları (constraints) geri ekle
ALTER TABLE users ALTER COLUMN phone SET NOT NULL;
ALTER TABLE users ALTER COLUMN password SET NOT NULL;

-- Eklenen OAuth sütunlarını kaldır
ALTER TABLE users DROP COLUMN oauth_id;
ALTER TABLE users DROP COLUMN oauth_provider;