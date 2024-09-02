package gjwt

import (
	"bytes"
	"crypto/aes"
	"crypto/cipher"
	"crypto/md5"
	"encoding/base64"
	"errors"
	"fmt"
)

func sign(v string) string {
	return fmt.Sprintf("%x", md5.Sum([]byte(v)))
}

// AES CBC 加密
func aesEncodeCBC(plaintext, key string) (ciphertext string, err error) {
	block, err := aes.NewCipher([]byte(key))
	if err != nil {
		err = errors.New("aes.NewCipher error " + err.Error())
		return
	}
	plainText := padding([]byte(plaintext), block.BlockSize())
	blockMode := cipher.NewCBCEncrypter(block, []byte(key[0:block.BlockSize()]))
	cipherText := make([]byte, len(plainText))
	blockMode.CryptBlocks(cipherText, plainText)
	ciphertext = base64.RawURLEncoding.EncodeToString(cipherText)
	return
}

// AES CBC 解密
func aesDecodeCBC(ciphertext, key string) (plaintext string, err error) {
	cipherText, err := base64.RawURLEncoding.DecodeString(ciphertext)
	if err != nil {
		err = errors.New("base64.Decode ciphertext  error" + err.Error())
		return
	}

	// AES 解密
	block, err := aes.NewCipher([]byte(key))
	if err != nil {
		err = errors.New("aes.NewCipher error" + err.Error())
		return
	}
	blockMode := cipher.NewCBCDecrypter(block, []byte(key[0:block.BlockSize()]))
	plainText := make([]byte, len(cipherText))
	blockMode.CryptBlocks(plainText, cipherText)
	plainText = unPadding(plainText)

	// 明文结果
	plaintext = string(plainText)
	return
}

func padding(plainText []byte, blockSize int) []byte {
	n := blockSize - len(plainText)%blockSize //计算要填充的长度
	temp := bytes.Repeat([]byte{byte(n)}, n)  //对原来的明文填充n个n
	plainText = append(plainText, temp...)
	return plainText
}

func unPadding(cipherText []byte) []byte {
	end := cipherText[len(cipherText)-1]               //取出密文最后一个字节end
	cipherText = cipherText[:len(cipherText)-int(end)] //删除填充
	return cipherText
}
