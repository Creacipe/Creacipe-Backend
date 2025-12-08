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

// ============ ADD TAG TO MENU TESTS ============

func TestAddTagToMenu_Success(t *testing.T) {
	db := SetupTestDB()
	router := SetupRouter()

	user := createTestUser(db, "user@example.com", 3)
	
	// Create category first
	category := models.Category{CategoryName: "Italian"}
	db.Create(&category)
	db.Where("category_name = ?", "Italian").First(&category)
	
	// Create tag
	tag := models.Tag{
		TagName:    "Spicy",
		CategoryID: category.CategoryID,
	}
	db.Create(&tag)
	db.Where("tag_name = ?", "Spicy").First(&tag)
	
	// Create menu
	menu := models.Menu{
		UserID:       user.UserID,
		Title:        "Test Recipe",
		Description:  "Test Description",
		Ingredients:  "Test Ingredients",
		Instructions: "Test Instructions",
		Status:       "approved",
	}
	db.Create(&menu)
	db.Where("title = ?", "Test Recipe").First(&menu)

	router.POST("/menus/:id/tags", createAuthContext(user), controllers.AddTagToMenu)

	input := map[string]interface{}{
		"tag_id": tag.TagID,
	}

	body, _ := json.Marshal(input)
	req, _ := http.NewRequest("POST", fmt.Sprintf("/menus/%d/tags", menu.MenuID), bytes.NewBuffer(body))
	req.Header.Set("Content-Type", "application/json")
	resp := httptest.NewRecorder()

	router.ServeHTTP(resp, req)

	assert.Equal(t, http.StatusOK, resp.Code)

	var response map[string]interface{}
	json.Unmarshal(resp.Body.Bytes(), &response)
	assert.Equal(t, "Tag berhasil ditambahkan ke resep", response["message"])
}

func TestAddTagToMenu_MenuNotFound(t *testing.T) {
	db := SetupTestDB()
	router := SetupRouter()

	user := createTestUser(db, "user@example.com", 3)

	router.POST("/menus/:id/tags", createAuthContext(user), controllers.AddTagToMenu)

	input := map[string]interface{}{
		"tag_id": 1,
	}

	body, _ := json.Marshal(input)
	req, _ := http.NewRequest("POST", "/menus/99999/tags", bytes.NewBuffer(body))
	req.Header.Set("Content-Type", "application/json")
	resp := httptest.NewRecorder()

	router.ServeHTTP(resp, req)

	assert.Equal(t, http.StatusNotFound, resp.Code)
}

func TestAddTagToMenu_TagNotFound(t *testing.T) {
	db := SetupTestDB()
	router := SetupRouter()

	user := createTestUser(db, "user@example.com", 3)
	
	menu := models.Menu{
		UserID:       user.UserID,
		Title:        "Test Recipe 2",
		Description:  "Test Description",
		Ingredients:  "Test Ingredients",
		Instructions: "Test Instructions",
		Status:       "approved",
	}
	db.Create(&menu)
	db.Where("title = ?", "Test Recipe 2").First(&menu)

	router.POST("/menus/:id/tags", createAuthContext(user), controllers.AddTagToMenu)

	input := map[string]interface{}{
		"tag_id": 99999,
	}

	body, _ := json.Marshal(input)
	req, _ := http.NewRequest("POST", fmt.Sprintf("/menus/%d/tags", menu.MenuID), bytes.NewBuffer(body))
	req.Header.Set("Content-Type", "application/json")
	resp := httptest.NewRecorder()

	router.ServeHTTP(resp, req)

	assert.Equal(t, http.StatusNotFound, resp.Code)
}

// ============ BOOKMARK TESTS ============

func TestBookmarkMenu_Success(t *testing.T) {
	db := SetupTestDB()
	router := SetupRouter()

	user := createTestUser(db, "user@example.com", 3)
	
	menu := models.Menu{
		UserID:       user.UserID,
		Title:        "Test Recipe 3",
		Description:  "Test Description",
		Ingredients:  "Test Ingredients",
		Instructions: "Test Instructions",
		Status:       "approved",
	}
	db.Create(&menu)
	db.Where("title = ?", "Test Recipe 3").First(&menu)

	router.POST("/menus/:id/bookmark", createAuthContext(user), controllers.BookmarkMenu)

	req, _ := http.NewRequest("POST", fmt.Sprintf("/menus/%d/bookmark", menu.MenuID), nil)
	resp := httptest.NewRecorder()

	router.ServeHTTP(resp, req)

	assert.Equal(t, http.StatusOK, resp.Code)

	var response map[string]interface{}
	json.Unmarshal(resp.Body.Bytes(), &response)
	assert.Equal(t, "Resep berhasil di-bookmark", response["message"])
}

func TestBookmarkMenu_NotFound(t *testing.T) {
	db := SetupTestDB()
	router := SetupRouter()

	user := createTestUser(db, "user@example.com", 3)

	router.POST("/menus/:id/bookmark", createAuthContext(user), controllers.BookmarkMenu)

	req, _ := http.NewRequest("POST", "/menus/99999/bookmark", nil)
	resp := httptest.NewRecorder()

	router.ServeHTTP(resp, req)

	assert.Equal(t, http.StatusNotFound, resp.Code)
}

func TestUnbookmarkMenu_Success(t *testing.T) {
	db := SetupTestDB()
	router := SetupRouter()

	user := createTestUser(db, "user@example.com", 3)
	
	menu := models.Menu{
		UserID:       user.UserID,
		Title:        "Test Recipe 4",
		Description:  "Test Description",
		Ingredients:  "Test Ingredients",
		Instructions: "Test Instructions",
		Status:       "approved",
	}
	db.Create(&menu)
	db.Where("title = ?", "Test Recipe 4").First(&menu)

	// First bookmark it
	db.Exec("INSERT INTO user_bookmarks (user_id, menu_id) VALUES (?, ?)", user.UserID, menu.MenuID)

	router.DELETE("/menus/:id/bookmark", createAuthContext(user), controllers.UnbookmarkMenu)

	req, _ := http.NewRequest("DELETE", fmt.Sprintf("/menus/%d/bookmark", menu.MenuID), nil)
	resp := httptest.NewRecorder()

	router.ServeHTTP(resp, req)

	assert.Equal(t, http.StatusOK, resp.Code)

	var response map[string]interface{}
	json.Unmarshal(resp.Body.Bytes(), &response)
	assert.Equal(t, "Bookmark berhasil dihapus", response["message"])
}

// ============ LIKE/DISLIKE TESTS ============

func TestLikeMenu_Success(t *testing.T) {
	db := SetupTestDB()
	router := SetupRouter()

	user := createTestUser(db, "user@example.com", 3)
	
	menu := models.Menu{
		UserID:       user.UserID,
		Title:        "Test Recipe 5",
		Description:  "Test Description",
		Ingredients:  "Test Ingredients",
		Instructions: "Test Instructions",
		Status:       "approved",
	}
	db.Create(&menu)
	db.Where("title = ?", "Test Recipe 5").First(&menu)

	router.POST("/menus/:id/like", createAuthContext(user), controllers.LikeMenu)

	req, _ := http.NewRequest("POST", fmt.Sprintf("/menus/%d/like", menu.MenuID), nil)
	resp := httptest.NewRecorder()

	router.ServeHTTP(resp, req)

	assert.Equal(t, http.StatusOK, resp.Code)

	var response map[string]interface{}
	json.Unmarshal(resp.Body.Bytes(), &response)
	assert.Equal(t, "Resep berhasil di-like", response["message"])
	assert.Equal(t, true, response["is_liked"])
}

func TestLikeMenu_ToggleOff(t *testing.T) {
	db := SetupTestDB()
	router := SetupRouter()

	user := createTestUser(db, "user@example.com", 3)
	
	menu := models.Menu{
		UserID:       user.UserID,
		Title:        "Test Recipe 6",
		Description:  "Test Description",
		Ingredients:  "Test Ingredients",
		Instructions: "Test Instructions",
		Status:       "approved",
	}
	db.Create(&menu)
	db.Where("title = ?", "Test Recipe 6").First(&menu)

	// Create existing like
	vote := models.MenuVote{
		UserID:     user.UserID,
		MenuID:     menu.MenuID,
		LikesCount: 1,
	}
	db.Create(&vote)

	router.POST("/menus/:id/like", createAuthContext(user), controllers.LikeMenu)

	req, _ := http.NewRequest("POST", fmt.Sprintf("/menus/%d/like", menu.MenuID), nil)
	resp := httptest.NewRecorder()

	router.ServeHTTP(resp, req)

	assert.Equal(t, http.StatusOK, resp.Code)

	var response map[string]interface{}
	json.Unmarshal(resp.Body.Bytes(), &response)
	assert.Equal(t, "Like berhasil dibatalkan", response["message"])
	assert.Equal(t, false, response["is_liked"])
}

func TestDislikeMenu_Success(t *testing.T) {
	db := SetupTestDB()
	router := SetupRouter()

	user := createTestUser(db, "user@example.com", 3)
	
	menu := models.Menu{
		UserID:       user.UserID,
		Title:        "Test Recipe 7",
		Description:  "Test Description",
		Ingredients:  "Test Ingredients",
		Instructions: "Test Instructions",
		Status:       "approved",
	}
	db.Create(&menu)
	db.Where("title = ?", "Test Recipe 7").First(&menu)

	router.POST("/menus/:id/dislike", createAuthContext(user), controllers.DislikeMenu)

	req, _ := http.NewRequest("POST", fmt.Sprintf("/menus/%d/dislike", menu.MenuID), nil)
	resp := httptest.NewRecorder()

	router.ServeHTTP(resp, req)

	assert.Equal(t, http.StatusOK, resp.Code)

	var response map[string]interface{}
	json.Unmarshal(resp.Body.Bytes(), &response)
	assert.Equal(t, "Resep berhasil di-dislike", response["message"])
	assert.Equal(t, true, response["is_disliked"])
}

func TestDislikeMenu_SwitchFromLike(t *testing.T) {
	db := SetupTestDB()
	router := SetupRouter()

	user := createTestUser(db, "user@example.com", 3)
	
	menu := models.Menu{
		UserID:       user.UserID,
		Title:        "Test Recipe 8",
		Description:  "Test Description",
		Ingredients:  "Test Ingredients",
		Instructions: "Test Instructions",
		Status:       "approved",
	}
	db.Create(&menu)
	db.Where("title = ?", "Test Recipe 8").First(&menu)

	// Create existing like
	vote := models.MenuVote{
		UserID:     user.UserID,
		MenuID:     menu.MenuID,
		LikesCount: 1,
	}
	db.Create(&vote)

	router.POST("/menus/:id/dislike", createAuthContext(user), controllers.DislikeMenu)

	req, _ := http.NewRequest("POST", fmt.Sprintf("/menus/%d/dislike", menu.MenuID), nil)
	resp := httptest.NewRecorder()

	router.ServeHTTP(resp, req)

	assert.Equal(t, http.StatusOK, resp.Code)

	var response map[string]interface{}
	json.Unmarshal(resp.Body.Bytes(), &response)
	assert.Equal(t, "Resep berhasil di-dislike", response["message"])
	assert.Equal(t, true, response["is_disliked"])
	assert.Equal(t, false, response["is_liked"])
}

// ============ GET INTERACTION STATUS TESTS ============

func TestGetUserInteractionStatus_Success(t *testing.T) {
	db := SetupTestDB()
	router := SetupRouter()

	user := createTestUser(db, "user@example.com", 3)
	
	menu := models.Menu{
		UserID:       user.UserID,
		Title:        "Test Recipe 9",
		Description:  "Test Description",
		Ingredients:  "Test Ingredients",
		Instructions: "Test Instructions",
		Status:       "approved",
	}
	db.Create(&menu)
	db.Where("title = ?", "Test Recipe 9").First(&menu)

	// Create like and bookmark
	vote := models.MenuVote{
		UserID:     user.UserID,
		MenuID:     menu.MenuID,
		LikesCount: 1,
	}
	db.Create(&vote)
	db.Exec("INSERT INTO user_bookmarks (user_id, menu_id) VALUES (?, ?)", user.UserID, menu.MenuID)

	router.GET("/menus/:id/interaction", createAuthContext(user), controllers.GetUserInteractionStatus)

	req, _ := http.NewRequest("GET", fmt.Sprintf("/menus/%d/interaction", menu.MenuID), nil)
	resp := httptest.NewRecorder()

	router.ServeHTTP(resp, req)

	assert.Equal(t, http.StatusOK, resp.Code)

	var response map[string]interface{}
	json.Unmarshal(resp.Body.Bytes(), &response)
	assert.Equal(t, true, response["is_liked"])
	assert.Equal(t, false, response["is_disliked"])
	assert.Equal(t, true, response["is_bookmarked"])
}

func TestGetUserInteractionStatus_NoInteractions(t *testing.T) {
	db := SetupTestDB()
	router := SetupRouter()

	user := createTestUser(db, "user@example.com", 3)
	
	menu := models.Menu{
		UserID:       user.UserID,
		Title:        "Test Recipe 10",
		Description:  "Test Description",
		Ingredients:  "Test Ingredients",
		Instructions: "Test Instructions",
		Status:       "approved",
	}
	db.Create(&menu)
	db.Where("title = ?", "Test Recipe 10").First(&menu)

	router.GET("/menus/:id/interaction", createAuthContext(user), controllers.GetUserInteractionStatus)

	req, _ := http.NewRequest("GET", fmt.Sprintf("/menus/%d/interaction", menu.MenuID), nil)
	resp := httptest.NewRecorder()

	router.ServeHTTP(resp, req)

	assert.Equal(t, http.StatusOK, resp.Code)

	var response map[string]interface{}
	json.Unmarshal(resp.Body.Bytes(), &response)
	assert.Equal(t, false, response["is_liked"])
	assert.Equal(t, false, response["is_disliked"])
	assert.Equal(t, false, response["is_bookmarked"])
}
