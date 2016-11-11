package main
//rsrc -manifest test.manifest -o rsrc.syso 已经做过无需在做
import (
    "github.com/lxn/walk"
    . "github.com/lxn/walk/declarative"
    "operation"
)

func main() {
    var inTE, outTE *walk.TextEdit
    MainWindow{
        Title:   "第一个APP",
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
                Text: "到右边去吧",
                OnClicked: func(){
                    operation.SearchViedos(inTE, outTE) //寻找视频操作
                },
            },
        },
    }.Run()
}

