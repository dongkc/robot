package main

import (
	"fmt"
	"time"
	// "math"
	"encoding/hex"
	"encoding/binary"
	"github.com/sciter-sdk/go-sciter"
)

func NowAsUnixMilli() int64 {
	t := time.Date(2019, time.May, 27, 0, 0, 0, 0, time.UTC)
    return (time.Now().UnixNano() - t.UnixNano()) / 1e6
}

func parse(buf []byte) {
	fmt.Println("buf: ", hex.EncodeToString(buf))
}

func send_cmd(cmd []byte) {
	_, err := g_port.Write(cmd)
	if err != nil {
		fmt.Println("Write to serial port failed")
		return
	}

	fmt.Println("send: ", hex.EncodeToString(cmd))
}

func send_data(cmd1 string) {
	frame := []byte {0x0A, 0x01, 0x28, 0x00, 0xDD, 0xEE}

	switch cmd1 {
	case "set0":
		frame[2] = 0x20
	case "set1":
		frame[2] = 0x21
	case "set2":
		frame[2] = 0x22
	case "set3":
		frame[2] = 0x23
	case "cfg":
		frame[2] = 0x17
	case "mode_config":
		frame[2] = 0x52
	case "mode_test":
		frame[2] = 0x51
	case "mode_work":
		frame[2] = 0x50
	case "start_collect":
		frame = []byte {0x0A, 0x04, 0x60, 0x01, 0x02, 0x03, 0x00, 0xDD, 0xEE}
	case "report_confirm":
		frame = []byte {0x0A, 0x04, 0x61, 0x01, 0x02, 0x03, 0x00, 0xDD, 0xEE}
	}

	send_cmd(frame)
}

func send_data1(cmd1 string, val int) {
	frame := []byte {0x0A, 0x03, 0x28, 0x00, 0x00, 0x00, 0xDD, 0xEE}
	binary.LittleEndian.PutUint16(frame[3:], uint16(val))

	switch cmd1 {
	case "set0":
		frame[2] = 0x24
	case "set1":
		frame[2] = 0x25
	case "set2":
		frame[2] = 0x26
	case "set3":
		frame[2] = 0x27
	case "addr":
		frame = []byte {0x0A, 0x04, 0x10, 0x00, 0x00, 0x00, 0x00, 0xDD, 0xEE}
		binary.LittleEndian.PutUint32(frame[3:], uint32(val))
	case "freq":
		frame = []byte {0x0A, 0x04, 0x11, 0x00, 0x00, 0x00, 0x00, 0xDD, 0xEE}
		binary.LittleEndian.PutUint32(frame[3:], uint32(val))
	case "power":
		frame = []byte {0x0A, 0x02, 0x12, 0x00, 0x00, 0xDD, 0xEE}
		frame[3] = uint8(val)
	case "tmv":
		frame = []byte {0x0A, 0x02, 0x13, 0x00, 0x00, 0xDD, 0xEE}
		frame[3] = uint8(val)
	case "loose1":
		frame = []byte {0x0A, 0x02, 0x14, 0x00, 0x00, 0xDD, 0xEE}
		frame[3] = uint8(val)
	case "loose2":
		frame = []byte {0x0A, 0x02, 0x15, 0x00, 0x00, 0xDD, 0xEE}
		frame[3] = uint8(val)
	case "opressure":
		frame = []byte {0x0A, 0x02, 0x16, 0x00, 0x00, 0xDD, 0xEE}
		frame[3] = uint8(val)
	}

	send_cmd(frame)
}

func process(root *sciter.Element, data []byte) {
	// fmt.Println("msg: ", hex.EncodeToString(data))

	if len(data) < 7 {
		return
	}

	// addr := binary.LittleEndian.Uint32(data)
	switch data[1] {
	case 0x0b:
		buf := make([]byte, 4)
		// addr
		copy(buf, data[7:11])
		
		// gyr_x := math.Float32frombits(binary.LittleEndian.Uint32(data[11:15]))
		// gyr_y := math.Float32frombits(binary.LittleEndian.Uint32(data[15:19]))
		// gyr_z := math.Float32frombits(binary.LittleEndian.Uint32(data[19:23]))

		angel   := int(int16(binary.LittleEndian.Uint16(data[4:6])))
		angel_a := float32(int16(binary.LittleEndian.Uint16(data[6:8])))
		gyr_x   := float32(int16(binary.LittleEndian.Uint16(data[8:10])))
		gyr_y   := float32(int16(binary.LittleEndian.Uint16(data[10:12])))
		gyr_z   := float32(int16(binary.LittleEndian.Uint16(data[12:14])))
		fmt.Printf("%d %d %.2f %.2f %.2f %.2f\n", NowAsUnixMilli(), angel / 100, angel_a / 50, gyr_x / 1000, gyr_y / 1000, gyr_z / 1000)

		// root.CallFunction("sensor_work_report",
		// 	sciter.NewValue(int(addr)),
		// 	sciter.NewValue(int(force)))
	}
}
