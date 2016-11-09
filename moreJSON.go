package main

import (
    "github.com/gin-gonic/gin"
    "net/http"
)

func main() {
    r := gin.Default()
    r.GET("/someJSON", func(c *gin.Context) {
        c.JSON(http.StatusOK, gin.H{"message" : "hey", "status" : http.StatusOK})
    });

    r.GET("/moreJSON", func(c *gin.Context) {
        var msg struct {
            Name string `json:"user"`
            Message string
            Number int
        }
        msg.Name = "Lena"
        msg.Message = "hey"
        msg.Number = 123
        c.JSON(http.StatusOK, msg)
    })
    r.Run(":8083")
}
