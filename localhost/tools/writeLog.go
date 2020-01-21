package tools

import (
	"bytes"
	"crypto/aes"
	"crypto/cipher"
	"crypto/md5"
	"crypto/rand"
	"encoding/json"
	"github.com/micro/go-micro/config"
	"io"
	"io/ioutil"
	"strings"
	"time"
)


//使用PKCS7进行填充，IOS也是7
func PKCS7Padding(ciphertext []byte, blockSize int) []byte {
	padding := blockSize - len(ciphertext) % blockSize
	padtext := bytes.Repeat([]byte{byte(padding)}, padding)
	return append(ciphertext, padtext...)
}

func PKCS7UnPadding(origData []byte) []byte {
	length := len(origData)
	unpadding := int(origData[length-1])
	return origData[:(length - unpadding)]
}

//aes加密，填充秘钥key的16位，24,32分别对应AES-128, AES-192, or AES-256.
func AesCBCEncrypt(rawData,key []byte) ([]byte, error) {
	block, err := aes.NewCipher(key)
	if err != nil {
		panic(err)
	}

	//填充原文
	blockSize := block.BlockSize()
	rawData = PKCS7Padding(rawData, blockSize)
	//初始向量IV必须是唯一，但不需要保密
	cipherText := make([]byte,blockSize+len(rawData))
	//block大小 16
	iv := cipherText[:blockSize]
	if _, err := io.ReadFull(rand.Reader,iv); err != nil {
		panic(err)
	}

	//block大小和初始向量大小一定要一致
	mode := cipher.NewCBCEncrypter(block,iv)
	mode.CryptBlocks(cipherText[blockSize:],rawData)

	return cipherText, nil
}

func AesCBCDncrypt(encryptData, key []byte) ([]byte,error) {
	block, err := aes.NewCipher(key)
	if err != nil {
		panic(err)
	}

	blockSize := block.BlockSize()

	if len(encryptData) < blockSize {
		panic("ciphertext too short")
	}
	iv := encryptData[:blockSize]
	encryptData = encryptData[blockSize:]

	// CBC mode always works in whole blocks.
	if len(encryptData)%blockSize != 0 {
		panic("ciphertext is not a multiple of the block size")
	}

	mode := cipher.NewCBCDecrypter(block, iv)

	// CryptBlocks can work in-place if the two arguments are the same.
	mode.CryptBlocks(encryptData, encryptData)
	//解填充
	encryptData = PKCS7UnPadding(encryptData)
	return encryptData,nil
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
	x1, err := AesCBCEncrypt(d, key)
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
