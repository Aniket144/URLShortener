package handlers

import (
	"URLShortener/migrations"
	"crypto/md5"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/jinzhu/gorm"
	"io/ioutil"
	"net/http"
	"os"
	"sync"
	_ "github.com/go-sql-driver/mysql"
)

type Link struct {
	hash string
	URL string
	hits int
}

/* For retrieving URLs through JSON */
type URLS[] string


/* Setting Default Domain Name */
var siteDomain = "localhost:8080"

/* Setting Database credentials */
var databaseName = "URLShortener"
var username = "root"
var password = "root"
var address = "127.0.0.1:3306"
//var databaseSource = username + ":" + password + "@" + "tcp(" + address + ")/" + databaseName
var databaseSource = username + ":" + password + "@/" + databaseName + "?charset=utf8&parseTime=True&loc=Local"


var DB *gorm.DB

func init() {
	fmt.Println("Here")
	DB, _ = gorm.Open("mysql", "root:root@/URLShortener?charset=utf8&parseTime=True&loc=Local")
	fmt.Println(DB)
}

/* Home Page */
func Home(c *gin.Context) {
	c.HTML(http.StatusOK, "home.tmpl", gin.H{
		"title":   "Home Page",
		"message": "Get your URL Shorten Here",
	})
}

/* File Upload Page */
func FileUpload(c *gin.Context) {
	c.HTML(http.StatusOK, "fileUpload.tmpl", gin.H{
		"title": "Upload URL File",
	})
}

/* File Parsing */
func FileParsing(c *gin.Context) {
	r := c.Request
	r.ParseMultipartForm(10 << 20)

	file, handler, err := r.FormFile("file")
	if err != nil {
		fmt.Println(err)
		return
	}
	defer file.Close()

	fmt.Println("file name => ", handler.Filename)

	tempFile, err := ioutil.TempFile("saved", "*.json")
	//fmt.Println(tempFile.Name())
	if err != nil {
		fmt.Println(err)
		return
	}
	defer tempFile.Close()

	fileBytes, err := ioutil.ReadAll(file)
	if err != nil {
		fmt.Println(err)
	}
	tempFile.Write(fileBytes)
	fmt.Println("Done", tempFile.Name())

	parseJSONFile(tempFile.Name())

	c.String(http.StatusOK, "File Uploaded Successfully")
}

func parseJSONFile(fileName string) {
	path := fileName
	jsonFile, err := os.Open(path)
	if err != nil {
		fmt.Println(jsonFile)
	}
	byteValue, _ := ioutil.ReadAll(jsonFile)
	var fileURLs URLS
	err = json.Unmarshal([]byte(byteValue), &fileURLs)
	if err != nil {
		fmt.Println(err)
	}

	createShortLinks(fileURLs)
}

func createShortLinks(Urls []string) {
	var wg sync.WaitGroup
	for _, url := range Urls {
		wg.Add(1)
		go func(temp_url string) {
			hash := generateHash(temp_url)
			getElseCreateShortLink(hash, temp_url)
			wg.Done()
		}(url)
	}
	wg.Wait()
}

var mutex = &sync.Mutex{}

/* Checks whether a short link exists else create it */
func CreateShortLink(c *gin.Context) {
	fmt.Println("here 2", DB)
	url := c.PostForm("url")
	hash := generateHash(url)

	/* Default Message */
	message := "Successfully Generated"

	/* Check whether the url already has a short link */
	shortLink, alreadyExist := getElseCreateShortLink(hash, url)
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
	originalLink, found, hits := getLongLink(hash)

	if found {
		go increaseHits(hash, hits)
		c.Redirect(http.StatusFound, "http://"+originalLink)
	} else {
		c.Redirect(http.StatusFound, "/")
	}
}

/* To increase the count of hits a short URL receives */
func increaseHits(hash string, originalHits int) {
	var link migrations.Link
	DB.Where("hash = ?", hash).First(&link)

	/* If there's no retrieval */
	if len(link.URL) == 0 {
		return
	}
	link.Hits = originalHits + 1
	DB.Save(&link)
	DB.Close()
}

/* Retrieve the Long URL by matching the Hash */
func getLongLink(hash string) (string, bool, int) {
	var link migrations.Link
	DB.Where("hash = ?", hash).First(&link)
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
func getElseCreateShortLink(hash string, url string) (string, bool) {
	shortLink := siteDomain + "/h/" + hash
	_, alreadyExist, _ := getLongLink(hash)
	if alreadyExist {
		return shortLink, true
	}
	mutex.Lock()
	DB.Create(&migrations.Link{Hash:hash, URL:url, Hits: 0})
	mutex.Unlock()
	return shortLink, false
}