package web

import (
	"ACT_GO/bngx"
	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	"log"
	"net/http"
	"time"
)

func Start_Web() {
	port := "8080"
	// Starts a new Gin instance with no middle-ware
	router := gin.Default()

	config := cors.DefaultConfig()
	config.AddAllowHeaders("Authorization")
	config.AllowCredentials = true
	config.AllowAllOrigins = false
	// I think you should whitelist a limited origins instead:
	//  config.AllowAllOrigins = []{"xxxx", "xxxx"}
	config.AllowOriginFunc = func(origin string) bool {
		return true
	}
	router.Use(cors.New(config))

	//store := cookie.NewStore([]byte("your-secret-key"))
	//store.Options(sessions.Options{MaxAge: 60 * 60 * 24})
	//router.Use(sessions.Sessions("sessions", store))

	router.Use(cors.New(cors.Config{
		//AllowOrigins:     []string{"http://localhost:5173"},
		AllowMethods:     []string{"GET", "POST", "OPTIONS"},
		AllowHeaders:     []string{"Origin", "Authorization"},
		ExposeHeaders:    []string{"Content-Length"},
		AllowCredentials: true,
		AllowOriginFunc: func(origin string) bool {
			return origin == "http://localhost:5173"
		},
		MaxAge: 12 * time.Hour,
	}))

	// Define handlers
	router.GET("/", func(c *gin.Context) {
		c.String(http.StatusOK, "Hello World!")
	})

	router.GET("/ping", func(c *gin.Context) {
		c.String(http.StatusOK, "pong")
	})

	router.GET("/start-long", func(c *gin.Context) {
		go bngx.Start_Bot1_Long()
		c.String(http.StatusOK, "Started")
	})

	router.GET("/stop-long", func(c *gin.Context) {
		go bngx.Stop_Bot1_Long()
		c.String(http.StatusOK, "Stopped")
	})

	router.GET("/start-short", func(c *gin.Context) {
		go bngx.Start_Bot1_Short()
		c.String(http.StatusOK, "Started")
	})

	router.GET("/stop-short", func(c *gin.Context) {
		go bngx.Stop_Bot1_Short()
		c.String(http.StatusOK, "Stopped")
	})

	router.POST("/login", func(c *gin.Context) {
		var json struct {
			Username string `json:"username"`
			Pwd      string `json:"pwd"`
		}
		if err := c.ShouldBindJSON(&json); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}
		user_ok := auth(json.Username, json.Pwd)
		if user_ok {
			c.JSON(http.StatusOK, gin.H{"success": true})
		} else {
			c.JSON(http.StatusOK, gin.H{"success": false})
		}
	})

	// Listen and serve on defined port
	log.Printf("Listening on port %s", port)
	router.Run(":" + port)
}
