package util

import (
	"bytes"
	"compress/gzip"
	"encoding/ascii85"
	"io"
)

func Compress(s []byte) ([]byte, error) {
	buff := bytes.Buffer{}
	compressed := gzip.NewWriter(&buff)
	if _, err := compressed.Write(s); err != nil {
		return nil, err
	}

	compressed.Close()
	return buff.Bytes(), nil
}

func Decompress(s []byte) ([]byte, error) {
	r, err := gzip.NewReader(bytes.NewReader(s))
	if err != nil {
		return nil, err
	}

	buff, err := io.ReadAll(r)
	if err != nil {
		return nil, err
	}
	r.Close()
	return buff, nil
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
