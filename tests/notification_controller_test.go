package tests

import (
	"creacipe-backend/controllers"
	"creacipe-backend/models"
	"encoding/json"
	"fmt"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/stretchr/testify/assert"
)

// ============ GET MY NOTIFICATIONS TESTS ============

func TestGetMyNotifications_Success(t *testing.T) {
	db := SetupTestDB()
	router := SetupRouter()

	user := createTestUser(db, "user@example.com", 3)

	// Create notifications
	notif1 := models.Notification{
		UserID:  user.UserID,
		Title:   "Test Notification 1",
		Message: "Message 1",
		Type:    "info",
		IsRead:  false,
	}
	db.Create(&notif1)

	notif2 := models.Notification{
		UserID:  user.UserID,
		Title:   "Test Notification 2",
		Message: "Message 2",
		Type:    "success",
		IsRead:  true,
	}
	db.Create(&notif2)

	router.GET("/notifications", createAuthContext(user), controllers.GetMyNotifications)

	req, _ := http.NewRequest("GET", "/notifications", nil)
	resp := httptest.NewRecorder()

	router.ServeHTTP(resp, req)

	assert.Equal(t, http.StatusOK, resp.Code)

	var response map[string]interface{}
	json.Unmarshal(resp.Body.Bytes(), &response)
	data := response["data"].([]interface{})
	assert.GreaterOrEqual(t, len(data), 2)
}

func TestGetMyNotifications_WithFilter(t *testing.T) {
	db := SetupTestDB()
	router := SetupRouter()

	user := createTestUser(db, "user@example.com", 3)

	// Create notifications
	notif1 := models.Notification{
		UserID:  user.UserID,
		Title:   "Unread Notification",
		Message: "Message",
		Type:    "info",
		IsRead:  false,
	}
	db.Create(&notif1)

	notif2 := models.Notification{
		UserID:  user.UserID,
		Title:   "Read Notification",
		Message: "Message",
		Type:    "info",
		IsRead:  true,
	}
	db.Create(&notif2)

	router.GET("/notifications", createAuthContext(user), controllers.GetMyNotifications)

	req, _ := http.NewRequest("GET", "/notifications?is_read=false", nil)
	resp := httptest.NewRecorder()

	router.ServeHTTP(resp, req)

	assert.Equal(t, http.StatusOK, resp.Code)

	var response map[string]interface{}
	json.Unmarshal(resp.Body.Bytes(), &response)
	data := response["data"].([]interface{})
	
	// Should only return unread notifications
	for _, item := range data {
		notif := item.(map[string]interface{})
		assert.Equal(t, false, notif["is_read"])
	}
}

func TestGetMyNotifications_WithPagination(t *testing.T) {
	db := SetupTestDB()
	router := SetupRouter()

	user := createTestUser(db, "user@example.com", 3)

	// Create multiple notifications
	for i := 1; i <= 5; i++ {
		notif := models.Notification{
			UserID:  user.UserID,
			Title:   fmt.Sprintf("Notification %d", i),
			Message: "Message",
			Type:    "info",
			IsRead:  false,
		}
		db.Create(&notif)
	}

	router.GET("/notifications", createAuthContext(user), controllers.GetMyNotifications)

	req, _ := http.NewRequest("GET", "/notifications?page=1&limit=2", nil)
	resp := httptest.NewRecorder()

	router.ServeHTTP(resp, req)

	assert.Equal(t, http.StatusOK, resp.Code)

	var response map[string]interface{}
	json.Unmarshal(resp.Body.Bytes(), &response)
	data := response["data"].([]interface{})
	pagination := response["pagination"].(map[string]interface{})
	
	assert.LessOrEqual(t, len(data), 2)
	assert.Equal(t, float64(1), pagination["page"])
	assert.Equal(t, float64(2), pagination["limit"])
}

// ============ GET UNREAD COUNT TESTS ============

func TestGetUnreadNotificationCount_Success(t *testing.T) {
	db := SetupTestDB()
	router := SetupRouter()

	user := createTestUser(db, "user@example.com", 3)

	// Create notifications
	for i := 1; i <= 3; i++ {
		notif := models.Notification{
			UserID:  user.UserID,
			Title:   fmt.Sprintf("Unread %d", i),
			Message: "Message",
			Type:    "info",
			IsRead:  false,
		}
		db.Create(&notif)
	}

	// Create read notification
	readNotif := models.Notification{
		UserID:  user.UserID,
		Title:   "Read",
		Message: "Message",
		Type:    "info",
		IsRead:  true,
	}
	db.Create(&readNotif)

	router.GET("/notifications/unread-count", createAuthContext(user), controllers.GetUnreadNotificationCount)

	req, _ := http.NewRequest("GET", "/notifications/unread-count", nil)
	resp := httptest.NewRecorder()

	router.ServeHTTP(resp, req)

	assert.Equal(t, http.StatusOK, resp.Code)

	var response map[string]interface{}
	json.Unmarshal(resp.Body.Bytes(), &response)
	assert.GreaterOrEqual(t, int(response["unread_count"].(float64)), 3)
}

func TestGetUnreadNotificationCount_NoUnread(t *testing.T) {
	db := SetupTestDB()
	router := SetupRouter()

	user := createTestUser(db, "nouunread@example.com", 3)

	router.GET("/notifications/unread-count", createAuthContext(user), controllers.GetUnreadNotificationCount)

	req, _ := http.NewRequest("GET", "/notifications/unread-count", nil)
	resp := httptest.NewRecorder()

	router.ServeHTTP(resp, req)

	assert.Equal(t, http.StatusOK, resp.Code)

	var response map[string]interface{}
	json.Unmarshal(resp.Body.Bytes(), &response)
	assert.Equal(t, float64(0), response["unread_count"])
}

// ============ MARK AS READ TESTS ============

func TestMarkNotificationAsRead_Success(t *testing.T) {
	db := SetupTestDB()
	router := SetupRouter()

	user := createTestUser(db, "user@example.com", 3)

	notif := models.Notification{
		UserID:  user.UserID,
		Title:   "Test Notification",
		Message: "Message",
		Type:    "info",
		IsRead:  false,
	}
	db.Create(&notif)
	db.Where("title = ?", "Test Notification").First(&notif)

	router.PATCH("/notifications/:id/read", createAuthContext(user), controllers.MarkNotificationAsRead)

	req, _ := http.NewRequest("PATCH", fmt.Sprintf("/notifications/%d/read", notif.NotificationID), nil)
	resp := httptest.NewRecorder()

	router.ServeHTTP(resp, req)

	assert.Equal(t, http.StatusOK, resp.Code)

	var response map[string]interface{}
	json.Unmarshal(resp.Body.Bytes(), &response)
	assert.Equal(t, "Notifikasi ditandai sudah dibaca", response["message"])

	// Verify notification is marked as read
	var updatedNotif models.Notification
	db.First(&updatedNotif, notif.NotificationID)
	assert.True(t, updatedNotif.IsRead)
}

func TestMarkNotificationAsRead_NotFound(t *testing.T) {
	db := SetupTestDB()
	router := SetupRouter()

	user := createTestUser(db, "user@example.com", 3)

	router.PATCH("/notifications/:id/read", createAuthContext(user), controllers.MarkNotificationAsRead)

	req, _ := http.NewRequest("PATCH", "/notifications/99999/read", nil)
	resp := httptest.NewRecorder()

	router.ServeHTTP(resp, req)

	assert.Equal(t, http.StatusNotFound, resp.Code)
}

func TestMarkNotificationAsRead_Unauthorized(t *testing.T) {
	db := SetupTestDB()
	router := SetupRouter()

	user1 := createTestUser(db, "user1@example.com", 3)
	user2 := createTestUser(db, "user2@example.com", 3)

	// User1's notification
	notif := models.Notification{
		UserID:  user1.UserID,
		Title:   "Test Notification",
		Message: "Message",
		Type:    "info",
		IsRead:  false,
	}
	db.Create(&notif)
	db.Where("title = ?", "Test Notification").First(&notif)

	// User2 tries to mark user1's notification as read
	router.PATCH("/notifications/:id/read", createAuthContext(user2), controllers.MarkNotificationAsRead)

	req, _ := http.NewRequest("PATCH", fmt.Sprintf("/notifications/%d/read", notif.NotificationID), nil)
	resp := httptest.NewRecorder()

	router.ServeHTTP(resp, req)

	assert.Equal(t, http.StatusNotFound, resp.Code)
}

// ============ MARK ALL AS READ TESTS ============

func TestMarkAllNotificationsAsRead_Success(t *testing.T) {
	db := SetupTestDB()
	router := SetupRouter()

	user := createTestUser(db, "user@example.com", 3)

	// Create unread notifications
	for i := 1; i <= 3; i++ {
		notif := models.Notification{
			UserID:  user.UserID,
			Title:   fmt.Sprintf("Unread %d", i),
			Message: "Message",
			Type:    "info",
			IsRead:  false,
		}
		db.Create(&notif)
	}

	router.PATCH("/notifications/read-all", createAuthContext(user), controllers.MarkAllNotificationsAsRead)

	req, _ := http.NewRequest("PATCH", "/notifications/read-all", nil)
	resp := httptest.NewRecorder()

	router.ServeHTTP(resp, req)

	assert.Equal(t, http.StatusOK, resp.Code)

	var response map[string]interface{}
	json.Unmarshal(resp.Body.Bytes(), &response)
	assert.Equal(t, "Semua notifikasi ditandai sudah dibaca", response["message"])

	// Verify all notifications are marked as read
	var unreadCount int64
	db.Model(&models.Notification{}).Where("user_id = ? AND is_read = false", user.UserID).Count(&unreadCount)
	assert.Equal(t, int64(0), unreadCount)
}

func TestMarkAllNotificationsAsRead_NoUnread(t *testing.T) {
	db := SetupTestDB()
	router := SetupRouter()

	user := createTestUser(db, "user@example.com", 3)

	router.PATCH("/notifications/read-all", createAuthContext(user), controllers.MarkAllNotificationsAsRead)

	req, _ := http.NewRequest("PATCH", "/notifications/read-all", nil)
	resp := httptest.NewRecorder()

	router.ServeHTTP(resp, req)

	assert.Equal(t, http.StatusOK, resp.Code)

	var response map[string]interface{}
	json.Unmarshal(resp.Body.Bytes(), &response)
	assert.Equal(t, "Semua notifikasi ditandai sudah dibaca", response["message"])
}
