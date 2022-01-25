package kfile

import (
	"fmt"
	"testing"
)

func TestGetCurrentDirectory(t *testing.T) {
	fmt.Println("GetCurrentDirectory:", GetCurrentDirectory())
	fmt.Println("LeftAddPathPos:", LeftAddPathPos("abc/scs"))
	fmt.Println("RightAddPathPos:", RightAddPathPos("abc/scs"))
	// fmt.Println("GetCurrentPackage:",GetCurrentPackage())
	//CreateDir(GetCurrentDirectory()+"/demo")
}
