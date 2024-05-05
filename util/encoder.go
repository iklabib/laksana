package util

import (
	"bytes"
	"encoding/ascii85"
	"io"
)


func EncodeAscii85(data []byte) string {
	var buff bytes.Buffer
	encoder := ascii85.NewEncoder(&buff)
	encoder.Write(data)
	encoder.Close()
	return buff.String()
}

func DecodeAscii85(data []byte) []byte {
	buff := bytes.NewBuffer(data)
	decoder := ascii85.NewDecoder(buff)
	decoded, _ := io.ReadAll(decoder)
	return decoded
}
