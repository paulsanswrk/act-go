package web

import (
	"ACT_GO/bngx"
	"github.com/gin-gonic/gin"
	"log"
	"net/http"
)

func Start_Web() {
	port := "8080"
	// Starts a new Gin instance with no middle-ware
	r := gin.New()

	// Define handlers
	r.GET("/", func(c *gin.Context) {
		c.String(http.StatusOK, "Hello World!")
	})

	r.GET("/ping", func(c *gin.Context) {
		c.String(http.StatusOK, "pong")
	})

	r.GET("/start-long", func(c *gin.Context) {
		go bngx.Start_Bot1_Long()
		c.String(http.StatusOK, "Started")
	})

	r.GET("/stop-long", func(c *gin.Context) {
		go bngx.Stop_Bot1_Long()
		c.String(http.StatusOK, "Stopped")
	})

	r.GET("/start-short", func(c *gin.Context) {
		go bngx.Start_Bot1_Short()
		c.String(http.StatusOK, "Started")
	})

	r.GET("/stop-short", func(c *gin.Context) {
		go bngx.Stop_Bot1_Short()
		c.String(http.StatusOK, "Stopped")
	})

	// Listen and serve on defined port
	log.Printf("Listening on port %s", port)
	r.Run(":" + port)
}
