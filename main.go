package main

//#cgo CFLAGS: -I/usr/include/libappindicator3-0.1
//#cgo LDFLAGS: -lappindicator3
//#cgo pkg-config: gtk+-3.0
//#include <stdlib.h>
//#include <gtk/gtk.h>
//#include <libappindicator/app-indicator.h>
import "C"

import "unsafe"

func main() {
	C.gtk_init(nil, nil)

	menu := C.gtk_menu_new()

	quitString := (*C.gchar)(unsafe.Pointer(C.CString("Quit")))
	defer C.free(unsafe.Pointer(quitString))
	menuQuit := C.gtk_menu_item_new_with_label(quitString)

	C.gtk_menu_attach((*C.GtkMenu)(unsafe.Pointer(menu)), menuQuit, 0, 1, 0, 1)

	titleString := (*C.gchar)(unsafe.Pointer(C.CString("test-indicator")))
	defer C.free(unsafe.Pointer(titleString))
	iconString := (*C.gchar)(unsafe.Pointer(C.CString("indicator-messages")))
	defer C.free(unsafe.Pointer(iconString))

	indicator := C.app_indicator_new(titleString, iconString, C.APP_INDICATOR_CATEGORY_APPLICATION_STATUS)
	C.app_indicator_set_status(indicator, C.APP_INDICATOR_STATUS_ACTIVE)
	C.app_indicator_set_menu(indicator, (*C.GtkMenu)(unsafe.Pointer(menu)))

	C.gtk_main()
}
