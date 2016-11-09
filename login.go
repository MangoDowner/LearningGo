package main

import (
    "github.com/gin-gonic/gin"
    "net/http"
    "model/mysql"
)

func main() {
    router := gin.Default()

    router.POST("login", login);

    router.Run(":8082")
}

func login(c *gin.Context) {
    var form mysql.LoginForm
    if c.Bind(&form) == nil {
        if form.User == "mango" && form.Password == "123" {
            c.JSON(
                http.StatusOK,
                gin.H{"status" : "StatusOK"},
            )
        } else {
            c.JSON(
                http.StatusUnauthorized,
                gin.H{"status" : "StatusUnauthorized"},
            )
        }
    }
}

//curl -v --form user=mango --form password=123 http://localhost:8082/login

































