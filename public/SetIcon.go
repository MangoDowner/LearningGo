package public

import (
	"github.com/lxn/walk"
)

/**
  设置软件图标
 */
func SetIcon(mw *walk.MainWindow, file_path string) {
	if file_path == "" {
		file_path = "../resources/img/ico/4/16h.ico";
	}
	ic, err := walk.NewIconFromFile(file_path)
	if err != nil {
		return
	}
	mw.SetIcon(ic)
}
