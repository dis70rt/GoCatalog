package database

import (
	"fmt"
	"database/sql"
	"github.com/dis70rt/streamoid/internal/models"
)

type ProductService struct {}

type ProductRequest struct {
	Brand 		string	`form:"brand"`
	Color 		string	`form:"color"`
	MinPrice 	float64	`form:"minPrice"`
	MaxPrice 	float64	`form:"maxPrice"`
	Page 		int		`form:"page"`
	Limit 		int		`form:"limit"`
}

func (pr *ProductRequest) Default() {
	if pr.Page <= 0 {pr.Page = 1}
	if pr.Limit <= 0 {pr.Limit = 10}
}

func (d *DatabaseService) GetAll(page int, limit int) (*[]models.Product, error) {
	if page <= 0 {page = 1}
	if limit <= 0 {limit = 10}

	query := `
        SELECT sku, name, brand, color, size, mrp, price, quantity
        FROM products
        ORDER BY sku
		LIMIT $1 OFFSET $2
    `
	offset := (page - 1) * limit
	rows, err := d.DB.Query(query, limit, offset)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var products []models.Product
	for rows.Next() {
		var p models.Product
		if err := rows.Scan(&p.SKU, &p.Name, &p.Brand, &p.Color, &p.Size, &p.MRP, &p.Price, &p.Quantity); err != nil {
			return nil, err
		}
		products = append(products, p)
	}

	return &products, nil
}


func (d *DatabaseService) GetBySKU(sku string) (*models.Product, error) {
	query := `
        SELECT sku, name, brand, color, size, mrp, price, quantity
        FROM products
        WHERE sku = $1
    `

	var p models.Product
	err := d.DB.QueryRow(query, sku).Scan(
		&p.SKU,
		&p.Name,
		&p.Brand,
		&p.Color,
		&p.Size,
		&p.MRP,
		&p.Price,
		&p.Quantity,
	)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, nil
		}
		return nil, err
	}

	return &p, nil
}

func (d *DatabaseService) Search(request *ProductRequest) (*[]models.Product, error) {
	request.Default()

	query := `
        SELECT sku, name, brand, color, size, mrp, price, quantity
        FROM products
        WHERE 1=1
    `
    args := []any{}
    argIndex := 1

	if request.Brand != "" {
        query += fmt.Sprintf(" AND brand = $%d", argIndex)
        args = append(args, request.Brand)
        argIndex++
    }

    if request.Color != "" {
        query += fmt.Sprintf(" AND color = $%d", argIndex)
        args = append(args, request.Color)
        argIndex++
    }

    if request.MinPrice > 0 {
        query += fmt.Sprintf(" AND price >= $%d", argIndex)
        args = append(args, request.MinPrice)
        argIndex++
    }

    if request.MaxPrice > 0 {
        query += fmt.Sprintf(" AND price <= $%d", argIndex)
        args = append(args, request.MaxPrice)
        argIndex++
    }

	offset := (request.Page - 1) * request.Limit
    query += fmt.Sprintf(" ORDER BY sku LIMIT $%d OFFSET $%d", argIndex, argIndex+1)
    args = append(args, request.Limit, offset)

	rows, err := d.DB.Query(query, args...)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var products []models.Product
	for rows.Next() {
		var p models.Product
        if err := rows.Scan(&p.SKU, &p.Name, &p.Brand, &p.Color, &p.Size, &p.MRP, &p.Price, &p.Quantity); err != nil {
            return nil, err
        }
        products = append(products, p)

	}

	return &products, nil
}