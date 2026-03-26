//go:build darwin

package notifier

/*
#cgo LDFLAGS: -framework Foundation -framework UserNotifications
#include <stdlib.h>

char *focusbar_notify(const char *title, const char *message);
*/
import "C"

import (
	"errors"
	"unsafe"

	"focusbar/internal/logx"
)

type Native struct{}

func New() Sender {
	return &Native{}
}

func (n *Native) Notify(title, message string) error {
	logx.Infof("notifier invoking native bridge title=%q message=%q", title, message)

	cTitle := C.CString(title)
	cMessage := C.CString(message)
	defer C.free(unsafe.Pointer(cTitle))
	defer C.free(unsafe.Pointer(cMessage))

	errPtr := C.focusbar_notify(cTitle, cMessage)
	if errPtr == nil {
		logx.Infof("notifier native bridge completed successfully")
		return nil
	}
	defer C.free(unsafe.Pointer(errPtr))

	errMsg := C.GoString(errPtr)
	logx.Errorf("notifier native bridge failed: %s", errMsg)
	return errors.New(errMsg)
}
