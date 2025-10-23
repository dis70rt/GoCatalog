package models

import (
	"fmt"
	"strings"
)

type Product struct {
	SKU 		string 	`json:"sku"`
	Name 		string 	`json:"name"`
	Brand 		string 	`json:"brand"`
	Size 		string 	`json:"size"`
    Color       string  `json:"color"`
	MRP 		float64 `json:"mrp"`
	Price 		float64 `json:"price"`
	Quantity 	int32 	`json:"quantity"`
}

type FailedProduct struct {
	SKU 		string 	`json:"sku"`
	Name 		string 	`json:"name"`
	Brand 		string 	`json:"brand"`
	Size 		string 	`json:"size"`
    Color       string  `json:"color"`
	MRP 		float64 `json:"mrp"`
	Price 		float64 `json:"price"`
	Quantity 	int32 	`json:"quantity"`
    FailedReason string `json:"failed_reason"`
}

func (p *Product) Validate() error {
	var missing []string

    if strings.TrimSpace(p.SKU) == "" {
        missing = append(missing, "sku")
    }
    if strings.TrimSpace(p.Name) == "" {
        missing = append(missing, "name")
    }
    if strings.TrimSpace(p.Brand) == "" {
        missing = append(missing, "brand")
    }
    if p.MRP <= 0 {
        missing = append(missing, "mrp")
    }
    if p.Price <= 0 {
        missing = append(missing, "price")
    }

    if len(missing) > 0 {
        return fmt.Errorf("missing or invalid required fields: %s", strings.Join(missing, ", "))
    }

	if p.Price > p.MRP {
		return fmt.Errorf("price (%.2f) cannot be greater than mrp (%.2f)", p.Price, p.MRP)
	}

	if p.Quantity < 0 {
		return fmt.Errorf("quantity (%d) cannot be negative", p.Quantity)
	}

	return nil
}