package operation

import (
    "github.com/lxn/walk"
    "fmt"
    "os"
    "scripts"
    "strings"
)

//搜寻视频操作
func SearchViedos(inTE, outTE *walk.TextEdit) {
    fmt.Println(inTE.Text())
    resultFolderName := "E:\\SearchResult" //存放搜索结果的文件夹
    os.Mkdir(resultFolderName, 777)
    quickDirName := resultFolderName + "\\Quick" //存放搜索结果快捷方式的文件夹
    os.Mkdir(quickDirName, 777)
    resultTextName := resultFolderName + "\\文件清单.txt"
    //var suffixArr = []string{"mpeg", "avi", "mov", "wmv", "mkv", "mp4"}
    var suffixArr = []string{"jpg", "gif"}
    searchPath := inTE.Text()
    files, names, _ := scripts.WalkDir(string(searchPath), suffixArr)
    f, _ := os.OpenFile(resultTextName, os.O_RDWR|os.O_CREATE|os.O_APPEND, 0x644)
    for k, v := range files {
        f.WriteString(v + "\r\n")
        scripts.CreateQuickRef(v, names[k], quickDirName)
        fmt.Println(v)
    }
    fmt.Println("\r\n---------搜索结果储存于 E:\\SearchResult----------\r\n")
    outTE.SetText(strings.ToUpper(inTE.Text()))
}
