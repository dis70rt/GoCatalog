package database

import (
    "encoding/csv"
    "fmt"
    "io"
    "strconv"
    "strings"
)

func (d *DatabaseService) ProcessCSVFromReader(r io.Reader) (int, []string, error) {
    csvReader := csv.NewReader(r)
    
    headers, err := csvReader.Read()
    if err != nil {
        return 0, nil, fmt.Errorf("failed to read CSV header: %w", err)
    }
    
    for i := range headers {
        headers[i] = strings.ToLower(strings.TrimSpace(headers[i]))
    }

    successCount := 0
    var failedProducts []string

    for {
        record, err := csvReader.Read()
        if err == io.EOF {
            break
        }
        if err != nil {
            failedProducts = append(failedProducts, fmt.Sprintf("Row error: %v", err))
            continue
        }

        product := make(map[string]string)
        for i, value := range record {
            if i < len(headers) {
                product[headers[i]] = value
            }
        }

        sku := product["sku"]
        if sku == "" {
            failedProducts = append(failedProducts, "Missing SKU")
            continue
        }

        mrp, err := strconv.ParseFloat(product["mrp"], 64)
        if err != nil {
            failedProducts = append(failedProducts, sku+": Invalid MRP")
            continue
        }

        price, err := strconv.ParseFloat(product["price"], 64)
        if err != nil {
            failedProducts = append(failedProducts, sku+": Invalid price")
            continue
        }

        quantity, err := strconv.Atoi(product["quantity"])
        if err != nil || quantity < 0 {
            failedProducts = append(failedProducts, sku+": Invalid quantity")
            continue
        }

        _, err = d.DB.Exec(`
            INSERT INTO products (sku, name, brand, color, size, mrp, price, quantity)
            VALUES ($1, $2, $3, $4, $5, $6, $7, $8)
        `,
            sku,
            product["name"],
            product["brand"],
            product["color"],
            product["size"],
            mrp,
            price,
            quantity,
        )

        if err != nil {
            failedProducts = append(failedProducts, sku+": "+err.Error())
            continue
        }

        successCount++
    }

    return successCount, failedProducts, nil
}