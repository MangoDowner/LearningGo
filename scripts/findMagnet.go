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
    var inTE *walk.LineEdit
    var outTE *walk.TextEdit

    var fontYahei Font
    fontYahei.Family = "Consolas"
    fontYahei.Create()

    fmt.Println(Font{})
    if err := (MainWindow{
        AssignTo: &mainWindow,
        Title:   "查找磁性链接",
        MinSize: Size{600, 400},
        Layout:  VBox{},
        Children: []Widget{
            Composite{
                Layout: Grid{Columns: 2},
                Children: []Widget{
                    LineEdit{
                        Font: fontYahei,
                        AssignTo: &inTE,
                    },
                    PushButton{
                        Font: fontYahei,
                        Text: "开始查找",
                        OnClicked: func() {
                            outTE.SetText(strings.ToUpper(inTE.Text()))
                        },
                    },
                    TextEdit{
                        Font: fontYahei,
                        ColumnSpan: 2,
                        AssignTo: &outTE,
                        ReadOnly: true,
                    },
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
