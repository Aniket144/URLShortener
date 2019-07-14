package handlers

import (
	"crypto/md5"
	"database/sql"
	"encoding/hex"
	"fmt"
	"net/http"
	"strconv"
	"github.com/gin-gonic/gin"
)

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
	db, err := sql.Open("mysql", "root:root@tcp(127.0.0.1:3306)/URLShortner")
	if err != nil {
		fmt.Println(err)
		return
	}
	defer db.Close()

	link := c.PostForm("link")
	hash := generateHash(link)

	/* Default Message */
	message := "Successfully Generated"

	/* Check whether the Link already has a short link */
	shortLink, alreadyExist := getShortLink(db, hash, link)
	if alreadyExist {
		message = "The Short Link Already Exist"
	}

	/* Rendering Successful Creation Page OR Already Exist*/
	c.HTML(http.StatusOK, "posting.tmpl", gin.H{
		"title":   "URLShortner Page",
		"message": message,
		"link":    shortLink,
	})
}

/* Redirect to the Main(Long) URL */
func ShortLinkRedirect(c *gin.Context) {
	hash := c.Params.ByName("hash")

	/* Connecting to Database */
	db, err := sql.Open("mysql", "root:root@tcp(127.0.0.1:3306)/URLShortner")
	if err != nil {
		fmt.Println(err)
		return
	}
	defer db.Close()

	originalLink, found, hits := getLongLink(db, hash)

	if found {
		increaseHits(db, hash, hits)
		c.Redirect(http.StatusMovedPermanently, "http://"+originalLink)
	} else {
		c.Redirect(301, "/")
	}
}

/* To increase the count of hits a short URL receives */
func increaseHits(db *sql.DB, hash string, originalHits int) {
	query := "UPDATE links SET Hits = " + strconv.Itoa(originalHits+1) + " WHERE Hash = '" + hash + "';"
	_, err := db.Query(query)
	if err != nil {
		fmt.Println(err)
	}
}

/* Retrieve the Long URL by matching the Hash */
func getLongLink(db *sql.DB, hash string) (string, bool, int) {
	query := "SELECT * FROM links WHERE Hash = '" + hash + "'"
	res, err := db.Query(query)
	defer res.Close()

	var retrivedHash string
	var retrivedLink string
	var retrivedHits int

	/* If there's an error or no rows returned */
	if err != nil {
		return "", false, 0
	}
	if !res.Next() {
		return "", false, 0
	}

	res.Scan(&retrivedHash, &retrivedLink, &retrivedHits)
	return retrivedLink, true, retrivedHits
}

/* Generate Hash of the Long URL using md5 algorithm */
func generateHash(link string) string {
	hasher := md5.New()
	hasher.Write([]byte(link))
	return hex.EncodeToString(hasher.Sum(nil))[:10]
}

/* Get the shirt link of the by searching in DB using hash as key */
func getShortLink(db *sql.DB, hash string, link string) (string, bool) {
	shortLink := "localhost:8080/h/" + hash
	_, alreadyExist, _ := getLongLink(db, hash)
	if alreadyExist {
		return shortLink, true
	}
	query := "INSERT INTO links VALUES ('" + hash + "','" + link + "', " + strconv.Itoa(0) + ");"
	_, err := db.Query(query)
	if err != nil {
		fmt.Println(err)
	}
	return shortLink, false
}


