package main

import (
	"errors"
	"fmt"
	"go-crypt-message/internal/core/users"
	"go-crypt-message/internal/ui"
	manager "go-crypt-message/pkg/rsa"
	"os"
	"os/signal"
	"syscall"
)

func main() {

	var managers *manager.ManagerRSA
	if _, err := os.Stat("./etc"); errors.Is(err, os.ErrNotExist) {
		if err = os.Mkdir("./etc", os.ModePerm); err != nil {
			panic(err)
		}
		managers = manager.NewManagerRSA(nil)
	} else {
		if _, err := os.Stat("./etc/private"); errors.Is(err, os.ErrNotExist) {
			fmt.Println(err)
			managers = manager.NewManagerRSA(nil)
		} else {
			managers = manager.NewManagerRSA(manager.OptionString("./etc/private"))
		}
	}

	userManager := users.NewUserManager()
	err := userManager.Load()
	if err != nil {
		panic(err)
	}

	sigChan := make(chan os.Signal, 1)
	signal.Notify(sigChan,
		syscall.SIGHUP,
		syscall.SIGINT,
		syscall.SIGTERM,
		syscall.SIGQUIT)
	go func() {
		<-sigChan
		err := userManager.Save()
		fmt.Println(err)
	}()

	defer func() {
		e := userManager.Save()
		if e != nil {
			panic(e)
		}
	}()

	ui.RunUI(managers, userManager)
}
