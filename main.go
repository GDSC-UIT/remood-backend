package main

import (
	"context"
	"fmt"
	"os"

	"remood/routes"

	"github.com/gin-gonic/gin"
	_ "github.com/joho/godotenv/autoload"

	"remood/pkg/database"
)

func main() {
	//DATABASE CONNECTION
	db := database.GetMongoInstance()
	defer db.Client.Connect(context.Background())
	fmt.Println("MONGODB CONNECTED")

	//GIN DEFINE
	router := gin.Default()
	api := router.Group("/api")
	{
		api.GET("/", func(ctx *gin.Context) {
			ctx.JSON(200, gin.H{
				"message": "Ping successful",
			})
		})
	}

	//ROUTER DEFINE
	routes.UserRouter(api)
	routes.DiaryNoteRouter(api)
	routes.ReviewNoteRouter(api)
	routes.DayReviewRouter(api)

	port := fmt.Sprintf(":%s", os.Getenv("APP_PORT"))
	router.Run(port)
}
