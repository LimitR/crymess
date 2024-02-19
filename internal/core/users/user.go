package users

import (
	"bufio"
	"crypto/rand"
	"crypto/rsa"
	"crypto/sha512"
	"crypto/x509"
	"encoding/asn1"
	"encoding/base64"
	"encoding/pem"
	"errors"
	"io/ioutil"
	"os"
	"strconv"
	"strings"
	"time"
)

type Message struct {
	MyMessage bool
	Id        string
	Text      string
	TimeStamp time.Time
}

type User struct {
	Name          string
	PasswordHash  string
	MessageList   []*Message
	MyMessageList []*Message
	PubKey        *rsa.PublicKey
}

func NewUser(name string, pubKey string) User {
	err := ioutil.WriteFile("./"+DIR+"/"+USER_PUB_DIR+"/"+name+".pub", []byte(pubKey), 0777)
	if err != nil {
		panic(err)
	}
	u := User{
		Name:          name,
		MessageList:   make([]*Message, 0, 10),
		MyMessageList: make([]*Message, 0, 10),
	}
	block, _ := pem.Decode([]byte(pubKey))
	pub, err := x509.ParsePKCS1PublicKey(block.Bytes)
	if err != nil {
		panic(err)
	}

	u.PubKey = pub
	return u
}

func (u *User) load(file []byte) {
	block, _ := pem.Decode(file)
	pubKey, err := x509.ParsePKCS1PublicKey(block.Bytes)
	if err != nil {
		panic(err)
	}

	u.PubKey = pubKey
}

func (u *User) loadMessage(filePath string) {
	file, err := os.Open(filePath)
	if err != nil {
		ioutil.WriteFile(filePath, []byte{}, 0777)
		return
	}
	defer file.Close()

	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		txt := scanner.Text()
		sliceText := strings.Split(txt, " ")
		unixNum, _ := strconv.ParseInt(sliceText[0], 10, 64)
		if sliceText[1] == "my" {
			u.MessageList = append(u.MessageList, &Message{
				Text:      getNormalString(sliceText),
				TimeStamp: time.Unix(unixNum, 0),
				MyMessage: true,
			})
		} else {
			u.MessageList = append(u.MessageList, &Message{
				Text:      getNormalString(sliceText),
				TimeStamp: time.Unix(unixNum, 0),
				MyMessage: false,
			})
		}
	}
	if err := scanner.Err(); err != nil {
		panic(err)
	}
}

func getNormalString(txt []string) string {
	res := ""

	for i, s := range txt {
		if i > 1 {
			res += s
		}
	}

	return res
}

func (u *User) Save() error {
	asn1Bytes, err := asn1.Marshal(*u.PubKey)
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

	if _, err := os.Stat("./" + DIR + "/" + MESSAGE_DIR); errors.Is(err, os.ErrNotExist) {
		os.Mkdir("./"+DIR+"/"+MESSAGE_DIR, 0777)
	}

	if _, err := os.Stat("./" + DIR + "/" + MESSAGE_DIR + "/" + u.Name + ".message"); errors.Is(err, os.ErrNotExist) {
		os.Create("./" + DIR + "/" + MESSAGE_DIR + "/" + u.Name + ".message")
	}

	file := []byte{}

	for _, msg := range u.MessageList {
		file = append(file, []byte(strconv.Itoa(int(msg.TimeStamp.Unix())))...)
		file = append(file, []byte(" ")...)
		if msg.MyMessage {
			file = append(file, []byte("my")...)
			file = append(file, []byte(" ")...)
		} else {
			file = append(file, []byte("you")...)
			file = append(file, []byte(" ")...)
		}
		file = append(file, []byte(msg.Text)...)
		file = append(file, []byte("\n")...)
	}

	err = ioutil.WriteFile("./"+DIR+"/"+MESSAGE_DIR+"/"+u.Name+".message", file, 0777)

	return err
}

func (m *User) GetMessage() []string {
	msgList := make([]string, 0, len(m.MessageList))

	for _, msg := range m.MessageList {
		pubMsg := m.Encrypt(msg.Text)
		msgList = append(msgList, pubMsg)
	}

	return msgList
}

func (m *User) AddMessage(text string, myMessage bool) {
	m.MessageList = append(m.MessageList, &Message{
		Text:      text,
		TimeStamp: time.Now(),
		MyMessage: myMessage,
	})
}

func (m *User) Encrypt(secretMessage string) string {
	if m.PubKey == nil {
		panic("PubKey not found")
	}
	label := []byte("")
	rng := rand.Reader
	ciphertext, _ := rsa.EncryptOAEP(sha512.New(), rng, m.PubKey, []byte(secretMessage), label)

	return base64.StdEncoding.EncodeToString(ciphertext)
}
