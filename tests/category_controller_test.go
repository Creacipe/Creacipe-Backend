package tests

import (
	"bytes"
	"creacipe-backend/controllers"
	"creacipe-backend/models"
	"encoding/json"
	"fmt"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/stretchr/testify/assert"
)

// ============ CREATE CATEGORY TESTS ============

func TestCreateCategory_Success(t *testing.T) {
	db := SetupTestDB()
	router := SetupRouter()

	// Create admin user
	admin := createTestUser(db, "categoryadmin@example.com", 1)

	router.POST("/admin/categories", createAuthContext(admin), controllers.CreateCategory)

	body := bytes.NewBufferString(`{"category_name": "Italian Cuisine"}`)
	req, _ := http.NewRequest("POST", "/admin/categories", body)
	req.Header.Set("Content-Type", "application/json")
	resp := httptest.NewRecorder()

	router.ServeHTTP(resp, req)

	assert.Equal(t, http.StatusCreated, resp.Code)

	var response map[string]interface{}
	json.Unmarshal(resp.Body.Bytes(), &response)
	if response["data"] != nil {
		category := response["data"].(map[string]interface{})
		assert.Equal(t, "Italian Cuisine", category["category_name"])
	}
}

func TestCreateCategory_MissingName(t *testing.T) {
	db := SetupTestDB()
	router := SetupRouter()

	admin := createTestUser(db, "categoryadmin2@example.com", 1)

	router.POST("/admin/categories", createAuthContext(admin), controllers.CreateCategory)

	body := bytes.NewBufferString(`{}`)
	req, _ := http.NewRequest("POST", "/admin/categories", body)
	req.Header.Set("Content-Type", "application/json")
	resp := httptest.NewRecorder()

	router.ServeHTTP(resp, req)

	assert.Equal(t, http.StatusBadRequest, resp.Code)
}

func TestCreateCategory_Duplicate(t *testing.T) {
	db := SetupTestDB()
	router := SetupRouter()

	admin := createTestUser(db, "categoryadmin3@example.com", 1)

	// Create existing category
	db.Create(&models.Category{CategoryName: "Japanese"})

	router.POST("/admin/categories", createAuthContext(admin), controllers.CreateCategory)

	body := bytes.NewBufferString(`{"category_name": "Japanese"}`)
	req, _ := http.NewRequest("POST", "/admin/categories", body)
	req.Header.Set("Content-Type", "application/json")
	resp := httptest.NewRecorder()

	router.ServeHTTP(resp, req)

	// Controller returns 500 for all DB errors including duplicates
	assert.Equal(t, http.StatusInternalServerError, resp.Code)
}

// ============ GET ALL CATEGORIES TESTS ============

func TestGetAllCategories_Success(t *testing.T) {
	db := SetupTestDB()
	router := SetupRouter()

	// Create categories
	db.Create(&models.Category{CategoryName: "Asian"})
	db.Create(&models.Category{CategoryName: "European"})
	db.Create(&models.Category{CategoryName: "American"})

	router.GET("/categories", controllers.GetAllCategories)

	req, _ := http.NewRequest("GET", "/categories", nil)
	resp := httptest.NewRecorder()

	router.ServeHTTP(resp, req)

	assert.Equal(t, http.StatusOK, resp.Code)

	var response map[string]interface{}
	json.Unmarshal(resp.Body.Bytes(), &response)
	if response["data"] != nil {
		categories := response["data"].([]interface{})
		assert.GreaterOrEqual(t, len(categories), 3)
	}
}

func TestGetAllCategories_Empty(t *testing.T) {
	SetupTestDB()
	router := SetupRouter()

	router.GET("/categories", controllers.GetAllCategories)

	req, _ := http.NewRequest("GET", "/categories", nil)
	resp := httptest.NewRecorder()

	router.ServeHTTP(resp, req)

	assert.Equal(t, http.StatusOK, resp.Code)
}

// ============ UPDATE CATEGORY TESTS ============

func TestUpdateCategory_Success(t *testing.T) {
	db := SetupTestDB()
	router := SetupRouter()

	admin := createTestUser(db, "categoryupdate@example.com", 1)

	// Create category
	category := models.Category{CategoryName: "Old Name"}
	db.Create(&category)
	db.Where("category_name = ?", "Old Name").First(&category)

	router.PUT("/admin/categories/:id", createAuthContext(admin), controllers.UpdateCategory)

	body := bytes.NewBufferString(`{"category_name": "New Name"}`)
	req, _ := http.NewRequest("PUT", fmt.Sprintf("/admin/categories/%d", category.CategoryID), body)
	req.Header.Set("Content-Type", "application/json")
	resp := httptest.NewRecorder()

	router.ServeHTTP(resp, req)

	assert.Equal(t, http.StatusOK, resp.Code)

	// Verify update
	var updated models.Category
	db.First(&updated, category.CategoryID)
	assert.Equal(t, "New Name", updated.CategoryName)
}

func TestUpdateCategory_NotFound(t *testing.T) {
	db := SetupTestDB()
	router := SetupRouter()

	admin := createTestUser(db, "categoryupdate2@example.com", 1)

	router.PUT("/admin/categories/:id", createAuthContext(admin), controllers.UpdateCategory)

	body := bytes.NewBufferString(`{"category_name": "Test"}`)
	req, _ := http.NewRequest("PUT", "/admin/categories/99999", body)
	req.Header.Set("Content-Type", "application/json")
	resp := httptest.NewRecorder()

	router.ServeHTTP(resp, req)

	assert.Equal(t, http.StatusNotFound, resp.Code)
}

func TestUpdateCategory_InvalidInput(t *testing.T) {
	db := SetupTestDB()
	router := SetupRouter()

	admin := createTestUser(db, "categoryupdate3@example.com", 1)

	// Create category
	category := models.Category{CategoryName: "Test Category"}
	db.Create(&category)

	router.PUT("/admin/categories/:id", createAuthContext(admin), controllers.UpdateCategory)

	body := bytes.NewBufferString(`{}`)
	req, _ := http.NewRequest("PUT", fmt.Sprintf("/admin/categories/%d", category.CategoryID), body)
	req.Header.Set("Content-Type", "application/json")
	resp := httptest.NewRecorder()

	router.ServeHTTP(resp, req)

	assert.Equal(t, http.StatusBadRequest, resp.Code)
}

// ============ DELETE CATEGORY TESTS ============

func TestDeleteCategory_Success(t *testing.T) {
	db := SetupTestDB()
	router := SetupRouter()

	admin := createTestUser(db, "categorydelete@example.com", 1)

	// Create category
	category := models.Category{CategoryName: "To Delete"}
	db.Create(&category)
	db.Where("category_name = ?", "To Delete").First(&category)

	router.DELETE("/admin/categories/:id", createAuthContext(admin), controllers.DeleteCategory)

	req, _ := http.NewRequest("DELETE", fmt.Sprintf("/admin/categories/%d", category.CategoryID), nil)
	resp := httptest.NewRecorder()

	router.ServeHTTP(resp, req)

	assert.Equal(t, http.StatusOK, resp.Code)

	// Verify deletion
	var deleted models.Category
	err := db.First(&deleted, category.CategoryID).Error
	assert.Error(t, err)
}

func TestDeleteCategory_NotFound(t *testing.T) {
	db := SetupTestDB()
	router := SetupRouter()

	admin := createTestUser(db, "categorydelete2@example.com", 1)

	router.DELETE("/admin/categories/:id", createAuthContext(admin), controllers.DeleteCategory)

	req, _ := http.NewRequest("DELETE", "/admin/categories/99999", nil)
	resp := httptest.NewRecorder()

	router.ServeHTTP(resp, req)

	assert.Equal(t, http.StatusNotFound, resp.Code)
}
