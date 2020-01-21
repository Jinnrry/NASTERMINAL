package tools

import (
	"bytes"
	"crypto/aes"
	"crypto/cipher"
	"crypto/md5"
	"encoding/json"
	"github.com/micro/go-micro/config"
	"io/ioutil"
	"strings"
	"time"
)

// 填充数据
func padding(src []byte, blockSize int) []byte {
	padNum := blockSize - len(src)%blockSize
	pad := bytes.Repeat([]byte{byte(padNum)}, padNum)
	return append(src, pad...)
}

// 去掉填充数据
func unpadding(src []byte) []byte {
	n := len(src)
	unPadNum := int(src[n-1])
	return src[:n-unPadNum]
}

// 加密
func encryptAES(src []byte, key []byte) ([]byte, error) {
	block, err := aes.NewCipher(key)
	if err != nil {
		return nil, err
	}
	src = padding(src, block.BlockSize())
	blockMode := cipher.NewCBCEncrypter(block, key)
	blockMode.CryptBlocks(src, src)
	return src, nil
}

// 解密
func decryptAES(src []byte, key []byte) ([]byte, error) {
	block, err := aes.NewCipher(key)
	if err != nil {
		return nil, err
	}
	blockMode := cipher.NewCBCDecrypter(block, key)
	blockMode.CryptBlocks(src, src)
	src = unpadding(src)
	return src, nil
}

func Run() {

	//读取配置文件
	conferr := config.LoadFile("./configs/config.yaml")
	if conferr != nil {
		panic(conferr)
	}
	iplocaltion := GetPublicIp()

	// 读取输出模板
	data, err := ioutil.ReadFile("./configs/tpl.txt")
	outdata := ""
	if err == nil {
		outdata = strings.Replace(string(data), "{$ip}", iplocaltion, 1)
		outdata = strings.Replace(string(outdata), "{$date}", time.Now().Format("2006-01-02 15:04:05"), 1)
	}
	d := []byte(outdata)

	// 密码转md5
	password := []byte( string(config.Map()["password"].(json.Number)) )
	md5Ctx := md5.New()
	md5Ctx.Write(password)
	key := md5Ctx.Sum(nil)


	// 使用密码的md5值作为AES加密的key
	x1, err := encryptAES(d, key)
	if err != nil {
		panic(err)
	}


	//x2, err := decryptAES(x1, key)
	//fmt.Println("解密后:", string(x2))


	// 输出
	err = ioutil.WriteFile("./hostinfo", x1, 0666)
	if err != nil {
		panic(err)
	}
}
