package scripts

import (
"github.com/lxn/walk"
. "github.com/lxn/walk/declarative"
"strings"
    "public"
    "log"
    "fmt"
    "encoding/xml"
    "time"
    "net/http"
    "io/ioutil"
    "regexp"
    "math/rand"
)

type MagnetItem struct {
    title  string  //标题
    path   string  //地址
}

type MagnetModel struct {
    walk.ListModelBase
    items []MagnetItem
}

func (m *MagnetModel) ItemCount() int {
    return len(m.items)
}

func (m *MagnetModel) Value(index int) interface{} {
    return m.items[index].title
}

func CreateFindMagnetFrame() {
    var mainWindow *walk.MainWindow
    var inTE *walk.LineEdit
    var outTE *walk.TextEdit
    var listBox *walk.ListBox

    var fontYahei Font
    fontYahei.Family = "Consolas"
    fontYahei.Create()

    magnetModel := new(MagnetModel)
    //magnetModel := GetResponseFromTorrentKitty("iptd-651")

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
                            searchWord := inTE.Text()
                            GetResponseFromTorrentKitty(searchWord, magnetModel)
                        },
                    },
                    ListBox{
                        AssignTo: &listBox,
                        ColumnSpan: 2,
                        Model: magnetModel,
                        OnCurrentIndexChanged: func() {
                            var url string
                            if index := listBox.CurrentIndex(); index > -1 {
                                path := magnetModel.items[index].path
                                url = path
                            }
                            outTE.SetText(url)
                        },
                    },
                    TextEdit{
                        AssignTo: &outTE,
                        Font: fontYahei,
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

//随机伪造User-Agent
func GetRandomUserAgent() string {
    var userAgent = [...]string{
        "Mozilla/5.0 (compatible, MSIE 10.0, Windows NT, DigExt)",
        "Mozilla/4.0 (compatible, MSIE 7.0, Windows NT 5.1, 360SE)",
        "Mozilla/4.0 (compatible, MSIE 8.0, Windows NT 6.0, Trident/4.0)",
        "Mozilla/5.0 (compatible, MSIE 9.0, Windows NT 6.1, Trident/5.0,",
        "Opera/9.80 (Windows NT 6.1, U, en) Presto/2.8.131 Version/11.11",
        "Mozilla/4.0 (compatible, MSIE 7.0, Windows NT 5.1, TencentTraveler 4.0)",
        "Mozilla/5.0 (Windows, U, Windows NT 6.1, en-us) AppleWebKit/534.50 (KHTML, like Gecko) Version/5.1 Safari/534.50",
        "Mozilla/5.0 (Macintosh, Intel Mac OS X 10_7_0) AppleWebKit/535.11 (KHTML, like Gecko) Chrome/17.0.963.56 Safari/535.11",
        "Mozilla/5.0 (Macintosh, U, Intel Mac OS X 10_6_8, en-us) AppleWebKit/534.50 (KHTML, like Gecko) Version/5.1 Safari/534.50",
        "Mozilla/5.0 (Linux, U, Android 3.0, en-us, Xoom Build/HRI39) AppleWebKit/534.13 (KHTML, like Gecko) Version/4.0 Safari/534.13",
        "Mozilla/5.0 (iPad, U, CPU OS 4_3_3 like Mac OS X, en-us) AppleWebKit/533.17.9 (KHTML, like Gecko) Version/5.0.2 Mobile/8J2 Safari/6533.18.5",
        "Mozilla/4.0 (compatible, MSIE 7.0, Windows NT 5.1, Trident/4.0, SE 2.X MetaSr 1.0, SE 2.X MetaSr 1.0, .NET CLR 2.0.50727, SE 2.X MetaSr 1.0)",
        "Mozilla/5.0 (iPhone, U, CPU iPhone OS 4_3_3 like Mac OS X, en-us) AppleWebKit/533.17.9 (KHTML, like Gecko) Version/5.0.2 Mobile/8J2 Safari/6533.18.5",
        "MQQBrowser/26 Mozilla/5.0 (Linux, U, Android 2.3.7, zh-cn, MB200 Build/GRJ22, CyanogenMod-7) AppleWebKit/533.1 (KHTML, like Gecko) Version/4.0 Mobile Safari/533.1"}
    var r = rand.New(rand.NewSource(time.Now().UnixNano()))
    return userAgent[r.Intn(len(userAgent))]
}

//获取云播网返回的内容
func (m *FileInfoModel) GetResponseFromYunbo(searchWord string) (htmlText string) {
    url := "http://www.yunbosou.cc/s/" + searchWord +".html"
    req, _ := http.NewRequest("GET", url, nil)
    req.Header.Set("User-Agent", GetRandomUserAgent())
    client := http.DefaultClient
    res, e := client.Do(req)
    if e != nil {
        fmt.Errorf("Get请求地址%s返回错误:%s", url, e)
        return htmlText
    }
    if res.StatusCode == 200 {
        body := res.Body
        defer body.Close()
        bodyByte, _ := ioutil.ReadAll(body)
        oriText := string(bodyByte)
        //var regExpress = regexp.MustCompile(`<a[^>]+[(href)|(HREF)]\s*\t*\n*=\s*\t*\n*[(".+")|('.+')][^>]*>[^<]*</a>`)
        var regExpress = regexp.MustCompile(`<a href="/magnet/detail/(\w+)[^>]*>([^<]*)</a>`) //以Must前缀的方法或函数都是必须保证一定能执行成功的,否则将引发一次panic
        validUrls := regExpress.FindAllStringSubmatch(oriText, -1)
        for _, url := range validUrls {
            item := &FileInfo{
                Name:     url[2],
                Path:     url[1],
            }
            m.items = append(m.items, item)
            htmlText = htmlText + url[2] + "\r\nmagnet:?xt=urn:btih:" + url[1] + "\r\n\r\n"
        }
        m.PublishRowsReset()
        return htmlText
    }
    return htmlText
}


//获取torrent-kitty返回的内容
func GetResponseFromTorrentKitty(searchWord string, magnetModel *MagnetModel) bool {
    url := "https://www.torrentkitty.tv/search/" + searchWord +"/"
    req, _ := http.NewRequest("GET", url, nil)
    req.Header.Set("User-Agent", GetRandomUserAgent())
    client := http.DefaultClient
    res, e := client.Do(req)
    if e != nil {
        fmt.Errorf("Get请求地址%s返回错误:%s", url, e)
        return false
    }
    if res.StatusCode == 200 {
        body := res.Body
        defer body.Close()
        bodyByte, _ := ioutil.ReadAll(body)
        oriText := string(bodyByte)
        var regExpress = regexp.MustCompile(`<td class="name">([^<]*)</td><td class="size">.*<a href="/information/(\w+)`)
        validUrls := regExpress.FindAllStringSubmatch(oriText, -1)
        magnetModel.items = make([]MagnetItem, len(validUrls))
        for k, url := range validUrls {
            path :=  "magnet:?xt=urn:btih:" + url[2]
            magnetModel.items[k] = MagnetItem{url[1], path}
        }
        magnetModel.PublishItemsReset()
        return true
    }
    return false
}

//解析a元素 TODO:值得研究的xml函数库
func GetHref(url string) (href, content string) {
    inputReader := strings.NewReader(url)
    decoder := xml.NewDecoder(inputReader)
    for t, err := decoder.Token(); err == nil; t, err = decoder.Token() {
        switch token := t.(type) {
        // 处理元素开始（标签）
        case xml.StartElement:
            for _, attr := range token.Attr {
                attrName := attr.Name.Local
                attrValue := attr.Value
                if(strings.EqualFold(attrName, "href") || strings.EqualFold(attrName, "HREF")){
                    href = attrValue
                }
            }
        // 处理元素结束（标签）
        case xml.EndElement:
        // 处理字符数据（这里就是元素的文本）
        case xml.CharData:
            content = string([]byte(token))
        default:
            href = ""
            content = ""
        }
    }
    return href, content
}