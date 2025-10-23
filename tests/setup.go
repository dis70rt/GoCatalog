package tests

import (
	"bytes"
	"database/sql"
	"io"
	"mime/multipart"
	"net/http"
	"net/http/httptest"
	"os"
	"path/filepath"
	"strings"
	"testing"

	"github.com/dis70rt/streamoid/internal/database"
	"github.com/dis70rt/streamoid/logger"
	"github.com/dis70rt/streamoid/routes"
	"github.com/gin-gonic/gin"
)

type Test struct {
	T 		 *testing.T
	Router 	 *gin.Engine
	DB 	 	 *sql.DB
	Recorder *httptest.ResponseRecorder
}

func TestRouter(t *testing.T) *Test {
	gin.SetMode(gin.TestMode)
	logger.Init()

	router := gin.New()
	testDbConfig, err := database.TestPSQL()
	if err != nil {
		t.Fatalf("Couldn't connect to test-database: %v\n", err)
	}

	testDB, err := testDbConfig.Connect()
	if err != nil {
		t.Fatalf("Couldn't open test database: %v", err)
	}
	routes.RegisterRoutes(router, testDB)
	return &Test{
		T: t,
		Router: router,
		DB: testDB,
		Recorder: nil,
	}
}

func (test *Test) PerformRequest(method, path string) {
	req, _ := http.NewRequest(method, path, nil)
	test.Recorder = httptest.NewRecorder()
	test.Router.ServeHTTP(test.Recorder, req)
}

func (test Test) AssertResponse(expectedCode int, expectedBody string) {
	if test.Recorder.Code != expectedCode {
		test.T.Errorf("Expected status %d, got %d", expectedCode, test.Recorder.Code)
	}
	if strings.TrimSpace(test.Recorder.Body.String()) != expectedBody {
		test.T.Errorf("Expected body %s, got %s", expectedBody, test.Recorder.Body.String())
	}
}

func (test *Test) PerformUploadRequest(method, path, filePath, formFieldName string) {
    file, err := os.Open(filePath)
    if err != nil {
        test.T.Fatalf("Failed to open file: %v", err)
    }
    defer file.Close()

    body := &bytes.Buffer{}
    writer := multipart.NewWriter(body)
    part, err := writer.CreateFormFile(formFieldName, filepath.Base(filePath))
    if err != nil {
        test.T.Fatalf("Failed to create form file: %v", err)
    }
    _, err = io.Copy(part, file)
    if err != nil {
        test.T.Fatalf("Failed to copy file content: %v", err)
    }
    writer.Close()

    req, _ := http.NewRequest(method, path, body)
    req.Header.Set("Content-Type", writer.FormDataContentType())

    test.Recorder = httptest.NewRecorder()
    test.Router.ServeHTTP(test.Recorder, req)
}