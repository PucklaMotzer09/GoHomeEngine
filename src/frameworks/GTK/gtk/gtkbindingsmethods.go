package gtk

/*
	#cgo pkg-config: gtk+-3.0
	#include "includes.h"
*/
import "C"
import (
	"errors"
	"unsafe"
)

func (this Window) ConfigureParametersAdv(width, height uint32, title string) {
	ctitle := C.CString(title)
	defer C.free(unsafe.Pointer(ctitle))

	C.configureWindowParameters(this.Handle, C.uint(width), C.uint(height), ctitle)
}

func (this Window) ConfigureParameters() {
	this.ConfigureParametersAdv(0, 0, "")
}

func (this Window) ConnectSignals() {
	C.connectWindowSignals(this.Handle)
}

func (this Window) ToWidget() Widget {
	return Widget{C.windowToWidget(this.Handle)}
}

func (this Container) Add(widget Widget) {
	C.gtk_container_add(this.Handle, widget.Handle)
	C.gtk_widget_show(widget.Handle)
}

func (this Container) Remove(widget Widget) {
	C.gtk_container_remove(this.Handle, widget.Handle)
}

func (this Container) GetChildren() GList {
	return GList{C.gtk_container_get_children(this.Handle)}
}

func (this Box) ToContainer() Container {
	return Container{C.boxToContainer(this.Handle)}
}

func (this Button) ToContainer() Container {
	return Container{C.buttonToContainer(this.Handle)}
}

func (this Box) ToWidget() Widget {
	return Widget{C.boxToWidget(this.Handle)}
}

func (this GLArea) ToWidget() Widget {
	return Widget{C.glareaToWidget(this.Handle)}
}

func (this GLArea) Configure() {
	C.configureGLArea(this.Handle)
}

func (this Button) ToWidget() Widget {
	return Widget{C.buttonToWidget(this.Handle)}
}

func (this GObject) ToWidget() Widget {
	return Widget{C.gobjectToWidget(this.Handle)}
}

func (this GObject) ToGLArea() GLArea {
	return GLArea{C.gobjectToGLArea(this.Handle)}
}

func (this Widget) ToBox() Box {
	return Box{C.widgetToBox(this.Handle)}
}

func (this Window) ToContainer() Container {
	return Container{C.windowToContainer(this.Handle)}
}

func (this Widget) ToWindow() Window {
	return Window{C.widgetToWindow(this.Handle)}
}

func (this Widget) ToContainer() Container {
	return Container{C.widgetToContainer(this.Handle)}
}

func (this Widget) ToGrid() Grid {
	return Grid{C.widgetToGrid(this.Handle)}
}

func (this Widget) ShowAll() {
	C.gtk_widget_show_all(this.Handle)
}

func (this Widget) Show() {
	C.gtk_widget_show(this.Handle)
}

func (this Widget) SetSizeRequest(width, height int) {
	C.gtk_widget_set_size_request(this.Handle, C.gint(width), C.gint(height))
}

func (this Widget) GetSizeRequest() (int32, int32) {
	var width, height C.gint

	C.gtk_widget_get_size_request(this.Handle, &width, &height)

	return int32(width), int32(height)
}

func (this Widget) GetSize() (int32, int32) {
	var width, height C.gint
	C.widgetGetSize(this.Handle, &width, &height)
	return int32(width), int32(height)
}

func (this Widget) GetParent() Widget {
	return Widget{C.gtk_widget_get_parent(this.Handle)}
}

func (this Widget) IsNULL() bool {
	return this.Handle == nil
}

func (this Widget) Realize() {
	C.gtk_widget_realize(this.Handle)
}

func (this Widget) GetName() string {
	name := C.gtk_widget_get_name(this.Handle)
	return C.GoString(name)
}

func (this Builder) GetObject(name string) GObject {
	cstr := C.CString(name)
	defer C.free(unsafe.Pointer(cstr))

	handle := C.gtk_builder_get_object(this.Handle, cstr)

	return GObject{handle}
}

func (this Builder) AddFromFile(file string) error {
	cfile := C.CString(file)
	defer C.free(unsafe.Pointer(cfile))

	if err := C.gtk_builder_add_from_file(this.Handle, cfile, nil); err == 0 {
		return errors.New("Error while loading file")
	}

	return nil
}

func (this GList) Next() GList {
	return GList{this.Handle.next}
}

func (this GList) Prev() GList {
	return GList{this.Handle.prev}
}

func (this GList) Data() GPointer {
	return GPointer{this.Handle.data}
}

func (this GList) Equals(other GList) bool {
	return this.Handle == other.Handle
}

func (this Widget) Equals(other Widget) bool {
	return this.Handle == other.Handle
}

func (this GPointer) ToWidget() Widget {
	return Widget{C.gpointerToWidget(this.Handle)}
}

func (this Grid) Attach(child Widget, left, top, width, height int32) {
	C.gtk_grid_attach(this.Handle, child.Handle, C.gint(left), C.gint(top), C.gint(width), C.gint(height))
}

func (this Button) SignalConnect(signal string, callback ButtonSignalCallback) {

	if buttonSignalCallbacks == nil {
		buttonSignalCallbacks = make(map[int]map[string]ButtonSignalCallback)
	}
	if buttonSignalCallbacks[this.ID] == nil {
		buttonSignalCallbacks[this.ID] = make(map[string]ButtonSignalCallback)
	}
	var alreadyConnected = false
	if _, ok := buttonSignalCallbacks[this.ID]; ok {
		if _, ok1 := buttonSignalCallbacks[this.ID][signal]; ok1 {
			alreadyConnected = true
		}
	}
	if !alreadyConnected {
		signalcs := C.CString(signal)
		C.signalConnectButton(this.Handle, signalcs, C.int(this.ID))
		C.free(unsafe.Pointer(signalcs))
	}

	buttonSignalCallbacks[this.ID][signal] = callback
}

func (this Widget) SignalConnect(signal string, callback WidgetSignalCallback) {
	if widgetSignalCallbacks == nil {
		widgetSignalCallbacks = make(map[string]map[string]WidgetSignalCallback)
	}
	name := this.GetName()
	if widgetSignalCallbacks[name] == nil {
		widgetSignalCallbacks[name] = make(map[string]WidgetSignalCallback)
	}
	var alreadyConnected = false
	if _, ok := widgetSignalCallbacks[name]; ok {
		if _, ok1 := widgetSignalCallbacks[name][signal]; ok1 {
			alreadyConnected = true
		}
	}
	if !alreadyConnected {
		signalcs := C.CString(signal)
		namecs := C.CString(name)

		if signal == "size-allocate" {
			C.sizeAllocateSignalConnectWidget(this.Handle, signalcs, namecs)
		}

		C.free(unsafe.Pointer(signalcs))
		C.free(unsafe.Pointer(namecs))
	}

	widgetSignalCallbacks[name][signal] = callback
}
