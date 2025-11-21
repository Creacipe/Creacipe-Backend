// controllers/recommendation_controller.go
package controllers

import (
	"creacipe-backend/config"
	"creacipe-backend/models"
	"encoding/json"
	"io"
	"log" // Import log
	"net/http"
	"net/url" // Import untuk URL encoding
	"strings"

	"github.com/gin-gonic/gin"
)

// GetRecommendations (by Title) - Mendapatkan resep serupa.
func GetRecommendations(c *gin.Context) {
	// 1. Ambil menu_id dari URL.
	menuID := c.Param("id")
	// log.Printf("[DEBUG] Menerima permintaan rekomendasi untuk menu_id: %s", menuID) // Log 1

	// 2. Ambil title dari database berdasarkan menu_id.
	var menu models.Menu
	if err := config.DB.Select("title").First(&menu, menuID).Error; err != nil {
		log.Printf("Error finding menu title for ID %s: %v", menuID, err)
		c.JSON(http.StatusNotFound, gin.H{"error": "Resep tidak ditemukan di database"})
		return
	}
	title := menu.Title
	// log.Printf("[DEBUG] Judul ditemukan di DB: '%s'", title)
	if title == "" {
		c.JSON(http.StatusNotFound, gin.H{"error": "Judul resep kosong"})
		return
	}

	// 3. Encode title agar aman dikirim di URL.
	encodedTitle := url.PathEscape(title)
	// log.Printf("[DEBUG] Judul diencode: '%s'", encodedTitle) // Log 3

	// 4. Panggil service Python dengan endpoint /recommend/title/<encoded_title>.
	pyServiceURL := "http://localhost:5000/recommend/title/" + encodedTitle
	// log.Printf("[DEBUG] Memanggil Python service: %s", pyServiceURL) // Log 4
	resp, err := http.Get(pyServiceURL)
	if err != nil {
		log.Printf("Error contacting Python service at %s: %v", pyServiceURL, err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Gagal menghubungi service rekomendasi"})
		return
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		log.Printf("Error reading Python service response body: %v", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Gagal membaca respons dari service"})
		return
	}

	// Teruskan error dari Python jika ada.
	if resp.StatusCode != http.StatusOK {
		log.Printf("Python service returned status %d: %s", resp.StatusCode, string(body))
		// log.Printf("[ERROR] Python service error (%d): %s", resp.StatusCode, string(body)) // Log 5
		c.Data(resp.StatusCode, "application/json", body)
		return
	}

	// 5. Proses hasil (daftar JUDUL resep) dari Python.
	var recommendedTitles []string
	if err := json.Unmarshal(body, &recommendedTitles); err != nil {
		// log.Printf("[ERROR] Gagal unmarshal judul rekomendasi: %v", err) // Log 6
		log.Printf("Error unmarshalling recommended titles: %v", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Gagal memproses hasil rekomendasi (judul)"})
		return
	}
	log.Printf("[DEBUG] Judul rekomendasi diterima dari Python: %v", recommendedTitles) // Log 7

	// 6. Ambil detail resep lengkap dari database berdasarkan JUDUL yang direkomendasikan.
	// 6. Ambil detail resep lengkap dari database berdasarkan JUDUL yang direkomendasikan dengan statistik.
	type MenuWithStats struct {
		models.Menu
		TotalLikes     int `json:"total_likes"`
		TotalDislikes  int `json:"total_dislikes"`
		TotalBookmarks int `json:"total_bookmarks"`
	}
	
	var recommendedMenus []MenuWithStats
	if len(recommendedTitles) > 0 {
		// log.Printf("[DEBUG] Mencari menu di DB dengan judul: %v", recommendedTitles) // Log 8
		// Convert semua recommendedTitles ke lowercase untuk matching
		lowercaseTitles := make([]string, len(recommendedTitles))
		for i, title := range recommendedTitles {
			lowercaseTitles[i] = strings.ToLower(title)
		}
		// Query dengan statistik vote dan bookmark
		query := `
			SELECT 
				m.menu_id,
				m.user_id,
				m.title,
				m.description,
				m.ingredients,
				m.instructions,
				m.image_url,
				m.status,
				m.created_at,
				m.updated_at,
				COALESCE(SUM(mv.likes_count), 0) as total_likes,
				COALESCE(SUM(mv.dislikes_count), 0) as total_dislikes,
				(SELECT COUNT(*) FROM user_bookmarks WHERE menu_id = m.menu_id) as total_bookmarks
			FROM menus m
			LEFT JOIN menu_votes mv ON m.menu_id = mv.menu_id
			WHERE LOWER(m.title) IN ? AND m.status = 'approved'
			GROUP BY m.menu_id, m.user_id, m.title, m.description, m.ingredients, m.instructions,
			         m.image_url, m.status, m.created_at, m.updated_at
		`
		
		if err := config.DB.Raw(query, lowercaseTitles).Scan(&recommendedMenus).Error; err != nil {
			// log.Printf("[ERROR] Gagal mengambil detail menu rekomendasi: %v", err) // Log 9
			log.Printf("Error fetching recommended menus by title: %v", err)
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Gagal mengambil detail resep dari judul rekomendasi"})
			return
		}
		log.Printf("[DEBUG] Menu rekomendasi ditemukan di DB: %d", len(recommendedMenus)) // Log 10
	} else {
		log.Println("Python service returned empty recommendation list.")
	}

	c.JSON(http.StatusOK, gin.H{"data": recommendedMenus})
}

// GetPersonalRecommendations (by Titles) - Mendapatkan rekomendasi personal.
func GetPersonalRecommendations(c *gin.Context) {
	// 1. Ambil user yang sedang login.
	user := c.MustGet("user").(models.User)

	// 2. Ambil ID resep favorit (bookmarks & likes).
	var bookmarkedMenuIDs []uint
	config.DB.Table("user_bookmarks").Where("user_id = ?", user.UserID).Pluck("menu_id", &bookmarkedMenuIDs)
	var likedMenuIDs []uint
	config.DB.Table("menu_votes").Where("user_id = ? AND vote_type = ?", user.UserID, 1).Pluck("menu_id", &likedMenuIDs)

	favoriteMenuIDs := append(bookmarkedMenuIDs, likedMenuIDs...)
	idMap := make(map[uint]bool)
	uniqueIDs := []uint{}
	for _, id := range favoriteMenuIDs {
		if _, value := idMap[id]; !value {
			idMap[id] = true
			uniqueIDs = append(uniqueIDs, id)
		}
	}

	if len(uniqueIDs) == 0 {
		c.JSON(http.StatusOK, gin.H{"data": []models.Menu{}})
		return
	}

	// 3. Ambil JUDUL dari database untuk ID favorit tersebut.
	var favoriteTitles []string
	if err := config.DB.Model(&models.Menu{}).Where("menu_id IN ?", uniqueIDs).Pluck("title", &favoriteTitles).Error; err != nil {
		log.Printf("Error fetching favorite titles for user %d: %v", user.UserID, err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Gagal mengambil judul resep favorit"})
		return
	}
	if len(favoriteTitles) == 0 {
		log.Printf("No valid titles found for favorite IDs of user %d", user.UserID)
		c.JSON(http.StatusOK, gin.H{"data": []models.Menu{}})
		return
	}

	// 4. Encode setiap judul dan gabungkan menjadi string query parameter.
	encodedTitles := []string{}
	for _, title := range favoriteTitles {
		encodedTitles = append(encodedTitles, url.QueryEscape(title))
	}
	titlesQueryParam := strings.Join(encodedTitles, ",")

	// 5. Panggil service Python dengan endpoint /recommend/profile?titles=....
	pyServiceURL := "http://localhost:5000/recommend/profile?titles=" + titlesQueryParam
	resp, err := http.Get(pyServiceURL)
	if err != nil {
		log.Printf("Error contacting Python service at %s: %v", pyServiceURL, err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Gagal menghubungi service rekomendasi"})
		return
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		log.Printf("Error reading Python service response body: %v", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Gagal membaca respons"})
		return
	}
	if resp.StatusCode != http.StatusOK {
		log.Printf("Python service returned status %d: %s", resp.StatusCode, string(body))
		c.Data(resp.StatusCode, "application/json", body)
		return
	}

	// 6. Proses hasil (daftar JUDUL resep) dari Python.
	var recommendedTitles []string
	if err := json.Unmarshal(body, &recommendedTitles); err != nil {
		log.Printf("Error unmarshalling recommended titles: %v", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Gagal memproses hasil rekomendasi (judul)"})
		return
	}

	// 7. Ambil detail resep lengkap dari database berdasarkan JUDUL yang direkomendasikan.
	var recommendedMenus []models.Menu
	if len(recommendedTitles) > 0 {
		if err := config.DB.Preload("User").Where("LOWER(title) IN ?", recommendedTitles).Find(&recommendedMenus).Error; err != nil {
			log.Printf("Error fetching recommended menus by title: %v", err)
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Gagal mengambil detail resep dari judul rekomendasi"})
			return
		}
	} else {
		log.Printf("Python service returned empty recommendation list for user %d", user.UserID)
	}

	c.JSON(http.StatusOK, gin.H{"data": recommendedMenus})
}