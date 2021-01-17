package service

import (
	"crypto/rand"
	"crypto/rsa"
	"crypto/x509"
	"encoding/base64"
	"fmt"
	mrand "math/rand"
)

type rsaKeyPar struct {
	PublicKey  string
	PrivateKey string
}

type RsaKeyPool struct {
	keyPar map[string]rsaKeyPar
	index  []string
	size   int
}

var pool *RsaKeyPool

func GetRsaKeyPool() *RsaKeyPool {
	return pool
}

func (pool *RsaKeyPool)GetPublicKey() string {
	index := mrand.Intn(pool.size - 1)
	return pool.keyPar[pool.index[index]].PublicKey
}

func (pool *RsaKeyPool) GetPrivateKey(publicKey string) string {
	return pool.keyPar[publicKey].PrivateKey
}

func RefreshPool(size int) *RsaKeyPool {
	if size < 1 {
		size = 50
	}
	pool = &RsaKeyPool{
		make(map[string]rsaKeyPar),
		make([]string, 0, 0),
		0,
	}
	n := mrand.Intn(size) + 3
	for i := 0; i < n; i++ {
		pub, pri, err := RSAKeyGen(1024)
		if err != nil {
			continue
		}
		pool.size++
		pool.keyPar[*pub] = rsaKeyPar{*pub, *pri}
		pool.index = append(pool.index, *pub)
	}
	return pool
}

func init() {
	RefreshPool(50)
}


func RSAKeyGen(bits int) (PublicKey *string, PrivateKey *string, err error) {
	privatekey, err := rsa.GenerateKey(rand.Reader, bits)
	if err != nil {
		return nil, nil, err
		fmt.Println("私钥文件生成失败")
	}
	//fmt.Println("私钥为：", privatekey)
	derStream := x509.MarshalPKCS1PrivateKey(privatekey)
	PriKey := base64.StdEncoding.EncodeToString(derStream)
	publickey := &privatekey.PublicKey
	//fmt.Println("公钥为：", publickey)
	derpkix, err := x509.MarshalPKIXPublicKey(publickey)
	if err != nil {
		fmt.Println("公钥文件生成失败")
		return nil, nil, err
	}
	PubKey := base64.StdEncoding.EncodeToString(derpkix)
	return &PubKey, &PriKey, nil
}
