// Copyright 2011 The Walk Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package main

import (
    "log"
)

import (
    "github.com/lxn/walk"
)

func main() {
    //我们需要walk.MainWindow或者walk.Dialog来进行消息循环
    //在这个例子里面，我们让它不可见
    mw, err := walk.NewMainWindow()
    if err != nil {
        log.Fatal(err)
    }
    //从文件中读取icon
    icon, err := walk.NewIconFromFile("../resources/img/ico/4/16h.ico")
    if err != nil {
        log.Fatal(err)
    }
    //创建通知icon，并且确保退出时候清除掉它
    ni, err := walk.NewNotifyIcon()
    if err != nil {
        log.Fatal(err)
    }
    defer ni.Dispose()
    //设置icon和工具提示文本
    if err := ni.SetIcon(icon); err != nil {
        log.Fatal(err)
    }
    if err := ni.SetToolTip("小夏工具包"); err != nil {
        log.Fatal(err)
    }
    //当鼠标左键点击时，出现气球
    ni.MouseDown().Attach(func(x, y int, button walk.MouseButton) {
        if button != walk.LeftButton {
            return
        }
        if err := ni.ShowCustom("通知", "早睡早起，方能以弱胜强"); err != nil {
            log.Fatal(err)
        }
    })

    /**
     *---------------------------------------
     * 菜单区
     *---------------------------------------
     */
    //发现视频功能
    searchAction := walk.NewAction()
    if err := searchAction.SetText("发现视频"); err != nil {
        log.Fatal(err)
    }
    searchAction.Triggered().Attach(func() { walk.App().Exit(0) })
    if err := ni.ContextMenu().Actions().Add(searchAction); err != nil {
        log.Fatal(err)
    }
    //在内容菜单里添加一个退出动作
    exitAction := walk.NewAction()
    if err := exitAction.SetText("退出软件"); err != nil {
        log.Fatal(err)
    }
    exitAction.Triggered().Attach(func() { walk.App().Exit(0) })
    if err := ni.ContextMenu().Actions().Add(exitAction); err != nil {
        log.Fatal(err)
    }

    //通知icon初始是隐藏的，我们需要让他显示出来
    if err := ni.SetVisible(true); err != nil {
        log.Fatal(err)
    }
    //icon显示出来后，将气泡显示出来
    if err := ni.ShowInfo("通知示范", "点击icon再次展现."); err != nil {
        log.Fatal(err)
    }
    //开始消息循环
    mw.Run()
}
