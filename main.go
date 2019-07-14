package main

import (
	controller "URLShortener/handlers"
	"github.com/gin-gonic/gin"
	_ "github.com/go-sql-driver/mysql"
)

func main() {
	router := gin.Default()
	router.LoadHTMLGlob("templates/*")

	/* Routes */
	router.GET("/", controller.Home)
	router.POST("/", controller.CreateShortLink)
	router.GET("/h/:hash", controller.ShortLinkRedirect)

	/* Server running at default Port 8080 */
	router.Run()
}
