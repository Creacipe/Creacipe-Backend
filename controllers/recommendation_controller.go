package controllers

import (
	"creacipe-backend/config"
	"creacipe-backend/models"
	"encoding/json"
	"io"
	"net/http"

	"strings"
  "strconv"

	"github.com/gin-gonic/gin"
)

// GetRecommendations mengambil rekomendasi resep serupa dari service Python.
func GetRecommendations(c *gin.Context) {
	menuID := c.Param("id")

	// Panggil service Python yang berjalan di port 5000.
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

	// Teruskan error dari service Python jika ada (misal: 404 Not Found).
	if resp.StatusCode != http.StatusOK {
		c.Data(resp.StatusCode, "application/json", body)
		return
	}

	// Proses hasil (daftar ID resep) dari Python.
	var recommendedIDs []int
	if err := json.Unmarshal(body, &recommendedIDs); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Gagal memproses hasil rekomendasi"})
		return
	}
	
	// Ambil detail lengkap resep dari database berdasarkan ID yang didapat.
	var recommendedMenus []models.Menu
	if len(recommendedIDs) > 0 {
		config.DB.Where("menu_id IN ?", recommendedIDs).Find(&recommendedMenus)
	}

	c.JSON(http.StatusOK, gin.H{"data": recommendedMenus})
}

// REKOMENDASI PERSONAL
// GetPersonalRecommendations mengambil rekomendasi personal berdasarkan riwayat user.
func GetPersonalRecommendations(c *gin.Context) {
	// 1. Ambil user yang sedang login dari context.
	userInterface, _ := c.Get("user")
	user := userInterface.(models.User)

	// 2. Ambil semua menu_id yang pernah di-bookmark atau diberi rating tinggi (>= 4).
	var favoriteMenuIDs []int
	// Query untuk bookmark
	config.DB.Table("user_bookmarks").Where("user_id = ?", user.UserID).Pluck("menu_id", &favoriteMenuIDs)
	// Query untuk rating tinggi
	var highlyRatedMenuIDs []int
	config.DB.Table("menu_ratings").Where("user_id = ? AND rating >= ?", user.UserID, 4).Pluck("menu_id", &highlyRatedMenuIDs)
	
	// Gabungkan semua ID dan hilangkan duplikat.
	favoriteMenuIDs = append(favoriteMenuIDs, highlyRatedMenuIDs...)
	idMap := make(map[int]bool)
	uniqueIDs := []int{}
	for _, id := range favoriteMenuIDs {
		if _, value := idMap[id]; !value {
			idMap[id] = true
			uniqueIDs = append(uniqueIDs, id)
		}
	}

	// Jika user belum punya riwayat, kembalikan daftar kosong.
	if len(uniqueIDs) == 0 {
		c.JSON(http.StatusOK, gin.H{"data": []models.Menu{}})
		return
	}

	// 3. Ubah slice of int menjadi string "101,102,103" untuk dikirim sebagai parameter.
	idsStrSlice := []string{}
	for _, id := range uniqueIDs {
		idsStrSlice = append(idsStrSlice, strconv.Itoa(id))
	}
	idsQueryParam := strings.Join(idsStrSlice, ",")

	// 4. Panggil service Python dengan daftar ID selera pengguna.
	resp, err := http.Get("http://localhost:5000/recommend/profile?ids=" + idsQueryParam)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Gagal menghubungi service rekomendasi"})
		return
	}
	defer resp.Body.Close()

    // ... (sisa kode sama seperti GetRecommendations, untuk memproses hasil)
    body, err := io.ReadAll(resp.Body)
	if err != nil { // ... handle error
	}
	if resp.StatusCode != http.StatusOK { // ... handle error
	}
	var recommendedIDs []int
	json.Unmarshal(body, &recommendedIDs)

	var recommendedMenus []models.Menu
	if len(recommendedIDs) > 0 {
		config.DB.Where("menu_id IN ?", recommendedIDs).Find(&recommendedMenus)
	}

	c.JSON(http.StatusOK, gin.H{"data": recommendedMenus})
}