package users

import (
	"bufio"
	"crypto/rand"
	"crypto/rsa"
	"crypto/sha256"
	"crypto/x509"
	"encoding/asn1"
	"encoding/base64"
	"encoding/pem"
	"os"
	"time"
)

type Message struct {
	Load      bool
	Id        string
	Text      string
	TimeStamp time.Time
}

type User struct {
	Name         string
	PasswordHash string
	MessageList  []*Message
	PubKey       rsa.PublicKey
}

func NewUser(name string) User {
	return User{
		Name: name,
	}
}

func (u *User) load(file []byte) {
	block, _ := pem.Decode(file)
	pubKey, err := x509.ParsePKCS1PublicKey(block.Bytes)
	if err != nil {
		panic(err)
	}

	u.PubKey = *pubKey
}

func (u *User) loadMessage(filePath string) {
	file, err := os.Open(filePath)
	if err != nil {
		panic(err)
	}
	defer file.Close()

	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		u.MessageList = append(u.MessageList, &Message{
			Text: scanner.Text(),
		})
	}

	if err := scanner.Err(); err != nil {
		panic(err)
	}
}

func (u *User) Save() error {
	asn1Bytes, err := asn1.Marshal(u.PubKey)
	if err != nil {
		return err
	}

	var pemkey = &pem.Block{
		Type:  "PUBLIC KEY",
		Bytes: asn1Bytes,
	}

	pemfile, err := os.Create("./" + DIR + "/" + USER_PUB_DIR + "/" + u.Name + ".pub")
	if err != nil {
		return err
	}
	defer pemfile.Close()

	err = pem.Encode(pemfile, pemkey)
	if err != nil {
		return err
	}

	f, _ := os.Create("./" + DIR + "/" + MESSAGE_DIR + "/" + u.Name + ".message")
	defer f.Close()

	w := bufio.NewWriter(f)

	for _, msg := range u.MessageList {
		w.WriteString(msg.Text + "\n")
	}

	return nil
}

func (m *User) GetMessage() []string {
	msgList := make([]string, 0, len(m.MessageList))

	for _, msg := range m.MessageList {
		if !msg.Load {
			pubMsg := m.encrypt(msg.Text)
			msg.Load = true
			msgList = append(msgList, pubMsg)
		}
	}

	return msgList
}

func (m *User) AddMessage(text string) {
	m.MessageList = append(m.MessageList, &Message{
		Text:      text,
		TimeStamp: time.Now(),
	})
}

func (m *User) encrypt(secretMessage string) string {
	label := []byte("")
	rng := rand.Reader
	ciphertext, _ := rsa.EncryptOAEP(sha256.New(), rng, &m.PubKey, []byte(secretMessage), label)
	return base64.StdEncoding.EncodeToString(ciphertext)
}
