package scripts

import (
	"github.com/yanyiwu/gojieba"
	"strings"
	"os/exec"
)

func GetQueryString(question string) string {
	var (
		words []string
		questionWords []string
		searchStr string
	)
	x := gojieba.NewJieba()
	defer x.Free()
	words = x.Tag(question)
	//词性: 1799/m,年/m,乾隆皇帝/nr,去世/t,，/x,同年/t,去世/t,的/uj,另外/c,一位/m,美国/ns,总统/n,是/v,谁/r
	for _, v := range words {
		wordSplice := strings.Split(v, "/")
		if wordSplice[1] == v || wordSplice[1] == "uj" || wordSplice[1] == "x" || wordSplice[1] == "c"  {
			continue
		}
		questionWords = append(questionWords, wordSplice[0])
	}
	searchStr = strings.Join(questionWords, "+")
	return searchStr
}

//打开系统默认浏览器(目前只支持WINDOWS)
func OpenUrl(uri string) error {
	//l2, err2 := url.ParseRequestURI(uri)
	//fmt.Println(l2, err2)
	//fmt.Println(l2.RawPath)
	cmd := exec.Command("cmd", "/C", "start", uri)
	return cmd.Start()
}

