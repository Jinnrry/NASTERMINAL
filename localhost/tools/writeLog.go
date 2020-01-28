package tools

import (
	"bytes"
	"crypto/aes"
	"crypto/cipher"
	"crypto/md5"
	"encoding/base64"
	"encoding/hex"
	"github.com/micro/go-micro/config"
	"io/ioutil"
	"log"
	"strings"
	"time"
)

var iv = []byte("IV:NasTerminal*!")

func Encrypt(text []byte, key []byte) (string, error) {
	//生成cipher.Block 数据块
	block, err := aes.NewCipher(key)
	if err != nil {
		log.Println("错误 -" + err.Error())
		return "", err
	}
	//填充内容，如果不足16位字符
	blockSize := block.BlockSize()
	originData := pad(text, blockSize)
	//加密方式
	blockMode := cipher.NewCBCEncrypter(block, iv)
	//加密，输出到[]byte数组
	crypted := make([]byte, len(originData))
	blockMode.CryptBlocks(crypted, originData)
	return base64.StdEncoding.EncodeToString(crypted), nil
}

func pad(ciphertext []byte, blockSize int) []byte {
	padding := blockSize - len(ciphertext)%blockSize
	padtext := bytes.Repeat([]byte{byte(padding)}, padding)
	return append(ciphertext, padtext...)
}

func Decrypt(text string, key []byte) (string, error) {
	decode_data, err := base64.StdEncoding.DecodeString(text)
	if err != nil {
		return "", nil
	}
	//生成密码数据块cipher.Block
	block, _ := aes.NewCipher(key)
	//解密模式
	blockMode := cipher.NewCBCDecrypter(block, iv)
	//输出到[]byte数组
	origin_data := make([]byte, len(decode_data))
	blockMode.CryptBlocks(origin_data, decode_data)
	//去除填充,并返回
	return string(unpad(origin_data)), nil
}

func unpad(ciphertext []byte) []byte {
	length := len(ciphertext)
	//去掉最后一次的padding
	unpadding := int(ciphertext[length-1])
	return ciphertext[:(length - unpadding)]
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
	password := []byte(string(config.Map()["password"].(string)))
	h := md5.New()
	h.Write(password)
	skey := hex.EncodeToString(h.Sum(nil))

	key := []byte(skey)

	// 使用密码的md5值作为AES加密的key
	x1, err := Encrypt(d, key)
	if err != nil {
		panic(err)
	}

	//x2, err := decryptAES(x1, key)
	//fmt.Println("解密后:", string(x2))

	// 输出
	err = ioutil.WriteFile("./hostinfo", []byte(x1), 0666)
	if err != nil {
		panic(err)
	}
}
