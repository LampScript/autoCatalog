package main

import (
	"fmt"
	"io"
	"io/ioutil"
	"os"
	"path/filepath"
	"strings"
)

var ProjectPath = GetCurrentDirectory()

func GetCurrentDirectory() string {
	dir, err := filepath.Abs(filepath.Dir(os.Args[0])) //返回绝对路径  filepath.Dir(os.Args[0])去除最后一个元素的路径
	if err != nil {
		panic(err.Error())
	}
	return strings.Replace(dir, "\\", "/", -1) //将\替换成/
}

// CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build main.go

func main() {
	currentPath := GetCurrentDirectory()
	mdFile := "catalog.md"
	if checkFileIsExist(mdFile) {
		err := os.Remove(mdFile)
		check(err)
	}
	listFile(mdFile, currentPath, 0)
}

func listFile(mdFile, currentPath string, deep int) {
	files, err := ioutil.ReadDir(currentPath)
	if err != nil {
		panic(err.Error())
	}
	for _, file := range files {
		if file.IsDir() {
			if !shouldSkip(file.Name()) {
				dirName := currentPath + "/" + file.Name()
				showDirName := strings.Replace(dirName, ProjectPath, ".", -1)
				fmt.Println(getPrefix(deep) + fmt.Sprintf("[%s](%s)", file.Name(), showDirName))
				writ(mdFile, getPrefix(deep)+fmt.Sprintf("[%s](%s)", file.Name(), showDirName))
				deep++
				listFile(mdFile, dirName, deep)
				deep--
			}
		} else {
			if !shouldSkip(file.Name()) {
				fileName := currentPath + "/" + file.Name()
				if strings.Contains(fileName, " ") {
					newName := strings.Replace(fileName, " ", "", -1)
					err := os.Rename(fileName, newName)
					if err != nil {
						fmt.Println(err.Error())
					}
				}
				md := getPrefix(deep) + fmt.Sprintf("[%s](%s)", file.Name(), fileName)
				showFileName := strings.Replace(md, ProjectPath, ".", -1)
				fmt.Println(showFileName)
				writ(mdFile, showFileName)
			}
		}
	}
}

func shouldSkip(str string) bool {
	if strings.Contains(str, `.git`) {
		return true
	} else if strings.Contains(str, `.idea`) {
		return true
	} else if strings.Contains(str, `DS_Store`) {
		return true
	} else if strings.Contains(str, `autoCatalog`) {
		return true
	}
	return false
}

func getPrefix(deep int) string {
	if deep < 0 {
		deep = 0
	}
	temp := strings.Repeat(`    `, deep)
	return temp + "- "
}

func writ(fileName, input string) {
	var f *os.File
	var err error
	if checkFileIsExist(fileName) { //如果文件存在
		f, err = os.OpenFile(fileName, os.O_APPEND|os.O_WRONLY, 0666) //打开文件
	} else {
		f, err = os.Create(fileName) //创建文件
	}
	check(err)
	_, err = io.WriteString(f, input+"\n") //写入文件(字符串)
	check(err)
}

/**
* 判断文件是否存在  存在返回 true 不存在返回false
 */

func checkFileIsExist(filename string) bool {
	var exist = true
	if _, err := os.Stat(filename); os.IsNotExist(err) {
		exist = false
	}
	return exist
}

func check(e error) {
	if e != nil {
		panic(e)
	}
}
