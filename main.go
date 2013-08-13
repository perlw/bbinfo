package main

import (
	"github.com/conformal/gotk3/gtk"
	"github.com/perlw/appindicator"
	"github.com/perlw/appindicator/gtk-extensions/gotk3"
)

func main() {
	gtk.Init(nil)

	menu, _ := gtk.MenuNew()
	menuQuit, _ := gtk.MenuItemNewWithLabel("Quit")
	menu.Append(menuQuit)

	indicator := gotk3.NewAppIndicator("test-indicator", "indicator-messages", appindicator.CategoryApplicationStatus)
	indicator.SetStatus(appindicator.StatusActive)
	indicator.SetMenu(menu)

	gtk.Main()
}
