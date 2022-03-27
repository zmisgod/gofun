package cipher

import (
	"fmt"
	"log"
	"testing"
)

func TestEncode(t *testing.T) {
	key := []byte("*(!&*(SHHSsdfasd")
	cipher := NewCipher(key)
	res, err := cipher.Aes([]byte("121213"))
	if err != nil {
		log.Fatalln(err)
	} else {
		encodeStr := cipher.B64Encode(res)
		fmt.Println(encodeStr)
	}
}

func TestDecode(t *testing.T) {
	key := []byte("*(!&*(SHHSsdfasd")
	code := "lmo8J1EPIuTVTZe8rh6tUg=="
	cipher := NewCipher(key)
	resByte, err := cipher.B64Decode(code)
	if err != nil {
		log.Fatalln(err)
	}
	res, err := cipher.Dec(resByte)
	if err != nil {
		log.Fatalln(err)
	} else {
		fmt.Println(string(res))
	}
}

func TestMd5(t *testing.T) {
	fmt.Println(NewCipher([]byte("")).Md5("2222") == "934b535800b1cba8f96a5d72f72f1611")
}
