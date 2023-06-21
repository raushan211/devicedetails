package main

import (
	"database/sql"
	"fmt"
	"net/http"

	//"net/url"
	"time"

	"github.com/gin-gonic/gin"
)

var DB *sql.DB

type Device struct {
	ID                int       `json:"id"`
	TYPE              string    `json:"type"`
	BROWSER           string    `json:"browser"`
	BROWSER_VERSION   string    `json:"browser_version"`
	TIME_STAMP        time.Time `json:"time_stamp"`
	SCREEN_RESOLUTION string    `json:"screen_resolution"`
}

func main() {
	createDBConnection()
	defer DB.Close()
	r := gin.Default()
	r.Use(CORSMiddleware())
	setupRoutes(r)
	r.Run() // listen and serve on 0.0.0.0:8080 (for windows "localhost:8080")
}
func setupRoutes(r *gin.Engine) {

	r.POST("/device_details", SaveLongLink)
	r.GET("/device_details/all", GetAllDevices)
	r.NoRoute(func(c *gin.Context) {
		c.JSON(404, gin.H{"code": "PAGE_NOT_FOUND", "message": "Page not found"})
	})

}
func CORSMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		c.Writer.Header().Set("Access-Control-Allow-Origin", "*")
		c.Writer.Header().Set("Access-Control-Allow-Credentials", "true")
		c.Writer.Header().Set("Access-Control-Allow-Headers", "Content-Type, Content-Length, Accept-Encoding, X-CSRF-Token, Authorization, accept, origin, Cache-Control, X-Requested-With")
		c.Writer.Header().Set("Access-Control-Allow-Methods", "POST, OPTIONS, GET, PUT")

		if c.Request.Method == "OPTIONS" {
			c.AbortWithStatus(204)
			return
		}

		c.Next()
	}
}

// POST
func SaveLongLink(c *gin.Context) {
	reqBody := Device{}
	err := c.Bind(&reqBody)
	if err != nil {
		res := gin.H{
			"error": "invalid request body",
		}
		c.Writer.Header().Set("Content-Type", "application/json")
		c.JSON(http.StatusBadRequest, res)

		return
	}

	//reqBody.ValidUrl = validurl(reqBody.URL)

	// Data[lastID] = reqBody
	reqBody.TIME_STAMP = time.Now()
	fmt.Println(reqBody)
	res, err := DB.Exec(`INSERT INTO "device_details" ( "type", "browser", "browser_version", "time_stamp", "screen_resolution")
	VALUES ( $1, $2, $3, $4, $5)`, reqBody.TYPE, reqBody.BROWSER, reqBody.BROWSER_VERSION, reqBody.TIME_STAMP, reqBody.SCREEN_RESOLUTION)
	if err != nil {
		fmt.Println("err inserting data: ", err)
		c.Writer.Header().Set("Content-Type", "application/json")
		c.JSON(http.StatusBadRequest, gin.H{
			"error": err.Error(),
		})
		return
	}
	lastInsID, _ := res.LastInsertId()
	reqBody.ID = int(lastInsID)
	fmt.Println("res: ", lastInsID)
	c.JSON(http.StatusOK, reqBody)
	c.Writer.Header().Set("Content-Type", "application/json")
}

// GET
func GetAllDevices(c *gin.Context) {
	rows, err := DB.Query("SELECT id, type, browser, browser_version, time_stamp, screen_resolution FROM device_details order by id desc")
	if err != nil {
		panic(err)
	}
	defer rows.Close()
	devices := []Device{}
	for rows.Next() {
		d := Device{}
		err := rows.Scan(&d.ID, &d.TYPE, &d.BROWSER, &d.BROWSER_VERSION, &d.TIME_STAMP, &d.SCREEN_RESOLUTION)
		if err != nil {
			panic(err)
		}
		devices = append(devices, d)
	}
	err = rows.Err()
	if err != nil {
		panic(err)
	}
	res := gin.H{
		"data": devices,
	}
	c.JSON(http.StatusOK, res)
}
