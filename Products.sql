CREATE DATABASE streamoid;

CREATE TABLE products (
    sku VARCHAR(50) PRIMARY KEY,
    name VARCHAR(50) NOT NULL,
    brand VARCHAR(50) NOT NULL,
    color VARCHAR(50),
    size VARCHAR(20),
    mrp NUMERIC(10,2) CHECK (mrp > 0),
    price NUMERIC(10,2) CHECK (price <= mrp),
    quantity INT DEFAULT 0 CHECK (quantity >= 0)
);

CREATE INDEX idx_products_brand ON products(brand);
CREATE INDEX idx_products_color ON products(color);
CREATE INDEX idx_products_price ON products(price);

CREATE TEMP TABLE products_staging (
    sku VARCHAR(50),
    name VARCHAR(50),
    brand VARCHAR(50),
    color VARCHAR(50),
    size VARCHAR(20),
    mrp NUMERIC,
    price NUMERIC,
    quantity INT
) ON COMMIT DROP;

COPY products_staging(sku, name, brand, color, size, mrp, price, quantity)
FROM '%s' DELIMITER ',' CSV HEADER;

WITH inserted AS (
    INSERT INTO products (sku, name, brand, color, size, mrp, price, quantity)
    SELECT sku, name, brand, color, size, mrp, price, quantity
    FROM products_staging
    WHERE sku IS NOT NULL
        AND name IS NOT NULL
        AND brand IS NOT NULL
        AND mrp > 0
        AND price <= mrp
        AND quantity >= 0
    ON CONFLICT (sku) DO NOTHING
    RETURNING sku
)
SELECT COUNT(*) FROM inserted;