package users

import (
	"io/ioutil"
	"os"
	"strings"
)

const (
	DIR          = "etc"
	USER_PUB_DIR = "users"
	MESSAGE_DIR  = "message"
)

type UserManager struct {
	UserList map[string]*User
}

func NewUserManager() *UserManager {
	return &UserManager{
		UserList: make(map[string]*User),
	}
}

func (u *UserManager) GetNameList() []string {
	res := make([]string, 0, 10)

	for _, user := range u.UserList {
		res = append(res, user.Name)
	}

	return res
}

func (u *UserManager) Load() error {
	_, err := ioutil.ReadDir("./" + DIR)

	if err != nil {
		if err = os.Mkdir("./"+DIR, 0777); err != nil {
			panic(err)
		}
		if err = os.Mkdir("./"+DIR+"/"+USER_PUB_DIR, 0777); err != nil {
			panic(err)
		}

		if err = os.Mkdir("./"+DIR+"/"+MESSAGE_DIR, 0777); err != nil {
			panic(err)
		}
	}

	dirPub, err := ioutil.ReadDir("./" + DIR + "/" + USER_PUB_DIR)
	if err != nil {
		err := os.Mkdir("./"+DIR+"/"+USER_PUB_DIR, 0777)
		if err != nil {
			return err
		}
	}

	for _, file := range dirPub {
		name := strings.Split(file.Name(), ".")[0]
		bFile, err := ioutil.ReadFile("./" + DIR + "/" + USER_PUB_DIR + "/" + file.Name())
		if err != nil {
			panic(err)
		}

		user := NewUser(name, string(bFile))

		user.loadMessage("./" + DIR + "/" + MESSAGE_DIR + "/" + name + ".message")

		u.UserList[name] = &user
	}

	return nil
}

func (u *UserManager) Save() error {
	var err error
	for _, user := range u.UserList {
		err = user.Save()
	}
	return err
}
