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

// ============ CREATE COMMENT TESTS ============

func TestCreateComment_Success(t *testing.T) {
	db := SetupTestDB()
	router := SetupRouter()

	// Create user and menu
	user := createTestUser(db, "user@example.com", 3)
	
	menu := models.Menu{
		UserID:      user.UserID,
		Title:       "Test Recipe",
		Description: "Test Description",
		Ingredients: "Test Ingredients",
		Instructions: "Test Instructions",
		Status:      "approved",
	}
	db.Create(&menu)
	db.Where("title = ?", "Test Recipe").First(&menu)

	router.POST("/menus/:id/comments", createAuthContext(user), controllers.CreateComment)

	input := map[string]interface{}{
		"comment_text": "This is a great recipe!",
	}

	body, _ := json.Marshal(input)
	req, _ := http.NewRequest("POST", fmt.Sprintf("/menus/%d/comments", menu.MenuID), bytes.NewBuffer(body))
	req.Header.Set("Content-Type", "application/json")
	resp := httptest.NewRecorder()

	router.ServeHTTP(resp, req)

	if resp.Code != http.StatusCreated {
		t.Logf("Response body: %s", resp.Body.String())
	}
	assert.Equal(t, http.StatusCreated, resp.Code)

	var response map[string]interface{}
	json.Unmarshal(resp.Body.Bytes(), &response)
	assert.Equal(t, "Komentar berhasil ditambahkan", response["message"])
}

func TestCreateComment_ReplyToComment(t *testing.T) {
	db := SetupTestDB()
	router := SetupRouter()

	// Create users and menu
	user1 := createTestUser(db, "user1@example.com", 3)
	user2 := createTestUser(db, "user2@example.com", 3)
	
	menu := models.Menu{
		UserID:      user1.UserID,
		Title:       "Test Recipe 2",
		Description: "Test Description 2",
		Ingredients: "Test Ingredients",
		Instructions: "Test Instructions",
		Status:      "approved",
	}
	db.Create(&menu)
	db.Where("title = ?", "Test Recipe 2").First(&menu)

	// Create parent comment
	parentComment := models.Comment{
		MenuID:      menu.MenuID,
		UserID:      user1.UserID,
		CommentText: "Original comment",
	}
	db.Create(&parentComment)
	db.Where("comment_text = ?", "Original comment").First(&parentComment)

	router.POST("/menus/:id/comments", createAuthContext(user2), controllers.CreateComment)

	input := map[string]interface{}{
		"comment_text": "This is a reply!",
		"parent_id":    parentComment.CommentID,
	}

	body, _ := json.Marshal(input)
	req, _ := http.NewRequest("POST", fmt.Sprintf("/menus/%d/comments", menu.MenuID), bytes.NewBuffer(body))
	req.Header.Set("Content-Type", "application/json")
	resp := httptest.NewRecorder()

	router.ServeHTTP(resp, req)

	assert.Equal(t, http.StatusCreated, resp.Code)

	// Verify notification was created
	var notification models.Notification
	err := db.Where("user_id = ? AND type = ?", user1.UserID, "info").First(&notification).Error
	assert.NoError(t, err)
}

func TestCreateComment_MissingText(t *testing.T) {
	db := SetupTestDB()
	router := SetupRouter()

	user := createTestUser(db, "user@example.com", 3)
	
	menu := models.Menu{
		UserID:      user.UserID,
		Title:       "Test Recipe 3",
		Description: "Test Description 3",
		Ingredients: "Test Ingredients",
		Instructions: "Test Instructions",
		Status:      "approved",
	}
	db.Create(&menu)
	db.Where("title = ?", "Test Recipe 3").First(&menu)

	router.POST("/menus/:id/comments", createAuthContext(user), controllers.CreateComment)

	input := map[string]interface{}{}

	body, _ := json.Marshal(input)
	req, _ := http.NewRequest("POST", fmt.Sprintf("/menus/%d/comments", menu.MenuID), bytes.NewBuffer(body))
	req.Header.Set("Content-Type", "application/json")
	resp := httptest.NewRecorder()

	router.ServeHTTP(resp, req)

	assert.Equal(t, http.StatusBadRequest, resp.Code)
}

func TestCreateComment_MenuNotFound(t *testing.T) {
	db := SetupTestDB()
	router := SetupRouter()

	user := createTestUser(db, "user@example.com", 3)

	router.POST("/menus/:id/comments", createAuthContext(user), controllers.CreateComment)

	input := map[string]interface{}{
		"comment_text": "This is a comment",
	}

	body, _ := json.Marshal(input)
	req, _ := http.NewRequest("POST", "/menus/99999/comments", bytes.NewBuffer(body))
	req.Header.Set("Content-Type", "application/json")
	resp := httptest.NewRecorder()

	router.ServeHTTP(resp, req)

	assert.Equal(t, http.StatusNotFound, resp.Code)
}

func TestCreateComment_ParentNotFound(t *testing.T) {
	db := SetupTestDB()
	router := SetupRouter()

	user := createTestUser(db, "user@example.com", 3)
	
	menu := models.Menu{
		UserID:      user.UserID,
		Title:       "Test Recipe 4",
		Description: "Test Description 4",
		Ingredients: "Test Ingredients",
		Instructions: "Test Instructions",
		Status:      "approved",
	}
	db.Create(&menu)
	db.Where("title = ?", "Test Recipe 4").First(&menu)

	router.POST("/menus/:id/comments", createAuthContext(user), controllers.CreateComment)

	input := map[string]interface{}{
		"comment_text": "This is a reply",
		"parent_id":    99999,
	}

	body, _ := json.Marshal(input)
	req, _ := http.NewRequest("POST", fmt.Sprintf("/menus/%d/comments", menu.MenuID), bytes.NewBuffer(body))
	req.Header.Set("Content-Type", "application/json")
	resp := httptest.NewRecorder()

	router.ServeHTTP(resp, req)

	assert.Equal(t, http.StatusNotFound, resp.Code)
}

// ============ GET COMMENTS TESTS ============

func TestGetCommentsByMenu_Success(t *testing.T) {
	db := SetupTestDB()
	router := SetupRouter()

	// Create user and menu
	user := createTestUser(db, "user@example.com", 3)
	
	menu := models.Menu{
		UserID:      user.UserID,
		Title:       "Test Recipe 5",
		Description: "Test Description 5",
		Ingredients: "Test Ingredients",
		Instructions: "Test Instructions",
		Status:      "approved",
	}
	db.Create(&menu)
	db.Where("title = ?", "Test Recipe 5").First(&menu)

	// Create comments
	comment1 := models.Comment{
		MenuID:      menu.MenuID,
		UserID:      user.UserID,
		CommentText: "First comment",
	}
	db.Create(&comment1)
	db.Where("comment_text = ?", "First comment").First(&comment1)

	comment2 := models.Comment{
		MenuID:      menu.MenuID,
		UserID:      user.UserID,
		CommentText: "Second comment",
	}
	db.Create(&comment2)
	db.Where("comment_text = ?", "Second comment").First(&comment2)

	router.GET("/menus/:id/comments", controllers.GetCommentsByMenu)

	req, _ := http.NewRequest("GET", fmt.Sprintf("/menus/%d/comments", menu.MenuID), nil)
	resp := httptest.NewRecorder()

	router.ServeHTTP(resp, req)

	assert.Equal(t, http.StatusOK, resp.Code)

	var response map[string]interface{}
	json.Unmarshal(resp.Body.Bytes(), &response)
	assert.GreaterOrEqual(t, int(response["total"].(float64)), 2)
}

func TestGetCommentsByMenu_WithReplies(t *testing.T) {
	db := SetupTestDB()
	router := SetupRouter()

	// Create users and menu
	user1 := createTestUser(db, "user1@example.com", 3)
	user2 := createTestUser(db, "user2@example.com", 3)
	
	menu := models.Menu{
		UserID:      user1.UserID,
		Title:       "Test Recipe 6",
		Description: "Test Description 6",
		Ingredients: "Test Ingredients",
		Instructions: "Test Instructions",
		Status:      "approved",
	}
	db.Create(&menu)
	db.Where("title = ?", "Test Recipe 6").First(&menu)

	// Create parent comment
	parentComment := models.Comment{
		MenuID:      menu.MenuID,
		UserID:      user1.UserID,
		CommentText: "Parent comment",
	}
	db.Create(&parentComment)
	db.Where("comment_text = ?", "Parent comment").First(&parentComment)

	// Create reply
	replyComment := models.Comment{
		MenuID:      menu.MenuID,
		UserID:      user2.UserID,
		ParentID:    &parentComment.CommentID,
		CommentText: "Reply comment",
	}
	db.Create(&replyComment)

	router.GET("/menus/:id/comments", controllers.GetCommentsByMenu)

	req, _ := http.NewRequest("GET", fmt.Sprintf("/menus/%d/comments", menu.MenuID), nil)
	resp := httptest.NewRecorder()

	router.ServeHTTP(resp, req)

	assert.Equal(t, http.StatusOK, resp.Code)
}

// ============ DELETE COMMENT TESTS ============

func TestDeleteComment_Success(t *testing.T) {
	db := SetupTestDB()
	router := SetupRouter()

	user := createTestUser(db, "user@example.com", 3)
	
	menu := models.Menu{
		UserID:      user.UserID,
		Title:       "Test Recipe 7",
		Description: "Test Description 7",
		Ingredients: "Test Ingredients",
		Instructions: "Test Instructions",
		Status:      "approved",
	}
	db.Create(&menu)
	db.Where("title = ?", "Test Recipe 7").First(&menu)

	comment := models.Comment{
		MenuID:      menu.MenuID,
		UserID:      user.UserID,
		CommentText: "Comment to delete",
	}
	db.Create(&comment)
	db.Where("comment_text = ?", "Comment to delete").First(&comment)

	router.DELETE("/comments/:id", createAuthContext(user), controllers.DeleteComment)

	req, _ := http.NewRequest("DELETE", fmt.Sprintf("/comments/%d", comment.CommentID), nil)
	resp := httptest.NewRecorder()

	router.ServeHTTP(resp, req)

	assert.Equal(t, http.StatusOK, resp.Code)

	// Verify comment was deleted
	var deletedComment models.Comment
	err := db.First(&deletedComment, comment.CommentID).Error
	assert.Error(t, err) // Should return error because comment is deleted
}

func TestDeleteComment_NotFound(t *testing.T) {
	db := SetupTestDB()
	router := SetupRouter()

	user := createTestUser(db, "user@example.com", 3)

	router.DELETE("/comments/:id", createAuthContext(user), controllers.DeleteComment)

	req, _ := http.NewRequest("DELETE", "/comments/99999", nil)
	resp := httptest.NewRecorder()

	router.ServeHTTP(resp, req)

	assert.Equal(t, http.StatusNotFound, resp.Code)
}

func TestDeleteComment_Forbidden(t *testing.T) {
	db := SetupTestDB()
	router := SetupRouter()

	// Create two users
	user1 := createTestUser(db, "user1@example.com", 3)
	user2 := createTestUser(db, "user2@example.com", 3)
	
	menu := models.Menu{
		UserID:      user1.UserID,
		Title:       "Test Recipe 8",
		Description: "Test Description 8",
		Ingredients: "Test Ingredients",
		Instructions: "Test Instructions",
		Status:      "approved",
	}
	db.Create(&menu)
	db.Where("title = ?", "Test Recipe 8").First(&menu)

	// User1 creates comment
	comment := models.Comment{
		MenuID:      menu.MenuID,
		UserID:      user1.UserID,
		CommentText: "User1's comment",
	}
	db.Create(&comment)
	db.Where("comment_text = ?", "User1's comment").First(&comment)

	// User2 tries to delete user1's comment
	router.DELETE("/comments/:id", createAuthContext(user2), controllers.DeleteComment)

	req, _ := http.NewRequest("DELETE", fmt.Sprintf("/comments/%d", comment.CommentID), nil)
	resp := httptest.NewRecorder()

	router.ServeHTTP(resp, req)

	assert.Equal(t, http.StatusForbidden, resp.Code)
}
