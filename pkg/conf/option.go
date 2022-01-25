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

// TagName
//  @Description  设置Tag
//  @Param tag
//  @Return GetOption
func TagName(tag string) GetOption {
	return func(o *GetOptions) {
		o.TagName = tag
	}
}
