// +build development

package log

import (
	"sync/atomic"
	"unsafe"
)

func init() {
	atomic.StorePointer(&defaultLogger, unsafe.Pointer(NewDevelopmentLogger()))
}
