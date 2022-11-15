package util

import (
	"crypto"
	"encoding/binary"
	"io"
	"log"
	"os"

	"github.com/zeebo/xxh3"
)

func GetFileMD5(path string) uint64 {
	rom, err := os.Open(path)
	if err != nil {
		log.Fatal(err)
	}

	hash := crypto.MD5.New()
	_, err = io.Copy(hash, rom)

	if err != nil {
		log.Fatal(err)
	}

	return binary.LittleEndian.Uint64(hash.Sum(nil))
}

func GetFileXXH3(path string) uint64 {
	rom, err := os.Open(path)
	if err != nil {
		log.Fatal(err)
	}

	hash := xxh3.New()
	_, err = io.Copy(hash, rom)

	if err != nil {
		log.Fatal(err)
	}

	//bad
	return hash.Sum64()

}
