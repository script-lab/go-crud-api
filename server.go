package main

import (
	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"

	"github.com/script-lab/go-crud-api/database"
	"gorm.io/gorm"

	"net/http"
	"time"
)

type (
	Post struct {
		ID        uint   `gorm:"primaryKey" json:"id" param:"id"`
		Title     string `json:"title"`
		Content   string `json:"content"`
		CreatedAt time.Time
		UpdatedAt time.Time
		DeletedAt gorm.DeletedAt `gorm:"index"`
	}
)

func createPost(c echo.Context) error {
	post := Post{}
	if err := c.Bind(&post); err != nil {
		return err
	}
	database.Db.Create(&post)
	return c.JSON(http.StatusCreated, post)
}

func getAllPosts(c echo.Context) error {
	posts := []Post{}
	database.Db.Order("updated_at desc").Find(&posts)
	return c.JSON(http.StatusOK, posts)
}

func getPost(c echo.Context) error {
	post := Post{}
	if err := c.Bind(&post); err != nil {
		return err
	}
	database.Db.Take(&post)
	return c.JSON(http.StatusOK, post)
}

func updatePost(c echo.Context) error {
	post := Post{}
	if err := c.Bind(&post); err != nil {
		return err
	}
	database.Db.Save(&post)
	return c.JSON(http.StatusOK, post)
}

func deletePost(c echo.Context) error {
	post := Post{}
	if err := c.Bind(&post); err != nil {
		return err
	}
	database.Db.Unscoped().Delete(&post)
	return c.NoContent(http.StatusNoContent)
}

func searchPost(c echo.Context) error {
	posts := []Post{}
	key := c.QueryParam("keyword")
	database.Db.Where("title LIKE ?", "%"+key+"%").Find(&posts)
	return c.JSON(http.StatusOK, posts)
}

func main() {

	e := echo.New()

	// cors
	e.Use(middleware.CORSWithConfig(middleware.CORSConfig{
		AllowOrigins: []string{"http://localhost:3000"},
		AllowHeaders: []string{echo.HeaderOrigin, echo.HeaderContentType, echo.HeaderAccept},
	}))

	// Database
	database.Connect()
	if database.Db.Migrator().HasTable(&Post{}) == false {
		database.Db.Migrator().CreateTable(&Post{})
	}
	sqlDB, _ := database.Db.DB()
	defer sqlDB.Close()

	// Middleware
	e.Use(middleware.Logger())
	e.Use(middleware.Recover())

	// Routing
	e.POST("/post", createPost)
	e.GET("/posts", getAllPosts)
	e.GET("/post/:id", getPost)
	e.PUT("/post/:id", updatePost)
	e.DELETE("/post/:id", deletePost)
	e.GET("posts/search", searchPost)

	// Start server
	e.Logger.Fatal(e.Start(":1323"))
}
