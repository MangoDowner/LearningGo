package mysql


type LoginForm struct {
    User        string `form:"user" json:"user" binding:"required"`
    Password    string `form:"password" json:"password" binding:"required"`
}
