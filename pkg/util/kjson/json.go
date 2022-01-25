package kjson

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"strings"
)

// Encode
//  @Description:对变量进行 JSON 编码 返回byte
//  @Param v
//  @Return []byte
//  @Return error
func Encode(v interface{}) ([]byte, error) {
	return json.Marshal(v)
}

// EncodeToString
//  @Description:对变量进行 JSON 编码 返回字符串
//  @Param v
//  @Return string
//  @Return error
func EncodeToString(v interface{}) (string, error) {
	jsons, err := json.Marshal(v)
	return string(jsons), err
}

// Decode
//  @Description: 对 JSON byte 进行解码
//  @Param data
//  @Param v
//  @Return error
func Decode(data []byte, v interface{}) error {
	decoder := json.NewDecoder(bytes.NewReader(data))
	if err := decodeUseNumber(decoder, v); err != nil {
		return formatError(string(data), err)
	}

	return nil
}

// DecodeToMap
//  @Description:对 JSON 格式的字符串进行解码 并返回map
//  @Param data
//  @Return map[string]interface{}
//  @Return error
func DecodeToMap(data string) (map[string]interface{}, error) {
	var dat map[string]interface{}
	err := json.Unmarshal([]byte(data), &dat)
	return dat, err
}

// DecodeFromString
// 	@Description 对字符串进行JSON接码
//	@Param str
//	@Param v
// 	@Return error
func DecodeFromString(str string, v interface{}) error {
	decoder := json.NewDecoder(strings.NewReader(str))
	if err := decodeUseNumber(decoder, v); err != nil {
		return formatError(str, err)
	}

	return nil
}

// DecodeFromReader
// 	@Description unmarshals v from reader.
//	@Param reader
//	@Param v
// 	@Return error
func DecodeFromReader(reader io.Reader, v interface{}) error {
	var buf strings.Builder
	teeReader := io.TeeReader(reader, &buf)
	decoder := json.NewDecoder(teeReader)
	if err := decodeUseNumber(decoder, v); err != nil {
		return formatError(buf.String(), err)
	}

	return nil
}

func decodeUseNumber(decoder *json.Decoder, v interface{}) error {
	decoder.UseNumber()
	return decoder.Decode(v)
}

func formatError(v string, err error) error {
	return fmt.Errorf("string: `%s`, error: `%s`", v, err.Error())
}
