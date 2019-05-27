package main

import (
	"fmt"
	"log"
	"os"
	"path"
	// "time"
	// "math"
	// "encoding/hex"
	"encoding/binary"
	"github.com/sciter-sdk/go-sciter"
)

func send_cmd2(cmd []byte) {
	_, err := g_port2.Write(cmd)
	if err != nil {
		fmt.Println("Write to serial port failed")
		return
	}

	// fmt.Println("send: ", hex.EncodeToString(cmd))
}

func send_data2(cmd1 string) {
	frame := []byte {0xBB, 0x55, 0x06, 0x01, 0x04, 0x00, 0x00, 0x00, 0x00, 0x34}

	send_cmd2(frame)
}

func process2(root *sciter.Element, data []byte) {
	// fmt.Println("msg: ", hex.EncodeToString(data))

	if len(data) < 10 {
		return
	}

	// addr := binary.LittleEndian.Uint32(data)
	switch data[24] {
	case 0x04:
		if data[25] != 0x07 {
			return;
		}

		angel   := float32(int16(binary.LittleEndian.Uint16(data[26:28]))) * 180 / 32768
		angel_a := float32(int16(binary.LittleEndian.Uint16(data[28:30]))) * 2000 / 32768
		gyr_z   := float32(int16(binary.LittleEndian.Uint16(data[30:32]))) * 16 / 32768
		s := fmt.Sprintf("%d %.0f %.2f %.2f\n", NowAsUnixMilli(), angel, angel_a, gyr_z)

		write_data("sensor1.dat", s)

		root.CallFunction("sensor_report",
			sciter.NewValue(s))
	}
}

func write_data(name string, s string) {
	aaa, _ := os.Executable()
	data_1_path := path.Dir(aaa)+ "/" + name
	f, err := os.OpenFile(data_1_path, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
	if err != nil {
		log.Fatal(err)
	}
	defer f.Close()

	f.WriteString(s)
}