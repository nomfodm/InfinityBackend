package main

import (
	"backend/internal/entity"
	"backend/internal/handler/auth"
	"backend/internal/handler/game"
	"backend/internal/handler/user"
	postgresRepository "backend/internal/infrastructure/repository/postgres"
	"backend/internal/usecase"
	"log"
	"os"

	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

func CORSMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		c.Writer.Header().Set("Access-Control-Allow-Origin", os.Getenv("FRONTEND_URL"))
		c.Writer.Header().Set("Access-Control-Allow-Credentials", "true")
		c.Writer.Header().Set("Access-Control-Allow-Headers", "Content-Type, Content-Length, Accept-Encoding, X-CSRF-Token, Authorization, accept, origin, Cache-Control, X-Requested-With")
		c.Writer.Header().Set("Access-Control-Allow-Methods", "POST, OPTIONS, GET, PUT")

		if c.Request.Method == "OPTIONS" {
			c.AbortWithStatus(204)
			return
		}

		c.Next()
	}
}

func main() {
	err := godotenv.Load()
	if err != nil {
		log.Print(err)
	}

	dsn := os.Getenv("POSTGRES_DSN")
	db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{})
	if err != nil {
		log.Fatal(err)
	}

	db.AutoMigrate(&entity.Skin{}, &entity.Cape{}, &entity.MinecraftCredential{}, &entity.User{}, &entity.RefreshToken{})

	router := gin.Default()

	router.Use(CORSMiddleware())

	postgresUserRepository := postgresRepository.NewPostgresUserRepository(db)
	authUseCaseImpl := usecase.NewAuthUseCaseImpl(postgresUserRepository)
	authHandler := auth.NewAuthHandler(authUseCaseImpl)
	authGroup := router.Group("/auth")
	{
		authGroup.POST("/signup", authHandler.SignUp)
		authGroup.POST("/signin", authHandler.SignIn)
		authGroup.POST("/logout", authHandler.Logout)
		authGroup.POST("/refresh", authHandler.Refresh)
		authGroup.GET("/activate", authHandler.Activate)
	}

	authMiddleware := auth.NewAuthMiddleware(authUseCaseImpl)

	postgresTextureRepository := postgresRepository.NewPostgresTextureRepository(db)
	userUseCaseImpl := usecase.NewUserUseCaseImpl(postgresTextureRepository, postgresUserRepository)
	userHandler := user.NewUserHandler(userUseCaseImpl)

	userGroup := router.Group("/user", authMiddleware)
	{
		userGroup.GET("/me", userHandler.Me)
		userGroup.POST("/skin", userHandler.Skin)
		userGroup.POST("/cape", userHandler.Cape)
		userGroup.GET("/nickname", userHandler.Nickname)
	}

	postgresGameRepository := postgresRepository.NewPostgresGameRepository(db)
	gameUseCaseImpl := usecase.NewGameUseCaseImpl(postgresGameRepository, postgresTextureRepository)
	gameHandler := game.NewGameHandler(gameUseCaseImpl)

	gameGroup := router.Group("/game")
	{
		gameGroup.GET("/launcher", authMiddleware, gameHandler.Launcher)
		gameGroup.POST("/join", gameHandler.Join)
		gameGroup.GET("/profile/:uuid", gameHandler.Profile)
		gameGroup.GET("/hasJoined", gameHandler.HasJoined)
	}

	router.Run(":8000")
}
