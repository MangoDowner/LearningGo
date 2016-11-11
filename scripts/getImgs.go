package scripts

import (
    "io/ioutil"
    "os"
    "os/exec"
    "path/filepath"
    "strings"
)


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
func CreateQuickRef(path string, name string, destPath string) {
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

