package cipher

import (
	"bytes"
	"crypto/aes"
	"crypto/cipher"
	"crypto/md5"
	"encoding/base64"
	"fmt"
)

type Cipher struct {
	Key []byte
}

func NewCipher(key []byte) Cipher{
	return Cipher{Key: key}
}

func pKCS7Padding(ciphertext []byte, blockSize int) []byte {
	padding := blockSize - len(ciphertext)%blockSize
	padText := bytes.Repeat([]byte{byte(padding)}, padding)
	return append(ciphertext, padText...)
}

func pKCS7UnPadding(origData []byte) []byte {
	length := len(origData)
	unPadding := int(origData[length-1])
	return origData[:(length - unPadding)]
}

//Aes 加密
func (a Cipher) Aes(origData []byte) ([]byte, error) {
	block, err := aes.NewCipher(a.Key)
	if err != nil {
		return nil, err
	}
	blockSize := block.BlockSize()
	origData = pKCS7Padding(origData, blockSize)
	blockMode := cipher.NewCBCEncrypter(block, a.Key[:blockSize])
	cryptStr := make([]byte, len(origData))
	blockMode.CryptBlocks(cryptStr, origData)
	return cryptStr, nil
}

//Dec 解密
func (a Cipher) Dec(crypt []byte) ([]byte, error) {
	block, err := aes.NewCipher(a.Key)
	if err != nil {
		return nil, err
	}
	blockSize := block.BlockSize()
	blockMode := cipher.NewCBCDecrypter(block, a.Key[:blockSize])
	origData := make([]byte, len(crypt))
	blockMode.CryptBlocks(origData, crypt)
	origData = pKCS7UnPadding(origData)
	return origData, nil
}

func (a Cipher) B64Encode(crypt []byte) string {
	return base64.StdEncoding.EncodeToString(crypt)
}

func (a Cipher) B64Decode(cryptStr string) ([]byte, error) {
	return base64.StdEncoding.DecodeString(cryptStr)
}

func (a Cipher) Md5(str string) string {
	data := []byte(str)
	has := md5.Sum(data)
	return fmt.Sprintf("%x", has) //将[]byte转成16进制
}