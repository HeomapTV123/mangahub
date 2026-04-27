package main

import (
	"log"

	"mangahub/internal/auth"
	mangaFeature "mangahub/internal/manga"
	"mangahub/internal/tcp"
	userFeature "mangahub/internal/user"
	"mangahub/pkg/database"

	"github.com/gin-gonic/gin"
)

func main() {
	db, err := database.InitSQLite("data/mangahub.db")
	if err != nil {
		log.Fatal("failed to connect database: ", err)
	}
	defer db.Close()

	empty, err := database.IsMangaTableEmpty(db)
	if err != nil {
		log.Fatal(err)
	}

	if empty {
		log.Println("Seeding manga data...")

		mangaList, err := database.LoadMangaJSON("data/manga.json")
		if err != nil {
			log.Fatal("failed to load manga json: ", err)
		}

		if err := database.SeedManga(db, mangaList); err != nil {
			log.Fatal("failed to seed manga: ", err)
		}
	}

	authService := auth.NewService(db)
	authHandler := auth.NewHandler(authService)

	mangaRepo := mangaFeature.NewRepository(db)
	mangaService := mangaFeature.NewService(mangaRepo)
	mangaHandler := mangaFeature.NewHandler(mangaService)

	userRepo := userFeature.NewRepository(db)
	userService := userFeature.NewService(userRepo)
	userHandler := userFeature.NewHandler(userService)

	router := gin.Default()

	router.GET("/", func(c *gin.Context) {
		c.JSON(200, gin.H{
			"message": "MangaHub API Server is running",
		})
	})

	router.POST("/auth/register", authHandler.Register)
	router.POST("/auth/login", authHandler.Login)

	router.GET("/manga", mangaHandler.GetAll)
	router.GET("/manga/:id", mangaHandler.GetByID)

	protected := router.Group("/")
	protected.Use(auth.AuthMiddleware())

	protected.POST("/users/library", userHandler.AddToLibrary)
	protected.GET("/users/library", userHandler.GetLibrary)
	protected.PUT("/users/progress", userHandler.UpdateProgress)

	protected.POST("/manga", mangaHandler.Create)
	protected.PUT("/manga/:id", mangaHandler.Update)
	protected.DELETE("/manga/:id", mangaHandler.Delete)

	protected.POST("/admin/import/mangadex", mangaHandler.ImportFromMangaDex)

	go tcp.StartTCPServer(":9090", db)
	// log.Println("TCP Progress Sync Server running on localhost:9090")

	log.Println("API server running on http://localhost:8080")
	if err := router.Run(":8080"); err != nil {
		log.Fatal("failed to start server: ", err)
	}
}
