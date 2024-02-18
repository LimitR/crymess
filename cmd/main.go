package main

import (
	"go-crypt-message/internal/core/users"
	"go-crypt-message/internal/ui"
	manager "go-crypt-message/pkg/rsa"
)

func main() {

	// manager := manager.NewManagerRSA(nil)
	manager := manager.NewManagerRSA(manager.OptionString("./etc/private"))

	userManager := users.NewUserManager()
	err := userManager.Load()
	if err != nil {
		panic(err)
	}

	ui.RunUI(manager, userManager)
}
