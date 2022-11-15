package util

import (
	"encoding/binary"
	"log"
	"os"
)

func ReadUint8(file *os.File, offset int) uint8 {
	file.Seek(int64(offset), 0)
	data := make([]byte, 1)
	_, err := file.Read(data)

	if err != nil {
		log.Fatal(err)
	}

	return data[0]
}

func ReadUint16(file *os.File, offset int) uint16 {
	file.Seek(int64(offset), 0)
	data := make([]byte, 2)
	_, err := file.Read(data)

	if err != nil {
		log.Fatal(err)
	}

	return binary.BigEndian.Uint16(data)

}

func ReadUint32(file *os.File, offset int) uint32 {
	file.Seek(int64(offset), 0)
	data := make([]byte, 4)
	_, err := file.Read(data)

	if err != nil {
		log.Fatal(err)
	}

	return binary.BigEndian.Uint32(data)
}

func ReadNullTerminatedString(file *os.File, offset int) string {
	str := ""
	index := offset

	for {
		ch := ReadUint8(file, index)
		if ch == 0 {
			break
		}
		str = str + string(ch)
		index++
	}
	return str
}

func ReadData(file *os.File, offset int, size int) []byte {
	var data = make([]byte, size)
	for i := 0; i < size; i++ {
		data[i] = ReadUint8(file, offset+i)
	}
	return data
}
