package bootstrap

import (
	"bufio"
	"fmt"
	"io"
	"io/ioutil"
	"net/http"
	"os"
)

type BootStrap struct {
	filepath string
	bootlist []string
	err      error
}

var B BootStrap

func InitBootStrap(file string) {
	B = BootStrap{filepath: file}
	getBootList()
}

func getBootList() {
	file, err := os.Open(B.filepath)
	if err != nil {
		B.err = err
		return
	}

	var result []string
	bufReader := bufio.NewReader(file)
	for {
		a, _, c := bufReader.ReadLine()
		if c == io.EOF {
			break
		}
		result = append(result, string(a))
	}

	B.bootlist = result
	removeLocal()
}

func removeLocal() {
	localIp := getLocalIP()

	list := B.bootlist
	for k, v := range list {
		if v == localIp {
			list = append(list[:k], list[k+1:]...)
		}
	}

	B.bootlist = list
}

func getLocalIP() string {
	resp, err := http.Get("http://myexternalip.com/raw")
	if err != nil {
		return ""
	}
	defer resp.Body.Close()

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return ""
	}

	return string(body)
}
