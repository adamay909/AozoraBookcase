package main

import (
	"syscall/js"

	ac "github.com/adamay909/AozoraConvert/v2"
)

func settingsMenu(event js.Value, param ...any) {

	coverScreen(30)

	showSettingsMenu()

	activateMenu()

}

func handleSettings(event js.Value, param ...any) {

	kidsB, _ := getElementByID("kidslib")

	kidsNew := kidsB.Get("checked").Bool()

	if kidsNew != globalSettings.kids {

		globalSettings.kids = kidsNew

		globalLib = nil

		go reloadSvc()

		return

	}

	c, _ := getElementByID("jis0213")

	jis0213 := c.Get("checked").Bool()

	if jis0213 {

		ac.SetJIS0213()

		globalSettings.jis0213 = true

	} else {

		ac.SetFullUnicode()

		globalSettings.jis0213 = false
	}

	removeMenu()

	return
}

func removeMenu() {

	mn, _ := getElementByID("x-menu")

	mn.Call("remove")

	uncoverScreen()

}

func reloadSvc() {

	replaceBody(string(readFromResources("loading.html")))

	coverAndWait(domBody, 10)

	initLibrary()

	loadMainPage()

}

func showSettingsMenu() {

	elem := createElement("div", string(readFromResources("menu.html")))

	domBody.Call("append", elem)

	if globalSettings.kids {
		elem, _ = getElementByID("kidslib")
		elem.Call("setAttribute", "checked", true)
	}

	if globalSettings.jis0213 {
		elem, _ = getElementByID("jis0213")
		elem.Call("setAttribute", "checked", true)
	}

}

func activateMenu() {

	submitb, _ := getElementByID("submit")

	addEventListener(submitb, "click", handleSettings)

}
