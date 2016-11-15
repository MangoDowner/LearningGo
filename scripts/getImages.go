package scripts

import (
    "time"
    "github.com/lxn/walk"
    "path/filepath"
    "os"
    "strings"
    "fmt"
)

//文件属性
type FileInfo struct {
    Name     string
    Size     string
    Modified time.Time
}

type FileInfoModel struct {
    walk.SortedReflectTableModelBase
    dirPath string
    items   []*FileInfo
}

func NewFileInfoModel() *FileInfoModel {
    return new(FileInfoModel)
}

func (m *FileInfoModel) Items() interface{} {
    return m.items
}

//获取指定目录及所有子目录下的所有文件，可以匹配后缀过滤。
func (m *FileInfoModel) WalkDir(dirPth string, suffixArr []string) (files []string, names []string, err error) {
    m.dirPath = dirPth
    m.items = nil

    files = make([]string, 0, 100)                                                       //文件路径
    names = make([]string, 0, 100)
    //文件名称

    err = filepath.Walk(dirPth, func(filename string, info os.FileInfo, err error) error { //遍历目录
        if err != nil { //忽略错误
            return err
        }
        if info.IsDir() { // 忽略目录
            return nil
        }
        //忽略后缀匹配的大小写
        for _, v := range suffixArr {
            suffix := "." + strings.ToUpper(v) //忽略后缀匹配的大小写
            if strings.HasSuffix(strings.ToUpper(info.Name()), suffix) {
                files = append(files, filename)
                names = append(names, info.Name())
                size := ""
                if  int(info.Size() / (1024 * 1024)) > 1  {
                    size = fmt.Sprintf( "%d GB", int(info.Size() / (1024 * 1024)) )
                } else {
                    size = fmt.Sprintf( "%d MB", int(info.Size() / 1024) )
                }
                item := &FileInfo{
                    Name:     info.Name(),
                    Size:     size,
                    Modified: info.ModTime(),
                }
                m.items = append(m.items, item)
            }
        }
        m.PublishRowsReset()
        return nil
    })
    return files, names, err
}
