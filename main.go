package main

import (
	"github.com/conformal/gotk3/gtk"
	"github.com/perlw/appindicator"
	"unsafe"
)

func main() {
	gtk.Init(nil)

	menu, _ := gtk.MenuNew()
	menuQuit, _ := gtk.MenuItemNewWithLabel("Quit")
	menu.Append(menuQuit)

	indicator := appindicator.NewAppIndicator("test-indicator", "indicator-messages", appindicator.CategoryApplicationStatus)
	indicator.SetStatus(appindicator.StatusActive)
	indicator.C_SetMenu(unsafe.Pointer(menu.Native()))

	gtk.Main()
}
