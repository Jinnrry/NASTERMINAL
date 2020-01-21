package tools

import (
	"fmt"
	"io/ioutil"
	"net/http"
	"regexp"
	"strings"
)

//调用第三方接口获取公网ip地址
func GetPublicIp() string {

	//发起get请求
	resp, err1 := http.Get("http://myhttpheader.com/")
	if err1 != nil || resp.StatusCode != http.StatusOK {
		fmt.Println("get fail:", err1)
		return ""
	}
	defer resp.Body.Close()
	//读取响应体
	body, err2 := ioutil.ReadAll(resp.Body)
	if err2 != nil {
		fmt.Println("read body fail")
		return ""
	}

	sbody := string(body)
	sbody = strings.Replace(sbody, "\n", "", -1)
	sbody = strings.Replace(sbody, "\r", "", -1)

	re, _ := regexp.Compile(`(?U)X-Real-Ip\:<\/div>.*<\/div>`)

	all := re.FindAll([]byte(sbody), 1)

	stmp := ""

	for _, item := range all {
		stmp = string(item)
	}

	if stmp != "" {
		stmp = strings.ReplaceAll(stmp, "\t", "")
		stmp = strings.ReplaceAll(stmp, " ", "")
		re, _ = regexp.Compile(`[\d\.]+`)

		return string(re.Find([]byte(stmp)))
	}

	return ""
}
