package gcrypt

import (
	"crypto/rand"
	"encoding/base64"
	"encoding/hex"
	"github.com/tjfoc/gmsm/sm2"
	"github.com/tjfoc/gmsm/x509"
	"io"
)

// NewSm2 国密Sm2
func NewSm2() *Sm2Service {
	return &Sm2Service{
		mode:   C1C3C2,
		random: rand.Reader,
	}
}

type Sm2Service struct {
	publicKey  *sm2.PublicKey
	privateKey *sm2.PrivateKey
	random     io.Reader
	mode       Mode
}

type Mode int

var (
	C1C3C2 = Mode(sm2.C1C3C2)
	C1C2C3 = Mode(sm2.C1C2C3)
)

func (s *Sm2Service) SetMode(v Mode) {
	s.mode = v
}
func (s *Sm2Service) SetRandom(v io.Reader) {
	s.random = v
}

func (s *Sm2Service) SetPublicKey(v *sm2.PublicKey) {
	s.publicKey = v
}

func (s *Sm2Service) SetPublicKeyByByte(v []byte) (err error) {
	p, err := x509.ParseSm2PublicKey(v)
	if err != nil {
		return
	}
	s.SetPublicKey(p)
	return
}

func (s *Sm2Service) SetPublicKeyByHex(v string) (err error) {
	b, err := hex.DecodeString(v)
	if err != nil {
		return
	}
	err = s.SetPublicKeyByByte(b)
	if err != nil {
		return
	}
	return
}

func (s *Sm2Service) SetPublicKeyByBase64(v string) (err error) {
	b, err := base64.StdEncoding.DecodeString(v)
	if err != nil {
		return
	}
	err = s.SetPublicKeyByByte(b)
	if err != nil {
		return
	}
	return
}

// Encrypt 加密
func (s *Sm2Service) Encrypt(data []byte) ([]byte, error) {
	return sm2.Encrypt(s.publicKey, data, s.random, int(s.mode))
}

// Decrypt 解密
func (s *Sm2Service) Decrypt(data []byte) ([]byte, error) {
	return sm2.Decrypt(s.privateKey, data, int(s.mode))
}
