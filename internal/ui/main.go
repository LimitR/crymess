package ui

import (
	"go-crypt-message/internal/core/users"
	manager "go-crypt-message/pkg/rsa"

	"github.com/gdamore/tcell/v2"
	"github.com/rivo/tview"

	"golang.design/x/clipboard"
)

const (
	MY_MESSAGE   = 1
	SOME_MESSAGE = 0
)

func RunUI(manager *manager.ManagerRSA, userManager *users.UserManager) {
	app := tview.NewApplication()
	box := tview.NewFlex()

	column := 0

	btnGetMyPublicKey := tview.NewButton("Save public key in clipboard")
	btnGetMyPublicKey.SetTitleAlign(1)
	lableUserName := tview.NewTextView().SetText("User is not selected").SetTextAlign(1)
	listUser := tview.NewList()
	inputText := tview.NewTextArea().SetPlaceholder("Input text...")
	tableMessage := tview.NewTable().
		SetBorders(false).SetFixed(2, 10)

	pages := tview.NewPages()

	userPress := ""

	if len(userManager.UserList) != 0 {
		for _, user := range userManager.UserList {
			n := user.Name
			listUser.AddItem(n, "", 20, func() {
				lableUserName.SetText(n)
				userPress = n
			})
		}
	}

	flexUserList := tview.NewFlex().SetDirection(tview.FlexRow)
	flexUserList.AddItem(listUser, 0, 5, false)

	nameUser := ""
	pubKey := ""
	form := tview.NewForm().
		AddInputField("Name", "", 20, nil, func(text string) {
			nameUser = text
		}).
		AddTextArea("Public key", "", 20, 0, 10000, func(text string) {
			pubKey = text
		}).
		AddButton("Save", func() {
			if nameUser != "" && pubKey != "" {
				u := users.NewUser(nameUser, pubKey)
				userManager.UserList[nameUser] = u
				go func() {
					u.Save()
				}()
				listUser.AddItem(nameUser, "", 20, func() {
					lableUserName.SetText(nameUser)
				})
				pages.SwitchToPage("main")
			}
		}).
		AddButton("Quit", func() {
			pages.SwitchToPage("main")
		})

	pages.AddPage("add-user", form, true, true)

	btnEncrypt := tview.NewButton("Encrypt")
	btnCrypto := tview.NewButton("Decrypt")
	btnAddUser := tview.NewButton("Add User")
	listUser.SetShortcutColor(tcell.Color150)

	flexUserList.AddItem(btnAddUser, 0, 1, false)

	btnAddUser.SetSelectedFunc(func() {
		pages.SwitchToPage("add-user")
	})

	btnGetMyPublicKey.SetSelectedFunc(func() {
		msg := manager.GetPublicKey()
		clipboard.Write(clipboard.FmtText, msg)
	})

	btnEncrypt.SetSelectedFunc(func() {
		msg := inputText.GetText()
		usr := userManager.UserList[userPress]
		newMsg := usr.Encrypt(msg)
		clipboard.Write(clipboard.FmtText, []byte(newMsg))

		inputText.SetText("", true)
		column += 1
		tableMessage.SetCell(column, MY_MESSAGE, tview.NewTableCell(msg))
	})
	btnCrypto.SetSelectedFunc(func() {
		msg := inputText.GetText()
		inputText.SetText("", true)
		column += 1
		newMsg := manager.Decrypt(msg)
		// usr := userManager.UserList[userPress]
		// usr.AddMessage(msg)
		tableMessage.SetCell(column, SOME_MESSAGE, tview.NewTableCell(userPress+": "+newMsg))
		tableMessage.Draw(tcell.NewSimulationScreen(""))
	})

	box.AddItem(
		tview.NewFlex().SetDirection(tview.FlexRow).
			AddItem(
				tview.NewFlex().SetDirection(tview.FlexColumn).
					AddItem(tview.NewFlex().SetDirection(tview.FlexRow).AddItem(btnGetMyPublicKey, 0, 1, false).
						AddItem(flexUserList, 0, 5, false), 0, 1, false).
					AddItem(tview.NewFlex().SetDirection(tview.FlexRow).AddItem(lableUserName, 0, 1, false).AddItem(tableMessage, 0, 1, false), 0, 1, false),
				0, 3, false).
			AddItem(tview.NewFlex().SetDirection(tview.FlexColumn).AddItem(inputText, 0, 5, true).AddItem(btnEncrypt, 0, 1, false).AddItem(btnCrypto, 0, 1, false), 0, 1, true), 0, 1, false)

	pages.AddPage("main", box, true, true)
	pages.SwitchToPage("main")

	if err := app.SetRoot(pages, true).EnableMouse(true).SetFocus(pages).Run(); err != nil {
		panic(err)
	}
}
