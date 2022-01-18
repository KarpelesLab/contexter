package contexter

import (
	"context"
	"reflect"
	"regexp"
	"runtime/debug"
	"strconv"
	"unsafe"
)

var (
	matchHex    = regexp.MustCompile("[^+]0x([a-fA-F0-9]{2,16}), 0x([a-fA-F0-9]{2,16})")
	intfType    = reflect.TypeOf((*context.Context)(nil)).Elem()
	intfTypeVal = intfValue(intfType)
)

type itab struct {
	inter uintptr
	typ   uintptr
}

type interfaceHeader struct {
	tab  *itab
	data unsafe.Pointer
}

func intfValue(v interface{}) uintptr {
	x := (*interfaceHeader)(unsafe.Pointer(&v))
	return uintptr(x.data)
}

func testPtr(p1, p2 uintptr) context.Context {
	defer func() {
		recover()
	}()

	// we expect itab in p1
	tab := (*itab)(unsafe.Pointer(p1))

	if tab == nil || tab.inter != intfTypeVal {
		// not a context.Context interface
		return nil
	}

	var res context.Context
	ti := (*interfaceHeader)(unsafe.Pointer(&res))
	ti.data = unsafe.Pointer(p2)
	ti.tab = tab

	return res
}

func testTab(tab, expect uintptr) bool {
	defer func() {
		recover()
	}()

	t := (*itab)(unsafe.Pointer(tab))
	return t.inter == expect
}

// Context will return the first instance of context.Context found in the
// calling stack
func Context() context.Context {
	// argp := unsafe.Pointer(frame.argp)
	// The "instruction" of argument printing is encoded in _FUNCDATA_ArgInfo.
	// See cmd/compile/internal/ssagen.emitArgInfo for the description of the
	// encoding
	// _FUNCDATA_ArgInfo            = 5
	// ...
	// Unfortunately(?) go provides no way to get this kind of data from the
	// stack except by calling runtime.Stack(), in which case pointers are
	// encoded in hex as string. This is not optimal, but at least it works.

	v := debug.Stack()
	res := matchHex.FindAllSubmatch(v, -1)
	for _, h := range res {
		i1, _ := strconv.ParseUint(string(h[1]), 16, 64)
		i2, _ := strconv.ParseUint(string(h[2]), 16, 64)

		if i1 == 0 || i2 == 0 {
			continue
		}

		res := testPtr(uintptr(i1), uintptr(i2))
		if res != nil {
			return res
		}
	}

	return nil
}

func Find(i interface{}) bool {
	// try to find object of type intf
	v := reflect.ValueOf(i).Elem()
	searchType := intfValue(v.Type())

	// when an interface is called, it passes two pointers
	res := matchHex.FindAllSubmatch(debug.Stack(), -1)
	for _, h := range res {
		i1, _ := strconv.ParseUint(string(h[1]), 16, 64)
		i2, _ := strconv.ParseUint(string(h[2]), 16, 64)

		if i1 == 0 || i2 == 0 {
			continue
		}

		if !testTab(uintptr(i1), searchType) {
			continue
		}
		// found our value
		ti := (*interfaceHeader)(unsafe.Pointer(v.UnsafeAddr()))
		ti.tab = (*itab)(unsafe.Pointer(uintptr(i1)))
		ti.data = unsafe.Pointer(uintptr(i2))
		return true
	}
	return false
}
