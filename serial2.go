package main

import (
	// "fmt"
	"time"
	// "encoding/hex"
	// "encoding/binary"
	// "github.com/sciter-sdk/go-sciter/window"
	"github.com/sciter-sdk/go-sciter"
	"go.bug.st/serial.v1"
	// "go.bug.st/serial.v1/enumerator"
	// "github.com/bugst/go-serial"
	// "github.com/bugst/go-serial/enumerator"
)

// 全局变量,用来保存选定的串口
var g_port2 serial.Port

func serial_open2(name string) (serial.Port, error) {
	mode := &serial.Mode{
		BaudRate: 115200,
		Parity:   serial.NoParity,
		DataBits: 8,
		StopBits: serial.OneStopBit,
	}

	return serial.Open(name, mode)
}

func serial_read2(port serial.Port, root *sciter.Element) {
	out := make(chan []byte, 32)
	stop := make(chan string)

   send2 := func() { send_data2("") }

    stop_send := schedule(send2, 20*time.Millisecond)
	// 启动接收串口协程
	go recv2(g_port2, out, stop)

	// 接收串口数据帧,并做处理
	for {
		select {
		case msg := <-out:
			process2(root, msg)
		case <-stop:
			stop_send <- true;

			return
		}
	}
}

func recv2(port serial.Port, out chan <- []byte, stop chan <- string) {
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
		if frame[0] != 0xAA {
			for i := 0; i < len(frame); i++ {
				if (frame[i] == 0xAA) {
					frame = frame[i:len(frame)]
					break;
				}
			}
		}

		if len(frame) < 7 {
			// 收到的数据不够一桢
			continue
		}

		body_len := frame[2]
		if len(frame) < int(body_len + 4) {
			// fmt.Println("len: ", len(frame), " body_len ", int(body_len + 10))
			continue
		}

		// fmt.Println("len: ", body_len, " frame: ", hex.EncodeToString(frame[:body_len + 7]))
			// 发送出去
		out <- frame[:body_len + 5]

		// fmt.Println("frame2: ", hex.EncodeToString(frame[:body_len + 5]))

		// 处理下一桢数据
		frame = frame[body_len + 4 :]
	}
}

func schedule(what func(), delay time.Duration) chan bool {
    stop := make(chan bool)

    go func() {
        for {
            what()
            select {
            case <-time.After(delay):
            case <-stop:
                return
            }
        }
    }()

    return stop
}
