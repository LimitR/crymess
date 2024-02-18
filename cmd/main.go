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
	if _, err := os.Stat("./etc/private"); errors.Is(err, os.ErrNotExist) {
		managers = manager.NewManagerRSA(nil)
	} else {
		managers = manager.NewManagerRSA(manager.OptionString("./etc/private"))
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

	defer userManager.Save()

	ui.RunUI(managers, userManager)
}
