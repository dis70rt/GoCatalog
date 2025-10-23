package database

import (
	"fmt"

	"github.com/dis70rt/streamoid/internal/models"
	"github.com/dis70rt/streamoid/logger"
)

func (db *DatabaseService) UploadCSV(filePath string) (int, []models.FailedProduct, error) {
	tx, err := db.DB.Begin()
	if err != nil {
		logger.Error("cannot start transaction", err)
		return 0, nil, err
	}
	defer tx.Rollback()

	staging_query := `
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
	`

	if _, err := tx.Exec(staging_query); err != nil {
		logger.Error("cannot create staging products", err)
		return 0, nil, err
	}

	copyStaging := fmt.Sprintf(`
		COPY products_staging(sku, name, brand, color, size, mrp, price, quantity)
		FROM '%s' DELIMITER ',' CSV HEADER;
	`, filePath)

	if _, err := tx.Exec(copyStaging); err != nil {
		logger.Error("cannot copy csv into database", err)
		return 0, nil, err
	}

	countQuery := `
		WITH inserted AS (
			INSERT INTO products(sku, name, brand, color, size, mrp, price, quantity)
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
		) SELECT COUNT(*) FROM inserted;
	`

	var successCount int
	if err := tx.QueryRow(countQuery).Scan(&successCount); err != nil {
		logger.Error("error inserting into products", err)
		return 0, nil, err
	}

	failedProductsQuery := `
		SELECT 
			ps.sku,
			ps.name,
			ps.brand,
			ps.color,
			ps.size,
			ps.mrp,
			ps.price,
			ps.quantity,
			CASE
				WHEN ps.sku IS NULL THEN 'SKU missing'
				WHEN ps.name IS NULL THEN 'Name missing'
				WHEN ps.brand IS NULL THEN 'Brand missing'
				WHEN ps.mrp <= 0 THEN 'MRP invalid'
				WHEN ps.price > ps.mrp THEN 'Price > MRP'
				WHEN ps.quantity < 0 THEN 'Quantity invalid'
			END AS failed_reason
		FROM products_staging ps
		WHERE 
			ps.sku IS NULL
			OR ps.name IS NULL
			OR ps.brand IS NULL
			OR ps.mrp <= 0
			OR ps.price > ps.mrp
			OR ps.quantity < 0;
	`

	var failedProducts []models.FailedProduct

	rows, err := tx.Query(failedProductsQuery)
	if err != nil {
		logger.Error("error retriving failed products while inserting", err)
		return 0, nil, nil
	}

	for rows.Next() {
		var fp models.FailedProduct
		err := rows.Scan(
			&fp.SKU,
			&fp.Name,
			&fp.Brand,
			&fp.Color,
			&fp.Size,
			&fp.MRP,
			&fp.Price,
			&fp.Quantity,
			&fp.FailedReason,
		)
		if err != nil {
			logger.Error("failed to map into FailedProduct struct", err)
			return 0, nil, err
		}
		failedProducts = append(failedProducts, fp)
	}


	if err := tx.Commit(); err != nil {
		logger.Error("cannot commit transaction", err)
		return 0, nil, err
	}

	logger.Info(fmt.Sprintf("Successfully added %d new products", successCount))
	return successCount, failedProducts, nil
}