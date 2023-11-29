package util

import (
	"bytes"
	"encoding/ascii85"
	"io"

	"github.com/klauspost/compress/zstd"
)

var encoder, _ = zstd.NewWriter(nil)
var decoder, _ = zstd.NewReader(nil)

func Compress(src []byte) []byte {
	return encoder.EncodeAll(src, make([]byte, 0, len(src)))
}

func Decompress(src []byte) ([]byte, error) {
	return decoder.DecodeAll(src, nil)
}

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
