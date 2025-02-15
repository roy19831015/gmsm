package main

import (
	"bufio"
	"bytes"
	"crypto/des"
	"errors"
	"fmt"
	"github.com/roy19831015/gmsm/ucapp4go"
	"os"
	"strings"
)

func main() {
	var str string
	for {
		fmt.Println("请输入待加密的数据库密码：")
		str = GetUserInputFromOsStdIn()
		if len(str) <= 0 {
			fmt.Println("输入错误请重试")
			continue
		}
		key := "2Qh/lOwL5Zjx1pKnvKGPXVvvc9lMaK7p"
		decode, err := ucapp4go.Base64Decode(key)
		if err != nil {
			fmt.Println("输入错误请重试")
			continue
		}
		desEncrypt, err := TripleEcbDesEncrypt([]byte(str), decode)
		if err != nil {
			fmt.Println("输入错误请重试")
			continue
		}
		encode, err := ucapp4go.Base64Encode(desEncrypt)
		if err != nil {
			fmt.Println("输入错误请重试")
			continue
		}
		fmt.Println("加密后的密码为:", encode)
	}
}

func GetUserInputFromOsStdIn() string {
	var bs []byte
	var str string
	buf := bufio.NewReader(os.Stdin)
	fmt.Print("> ")
	bs, _ = buf.ReadBytes('\n')
	str = string(bs)
	str = strings.ReplaceAll(str, "\r", "")
	str = strings.ReplaceAll(str, "\n", "")
	return str
}

//ECB PKCS5Padding
func PKCS5Padding(ciphertext []byte, blockSize int) []byte {
	padding := blockSize - len(ciphertext)%blockSize
	padtext := bytes.Repeat([]byte{byte(padding)}, padding)
	return append(ciphertext, padtext...)
}

//ECB PKCS5Unpadding
func PKCS5Unpadding(origData []byte) []byte {
	length := len(origData)
	unpadding := int(origData[length-1])
	return origData[:(length - unpadding)]
}

//Des加密
func encrypt(origData, key []byte) ([]byte, error) {
	if len(origData) < 1 || len(key) < 1 {
		return nil, errors.New("wrong data or key")
	}
	block, err := des.NewCipher(key)
	if err != nil {
		return nil, err
	}
	bs := block.BlockSize()
	if len(origData)%bs != 0 {
		return nil, errors.New("wrong padding")
	}
	out := make([]byte, len(origData))
	dst := out
	for len(origData) > 0 {
		block.Encrypt(dst, origData[:bs])
		origData = origData[bs:]
		dst = dst[bs:]
	}
	return out, nil
}

//Des解密
func decrypt(crypted, key []byte) ([]byte, error) {
	if len(crypted) < 1 || len(key) < 1 {
		return nil, errors.New("wrong data or key")
	}
	block, err := des.NewCipher(key)
	if err != nil {
		return nil, err
	}
	out := make([]byte, len(crypted))
	dst := out
	bs := block.BlockSize()
	if len(crypted)%bs != 0 {
		return nil, errors.New("wrong crypted size")
	}

	for len(crypted) > 0 {
		block.Decrypt(dst, crypted[:bs])
		crypted = crypted[bs:]
		dst = dst[bs:]
	}

	return out, nil
}

//[golang ECB 3DES Encrypt]
func TripleEcbDesEncrypt(origData, key []byte) ([]byte, error) {
	tkey := make([]byte, 24, 24)
	copy(tkey, key)
	k1 := tkey[:8]
	k2 := tkey[8:16]
	k3 := tkey[16:]

	block, err := des.NewCipher(k1)
	if err != nil {
		return nil, err
	}
	bs := block.BlockSize()
	origData = PKCS5Padding(origData, bs)

	buf1, err := encrypt(origData, k1)
	if err != nil {
		return nil, err
	}
	buf2, err := decrypt(buf1, k2)
	if err != nil {
		return nil, err
	}
	out, err := encrypt(buf2, k3)
	if err != nil {
		return nil, err
	}
	return out, nil
}

//[golang ECB 3DES Decrypt]
func TripleEcbDesDecrypt(crypted, key []byte) ([]byte, error) {
	tkey := make([]byte, 24, 24)
	copy(tkey, key)
	k1 := tkey[:8]
	k2 := tkey[8:16]
	k3 := tkey[16:]
	buf1, err := decrypt(crypted, k3)
	if err != nil {
		return nil, err
	}
	buf2, err := encrypt(buf1, k2)
	if err != nil {
		return nil, err
	}
	out, err := decrypt(buf2, k1)
	if err != nil {
		return nil, err
	}
	out = PKCS5Unpadding(out)
	return out, nil
}
