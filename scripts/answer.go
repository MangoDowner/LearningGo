package scripts

import (
	"github.com/lxn/walk"
	. "github.com/lxn/walk/declarative"
	"public"
	"log"
	"fmt"
	"io/ioutil"
	"net/http"
	"regexp"
)

//创建冲顶大会界面
func CreateAnswerFrame() {
	var mainWindow *walk.MainWindow
	var inTE *walk.LineEdit
	var outTE *walk.TextEdit

	var fontYahei Font
	fontYahei.Family = "Consolas"
	fontYahei.Create()

	var bigFontYahei Font
	bigFontYahei.PointSize = 160
	bigFontYahei.Create()


	if err := (MainWindow{
		AssignTo: &mainWindow,
		Title:   "冲顶大会",
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
						Text: "开始答题",
						OnClicked: func() {
							searchWord := inTE.Text()
							data := GetResponseNumFromSogou(searchWord)
							outTE.SetText(data)
						},
					},
					TextEdit{
						AssignTo: &outTE,
						Font: bigFontYahei,
						ColumnSpan: 2,
						ReadOnly: true,
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


//获取百度返回的内容数目
func GetResponseNumFromSogou(searchWord string) (data string) {
	searchWord = GetQueryString(searchWord)
	//url := "https://www.baidu.com/s?wd=" + searchWord
	url := "https://www.sogou.com/web?query=" + searchWord
	req, _ := http.NewRequest("GET", url, nil)
	req.Header.Set("User-Agent", GetRandomUserAgent())
	client := http.DefaultClient
	res, e := client.Do(req)
	if e != nil {
		fmt.Errorf("Get请求地址%s返回错误:%s", url, e)
		return data
	}
	if res.StatusCode == 200 {
		body := res.Body
		defer body.Close()
		bodyByte, _ := ioutil.ReadAll(body)
		oriText := string(bodyByte)
		var regExpress = regexp.MustCompile(`搜狗已为您找到约([\w,]+)条相关结果`)
		validArr := regExpress.FindAllStringSubmatch(oriText, -1)
		for _, v := range validArr {
			return v[1]
		}
	}
	return data
}