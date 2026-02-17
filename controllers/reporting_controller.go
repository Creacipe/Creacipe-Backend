// controllers/reporting_controller.go
package controllers

import (
	"creacipe-backend/config"
	"creacipe-backend/models"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
)


func parseDateRange(c *gin.Context) (time.Time, time.Time) {
	
	endDate := time.Now()
	startDate := endDate.AddDate(0, 0, -29) 

	
	if val, ok := c.GetQuery("startDate"); ok {
		if parsed, err := time.Parse("2006-01-02", val); err == nil {
			startDate = parsed
		}
	}
	if val, ok := c.GetQuery("endDate"); ok {
		if parsed, err := time.Parse("2006-01-02", val); err == nil {
			
			endDate = parsed.Add(23*time.Hour + 59*time.Minute + 59*time.Second)
		}
	}

	return startDate, endDate
}



func GetRecipeStats(c *gin.Context) {
	var pendingCount, approvedCount, rejectedCount int64

	
	config.DB.Model(&models.Menu{}).Where("status = ?", "pending").Count(&pendingCount)
	config.DB.Model(&models.Menu{}).Where("status = ?", "approved").Count(&approvedCount)
	config.DB.Model(&models.Menu{}).Where("status = ?", "rejected").Count(&rejectedCount)

	stats := gin.H{
		"pending":  pendingCount,
		"approved": approvedCount,
		"rejected": rejectedCount,
	}

	c.JSON(http.StatusOK, gin.H{"data": stats})
}

func GetGrowthStats(c *gin.Context) {
	startDate, endDate := parseDateRange(c)

	type DailyStat struct {
		Date       string `json:"date"`
		NewUsers   int    `json:"new_users"`
		NewRecipes int    `json:"new_recipes"`
	}

	statsMap := make(map[string]*DailyStat)
	for d := startDate; !d.After(endDate); d = d.AddDate(0, 0, 1) {
		dateStr := d.Format("2006-01-02")
		statsMap[dateStr] = &DailyStat{Date: dateStr}
	}

	
	var userGrowth []struct {
		Date  time.Time
		Count int
	}
	config.DB.Model(&models.User{}).
		Select("DATE(created_at) as date, COUNT(user_id) as count").
		Where("created_at BETWEEN ? AND ?", startDate, endDate).
		Group("date").
		Scan(&userGrowth)

	for _, stat := range userGrowth {
		dateStr := stat.Date.Format("2006-01-02")
		if entry, ok := statsMap[dateStr]; ok {
			entry.NewUsers = stat.Count
		}
	}


	var recipeGrowth []struct {
		Date  time.Time
		Count int
	}
	config.DB.Model(&models.Menu{}).
		Select("DATE(created_at) as date, COUNT(menu_id) as count").
		Where("created_at BETWEEN ? AND ?", startDate, endDate).
		Group("date").
		Scan(&recipeGrowth)

	for _, stat := range recipeGrowth {
		dateStr := stat.Date.Format("2006-01-02")
		if entry, ok := statsMap[dateStr]; ok {
			entry.NewRecipes = stat.Count
		}
	}


	var results []DailyStat
	for d := startDate; !d.After(endDate); d = d.AddDate(0, 0, 1) {
		dateStr := d.Format("2006-01-02")
		results = append(results, *statsMap[dateStr])
	}

	c.JSON(http.StatusOK, gin.H{"data": results})
}



func GetTopTags(c *gin.Context) {
	type TagStat struct {
		TagName string `json:"tag_name"`
		Count   int    `json:"count"`
	}

	var results []TagStat

	err := config.DB.Table("tags").
		Select("tags.tag_name, COUNT(menu_tags.menu_id) as count").
		Joins("join menu_tags on tags.tag_id = menu_tags.tag_id").
		Group("tags.tag_name").
		Order("count DESC").
		Limit(10).
		Scan(&results).Error

	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to retrieve top tags"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"data": results})
}


func GetActivityLogStats(c *gin.Context) {
	startDate, endDate := parseDateRange(c)

	type DailyActivity struct {
		Date          string `json:"date"`
		ActivityCount int    `json:"activity_count"`
	}
	
	
	statsMap := make(map[string]*DailyActivity)
	for d := startDate; !d.After(endDate); d = d.AddDate(0, 0, 1) {
		dateStr := d.Format("2006-01-02")
		statsMap[dateStr] = &DailyActivity{Date: dateStr}
	}

	var dbResults []struct {
		Date time.Time
		ActivityCount int
	}

	err := config.DB.Model(&models.LogActivity{}).
		Select("DATE(created_at) as date, COUNT(activity_id) as activity_count").
		Where("created_at BETWEEN ? AND ?", startDate, endDate).
		Group("date").
		Order("date ASC").
		Scan(&dbResults).Error

	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to retrieve activity log stats"})
		return
	}

	for _, res := range dbResults {
		dateStr := res.Date.Format("2006-01-02")
		if entry, ok := statsMap[dateStr]; ok {
			entry.ActivityCount = res.ActivityCount
		}
	}
	
	
	var results []DailyActivity
	for d := startDate; !d.After(endDate); d = d.AddDate(0, 0, 1) {
		dateStr := d.Format("2006-01-02")
		results = append(results, *statsMap[dateStr])
	}

	c.JSON(http.StatusOK, gin.H{"data": results})
}



func GetTopLikedRecipes(c *gin.Context) {
	type RecipeStat struct {
		Title      string `json:"title"`
		TotalLikes int    `json:"total_likes"`
	}

	var results []RecipeStat

	err := config.DB.Table("menus").
		Select("menus.title, SUM(menu_votes.likes_count) as total_likes").
		Joins("join menu_votes on menus.menu_id = menu_votes.menu_id").
		Group("menus.menu_id").
		Order("total_likes DESC").
		Limit(5).
		Scan(&results).Error

	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to retrieve top liked recipes"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"data": results})
}


func GetTopBookmarkedRecipes(c *gin.Context) {
	type RecipeStat struct {
		Title          string `json:"title"`
		TotalBookmarks int    `json:"total_bookmarks"`
	}

	var results []RecipeStat

	err := config.DB.Table("menus").
		Select("menus.title, COUNT(user_bookmarks.user_id) as total_bookmarks").
		Joins("join user_bookmarks on menus.menu_id = user_bookmarks.menu_id").
		Group("menus.menu_id").
		Order("total_bookmarks DESC").
		Limit(5).
		Scan(&results).Error

	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to retrieve top bookmarked recipes"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"data": results})
}