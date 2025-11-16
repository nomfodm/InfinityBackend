package main

import (
	"context"
	"errors"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/nomfodm/InfinityBackend/internal/entity"
	"github.com/nomfodm/InfinityBackend/internal/handler/auth"
	"github.com/nomfodm/InfinityBackend/internal/handler/game"
	"github.com/nomfodm/InfinityBackend/internal/handler/healthstate"
	"github.com/nomfodm/InfinityBackend/internal/handler/launcher"
	"github.com/nomfodm/InfinityBackend/internal/handler/user"
	postgresRepository "github.com/nomfodm/InfinityBackend/internal/infrastructure/repository/postgres"
	"github.com/nomfodm/InfinityBackend/internal/usecase"

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

	err = db.AutoMigrate(&entity.Skin{}, &entity.Cape{}, &entity.MinecraftCredential{}, &entity.User{}, &entity.RefreshToken{}, &entity.LauncherVersion{}, &entity.HealthState{})
	if err != nil {
		log.Fatal(err)
	}

	router := gin.Default()

	router.Use(CORSMiddleware())

	postgresHealthStateRepository := postgresRepository.NewPostgresHealthStateRepository(db)
	healthStateUseCaseImpl := usecase.NewHealthStateUseCaseImpl(postgresHealthStateRepository)
	healthStateHander := healthstate.NewHealthStateHandler(healthStateUseCaseImpl)

	healthStateMiddleware := healthstate.NewHealthStateMiddleware(healthStateUseCaseImpl)

	err = healthStateUseCaseImpl.InitHealthState()
	if err != nil {
		log.Fatal(err)
	}

	router.GET("/checkConnection", healthStateMiddleware, healthStateHander.Index)
	router.GET("/health", healthStateMiddleware, healthStateHander.Index)

	postgresUserRepository := postgresRepository.NewPostgresUserRepository(db)
	authUseCaseImpl := usecase.NewAuthUseCaseImpl(postgresUserRepository)
	authHandler := auth.NewAuthHandler(authUseCaseImpl)
	authGroup := router.Group("/auth", healthStateMiddleware)

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

	userGroup := router.Group("/user", healthStateMiddleware, authMiddleware)
	{
		userGroup.GET("/me", userHandler.Me)
		userGroup.POST("/skin", userHandler.Skin)
		userGroup.POST("/cape", userHandler.Cape)
		userGroup.GET("/nickname", userHandler.Nickname)
	}

	postgresGameRepository := postgresRepository.NewPostgresGameRepository(db)
	gameUseCaseImpl := usecase.NewGameUseCaseImpl(postgresGameRepository, postgresTextureRepository)
	gameHandler := game.NewGameHandler(gameUseCaseImpl)

	gameGroup := router.Group("/game", healthStateMiddleware)
	{
		gameGroup.GET("/launcher", authMiddleware, gameHandler.Launcher)
		gameGroup.POST("/join", gameHandler.Join)
		gameGroup.GET("/profile/:uuid", gameHandler.Profile)
		gameGroup.GET("/hasJoined", gameHandler.HasJoined)
	}

	launcherRepository := postgresRepository.NewPostgresLauncherRepository(db)
	launcherUseCaseImpl := usecase.NewLauncherUseCaseImpl(launcherRepository)
	launcherHandler := launcher.NewLauncherHandler(launcherUseCaseImpl)

	adminAccessMiddleware := launcher.NewAdminAccessMiddleware()
	router.GET("/health/set", adminAccessMiddleware, healthStateHander.SetStatus)

	launcherGroup := router.Group("/launcher", healthStateMiddleware)
	{
		launcherGroup.GET("/download", launcherHandler.DownloadLauncher)
		launcherGroup.GET("/updates", func(c *gin.Context) {
			c.JSON(302, gin.H{
				"error":  "Page Moved",
				"detail": "Была переработана система версий лаунчера, скачайте актуальную версию вручную: infinityserver.ru/launcher",
			})
		})
	}

	updateGroup := router.Group("/launcher/update", healthStateMiddleware)
	{
		updateGroup.GET("/actual", launcherHandler.ActualVersion)
		updateGroup.POST("/register", adminAccessMiddleware, launcherHandler.RegisterUpdate)
		updateGroup.GET("/lastmandatory", launcherHandler.LastMandatory)
	}

	srv := &http.Server{
		Addr:    ":8000",
		Handler: router.Handler(),
	}

	go func() {
		if err := srv.ListenAndServe(); err != nil && !errors.Is(err, http.ErrServerClosed) {
			log.Fatalf("listen: %s\n", err)
		}
	}()

	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit
	log.Println("Shutdown Server ...")

	ctx, cancel := context.WithTimeout(context.Background(), 2*time.Second)
	defer cancel()
	if err := srv.Shutdown(ctx); err != nil {
		log.Println("Server Shutdown:", err)
	}

	<-ctx.Done()
	log.Println("timeout of 2 seconds.")
	log.Println("Server exiting")
}
