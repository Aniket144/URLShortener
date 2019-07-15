package main

import (
	controller "URLShortener/handlers"
	"URLShortener/migrations"
	"github.com/gin-gonic/gin"
	_ "github.com/go-sql-driver/mysql"
	"github.com/jinzhu/gorm"
)

func main() {
	router := gin.Default()

	/* Loading HTML Templates Folder */
	router.LoadHTMLGlob("templates/*")

	/* Migrating Table */

		db, err := gorm.Open("mysql", "root:root@/URLShortener?charset=utf8&parseTime=True&loc=Local")
		defer db.Close()
		if err != nil {
			println(err)
		}
		db.AutoMigrate(&migrations.Link{})

	///* Routes */
	router.GET("/", controller.Home)
	router.POST("/", controller.CreateShortLink)
	router.GET("/h/:hash", controller.ShortLinkRedirect)

	/* Server running at default Port 8080 */
	router.Run()
}
