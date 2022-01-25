package karray

import (
	"bytes"
	"errors"
	"math"
	"math/rand"
	"reflect"
	"strings"
	"time"
)

// InArray
//  @Description 查找切面、数组、或map中是否存在某个值(Support types: slice, array or map)
//  @Param needle 查找的数值
//  @Param hayStack 被查找的集合
func InArray(needle interface{}, hayStack interface{}) error {
	val := reflect.ValueOf(hayStack)
	switch val.Kind() {
	case reflect.Slice, reflect.Array:
		for i := 0; i < val.Len(); i++ {
			if reflect.DeepEqual(needle, val.Index(i).Interface()) {
				return nil
			}
		}
	case reflect.Map:
		for _, k := range val.MapKeys() {
			if reflect.DeepEqual(needle, val.MapIndex(k).Interface()) {
				return nil
			}
		}
	default:
		errors.New("hayStack type must be slice, array or map")
	}

	return errors.New("unknown error")
}

// ArrayDiff
//  @Description 比较两个集合的值，并返回差集
//  @Param a1 集合1
//  @Param a2 集合2
//  @Return 被比较的两个集合差集
func ArrayDiff(a1 []interface{}, a2 ...[]interface{}) ([]interface{}, error) {
	if len(a1) == 0 {
		return []interface{}{}, nil
	}
	if len(a1) > 0 && len(a2) == 0 {
		return a1, nil
	}
	var tmp = make(map[interface{}]int, len(a1))
	for _, v := range a1 {
		tmp[v] = 1
	}
	for _, param := range a2 {
		for _, arg := range param {
			if tmp[arg] != 0 {
				tmp[arg]++
			}
		}
	}
	var res = make([]interface{}, 0, len(tmp))
	for k, v := range tmp {
		if v == 1 {
			res = append(res, k)
		}
	}
	return res, nil
}

// ArrayFlip
//  @Description 反转Map的键和值
//  @Param m 被反转的map
//  @Return 反传后的map
func ArrayFlip(m map[interface{}]interface{}) map[interface{}]interface{} {
	n := make(map[interface{}]interface{})
	for i, v := range m {
		n[v] = i
	}
	return n
}

// ArrayKeys
//  @Description 获取指定Map的所有键集合
//  @Param m map集合
//  @Return 键数组
func ArrayKeys(m map[interface{}]interface{}) []interface{} {
	i, keys := 0, make([]interface{}, len(m))
	for key := range m {
		keys[i] = key
		i++
	}
	return keys
}

// ArrayValues
//  @Description 获取指定Map的所有值集合
//  @Param m map集合
//  @Return 值数组
func ArrayValues(m map[interface{}]interface{}) []interface{} {
	i, values := 0, make([]interface{}, len(m))
	for _, value := range m {
		values[i] = value
		i++
	}
	return values
}

// ArrayMerge
//  @Description 合并集合
//  @Param  as 集合列表
//  @Return 合并后的集合
func ArrayMerge(as ...[]interface{}) []interface{} {
	n := 0
	for _, v := range as {
		n += len(v)
	}
	s := make([]interface{}, 0, n)
	for _, v := range as {
		s = append(s, v...)
	}
	return s
}

// ArrayChunk
//  @Description 将一个数组分割成多个
//  @Param  a 被分割的数组
//  @Return 分割后的集合
func ArrayChunk(a []interface{}, size int) ([][]interface{}, error) {
	if size < 1 {
		return nil, errors.New("size: cannot be less than 1")
	}
	length := len(a)
	chunks := int(math.Ceil(float64(length) / float64(size)))
	var n [][]interface{}
	for i, end := 0, 0; chunks > 0; chunks-- {
		end = (i + 1) * size
		if end > length {
			end = length
		}
		n = append(n, a[i*size:end])
		i++
	}
	return n, nil
}

// ArrayPad
//  @Description 以指定长度将一个值填充进数组
//  @Param  a 需要被填充的原始数组
//  @Param  size 新数组的长度
//  @Return 返回 array 用 value 填充到 size 指定的长度之后的一个副本。 如果 size 为正，则填补到数组的右侧，如果为负则从左侧开始填补。 如果 size 的绝对值小于或等于 array 数组的长度则没有任何填补
func ArrayPad(a []interface{}, size int, value interface{}) []interface{} {
	if size == 0 || (size > 0 && size < len(a)) || (size < 0 && size > -len(a)) {
		return a
	}
	n := size
	if size < 0 {
		n = -size
	}
	n -= len(a)
	tmp := make([]interface{}, n)
	for i := 0; i < n; i++ {
		tmp[i] = value
	}
	if size > 0 {
		return append(a, tmp...)
	}
	return append(tmp, a...)
}

// ArraySlice
//  @Description 从数组中取出一段
//  @Param  a 输入的数组
//  @Param  offset 位置
//  @Param  length 获取长度
//  @Return 返回指定的数组段集合
func ArraySlice(a []interface{}, offset, length uint) ([]interface{}, error) {
	if offset > uint(len(a)) {
		return nil, errors.New("offset: the offset is less than the length of s")
	}
	end := offset + length
	if end < uint(len(a)) {
		return a[offset:end], nil
	}
	return a[offset:], nil
}

// ArrayRand
//  @Description 打乱数组
//  @Param  a 输入的数组
//  @Return 打乱后的数组
func ArrayRand(a []interface{}) []interface{} {
	r := rand.New(rand.NewSource(time.Now().UnixNano()))
	n := make([]interface{}, len(a))
	for i, v := range r.Perm(len(a)) {
		n[i] = a[v]
	}
	return n
}

// ArrayColumn
//  @Description 获取二维Map的二维key对应的值集合
//  @Param  m 输入的二维Map
//  @Param  columnKey  二维key
//  @Return 二维key对应的值集合
func ArrayColumn(m map[string]map[string]interface{}, columnKey string) []interface{} {
	columns := make([]interface{}, 0, len(m))
	for _, val := range m {
		if v, ok := val[columnKey]; ok {
			columns = append(columns, v)
		}
	}
	return columns
}

// ArrayPush
//  @Description 将一个或多个单元压入数组的末尾（入栈）
//  @Param  a 输入的数组
//  @Param  elements  要压入的元素集合
//  @Return 压入后数组的长度
func ArrayPush(a *[]interface{}, elements ...interface{}) int {
	*a = append(*a, elements...)
	return len(*a)
}

// ArrayPop
//  @Description 弹出数组最后一个单元（出栈）
//  @Param  a 需要弹出栈的数组
//  @Return 返回数组的最后一个值
func ArrayPop(a *[]interface{}) interface{} {
	if len(*a) == 0 {
		return nil
	}
	ep := len(*a) - 1
	e := (*a)[ep]
	*a = (*a)[:ep]
	return e
}

// ArrayUnshift
//  @Description 将一个或多个单元压入数组的头部（入栈）
//  @Param  a 输入的数组
//  @Param  elements  要压入的元素集合
//  @Return 压入后数组的长度
func ArrayUnshift(s *[]interface{}, elements ...interface{}) int {
	*s = append(elements, *s...)
	return len(*s)
}

// ArrayShift
//  @Description 将数组开头的单元移出数组
//  @Param  a 需要弹出栈的数组
//  @Return 返回移出的顶部的值
func ArrayShift(s *[]interface{}) interface{} {
	if len(*s) == 0 {
		return nil
	}
	f := (*s)[0]
	*s = (*s)[1:]
	return f
}

// ArrayKeyExists
//  @Description 判断Map中否存在指定的key
//  @Param  key 判断的key值
//  @Param m 输入Map
//  @Return 判断结果
func ArrayKeyExists(key interface{}, m map[interface{}]interface{}) bool {
	_, ok := m[key]
	return ok
}

// ArrayCombine
//  @Description 创建一个Map，用一个数组的值作为其键名，另一个数组的值作为其值
//  @Param  a1 输入数组1
//  @Param  a2 输入数组2
//  @Return 创建的map
func ArrayCombine(a1, a2 []interface{}) (map[interface{}]interface{}, error) {
	if len(a1) != len(a2) {
		return nil, errors.New("the number of elements for each slice isn't equal")
	}
	m := make(map[interface{}]interface{}, len(a1))
	for i, v := range a1 {
		m[v] = a2[i]
	}
	return m, nil
}

// ArrayReverse
//  @Description 返回单元顺序相反的数组
//  @Param  a 输入数组
//  @Return 创建的map
func ArrayReverse(a []interface{}) []interface{} {
	for i, j := 0, len(a)-1; i < j; i, j = i+1, j-1 {
		a[i], a[j] = a[j], a[i]
	}
	return a
}

// Implode
//  @Description 用指定字符连接数组元素
//  @Param  s 连接的字符
//  @Param  a 输入数组
//  @Return 连接后的字符串
func Implode(s string, a []string) string {
	var buf bytes.Buffer
	l := len(a)
	for _, str := range a {
		buf.WriteString(str)
		if l--; l > 0 {
			buf.WriteString(s)
		}
	}
	return buf.String()
}

// Explode
//  @Description 用指定字符分割字符串
//  @Param  separator 分隔符
//  @Param  str 要分割的字符串
//  @Return 切片
func Explode(separator string, str string) []string {
	return strings.Split(str, separator)
}

// ExplodeReturnMap
//  @Description 用指定字符分割字符串,返回map格式
//  @Param  separator 分隔符
//  @Param  str 要分割的字符串
//  @Return map[int]string
func ExplodeReturnMap(separator string, str string) map[int]string {
	slice := strings.Split(str, separator)
	mapS := make(map[int]string, 0)
	for i, s := range slice {
		mapS[i] = s
	}
	return mapS
}
