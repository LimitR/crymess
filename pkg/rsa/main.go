package manager

import (
	"crypto/rand"
	"crypto/rsa"
	"crypto/sha512"
	"crypto/x509"
	"encoding/asn1"
	"encoding/base64"
	"encoding/pem"
	"fmt"
	"io/ioutil"
	"os"
)

type ManagerRSA struct {
	privateKey *rsa.PrivateKey
	publicKey  rsa.PublicKey
}

func OptionString(s string) *string {
	return &s
}

func NewManagerRSA(pathToPrivateKey *string) *ManagerRSA {
	if pathToPrivateKey == nil {
		privateKey, _ := rsa.GenerateKey(rand.Reader, 4096)
		m := &ManagerRSA{
			privateKey: privateKey,
			publicKey:  privateKey.PublicKey,
		}
		m.Save()
		return m
	} else {
		m := ManagerRSA{}
		err := m.load(*pathToPrivateKey)
		if err != nil {
			panic(err)
		}
		return &m
	}
}

func (m *ManagerRSA) load(pathToPrivateKey string) error {
	b, err := ioutil.ReadFile(pathToPrivateKey)

	if err != nil {
		return err
	}

	block, _ := pem.Decode(b)
	key, err := x509.ParsePKCS1PrivateKey(block.Bytes)

	if err != nil {
		return err
	}

	m.privateKey = key

	m.publicKey = key.PublicKey

	return nil
}

func (m *ManagerRSA) Encrypt(secretMessage string) string {
	label := []byte("")
	rng := rand.Reader
	ciphertext, _ := rsa.EncryptOAEP(sha512.New(), rng, &m.publicKey, []byte(secretMessage), label)
	return base64.StdEncoding.EncodeToString(ciphertext)
}

func (m *ManagerRSA) Decrypt(cipherText string) string {
	ct, _ := base64.StdEncoding.DecodeString(cipherText)
	label := []byte("")
	rng := rand.Reader
	plaintext, _ := rsa.DecryptOAEP(sha512.New(), rng, m.privateKey, ct, label)
	return string(plaintext)
}

func (m *ManagerRSA) GetPublicKey() []byte {
	f, err := ioutil.ReadFile("./etc/public.pub")
	if err != nil {
		panic(err)
	}
	return f
}

func (m *ManagerRSA) Save() {
	savePEMKey("private", m.privateKey)
	savePublicPEMKey("public.pub", m.publicKey)
}

func savePEMKey(fileName string, key *rsa.PrivateKey) {
	outFile, err := os.Create("./etc/" + fileName)
	checkError(err)
	defer outFile.Close()

	var privateKey = &pem.Block{
		Type:  "PRIVATE KEY",
		Bytes: x509.MarshalPKCS1PrivateKey(key),
	}

	err = pem.Encode(outFile, privateKey)
	checkError(err)
}

func savePublicPEMKey(fileName string, pubkey rsa.PublicKey) {
	asn1Bytes, err := asn1.Marshal(pubkey)
	checkError(err)

	var pemkey = &pem.Block{
		Type:  "PUBLIC KEY",
		Bytes: asn1Bytes,
	}

	pemfile, err := os.Create("./etc/" + fileName)
	if err != nil {
		if err = os.Mkdir("./etc/", os.ModePerm); err != nil {
			panic(err)
		}
		pemfile, _ = os.Create("./etc/" + fileName)
	}
	defer pemfile.Close()

	err = pem.Encode(pemfile, pemkey)
	checkError(err)
}

func checkError(err error) {
	if err != nil {
		fmt.Println("Fatal error ", err.Error())
		os.Exit(1)
	}
}
