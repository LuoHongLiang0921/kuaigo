package kcolor

import (
	"fmt"
	"testing"
)

func TestSetColor(t *testing.T) {
	fmt.Println(Red("我是红色"))
	fmt.Println(Green("我是绿色"))
	fmt.Println(Blue("我是蓝色"))
	fmt.Println(Cyan("我是青蓝色"))
	fmt.Println(Magenta("我是紫红色"))
	fmt.Println(White("我是白色"))
	fmt.Println(Redf("我是红色格式化","文本"))

	fmt.Println(Error("我是错误"))
	fmt.Println(Debug("我是调试"))
	fmt.Println(Warn("我是警告"))
	for b := 40; b <= 47; b++ { // 背景色彩 = 40-47
		for f := 30; f <= 37; f++ { // 前景色彩 = 30-37
			for d := range []int{0, 1, 4, 5, 7, 8} { // 显示方式 = 0,1,4,5,7,8
				txtColor := fmt.Sprintf(" (d=%d,b=%d,f=%d) ",d, b, f)
				fmt.Printf(" "+SetColor(txtColor,d,b,f)+" ")
			}
			fmt.Println("")
		}
		fmt.Println("")
	}
}
