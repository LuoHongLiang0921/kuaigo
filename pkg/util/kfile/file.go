package kfile

import (
	"bufio"
	"fmt"
	"io"
	"io/ioutil"
	"os"
	"path/filepath"
	"runtime"
	"strings"
)

// FileInfo  describes a configuration file and is returned by fileStat.
type FileInfo struct {
	Uid  uint32
	Gid  uint32
	Mode os.FileMode
	Md5  string
}

// Exists
//  @Description  检测文件是否存在
//  @Param fpath
//  @Return bool
func Exists(fpath string) bool {
	if _, err := os.Stat(fpath); os.IsNotExist(err) {
		return false
	}

	return true
}

// ListFiles
//  @Description  列出目录内所有文件
//  @Param dir
//  @Param ext
//  @Return []string
func ListFiles(dir string, ext string) []string {
	fs, err := ioutil.ReadDir(dir)
	if err != nil {
		return []string{}
	}

	var ret []string
	for _, fp := range fs {
		if fp.IsDir() {
			continue
		}

		if ext != "" && filepath.Ext(fp.Name()) != ext {
			continue
		}

		ret = append(ret, dir+"/"+fp.Name())
	}

	return ret
}

// IsFileChanged reports whether src and dest config files are equal.
// Two config files are equal when they have the same file contents and
// Unix permissions. The owner, group, and mode must match.
// Returns false in other cases.
func IsFileChanged(src, dest string) (bool, error) {
	if !Exists(dest) {
		return true, nil
	}
	d, err := FileStat(dest)
	if err != nil {
		return true, err
	}
	s, err := FileStat(src)
	if err != nil {
		return true, err
	}

	if d.Uid != s.Uid || d.Gid != s.Gid || d.Mode != s.Mode || d.Md5 != s.Md5 {
		return true, nil
	}
	return false, nil
}

// IsDirectory
//  @Description  是否是目录
//  @Param path 检测路径
//  @Return bool
//  @Return error
func IsDirectory(path string) (bool, error) {
	f, err := os.Stat(path)
	if err != nil {
		return false, err
	}
	switch mode := f.Mode(); {
	case mode.IsDir():
		return true, nil
	case mode.IsRegular():
		return false, nil
	}
	return false, nil
}

// RecursiveFilesLookup ...
func RecursiveFilesLookup(root string, pattern string) ([]string, error) {
	return recursiveLookup(root, pattern, false)
}

// RecursiveDirsLookup ...
func RecursiveDirsLookup(root string, pattern string) ([]string, error) {
	return recursiveLookup(root, pattern, true)
}

func recursiveLookup(root string, pattern string, dirsLookup bool) ([]string, error) {
	var result []string
	isDir, err := IsDirectory(root)
	if err != nil {
		return nil, err
	}
	if isDir {
		err := filepath.Walk(root, func(root string, f os.FileInfo, err error) error {
			match, err := filepath.Match(pattern, f.Name())
			if err != nil {
				return err
			}
			if match {
				isDir, err := IsDirectory(root)
				if err != nil {
					return err
				}
				if isDir && dirsLookup {
					result = append(result, root)
				} else if !isDir && !dirsLookup {
					result = append(result, root)
				}
			}
			return nil
		})
		if err != nil {
			return nil, err
		}
	} else {
		if !dirsLookup {
			result = append(result, root)
		}
	}
	return result, nil
}

func substr(s string, pos, length int) string {
	runes := []rune(s)
	l := pos + length
	if l > len(runes) {
		l = len(runes)
	}
	return string(runes[pos:l])
}

func getParentDirectory(dirctory string) string {
	if runtime.GOOS == "windows" {
		dirctory = strings.Replace(dirctory, "\\", "/", -1)
	}
	return substr(dirctory, 0, strings.LastIndex(dirctory, "/"))
}

// CheckAndGetParentDir ...
func CheckAndGetParentDir(path string) string {
	// check path is the directory
	isDir, err := IsDirectory(path)
	if err != nil || isDir {
		return path
	}
	return getParentDirectory(path)
}

// MkdirIfNecessary ...
func MkdirIfNecessary(createDir string) error {
	var path string
	var err error
	//前边的判断是否是系统的分隔符
	if os.IsPathSeparator('\\') {
		path = "\\"
	} else {
		path = "/"
	}

	s := strings.Split(createDir, path)
	startIndex := 0
	dir := ""
	if s[0] == "" {
		startIndex = 1
	} else {
		dir, _ = os.Getwd() //当前的目录
	}
	for i := startIndex; i < len(s); i++ {
		d := dir + path + strings.Join(s[startIndex:i+1], path)
		if _, e := os.Stat(d); os.IsNotExist(e) {
			//在当前目录下生成md目录
			err = os.Mkdir(d, os.ModePerm)
			if err != nil {
				break
			}
		}
	}
	return err
}

// IsPathExist
//  @Description  路径是否存在
//  @Param path
//  @Return bool
//  @Return error
func IsPathExist(path string) (bool, error) {
	_, err := os.Stat(path)
	if err == nil {
		return true, nil
	}
	if os.IsNotExist(err) {
		return false, nil
	}

	return true, err
}

// SaveFile
//  @Description  将数据保存到文件
//  @Param path
//  @Param value
func SaveFile(path string, value string) {
	fileHandler, err := os.OpenFile(path, os.O_RDWR|os.O_CREATE|os.O_APPEND, 0666)
	if err != nil {
		return
	}
	defer fileHandler.Close()

	buf := bufio.NewWriter(fileHandler)

	fmt.Fprintln(buf, value)

	buf.Flush()
}

// ReadFile
//  @Description  一次性读取文件 string
//  @Param filePth
//  @Return string
func ReadFile(filePth string) string {
	f, err := os.Open(filePth)
	defer f.Close()
	if err != nil {
		return ""
	}
	r, err := ioutil.ReadAll(f)
	return string(r)
}

// ReadFileByte
//  @Description   一次性读取文件 []byte
//  @Param filePth
//  @Return []byte
func ReadFileByte(filePth string) []byte {
	f, err := os.Open(filePth)
	defer f.Close()
	if err != nil {
		return nil
	}
	r, err := ioutil.ReadAll(f)
	return r
}

// WalkDir
//  @Description  递归读取目录下所有文件名,存在 *fileList
//  @Param dirpath
//  @Param fileList
func WalkDir(dirpath string, fileList *[]string) {
	files, err := ioutil.ReadDir(dirpath)
	if err != nil {
		return
	}
	for _, file := range files {
		if file.IsDir() {
			WalkDir(dirpath+"/"+file.Name(), fileList)
			continue
		} else {
			*fileList = append(*fileList, dirpath+"/"+file.Name())
		}
	}
}

// ReadListFile
//  @Description  按行读文件
//  @Param filename
//  @Return []string
func ReadListFile(filename string) []string {
	result := make([]string, 0)
	f, err := os.Open(filename)
	if err != nil {
		return nil
	}
	defer f.Close()
	rd := bufio.NewReader(f)

	for {
		line, err := rd.ReadString('\n')
		line = strings.TrimSpace(line)

		if line != "" {
			result = append(result, line)
		}

		if err == io.EOF {
			break
		}
	}
	return result
}

// WriteListFile
//  @Description  按行存储文件
//  @Param filename
//  @Param value
func WriteListFile(filename string, value []string) {
	fileHandler, err := os.Create(filename)
	defer fileHandler.Close()
	if err != nil {
		return
	}
	buf := bufio.NewWriter(fileHandler)

	for _, v := range value {
		fmt.Fprintln(buf, v)
	}
	buf.Flush()
}

// ListDir
//  @Description  获取指定目录下的所有文件，不进入下一级目录搜索，可以匹配后缀过滤。
//  @Param dirPth
//  @Return files
//  @Return err
func ListDir(dirPth string) (files []string, err error) {
	files = make([]string, 0, 10)
	dir, err := ioutil.ReadDir(dirPth)
	if err != nil {
		return nil, err
	}
	PthSep := string(os.PathSeparator)

	for _, fi := range dir {
		if fi.IsDir() { // 忽略目录
			continue
		}

		files = append(files, dirPth+PthSep+fi.Name())
	}
	return files, nil
}

// SimpleCopyFile
//  @Description  简单复制文件，用于小文件复制
//  @Param src 源文件路径
//  @Param dst 目标文件路径
//  @Return error
func SimpleCopyFile(src string, dst string) error {
	bytes, err := ioutil.ReadFile(src)
	if err != nil {
		return err
	}

	err = ioutil.WriteFile(dst, bytes, 0666)
	if err != nil {
		return err
	}

	return nil
}

// CopyFile
//  @Description  带缓冲的文件复制，用于大文件
//  @Param src 源文件路径
//  @Param dst 目标文件路径
//  @Return error
func CopyFile(src string, dst string) error {
	// 只读模式打开文件
	srcFile, err := os.OpenFile(src, os.O_RDONLY, os.ModePerm)
	if err != nil {
		return err
	}
	defer srcFile.Close()

	// 只写模式打开创建文件
	dstFile, err := os.OpenFile(dst, os.O_CREATE|os.O_WRONLY, os.ModePerm)
	if err != nil {
		return err
	}
	defer dstFile.Close()

	// 创建缓冲区[]byte
	var n int = 0
	bufbytes := make([]byte, 1024)
	reader := bufio.NewReader(srcFile)
	writer := bufio.NewWriter(dstFile)

	for {
		n++
		_, err = reader.Read(bufbytes)
		if err != nil {
			if err == io.EOF {
				break
			}
			return err
		}

		_, err := writer.Write(bufbytes)
		if err != nil {
			return err
		}
	}

	return nil
}
