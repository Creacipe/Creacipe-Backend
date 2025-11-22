// main.go
package main

import (
	"creacipe-backend/config"
	"creacipe-backend/controllers"
	"creacipe-backend/middlewares"
	"log"
	"os"

	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"
)

func init() {
	err := godotenv.Load()
	if err != nil {
		log.Fatal("Gagal memuat file .env")
	}
	config.ConnectDatabase()
}

func main() {
	r := gin.Default()

	// --- 2. TAMBAHKAN DUA BARIS INI ---
	// Pastikan folder 'assets' ada
	os.MkdirAll("./assets", os.ModePerm) 
	// Sajikan folder 'assets' di URL '/assets'
	r.Static("/assets", "./assets") 
	// ---------------------------------

	r.Use(config.CORSMiddleware())

	api := r.Group("/api")
	{
		// --- Rute Publik ---
		api.POST("/setup/first-admin", controllers.SetupFirstAdmin)
		api.POST("/register", controllers.Register)
		api.POST("/login", controllers.Login)
		api.POST("/forgot-password", controllers.ForgotPasswordRequest)
		api.POST("/forgot-password/verify", controllers.ForgotPasswordVerify)
		api.GET("/tags", controllers.GetAllTags)
		api.GET("/menus", controllers.GetAllMenus)
		api.GET("/menus/popular", controllers.GetPopularMenus)
		api.GET("/menus/:id", controllers.GetMenuByID)
		api.GET("/menus/:id/recommendations", controllers.GetRecommendations)
		api.GET("/categories", controllers.GetAllCategories)
		api.POST("/auth/refresh", controllers.RefreshToken)

		// --- Rute Terotentikasi ---
		auth := api.Group("/")
		auth.Use(middlewares.RequireAuth)
		{
			// Rute Member
			auth.POST("/menus", controllers.CreateMenu)
			auth.PUT("/menus/:id", controllers.UpdateMenu)
			auth.DELETE("/menus/:id", controllers.DeleteMenu)
			auth.POST("/menus/:id/tags", controllers.AddTagToMenu)
			auth.POST("/menus/:id/bookmark", controllers.BookmarkMenu)
			auth.DELETE("/menus/:id/bookmark", controllers.UnbookmarkMenu)
			auth.POST("/menus/:id/like", controllers.LikeMenu)
			auth.POST("/menus/:id/dislike", controllers.DislikeMenu)
			auth.GET("/menus/:id/interaction-status", controllers.GetUserInteractionStatus)
			auth.GET("/me/recommendations", controllers.GetPersonalRecommendations)
			auth.GET("/me", controllers.GetMyProfile)
			auth.PUT("/me", controllers.UpdateMyProfile)
			auth.GET("/me/menus", controllers.GetMyMenus)
			auth.GET("/me/bookmarks", controllers.GetMyBookmarks)
			auth.GET("/me/collection", controllers.GetMyCollection)

			// Rute Reset Password & Email Change
			auth.POST("/me/request-password-reset", controllers.RequestPasswordReset)
			auth.POST("/me/verify-reset-password", controllers.VerifyAndResetPassword)
			auth.POST("/me/request-email-change", controllers.RequestEmailChange)
			auth.POST("/me/verify-email-change", controllers.VerifyAndChangeEmail)

			// Notifications
			auth.GET("/me/notifications", controllers.GetMyNotifications)
			auth.GET("/me/notifications/unread-count", controllers.GetUnreadNotificationCount)
			auth.PATCH("/me/notifications/:id/read", controllers.MarkNotificationAsRead)
			auth.PATCH("/me/notifications/mark-all-read", controllers.MarkAllNotificationsAsRead)

			// Comments
			auth.POST("/menus/:id/comments", controllers.CreateComment)
			auth.GET("/menus/:id/comments", controllers.GetCommentsByMenu)
			auth.DELETE("/comments/:id", controllers.DeleteComment)

			// Rute Moderasi (Editor & Admin)
			editorRoutes := auth.Group("/editor")
			editorRoutes.Use(middlewares.AuthorizeRole("admin", "editor"))
			{
				// Dashboard Stats
				editorRoutes.GET("/dashboard/stats", controllers.GetDashboardStats)

				// manajement resep
				editorRoutes.GET("/menus", controllers.GetAllMenusForModeration)
				editorRoutes.GET("/menus/pending", controllers.GetPendingMenus)
				editorRoutes.PATCH("/menus/:id/status", controllers.UpdateMenuStatus)
				editorRoutes.POST("/tags", controllers.CreateTag)
        editorRoutes.PUT("/tags/:id", controllers.UpdateTag)
        editorRoutes.DELETE("/tags/:id", controllers.DeleteTag)
				editorRoutes.DELETE("/menus/:id", controllers.DeleteMenu)

				// --- manajemen kategori ---
        editorRoutes.POST("/categories", controllers.CreateCategory)
        editorRoutes.PUT("/categories/:id", controllers.UpdateCategory)
        editorRoutes.DELETE("/categories/:id", controllers.DeleteCategory)
        // ------------------------------------------
			}

			// Rute Manajemen User (Hanya Admin)
			adminRoutes := auth.Group("/admin")
			adminRoutes.Use(middlewares.AuthorizeRole("admin"))
			{
				// Dashboard Stats (Admin dapat akses semua stats termasuk user)
				adminRoutes.GET("/dashboard/stats", controllers.GetDashboardStats)
				// Manajemen User
				adminRoutes.POST("/users", controllers.AdminCreateUser)
				adminRoutes.GET("/users", controllers.GetAllUsers)
				adminRoutes.GET("/users/:id", controllers.GetUserByID)
				adminRoutes.PUT("/users/:id", controllers.UpdateUser)
				adminRoutes.PATCH("/users/:id/role", controllers.UpdateUserRole)
				adminRoutes.PATCH("/users/:id/deactivate", controllers.DeactivateUser)
				adminRoutes.PATCH("/users/:id/activate", controllers.ActivateUser)
				adminRoutes.DELETE("/users/:id", controllers.DeleteUser)
				// Log Aktivitas
				adminRoutes.GET("/logs", controllers.GetActivityLogs)
				// Roles
				adminRoutes.GET("/roles", controllers.GetAllRoles)
			}
		}
	}

	log.Println("Server dimulai di http://localhost:8080")
	r.Run()
}