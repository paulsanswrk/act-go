package utils

import (
	"bytes"
	"compress/gzip"
	"crypto/hmac"
	"crypto/sha256"
	"encoding/hex"
	"encoding/json"
	"github.com/google/uuid"
	"github.com/mitchellh/mapstructure"
	"io"
	"math"
	"net/http"
	"strconv"
	"strings"
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

func ComputeHmac256(strMessage string, strSecret string) string {
	key := []byte(strSecret)
	h := hmac.New(sha256.New, key)
	h.Write([]byte(strMessage))
	return hex.EncodeToString(h.Sum(nil))
}

func HTTP_Request(url string, method string, headers map[string]string, dest interface{}) (response_string string, err error) {
	var req *http.Request
	req, err = http.NewRequest(method, url, nil)
	if err != nil {
		return
	}
	client := &http.Client{}
	for k, v := range headers {
		req.Header.Add(k, v)
	}

	var res *http.Response
	res, err = client.Do(req)
	if err != nil {
		return
	}

	var body []byte
	body, err = io.ReadAll(res.Body)
	if err != nil {
		return
	}

	response_string = string(body)

	if dest == nil {
		return
	}

	if err = json.Unmarshal(body, &dest); err != nil {
		return
	}

	return
}

func Struct_To_Map(input interface{}) (output map[string]interface{}, err error) {
	config := &mapstructure.DecoderConfig{
		Result:  &output,
		TagName: "json",
	}
	decoder, _ := mapstructure.NewDecoder(config)
	err = decoder.Decode(input)
	return
}

func Guid() string {
	return strings.Replace(uuid.Must(uuid.NewRandom()).String(), "-", "", -1)
}

func Range(from float64, to float64, step float64) (result []float64) {
	//result := make([]float64, 0)
	step = math.Abs(step)

	cnt := math.Floor(math.Abs((to-from)/step)) + 1

	if to < from {
		step = -step
	}

	for i := range int(cnt) {
		result = append(result, from+float64(i)*step)
	}
	return
}
