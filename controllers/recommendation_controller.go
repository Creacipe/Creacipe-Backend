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
	log.Printf("[DEBUG-REC] ====== START GetRecommendations for menu_id=%s ======", menuID)

	// 2. Ambil title dari database berdasarkan menu_id.
	var menu models.Menu
	if err := config.DB.Select("title").First(&menu, menuID).Error; err != nil {
		log.Printf("[DEBUG-REC] Error finding menu title for ID %s: %v", menuID, err)
		c.JSON(http.StatusNotFound, gin.H{"error": "Resep tidak ditemukan di database"})
		return
	}
	title := menu.Title
	log.Printf("[DEBUG-REC] Title dari DB: '%s'", title)
	
	if title == "" {
		c.JSON(http.StatusNotFound, gin.H{"error": "Judul resep kosong"})
		return
	}

	// 3. Encode title agar aman dikirim di URL.
	encodedTitle := url.PathEscape(title)
	log.Printf("[DEBUG-REC] Title encoded: '%s'", encodedTitle)

	// 4. Panggil service Python dengan endpoint /recommend/title/<encoded_title>.
	pyServiceURL := "http://localhost:5000/recommend/title/" + encodedTitle
	log.Printf("[DEBUG-REC] Calling Python service: %s", pyServiceURL)
	
	resp, err := http.Get(pyServiceURL)
	if err != nil {
		log.Printf("[DEBUG-REC] Error contacting Python service: %v", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Gagal menghubungi service rekomendasi"})
		return
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		log.Printf("[DEBUG-REC] Error reading response: %v", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Gagal membaca respons dari service"})
		return
	}
	
	log.Printf("[DEBUG-REC] Python response status: %d", resp.StatusCode)
	log.Printf("[DEBUG-REC] Python response body: %s", string(body))

	// Teruskan error dari Python jika ada.
	if resp.StatusCode != http.StatusOK {
		log.Printf("[DEBUG-REC] Python returned error: %s", string(body))
		c.Data(resp.StatusCode, "application/json", body)
		return
	}

	// 5. Proses hasil (daftar JUDUL resep) dari Python.
	var recommendedTitles []string
	if err := json.Unmarshal(body, &recommendedTitles); err != nil {
		log.Printf("[DEBUG-REC] Error unmarshalling: %v", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Gagal memproses hasil rekomendasi (judul)"})
		return
	}
	log.Printf("[DEBUG-REC] Titles dari Python: %d items", len(recommendedTitles))
	for i, t := range recommendedTitles {
		if i < 5 { // Log first 5
			log.Printf("[DEBUG-REC]   [%d] '%s'", i, t)
		}
	}

	// 6. Ambil detail resep lengkap dari database berdasarkan JUDUL yang direkomendasikan dengan statistik.
	type MenuWithStats struct {
		models.Menu
		TotalLikes     int `json:"total_likes"`
		TotalDislikes  int `json:"total_dislikes"`
		TotalBookmarks int `json:"total_bookmarks"`
	}
	
	var recommendedMenus []MenuWithStats
	if len(recommendedTitles) > 0 {
		lowercaseTitles := make([]string, len(recommendedTitles))
		for i, title := range recommendedTitles {
			lowercaseTitles[i] = strings.ToLower(strings.TrimSpace(title))
		}
		
		log.Printf("[DEBUG-REC] Searching DB with %d lowercase titles", len(lowercaseTitles))
		for i, t := range lowercaseTitles {
			if i < 5 { // Log first 5
				log.Printf("[DEBUG-REC]   Search[%d]: '%s'", i, t)
			}
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
			log.Printf("[DEBUG-REC] DB Query Error: %v", err)
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Gagal mengambil detail resep dari judul rekomendasi"})
			return
		}
		
		log.Printf("[DEBUG-REC] Found in DB: %d menus (dari %d titles)", len(recommendedMenus), len(lowercaseTitles))
		for i, m := range recommendedMenus {
			if i < 5 {
				log.Printf("[DEBUG-REC]   Found[%d]: ID=%d, Title='%s'", i, m.MenuID, m.Title)
			}
		}
		
		// Debug: Cek title yang tidak match
		if len(recommendedMenus) < len(lowercaseTitles) {
			foundTitles := make(map[string]bool)
			for _, menu := range recommendedMenus {
				foundTitles[strings.ToLower(menu.Title)] = true
			}
			
			notFound := []string{}
			for _, title := range lowercaseTitles {
				if !foundTitles[title] {
					notFound = append(notFound, title)
				}
			}
			
			if len(notFound) > 0 {
				log.Printf("[DEBUG-REC] ⚠️ NOT FOUND in DB (%d items):", len(notFound))
				for i, t := range notFound {
					if i < 10 {
						log.Printf("[DEBUG-REC]   Missing[%d]: '%s'", i, t)
					}
				}
			}
		}
	} else {
		log.Println("[DEBUG-REC] Python returned empty list")
	}

	log.Printf("[DEBUG-REC] ====== END: Returning %d recommendations ======", len(recommendedMenus))
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

	// 7. Ambil detail resep lengkap dari database berdasarkan JUDUL yang direkomendasikan dengan statistik.
	type MenuWithStats struct {
		models.Menu
		TotalLikes     int `json:"total_likes"`
		TotalDislikes  int `json:"total_dislikes"`
		TotalBookmarks int `json:"total_bookmarks"`
	}
	
	var recommendedMenus []MenuWithStats
	if len(recommendedTitles) > 0 {
		// Convert semua recommendedTitles ke lowercase untuk matching
		lowercaseTitles := make([]string, len(recommendedTitles))
		for i, title := range recommendedTitles {
			lowercaseTitles[i] = strings.ToLower(strings.TrimSpace(title))
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
			log.Printf("Error fetching recommended menus by title with stats: %v", err)
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Gagal mengambil detail resep dari judul rekomendasi"})
			return
		}
		
		// Preload User untuk setiap menu
		menuIDs := make([]uint, len(recommendedMenus))
		for i, menu := range recommendedMenus {
			menuIDs[i] = menu.MenuID
		}
		
		var users []models.User
		config.DB.Where("user_id IN (SELECT user_id FROM menus WHERE menu_id IN ?)", menuIDs).Find(&users)
		userMap := make(map[uint]models.User)
		for _, u := range users {
			userMap[u.UserID] = u
		}
		
		for i := range recommendedMenus {
			if user, ok := userMap[recommendedMenus[i].UserID]; ok {
				recommendedMenus[i].User = user
			}
		}
	} else {
		log.Printf("Python service returned empty recommendation list for user %d", user.UserID)
	}

	c.JSON(http.StatusOK, gin.H{"data": recommendedMenus})
}

// GetEvaluationLogs - Endpoint untuk mengambil log evaluasi real-time.
func GetEvaluationLogs(c *gin.Context) {
	// 1. Ambil parameter limit dari query string (opsional)
	limit := c.DefaultQuery("limit", "50")
	
	// 2. Panggil endpoint logs di ML Service Python
	pyServiceURL := "http://localhost:5000/admin/logs?limit=" + limit
	
	resp, err := http.Get(pyServiceURL)
	if err != nil {
		log.Printf("Error contacting Python logs service: %v", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Gagal menghubungi service evaluasi"})
		return
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Gagal membaca respons logs"})
		return
	}

	// 3. Forward respons langsung ke Frontend
	c.Data(resp.StatusCode, "application/json", body)
}