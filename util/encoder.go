package util

import (
	"bytes"
	"compress/gzip"
	"io"
  "os"
)

func CreateTempDir() (string, error) {
	tempDir, err := os.MkdirTemp("", "box_")
	if err != nil {
    return "", err
	}
  return tempDir, nil
}

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
