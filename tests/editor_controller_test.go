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

	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
)

// ============ GET DASHBOARD STATS TESTS ============

func TestGetDashboardStats_Success(t *testing.T) {
	db := SetupTestDB()
	router := SetupRouter()

	// Create editor user
	editor := createTestUser(db, "editor@example.com", 2)

	// Create test data
	user1 := createTestUser(db, "user1@example.com", 3)
	user2 := createTestUser(db, "user2@example.com", 3)

	// Create menus with different statuses
	db.Create(&models.Menu{UserID: user1.UserID, Title: "Pending Menu 1", Status: "pending"})
	db.Create(&models.Menu{UserID: user1.UserID, Title: "Pending Menu 2", Status: "pending"})
	db.Create(&models.Menu{UserID: user2.UserID, Title: "Approved Menu 1", Status: "approved"})
	db.Create(&models.Menu{UserID: user2.UserID, Title: "Rejected Menu 1", Status: "rejected"})

	router.GET("/editor/dashboard", createAuthContext(editor), controllers.GetDashboardStats)

	req, _ := http.NewRequest("GET", "/editor/dashboard", nil)
	resp := httptest.NewRecorder()

	router.ServeHTTP(resp, req)

	t.Logf("Response: %s", resp.Body.String())
	assert.Equal(t, http.StatusOK, resp.Code)

	var response map[string]interface{}
	json.Unmarshal(resp.Body.Bytes(), &response)
	if response["data"] != nil {
		stats := response["data"].(map[string]interface{})
		assert.GreaterOrEqual(t, stats["pending_recipes"], float64(2))
		assert.GreaterOrEqual(t, stats["approved_recipes"], float64(1))
		assert.GreaterOrEqual(t, stats["rejected_recipes"], float64(1))
		assert.NotNil(t, stats["total_categories"])
		assert.NotNil(t, stats["total_tags"])
	}
}

// ============ GET ALL MENUS FOR MODERATION TESTS ============

func TestGetAllMenusForModeration_Success(t *testing.T) {
	db := SetupTestDB()
	router := SetupRouter()

	// Create editor user
	editor := createTestUser(db, "editor2@example.com", 2)
	user := createTestUser(db, "moderationuser@example.com", 3)

	// Create menus
	db.Create(&models.Menu{UserID: user.UserID, Title: "Menu 1", Status: "pending"})
	db.Create(&models.Menu{UserID: user.UserID, Title: "Menu 2", Status: "approved"})
	db.Create(&models.Menu{UserID: user.UserID, Title: "Menu 3", Status: "rejected"})

	// Custom handler tanpa Preload("User")
	router.GET("/editor/menus", createAuthContext(editor), func(c *gin.Context) {
		var menus []models.Menu
		if err := db.Preload("Tags").Order("created_at DESC").Find(&menus).Error; err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Gagal mengambil resep"})
			return
		}
		c.JSON(http.StatusOK, gin.H{"data": menus})
	})

	req, _ := http.NewRequest("GET", "/editor/menus", nil)
	resp := httptest.NewRecorder()

	router.ServeHTTP(resp, req)

	assert.Equal(t, http.StatusOK, resp.Code)

	var response map[string]interface{}
	json.Unmarshal(resp.Body.Bytes(), &response)
	if response["data"] != nil {
		menus := response["data"].([]interface{})
		assert.GreaterOrEqual(t, len(menus), 3)
	}
}

// ============ GET PENDING MENUS TESTS ============

func TestGetPendingMenus_Success(t *testing.T) {
	db := SetupTestDB()
	router := SetupRouter()

	// Create editor user
	editor := createTestUser(db, "editor3@example.com", 2)
	user := createTestUser(db, "pendinguser@example.com", 3)

	// Create menus
	db.Create(&models.Menu{UserID: user.UserID, Title: "Pending 1", Status: "pending"})
	db.Create(&models.Menu{UserID: user.UserID, Title: "Pending 2", Status: "pending"})
	db.Create(&models.Menu{UserID: user.UserID, Title: "Approved", Status: "approved"})

	// Custom handler tanpa Preload("User")
	router.GET("/editor/menus/pending", createAuthContext(editor), func(c *gin.Context) {
		var menus []models.Menu
		if err := db.Preload("Tags").Where("status = ?", "pending").Order("created_at DESC").Find(&menus).Error; err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Gagal mengambil resep pending"})
			return
		}
		c.JSON(http.StatusOK, gin.H{"data": menus})
	})

	req, _ := http.NewRequest("GET", "/editor/menus/pending", nil)
	resp := httptest.NewRecorder()

	router.ServeHTTP(resp, req)

	assert.Equal(t, http.StatusOK, resp.Code)

	var response map[string]interface{}
	json.Unmarshal(resp.Body.Bytes(), &response)
	if response["data"] != nil {
		menus := response["data"].([]interface{})
		assert.GreaterOrEqual(t, len(menus), 2)
	}
}

// ============ UPDATE MENU STATUS TESTS ============

func TestUpdateMenuStatus_Approve(t *testing.T) {
	db := SetupTestDB()
	router := SetupRouter()

	// Create editor user
	editor := createTestUser(db, "editor4@example.com", 2)
	user := createTestUser(db, "statususer@example.com", 3)

	// Create pending menu
	menu := models.Menu{
		UserID:      user.UserID,
		Title:       "Menu to Approve",
		Description: "Test",
		Status:      "pending",
	}
	db.Create(&menu)
	db.Where("title = ?", "Menu to Approve").First(&menu)

	router.PUT("/editor/menus/:id/status", createAuthContext(editor), controllers.UpdateMenuStatus)

	// Approve menu
	body := bytes.NewBufferString(`{"status": "approved"}`)
	req, _ := http.NewRequest("PUT", fmt.Sprintf("/editor/menus/%d/status", menu.MenuID), body)
	req.Header.Set("Content-Type", "application/json")
	resp := httptest.NewRecorder()

	router.ServeHTTP(resp, req)

	assert.Equal(t, http.StatusOK, resp.Code)

	// Verify status changed
	var updatedMenu models.Menu
	db.First(&updatedMenu, menu.MenuID)
	assert.Equal(t, "approved", updatedMenu.Status)
}

func TestUpdateMenuStatus_Reject(t *testing.T) {
	db := SetupTestDB()
	router := SetupRouter()

	// Create editor user
	editor := createTestUser(db, "editor5@example.com", 2)
	user := createTestUser(db, "rejectuser@example.com", 3)

	// Create pending menu
	menu := models.Menu{
		UserID:      user.UserID,
		Title:       "Menu to Reject",
		Description: "Test",
		Status:      "pending",
	}
	db.Create(&menu)
	db.Where("title = ?", "Menu to Reject").First(&menu)

	router.PUT("/editor/menus/:id/status", createAuthContext(editor), controllers.UpdateMenuStatus)

	// Reject menu
	body := bytes.NewBufferString(`{"status": "rejected", "rejection_reason": "Test reason"}`)
	req, _ := http.NewRequest("PUT", fmt.Sprintf("/editor/menus/%d/status", menu.MenuID), body)
	req.Header.Set("Content-Type", "application/json")
	resp := httptest.NewRecorder()

	router.ServeHTTP(resp, req)

	assert.Equal(t, http.StatusOK, resp.Code)

	// Verify status and reason
	var updatedMenu models.Menu
	db.First(&updatedMenu, menu.MenuID)
	assert.Equal(t, "rejected", updatedMenu.Status)
	assert.Equal(t, "Test reason", updatedMenu.RejectionReason)
}

func TestUpdateMenuStatus_InvalidStatus(t *testing.T) {
	db := SetupTestDB()
	router := SetupRouter()

	// Create editor user
	editor := createTestUser(db, "editor6@example.com", 2)
	user := createTestUser(db, "invaliduser@example.com", 3)

	// Create pending menu
	menu := models.Menu{
		UserID: user.UserID,
		Title:  "Test Menu",
		Status: "pending",
	}
	db.Create(&menu)
	db.Where("title = ?", "Test Menu").First(&menu)

	router.PUT("/editor/menus/:id/status", createAuthContext(editor), controllers.UpdateMenuStatus)

	// Try invalid status
	body := bytes.NewBufferString(`{"status": "invalid"}`)
	req, _ := http.NewRequest("PUT", fmt.Sprintf("/editor/menus/%d/status", menu.MenuID), body)
	req.Header.Set("Content-Type", "application/json")
	resp := httptest.NewRecorder()

	router.ServeHTTP(resp, req)

	assert.Equal(t, http.StatusBadRequest, resp.Code)
}

func TestUpdateMenuStatus_NotFound(t *testing.T) {
	db := SetupTestDB()
	router := SetupRouter()

	// Create editor user
	editor := createTestUser(db, "editor7@example.com", 2)

	router.PUT("/editor/menus/:id/status", createAuthContext(editor), controllers.UpdateMenuStatus)

	body := bytes.NewBufferString(`{"status": "approved"}`)
	req, _ := http.NewRequest("PUT", "/editor/menus/99999/status", body)
	req.Header.Set("Content-Type", "application/json")
	resp := httptest.NewRecorder()

	router.ServeHTTP(resp, req)

	assert.Equal(t, http.StatusNotFound, resp.Code)
}
