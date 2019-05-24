package main

import (
	"fmt"
	// "encoding/hex"
	// "encoding/binary"
	// "github.com/sciter-sdk/go-sciter/window"
	"github.com/sciter-sdk/go-sciter"

	"go.bug.st/serial.v1"
	"go.bug.st/serial.v1/enumerator"
	// "github.com/bugst/go-serial"
	// "github.com/bugst/go-serial/enumerator"
)

// 全局变量,用来保存选定的串口
var g_port serial.Port

func listSerialPort() ([]string, error) {
	names := make([]string, 0)

	ports, err := enumerator.GetDetailedPortsList()
	if err != nil {
		return names, err
	}

	if len(ports) == 0 {
		fmt.Println("No serial ports found!")
		return names, err
	}

	for _, port := range ports {
		fmt.Printf("Found port: %s\n", port.Name)
		if port.IsUSB {
			fmt.Printf("   USB ID     %s:%s\n", port.VID, port.PID)
			fmt.Printf("   USB serial %s\n", port.SerialNumber)
		}

		names = append(names, port.Name)
	}

	return names, err
}

func serial_open(name string) (serial.Port, error) {
	mode := &serial.Mode{
		BaudRate: 115200,
		Parity:   serial.NoParity,
		DataBits: 8,
		StopBits: serial.OneStopBit,
	}

	return serial.Open(name, mode)
}

func serial_read(port serial.Port, root *sciter.Element) {
	out := make(chan []byte, 32)
	stop := make(chan string)

	// 启动接收串口协程
	go recv(g_port, out, stop)

	// 接收串口数据帧,并做处理
	for {
		select {
		case msg := <-out:
			process(root, msg)
		case <-stop:
			return
		}
	}
}

func recv(port serial.Port, out chan <- []byte, stop chan <- string) {
	frame := make([]byte, 0)
	buf := make([]byte, 10000000)

	for {
		// time.Sleep(3000* time.Millisecond)
		n, err := port.Read(buf)
		if err != nil {
			stop <- "stop"
			return
		}

		frame = append(frame, buf[:n]...)

		// 分桢 桢开头处理
		if frame[0] != 0x3A {
			for i := 0; i < len(frame); i++ {
				if (frame[i] == 0x3A) {
					frame = frame[i:len(frame)]
					break;
				}
			}
		}

		if len(frame) < 6 {
			// 收到的数据不够一桢
			continue
		}

		body_len := frame[3]
		if len(frame) < int(body_len + 7) {
			// fmt.Println("len: ", len(frame), " body_len ", int(body_len + 10))
			continue
		}

		// fmt.Println("len: ", body_len, " frame: ", hex.EncodeToString(frame[:body_len + 7]))
		// 检查桢尾
		if frame[body_len + 5] == 0x0D &&
			frame[body_len + 6] == 0x0A {
			// 发送出去
			out <- frame[:body_len + 7]
		}
		// fmt.Println("frame2: ", hex.EncodeToString(frame[:body_len + 11]))

		// 处理下一桢数据
		frame = frame[body_len + 6 :]
	}
}