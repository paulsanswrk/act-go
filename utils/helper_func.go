package utils

import (
	"bytes"
	"compress/gzip"
	"io"
	"strconv"
)

func First[T, U any](val T, _ U) T {
	return val
}

func StrToFloat64(s string) float64 {
	res, _ := strconv.ParseFloat(s, 64)
	return res
}

func DecodeGzip(data []byte) (string, error) {
	reader, err := gzip.NewReader(bytes.NewReader(data))
	if err != nil {
		return "", err
	}
	defer reader.Close()

	decodedMsg, err := io.ReadAll(reader)
	if err != nil {
		return "", err
	}

	return string(decodedMsg), nil
}
