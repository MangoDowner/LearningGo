package scripts

import (
    "io/ioutil"
    "os"
    "os/exec"
    "path/filepath"
    "strings"
    "github.com/lxn/walk"
    . "github.com/lxn/walk/declarative"
    "public"
    "log"
)


type MyMainWindow struct {
    *walk.MainWindow
}

//展现主框体
func CreateSearchViedoFrame() {
    var mainWindow *walk.MainWindow
    var textEdit *walk.TextEdit //拖拽文件目录匡
    var startSearchBtn *walk.PushButton //触发时间的按钮
    var tableView *walk.TableView
    var webView *walk.WebView //网页浏览器

    tableModel := NewFileInfoModel()

    if err := (MainWindow{
        AssignTo: &mainWindow,
        Title:   "发现视频",
        MinSize:  Size{600, 400},
        Size:     Size{1024, 640},
        Layout:   HBox{MarginsZero: true},
        OnDropFiles: func(files []string) {
            textEdit.SetText(strings.Join(files, "\r\n"))
        },
        Children: []Widget{
            HSplitter{
                Children: []Widget{
                    VSplitter{
                        Children: []Widget{
                            TextEdit{
                                AssignTo: &textEdit,
                                ReadOnly: true,
                                Text:     "将 文 件 夹 拉 到 这 里",
                            },
                            PushButton {
                                AssignTo: &startSearchBtn,
                                Text:     "我 们 开 始 搜 索 了 ！",
                                OnClicked: func() { SearchViedo(textEdit, tableModel) },
                            },
                        },
                    },
                    TableView{
                        AssignTo:      &tableView,
                        StretchFactor: 2,
                        Columns: []TableViewColumn{
                            TableViewColumn{
                                Title: "名称",
                                DataMember: "Name",
                                Width: 300,
                            },
                            //TableViewColumn{
                            //    Title: "大小",
                            //    DataMember: "Size",
                            //    Width: 100,
                            //},
                            //TableViewColumn{
                            //    Title: "最后修改时间",
                            //    DataMember: "Modified",
                            //    Format:     "2006-01-02 15:04:05",
                            //    Width: 150,
                            //},
                        },
                        Model: tableModel,
                        OnCurrentIndexChanged: func() {
                            var url string
                            if index := tableView.CurrentIndex(); index > -1 {
                                path := tableModel.items[index].Path
                            //    //dir := treeView.CurrentItem().(*Directory)
                            //    //url = filepath.Join(dir.Path(), name)
                                url = path
                            }
                            webView.SetURL(url)
                        },
                    },
                    WebView {
                        AssignTo: &webView,
                        StretchFactor: 2,
                    },
                },
            },
        },
    }.Create()); err != nil {
        log.Fatal(err)
    }
    public.SetIcon(mainWindow, "")
    mainWindow.Run()
}



//搜寻视频操作
func SearchViedo(inTE *walk.TextEdit, tableModel *FileInfoModel) {
    //resultFolderName := "E:\\SearchResult" //存放搜索结果的文件夹
    //os.Mkdir(resultFolderName, 777)
    //quickDirName := resultFolderName + "\\Quick" //存放搜索结果快捷方式的文件夹
    //os.Mkdir(quickDirName, 777)
    //resultTextName := resultFolderName + "\\文件清单.txt"
    var suffixArr = []string{"mpeg", "avi", "mov", "wmv", "mkv", "mp4"}
    //var suffixArr = []string{"jpg", "gif"}
    searchPath := inTE.Text()
    tableModel.WalkDir(string(searchPath), suffixArr)
    //files, _, _ := WalkDir(string(searchPath), suffixArr)
    //f, _ := os.OpenFile(resultTextName, os.O_RDWR|os.O_CREATE|os.O_APPEND, 0x644)
    //for _, v := range files {
    //    f.WriteString(v + "\r\n")
        //CreateQuickRef(v, names[k], quickDirName) //建立文本的快捷方式
    //}
    //tableView.SetModel(tableModel)
    //inTE.SetText("搜索结果储存于 E:\\SearchResult")
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