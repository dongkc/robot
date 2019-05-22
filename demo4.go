package main

import (
	"fmt"
	"errors"
	"log"
	"github.com/sciter-sdk/go-sciter"
	"github.com/sciter-sdk/go-sciter/window"
)

var flag_connect int
var g_root *sciter.Element
//设置元素的处理程序
func setElementHandlers(root *sciter.Element) {
	ports, err := listSerialPort()
	if err != nil {
	}

	for _, port := range ports {
		fmt.Println("port: ", port)
		root.CallFunction("add_ports", sciter.NewValue(port))
	}
}


//定义函数
func setWinHandler(w *window.Window) {
	//定义函数，在tis脚本中需要通过view对象调用
	w.DefineFunction("sendCmd", func(args ...*sciter.Value) *sciter.Value {
		if len(args) == 1 {
			send_data(args[0].String())
		} else if len(args) == 2 {
			send_data1(args[0].String(), args[1].Int())
		}

		return sciter.NewValue(1)
	})

	w.DefineFunction("openPort", func(args ...*sciter.Value) *sciter.Value {
		err := errors.New("initialize")

		g_port, err = serial_open(args[0].String())
		if err != nil {
			return sciter.NewValue(-1)
		}

		go serial_read(g_port, g_root)

		return sciter.NewValue(1)
	})

	w.DefineFunction("closePort", func(args ...*sciter.Value) *sciter.Value {
		g_port.Close()
		return sciter.NewValue(1)
	})
}

func main() {
	//创建一个新窗口
	w, err := window.New(sciter.DefaultWindowCreateFlag,
		&sciter.Rect{Left: 0, Top: 0, Right: 700, Bottom: 600})
	if err != nil {
		log.Fatal(err)
	}

	w.SetResourceArchive(resources)

	w.LoadFile("this://app/htdocs/table.htm")

	//设置标题
	w.SetTitle("法智达智能传感器数据采集工具")

	w.SetOption(sciter.SCITER_SET_DEBUG_MODE, 1)
	//获取根元素
	root, _ := w.GetRootElement()
	g_root = root

	//设置元素处理程序
	setElementHandlers(root)
	//设置窗口处理程序
	setWinHandler(w)

	//显示窗口
	w.Show()

	//运行窗口，进入消息循环
	w.Run()
}