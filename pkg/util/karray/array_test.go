package karray

import (
	"fmt"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

// TestArrayPad
//  @Description 测试切片填充数据
func TestArrayPad(t *testing.T) {
	arr := []interface{}{"a", "b"}
	//往后面添加
	newArray := ArrayPad(arr, 5, "-")
	fmt.Println(newArray)
	//往前面添加
	newArray = ArrayPad(arr, -5, "-")
	fmt.Println(newArray)

	arr2 := []interface{}{1, 2, "test"}
	newArray = ArrayPad(arr2, 5, 0)
	fmt.Println(newArray)
	//并发测试
	for i := 0; i < 100; i++ {
		go func() {
			arr := []interface{}{"A", "B"}
			newArray2 := ArrayPad(arr, 5, "-")
			fmt.Println(newArray2)
		}()
	}
	time.Sleep(1 * time.Second)
}

// TestInArray
//  @Description 测试数组、切片、map中是否存在某值
func TestInArray(t *testing.T) {
	array := [...]int{1, 2, 3, 4, 5}
	err := InArray(2, array)
	assert.Equal(t, err, nil)

	slice := []string{"demo", "测试"}
	err = InArray("demo", slice)
	assert.Equal(t, err, nil)

	mapS := make(map[int]string, 0)
	mapS[0] = "测试"
	mapS[1] = "姓名"
	err = InArray("测试", mapS)
	assert.Equal(t, err, nil)

}

// TestExplode
//  @Description 测试用指定字符分割字符串
func TestExplode(t *testing.T) {
	slice := Explode(",", "中国,北京,朝阳区")
	assert.Equal(t, slice[0], "中国")
	assert.Equal(t, slice[1], "北京")
	assert.Equal(t, slice[2], "朝阳区")
}

// TestExplodeReturnMap
//  @Description 测试用指定字符分割字符串,返回map格式
func TestExplodeReturnMap(t *testing.T) {
	maps := ExplodeReturnMap(",", "中国,北京,朝阳区")
	if v, ok := maps[0]; ok {
		assert.Equal(t, v, "中国")
	}
	if v, ok := maps[1]; ok {
		assert.Equal(t, v, "北京")
	}
	if v, ok := maps[2]; ok {
		assert.Equal(t, v, "朝阳区")
	}
}
