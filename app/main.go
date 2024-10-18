package main

import (
	"backend/internal/entity"
	"backend/internal/handler/auth"
	"backend/internal/handler/game"
	"backend/internal/handler/user"
	postgresRepository "backend/internal/infrastructure/repository/postgres"
	"backend/internal/usecase"
	"log"

	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

func main() {
	err := godotenv.Load()
	if err != nil {
		log.Fatal(err)
	}

	dsn := "host=localhost user=postgres password=postgres dbname=postgres port=5432 sslmode=disable TimeZone=Europe/Moscow"
	db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{})
	if err != nil {
		log.Fatal(err)
	}

	db.AutoMigrate(&entity.Skin{}, &entity.Cape{}, &entity.MinecraftCredential{}, &entity.User{}, &entity.RefreshToken{})

	router := gin.Default()

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
	userUseCaseImpl := usecase.NewUserUseCaseImpl(postgresTextureRepository)
	userHandler := user.NewUserHandler(userUseCaseImpl)

	userGroup := router.Group("/user", authMiddleware)
	{
		userGroup.GET("/me", userHandler.Me)
		userGroup.POST("/skin", userHandler.Skin)
		userGroup.POST("/cape", userHandler.Cape)
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
