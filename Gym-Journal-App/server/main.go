package main

import (
	"os"

	"github.com/alaiy95/go-gym-joural-app/routes"
	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
)

// This code reads the server port from the environment variable PORT. If it is not set, it defaults to 8000.

func main() {
	port := os.Getenv("PORT")
	if port == "" {
		port = "8000"
	}

	router := gin.New()        //  creates a new Gin router instance, which handles incoming requests and routes them to the handlers.
	router.Use(gin.Logger())   // middleware that logs incoming requests and responses to the console.
	router.Use(cors.Default()) // middleware that enables Cross-Origin Resource Sharing (CORS) for the server, which allows web pages from other domains to make requests to the server.

	router.POST("/entry/create", routes.AddEntry)
	router.GET("/entries", routes.GetEntries)
	router.GET("/entry/:id/", routes.GetEntryById)
	router.GET("/exercise/:exercise", routes.GetEntriesByExercise)

	router.PUT("/entry/update/:id", routes.UpdateEntry)
	router.PUT("/exercise/update/:id", routes.UpdateExercise)
	router.DELETE("/entry/delete/:id", routes.DeleteEntry)
	router.Run(":" + port)
}
