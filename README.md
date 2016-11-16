# LearningGo

###需要的库
go get github.com/gin-gonic/gin

go get github.com/lxn/walk

go get github.com/akavel/rsrc
###walk的文档
https://godoc.org/github.com/lxn/walk#pkg-index

###生成没有命令行界面的exe指令
go build -ldflags="-H windowsgui"