package kfile

import (
	"log"
	"os"
	"path/filepath"
	"strings"

	"go.uber.org/multierr"
)

// GetCurrentDirectory
//  @Description  获得当前绝对路径
//  @Return string
func GetCurrentDirectory() string {
	if dir, err := os.Getwd(); err == nil {
		return dir
	}
	dir, err := filepath.Abs(filepath.Dir(os.Args[0]))
	if err != nil {
		log.Fatal("err", err)
	}

	return strings.Replace(dir, "\\", "/", -1)
}

// LeftAddPathPos
//  @Description  检测并补全路径左边的反斜杠
//  @Param path
//  @Return string
func LeftAddPathPos(path string) string {
	if path[:0] != "/" {
		path = "/" + path
	}
	return path
}

// RightAddPathPos
//  @Description 检测并补全路径右边的反斜杠
//  @Param path
//  @Return string
func RightAddPathPos(path string) string {
	if path[len(path)-1:] != "/" {
		path = path + "/"
	}
	return path
}

// CreateDir
//  @Description  不存在则创建目录
//  @Param folderPath
func CreateDir(folderPath string) {
	if _, err := os.Stat(folderPath); os.IsNotExist(err) {
		// 必须分成两步：先创建文件夹、再修改权限
		os.MkdirAll(folderPath, 0777) //0777也可以os.ModePerm
		os.Chmod(folderPath, 0777)
	}
}

// MakeDirectory
//  @Description 创建目录 支持多个目录
//  @Param dirs
//  @Return error
func MakeDirectory(dirs ...string) error {
	var errs error
	for _, dir := range dirs {
		if !Exists(dir) {
			if err := os.MkdirAll(dir, 0755); err != nil {
				errs = multierr.Append(errs, err)
			}
		}
	}
	return errs
}

// IsExists
//  @Description  检测文件或者目录是否存在
//  @Param name
//  @Return bool
func IsExists(name string) bool {
	_, err := os.Stat(name)
	return err == nil || os.IsExist(err)
}
