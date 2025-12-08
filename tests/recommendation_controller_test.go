package tests

import (
	"creacipe-backend/controllers"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/stretchr/testify/assert"
)

// ============ GET RECOMMENDATIONS BY TITLE TESTS ============

func TestGetRecommendations_MenuNotFound(t *testing.T) {
	db := SetupTestDB()
	router := SetupRouter()

	user := createTestUser(db, "user@example.com", 3)

	router.GET("/recommendations/:id", createAuthContext(user), controllers.GetRecommendations)

	req, _ := http.NewRequest("GET", "/recommendations/99999", nil)
	resp := httptest.NewRecorder()

	router.ServeHTTP(resp, req)

	assert.Equal(t, http.StatusNotFound, resp.Code)

	var response map[string]interface{}
	json.Unmarshal(resp.Body.Bytes(), &response)
	assert.Contains(t, response["error"], "tidak ditemukan")
}

func TestGetRecommendations_EmptyTitle(t *testing.T) {
	db := SetupTestDB()
	router := SetupRouter()

	user := createTestUser(db, "user2@example.com", 3)

	// Create menu with empty title
	db.Exec("INSERT INTO menus (user_id, title, description, ingredients, instructions, status) VALUES (?, '', 'desc', 'ing', 'inst', 'approved')", user.UserID)

	var menuID uint
	db.Raw("SELECT menu_id FROM menus WHERE title = ''").Scan(&menuID)

	router.GET("/recommendations/:id", createAuthContext(user), controllers.GetRecommendations)

	req, _ := http.NewRequest("GET", "/recommendations/"+string(rune(menuID)), nil)
	resp := httptest.NewRecorder()

	router.ServeHTTP(resp, req)

	// Should return error for empty title
	assert.Equal(t, http.StatusNotFound, resp.Code)
}

// ============ GET PERSONAL RECOMMENDATIONS TESTS ============

func TestGetPersonalRecommendations_NoFavorites(t *testing.T) {
	db := SetupTestDB()
	router := SetupRouter()

	// Create user with no bookmarks or likes
	user := createTestUser(db, "nofav@example.com", 3)

	router.GET("/recommendations/personal", createAuthContext(user), controllers.GetPersonalRecommendations)

	req, _ := http.NewRequest("GET", "/recommendations/personal", nil)
	resp := httptest.NewRecorder()

	router.ServeHTTP(resp, req)

	// Should return 200 with empty or error message
	assert.Contains(t, []int{http.StatusOK, http.StatusNotFound}, resp.Code)

	var response map[string]interface{}
	json.Unmarshal(resp.Body.Bytes(), &response)
	
	// Either empty data or error message
	if resp.Code == http.StatusOK {
		if response["data"] != nil {
			data := response["data"].([]interface{})
			assert.Equal(t, 0, len(data))
		}
	} else {
		assert.NotNil(t, response["error"])
	}
}
