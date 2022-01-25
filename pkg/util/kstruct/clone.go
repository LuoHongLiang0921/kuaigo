// @Description struct 复制

package kstruct

import "reflect"

// CloneStruct
//  @Description  克隆结构体
//  @Param src
//  @Param dst
func CloneStruct(src, dst interface{}) {
	srcVal := reflect.ValueOf(src).Elem()
	dstVal := reflect.ValueOf(dst).Elem()

	for i := 0; i < srcVal.NumField(); i++ {
		value := srcVal.Field(i)
		name := srcVal.Type().Field(i).Name

		dstValue := dstVal.FieldByName(name)
		if dstValue.IsValid() == false {
			continue
		}
		dstValue.Set(value)
	}
}
