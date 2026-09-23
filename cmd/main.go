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

	"survey-backend/internal/answer"
	"survey-backend/internal/category"
	"survey-backend/internal/country"
	"survey-backend/internal/question"
	"survey-backend/internal/region"
	"survey-backend/internal/savedsurvey"
	"survey-backend/internal/survey"
	"survey-backend/internal/user"
	"survey-backend/middleware"
	"survey-backend/pkg/config"
	"survey-backend/pkg/database"
	"survey-backend/pkg/jwt"
	"survey-backend/routes"

	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"
)

func main() {
	// A .env file is a local convenience. Hosting platforms inject environment
	// variables directly and ship no such file, so its absence is not an error.
	if err := godotenv.Load(); err != nil {
		log.Println("no .env file loaded, using the process environment")
	}

	if config.Get("APP_ENV", "development") == "production" {
		gin.SetMode(gin.ReleaseMode)
	}

	// Fail fast and name every missing variable at once.
	if err := config.RequireAll("DB_HOST", "DB_USER", "DB_NAME", "JWT_SECRET"); err != nil {
		log.Fatal(err)
	}
	if err := jwt.CheckSecret(); err != nil {
		log.Fatal(err)
	}
	if len(config.Get("JWT_SECRET", "")) < jwt.MinSecretLength {
		log.Printf("warning: JWT_SECRET is shorter than %d characters; use a longer random value", jwt.MinSecretLength)
	}

	// database
	if err := database.ConnectDB(); err != nil {
		log.Fatal(err)
	}
	defer database.Close()

	if shouldRunMigrations() {
		log.Println("applying schema migrations")
		if err := database.DB.AutoMigrate(models()...); err != nil {
			log.Fatal("Failed to migrate database:", err)
		}
	} else {
		log.Println("skipping migrations (set RUN_MIGRATIONS=true to apply schema changes)")

		// Catch the case where the schema was never created, rather than
		// letting seeding fail with a bare "table doesn't exist".
		if !database.DB.Migrator().HasTable(&user.User{}) {
			log.Fatal("database schema is missing; start once with RUN_MIGRATIONS=true to create it")
		}
	}

	// seed
	if err := category.SeedCategories(); err != nil {
		log.Fatal("Failed to seed categories:", err)
	}
	if err := country.SeedCountries(); err != nil {
		log.Fatal("Failed to seed countries:", err)
	}
	if err := region.SeedRegions(); err != nil {
		log.Fatal("Failed to seed regions:", err)
	}
	if err := user.SeedAdmin(); err != nil {
		log.Fatal("Failed to seed admin:", err)
	}

	r := gin.New()
	r.Use(gin.Logger(), gin.Recovery())

	// Only needed for a Flutter web build; an empty list leaves CORS off.
	if origins := config.List("CORS_ALLOWED_ORIGINS"); len(origins) > 0 {
		r.Use(middleware.CORSMiddleware(origins))
		log.Println("CORS enabled for:", origins)
	}

	// Platform health checks hit this before routing any traffic here, so it
	// must not touch the database.
	r.GET("/health", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{"status": "ok"})
	})

	// public routes - no token required
	routes.PublicRoutes(r.Group("/"))

	// protected routes
	protected := r.Group("/")
	protected.Use(middleware.AuthMiddleware())

	// Attach role-based groups
	routes.UserRoutes(protected)
	routes.RespondentRoutes(protected)
	routes.InterviewerRoutes(protected)
	routes.AdminRoutes(protected)

	// Every platform assigns the port through PORT; the fallback is for local runs.
	port := config.Get("PORT", "8080")

	server := &http.Server{
		Addr:              ":" + port,
		Handler:           r,
		ReadHeaderTimeout: 10 * time.Second,
	}

	go func() {
		log.Println("server listening on port", port)
		if err := server.ListenAndServe(); err != nil && !errors.Is(err, http.ErrServerClosed) {
			log.Fatal("server failed: ", err)
		}
	}()

	// Container platforms stop an instance with SIGTERM; drain in-flight
	// requests instead of cutting them off.
	stop := make(chan os.Signal, 1)
	signal.Notify(stop, os.Interrupt, syscall.SIGTERM)
	<-stop

	log.Println("shutting down")
	ctx, cancel := context.WithTimeout(context.Background(), 15*time.Second)
	defer cancel()

	if err := server.Shutdown(ctx); err != nil {
		log.Println("forced shutdown:", err)
	}
}

// models lists every table AutoMigrate manages.
func models() []any {
	return []any{
		&user.User{},
		&category.Category{},
		&country.Country{},
		&region.Region{},
		&survey.Survey{},
		&question.Question{},
		&question.Option{},
		&answer.Answer{},
		&answer.UserAnswer{},
		&savedsurvey.SavedSurvey{},
	}
}

// shouldRunMigrations defaults to on in development and off in production.
// AutoMigrate costs roughly 180 schema-introspection queries per boot, which is
// wasted latency and quota on a hosted database, and a restart is not the right
// moment to alter a live schema. Set RUN_MIGRATIONS explicitly to override.
func shouldRunMigrations() bool {
	return config.Bool("RUN_MIGRATIONS", config.Get("APP_ENV", "development") != "production")
}
