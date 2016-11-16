package scripts

import (
"github.com/lxn/walk"
. "github.com/lxn/walk/declarative"
"strings"
    "public"
    "log"
    "fmt"
)

func CreateFindMagnetFrame() {
    var mainWindow *walk.MainWindow
    var inTE, outTE *walk.TextEdit

    var fontYahei Font
    fontYahei.Family = "微软雅黑"
    fontYahei.Create()

    fmt.Println(Font{})
    if err := (MainWindow{
        AssignTo: &mainWindow,
        Title:   "查找磁性链接",
        MinSize: Size{600, 400},
        Layout:  VBox{},
        Children: []Widget{
            HSplitter{
                Children: []Widget{
                    TextEdit{AssignTo: &inTE},
                    TextEdit{AssignTo: &outTE, ReadOnly: true},
                },
            },
            PushButton{
                Font: fontYahei,
                Text: "开始查找",
                OnClicked: func() {
                    outTE.SetText(strings.ToUpper(inTE.Text()))
                },
            },
        },
    }.Create()); err != nil {
        log.Fatal(err)
    }
    public.SetIcon(mainWindow, "")
    mainWindow.Font()
    mainWindow.Run()
}
