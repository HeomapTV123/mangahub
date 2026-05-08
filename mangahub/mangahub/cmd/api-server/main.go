package main

import (
	"log"

	"mangahub/internal/auth"
	grpc "mangahub/internal/grpc"
	mangaFeature "mangahub/internal/manga"
	"mangahub/internal/tcp"
	"mangahub/internal/udp"
	userFeature "mangahub/internal/user"
	ws "mangahub/internal/websocket"
	"mangahub/pkg/database"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	"github.com/gorilla/websocket"
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
	udpServer := &udp.NotificationServer{
		Port: ":9091",
	}
	authService := auth.NewService(db)
	authHandler := auth.NewHandler(authService)

	mangaRepo := mangaFeature.NewRepository(db)
	mangaService := mangaFeature.NewService(mangaRepo)
	mangaHandler := mangaFeature.NewHandler(mangaService, udpServer)

	userRepo := userFeature.NewRepository(db)
	userService := userFeature.NewService(userRepo)
	userHandler := userFeature.NewHandler(userService)
	hub := &ws.ChatHub{
		Clients:    make(map[*websocket.Conn]string),
		Broadcast:  make(chan ws.ChatMessage),
		Register:   make(chan ws.ClientConnection),
		Unregister: make(chan *websocket.Conn),
	}
	router := gin.Default()

	router.Use(cors.New(cors.Config{
		AllowOrigins: []string{
			"http://localhost:3000", // React (CRA)
			"http://localhost:5173", // Vite
			"http://127.0.0.1:5173",
		},
		AllowMethods: []string{
			"GET", "POST", "PUT", "DELETE", "OPTIONS",
		},
		AllowHeaders: []string{
			"Origin", "Content-Type", "Authorization",
		},
		ExposeHeaders: []string{
			"Content-Length",
		},
		AllowCredentials: true,
	}))

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
	go hub.Run()
	router.GET("/ws", ws.ServeWS(hub))
	go tcp.StartTCPServer(":9090", db)
	go udp.StartUDPServer(udpServer)
	log.Println("[gRPC] starting from api-server...")
	go grpc.StartGRPCServer()
	// log.Println("TCP Progress Sync Server running on localhost:9090")

	log.Println("API server running on http://localhost:8080")
	if err := router.Run(":8080"); err != nil {
		log.Fatal("failed to start server: ", err)
	}
}
