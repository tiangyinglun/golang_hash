package main

import (
	"net/http"
	"fmt"
	"strconv"
	"os"
	"io"
)

var path string = "E:/gows/src/test/stream/node-v8.9.1-x64.msi"
var url string = "https://nodejs.org/dist/v8.9.1/node-v8.9.1-x64.msi"

func main() {
	//打开文件 地址没有就创建
	f, err := os.OpenFile(path, os.O_RDWR|os.O_CREATE|os.O_APPEND, 0666) //其实这里的 O_RDWR应该是 O_RDWR|O_CREATE，也就是文件不存在的情况下就建一个空文件，
	if err != nil {
		panic(err)
	}
	//获取文件信息 这就是状态
	stat, err := f.Stat()
	if err != nil {
		panic(err)
	}
	//把指针指向 文件大小位置
	f.Seek(stat.Size(), 0)
	//构建client
	client := &http.Client{}
	//构建请求
	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		fmt.Print(err)
	}
	//设置头信息
	req.Header.Set("User-Agent", "Mozilla/5.0 (Windows NT 10.0; WOW64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/63.0.3236.0 Safari/537.36")
	req.Header.Set("Range", "bytes="+strconv.FormatInt(stat.Size(), 10)+"-")
	resp, err := client.Do(req)
	defer resp.Body.Close()
	written, err := io.Copy(f, resp.Body)
	if err != nil {
		panic(err)
	}
	println("written: ", written)
}
