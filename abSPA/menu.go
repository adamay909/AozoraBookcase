package main

import (
	"log"
	"syscall/js"
)

func settingsMenu(event js.Value, param ...any) {

	//	ma, _ := getElementById("mainarea")

	coverScreen(30)

	showSettingsMenu()

	activateMenu()

}

func handleSettings(event js.Value, param ...any) {

	kidsB, _ := getElementById("kidslib")

	kidsNew := kidsB.Get("checked").Bool()

	log.Println("switch to kids library", kidsNew)

	if kidsNew != globalSettings.kids {

		log.Println("reloading library")

		globalSettings.kids = kidsNew

		globalLib = nil

		go reloadSvc()

		return

	}

	removeMenu()

	return
}

func removeMenu() {

	mn, _ := getElementById("x-menu")

	mn.Call("remove")

	uncoverScreen()

}

func reloadSvc() {

	f, _ := templateFiles.Open("resources/loading.html")

	replaceBody(string(readFrom(f)))

	coverScreen(0)

	initLibrary()

	loadMainPage()

}

func showSettingsMenu() {

	log.Println("showing settings")

	f, _ := templateFiles.Open("resources/menu.html")

	elem := createElement("div", string(readFrom(f)))

	domBody.Call("append", elem)

	if globalSettings.kids {
		elem, _ = getElementById("kidslib")
		elem.Call("setAttribute", "checked", true)
	}

	enableClick(elem)
}

func activateMenu() {

	submitb, _ := getElementById("submit")

	addEventListener(submitb, "click", handleSettings)

}
