// main.go
package main

import (
	"creacipe-backend/config"
	"creacipe-backend/controllers"
	"creacipe-backend/middlewares"
	"log"

	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"
)

// init() berjalan sekali secara otomatis sebelum main()
func init() {
	// Memuat konfigurasi dari file .env
	err := godotenv.Load()
	if err != nil {
		log.Fatal("Gagal memuat file .env")
	}
	// Membuat koneksi ke database
	config.ConnectDatabase()
}

func main() {
	// Inisialisasi Gin router
	r := gin.Default()

	// Menerapkan middleware CORS untuk semua rute
	r.Use(config.CORSMiddleware())

	// Grup utama untuk semua rute API dengan prefix /api
	api := r.Group("/api")
	{
		// --- Rute Publik (Tidak Perlu Login) ---
		api.POST("/setup/first-admin", controllers.SetupFirstAdmin)
		api.POST("/register", controllers.Register)
		api.POST("/login", controllers.Login)
		api.GET("/tags", controllers.GetAllTags)
		api.GET("/menus", controllers.GetAllMenus)
		api.GET("/menus/popular", controllers.GetPopularMenus)
		api.GET("/menus/:id", controllers.GetMenuByID)
		api.GET("/menus/:id/recommendations", controllers.GetRecommendations)

		//---untuk reset password---//
		api.POST("/forgot-password", controllers.ForgotPassword)
		api.POST("/reset-password", controllers.ResetPassword)

		// --- Rute Terotentikasi (Harus Login) ---
		// Grup 'auth' ini dilindungi oleh middleware RequireAuth.
		auth := api.Group("/")
		auth.Use(middlewares.RequireAuth)
		{
			// Rute untuk Member (CRUD & Interaksi)
			auth.POST("/menus", controllers.CreateMenu)
			auth.PUT("/menus/:id", controllers.UpdateMenu)
			auth.DELETE("/menus/:id", controllers.DeleteMenu)
			
			// Rute Interaksi (dipanggil dari interaction_controller)
			auth.POST("/menus/:id/tags", controllers.AddTagToMenu)
			auth.POST("/menus/:id/bookmark", controllers.BookmarkMenu)
			auth.DELETE("/menus/:id/bookmark", controllers.UnbookmarkMenu)
			auth.POST("/menus/:id/like", controllers.LikeMenu)
			auth.POST("/menus/:id/dislike", controllers.DislikeMenu)
			
			// Rute Rekomendasi Personal (dipanggil dari recommendation_controller)
			auth.GET("/me/recommendations", controllers.GetPersonalRecommendations)

			// rute profil pengguna
			auth.GET("/me", controllers.GetMyProfile)
			auth.PUT("/me", controllers.UpdateMyProfile)

			// Rute Moderasi (Editor & Admin)
			// Grup ini dilindungi oleh middleware AuthorizeRole.
			editorRoutes := auth.Group("/editor")
			editorRoutes.Use(middlewares.AuthorizeRole("admin", "editor"))
			{
				editorRoutes.GET("/menus", controllers.GetAllMenusForModeration)
				editorRoutes.GET("/menus/pending", controllers.GetPendingMenus)
				editorRoutes.PATCH("/menus/:id/status", controllers.UpdateMenuStatus)
				// Rute Manajemen Tag
				editorRoutes.POST("/tags", controllers.CreateTag)
                editorRoutes.PUT("/tags/:id", controllers.UpdateTag)
                editorRoutes.DELETE("/tags/:id", controllers.DeleteTag)
			}

			// Rute Manajemen User (Hanya Admin)
			adminRoutes := auth.Group("/admin")
			adminRoutes.Use(middlewares.AuthorizeRole("admin"))
			{
				adminRoutes.GET("/users", controllers.GetAllUsers)
				adminRoutes.PUT("/users/:id", controllers.UpdateUser)
				adminRoutes.PATCH("/users/:id/role", controllers.UpdateUserRole)
				adminRoutes.DELETE("/users/:id", controllers.DeleteUser)
			}
		}
	}

	// Menjalankan server pada port 8080
	log.Println("Server dimulai di http://localhost:8080")
	r.Run()
}