// @description
// @author yixia
// Copyright 2021 sndks.com. All rights reserved.
// @datetime 2021/1/14 5:21 下午
// @lastmodify 2021/1/14 5:21 下午

package conf

// GetOption ...
type (
	GetOption  func(o *GetOptions)
	GetOptions struct {
		TagName string
	}
)

var defaultGetOptions = GetOptions{
	TagName: "mapstructure",
}

// 设置Tag
func TagName(tag string) GetOption {
	return func(o *GetOptions) {
		o.TagName = tag
	}
}
