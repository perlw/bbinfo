package main

//#cgo pkg-config: gtk+-3.0
//#include <stdlib.h>
//#include <gtk/gtk.h>
import "C"

import "unsafe"
import "appindicator"

func main() {
	C.gtk_init(nil, nil)

	menu := C.gtk_menu_new()
	quitString := (*C.gchar)(unsafe.Pointer(C.CString("Quit")))
	defer C.free(unsafe.Pointer(quitString))
	menuQuit := C.gtk_menu_item_new_with_label(quitString)
	C.gtk_menu_attach((*C.GtkMenu)(unsafe.Pointer(menu)), menuQuit, 0, 1, 0, 1)

	indicator := appindicator.NewAppIndicator("test-indicator", "indicator-messages", appindicator.CategoryApplicationStatus)
	indicator.SetStatus(appindicator.StatusActive)
	indicator.SetMenu(unsafe.Pointer(menu))

	C.gtk_main()
}
