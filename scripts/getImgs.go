package main

import (
    "fmt"
    "io/ioutil"
    "os"
    "os/exec"
    "path/filepath"
    "strings"
    "bufio"
)

func main() {
    //fmt.Println(">>>>>输入要搜索的盘符，以回车结束")
    resultFolderName := "D:\\SearchResult" //存放搜索结果的文件夹
    os.Mkdir(resultFolderName, 777)
    quickDirName := resultFolderName + "\\Quick" //存放搜索结果快捷方式的文件夹
    os.Mkdir(quickDirName, 777)
    resultTextName := resultFolderName + "\\文件清单.txt"
    var suffixArr = []string{"mpeg", "avi", "mov", "wmv", "mkv", "mp4"}
    //var suffixArr = []string{"jpg", "gif"}
    reader := bufio.NewReader(os.Stdin)
    for {
        fmt.Println(">>>>>输入要搜索的目录，以回车结束")
        searchPath, _, _ := reader.ReadLine()
        files, names, _ := WalkDir(string(searchPath), suffixArr)
        f, _ := os.OpenFile(resultTextName, os.O_RDWR|os.O_CREATE|os.O_APPEND, 0x644)
        for k, v := range files {
            f.WriteString(v + "\r\n")
            createQuickRef(v, names[k], quickDirName)
            fmt.Println(v)
        }
        fmt.Println("\r\n---------搜索结果储存于 D:\\SearchResult----------\r\n")
    }
}


//获取指定目录下的所有文件，不进入下一级目录搜索，可以匹配后缀过滤。
func ListDir(dirPth string, suffix string) (files []string, err error) {
    files = make([]string, 0, 10)
    dir, err := ioutil.ReadDir(dirPth)
    if err != nil {
        return nil, err
    }
    PthSep := string(os.PathSeparator)
    suffix = strings.ToUpper(suffix) //忽略后缀匹配的大小写
    for _, fi := range dir {
        if fi.IsDir() { // 忽略目录
            continue
        }
        if strings.HasSuffix(strings.ToUpper(fi.Name()), suffix) { //匹配文件
            files = append(files, dirPth+PthSep+fi.Name())
        }
    }
    return files, nil
}

//获取指定目录及所有子目录下的所有文件，可以匹配后缀过滤。
func WalkDir(dirPth string, suffixArr []string) (files []string, names []string, err error) {
    files = make([]string, 0, 100)                                                       //文件路径
    names = make([]string, 0, 100)                                                       //文件名称
    err = filepath.Walk(dirPth, func(filename string, fi os.FileInfo, err error) error { //遍历目录
        if err != nil { //忽略错误
            return err
        }
        if fi.IsDir() { // 忽略目录
            return nil
        }
        //忽略后缀匹配的大小写
        for _, v := range suffixArr {
            suffix := "." + strings.ToUpper(v) //忽略后缀匹配的大小写
            if strings.HasSuffix(strings.ToUpper(fi.Name()), suffix) {
                files = append(files, filename)
                names = append(names, fi.Name())
            }
        }
        return nil
    })
    return files, names, err
}

//建立快捷方式
func createQuickRef(path string, name string, destPath string) {
    dest_path := destPath + "\\" + name
    c := exec.Command("cmd", "/C", "echo [InternetShortcut] >>"+dest_path+".url")
    c.Run()
    c = exec.Command("cmd", "/C", "echo URL="+path+" >>"+dest_path+".url")
    c.Run()
    c = exec.Command("cmd", "/C", "echo IconIndex=0 >>"+dest_path+".url")
    c.Run()
    c = exec.Command("cmd", "/C", "echo IconFile="+path+" >>"+dest_path+".url")
    c.Run()
}

