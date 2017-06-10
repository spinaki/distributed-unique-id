package main

import (
        "gopkg.in/gin-gonic/gin.v1"
        "gopkg.in/gin-contrib/cors.v1"
        "net/http"
        "time"
        "strconv"
)

var idGeneratorSettings *Settings
func main() {
        // build snowflake using the IdGenerator API
        idGeneratorSettings = &Settings{}
        idGeneratorSettings.StartTime = time.Date(2016, 1, 1, 0, 0, 0, 0, time.UTC)

        // build
        gin.SetMode(gin.ReleaseMode)
        router := gin.Default()
        corsConfig := cors.DefaultConfig()
        corsConfig.AllowAllOrigins = true
        corsConfig.ExposeHeaders = []string{"Content-Type, Content-Length, Accept-Encoding, X-CSRF-Token, Authorization, " +
                "accept, origin, Cache-Control, X-Requested-With"}
        corsConfig.AllowHeaders = []string{"Origin", "Content-Length", "Content-Type"}
        corsConfig.AllowMethods = []string{"GET"}
        newCors := cors.New(corsConfig)
        router.Use(newCors)
        router.Use(gin.Logger())

        router.GET("/status", statusHandler)
        router.GET("/stringids", stringIdsHandler)
        router.GET("/longids", longIdsHandler)
        router.GET("/longidrange", longIdRangeHandler)
        // have a status endpoint too
        router.Run(":8080")
}

func statusHandler(c *gin.Context) {
        c.String(http.StatusOK,"OK")
}

func stringIdsHandler(c *gin.Context) {
        // num of ids and length of ids
        num := c.DefaultQuery("num", "10") // default num of ids = 10
        len := c.DefaultQuery("len", "32") // number of bytes used to generate random id = 32
        // NOTE: the char size of base64 encoded string will be different from the num of bytes used.
        l, _ := strconv.Atoi(len)
        n, _ := strconv.Atoi(num)
        ids := GenerateRandomStringId(l, n)
        strIdList := &StringIDList{List:ids}
        c.JSON(http.StatusOK, strIdList)
}

func longIdsHandler(c *gin.Context) {
        idList, err := GenerateIDList(idGeneratorSettings)
        if err != nil {
                c.JSON(http.StatusInternalServerError, gin.H{"result": "Failed to generate unique integer id list"})
                return
        }
        c.JSON(http.StatusOK, idList)
}

func longIdRangeHandler(c *gin.Context) {
        idRange, err := GenerateIDRange(idGeneratorSettings)
        if err != nil {
                c.JSON(http.StatusInternalServerError, gin.H{"result": "Failed to generate unique integer id range"})
                return
        }
        c.JSON(http.StatusOK, idRange)
}

