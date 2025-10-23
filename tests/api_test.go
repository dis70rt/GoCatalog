package tests

import (
    "net/http"
    "testing"
)

func TestHealthHandler(t *testing.T) {
    test := TestRouter(t)
    test.PerformRequest(http.MethodGet, "/health")
    test.AssertResponse(http.StatusOK, `{"status":"healthy"}`)
}

func TestGetProductsHandler(t *testing.T) {
    test := TestRouter(t)
    test.DB.Exec("TRUNCATE TABLE products RESTART IDENTITY CASCADE;")
    test.PerformRequest(http.MethodGet, "/products")
    test.AssertResponse(http.StatusOK, `{"products":null}`)
}

func TestUploadHandler(t *testing.T) {
    test := TestRouter(t)
    test.DB.Exec("TRUNCATE TABLE products RESTART IDENTITY CASCADE;")
    test.PerformUploadRequest(http.MethodPost, "/upload", "../sample.csv", "file")
    test.AssertResponse(http.StatusOK, `{"failed":["INVALID-002: pq: new row for relation \"products\" violates check constraint \"products_check\"","INVALID-003: Invalid quantity"],"stored":2}`)
}

func TestProductsSearchHandler(t *testing.T) {
    test := TestRouter(t)
    test.DB.Exec("TRUNCATE TABLE products RESTART IDENTITY CASCADE;")
    test.PerformUploadRequest(http.MethodPost, "/upload", "../sample.csv", "file")

    searchPath := "/products/search?brand=StreamThreads&color=Red"
    test.PerformRequest(http.MethodGet, searchPath)

    expectedBody := `{"products":[{"sku":"TSHIRT-RED-001","name":"Classic Cotton T-Shirt","brand":"StreamThreads","size":"M","color":"Red","mrp":799,"price":499,"quantity":5}]}`
    test.AssertResponse(http.StatusOK, expectedBody)
}

func TestProductsSearchWithPriceRange(t *testing.T) {
    test := TestRouter(t)
    test.DB.Exec("TRUNCATE TABLE products RESTART IDENTITY CASCADE;")
    test.PerformUploadRequest(http.MethodPost, "/upload", "../sample.csv", "file")

    searchPath := "/products/search?minPrice=500&maxPrice=550"
    test.PerformRequest(http.MethodGet, searchPath)

    expectedBody := `{"products":[{"sku":"TSHIRT-BLK-002","name":"Classic Cotton T-Shirt","brand":"StreamThreads","size":"L","color":"Black","mrp":799,"price":549,"quantity":12}]}`
    test.AssertResponse(http.StatusOK, expectedBody)
}

func TestProductsSearchWithLimit(t *testing.T) {
    test := TestRouter(t)
    test.DB.Exec("TRUNCATE TABLE products RESTART IDENTITY CASCADE;")
    test.PerformUploadRequest(http.MethodPost, "/upload", "../sample.csv", "file")

    searchPath := "/products/search?brand=StreamThreads&limit=1"
    test.PerformRequest(http.MethodGet, searchPath)

    expectedBody := `{"products":[{"sku":"TSHIRT-BLK-002","name":"Classic Cotton T-Shirt","brand":"StreamThreads","size":"L","color":"Black","mrp":799,"price":549,"quantity":12}]}`
    test.AssertResponse(http.StatusOK, expectedBody)
}

func TestProductsSearchWithOffset(t *testing.T) {
    test := TestRouter(t)
    test.DB.Exec("TRUNCATE TABLE products RESTART IDENTITY CASCADE;")
    test.PerformUploadRequest(http.MethodPost, "/upload", "../sample.csv", "file")

    searchPath := "/products/search?brand=StreamThreads&offset=1&limit=1"
    test.PerformRequest(http.MethodGet, searchPath)

    expectedBody := `{"products":[{"sku":"TSHIRT-BLK-002","name":"Classic Cotton T-Shirt","brand":"StreamThreads","size":"L","color":"Black","mrp":799,"price":549,"quantity":12}]}`
    test.AssertResponse(http.StatusOK, expectedBody)
}

func TestProductsSearchCombinedFilters(t *testing.T) {
    test := TestRouter(t)
    test.DB.Exec("TRUNCATE TABLE products RESTART IDENTITY CASCADE;")
    test.PerformUploadRequest(http.MethodPost, "/upload", "../sample.csv", "file")

    searchPath := "/products/search?brand=StreamThreads&minPrice=400&maxPrice=600"
    test.PerformRequest(http.MethodGet, searchPath)

    expectedBody := `{"products":[{"sku":"TSHIRT-BLK-002","name":"Classic Cotton T-Shirt","brand":"StreamThreads","size":"L","color":"Black","mrp":799,"price":549,"quantity":12},{"sku":"TSHIRT-RED-001","name":"Classic Cotton T-Shirt","brand":"StreamThreads","size":"M","color":"Red","mrp":799,"price":499,"quantity":5}]}`
    test.AssertResponse(http.StatusOK, expectedBody)
}