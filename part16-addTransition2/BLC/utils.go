package BLC

import (
	"bytes"
	"encoding/binary"
	"encoding/json"
	"fmt"
	"log"
	"os"
)

func IntToHex(num int64) []byte {
	buff := new(bytes.Buffer)
	err := binary.Write(buff, binary.BigEndian, num)
	if err != nil {
		log.Panic(err)
	}
	return buff.Bytes()
}

// 将json转换为array
func JSONtoArray(jsonStr string) []string {
	log.Println(jsonStr)
	var array []string
	if err := json.Unmarshal([]byte(jsonStr), &array); err != nil {
		fmt.Println("json To array error...")
		log.Panic(err)
		os.Exit(1)
	}

	return array

}
