package main

import (
	"fmt"
	"math"
	// "encoding/hex"
	"encoding/binary"
	// "github.com/sciter-sdk/go-sciter/window"
	"github.com/sciter-sdk/go-sciter"

	"go.bug.st/serial.v1"
	// "go.bug.st/serial.v1/enumerator"
	// "github.com/bugst/go-serial"
	// "github.com/bugst/go-serial/enumerator"
)

// 全局变量,用来保存选定的串口
var g_port_lpms9 serial.Port

func serial_lpms9_open(name string) (serial.Port, error) {
	mode := &serial.Mode{
		BaudRate: 921600,
		Parity:   serial.NoParity,
		DataBits: 8,
		StopBits: serial.OneStopBit,
	}

	return serial.Open(name, mode)
}

func serial_lpms9_read(port serial.Port, root *sciter.Element) {
	out := make(chan []byte, 32)
	stop := make(chan string)

	// 启动接收串口协程
	go recv_lpms9(g_port_lpms9, out, stop)

	// 接收串口数据帧,并做处理
	for {
		select {
		case msg := <-out:
			process_lpms9(root, msg)
		case <-stop:
			return
		}
	}
}

func recv_lpms9(port serial.Port, out chan <- []byte, stop chan <- string) {
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

		if len(frame) < 11 {
			// 收到的数据不够一桢
			continue
		}

		body_len := binary.LittleEndian.Uint16(frame[5:7])
		if len(frame) < int(body_len + 11) {
			// fmt.Println("len: ", len(frame), " body_len ", int(body_len + 10))
			continue
		}

		// fmt.Println("len: ", body_len, " frame: ", hex.EncodeToString(frame[:body_len + 7]))
		// 检查桢尾
		if frame[body_len + 9] == 0x0D &&
			frame[body_len + 10] == 0x0A {
			// 发送出去
			out <- frame[:body_len + 11]
		}
		// fmt.Println("frame2: ", hex.EncodeToString(frame[:body_len + 11]))

		// 处理下一桢数据
		frame = frame[body_len + 11 :]
	}
}

func process_lpms9(root *sciter.Element, data []byte) {
	// fmt.Println("msg: ", hex.EncodeToString(data))
	if len(data) < 11 {
		return
	}

	switch data[3] {
	case 0x09:
		angel   := math.Float32frombits(binary.LittleEndian.Uint32(data[71:75])) * 180 / math.Pi
		angel_a := math.Float32frombits(binary.LittleEndian.Uint32(data[83:87]))
		gyr_z   := math.Float32frombits(binary.LittleEndian.Uint32(data[31:35]))

		s := fmt.Sprintf("%d %.0f %.2f %.2f\n", NowAsUnixMilli(), angel, angel_a, gyr_z)
		write_data("sensor_lpms9.dat", s)

		root.CallFunction("sensor_lpms9_report",
			sciter.NewValue(s))
	}
}
