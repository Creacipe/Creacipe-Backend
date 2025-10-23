// controllers/recommendation_controller.go
package controllers

import (
	"creacipe-backend/config"
	"creacipe-backend/models"
	"encoding/json"
	"io"
	"net/http"
	"strconv"
	"strings"

	"github.com/gin-gonic/gin"
)

// GetRecommendations mengambil rekomendasi resep serupa dari service Python.
func GetRecommendations(c *gin.Context) {
	menuID := c.Param("id")

	// Panggil service Python
	resp, err := http.Get("http://localhost:5000/recommend/" + menuID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Gagal menghubungi service rekomendasi"})
		return
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Gagal membaca respons dari service"})
		return
	}

	if resp.StatusCode != http.StatusOK {
		c.Data(resp.StatusCode, "application/json", body)
		return
	}

	var recommendedIDs []int
	if err := json.Unmarshal(body, &recommendedIDs); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Gagal memproses hasil rekomendasi"})
		return
	}
	
	var recommendedMenus []models.Menu
	if len(recommendedIDs) > 0 {
		// Tambahkan Preload("User") untuk data autor dan penanganan error
		if err := config.DB.Preload("User").Where("menu_id IN ?", recommendedIDs).Find(&recommendedMenus).Error; err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Gagal mengambil detail resep rekomendasi"})
			return
		}
	}

	c.JSON(http.StatusOK, gin.H{"data": recommendedMenus})
}


// GetPersonalRecommendations mengambil rekomendasi personal berdasarkan riwayat user.
func GetPersonalRecommendations(c *gin.Context) {
	// 1. Ambil user yang sedang login dari context.
	user := c.MustGet("user").(models.User)

	// --- PERBAIKAN LOGIKA PENGAMBILAN RIWAYAT ---
	// 2. Ambil semua menu_id yang di-bookmark oleh user.
	var bookmarkedMenuIDs []int
	config.DB.Table("user_bookmarks").Where("user_id = ?", user.UserID).Pluck("menu_id", &bookmarkedMenuIDs)

	// 3. Ambil semua menu_id yang di-like (vote_type = 1).
	var likedMenuIDs []int
	config.DB.Table("menu_votes").Where("user_id = ? AND vote_type = ?", user.UserID, 1).Pluck("menu_id", &likedMenuIDs)
	
	// Gabungkan semua ID dan hilangkan duplikat.
	favoriteMenuIDs := append(bookmarkedMenuIDs, likedMenuIDs...)
	idMap := make(map[int]bool)
	uniqueIDs := []int{}
	for _, id := range favoriteMenuIDs {
		if _, value := idMap[id]; !value {
			idMap[id] = true
			uniqueIDs = append(uniqueIDs, id)
		}
	}
	// ---------------------------------------------

	// Jika user belum punya riwayat, kembalikan daftar kosong.
	if len(uniqueIDs) == 0 {
		c.JSON(http.StatusOK, gin.H{"data": []models.Menu{}})
		return
	}

	// 4. Ubah slice of int menjadi string "101,102,103".
	idsStrSlice := []string{}
	for _, id := range uniqueIDs {
		idsStrSlice = append(idsStrSlice, strconv.Itoa(id))
	}
	idsQueryParam := strings.Join(idsStrSlice, ",")

	// 5. Panggil service Python.
	resp, err := http.Get("http://localhost:5000/recommend/profile?ids=" + idsQueryParam)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Gagal menghubungi service rekomendasi"})
		return
	}
	defer resp.Body.Close()

	// --- BAGIAN YANG HILANG DILENGKAPI ---
	// 6. Baca dan proses respons dari Python.
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Gagal membaca respons dari service"})
		return
	}

	if resp.StatusCode != http.StatusOK {
		c.Data(resp.StatusCode, "application/json", body)
		return
	}

	var recommendedIDs []int
	if err := json.Unmarshal(body, &recommendedIDs); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Gagal memproses hasil rekomendasi"})
		return
	}

	// 7. Ambil detail resep dan kirim hasilnya.
	var recommendedMenus []models.Menu
	if len(recommendedIDs) > 0 {
		if err := config.DB.Preload("User").Where("menu_id IN ?", recommendedIDs).Find(&recommendedMenus).Error; err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Gagal mengambil detail resep rekomendasi"})
			return
		}
	}

	c.JSON(http.StatusOK, gin.H{"data": recommendedMenus})
}