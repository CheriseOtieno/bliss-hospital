package main

import (
	"fmt"
	"log"
	"os"

	"bliss-backend/db"
	"bliss-backend/handlers"
	"bliss-backend/middleware"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"
)

func main() {

	// Load env
	if err := godotenv.Load(); err != nil {
		log.Println("No .env file found — using system environment variables")
	}

	// DB connection
	db.Connect()
	defer db.Close()

	// Gin mode
	if os.Getenv("ENV") == "production" {
		gin.SetMode(gin.ReleaseMode)
	}

	r := gin.Default()

	// CORS
	r.Use(cors.New(cors.Config{
		AllowOrigins: []string{
			"http://localhost:5173",
			"http://localhost:5174",
			"http://localhost:5175",
			"http://localhost:3000",
		},
		AllowMethods: []string{
			"GET", "POST", "PUT", "PATCH", "DELETE", "OPTIONS",
		},
		AllowHeaders: []string{
			"Origin",
			"Content-Type",
			"Authorization",
		},
		ExposeHeaders:    []string{"Content-Length"},
		AllowCredentials: true,
	}))

	r.OPTIONS("/*path", func(c *gin.Context) {
		c.AbortWithStatus(204)
	})

	// PUBLIC ROUTES
	public := r.Group("/api")
	{
		public.POST("/auth/register", handlers.Register)
		public.POST("/auth/login", handlers.Login)

		// ✅ FIX ADDED: branches endpoint
		public.GET("/branches", handlers.GetBranches)

		public.GET("/departments", handlers.GetDepartments)
		public.GET("/doctors", handlers.GetDoctors)
		public.GET("/slots", handlers.GetSlots)

		
	}

	// PROTECTED ROUTES
	protected := r.Group("/api")
	protected.Use(middleware.AuthRequired())
	{
		protected.GET("/auth/me", handlers.Me)

		protected.POST("/appointments", handlers.CreateAppointment)
		protected.GET("/appointments", handlers.GetAppointments)
		protected.GET("/appointments/:id", handlers.GetAppointment)
		protected.PATCH("/appointments/:id", handlers.UpdateAppointment)

		protected.POST("/queue/checkin", handlers.CheckIn)
		protected.GET("/queue/my", handlers.GetMyQueueStatus)
		protected.GET("/notifications", handlers.GetNotifications)
	}

	// STAFF ROUTES
	staff := r.Group("/api")
	staff.Use(middleware.AuthRequired(), middleware.StaffOnly())
	{
		staff.GET("/queue", handlers.GetQueue)
		staff.PATCH("/queue/:id", handlers.UpdateQueueEntry)
		staff.GET("/queue/stats", handlers.GetQueueStats)
	}

	// ADMIN ROUTES
	admin := r.Group("/api/admin")
	admin.Use(middleware.AuthRequired(), middleware.AdminOnly())
	{
		admin.GET("/dashboard", handlers.GetDashboardStats)
		admin.GET("/users", handlers.GetUsers)
		admin.POST("/doctors", handlers.CreateDoctor)
		admin.PATCH("/doctors/:id", handlers.UpdateDoctor)
	}

	// Health check
	r.GET("/health", func(c *gin.Context) {
		c.JSON(200, gin.H{
			"status":  "ok",
			"service": "Bliss Hospital API",
		})
	})

	// Start server
	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}

	fmt.Printf("\n🏥 Bliss Hospital API running on http://localhost:%s\n\n", port)

	if err := r.Run(":" + port); err != nil {
		log.Fatalf("Failed to start server: %v", err)
	}
}