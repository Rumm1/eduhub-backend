ALTER TABLE organizations
ADD COLUMN default_language VARCHAR(2) NOT NULL DEFAULT 'ru',
ADD COLUMN name_ru VARCHAR(255),
ADD COLUMN name_kk VARCHAR(255),
ADD COLUMN name_en VARCHAR(255);

UPDATE organizations
SET name_ru = name
WHERE name_ru IS NULL;

ALTER TABLE organizations
ADD CONSTRAINT organizations_default_language_check
CHECK (default_language IN ('ru', 'kk', 'en'));
