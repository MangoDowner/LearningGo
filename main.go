package main

import (
    "github.com/gin-gonic/gin"
    "net/http"
)

func main() {
    router := gin.Default()
    router.GET("/user/:name", func(c *gin.Context){
        name := c.Param("name")
        c.String(http.StatusOK, "你好啊, %s", name)
    })

    router.GET("/user/:name/:action", func(c *gin.Context){
        name := c.Param("name")
        action := c.Param("action")
        message := name + "在" + action
        c.JSON(200, gin.H{
            "code": 1001,
            "message": message,
        })
    })

    router.GET("/welcome", func(c *gin.Context){
        firstName := c.DefaultQuery("firstname", "Guest")
        lastName := c.Query("lastname")
        c.String(http.StatusOK, "Hello %s %s", firstName, lastName)
    })

    router.POST("/form_post", func(c *gin.Context){
        message := c.PostForm("message")
        nick := c.DefaultPostForm("nick", "anonymous")
        c.JSON(200, gin.H{
            "status": "posted",
            "message": message,
            "nick": nick,
        })
    })

    router.POST("/post", func(c *gin.Context){
        id := c.Query("id")
        page := c.DefaultQuery("page", "0")
        name := c.PostForm("name")
        message := c.PostForm("message")
        c.JSON(200, gin.H{
            "id": id,
            "page": page,
            "name": name,
            "message": message,
        })
    })

    v1 := router.Group("/v1")
    {
        v1.GET("/user/:name", func(c *gin.Context){
            name := c.Param("name")
            c.JSON(http.StatusOK, gin.H{
                "status": "posted",
                "message": "你好啊",
                "nick": name,
            })
        })
    }

    v2 := router.Group("/v2")
    {
        v2.GET("/user", UserGetUserInfo)
    }

    router.Run(":8080")
}



func UserGetUserInfo(c *gin.Context){
    name := c.Query("name")
    c.String(http.StatusOK, "我在这里, %s", name)
}



























