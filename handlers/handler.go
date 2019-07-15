package handlers

import (
	"URLShortener/migrations"
	"crypto/md5"
	"encoding/hex"
	"github.com/gin-gonic/gin"
	"github.com/jinzhu/gorm"
	"net/http"
)

type Link struct {
	hash string
	URL string
	hits int
}

/* Setting Default Domain Name */
var siteDomain = "localhost:8080"

/* Setting Database credentials */
var databaseName = "URLShortener"
var username = "root"
var password = "root"
var address = "127.0.0.1:3306"
//var databaseSource = username + ":" + password + "@" + "tcp(" + address + ")/" + databaseName
var databaseSource = username + ":" + password + "@/" + databaseName + "?charset=utf8&parseTime=True&loc=Local"


/* Home Page */
func Home(c *gin.Context) {
	c.HTML(http.StatusOK, "home.tmpl", gin.H{
		"title":   "Home Page",
		"message": "Get your URL Shorten Here",
	})
}

/* Checks whether a short link exists else create it */
func CreateShortLink(c *gin.Context) {
	/* Connecting to Database */
	db, err := gorm.Open("mysql", databaseSource)
	defer db.Close()
	if err != nil {
		println(err)
		return
	}

	url := c.PostForm("url")
	hash := generateHash(url)

	/* Default Message */
	message := "Successfully Generated"

	/* Check whether the url already has a short link */
	shortLink, alreadyExist := getShortLink(db, hash, url)
	if alreadyExist {
		message = "The Short Link Already Exist"
	}

	/* Rendering Successful Creation Page OR Already Exist*/
	c.HTML(http.StatusOK, "posting.tmpl", gin.H{
		"title":   "URLShortener Page",
		"message": message,
		"link":    shortLink,
	})
}

/* Redirect to the Main(Long) URL */
func ShortLinkRedirect(c *gin.Context) {
	hash := c.Params.ByName("hash")

	/* Connecting to Database */
	db, err := gorm.Open("mysql", databaseSource)
	if err != nil {
		println(err)
	}
	defer db.Close()

	originalLink, found, hits := getLongLink(db, hash)

	if found {
		increaseHits(db, hash, hits)
		c.Redirect(http.StatusFound, "http://"+originalLink)
	} else {
		c.Redirect(http.StatusFound, "/")
	}
}

/* To increase the count of hits a short URL receives */
func increaseHits(db *gorm.DB, hash string, originalHits int) {
	var link migrations.Link
	db.Where("hash = ?", hash).First(&link)

	/* If there's no retrieval */
	if len(link.URL) == 0 {
		return
	}
	link.Hits = originalHits + 1
	db.Save(&link)
}

/* Retrieve the Long URL by matching the Hash */
func getLongLink(db *gorm.DB, hash string) (string, bool, int) {
	var link migrations.Link
	db.Where("hash = ?", hash).First(&link)
	/* IF there's no retrieval */
	if len(link.URL) == 0 {
		return "", false, 0
	}

	return  link.URL, true, link.Hits
}

/* Generate Hash of the Long URL using md5 algorithm */
func generateHash(link string) string {
	hasher := md5.New()
	hasher.Write([]byte(link))
	return hex.EncodeToString(hasher.Sum(nil))[:10]
}

/* Get the shirt link of the by searching in DB using hash as key */
func getShortLink(db *gorm.DB, hash string, url string) (string, bool) {
	shortLink := siteDomain + "/h/" + hash
	_, alreadyExist, _ := getLongLink(db, hash)
	if alreadyExist {
		return shortLink, true
	}
	db.Create(&migrations.Link{Hash:hash, URL:url, Hits: 0})
	return shortLink, false
}