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

// ============ GET ALL TAGS TESTS ============

func TestGetAllTags_Success(t *testing.T) {
	db := SetupTestDB()
	router := SetupRouter()

	// Create category first
	category := models.Category{CategoryName: "Test Category"}
	db.Create(&category)

	// Create tags
	db.Create(&models.Tag{TagName: "Vegetarian", CategoryID: category.CategoryID})
	db.Create(&models.Tag{TagName: "Spicy", CategoryID: category.CategoryID})
	db.Create(&models.Tag{TagName: "Quick", CategoryID: category.CategoryID})

	router.GET("/tags", controllers.GetAllTags)

	req, _ := http.NewRequest("GET", "/tags", nil)
	resp := httptest.NewRecorder()

	router.ServeHTTP(resp, req)

	assert.Equal(t, http.StatusOK, resp.Code)

	var response map[string]interface{}
	json.Unmarshal(resp.Body.Bytes(), &response)
	if response["data"] != nil {
		tags := response["data"].([]interface{})
		assert.GreaterOrEqual(t, len(tags), 3)
	}
}

// ============ CREATE TAG TESTS ============

func TestCreateTag_Success(t *testing.T) {
	db := SetupTestDB()
	router := SetupRouter()

	// Create category first
	category := models.Category{CategoryName: "Test Category"}
	db.Create(&category)
	db.Where("category_name = ?", "Test Category").First(&category)
	
	if category.CategoryID == 0 {
		t.Fatalf("Category not created, CategoryID is 0")
	}

	admin := createTestUser(db, "tagadmin@example.com", 1)

	router.POST("/admin/tags", createAuthContext(admin), controllers.CreateTag)

	body := bytes.NewBufferString(fmt.Sprintf(`{"tag_name": "Low Carb", "category_id": %d}`, category.CategoryID))
	req, _ := http.NewRequest("POST", "/admin/tags", body)
	req.Header.Set("Content-Type", "application/json")
	resp := httptest.NewRecorder()

	router.ServeHTTP(resp, req)

	t.Logf("Request body: %s", fmt.Sprintf(`{"tag_name": "Low Carb", "category_id": %d}`, category.CategoryID))
	t.Logf("Response: %s", resp.Body.String())
	assert.Equal(t, http.StatusCreated, resp.Code)

	var response map[string]interface{}
	json.Unmarshal(resp.Body.Bytes(), &response)
	if response["data"] != nil {
		tag := response["data"].(map[string]interface{})
		assert.Equal(t, "Low Carb", tag["tag_name"])
	}
}

func TestCreateTag_MissingName(t *testing.T) {
	db := SetupTestDB()
	router := SetupRouter()

	admin := createTestUser(db, "tagadmin2@example.com", 1)

	router.POST("/admin/tags", createAuthContext(admin), controllers.CreateTag)

	body := bytes.NewBufferString(`{}`)
	req, _ := http.NewRequest("POST", "/admin/tags", body)
	req.Header.Set("Content-Type", "application/json")
	resp := httptest.NewRecorder()

	router.ServeHTTP(resp, req)

	assert.Equal(t, http.StatusBadRequest, resp.Code)
}

func TestCreateTag_Duplicate(t *testing.T) {
	db := SetupTestDB()
	router := SetupRouter()

	// Create category first
	category := models.Category{CategoryName: "Test Category"}
	db.Create(&category)

	admin := createTestUser(db, "tagadmin3@example.com", 1)

	// Create existing tag
	db.Create(&models.Tag{TagName: "Organic", CategoryID: category.CategoryID})

	router.POST("/admin/tags", createAuthContext(admin), controllers.CreateTag)

	body := bytes.NewBufferString(fmt.Sprintf(`{"tag_name": "Organic", "category_id": %d}`, category.CategoryID))
	req, _ := http.NewRequest("POST", "/admin/tags", body)
	req.Header.Set("Content-Type", "application/json")
	resp := httptest.NewRecorder()

	router.ServeHTTP(resp, req)

	// Controller returns 400/500 for all errors including duplicates
	assert.NotEqual(t, http.StatusCreated, resp.Code)
}

// ============ UPDATE TAG TESTS ============

func TestUpdateTag_Success(t *testing.T) {
	db := SetupTestDB()
	router := SetupRouter()

	// Create category first
	category := models.Category{CategoryName: "Test Category"}
	db.Create(&category)

	admin := createTestUser(db, "tagupdate@example.com", 1)

	// Create tag
	tag := models.Tag{TagName: "Old Tag", CategoryID: category.CategoryID}
	db.Create(&tag)
	db.Where("tag_name = ?", "Old Tag").First(&tag)

	router.PUT("/admin/tags/:id", createAuthContext(admin), controllers.UpdateTag)

	body := bytes.NewBufferString(`{"tag_name": "New Tag"}`)
	req, _ := http.NewRequest("PUT", fmt.Sprintf("/admin/tags/%d", tag.TagID), body)
	req.Header.Set("Content-Type", "application/json")
	resp := httptest.NewRecorder()

	router.ServeHTTP(resp, req)

	assert.Equal(t, http.StatusOK, resp.Code)

	// Verify update
	var updated models.Tag
	db.First(&updated, tag.TagID)
	assert.Equal(t, "New Tag", updated.TagName)
}

func TestUpdateTag_NotFound(t *testing.T) {
	db := SetupTestDB()
	router := SetupRouter()

	admin := createTestUser(db, "tagupdate2@example.com", 1)

	router.PUT("/admin/tags/:id", createAuthContext(admin), controllers.UpdateTag)

	body := bytes.NewBufferString(`{"name": "Test"}`)
	req, _ := http.NewRequest("PUT", "/admin/tags/99999", body)
	req.Header.Set("Content-Type", "application/json")
	resp := httptest.NewRecorder()

	router.ServeHTTP(resp, req)

	assert.Equal(t, http.StatusNotFound, resp.Code)
}

// ============ DELETE TAG TESTS ============

func TestDeleteTag_Success(t *testing.T) {
	db := SetupTestDB()
	router := SetupRouter()

	// Create category first
	category := models.Category{CategoryName: "Test Category"}
	db.Create(&category)

	admin := createTestUser(db, "tagdelete@example.com", 1)

	// Create tag
	tag := models.Tag{TagName: "To Delete", CategoryID: category.CategoryID}
	db.Create(&tag)
	db.Where("tag_name = ?", "To Delete").First(&tag)

	router.DELETE("/admin/tags/:id", createAuthContext(admin), controllers.DeleteTag)

	req, _ := http.NewRequest("DELETE", fmt.Sprintf("/admin/tags/%d", tag.TagID), nil)
	resp := httptest.NewRecorder()

	router.ServeHTTP(resp, req)

	assert.Equal(t, http.StatusOK, resp.Code)

	// Verify deletion
	var deleted models.Tag
	err := db.First(&deleted, tag.TagID).Error
	assert.Error(t, err)
}

func TestDeleteTag_NotFound(t *testing.T) {
	db := SetupTestDB()
	router := SetupRouter()

	admin := createTestUser(db, "tagdelete2@example.com", 1)

	router.DELETE("/admin/tags/:id", createAuthContext(admin), controllers.DeleteTag)

	req, _ := http.NewRequest("DELETE", "/admin/tags/99999", nil)
	resp := httptest.NewRecorder()

	router.ServeHTTP(resp, req)

	assert.Equal(t, http.StatusNotFound, resp.Code)
}
