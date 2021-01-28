package Modules

import (
	"io/ioutil"
	"os"
)

type file struct {}
var File file

func (*file)Exists(path string) bool {
	_, err := os.Stat(path)    //os.Stat获取文件信息
	if err != nil {
		if os.IsExist(err) {
			return true
		}
		return false
	}
	return true
}

func (*file)Read(Path string)([]byte,error){
	return ioutil.ReadFile(Path)
}

func (file)IsDir(path string) bool {
	s, err := os.Stat(path)
	if err != nil {
		return false
	}
	return s.IsDir()
}

func (*file)Remove(Path string)error{
	if File.IsDir(Path) {
		return os.RemoveAll(Path)
	}
	return os.Remove(Path)
}

func (*file)Write(Path string,content []byte)error{
	return ioutil.WriteFile(Path,content,600)
}